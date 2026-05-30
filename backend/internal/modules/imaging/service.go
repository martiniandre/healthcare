package imaging

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/google/uuid"
	"github.com/healthcare/backend/internal/shared/storage"
	"github.com/healthcare/backend/internal/shared/validator"
	"github.com/redis/go-redis/v9"
)

const MaxDICOMUploadBytes int64 = 2 << 30

type Service interface {
	UploadDICOMStream(ctx context.Context, patientFhirID, title, modality string, streamReader io.Reader) (*ImagingStudy, error)
	GetImagingStudy(ctx context.Context, studyID uuid.UUID) (*ImagingStudy, error)
	ListImagingStudies(ctx context.Context, patientFhirID string) ([]*ImagingStudy, error)
	GetDownloadURL(ctx context.Context, studyID uuid.UUID) (string, time.Time, error)
}

type service struct {
	dbRepository  Repository
	storageClient storage.StorageClient
	redisClient   *redis.Client
	bucketName    string
}

func NewService(dbRepository Repository, storageClient storage.StorageClient, redisClient *redis.Client, bucketName string) Service {
	return &service{
		dbRepository:  dbRepository,
		storageClient: storageClient,
		redisClient:   redisClient,
		bucketName:    bucketName,
	}
}

func (serviceInstance *service) UploadDICOMStream(ctx context.Context, patientFhirID, title, modality string, streamReader io.Reader) (*ImagingStudy, error) {
	if patientFhirID == "" {
		return nil, errors.New("patient fhir id is required")
	}
	if title == "" {
		return nil, errors.New("title is required")
	}
	if !validator.IsValidDICOMModality(modality) {
		return nil, errors.New("invalid dicom modality")
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
		return nil, fmt.Errorf("%w: preamble is too small", ErrInvalidDICOM)
	}

	magicBytesSignature := string(headerBytes[128:132])
	if magicBytesSignature != "DICM" {
		return nil, fmt.Errorf("%w: magic bytes DICM signature missing", ErrInvalidDICOM)
	}

	reconstructedReader := io.MultiReader(bytes.NewReader(headerBytes[:bytesRead]), limitedReader)
	boundedReader := &uploadLimitReader{reader: reconstructedReader, remainingBytes: MaxDICOMUploadBytes}

	studyID := uuid.New()
	objectPath := fmt.Sprintf("dicom/staging/%s/%s.dcm", patientFhirID, studyID.String())

	gcsStagingURL, uploadError := serviceInstance.storageClient.Upload(ctx, serviceInstance.bucketName, objectPath, boundedReader)
	if uploadError != nil {
		return nil, fmt.Errorf("failed to upload dicom to storage: %w", uploadError)
	}

	study := &ImagingStudy{
		ID:               studyID,
		PatientFhirID:    patientFhirID,
		Title:            title,
		Modality:         modality,
		GCSPath:          gcsStagingURL,
		StudyInstanceUID: "",
		Status:           ImagingStudyStatusPending,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	dbError := serviceInstance.dbRepository.CreateImagingStudy(ctx, study)
	if dbError != nil {
		return nil, fmt.Errorf("failed to persist imaging study operational record: %w", dbError)
	}

	if serviceInstance.redisClient != nil {
		enqueueError := serviceInstance.redisClient.LPush(ctx, "dicom_processing_queue", studyID.String()).Err()
		if enqueueError != nil {
			return nil, fmt.Errorf("failed to enqueue dicom processing job: %w", enqueueError)
		}
	}

	return study, nil
}

type uploadLimitReader struct {
	reader         io.Reader
	remainingBytes int64
}

func (reader *uploadLimitReader) Read(buffer []byte) (int, error) {
	if reader.remainingBytes <= 0 {
		return 0, ErrDICOMTooLarge
	}
	if int64(len(buffer)) > reader.remainingBytes {
		buffer = buffer[:int(reader.remainingBytes)]
	}
	bytesRead, readError := reader.reader.Read(buffer)
	reader.remainingBytes -= int64(bytesRead)
	if readError == nil && reader.remainingBytes == 0 {
		return bytesRead, ErrDICOMTooLarge
	}
	return bytesRead, readError
}

func (serviceInstance *service) GetImagingStudy(ctx context.Context, studyID uuid.UUID) (*ImagingStudy, error) {
	return serviceInstance.dbRepository.GetImagingStudy(ctx, studyID)
}

func (serviceInstance *service) ListImagingStudies(ctx context.Context, patientFhirID string) ([]*ImagingStudy, error) {
	if patientFhirID == "" {
		return nil, errors.New("patient fhir id is required")
	}
	return serviceInstance.dbRepository.ListImagingStudiesByPatient(ctx, patientFhirID)
}

func (serviceInstance *service) GetDownloadURL(ctx context.Context, studyID uuid.UUID) (string, time.Time, error) {
	study, dbError := serviceInstance.dbRepository.GetImagingStudy(ctx, studyID)
	if dbError != nil {
		return "", time.Time{}, dbError
	}

	expirationDuration := 15 * time.Minute
	downloadURL, presignError := serviceInstance.storageClient.GetPresignedURL(ctx, serviceInstance.bucketName, study.GCSPath, expirationDuration)
	if presignError != nil {
		return "", time.Time{}, fmt.Errorf("failed to generate presigned url: %w", presignError)
	}

	expiresAt := time.Now().Add(expirationDuration)
	return downloadURL, expiresAt, nil
}
