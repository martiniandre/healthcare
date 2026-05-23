package clinical

import (
	"github.com/healthcare/backend/internal/shared/healthcare"
	clinicalpb "github.com/healthcare/backend/internal/modules/clinical/pb"
	"google.golang.org/grpc"
)

func Register(grpcServer *grpc.Server, fhirClient healthcare.FHIRClient) Service {
	repo := NewRepository(fhirClient)
	svc := NewService(repo)
	handler := NewGRPCHandler(svc)
	clinicalpb.RegisterClinicalServiceServer(grpcServer, handler)
	return svc
}
