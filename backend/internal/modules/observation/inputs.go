package observation

import "time"

type CreateObservationInput struct {
	EncounterFHIRID string
	PatientFHIRID   string
	LoincCode       string
	CodeDisplay     string
	ValueQuantity   float64
	ValueUnit       string
	ObservedAt      *time.Time
}
