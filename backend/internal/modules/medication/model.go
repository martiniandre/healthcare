package medication

import (
	"errors"
	"time"
)

var ErrMedicationRequestNotFound = errors.New("medication request not found")

type Medication struct {
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
