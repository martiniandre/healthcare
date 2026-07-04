package pb

import "context"

type MedicationServiceServer interface {
	CreateMedicationRequest(ctx context.Context, req *CreateMedicationRequestRequest) (*CreateMedicationRequestResponse, error)
	GetMedicationRequests(ctx context.Context, req *GetMedicationRequestsRequest) (*GetMedicationRequestsResponse, error)
}

type CreateMedicationRequestRequest struct {
	EncounterFhirId    string
	PatientFhirId      string
	PractitionerFhirId string
	MedicationCode     string
	MedicationName     string
	DosageInstructions string
}

type CreateMedicationRequestResponse struct {
	MedicationRequestFhirId string
}

type GetMedicationRequestsRequest struct {
	EncounterFhirId string
}

type MedicationRequest struct {
	FhirId             string
	MedicationName     string
	DosageInstructions string
	Status             string
}

type GetMedicationRequestsResponse struct {
	MedicationRequests []*MedicationRequest
}

func RegisterMedicationServiceServer(_ interface{}, server MedicationServiceServer) {}
