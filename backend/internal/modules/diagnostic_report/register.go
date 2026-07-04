package diagnostic_report

import (
	diagnosticreportpb "github.com/healthcare/backend/internal/modules/diagnostic_report/pb"
	"github.com/healthcare/backend/internal/shared/healthcare"
	"google.golang.org/grpc"
)

type Dependency struct {
	FHIRClient healthcare.FHIRClient
}

func Register(grpcServer *grpc.Server, dep Dependency) Service {
	repo := NewRepository(dep.FHIRClient)
	svc := NewService(repo)
	handler := NewGRPCHandler(svc)
	diagnosticreportpb.RegisterDiagnosticReportServiceServer(grpcServer, handler)
	return svc
}
