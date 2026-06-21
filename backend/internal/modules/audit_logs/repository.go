package audit_logs

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	CreateAuditLog(contextVal context.Context, auditLog *AuditLog) error
	ListAuditLogs(contextVal context.Context, limit int, offset int) ([]*AuditLog, int, error)
}

type repository struct {
	dbPool *pgxpool.Pool
}

func NewRepository(dbPool *pgxpool.Pool) Repository {
	return &repository{dbPool: dbPool}
}

func (auditLogsRepository *repository) CreateAuditLog(contextVal context.Context, auditLog *AuditLog) error {
	query := `INSERT INTO audit_logs (id, correlation_id, caller_user_id, caller_role, method, access_granted, created_at)
			  VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, executionError := auditLogsRepository.dbPool.Exec(contextVal, query,
		auditLog.ID, auditLog.CorrelationID, auditLog.CallerUserID, auditLog.CallerRole, auditLog.Method, auditLog.AccessGranted, auditLog.CreatedAt,
	)
	return executionError
}

func (auditLogsRepository *repository) ListAuditLogs(contextVal context.Context, limit int, offset int) ([]*AuditLog, int, error) {
	countQuery := `SELECT COUNT(*) FROM audit_logs`
	var totalCount int
	countError := auditLogsRepository.dbPool.QueryRow(contextVal, countQuery).Scan(&totalCount)
	if countError != nil {
		return nil, 0, countError
	}

	query := `SELECT id, correlation_id, caller_user_id, caller_role, method, access_granted, created_at
			  FROM audit_logs
			  ORDER BY created_at DESC
			  LIMIT $1 OFFSET $2`

	rows, queryError := auditLogsRepository.dbPool.Query(contextVal, query, limit, offset)
	if queryError != nil {
		return nil, 0, queryError
	}
	defer rows.Close()

	logs := make([]*AuditLog, 0)
	for rows.Next() {
		auditLog := &AuditLog{}
		scanError := rows.Scan(
			&auditLog.ID, &auditLog.CorrelationID, &auditLog.CallerUserID, &auditLog.CallerRole, &auditLog.Method, &auditLog.AccessGranted, &auditLog.CreatedAt,
		)
		if scanError != nil {
			return nil, 0, scanError
		}
		logs = append(logs, auditLog)
	}

	return logs, totalCount, nil
}
