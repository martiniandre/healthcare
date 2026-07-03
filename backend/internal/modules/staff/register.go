package staff

import (
	"github.com/jackc/pgx/v5/pgxpool"
	staffpb "github.com/healthcare/backend/internal/modules/staff/pb"
	"google.golang.org/grpc"
)

type Dependency struct {
	DB *pgxpool.Pool
}

func Register(grpcServer *grpc.Server, dep Dependency) Service {
	repo := NewRepository(dep.DB)
	svc := NewService(repo)
	handler := NewGRPCHandler(svc)
	staffpb.RegisterStaffServiceServer(grpcServer, handler)
	return svc
}
