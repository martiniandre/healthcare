package allergy

import "time"

type Allergy struct {
	FHIRResourceID  string
	PatientFHIRID   string
	AllergenCode    string
	AllergenDisplay string
	ClinicalStatus  string
	Reaction        string
	RecordedAt      time.Time
}
