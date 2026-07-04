package telemetry

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	telemetrypb "github.com/healthcare/backend/internal/modules/telemetry/pb"
	"github.com/healthcare/backend/internal/shared/eventbus"
	"google.golang.org/grpc"
)

type Dependency struct {
	DB       *pgxpool.Pool
	EventBus eventbus.Bus
}

func Register(grpcServer *grpc.Server, dep Dependency) Service {
	repo := NewRepository(dep.DB)
	svc := NewService(repo, dep.EventBus)
	handler := NewGRPCHandler(svc)
	telemetrypb.RegisterTelemetryServiceServer(grpcServer, handler)
	return svc
}

func StartSimulator(ctx context.Context, dbPool *pgxpool.Pool, eventBus eventbus.Bus) *Simulator {
	repo := NewRepository(dbPool)
	simulator := NewSimulator(repo, eventBus)
	simulator.Start(ctx)
	return simulator
}
