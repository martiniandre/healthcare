package imaging

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/healthcare/backend/internal/modules/audit_logs"
	"github.com/healthcare/backend/internal/shared/apperrors"
	"github.com/healthcare/backend/internal/shared/ctxkeys"
	"github.com/healthcare/backend/internal/shared/dicom"
	"github.com/healthcare/backend/internal/shared/storage"
	"github.com/healthcare/backend/internal/shared/validator"
	"github.com/redis/go-redis/v9"
)

const MaxDICOMUploadBytes int64 = 2 << 30

var unsafeObjectPathChars = regexp.MustCompile(`[^A-Za-z0-9-.]`)

func sanitizeLeafName(rawName string) string {
	baseName := strings.TrimSpace(filepath.Base(rawName))
	if baseName == "" || baseName == "." || baseName == string(filepath.Separator) {
		return ""
	}
	return unsafeObjectPathChars.ReplaceAllString(baseName, "_")
}

type Service interface {
	UploadDICOMStream(ctx context.Context, input UploadDICOMInput, streamReader io.Reader) (*UploadDICOMResult, error)
	GetImagingStudy(ctx context.Context, studyID string) (*ImagingStudy, error)
	ListImagingStudies(ctx context.Context, patientFhirID string) ([]*ImagingStudy, error)
	GetDownloadURL(ctx context.Context, studyID string) (string, time.Time, error)
}

type service struct {
	dbRepository  Repository
	storageClient storage.StorageClient
	redisClient   *redis.Client
	auditService  audit_logs.Service

	patientIdentityReader PatientIdentityReader
}

func NewService(dbRepository Repository, storageClient storage.StorageClient, redisClient *redis.Client, auditService audit_logs.Service, patientIdentityReader PatientIdentityReader) Service {
	return &service{
		dbRepository:          dbRepository,
		storageClient:         storageClient,
		redisClient:           redisClient,
		auditService:          auditService,
		patientIdentityReader: patientIdentityReader,
	}
}

