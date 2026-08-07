package condition

import "time"

type Condition struct {
	FHIRResourceID  string
	EncounterFHIRID string
	PatientFHIRID   string
	ICD10Code       string
	CodeDisplay     string
	ClinicalStatus  string
	OnsetAt         time.Time
}
