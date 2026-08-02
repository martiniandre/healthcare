package tests

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/healthcare/backend/internal/modules/diagnostic_report"
	"github.com/healthcare/backend/internal/modules/diagnostic_report/mocks"
	"github.com/stretchr/testify/assert"
)

func TestCreateDiagnosticReport_AppliesDefaultsAndDelegatesCompleteEntity(t *testing.T) {
	var capturedReport *diagnostic_report.DiagnosticReport
	mockRepository := &mocks.MockDiagnosticReportRepository{
		CreateDiagnosticReportFn: func(ctx context.Context, entity *diagnostic_report.DiagnosticReport) (*diagnostic_report.DiagnosticReport, error) {
			capturedReport = entity
			return entity, nil
		},
	}
	diagnosticReportService := diagnostic_report.NewService(mockRepository)
	beforeCall := time.Now()

	createdReport, createErr := diagnosticReportService.CreateDiagnosticReport(context.Background(), diagnostic_report.CreateDiagnosticReportInput{
		EncounterFHIRID: "encounter-456",
		PatientFHIRID:   "patient-123",
		ReportCode:      "24323-8",
		ReportDisplay:   "Complete blood count",
		Conclusion:      "Normal values",
	})

	assert.NoError(t, createErr)
	assert.NotNil(t, createdReport)
	assert.Equal(t, "final", createdReport.Status)
	assert.False(t, createdReport.IssuedAt.IsZero())
	assert.False(t, createdReport.IssuedAt.Before(beforeCall))
	assert.WithinDuration(t, time.Now(), createdReport.IssuedAt, time.Minute)
	assert.NotNil(t, capturedReport)
	assert.Equal(t, "encounter-456", capturedReport.EncounterFHIRID)
	assert.Equal(t, "patient-123", capturedReport.PatientFHIRID)
	assert.Equal(t, "24323-8", capturedReport.ReportCode)
	assert.Equal(t, "Complete blood count", capturedReport.ReportDisplay)
	assert.Equal(t, "Normal values", capturedReport.Conclusion)
	assert.Equal(t, "final", capturedReport.Status)
	assert.False(t, capturedReport.IssuedAt.IsZero())
}

func TestCreateDiagnosticReport_RejectsEmptyReportCodeWithoutCallingRepository(t *testing.T) {
	repositoryCalled := false
	mockRepository := &mocks.MockDiagnosticReportRepository{
		CreateDiagnosticReportFn: func(ctx context.Context, entity *diagnostic_report.DiagnosticReport) (*diagnostic_report.DiagnosticReport, error) {
			repositoryCalled = true
			return entity, nil
		},
	}
	diagnosticReportService := diagnostic_report.NewService(mockRepository)

	createdReport, createErr := diagnosticReportService.CreateDiagnosticReport(context.Background(), diagnostic_report.CreateDiagnosticReportInput{
		EncounterFHIRID: "encounter-456",
		PatientFHIRID:   "patient-123",
	})

	assert.Nil(t, createdReport)
	assert.Error(t, createErr)
	assert.False(t, repositoryCalled)
}

func TestCreateDiagnosticReport_RejectsInvalidLOINCWithoutCallingRepository(t *testing.T) {
	repositoryCalled := false
	mockRepository := &mocks.MockDiagnosticReportRepository{
		CreateDiagnosticReportFn: func(ctx context.Context, entity *diagnostic_report.DiagnosticReport) (*diagnostic_report.DiagnosticReport, error) {
			repositoryCalled = true
			return entity, nil
		},
	}
	diagnosticReportService := diagnostic_report.NewService(mockRepository)

	createdReport, createErr := diagnosticReportService.CreateDiagnosticReport(context.Background(), diagnostic_report.CreateDiagnosticReportInput{
		EncounterFHIRID: "encounter-456",
		PatientFHIRID:   "patient-123",
		ReportCode:      "not-a-loinc",
	})

	assert.Nil(t, createdReport)
	assert.Error(t, createErr)
	assert.False(t, repositoryCalled)
}