func (serviceInstance *service) UploadDICOMStream(ctx context.Context, input UploadDICOMInput, streamReader io.Reader) (*UploadDICOMResult, error) {
	if fieldViolations := validateUploadDICOMInput(input); len(fieldViolations) > 0 {
		return nil, apperrors.InvalidArgument("invalid dicom upload input", fieldViolations)
	}

	limitedReader := &io.LimitedReader{
		R: streamReader,
		N: MaxDICOMUploadBytes + 1,
	}

	headerBytes := make([]byte, 132)
	bytesRead, readError := io.ReadFull(limitedReader, headerBytes)
	if readError != nil && !errors.Is(readError, io.EOF) && !errors.Is(readError, io.ErrUnexpectedEOF) {
		return nil, fmt.Errorf("failed to read dicom header: %w", readError)
	}

	if bytesRead < 132 {
		return nil, fmt.Errorf("%w: preamble is too small", apperrors.ErrInvalidDICOM)
	}

	magicBytesSignature := string(headerBytes[128:132])
	if magicBytesSignature != "DICM" {
		return nil, fmt.Errorf("%w: magic bytes DICM signature missing", apperrors.ErrInvalidDICOM)
	}

	inspection, inspectionError := inspectDICOMHeader(headerBytes[:bytesRead], limitedReader)
	if inspectionError != nil {
		return nil, inspectionError
	}

	validated, validationError := validateDICOMMetadata(inspection.metadata)
	if validationError != nil {
		return nil, validationError
	}

	resolvedModality, modalityWarnings, modalityError := resolveModality(validated.modality, input.Modality)
	if modalityError != nil {
		return nil, modalityError
	}
	validated.modality = resolvedModality
	validated.instanceIndex.Modality = resolvedModality

	identityWarnings, identityError := serviceInstance.verifyPatientIdentity(ctx, input.PatientFhirID, validated.source)
	if identityError != nil {
		return nil, identityError
	}

	if duplicateError := serviceInstance.rejectDuplicateStudy(ctx, validated.studyInstanceUID); duplicateError != nil {
		return nil, duplicateError
	}

	reconstructedReader := io.MultiReader(bytes.NewReader(inspection.consumedBytes), limitedReader)
	boundedReader := &uploadLimitReader{reader: reconstructedReader, remainingBytes: MaxDICOMUploadBytes}

	studyID := uuid.New()
	sanitizedPatientFhirID := unsafeObjectPathChars.ReplaceAllString(input.PatientFhirID, "_")
	sanitizedFileName := sanitizeLeafName(input.FileName)
	studyExtension := filepath.Ext(sanitizedFileName)
	if studyExtension == "" {
		studyExtension = ".dcm"
	}
	studyFileName := studyID.String() + studyExtension
	key := fmt.Sprintf("pacs/%s/%s", sanitizedPatientFhirID, studyFileName)

	uploadError := serviceInstance.storageClient.Upload(ctx, key, boundedReader, "application/dicom")
	if uploadError != nil {
		return nil, fmt.Errorf("failed to upload dicom to storage: %w", uploadError)
	}

	if sanitizedFileName == "" {
		sanitizedFileName = studyFileName
	}

	studyTitle := truncateRunes(strings.TrimSpace(input.Title), maximumTitleLength)
	if studyTitle == "" {
		studyTitle = buildStudyTitle(validated.modality, validated.studyDescription, validated.studyDateTime)
	}

	now := time.Now()
	study := &ImagingStudy{
		ID:               studyID,
		PatientFhirID:    input.PatientFhirID,
		Title:            studyTitle,
		Modality:         validated.modality,
		FileName:         sanitizedFileName,
		StudyInstanceUID: validated.studyInstanceUID,
		StudyDateTime:    validated.studyDateTime,
		AccessionNumber:  validated.accessionNumber,
		StudyDescription: validated.studyDescription,
		Status:           ImagingStudyStatusReceiving,
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	dbError := serviceInstance.dbRepository.CreateImagingStudy(ctx, study, validated.instanceIndex)
	if dbError != nil {
		serviceInstance.discardStoredObject(ctx, key)
		if errors.Is(dbError, apperrors.ErrDICOMDuplicateInstance) {
			return nil, dbError
		}
		return nil, fmt.Errorf("failed to persist imaging study operational record: %w", dbError)
	}

	if serviceInstance.redisClient != nil {
		enqueueError := serviceInstance.redisClient.LPush(ctx, "dicom_processing_queue", studyID.String()).Err()
		if enqueueError != nil {
			return nil, fmt.Errorf("failed to enqueue dicom processing job: %w", enqueueError)
		}
	}

	return &UploadDICOMResult{
		Study:             study,
		SeriesInstanceUID: validated.seriesInstanceUID,
		SOPInstanceUID:    validated.sopInstanceUID,
		Warnings:          append(modalityWarnings, identityWarnings...),
	}, nil
}

func resolveModality(fileModality string, userModality string) (string, []UploadWarning, error) {
	typedModality := strings.ToUpper(strings.TrimSpace(userModality))
	if fileModality == "" {
		if typedModality == "" {
			return "", nil, apperrors.WithDetail(apperrors.ErrInvalidDICOM, "Modality is missing")
		}
		return typedModality, nil, nil
	}
	if typedModality != "" && typedModality != fileModality {
		return fileModality, []UploadWarning{newUploadWarning(WarningCodeModalityOverridden)}, nil
	}
	return fileModality, nil, nil
}

func (serviceInstance *service) verifyPatientIdentity(ctx context.Context, patientFhirID string, source *dicom.InstanceMetadata) ([]UploadWarning, error) {
	var patientIdentity PatientIdentity
	if serviceInstance.patientIdentityReader != nil {
		readIdentity, readError := serviceInstance.patientIdentityReader.ReadPatientIdentity(ctx, patientFhirID)
		if readError != nil {
			return nil, readError
		}
		patientIdentity = readIdentity
	}

	return evaluatePatientMatch(patientFhirID, dicomPatientIdentity{
		patientID:   source.PatientID,
		patientName: source.PatientName,
	}, patientIdentity)
}

func (serviceInstance *service) rejectDuplicateStudy(ctx context.Context, studyInstanceUID string) error {
	_, findError := serviceInstance.dbRepository.FindImagingStudyByStudyInstanceUID(ctx, studyInstanceUID)
	if findError == nil {
		return apperrors.ErrDICOMDuplicateInstance
	}
	if errors.Is(findError, apperrors.ErrImagingStudyNotFound) {
		return nil
	}
	return fmt.Errorf("failed to check duplicate imaging study: %w", findError)
}

func (serviceInstance *service) discardStoredObject(ctx context.Context, key string) {
	if deleteError := serviceInstance.storageClient.Delete(context.WithoutCancel(ctx), key); deleteError != nil {
		slog.Warn("failed to discard orphaned dicom object", "error", deleteError)
	}
}

type uploadLimitReader struct {
	reader         io.Reader
	remainingBytes int64
}

func (reader *uploadLimitReader) Read(buffer []byte) (int, error) {
	if reader.remainingBytes <= 0 {
		return 0, apperrors.ErrPayloadTooLarge
	}
	if int64(len(buffer)) > reader.remainingBytes {
		buffer = buffer[:int(reader.remainingBytes)]
	}
	bytesRead, readError := reader.reader.Read(buffer)
	reader.remainingBytes -= int64(bytesRead)
	if readError == nil && reader.remainingBytes == 0 {
		return bytesRead, apperrors.ErrPayloadTooLarge
	}
	return bytesRead, readError
}

func (serviceInstance *service) GetImagingStudy(ctx context.Context, studyID string) (*ImagingStudy, error) {
	parsedStudyID, err := uuid.Parse(studyID)
	if err != nil {
		return nil, apperrors.InvalidArgument("invalid get imaging study input", map[string]string{"study_id": "must be a valid UUID"})
	}

	return serviceInstance.dbRepository.GetImagingStudy(ctx, parsedStudyID)
}

func (serviceInstance *service) ListImagingStudies(ctx context.Context, patientFhirID string) ([]*ImagingStudy, error) {
	if strings.TrimSpace(patientFhirID) == "" {
		return nil, apperrors.InvalidArgument("invalid list imaging studies input", map[string]string{"patient_fhir_id": "is required"})
	}

	return serviceInstance.dbRepository.ListImagingStudiesByPatient(ctx, patientFhirID)
}

func (serviceInstance *service) GetDownloadURL(ctx context.Context, studyID string) (string, time.Time, error) {
	parsedStudyID, err := uuid.Parse(studyID)
	if err != nil {
		return "", time.Time{}, apperrors.InvalidArgument("invalid download url input", map[string]string{"study_id": "must be a valid UUID"})
	}

	study, dbError := serviceInstance.dbRepository.GetImagingStudy(ctx, parsedStudyID)
	if dbError != nil {
		return "", time.Time{}, dbError
	}

	expirationDuration := 15 * time.Minute
	studyKey := serviceInstance.buildStorageKey(study)
	downloadURL, presignError := serviceInstance.storageClient.GetPresignedURL(ctx, studyKey, expirationDuration)
	if presignError != nil {
		return "", time.Time{}, fmt.Errorf("failed to generate presigned url: %w", presignError)
	}

	serviceInstance.auditDownloadURLIssuance(ctx, parsedStudyID)

	expiresAt := time.Now().Add(expirationDuration)
	return downloadURL, expiresAt, nil
}

func (serviceInstance *service) buildStorageKey(study *ImagingStudy) string {
	studyExtension := filepath.Ext(study.FileName)
	if studyExtension == "" {
		studyExtension = ".dcm"
	}
	sanitizedPatientFhirID := unsafeObjectPathChars.ReplaceAllString(study.PatientFhirID, "_")
	return fmt.Sprintf("pacs/%s/%s%s", sanitizedPatientFhirID, study.ID.String(), studyExtension)
}

func (serviceInstance *service) auditDownloadURLIssuance(ctx context.Context, studyID uuid.UUID) {
	if serviceInstance.auditService == nil {
		return
	}

	callerUserID, _ := ctx.Value(ctxkeys.UserIDKey).(string)
	callerRole, _ := ctx.Value(ctxkeys.RoleKey).(string)
	correlationID, _ := ctx.Value(ctxkeys.RequestIDKey).(string)

	go func() {
		_, auditError := serviceInstance.auditService.CreateResourceAuditLog(context.Background(), audit_logs.ResourceAuditLog{
			CorrelationID: correlationID,
			CallerUserID:  callerUserID,
			CallerRole:    callerRole,
			Method:        "imaging.GetDICOMDownloadURL",
			AccessGranted: true,
			ResourceType:  "imaging_study",
			ResourceID:    studyID.String(),
			Action:        "download_url_issued",
		})
		if auditError != nil {
			fmt.Printf("failed to persist imaging download audit log: %v\n", auditError)
		}
	}()
}

func validateUploadDICOMInput(input UploadDICOMInput) map[string]string {
	fieldViolations := make(map[string]string)
	if strings.TrimSpace(input.PatientFhirID) == "" {
		fieldViolations["patient_fhir_id"] = "is required"
	}
	if strings.TrimSpace(input.Modality) != "" && !validator.IsValidDICOMModality(input.Modality) {
		fieldViolations["modality"] = "invalid dicom modality"
	}
	return fieldViolations
}
