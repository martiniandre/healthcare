package mocks

import (
	"context"

	"github.com/healthcare/backend/internal/modules/diagnostic_report"
)

type MockDiagnosticReportRepository struct {
	CreateDiagnosticReportFn          func(ctx context.Context, entity *diagnostic_report.DiagnosticReport) (*diagnostic_report.DiagnosticReport, error)
	GetDiagnosticReportsByEncounterFn func(ctx context.Context, encounterFHIRID string) ([]*diagnostic_report.DiagnosticReport, error)
}

func (mockRepo *MockDiagnosticReportRepository) CreateDiagnosticReport(ctx context.Context, entity *diagnostic_report.DiagnosticReport) (*diagnostic_report.DiagnosticReport, error) {
	if mockRepo.CreateDiagnosticReportFn != nil {
		return mockRepo.CreateDiagnosticReportFn(ctx, entity)
	}
	return entity, nil
}

func (mockRepo *MockDiagnosticReportRepository) GetDiagnosticReportsByEncounter(ctx context.Context, encounterFHIRID string) ([]*diagnostic_report.DiagnosticReport, error) {
	if mockRepo.GetDiagnosticReportsByEncounterFn != nil {
		return mockRepo.GetDiagnosticReportsByEncounterFn(ctx, encounterFHIRID)
	}
	return []*diagnostic_report.DiagnosticReport{}, nil
}
