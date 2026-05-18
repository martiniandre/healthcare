package tests

import (
	"bytes"
	"context"
	"testing"

	"github.com/healthcare/backend/internal/modules/imaging"
	"github.com/healthcare/backend/internal/modules/imaging/mocks"
	"github.com/healthcare/backend/internal/shared/storage"
	"github.com/stretchr/testify/assert"
)

func TestService_UploadDICOMStream_Success(testingInstance *testing.T) {
	mockRepository := mocks.NewMockRepository()
	storageClient := storage.NewStorageClient()
	imagingService := imaging.NewService(mockRepository, storageClient, nil, "test-bucket")
	contextParam := context.Background()

	validDICOMBytes := make([]byte, 200)
	copy(validDICOMBytes[128:132], []byte("DICM"))

	streamReader := bytes.NewReader(validDICOMBytes)
	study, uploadError := imagingService.UploadDICOMStream(contextParam, "patient-123", "Brain MRI", "MR", streamReader)

	assert.NoError(testingInstance, uploadError)
	assert.NotNil(testingInstance, study)
	assert.Equal(testingInstance, "patient-123", study.PatientFhirID)
	assert.Equal(testingInstance, "MR", study.Modality)
	assert.Equal(testingInstance, "PENDING", study.Status)

	retrievedStudy, queryError := imagingService.GetImagingStudy(contextParam, study.ID)
	assert.NoError(testingInstance, queryError)
	assert.Equal(testingInstance, study.ID, retrievedStudy.ID)
}

func TestService_UploadDICOMStream_InvalidMagicBytes(testingInstance *testing.T) {
	mockRepository := mocks.NewMockRepository()
	storageClient := storage.NewStorageClient()
	imagingService := imaging.NewService(mockRepository, storageClient, nil, "test-bucket")
	contextParam := context.Background()

	invalidDICOMBytes := make([]byte, 200)
	copy(invalidDICOMBytes[128:132], []byte("JPEG"))

	streamReader := bytes.NewReader(invalidDICOMBytes)
	_, uploadError := imagingService.UploadDICOMStream(contextParam, "patient-123", "Brain MRI", "MR", streamReader)

	assert.Error(testingInstance, uploadError)
	assert.Contains(testingInstance, uploadError.Error(), "magic bytes DICM signature missing")
}

func TestService_UploadDICOMStream_TooSmall(testingInstance *testing.T) {
	mockRepository := mocks.NewMockRepository()
	storageClient := storage.NewStorageClient()
	imagingService := imaging.NewService(mockRepository, storageClient, nil, "test-bucket")
	contextParam := context.Background()

	tooSmallBytes := []byte("too-small")

	streamReader := bytes.NewReader(tooSmallBytes)
	_, uploadError := imagingService.UploadDICOMStream(contextParam, "patient-123", "Brain MRI", "MR", streamReader)

	assert.Error(testingInstance, uploadError)
	assert.Contains(testingInstance, uploadError.Error(), "file is too small")
}

func TestService_GetDownloadURL(testingInstance *testing.T) {
	mockRepository := mocks.NewMockRepository()
	storageClient := storage.NewStorageClient()
	imagingService := imaging.NewService(mockRepository, storageClient, nil, "test-bucket")
	contextParam := context.Background()

	validDICOMBytes := make([]byte, 200)
	copy(validDICOMBytes[128:132], []byte("DICM"))

	streamReader := bytes.NewReader(validDICOMBytes)
	study, uploadError := imagingService.UploadDICOMStream(contextParam, "patient-123", "Brain MRI", "MR", streamReader)
	assert.NoError(testingInstance, uploadError)

	downloadURL, expiresAt, queryError := imagingService.GetDownloadURL(contextParam, study.ID)
	assert.NoError(testingInstance, queryError)
	assert.NotEmpty(testingInstance, downloadURL)
	assert.True(testingInstance, expiresAt.After(study.CreatedAt))
}
