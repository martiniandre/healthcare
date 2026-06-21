package audit_logs

import (
	auditlogspb "github.com/healthcare/backend/internal/modules/audit_logs/pb"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
)

func Register(grpcServer *grpc.Server, dbPool *pgxpool.Pool) Service {
	auditLogsRepository := NewRepository(dbPool)
	auditLogsService := NewService(auditLogsRepository)
	auditLogsGRPCHandler := NewGRPCHandler(auditLogsService)
	auditlogspb.RegisterAuditLogsServiceServer(grpcServer, auditLogsGRPCHandler)
	return auditLogsService
}
