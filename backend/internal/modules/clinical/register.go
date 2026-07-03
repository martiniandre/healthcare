package clinical

import (
	"github.com/healthcare/backend/internal/shared/healthcare"
	clinicalpb "github.com/healthcare/backend/internal/modules/clinical/pb"
	"google.golang.org/grpc"
)

type Dependency struct {
	FHIRClient healthcare.FHIRClient
}

func Register(grpcServer *grpc.Server, dep Dependency) Service {
	repo := NewRepository(dep.FHIRClient)
	svc := NewService(repo)
	handler := NewGRPCHandler(svc)
	clinicalpb.RegisterClinicalServiceServer(grpcServer, handler)
	return svc
}
