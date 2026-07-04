package condition

import (
	"errors"
	"time"
)

var ErrConditionNotFound = errors.New("condition not found")

type Condition struct {
	FHIRResourceID  string
	EncounterFHIRID string
	PatientFHIRID   string
	ICD10Code       string
	CodeDisplay     string
	ClinicalStatus  string
	OnsetAt         time.Time
}
