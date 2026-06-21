package audit_logs

import (
	"time"

	"github.com/google/uuid"
)

type AuditLog struct {
	ID            uuid.UUID `db:"id"`
	CorrelationID string    `db:"correlation_id"`
	CallerUserID  string    `db:"caller_user_id"`
	CallerRole    string    `db:"caller_role"`
	Method        string    `db:"method"`
	AccessGranted bool      `db:"access_granted"`
	CreatedAt     time.Time `db:"created_at"`
}
