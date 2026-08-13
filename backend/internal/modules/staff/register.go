package staff

import (
	staffpb "github.com/healthcare/backend/internal/modules/staff/pb"
	"github.com/healthcare/backend/internal/shared/healthcare"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
)

type Dependency struct {
	DB         *pgxpool.Pool
	FHIRClient healthcare.FHIRClient
}

func Register(grpcServer *grpc.Server, dep Dependency) Service {
	repo := NewRepository(dep.DB)
	svc := NewService(repo, dep.FHIRClient)
	handler := NewGRPCHandler(svc)
	staffpb.RegisterStaffServiceServer(grpcServer, handler)
	return svc
}
