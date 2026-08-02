package mocks

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/healthcare/backend/internal/modules/diagnostic_report"
)

type MockDiagnosticReportRepository struct {
	CreateDiagnosticReportFn          func(ctx context.Context, entity *diagnostic_report.DiagnosticReport) (*diagnostic_report.DiagnosticReport, error)
	GetDiagnosticReportByIDFn         func(ctx context.Context, reportFHIRID string) (*diagnostic_report.DiagnosticReport, error)
	UpdateDiagnosticReportFn          func(ctx context.Context, reportFHIRID string, entity *diagnostic_report.DiagnosticReport) (*diagnostic_report.DiagnosticReport, error)
	GetDiagnosticReportsByEncounterFn func(ctx context.Context, encounterFHIRID string) ([]*diagnostic_report.DiagnosticReport, error)
}

func (mockRepo *MockDiagnosticReportRepository) CreateDiagnosticReport(ctx context.Context, entity *diagnostic_report.DiagnosticReport) (*diagnostic_report.DiagnosticReport, error) {
	if mockRepo.CreateDiagnosticReportFn != nil {
		return mockRepo.CreateDiagnosticReportFn(ctx, entity)
	}
	return entity, nil
}

func (mockRepo *MockDiagnosticReportRepository) GetDiagnosticReportByID(ctx context.Context, reportFHIRID string) (*diagnostic_report.DiagnosticReport, error) {
	if mockRepo.GetDiagnosticReportByIDFn != nil {
		return mockRepo.GetDiagnosticReportByIDFn(ctx, reportFHIRID)
	}
	return &diagnostic_report.DiagnosticReport{FHIRResourceID: reportFHIRID}, nil
}

func (mockRepo *MockDiagnosticReportRepository) UpdateDiagnosticReport(ctx context.Context, reportFHIRID string, entity *diagnostic_report.DiagnosticReport) (*diagnostic_report.DiagnosticReport, error) {
	if mockRepo.UpdateDiagnosticReportFn != nil {
		return mockRepo.UpdateDiagnosticReportFn(ctx, reportFHIRID, entity)
	}
	return entity, nil
}

func (mockRepo *MockDiagnosticReportRepository) GetDiagnosticReportsByEncounter(ctx context.Context, encounterFHIRID string) ([]*diagnostic_report.DiagnosticReport, error) {
	if mockRepo.GetDiagnosticReportsByEncounterFn != nil {
		return mockRepo.GetDiagnosticReportsByEncounterFn(ctx, encounterFHIRID)
	}
	return []*diagnostic_report.DiagnosticReport{}, nil
}

type MockVersionRepository struct {
	RecordVersionFn func(ctx context.Context, reportID string, version string, snapshot json.RawMessage, changedBy *uuid.UUID) (*diagnostic_report.DiagnosticReportVersion, error)
	ListVersionsFn  func(ctx context.Context, reportID string) ([]*diagnostic_report.DiagnosticReportVersion, error)
	GetVersionFn    func(ctx context.Context, reportID string, version string) (*diagnostic_report.DiagnosticReportVersion, error)
	RecordedReports []string
	RecordedVersion string
	Snapshots       []json.RawMessage
}

func (mockRepo *MockVersionRepository) RecordVersion(ctx context.Context, reportID string, version string, snapshot json.RawMessage, changedBy *uuid.UUID) (*diagnostic_report.DiagnosticReportVersion, error) {
	if mockRepo.RecordVersionFn != nil {
		return mockRepo.RecordVersionFn(ctx, reportID, version, snapshot, changedBy)
	}
	mockRepo.RecordedReports = append(mockRepo.RecordedReports, reportID)
	mockRepo.RecordedVersion = version
	mockRepo.Snapshots = append(mockRepo.Snapshots, snapshot)
	return &diagnostic_report.DiagnosticReportVersion{ReportID: reportID, Version: version, Snapshot: snapshot}, nil
}

func (mockRepo *MockVersionRepository) ListVersions(ctx context.Context, reportID string) ([]*diagnostic_report.DiagnosticReportVersion, error) {
	if mockRepo.ListVersionsFn != nil {
		return mockRepo.ListVersionsFn(ctx, reportID)
	}
	return []*diagnostic_report.DiagnosticReportVersion{}, nil
}

func (mockRepo *MockVersionRepository) GetVersion(ctx context.Context, reportID string, version string) (*diagnostic_report.DiagnosticReportVersion, error) {
	if mockRepo.GetVersionFn != nil {
		return mockRepo.GetVersionFn(ctx, reportID, version)
	}
	return &diagnostic_report.DiagnosticReportVersion{ReportID: reportID, Version: version}, nil
}
