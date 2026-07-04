package diagnostic_report

import (
	"errors"
	"time"
)

var ErrDiagnosticReportNotFound = errors.New("diagnostic report not found")

type DiagnosticReport struct {
	FHIRResourceID  string
	EncounterFHIRID string
	PatientFHIRID   string
	ReportCode      string
	ReportDisplay   string
	Status          string
	Conclusion      string
	IssuedAt        time.Time
}
