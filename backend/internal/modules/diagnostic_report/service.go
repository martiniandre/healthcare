package diagnostic_report

import (
	"context"
	"time"

	"github.com/healthcare/backend/internal/shared/apperrors"
	"github.com/healthcare/backend/internal/shared/eventbus"
	"github.com/healthcare/backend/internal/shared/validator"
)

type Service interface {
	CreateDiagnosticReport(ctx context.Context, input CreateDiagnosticReportInput) (*DiagnosticReport, error)
	GetDiagnosticReportsByEncounter(ctx context.Context, encounterFHIRID string) ([]*DiagnosticReport, error)
}

type service struct {
	repo     Repository
	eventBus eventbus.Bus
}

func NewService(repo Repository, eventBus eventbus.Bus) Service {
	return &service{repo: repo, eventBus: eventBus}
}

func (reportService *service) CreateDiagnosticReport(ctx context.Context, input CreateDiagnosticReportInput) (*DiagnosticReport, error) {
	fieldViolations := make(map[string]string)
	if input.EncounterFHIRID == "" {
		fieldViolations["encounter_fhir_id"] = "is required"
	}
	if input.PatientFHIRID == "" {
		fieldViolations["patient_fhir_id"] = "is required"
	}
	if input.ReportCode == "" {
		fieldViolations["report_code"] = "is required"
	} else if !validator.IsValidLOINC(input.ReportCode) {
		fieldViolations["report_code"] = "must be a valid LOINC code"
	}
	if len(fieldViolations) > 0 {
		return nil, apperrors.InvalidArgument("invalid diagnostic report input", fieldViolations)
	}

	report := &DiagnosticReport{
		EncounterFHIRID: input.EncounterFHIRID,
		PatientFHIRID:   input.PatientFHIRID,
		ReportCode:      input.ReportCode,
		ReportDisplay:   input.ReportDisplay,
		Conclusion:      input.Conclusion,
	}
	if report.Status == "" {
		report.Status = "final"
	}
	if report.IssuedAt.IsZero() {
		report.IssuedAt = time.Now()
	}
	createdReport, createErr := reportService.repo.CreateDiagnosticReport(ctx, report)
	if createErr != nil {
		return nil, createErr
	}

	if reportService.eventBus != nil {
		reportService.eventBus.Publish(ctx, eventbus.Event{
			Name: "report.ready",
			Data: map[string]any{
				"patient_id":    createdReport.PatientFHIRID,
				"report_id":     createdReport.FHIRResourceID,
				"version":       createdReport.Version,
				"title":         "Laudo Disponível",
				"body":          "O laudo do paciente " + createdReport.PatientFHIRID + " está pronto para consulta.",
				"resource_type": "diagnostic_report",
				"resource_id":   createdReport.FHIRResourceID,
			},
		})
	}

	return createdReport, nil
}

func (reportService *service) GetDiagnosticReportsByEncounter(ctx context.Context, encounterFHIRID string) ([]*DiagnosticReport, error) {
	return reportService.repo.GetDiagnosticReportsByEncounter(ctx, encounterFHIRID)
}
