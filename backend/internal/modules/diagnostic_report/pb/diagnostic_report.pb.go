package pb

import "context"

type DiagnosticReportServiceServer interface {
	CreateDiagnosticReport(ctx context.Context, req *CreateDiagnosticReportRequest) (*CreateDiagnosticReportResponse, error)
	GetDiagnosticReports(ctx context.Context, req *GetDiagnosticReportsRequest) (*GetDiagnosticReportsResponse, error)
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

func RegisterDiagnosticReportServiceServer(_ interface{}, server DiagnosticReportServiceServer) {}
