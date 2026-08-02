package tests

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/healthcare/backend/internal/modules/diagnostic_report"
	"github.com/healthcare/backend/internal/modules/diagnostic_report/mocks"
	"github.com/healthcare/backend/internal/shared/apperrors"
	"github.com/healthcare/backend/internal/shared/eventbus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockDiagnosticReportEventBus struct {
	PublishedEvents []eventbus.Event
}

func (mockBus *mockDiagnosticReportEventBus) Publish(ctx context.Context, event eventbus.Event) error {
	mockBus.PublishedEvents = append(mockBus.PublishedEvents, event)
	return nil
}

func (mockBus *mockDiagnosticReportEventBus) Subscribe(eventName string, handler eventbus.Handler) {}

func newTestDiagnosticReportService(mockRepository *mocks.MockDiagnosticReportRepository, mockVersionRepository *mocks.MockVersionRepository, mockEventBus *mockDiagnosticReportEventBus) diagnostic_report.Service {
	if mockVersionRepository == nil {
		mockVersionRepository = &mocks.MockVersionRepository{}
	}
	if mockEventBus == nil {
		mockEventBus = &mockDiagnosticReportEventBus{}
	}
	return diagnostic_report.NewService(mockRepository, mockVersionRepository, mockEventBus, nil)
}

