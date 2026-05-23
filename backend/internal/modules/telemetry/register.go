package telemetry

import (
	"github.com/jackc/pgx/v5/pgxpool"
	telemetrypb "github.com/healthcare/backend/internal/modules/telemetry/pb"
	"google.golang.org/grpc"
)

func Register(grpcServer *grpc.Server, dbPool *pgxpool.Pool) Service {
	repo := NewRepository(dbPool)
	svc := NewService(repo)
	handler := NewGRPCHandler(svc)
	telemetrypb.RegisterTelemetryServiceServer(grpcServer, handler)
	return svc
}
