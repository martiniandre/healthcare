package diagnostic_report

import (
	"context"
)

type Service interface {
	CreateDiagnosticReport(ctx context.Context, report *DiagnosticReport) (*DiagnosticReport, error)
	GetDiagnosticReportsByEncounter(ctx context.Context, encounterFHIRID string) ([]*DiagnosticReport, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (reportService *service) CreateDiagnosticReport(ctx context.Context, report *DiagnosticReport) (*DiagnosticReport, error) {
	if report.PatientFHIRID == "" || report.ReportCode == "" {
		return nil, ErrDiagnosticReportNotFound
	}
	if report.Status == "" {
		report.Status = "final"
	}
	return reportService.repo.CreateDiagnosticReport(ctx, report)
}

func (reportService *service) GetDiagnosticReportsByEncounter(ctx context.Context, encounterFHIRID string) ([]*DiagnosticReport, error) {
	return reportService.repo.GetDiagnosticReportsByEncounter(ctx, encounterFHIRID)
}
