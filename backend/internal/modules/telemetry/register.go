package telemetry

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	telemetrypb "github.com/healthcare/backend/internal/modules/telemetry/pb"
	"google.golang.org/grpc"
)

type Dependency struct {
	DB *pgxpool.Pool
}

func Register(grpcServer *grpc.Server, dep Dependency) Service {
	repo := NewRepository(dep.DB)
	svc := NewService(repo)
	handler := NewGRPCHandler(svc)
	telemetrypb.RegisterTelemetryServiceServer(grpcServer, handler)
	return svc
}

func StartSimulator(ctx context.Context, dbPool *pgxpool.Pool) *Simulator {
	repo := NewRepository(dbPool)
	simulator := NewSimulator(repo)
	simulator.Start(ctx)
	return simulator
}
