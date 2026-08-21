package pb

import "context"

type ObservationServiceServer interface {
	CreateObservation(ctx context.Context, req *CreateObservationRequest) (*CreateObservationResponse, error)
	CreateObservationBatch(ctx context.Context, req *CreateObservationBatchRequest) (*CreateObservationBatchResponse, error)
	GetObservations(ctx context.Context, req *GetObservationsRequest) (*GetObservationsResponse, error)
}

type CreateObservationRequest struct {
	EncounterFhirId string
	PatientFhirId   string
	LoincCode       string
	CodeDisplay     string
	ValueQuantity   float64
	ValueUnit       string
}

type CreateObservationResponse struct {
	ObservationFhirId string
}

type VitalSignsPanel struct {
	HeartRate              *float64
	BodyTemperature        *float64
	SystolicBloodPressure  *float64
	DiastolicBloodPressure *float64
	OxygenSaturation       *float64
	RespiratoryRate        *float64
	WeightKilograms        *float64
	HeightCentimeters      *float64
}

type CreateObservationBatchRequest struct {
	EncounterFhirId string
	PatientFhirId   string
	Panel           *VitalSignsPanel
}

type CreateObservationBatchResponse struct {
	Observations []*Observation
}

type GetObservationsRequest struct {
	EncounterFhirId string
}

type Observation struct {
	FhirId        string
	LoincCode     string
	CodeDisplay   string
	ValueQuantity float64
	ValueUnit     string
	NotPerformed  bool
}

type GetObservationsResponse struct {
	Observations []*Observation
}

func RegisterObservationServiceServer(_ interface{}, server ObservationServiceServer) {}
