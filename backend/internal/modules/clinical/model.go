package clinical

import (
	"errors"
	"time"
)

var (
	ErrEncounterNotFound         = errors.New("encounter not found")
	ErrObservationNotFound       = errors.New("observation not found")
	ErrConditionNotFound         = errors.New("condition not found")
	ErrAllergyNotFound           = errors.New("allergy intolerance not found")
	ErrMedicationRequestNotFound = errors.New("medication request not found")
	ErrDiagnosticReportNotFound  = errors.New("diagnostic report not found")
)

type Encounter struct {
	FHIRResourceID string
	PatientFHIRID  string
	PractitionerID string
	Status         string
	ReasonCode     string
	ReasonDisplay  string
	StartedAt      time.Time
	EndedAt        *time.Time
}

type Observation struct {
	FHIRResourceID  string
	EncounterFHIRID string
	PatientFHIRID   string
	LoincCode       string
	CodeDisplay     string
	ValueQuantity   float64
	ValueUnit       string
	ObservedAt      time.Time
}

type Condition struct {
	FHIRResourceID  string
	EncounterFHIRID string
	PatientFHIRID   string
	ICD10Code       string
	CodeDisplay     string
	ClinicalStatus  string
	OnsetAt         time.Time
}

type AllergyIntolerance struct {
	FHIRResourceID string
	PatientFHIRID  string
	AllergenCode   string
	AllergenDisplay string
	ClinicalStatus string
	Reaction       string
	RecordedAt     time.Time
}

type MedicationRequest struct {
	FHIRResourceID     string
	EncounterFHIRID    string
	PatientFHIRID      string
	PractitionerFHIRID string
	MedicationCode     string
	MedicationName     string
	DosageInstructions string
	Status             string
	IssuedAt           time.Time
}

type DiagnosticReport struct {
	FHIRResourceID  string
	EncounterFHIRID string
	PatientFHIRID   string
	ReportCode      string
	ReportDisplay   string
	Status          string
	Conclusion      string
	IssuedAt        time.Time
}
