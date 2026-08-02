package audit_logs

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Service interface {
	CreateAuditLog(contextVal context.Context, correlationID string, callerUserID string, callerRole string, method string, accessGranted bool) (*AuditLog, error)
	CreateResourceAuditLog(contextVal context.Context, resourceAuditLog ResourceAuditLog) (*AuditLog, error)
	ListAuditLogs(contextVal context.Context, limit int, offset int) ([]*AuditLog, int, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (auditLogsService *service) CreateAuditLog(contextVal context.Context, correlationID string, callerUserID string, callerRole string, method string, accessGranted bool) (*AuditLog, error) {
	auditLog := &AuditLog{
		ID:            uuid.New(),
		CorrelationID: correlationID,
		CallerUserID:  callerUserID,
		CallerRole:    callerRole,
		Method:        method,
		AccessGranted: accessGranted,
		CreatedAt:     time.Now(),
	}

	saveError := auditLogsService.repo.CreateAuditLog(contextVal, auditLog)
	if saveError != nil {
		return nil, saveError
	}

	return auditLog, nil
}

func (auditLogsService *service) CreateResourceAuditLog(contextVal context.Context, resourceAuditLog ResourceAuditLog) (*AuditLog, error) {
	auditLog := &AuditLog{
		ID:            uuid.New(),
		CorrelationID: resourceAuditLog.CorrelationID,
		CallerUserID:  resourceAuditLog.CallerUserID,
		CallerRole:    resourceAuditLog.CallerRole,
		Method:        resourceAuditLog.Method,
		AccessGranted: resourceAuditLog.AccessGranted,
		ResourceType:  resourceAuditLog.ResourceType,
		ResourceID:    resourceAuditLog.ResourceID,
		Action:        resourceAuditLog.Action,
		PayloadDiff:   resourceAuditLog.PayloadDiff,
		CreatedAt:     time.Now(),
	}

	saveError := auditLogsService.repo.CreateAuditLog(contextVal, auditLog)
	if saveError != nil {
		return nil, saveError
	}

	return auditLog, nil
}

func (auditLogsService *service) ListAuditLogs(contextVal context.Context, limit int, offset int) ([]*AuditLog, int, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}
	return auditLogsService.repo.ListAuditLogs(contextVal, limit, offset)
}
