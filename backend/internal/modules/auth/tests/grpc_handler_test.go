package tests

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/healthcare/backend/internal/modules/auth"
	"github.com/healthcare/backend/internal/modules/auth/mocks"
	"github.com/healthcare/backend/internal/modules/auth/pb"
	"github.com/healthcare/backend/internal/shared/role"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type mockServerTransportStream struct {
	grpc.ServerTransportStream
	header metadata.MD
}

func (mockStream *mockServerTransportStream) SetHeader(metadataParam metadata.MD) error {
	mockStream.header = metadataParam
	return nil
}

func TestGRPCHandler_Login(testingInstance *testing.T) {
	mockService := mocks.NewMockService()
	grpcHandler := auth.NewGRPCHandler(mockService)
	
	stream := &mockServerTransportStream{}
	contextParam := grpc.NewContextWithServerTransportStream(context.Background(), stream)

	mockService.User = &auth.User{
		ID:   uuid.New(),
		Role: role.RoleDoctor,
	}
	mockService.Token = "fake-jwt-token"

	request := &pb.LoginRequest{
		Email:    "doctor@example.com",
		Password: "password123",
	}

	response, loginError := grpcHandler.Login(contextParam, request)

	assert.NoError(testingInstance, loginError)
	assert.NotNil(testingInstance, response)
	assert.Equal(testingInstance, mockService.Token, response.Token)
	assert.Equal(testingInstance, string(role.RoleDoctor), response.Role)

	cookieHeader := stream.header.Get("set-cookie")
	assert.Len(testingInstance, cookieHeader, 2)
	assert.Contains(testingInstance, cookieHeader[0], "token=fake-jwt-token")
	assert.Contains(testingInstance, cookieHeader[1], "csrf_token=")
}
