package tests

import (
	"context"
	"testing"

	"github.com/healthcare/backend/internal/modules/health"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
)

func TestHealthCheck_Success(t *testing.T) {
	healthGRPCHandler := health.NewGRPCHandler(nil, nil)
	contextInstance := context.Background()

	response, checkError := healthGRPCHandler.Check(contextInstance, &grpc_health_v1.HealthCheckRequest{})

	assert.NoError(t, checkError)
	assert.NotNil(t, response)
	assert.Equal(t, grpc_health_v1.HealthCheckResponse_SERVING, response.Status)
}

func TestHealthWatch_Unimplemented(t *testing.T) {
	healthGRPCHandler := health.NewGRPCHandler(nil, nil)

	watchError := healthGRPCHandler.Watch(&grpc_health_v1.HealthCheckRequest{}, nil)

	assert.Error(t, watchError)
	grpcStatusError, ok := status.FromError(watchError)
	assert.True(t, ok)
	assert.Equal(t, codes.Unimplemented, grpcStatusError.Code())
}
