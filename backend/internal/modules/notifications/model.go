package notifications

import (
	"time"

	"github.com/google/uuid"
)

type NotificationType string

const (
	NotificationTypeTelemetryAlert  NotificationType = "telemetry_alert"
	NotificationTypeExamComplete    NotificationType = "exam_complete"
	NotificationTypeEncounterCreate NotificationType = "encounter_created"
	NotificationTypeEncounterUpdate NotificationType = "encounter_updated"
	NotificationTypePatientCreate   NotificationType = "patient_created"
	NotificationTypePatientUpdate   NotificationType = "patient_updated"
	NotificationTypeAuditAlert      NotificationType = "audit_alert"
	NotificationTypeSystem          NotificationType = "system"
)

type NotificationPriority string

const (
	PriorityCritical NotificationPriority = "critical"
	PriorityHigh     NotificationPriority = "high"
	PriorityMedium   NotificationPriority = "medium"
	PriorityLow      NotificationPriority = "low"
)

type Notification struct {
	ID           uuid.UUID            `db:"id"`
	Type         NotificationType     `db:"type"`
	Priority     NotificationPriority `db:"priority"`
	Title        string               `db:"title"`
	Body         string               `db:"body"`
	ActorID      *uuid.UUID           `db:"actor_id"`
	ResourceType string               `db:"resource_type"`
	ResourceID   string               `db:"resource_id"`
	IsRead       bool                 `db:"is_read"`
	CreatedAt    time.Time            `db:"created_at"`
}

type NotificationRecipient struct {
	NotificationID uuid.UUID  `db:"notification_id"`
	UserID         uuid.UUID  `db:"user_id"`
	IsRead         bool       `db:"is_read"`
	ReadAt         *time.Time `db:"read_at"`
}
