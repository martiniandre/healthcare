package pb

import "context"

type ObservationServiceServer interface {
	CreateObservation(ctx context.Context, req *CreateObservationRequest) (*CreateObservationResponse, error)
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

type GetObservationsRequest struct {
	EncounterFhirId string
}

type Observation struct {
	FhirId        string
	LoincCode     string
	CodeDisplay   string
	ValueQuantity float64
	ValueUnit     string
}

type GetObservationsResponse struct {
	Observations []*Observation
}

func RegisterObservationServiceServer(_ interface{}, server ObservationServiceServer) {}
