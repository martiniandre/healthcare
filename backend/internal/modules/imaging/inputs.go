package imaging

type UploadDICOMInput struct {
	PatientFhirID string
	Title         string
	Modality      string
	FileName      string
}
