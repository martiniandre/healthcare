package auth

import (
	"github.com/jackc/pgx/v5/pgxpool"
	authpb "github.com/healthcare/backend/internal/modules/auth/pb"
	"google.golang.org/grpc"
)

func Register(grpcServer *grpc.Server, dbPool *pgxpool.Pool) {
	repo := NewRepository(dbPool)
	svc := NewService(repo)
	handler := NewGRPCHandler(svc)
	authpb.RegisterAuthServiceServer(grpcServer, handler)
}
