package auth

import (
	authpb "github.com/healthcare/backend/internal/modules/auth/pb"
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
	authpb.RegisterAuthServiceServer(grpcServer, handler)
	return svc
}
