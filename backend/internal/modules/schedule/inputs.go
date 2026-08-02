package schedule

import (
	"time"

	"github.com/google/uuid"
)

type CreateAppointmentInput struct {
	PatientFHIRID  string
	StaffID        uuid.UUID
	StartsAt       time.Time
	EndsAt         time.Time
	Reason         string
	IdempotencyKey string
	RequestHash    string
}
