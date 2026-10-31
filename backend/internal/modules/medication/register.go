package medication

import (
	medicationpb "github.com/healthcare/backend/internal/modules/medication/pb"
	"github.com/healthcare/backend/internal/shared/healthcare"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
)

type Dependency struct {
	FHIRClient    healthcare.FHIRClient
	DB            *pgxpool.Pool
	AllergyReader AllergyReader
	AuditRecorder AuditRecorder
}

func Register(grpcServer *grpc.Server, dep Dependency) Service {
	repo := NewRepository(dep.FHIRClient)
	safetyRepository := NewSafetyRepository(dep.DB)
	svc := NewService(repo, SafetyDependencies{
		Allergies: dep.AllergyReader,
		Catalog:   safetyRepository,
		Overrides: safetyRepository,
		Audit:     dep.AuditRecorder,
	})
	handler := NewGRPCHandler(svc)
	medicationpb.RegisterMedicationServiceServer(grpcServer, handler)
	return svc
}
