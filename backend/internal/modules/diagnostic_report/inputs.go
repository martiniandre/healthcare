package diagnostic_report

type CreateDiagnosticReportInput struct {
	EncounterFHIRID string
	PatientFHIRID   string
	ReportCode      string
	ReportDisplay   string
	Conclusion      string
}

type UpdateDiagnosticReportInput struct {
	ReportCode    *string
	ReportDisplay *string
	Conclusion    *string
	Status        *string
}
