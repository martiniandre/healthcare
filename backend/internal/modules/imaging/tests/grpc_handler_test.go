package tests

import (
	"context"
	"io"
	"testing"

	"github.com/healthcare/backend/internal/modules/imaging"
	"github.com/healthcare/backend/internal/modules/imaging/mocks"
	"github.com/healthcare/backend/internal/modules/imaging/pb"
	"github.com/healthcare/backend/internal/shared/storage"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

type mockUploadDICOMServer struct {
	grpc.ServerStream
	ContextParam context.Context
	Requests     []*pb.UploadDICOMRequest
	Responses    []*pb.UploadDICOMResponse
	RequestIndex int
}

func (mockStream *mockUploadDICOMServer) Context() context.Context {
	return mockStream.ContextParam
}

func (mockStream *mockUploadDICOMServer) Recv() (*pb.UploadDICOMRequest, error) {
	if mockStream.RequestIndex >= len(mockStream.Requests) {
		return nil, io.EOF
	}
	request := mockStream.Requests[mockStream.RequestIndex]
	mockStream.RequestIndex++
	return request, nil
}

func (mockStream *mockUploadDICOMServer) Send(response *pb.UploadDICOMResponse) error {
	mockStream.Responses = append(mockStream.Responses, response)
	return nil
}

func TestGRPCHandler_UploadDICOM_ProgressTracking(testingInstance *testing.T) {
	mockRepository := mocks.NewMockRepository()
	storageClient := storage.NewStorageClient()
	imagingService := imaging.NewService(mockRepository, storageClient, nil, "test-bucket")
	grpcHandler := imaging.NewGRPCHandler(imagingService)

	validDICOMBytes := make([]byte, 200)
	copy(validDICOMBytes[128:132], []byte("DICM"))

	requests := []*pb.UploadDICOMRequest{
		{
			Payload: &pb.UploadDICOMRequest_Metadata{
				Metadata: &pb.UploadMetadata{
					PatientFhirId: "patient-123",
					Title:         "Brain Scan",
					Modality:      "MR",
				},
			},
		},
		{
			Payload: &pb.UploadDICOMRequest_Chunk{
				Chunk: validDICOMBytes[:100],
			},
		},
		{
			Payload: &pb.UploadDICOMRequest_Chunk{
				Chunk: validDICOMBytes[100:],
			},
		},
	}

	mockStream := &mockUploadDICOMServer{
		ContextParam: context.Background(),
		Requests:     requests,
	}

	uploadError := grpcHandler.UploadDICOM(mockStream)
	assert.NoError(testingInstance, uploadError)

	assert.NotEmpty(testingInstance, mockStream.Responses)
	assert.Len(testingInstance, mockStream.Responses, 3)

	assert.Equal(testingInstance, "UPLOADING", mockStream.Responses[0].Status)
	assert.Equal(testingInstance, int64(100), mockStream.Responses[0].BytesUploaded)

	assert.Equal(testingInstance, "UPLOADING", mockStream.Responses[1].Status)
	assert.Equal(testingInstance, int64(200), mockStream.Responses[1].BytesUploaded)

	assert.Equal(testingInstance, "PENDING", mockStream.Responses[2].Status)
	assert.Equal(testingInstance, int64(200), mockStream.Responses[2].BytesUploaded)
	assert.NotEmpty(testingInstance, mockStream.Responses[2].ImagingStudyId)
}
