package schedule

import (
	"time"

	"github.com/google/uuid"
)

type StaffUnavailability struct {
	ID        uuid.UUID
	StaffID   uuid.UUID
	StartsAt  time.Time
	EndsAt    time.Time
	Reason    string
	CreatedBy *uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
}
