package telemetry

import (
	"context"

	telemetrypb "github.com/healthcare/backend/internal/modules/telemetry/pb"
	"github.com/healthcare/backend/internal/shared/eventbus"
	"github.com/jackc/pgx/v5/pgxpool"
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

func ShouldRunSimulator(appEnv string, explicitlyEnabled bool) bool {
	return explicitlyEnabled || appEnv == "development" || appEnv == "test"
}

func StartSimulator(ctx context.Context, dbPool *pgxpool.Pool, eventBus eventbus.Bus) *Simulator {
	repo := NewRepository(dbPool)
	simulator := NewSimulator(repo, eventBus)
	simulator.Start(ctx)
	return simulator
}
