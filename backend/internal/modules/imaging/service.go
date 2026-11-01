package imaging

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
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

const (
	uploadLeadingBytesLength = 4
	uploadReadBufferBytes    = 64 << 10
	downloadURLExpiration    = 15 * time.Minute
	processingQueueKey       = "dicom_processing_queue"
)

type Service interface {
	UploadDICOMStream(ctx context.Context, input UploadDICOMInput, streamReader io.Reader) (*UploadDICOMResult, error)
	GetImagingStudy(ctx context.Context, studyID string) (*ImagingStudy, error)
	ListImagingStudies(ctx context.Context, patientFhirID string) ([]*ImagingStudy, error)
	GetDownloadURL(ctx context.Context, studyID string) (string, time.Time, error)
	ListStudyInstanceDownloads(ctx context.Context, studyID string) ([]StudyInstanceDownload, error)
}

type StudyInstanceDownload struct {
	Instance    *InstanceIndex
	DownloadURL string
	ExpiresAt   time.Time
}

type ServiceSettings struct {
	ArchiveLimits  ArchiveLimits
	SpoolDirectory string
}

func DefaultServiceSettings() ServiceSettings {
	return ServiceSettings{ArchiveLimits: DefaultArchiveLimits()}
}

type service struct {
	dbRepository  Repository
	storageClient storage.StorageClient
	redisClient   *redis.Client
	auditService  audit_logs.Service
	settings      ServiceSettings

	patientIdentityReader PatientIdentityReader
}

func NewService(dbRepository Repository, storageClient storage.StorageClient, redisClient *redis.Client, auditService audit_logs.Service, patientIdentityReader PatientIdentityReader) Service {
	return NewServiceWithSettings(dbRepository, storageClient, redisClient, auditService, patientIdentityReader, DefaultServiceSettings())
}

func NewServiceWithSettings(dbRepository Repository, storageClient storage.StorageClient, redisClient *redis.Client, auditService audit_logs.Service, patientIdentityReader PatientIdentityReader, settings ServiceSettings) Service {
	return &service{
		dbRepository:          dbRepository,
		storageClient:         storageClient,
		redisClient:           redisClient,
		auditService:          auditService,
		settings:              settings,
		patientIdentityReader: patientIdentityReader,
	}
}

func (serviceInstance *service) UploadDICOMStream(ctx context.Context, input UploadDICOMInput, streamReader io.Reader) (*UploadDICOMResult, error) {
	if fieldViolations := validateUploadDICOMInput(input); len(fieldViolations) > 0 {
		return nil, apperrors.InvalidArgument("invalid dicom upload input", fieldViolations)
	}

	bufferedReader := bufio.NewReaderSize(streamReader, uploadReadBufferBytes)
	leadingBytes, _ := bufferedReader.Peek(uploadLeadingBytesLength)
	if isArchiveUpload(leadingBytes, input.FileName) {
		return serviceInstance.uploadArchive(ctx, input, bufferedReader)
	}
	return serviceInstance.uploadSingleInstance(ctx, input, bufferedReader)
}

