package patients

import (
	"github.com/healthcare/backend/internal/shared/healthcare"
	patientspb "github.com/healthcare/backend/internal/modules/patients/pb"
	"google.golang.org/grpc"
)

type Dependency struct {
	FHIRClient healthcare.FHIRClient
}

func Register(grpcServer *grpc.Server, dep Dependency) Service {
	repo := NewRepository(dep.FHIRClient)
	svc := NewService(repo)
	handler := NewGRPCHandler(svc)
	patientspb.RegisterPatientServiceServer(grpcServer, handler)
	return svc
}
