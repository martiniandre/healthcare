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
	GCSPath          string    `db:"gcs_path"`
	StudyInstanceUID string    `db:"study_instance_uid"`
	Status           string    `db:"status"`
	CreatedAt        time.Time `db:"created_at"`
	UpdatedAt        time.Time `db:"updated_at"`
}
