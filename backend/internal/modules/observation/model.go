package observation

import (
	"errors"
	"time"
)

var ErrObservationNotFound = errors.New("observation not found")

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
