package diagnostic_report

import "time"

type DiagnosticReport struct {
	FHIRResourceID  string
	EncounterFHIRID string
	PatientFHIRID   string
	ReportCode      string
	ReportDisplay   string
	Status          string
	Conclusion      string
	IssuedAt        time.Time
	Version         string
}
