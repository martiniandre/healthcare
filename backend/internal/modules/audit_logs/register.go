package audit_logs

import (
	auditlogspb "github.com/healthcare/backend/internal/modules/audit_logs/pb"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
)

type Dependency struct {
	DB *pgxpool.Pool
}

func Register(grpcServer *grpc.Server, dep Dependency) Service {
	auditLogsRepository := NewRepository(dep.DB)
	auditLogsService := NewService(auditLogsRepository)
	auditLogsGRPCHandler := NewGRPCHandler(auditLogsService)
	auditlogspb.RegisterAuditLogsServiceServer(grpcServer, auditLogsGRPCHandler)
	return auditLogsService
}
