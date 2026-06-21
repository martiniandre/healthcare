package audit_logs

import (
	"time"

	"github.com/google/uuid"
)

type AuditLog struct {
	ID            uuid.UUID `db:"id" json:"id"`
	CorrelationID string    `db:"correlation_id" json:"correlation_id"`
	CallerUserID  string    `db:"caller_user_id" json:"caller_user_id"`
	CallerRole    string    `db:"caller_role" json:"caller_role"`
	Method        string    `db:"method" json:"method"`
	AccessGranted bool      `db:"access_granted" json:"access_granted"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
}
