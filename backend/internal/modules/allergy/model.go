package allergy

import (
	"errors"
	"time"
)

var ErrAllergyNotFound = errors.New("allergy intolerance not found")

type Allergy struct {
	FHIRResourceID  string
	PatientFHIRID   string
	AllergenCode    string
	AllergenDisplay string
	ClinicalStatus  string
	Reaction        string
	RecordedAt      time.Time
}
