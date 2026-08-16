package schedule

import (
	"time"

	"github.com/google/uuid"
)

type CreateUnavailabilityInput struct {
	StaffID  uuid.UUID
	StartsAt time.Time
	EndsAt   time.Time
	Reason   string
}

type ListUnavailabilityInput struct {
	StaffID uuid.UUID
	From    time.Time
	To      time.Time
}
