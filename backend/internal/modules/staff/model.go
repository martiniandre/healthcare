package staff

import (
	"time"

	"github.com/google/uuid"
	"github.com/healthcare/backend/internal/shared/role"
)

type Employee struct {
	ID             uuid.UUID  `db:"id"`
	FullName       string     `db:"full_name"`
	Email          string     `db:"email"`
	Role           role.Role  `db:"role"`
	CRMNumber      *string    `db:"crm_number"`
	FHIRResourceID *string    `db:"fhir_resource_id"`
	CreatedBy      *uuid.UUID `db:"created_by"`
	IsActive       bool       `db:"is_active"`
	CreatedAt      time.Time  `db:"created_at"`
	UpdatedAt      time.Time  `db:"updated_at"`
}
