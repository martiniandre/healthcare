package pb

import "context"

type EncounterServiceServer interface {
	CreateEncounter(ctx context.Context, req *CreateEncounterRequest) (*CreateEncounterResponse, error)
	GetEncounter(ctx context.Context, req *GetEncounterRequest) (*GetEncounterResponse, error)
	GetEncounters(ctx context.Context, req *GetEncountersRequest) (*GetEncountersResponse, error)
}

type CreateEncounterRequest struct {
	PatientFhirId  string
	PractitionerId string
	ReasonCode     string
	ReasonDisplay  string
}

type CreateEncounterResponse struct {
	EncounterFhirId string
}

type GetEncounterRequest struct {
	EncounterFhirId string
}

type GetEncounterResponse struct {
	FhirId        string
	PatientFhirId string
	Status        string
	ReasonDisplay string
}

type GetEncountersRequest struct {
	PatientFhirId string
}

type Encounter struct {
	FhirId        string
	PatientFhirId string
	Status        string
	ReasonDisplay string
}

type GetEncountersResponse struct {
	Encounters []*Encounter
}

func RegisterEncounterServiceServer(_ interface{}, server EncounterServiceServer) {}
