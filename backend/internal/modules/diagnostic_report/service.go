package diagnostic_report

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/healthcare/backend/internal/shared/apperrors"
	"github.com/healthcare/backend/internal/shared/ctxkeys"
	"github.com/healthcare/backend/internal/shared/eventbus"
	"github.com/healthcare/backend/internal/shared/validator"
)

type Service interface {
	CreateDiagnosticReport(ctx context.Context, input CreateDiagnosticReportInput) (*DiagnosticReport, error)
	UpdateDiagnosticReport(ctx context.Context, reportFHIRID string, input UpdateDiagnosticReportInput) (*DiagnosticReport, error)
	GetDiagnosticReportsByEncounter(ctx context.Context, encounterFHIRID string) ([]*DiagnosticReport, error)
	GetDiagnosticReportVersions(ctx context.Context, reportFHIRID string) ([]*DiagnosticReportVersion, error)
	GetDiagnosticReportVersion(ctx context.Context, reportFHIRID string, version string) (*DiagnosticReportVersion, error)
}

type service struct {
	repo              Repository
	versionRepository VersionRepository
	eventBus          eventbus.Bus
}

func NewService(repo Repository, versionRepository VersionRepository, eventBus eventbus.Bus) Service {
	return &service{repo: repo, versionRepository: versionRepository, eventBus: eventBus}
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

	if createdReport.Version == "" {
		createdReport.Version = "1"
	}
	reportService.recordVersion(ctx, createdReport, changedByFromContext(ctx))

	reportService.publishReportReady(ctx, createdReport)

	return createdReport, nil
}

func (reportService *service) UpdateDiagnosticReport(ctx context.Context, reportFHIRID string, input UpdateDiagnosticReportInput) (*DiagnosticReport, error) {
	currentReport, fetchErr := reportService.repo.GetDiagnosticReportByID(ctx, reportFHIRID)
	if fetchErr != nil {
		return nil, fetchErr
	}

	if input.ReportCode != nil && *input.ReportCode != "" && !validator.IsValidLOINC(*input.ReportCode) {
		return nil, apperrors.InvalidArgument("invalid diagnostic report input", map[string]string{"report_code": "must be a valid LOINC code"})
	}

	mergedReport := mergeDiagnosticReportInput(currentReport, input)
	updatedReport, updateErr := reportService.repo.UpdateDiagnosticReport(ctx, reportFHIRID, mergedReport)
	if updateErr != nil {
		return nil, updateErr
	}

	if updatedReport.Version == "" {
		updatedReport.Version = "1"
	}
	reportService.recordVersion(ctx, updatedReport, changedByFromContext(ctx))

	reportService.publishReportReady(ctx, updatedReport)

	return updatedReport, nil
}

func (reportService *service) GetDiagnosticReportsByEncounter(ctx context.Context, encounterFHIRID string) ([]*DiagnosticReport, error) {
	return reportService.repo.GetDiagnosticReportsByEncounter(ctx, encounterFHIRID)
}

func (reportService *service) GetDiagnosticReportVersions(ctx context.Context, reportFHIRID string) ([]*DiagnosticReportVersion, error) {
	return reportService.versionRepository.ListVersions(ctx, reportFHIRID)
}

func (reportService *service) GetDiagnosticReportVersion(ctx context.Context, reportFHIRID string, version string) (*DiagnosticReportVersion, error) {
	return reportService.versionRepository.GetVersion(ctx, reportFHIRID, version)
}

func mergeDiagnosticReportInput(currentReport *DiagnosticReport, input UpdateDiagnosticReportInput) *DiagnosticReport {
	mergedReport := *currentReport
	if input.ReportCode != nil {
		mergedReport.ReportCode = *input.ReportCode
	}
	if input.ReportDisplay != nil {
		mergedReport.ReportDisplay = *input.ReportDisplay
	}
	if input.Conclusion != nil {
		mergedReport.Conclusion = *input.Conclusion
	}
	if input.Status != nil {
		mergedReport.Status = *input.Status
	}
	return &mergedReport
}

func (reportService *service) recordVersion(ctx context.Context, report *DiagnosticReport, changedBy *uuid.UUID) {
	if reportService.versionRepository == nil || report.FHIRResourceID == "" {
		return
	}
	snapshot, marshalErr := NewReportSnapshot(report)
	if marshalErr != nil {
		return
	}
	_, recordErr := reportService.versionRepository.RecordVersion(ctx, report.FHIRResourceID, report.Version, snapshot, changedBy)
	if recordErr != nil {
		return
	}
}

func (reportService *service) publishReportReady(ctx context.Context, report *DiagnosticReport) {
	if reportService.eventBus == nil {
		return
	}
	reportService.eventBus.Publish(ctx, eventbus.Event{
		Name: "report.ready",
		Data: map[string]any{
			"patient_id":    report.PatientFHIRID,
			"report_id":     report.FHIRResourceID,
			"version":       report.Version,
			"title":         "Laudo Disponível",
			"body":          "O laudo do paciente " + report.PatientFHIRID + " está pronto para consulta.",
			"resource_type": "diagnostic_report",
			"resource_id":   report.FHIRResourceID,
		},
	})
}

func changedByFromContext(ctx context.Context) *uuid.UUID {
	userIDString, exists := ctx.Value(ctxkeys.UserIDKey).(string)
	if !exists || userIDString == "" {
		return nil
	}
	parsedUserID, parseErr := uuid.Parse(userIDString)
	if parseErr != nil {
		return nil
	}
	return &parsedUserID
}
