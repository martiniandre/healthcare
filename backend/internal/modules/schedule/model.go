package schedule

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type AppointmentStatus string

const (
	AppointmentStatusScheduled AppointmentStatus = "scheduled"
	AppointmentStatusConfirmed AppointmentStatus = "confirmed"
	AppointmentStatusCancelled AppointmentStatus = "cancelled"
	AppointmentStatusFinished  AppointmentStatus = "finished"
)

type Appointment struct {
	ID            uuid.UUID
	PatientFHIRID string
	StaffID       uuid.UUID
	StartsAt      time.Time
	EndsAt        time.Time
	Status        AppointmentStatus
	Reason        string
	Version       int
	CreatedBy     *uuid.UUID
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type IdempotencyKey struct {
	ID             uuid.UUID
	IdempotencyKey string
	RequestHash    string
	ResponseStatus int
	ResponseBody   json.RawMessage
	CreatedAt      time.Time
}
