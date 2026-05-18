package pb

import (
	"context"

	"google.golang.org/grpc"
)

type ClinicalServiceServer interface {
	CreateEncounter(ctx context.Context, req *CreateEncounterRequest) (*CreateEncounterResponse, error)
	GetEncounters(ctx context.Context, req *GetEncountersRequest) (*GetEncountersResponse, error)

	CreateObservation(ctx context.Context, req *CreateObservationRequest) (*CreateObservationResponse, error)
	GetObservations(ctx context.Context, req *GetObservationsRequest) (*GetObservationsResponse, error)

	CreateCondition(ctx context.Context, req *CreateConditionRequest) (*CreateConditionResponse, error)
	GetConditions(ctx context.Context, req *GetConditionsRequest) (*GetConditionsResponse, error)

	CreateAllergyIntolerance(ctx context.Context, req *CreateAllergyIntoleranceRequest) (*CreateAllergyIntoleranceResponse, error)
	GetAllergyIntolerances(ctx context.Context, req *GetAllergyIntolerancesRequest) (*GetAllergyIntolerancesResponse, error)

	CreateMedicationRequest(ctx context.Context, req *CreateMedicationRequestRequest) (*CreateMedicationRequestResponse, error)
	GetMedicationRequests(ctx context.Context, req *GetMedicationRequestsRequest) (*GetMedicationRequestsResponse, error)

	CreateDiagnosticReport(ctx context.Context, req *CreateDiagnosticReportRequest) (*CreateDiagnosticReportResponse, error)
	GetDiagnosticReports(ctx context.Context, req *GetDiagnosticReportsRequest) (*GetDiagnosticReportsResponse, error)
}

type UnimplementedClinicalServiceServer struct{}

func (UnimplementedClinicalServiceServer) CreateEncounter(_ context.Context, _ *CreateEncounterRequest) (*CreateEncounterResponse, error) {
	return nil, nil
}
func (UnimplementedClinicalServiceServer) GetEncounters(_ context.Context, _ *GetEncountersRequest) (*GetEncountersResponse, error) {
	return nil, nil
}
func (UnimplementedClinicalServiceServer) CreateObservation(_ context.Context, _ *CreateObservationRequest) (*CreateObservationResponse, error) {
	return nil, nil
}
func (UnimplementedClinicalServiceServer) GetObservations(_ context.Context, _ *GetObservationsRequest) (*GetObservationsResponse, error) {
	return nil, nil
}
func (UnimplementedClinicalServiceServer) CreateCondition(_ context.Context, _ *CreateConditionRequest) (*CreateConditionResponse, error) {
	return nil, nil
}
func (UnimplementedClinicalServiceServer) GetConditions(_ context.Context, _ *GetConditionsRequest) (*GetConditionsResponse, error) {
	return nil, nil
}
func (UnimplementedClinicalServiceServer) CreateAllergyIntolerance(_ context.Context, _ *CreateAllergyIntoleranceRequest) (*CreateAllergyIntoleranceResponse, error) {
	return nil, nil
}
func (UnimplementedClinicalServiceServer) GetAllergyIntolerances(_ context.Context, _ *GetAllergyIntolerancesRequest) (*GetAllergyIntolerancesResponse, error) {
	return nil, nil
}
func (UnimplementedClinicalServiceServer) CreateMedicationRequest(_ context.Context, _ *CreateMedicationRequestRequest) (*CreateMedicationRequestResponse, error) {
	return nil, nil
}
func (UnimplementedClinicalServiceServer) GetMedicationRequests(_ context.Context, _ *GetMedicationRequestsRequest) (*GetMedicationRequestsResponse, error) {
	return nil, nil
}
func (UnimplementedClinicalServiceServer) CreateDiagnosticReport(_ context.Context, _ *CreateDiagnosticReportRequest) (*CreateDiagnosticReportResponse, error) {
	return nil, nil
}
func (UnimplementedClinicalServiceServer) GetDiagnosticReports(_ context.Context, _ *GetDiagnosticReportsRequest) (*GetDiagnosticReportsResponse, error) {
	return nil, nil
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

type CreateDiagnosticReportRequest struct {
	EncounterFhirId string
	PatientFhirId   string
	ReportCode      string
	ReportDisplay   string
	Conclusion      string
}

type CreateDiagnosticReportResponse struct {
	DiagnosticReportFhirId string
}

type GetDiagnosticReportsRequest struct {
	EncounterFhirId string
}

type DiagnosticReport struct {
	FhirId        string
	ReportDisplay string
	Status        string
	Conclusion    string
}

type GetDiagnosticReportsResponse struct {
	DiagnosticReports []*DiagnosticReport
}

func RegisterClinicalServiceServer(server *grpc.Server, srv ClinicalServiceServer) {}
