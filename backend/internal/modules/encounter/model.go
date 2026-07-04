package encounter

import (
	"errors"
	"time"
)

var ErrEncounterNotFound = errors.New("encounter not found")

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
