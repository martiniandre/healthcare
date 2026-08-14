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
	query := `INSERT INTO audit_logs (id, correlation_id, caller_user_id, caller_role, method, access_granted, resource_type, resource_id, action, payload_diff, created_at)
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
	_, executionError := auditLogsRepository.dbPool.Exec(contextVal, query,
		auditLog.ID, auditLog.CorrelationID, auditLog.CallerUserID, auditLog.CallerRole, auditLog.Method, auditLog.AccessGranted,
		nullableString(auditLog.ResourceType), nullableString(auditLog.ResourceID), nullableString(auditLog.Action), nullableDiff(auditLog.PayloadDiff), auditLog.CreatedAt,
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

	query := `SELECT id, correlation_id, caller_user_id, caller_role, method, access_granted,
			         COALESCE(resource_type, ''), COALESCE(resource_id, ''), COALESCE(action, ''),
			         COALESCE(payload_diff, '{}'::jsonb), created_at
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
			&auditLog.ID, &auditLog.CorrelationID, &auditLog.CallerUserID, &auditLog.CallerRole, &auditLog.Method, &auditLog.AccessGranted,
			&auditLog.ResourceType, &auditLog.ResourceID, &auditLog.Action, &auditLog.PayloadDiff, &auditLog.CreatedAt,
		)
		if scanError != nil {
			return nil, 0, scanError
		}
		logs = append(logs, auditLog)
	}

	return logs, totalCount, nil
}

func nullableString(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}

func nullableDiff(payloadDiff map[string]any) any {
	if len(payloadDiff) == 0 {
		return nil
	}
	return payloadDiff
}
