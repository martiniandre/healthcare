package pb

import "context"

type AllergyServiceServer interface {
	CreateAllergyIntolerance(ctx context.Context, req *CreateAllergyIntoleranceRequest) (*CreateAllergyIntoleranceResponse, error)
	GetAllergyIntolerances(ctx context.Context, req *GetAllergyIntolerancesRequest) (*GetAllergyIntolerancesResponse, error)
}

type CreateAllergyIntoleranceRequest struct {
	PatientFhirId   string
	AllergenCode    string
	AllergenDisplay string
	ClinicalStatus  string
	Reaction        string
}

type CreateAllergyIntoleranceResponse struct {
	AllergyFhirId string
}

type GetAllergyIntolerancesRequest struct {
	PatientFhirId string
}

type AllergyIntolerance struct {
	FhirId          string
	AllergenDisplay string
	ClinicalStatus  string
	Reaction        string
}

type GetAllergyIntolerancesResponse struct {
	Allergies []*AllergyIntolerance
}

func RegisterAllergyServiceServer(_ interface{}, server AllergyServiceServer) {}
