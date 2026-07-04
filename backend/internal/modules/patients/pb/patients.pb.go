package pb

import "context"

type PatientServiceServer interface {
	CreatePatient(ctx context.Context, req *CreatePatientRequest) (*CreatePatientResponse, error)
	GetPatient(ctx context.Context, req *GetPatientRequest) (*GetPatientResponse, error)
	GetPatientByDocument(ctx context.Context, req *GetPatientByDocumentRequest) (*GetPatientResponse, error)
	ListPatients(ctx context.Context, req *ListPatientsRequest) (*ListPatientsResponse, error)
}

type CreatePatientRequest struct {
	FullName    string
	BirthDate   string
	DocumentID  string
	PhoneNumber string
}

type CreatePatientResponse struct {
	PatientId      string
	FhirResourceId string
}

type GetPatientRequest struct {
	FhirResourceId string
}

type GetPatientByDocumentRequest struct {
	DocumentId string
}

type GetPatientResponse struct {
	PatientId      string
	FhirResourceId string
	FullName       string
	BirthDate      string
	DocumentId     string
	PhoneNumber    string
}

type ListPatientsRequest struct {
	Search        string
	SortField     string
	SortDirection string
	Page          int32
	Limit         int32
}

type ListPatientsResponse struct {
	Patients []*GetPatientResponse
}

func RegisterPatientServiceServer(_ interface{}, server PatientServiceServer) {}