func (serviceInstance *service) uploadSingleInstance(ctx context.Context, input UploadDICOMInput, source io.Reader) (*UploadDICOMResult, error) {
	limitedReader := &io.LimitedReader{
		R: source,
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

	prepared, prepareError := prepareInstance(inspection.metadata, input.Modality)
	if prepareError != nil {
		return nil, singleFileUploadError(prepareError)
	}
	validated := prepared.validated

	identity, identityLoadError := serviceInstance.loadPatientIdentity(ctx, input.PatientFhirID)
	if identityLoadError != nil {
		return nil, identityLoadError
	}
	identityWarnings, identityError := evaluateInstancePatient(input.PatientFhirID, identity, validated.source)
	if identityError != nil {
		return nil, identityError
	}

	existingStudy, resolveError := serviceInstance.resolveExistingStudy(ctx, input.PatientFhirID, validated.studyInstanceUID)
	if resolveError != nil {
		return nil, resolveError
	}
	if existingStudy != nil {
		alreadyStored, lookupError := serviceInstance.dbRepository.HasInstance(ctx, existingStudy.ID, validated.sopInstanceUID)
		if lookupError != nil {
			return nil, fmt.Errorf("failed to check duplicate imaging instance: %w", lookupError)
		}
		if alreadyStored {
			return nil, apperrors.ErrDICOMDuplicateInstance
		}
	}

	studyID := uuid.New()
	if existingStudy != nil {
		studyID = existingStudy.ID
	}
	instance := validated.instanceIndex
	instance.ID = uuid.New()
	instance.StudyID = studyID
	instance.FileName = singleFileDisplayName(input.FileName, instance.ID)
	instance.ObjectKey = instanceObjectKey(input.PatientFhirID, studyID, instance.ID)

	boundedReader := &uploadLimitReader{
		reader:         io.MultiReader(bytes.NewReader(inspection.consumedBytes), limitedReader),
		remainingBytes: MaxDICOMUploadBytes,
	}
	uploadError := serviceInstance.storageClient.Upload(ctx, instance.ObjectKey, boundedReader, "application/dicom")
	if uploadError != nil {
		return nil, fmt.Errorf("failed to upload dicom to storage: %w", uploadError)
	}

	study, persistError := serviceInstance.persistInstances(ctx, input, existingStudy, studyID, validated, instance.FileName, []*InstanceIndex{instance})
	if persistError != nil {
		serviceInstance.discardStoredObjects(ctx, []string{instance.ObjectKey})
		return nil, persistError
	}

	if enqueueError := serviceInstance.enqueueStudyProcessing(ctx, study.ID); enqueueError != nil {
		return nil, enqueueError
	}

	return &UploadDICOMResult{
		Study:             study,
		StudyCreated:      existingStudy == nil,
		SeriesInstanceUID: validated.seriesInstanceUID,
		SOPInstanceUID:    validated.sopInstanceUID,
		Warnings:          appendUniqueWarnings(prepared.warnings, identityWarnings...),
		InstancesAccepted: 1,
		InstancesRejected: []RejectedInstance{},
	}, nil
}

func (serviceInstance *service) persistInstances(ctx context.Context, input UploadDICOMInput, existingStudy *ImagingStudy, studyID uuid.UUID, validated *validatedDICOMMetadata, fileName string, instances []*InstanceIndex) (*ImagingStudy, error) {
	now := time.Now()
	if existingStudy != nil {
		appendError := serviceInstance.dbRepository.AppendInstances(ctx, existingStudy.ID, instances, now)
		if appendError != nil {
			if errors.Is(appendError, apperrors.ErrDICOMDuplicateInstance) {
				return nil, appendError
			}
			return nil, fmt.Errorf("failed to append imaging instances: %w", appendError)
		}
		existingStudy.Status = ImagingStudyStatusReceiving
		existingStudy.UpdatedAt = now
		return existingStudy, nil
	}

	study := newImagingStudy(input, studyID, validated, fileName, now)
	createError := serviceInstance.dbRepository.CreateImagingStudy(ctx, study, instances)
	if createError != nil {
		if errors.Is(createError, apperrors.ErrDICOMDuplicateInstance) {
			return nil, createError
		}
		return nil, fmt.Errorf("failed to persist imaging study operational record: %w", createError)
	}
	return study, nil
}

func newImagingStudy(input UploadDICOMInput, studyID uuid.UUID, validated *validatedDICOMMetadata, fileName string, now time.Time) *ImagingStudy {
	studyTitle := truncateRunes(strings.TrimSpace(input.Title), maximumTitleLength)
	if studyTitle == "" {
		studyTitle = buildStudyTitle(validated.modality, validated.studyDescription, validated.studyDateTime)
	}
	return &ImagingStudy{
		ID:               studyID,
		PatientFhirID:    input.PatientFhirID,
		Title:            studyTitle,
		Modality:         validated.modality,
		FileName:         fileName,
		StudyInstanceUID: validated.studyInstanceUID,
		StudyDateTime:    validated.studyDateTime,
		AccessionNumber:  validated.accessionNumber,
		StudyDescription: validated.studyDescription,
		Status:           ImagingStudyStatusReceiving,
		CreatedAt:        now,
		UpdatedAt:        now,
	}
}

func singleFileDisplayName(rawFileName string, instanceID uuid.UUID) string {
	if sanitizedName := sanitizeLeafName(rawFileName); sanitizedName != "" {
		return truncateRunes(sanitizedName, maximumDisplayNameRunes)
	}
	return instanceID.String() + defaultDICOMExtension
}

func (serviceInstance *service) loadPatientIdentity(ctx context.Context, patientFhirID string) (PatientIdentity, error) {
	if serviceInstance.patientIdentityReader == nil {
		return PatientIdentity{}, nil
	}
	return serviceInstance.patientIdentityReader.ReadPatientIdentity(ctx, patientFhirID)
}

func evaluateInstancePatient(patientFhirID string, identity PatientIdentity, source *dicom.InstanceMetadata) ([]UploadWarning, error) {
	return evaluatePatientMatch(patientFhirID, dicomPatientIdentity{
		patientID:   source.PatientID,
		patientName: source.PatientName,
	}, identity)
}

func (serviceInstance *service) resolveExistingStudy(ctx context.Context, patientFhirID string, studyInstanceUID string) (*ImagingStudy, error) {
	existingStudy, findError := serviceInstance.dbRepository.FindImagingStudyByStudyInstanceUID(ctx, studyInstanceUID)
	if findError != nil {
		if errors.Is(findError, apperrors.ErrImagingStudyNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to check existing imaging study: %w", findError)
	}
	if existingStudy.PatientFhirID != patientFhirID {
		return nil, apperrors.ErrDICOMPatientMismatch
	}
	return existingStudy, nil
}

func (serviceInstance *service) enqueueStudyProcessing(ctx context.Context, studyID uuid.UUID) error {
	if serviceInstance.redisClient == nil {
		return nil
	}
	if enqueueError := serviceInstance.redisClient.LPush(ctx, processingQueueKey, studyID.String()).Err(); enqueueError != nil {
		return fmt.Errorf("failed to enqueue dicom processing job: %w", enqueueError)
	}
	return nil
}

func (serviceInstance *service) discardStoredObjects(ctx context.Context, keys []string) {
	for _, key := range keys {
		if deleteError := serviceInstance.storageClient.Delete(context.WithoutCancel(ctx), key); deleteError != nil {
			slog.Warn("failed to discard orphaned dicom object", "error", deleteError)
		}
	}
}

func appendUniqueWarnings(existing []UploadWarning, additions ...UploadWarning) []UploadWarning {
	for _, addition := range additions {
		alreadyPresent := false
		for _, current := range existing {
			if current.Code == addition.Code {
				alreadyPresent = true
				break
			}
		}
		if !alreadyPresent {
			existing = append(existing, addition)
		}
	}
	return existing
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

	storedObjects, objectsError := serviceInstance.storedObjects(ctx, study)
	if objectsError != nil {
		return "", time.Time{}, objectsError
	}
	if len(storedObjects) > 1 {
		return "", time.Time{}, apperrors.ErrDICOMMultipleInstances
	}

	downloadURL, expiresAt, presignError := serviceInstance.presignObject(ctx, storedObjects[0].ObjectKey)
	if presignError != nil {
		return "", time.Time{}, presignError
	}

	serviceInstance.auditDownloadURLIssuance(ctx, parsedStudyID)
	return downloadURL, expiresAt, nil
}

func (serviceInstance *service) ListStudyInstanceDownloads(ctx context.Context, studyID string) ([]StudyInstanceDownload, error) {
	parsedStudyID, err := uuid.Parse(studyID)
	if err != nil {
		return nil, apperrors.InvalidArgument("invalid study instances input", map[string]string{"study_id": "must be a valid UUID"})
	}

	study, dbError := serviceInstance.dbRepository.GetImagingStudy(ctx, parsedStudyID)
	if dbError != nil {
		return nil, dbError
	}

	storedObjects, objectsError := serviceInstance.storedObjects(ctx, study)
	if objectsError != nil {
		return nil, objectsError
	}

	downloads := make([]StudyInstanceDownload, 0, len(storedObjects))
	for _, storedObject := range storedObjects {
		downloadURL, expiresAt, presignError := serviceInstance.presignObject(ctx, storedObject.ObjectKey)
		if presignError != nil {
			return nil, presignError
		}
		downloads = append(downloads, StudyInstanceDownload{Instance: storedObject, DownloadURL: downloadURL, ExpiresAt: expiresAt})
	}

	serviceInstance.auditDownloadURLIssuance(ctx, parsedStudyID)
	return downloads, nil
}

func (serviceInstance *service) storedObjects(ctx context.Context, study *ImagingStudy) ([]*InstanceIndex, error) {
	instances, listError := serviceInstance.dbRepository.ListInstances(ctx, study.ID)
	if listError != nil {
		return nil, listError
	}
	if len(instances) == 0 {
		return []*InstanceIndex{{
			ID:        study.ID,
			StudyID:   study.ID,
			FileName:  study.FileName,
			ObjectKey: legacyStudyObjectKey(study),
			CreatedAt: study.CreatedAt,
		}}, nil
	}
	for _, instance := range instances {
		if instance.ObjectKey == "" {
			instance.ObjectKey = legacyStudyObjectKey(study)
		}
	}
	return instances, nil
}

func (serviceInstance *service) presignObject(ctx context.Context, objectKey string) (string, time.Time, error) {
	downloadURL, presignError := serviceInstance.storageClient.GetPresignedURL(ctx, objectKey, downloadURLExpiration)
	if presignError != nil {
		return "", time.Time{}, fmt.Errorf("failed to generate presigned url: %w", presignError)
	}
	return downloadURL, time.Now().Add(downloadURLExpiration), nil
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
