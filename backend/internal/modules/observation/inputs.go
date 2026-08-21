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

type CreateObservationBatchInput struct {
	EncounterFHIRID        string
	PatientFHIRID          string
	HeartRate              *float64
	BodyTemperature        *float64
	SystolicBloodPressure  *float64
	DiastolicBloodPressure *float64
	OxygenSaturation       *float64
	RespiratoryRate        *float64
	WeightKilograms        *float64
	HeightCentimeters      *float64
}