func TestCreateDiagnosticReport_AppliesDefaultsAndDelegatesCompleteEntity(t *testing.T) {
	var capturedReport *diagnostic_report.DiagnosticReport
	mockRepository := &mocks.MockDiagnosticReportRepository{
		CreateDiagnosticReportFn: func(ctx context.Context, entity *diagnostic_report.DiagnosticReport) (*diagnostic_report.DiagnosticReport, error) {
			capturedReport = entity
			return entity, nil
		},
	}
	diagnosticReportService := newTestDiagnosticReportService(mockRepository, nil, &mockDiagnosticReportEventBus{})
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

func TestCreateDiagnosticReport_PublishesReportReadyEvent(testingInstance *testing.T) {
	mockRepository := &mocks.MockDiagnosticReportRepository{
		CreateDiagnosticReportFn: func(ctx context.Context, entity *diagnostic_report.DiagnosticReport) (*diagnostic_report.DiagnosticReport, error) {
			entity.FHIRResourceID = "report-789"
			entity.Version = "3"
			return entity, nil
		},
	}
	eventBus := &mockDiagnosticReportEventBus{}
	diagnosticReportService := diagnostic_report.NewService(mockRepository, &mocks.MockVersionRepository{}, eventBus, nil)

	createdReport, createErr := diagnosticReportService.CreateDiagnosticReport(context.Background(), diagnostic_report.CreateDiagnosticReportInput{
		EncounterFHIRID: "encounter-456",
		PatientFHIRID:   "patient-123",
		ReportCode:      "24323-8",
	})

	assert.NoError(testingInstance, createErr)
	assert.NotNil(testingInstance, createdReport)
	require.Len(testingInstance, eventBus.PublishedEvents, 1)
	publishedEvent := eventBus.PublishedEvents[0]
	assert.Equal(testingInstance, "report.ready", publishedEvent.Name)
	assert.Equal(testingInstance, "patient-123", publishedEvent.Data["patient_id"])
	assert.Equal(testingInstance, "report-789", publishedEvent.Data["report_id"])
	assert.Equal(testingInstance, "3", publishedEvent.Data["version"])
	assert.Equal(testingInstance, "diagnostic_report", publishedEvent.Data["resource_type"])
	assert.Equal(testingInstance, "report-789", publishedEvent.Data["resource_id"])
}

func TestCreateDiagnosticReport_PublishesNoEventWhenRepositoryFails(testingInstance *testing.T) {
	repositoryFailure := errors.New("repository failure")
	mockRepository := &mocks.MockDiagnosticReportRepository{
		CreateDiagnosticReportFn: func(ctx context.Context, entity *diagnostic_report.DiagnosticReport) (*diagnostic_report.DiagnosticReport, error) {
			return nil, repositoryFailure
		},
	}
	eventBus := &mockDiagnosticReportEventBus{}
	diagnosticReportService := diagnostic_report.NewService(mockRepository, &mocks.MockVersionRepository{}, eventBus, nil)

	createdReport, createErr := diagnosticReportService.CreateDiagnosticReport(context.Background(), diagnostic_report.CreateDiagnosticReportInput{
		EncounterFHIRID: "encounter-456",
		PatientFHIRID:   "patient-123",
		ReportCode:      "24323-8",
	})

	assert.Nil(testingInstance, createdReport)
	assert.ErrorIs(testingInstance, createErr, repositoryFailure)
	assert.Empty(testingInstance, eventBus.PublishedEvents)
}

func TestCreateDiagnosticReport_RejectsEmptyReportCodeWithoutCallingRepository(t *testing.T) {
	repositoryCalled := false
	mockRepository := &mocks.MockDiagnosticReportRepository{
		CreateDiagnosticReportFn: func(ctx context.Context, entity *diagnostic_report.DiagnosticReport) (*diagnostic_report.DiagnosticReport, error) {
			repositoryCalled = true
			return entity, nil
		},
	}
	diagnosticReportService := newTestDiagnosticReportService(mockRepository, nil, &mockDiagnosticReportEventBus{})

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
	diagnosticReportService := newTestDiagnosticReportService(mockRepository, nil, &mockDiagnosticReportEventBus{})

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
			diagnosticReportService := newTestDiagnosticReportService(mockRepository, nil, &mockDiagnosticReportEventBus{})

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
	diagnosticReportService := newTestDiagnosticReportService(mockRepository, nil, &mockDiagnosticReportEventBus{})

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
	diagnosticReportService := newTestDiagnosticReportService(mockRepository, nil, &mockDiagnosticReportEventBus{})

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
	diagnosticReportService := newTestDiagnosticReportService(mockRepository, nil, &mockDiagnosticReportEventBus{})

	reports, getErr := diagnosticReportService.GetDiagnosticReportsByEncounter(context.Background(), "encounter-456")

	assert.Nil(t, reports)
	assert.ErrorIs(t, getErr, repositoryFailure)
}

func TestCreateDiagnosticReport_RecordsInitialVersionSnapshot(testingInstance *testing.T) {
	mockVersionRepository := &mocks.MockVersionRepository{}
	mockRepository := &mocks.MockDiagnosticReportRepository{
		CreateDiagnosticReportFn: func(ctx context.Context, entity *diagnostic_report.DiagnosticReport) (*diagnostic_report.DiagnosticReport, error) {
			entity.FHIRResourceID = "report-001"
			entity.Version = "1"
			return entity, nil
		},
	}
	diagnosticReportService := newTestDiagnosticReportService(mockRepository, mockVersionRepository, nil)

	createdReport, createErr := diagnosticReportService.CreateDiagnosticReport(context.Background(), diagnostic_report.CreateDiagnosticReportInput{
		EncounterFHIRID: "encounter-456",
		PatientFHIRID:   "patient-123",
		ReportCode:      "24323-8",
	})

	assert.NoError(testingInstance, createErr)
	assert.Equal(testingInstance, "report-001", createdReport.FHIRResourceID)
	require.Len(testingInstance, mockVersionRepository.RecordedReports, 1)
	assert.Equal(testingInstance, "report-001", mockVersionRepository.RecordedReports[0])
	assert.Equal(testingInstance, "1", mockVersionRepository.RecordedVersion)
	require.NotEmpty(testingInstance, mockVersionRepository.Snapshots[0])
}

func TestUpdateDiagnosticReport_MergesInputRecordsVersionAndPublishesEvent(testingInstance *testing.T) {
	conclusionUpdate := "Updated conclusion"
	statusUpdate := "amended"
	mockVersionRepository := &mocks.MockVersionRepository{}
	eventBus := &mockDiagnosticReportEventBus{}
	mockRepository := &mocks.MockDiagnosticReportRepository{
		GetDiagnosticReportByIDFn: func(ctx context.Context, reportFHIRID string) (*diagnostic_report.DiagnosticReport, error) {
			return &diagnostic_report.DiagnosticReport{
				FHIRResourceID:  reportFHIRID,
				EncounterFHIRID: "encounter-456",
				PatientFHIRID:   "patient-123",
				ReportCode:      "24323-8",
				ReportDisplay:   "Complete blood count",
				Conclusion:      "Original conclusion",
				Status:          "final",
				Version:         "1",
			}, nil
		},
		UpdateDiagnosticReportFn: func(ctx context.Context, reportFHIRID string, entity *diagnostic_report.DiagnosticReport) (*diagnostic_report.DiagnosticReport, error) {
			entity.FHIRResourceID = reportFHIRID
			entity.Version = "2"
			return entity, nil
		},
	}
	diagnosticReportService := diagnostic_report.NewService(mockRepository, mockVersionRepository, eventBus, nil)

	updatedReport, updateErr := diagnosticReportService.UpdateDiagnosticReport(context.Background(), "report-001", diagnostic_report.UpdateDiagnosticReportInput{
		Conclusion: &conclusionUpdate,
		Status:     &statusUpdate,
	})

	assert.NoError(testingInstance, updateErr)
	assert.Equal(testingInstance, "Updated conclusion", updatedReport.Conclusion)
	assert.Equal(testingInstance, "amended", updatedReport.Status)
	assert.Equal(testingInstance, "2", updatedReport.Version)
	assert.Equal(testingInstance, "24323-8", updatedReport.ReportCode)
	require.Len(testingInstance, mockVersionRepository.RecordedReports, 1)
	assert.Equal(testingInstance, "2", mockVersionRepository.RecordedVersion)
	require.Len(testingInstance, eventBus.PublishedEvents, 1)
	assert.Equal(testingInstance, "report.ready", eventBus.PublishedEvents[0].Name)
}

func TestUpdateDiagnosticReport_RejectsInvalidLOINC(testingInstance *testing.T) {
	invalidCode := "not-a-loinc"
	mockRepository := &mocks.MockDiagnosticReportRepository{
		GetDiagnosticReportByIDFn: func(ctx context.Context, reportFHIRID string) (*diagnostic_report.DiagnosticReport, error) {
			return &diagnostic_report.DiagnosticReport{FHIRResourceID: reportFHIRID, ReportCode: "24323-8", Status: "final"}, nil
		},
	}
	diagnosticReportService := newTestDiagnosticReportService(mockRepository, nil, nil)

	updatedReport, updateErr := diagnosticReportService.UpdateDiagnosticReport(context.Background(), "report-001", diagnostic_report.UpdateDiagnosticReportInput{
		ReportCode: &invalidCode,
	})

	assert.Nil(testingInstance, updatedReport)
	assert.Error(testingInstance, updateErr)
}

func TestUpdateDiagnosticReport_PropagatesNotFound(testingInstance *testing.T) {
	repositoryFailure := apperrors.ErrDiagnosticReportNotFound
	mockRepository := &mocks.MockDiagnosticReportRepository{
		GetDiagnosticReportByIDFn: func(ctx context.Context, reportFHIRID string) (*diagnostic_report.DiagnosticReport, error) {
			return nil, repositoryFailure
		},
	}
	diagnosticReportService := newTestDiagnosticReportService(mockRepository, nil, nil)

	updatedReport, updateErr := diagnosticReportService.UpdateDiagnosticReport(context.Background(), "report-missing", diagnostic_report.UpdateDiagnosticReportInput{})

	assert.Nil(testingInstance, updatedReport)
	assert.ErrorIs(testingInstance, updateErr, repositoryFailure)
}

func TestGetDiagnosticReportVersions_ReturnsHistoryFromRepository(testingInstance *testing.T) {
	expectedVersions := []*diagnostic_report.DiagnosticReportVersion{
		{ReportID: "report-001", Version: "2"},
		{ReportID: "report-001", Version: "1"},
	}
	mockVersionRepository := &mocks.MockVersionRepository{
		ListVersionsFn: func(ctx context.Context, reportID string) ([]*diagnostic_report.DiagnosticReportVersion, error) {
			return expectedVersions, nil
		},
	}
	diagnosticReportService := newTestDiagnosticReportService(&mocks.MockDiagnosticReportRepository{}, mockVersionRepository, nil)

	versions, listErr := diagnosticReportService.GetDiagnosticReportVersions(context.Background(), "report-001")

	assert.NoError(testingInstance, listErr)
	assert.Equal(testingInstance, expectedVersions, versions)
}

func TestGetDiagnosticReportVersion_ReturnsSpecificVersionFromRepository(testingInstance *testing.T) {
	expectedVersion := &diagnostic_report.DiagnosticReportVersion{ReportID: "report-001", Version: "1"}
	mockVersionRepository := &mocks.MockVersionRepository{
		GetVersionFn: func(ctx context.Context, reportID string, version string) (*diagnostic_report.DiagnosticReportVersion, error) {
			return expectedVersion, nil
		},
	}
	diagnosticReportService := newTestDiagnosticReportService(&mocks.MockDiagnosticReportRepository{}, mockVersionRepository, nil)

	versionEntry, getErr := diagnosticReportService.GetDiagnosticReportVersion(context.Background(), "report-001", "1")

	assert.NoError(testingInstance, getErr)
	assert.Equal(testingInstance, expectedVersion, versionEntry)
}

