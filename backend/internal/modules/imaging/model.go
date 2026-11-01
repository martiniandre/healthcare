package imaging

import (
	"time"

	"github.com/google/uuid"
)

type ImagingStudy struct {
	ID               uuid.UUID  `db:"id"`
	PatientFhirID    string     `db:"patient_fhir_id"`
	Title            string     `db:"title"`
	Modality         string     `db:"modality"`
	FileName         string     `db:"file_name"`
	StudyInstanceUID string     `db:"study_instance_uid"`
	StudyDateTime    *time.Time `db:"study_datetime"`
	AccessionNumber  string     `db:"accession_number"`
	StudyDescription string     `db:"study_description"`
	Status           string     `db:"status"`
	CreatedAt        time.Time  `db:"created_at"`
	UpdatedAt        time.Time  `db:"updated_at"`
}

type InstanceIndex struct {
	SeriesInstanceUID  string
	SeriesNumber       int
	SeriesDescription  string
	Modality           string
	SOPInstanceUID     string
	SOPClassUID        string
	TransferSyntaxUID  string
	InstanceNumber     int
	FrameCount         int
	Rows               int
	Columns            int
	PixelSpacingRow    float64
	PixelSpacingColumn float64
}
