package imaging

import (
	"time"

	"github.com/google/uuid"
)

type ImagingStudy struct {
	ID               uuid.UUID `db:"id"`
	PatientFhirID    string    `db:"patient_fhir_id"`
	Title            string    `db:"title"`
	Modality         string    `db:"modality"`
	FileName         string    `db:"file_name"`
	StudyInstanceUID string    `db:"study_instance_uid"`
	Status           string    `db:"status"`
	CreatedAt        time.Time `db:"created_at"`
	UpdatedAt        time.Time `db:"updated_at"`
}