func TestCreateDiagnosticReport_RejectsMissingPatientOrEncounter(t *testing.T) {
	testCases := []struct {
		name  string
		input diagnostic_report.CreateDiagnosticReportInput
	}{
		{
			name: "missing patient",
			input: diagnostic_report.CreateDiagnosticReportInput{
				EncounterFHIRID: "encounter-456",
				ReportCode:      "24323-8",
			},
		},
		{
			name: "missing encounter",
			input: diagnostic_report.CreateDiagnosticReportInput{
				PatientFHIRID: "patient-123",
				ReportCode:    "24323-8",
			},
		},
		{
			name: "missing both",
			input: diagnostic_report.CreateDiagnosticReportInput{
				ReportCode: "24323-8",
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			repositoryCalled := false
			mockRepository := &mocks.MockDiagnosticReportRepository{
				CreateDiagnosticReportFn: func(ctx context.Context, entity *diagnostic_report.DiagnosticReport) (*diagnostic_report.DiagnosticReport, error) {
					repositoryCalled = true
					return entity, nil
				},
			}
			diagnosticReportService := diagnostic_report.NewService(mockRepository)

			createdReport, createErr := diagnosticReportService.CreateDiagnosticReport(context.Background(), testCase.input)

			assert.Nil(t, createdReport)
			assert.Error(t, createErr)
			assert.False(t, repositoryCalled)
		})
	}
}

func TestCreateDiagnosticReport_PropagatesRepositoryError(t *testing.T) {
	repositoryFailure := errors.New("repository failure")
	mockRepository := &mocks.MockDiagnosticReportRepository{
		CreateDiagnosticReportFn: func(ctx context.Context, entity *diagnostic_report.DiagnosticReport) (*diagnostic_report.DiagnosticReport, error) {
			return nil, repositoryFailure
		},
	}
	diagnosticReportService := diagnostic_report.NewService(mockRepository)

	createdReport, createErr := diagnosticReportService.CreateDiagnosticReport(context.Background(), diagnostic_report.CreateDiagnosticReportInput{
		EncounterFHIRID: "encounter-456",
		PatientFHIRID:   "patient-123",
		ReportCode:      "24323-8",
	})

	assert.Nil(t, createdReport)
	assert.ErrorIs(t, createErr, repositoryFailure)
}

func TestGetDiagnosticReportsByEncounter_ReturnsReportsFromRepository(t *testing.T) {
	expectedReports := []*diagnostic_report.DiagnosticReport{
		{FHIRResourceID: "report-1", EncounterFHIRID: "encounter-456"},
		{FHIRResourceID: "report-2", EncounterFHIRID: "encounter-456"},
	}
	mockRepository := &mocks.MockDiagnosticReportRepository{
		GetDiagnosticReportsByEncounterFn: func(ctx context.Context, encounterFHIRID string) ([]*diagnostic_report.DiagnosticReport, error) {
			return expectedReports, nil
		},
	}
	diagnosticReportService := diagnostic_report.NewService(mockRepository)

	reports, getErr := diagnosticReportService.GetDiagnosticReportsByEncounter(context.Background(), "encounter-456")

	assert.NoError(t, getErr)
	assert.Equal(t, expectedReports, reports)
}

func TestGetDiagnosticReportsByEncounter_PropagatesRepositoryError(t *testing.T) {
	repositoryFailure := errors.New("repository failure")
	mockRepository := &mocks.MockDiagnosticReportRepository{
		GetDiagnosticReportsByEncounterFn: func(ctx context.Context, encounterFHIRID string) ([]*diagnostic_report.DiagnosticReport, error) {
			return nil, repositoryFailure
		},
	}
	diagnosticReportService := diagnostic_report.NewService(mockRepository)

	reports, getErr := diagnosticReportService.GetDiagnosticReportsByEncounter(context.Background(), "encounter-456")

	assert.Nil(t, reports)
	assert.ErrorIs(t, getErr, repositoryFailure)
}
