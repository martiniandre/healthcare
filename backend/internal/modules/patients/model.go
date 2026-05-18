package patients

import (
	"time"

	"github.com/google/uuid"
)

type Patient struct {
	ID             uuid.UUID
	FHIRResourceID string
	FullName       string
	BirthDate      time.Time
	DocumentID     string
	PhoneNumber    string
	IsActive       bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
