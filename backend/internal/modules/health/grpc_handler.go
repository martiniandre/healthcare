package health

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc/health/grpc_health_v1"

	"github.com/healthcare/backend/internal/shared/apperrors"
)

type GRPCHandler struct {
	grpc_health_v1.UnimplementedHealthServer
	databaseConnectionPool *pgxpool.Pool
	cacheConnectionClient  *redis.Client
}

func NewGRPCHandler(databaseConnectionPool *pgxpool.Pool, cacheConnectionClient *redis.Client) *GRPCHandler {
	return &GRPCHandler{
		databaseConnectionPool: databaseConnectionPool,
		cacheConnectionClient:  cacheConnectionClient,
	}
}

func (handler *GRPCHandler) Check(ctx context.Context, request *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	if handler.databaseConnectionPool != nil {
		databasePingError := handler.databaseConnectionPool.Ping(ctx)
		if databasePingError != nil {
			return &grpc_health_v1.HealthCheckResponse{
				Status: grpc_health_v1.HealthCheckResponse_NOT_SERVING,
			}, nil
		}
	}

	if handler.cacheConnectionClient != nil {
		cachePingError := handler.cacheConnectionClient.Ping(ctx).Err()
		if cachePingError != nil {
			return &grpc_health_v1.HealthCheckResponse{
				Status: grpc_health_v1.HealthCheckResponse_NOT_SERVING,
			}, nil
		}
	}

	return &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	}, nil
}

func (handler *GRPCHandler) Watch(request *grpc_health_v1.HealthCheckRequest, stream grpc_health_v1.Health_WatchServer) error {
	return apperrors.ErrNotImplemented.ToGRPC()
}
