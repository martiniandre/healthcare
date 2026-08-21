package observation

import (
	"time"
)

type Observation struct {
	FHIRResourceID  string
	EncounterFHIRID string
	PatientFHIRID   string
	LoincCode       string
	CodeDisplay     string
	ValueQuantity   float64
	ValueUnit       string
	NotPerformed    bool
	ObservedAt      time.Time
}
