package tests

import (
	"context"
	"io"
	"testing"

	"github.com/healthcare/backend/internal/modules/imaging"
	"github.com/healthcare/backend/internal/modules/imaging/pb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func buildUploadRequests(fileBytes []byte) []*pb.UploadDICOMRequest {
	return []*pb.UploadDICOMRequest{
		{
			Payload: &pb.UploadDICOMRequest_Metadata{
				Metadata: &pb.UploadMetadata{
					PatientFhirId: fixturePatientFhirID,
					Title:         "Brain Scan",
					Modality:      "MR",
				},
			},
		},
		{
			Payload: &pb.UploadDICOMRequest_Chunk{
				Chunk: fileBytes[:100],
			},
		},
		{
			Payload: &pb.UploadDICOMRequest_Chunk{
				Chunk: fileBytes[100:],
			},
		},
	}
}

func TestGRPCHandler_UploadDICOM_ProgressTracking(testingInstance *testing.T) {
	harness := newUploadHarness()
	grpcHandler := imaging.NewGRPCHandler(harness.service)
	fileBytes := newDICOMFixture().bytes()

	mockStream := &mockUploadDICOMServer{
		ContextParam: context.Background(),
		Requests:     buildUploadRequests(fileBytes),
	}

	uploadError := grpcHandler.UploadDICOM(mockStream)
	require.NoError(testingInstance, uploadError)

	require.Len(testingInstance, mockStream.Responses, 3)

	assert.Equal(testingInstance, "UPLOADING", mockStream.Responses[0].Status)
	assert.Equal(testingInstance, int64(100), mockStream.Responses[0].BytesUploaded)

	assert.Equal(testingInstance, "UPLOADING", mockStream.Responses[1].Status)
	assert.Equal(testingInstance, int64(len(fileBytes)), mockStream.Responses[1].BytesUploaded)

	assert.Equal(testingInstance, imaging.ImagingStudyStatusReceiving, mockStream.Responses[2].Status)
	assert.Equal(testingInstance, int64(len(fileBytes)), mockStream.Responses[2].BytesUploaded)
	assert.NotEmpty(testingInstance, mockStream.Responses[2].ImagingStudyId)
}

func TestGRPCHandler_UploadDICOM_PatientMismatchReturnsFailedPrecondition(testingInstance *testing.T) {
	harness := newUploadHarness()
	harness.identity.Identity = imaging.PatientIdentity{FullNames: []string{"Pedro Santos"}, DocumentIDs: []string{"11144477735"}}
	grpcHandler := imaging.NewGRPCHandler(harness.service)

	mockStream := &mockUploadDICOMServer{
		ContextParam: context.Background(),
		Requests:     buildUploadRequests(newDICOMFixture().bytes()),
	}

	uploadError := grpcHandler.UploadDICOM(mockStream)

	require.Error(testingInstance, uploadError)
	assert.Equal(testingInstance, codes.FailedPrecondition, status.Code(uploadError))
}

func TestGRPCHandler_UploadDICOM_DuplicateReturnsAlreadyExists(testingInstance *testing.T) {
	harness := newUploadHarness()
	grpcHandler := imaging.NewGRPCHandler(harness.service)
	fileBytes := newDICOMFixture().bytes()

	firstStream := &mockUploadDICOMServer{ContextParam: context.Background(), Requests: buildUploadRequests(fileBytes)}
	require.NoError(testingInstance, grpcHandler.UploadDICOM(firstStream))

	secondStream := &mockUploadDICOMServer{ContextParam: context.Background(), Requests: buildUploadRequests(fileBytes)}
	duplicateError := grpcHandler.UploadDICOM(secondStream)

	require.Error(testingInstance, duplicateError)
	assert.Equal(testingInstance, codes.AlreadyExists, status.Code(duplicateError))
}
