package pb

import "context"

type ConditionServiceServer interface {
	CreateCondition(ctx context.Context, req *CreateConditionRequest) (*CreateConditionResponse, error)
	GetConditions(ctx context.Context, req *GetConditionsRequest) (*GetConditionsResponse, error)
}

type CreateConditionRequest struct {
	EncounterFhirId string
	PatientFhirId   string
	Icd10Code       string
	CodeDisplay     string
	ClinicalStatus  string
}

type CreateConditionResponse struct {
	ConditionFhirId string
}

type GetConditionsRequest struct {
	PatientFhirId string
}

type Condition struct {
	FhirId         string
	Icd10Code      string
	CodeDisplay    string
	ClinicalStatus string
}

type GetConditionsResponse struct {
	Conditions []*Condition
}

func RegisterConditionServiceServer(_ interface{}, server ConditionServiceServer) {}
