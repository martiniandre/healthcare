package diagnostic_report

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

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

type ReportSnapshot struct {
	FHIRResourceID  string    `json:"fhir_id"`
	EncounterFHIRID string    `json:"encounter_fhir_id"`
	PatientFHIRID   string    `json:"patient_fhir_id"`
	ReportCode      string    `json:"report_code"`
	ReportDisplay   string    `json:"report_display"`
	Status          string    `json:"status"`
	Conclusion      string    `json:"conclusion"`
	IssuedAt        time.Time `json:"issued_at"`
	Version         string    `json:"version"`
}

func NewReportSnapshot(report *DiagnosticReport) (json.RawMessage, error) {
	return json.Marshal(ReportSnapshot{
		FHIRResourceID:  report.FHIRResourceID,
		EncounterFHIRID: report.EncounterFHIRID,
		PatientFHIRID:   report.PatientFHIRID,
		ReportCode:      report.ReportCode,
		ReportDisplay:   report.ReportDisplay,
		Status:          report.Status,
		Conclusion:      report.Conclusion,
		IssuedAt:        report.IssuedAt,
		Version:         report.Version,
	})
}

type DiagnosticReportVersion struct {
	ID        uuid.UUID       `db:"id"`
	ReportID  string          `db:"report_id"`
	Version   string          `db:"version"`
	Snapshot  json.RawMessage `db:"snapshot"`
	ChangedBy *uuid.UUID      `db:"changed_by"`
	ChangedAt time.Time       `db:"changed_at"`
}
