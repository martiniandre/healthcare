package audit_logs

import (
	"time"

	"github.com/google/uuid"
)

type AuditLog struct {
	ID            uuid.UUID     `db:"id" json:"id"`
	CorrelationID string        `db:"correlation_id" json:"correlation_id"`
	CallerUserID  string        `db:"caller_user_id" json:"caller_user_id"`
	CallerRole    string        `db:"caller_role" json:"caller_role"`
	Method        string        `db:"method" json:"method"`
	AccessGranted bool          `db:"access_granted" json:"access_granted"`
	ResourceType  string        `db:"resource_type" json:"resource_type,omitempty"`
	ResourceID    string        `db:"resource_id" json:"resource_id,omitempty"`
	Action        string        `db:"action" json:"action,omitempty"`
	PayloadDiff   map[string]any `db:"payload_diff" json:"payload_diff,omitempty"`
	CreatedAt     time.Time     `db:"created_at" json:"created_at"`
}

type ResourceAuditLog struct {
	CorrelationID string
	CallerUserID  string
	CallerRole    string
	Method        string
	AccessGranted bool
	ResourceType  string
	ResourceID    string
	Action        string
	PayloadDiff   map[string]any
}
