package diagnostic_report

import (
	diagnosticreportpb "github.com/healthcare/backend/internal/modules/diagnostic_report/pb"
	"github.com/healthcare/backend/internal/modules/audit_logs"
	"github.com/healthcare/backend/internal/shared/eventbus"
	"github.com/healthcare/backend/internal/shared/healthcare"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
)

type Dependency struct {
	FHIRClient   healthcare.FHIRClient
	EventBus     eventbus.Bus
	DB           *pgxpool.Pool
	AuditService audit_logs.Service
}

func Register(grpcServer *grpc.Server, dep Dependency) Service {
	repo := NewRepository(dep.FHIRClient)
	versionRepository := NewVersionRepository(dep.DB)
	svc := NewService(repo, versionRepository, dep.EventBus, dep.AuditService)
	handler := NewGRPCHandler(svc)
	diagnosticreportpb.RegisterDiagnosticReportServiceServer(grpcServer, handler)
	return svc
}
