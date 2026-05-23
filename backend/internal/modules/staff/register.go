package staff

import (
	"github.com/jackc/pgx/v5/pgxpool"
	staffpb "github.com/healthcare/backend/internal/modules/staff/pb"
	"google.golang.org/grpc"
)

func Register(grpcServer *grpc.Server, dbPool *pgxpool.Pool) Service {
	repo := NewRepository(dbPool)
	svc := NewService(repo)
	handler := NewGRPCHandler(svc)
	staffpb.RegisterStaffServiceServer(grpcServer, handler)
	return svc
}
