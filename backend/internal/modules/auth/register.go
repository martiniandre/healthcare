package auth

import (
	"github.com/jackc/pgx/v5/pgxpool"
	authpb "github.com/healthcare/backend/internal/modules/auth/pb"
	"google.golang.org/grpc"
)

type Dependency struct {
	DB *pgxpool.Pool
}

func Register(grpcServer *grpc.Server, dep Dependency) Service {
	repo := NewRepository(dep.DB)
	svc := NewService(repo)
	handler := NewGRPCHandler(svc)
	authpb.RegisterAuthServiceServer(grpcServer, handler)
	return svc
}
