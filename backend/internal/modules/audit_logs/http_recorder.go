package audit_logs

import (
	"context"

	"github.com/healthcare/backend/internal/api/middleware"
)

type HTTPAuditRecorder struct {
	service Service
}

func NewHTTPAuditRecorder(service Service) *HTTPAuditRecorder {
	return &HTTPAuditRecorder{service: service}
}

func (httpAuditRecorder *HTTPAuditRecorder) RecordHTTPAudit(contextVal context.Context, auditEntry middleware.AuditEntry) error {
	_, createError := httpAuditRecorder.service.CreateResourceAuditLog(contextVal, ResourceAuditLog{
		CorrelationID: auditEntry.CorrelationID,
		CallerUserID:  auditEntry.CallerUserID,
		CallerRole:    auditEntry.CallerRole,
		Method:        auditEntry.Route,
		AccessGranted: auditEntry.AccessGranted,
		ResourceType:  auditEntry.ResourceType,
		ResourceID:    auditEntry.ResourceID,
		Action:        auditEntry.Action,
	})
	return createError
}
