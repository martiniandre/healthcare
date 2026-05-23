package patients

import (
	"github.com/healthcare/backend/internal/shared/healthcare"
	patientspb "github.com/healthcare/backend/internal/modules/patients/pb"
	"google.golang.org/grpc"
)

func Register(grpcServer *grpc.Server, fhirClient healthcare.FHIRClient) Service {
	repo := NewRepository(fhirClient)
	svc := NewService(repo)
	handler := NewGRPCHandler(svc)
	patientspb.RegisterPatientServiceServer(grpcServer, handler)
	return svc
}
