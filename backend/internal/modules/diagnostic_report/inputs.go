package diagnostic_report

type CreateDiagnosticReportInput struct {
	EncounterFHIRID string
	PatientFHIRID   string
	ReportCode      string
	ReportDisplay   string
	Conclusion      string
}
