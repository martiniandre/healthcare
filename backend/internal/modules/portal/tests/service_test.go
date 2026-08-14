package tests

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/healthcare/backend/internal/modules/portal"
	"github.com/healthcare/backend/internal/modules/portal/mocks"
	"github.com/stretchr/testify/assert"
)

const testPatientFHIRID = "patient-fhir-1"

func buildDashboardFixture() *mocks.MockPortalRepository {
	mockRepository := mocks.NewMockPortalRepository()
	mockRepository.Patients[testPatientFHIRID] = &portal.PatientInfo{
		FHIRResourceID: testPatientFHIRID,
		FullName:       "Mariana Costa",
		BirthDate:      "1990-01-15",
		DocumentID:     "529.982.247-25",
	}
	mockRepository.Encounters[testPatientFHIRID] = []portal.PortalEncounter{
		{FHIRResourceID: "enc-1", Status: "planned", ReasonDisplay: "Retorno"},
		{FHIRResourceID: "enc-2", Status: "arrived", ReasonDisplay: "Emergência"},
		{FHIRResourceID: "enc-3", Status: "finished", ReasonDisplay: "Rotina"},
	}
	mockRepository.Observations[testPatientFHIRID] = []portal.PortalObservation{
		{FHIRResourceID: "obs-1", LoincCode: "8867-4", CodeDisplay: "Frequência cardíaca", ValueQuantity: 72, ValueUnit: "bpm"},
	}
	mockRepository.Conditions[testPatientFHIRID] = []portal.PortalCondition{
		{FHIRResourceID: "cond-1", ICD10Code: "E11", CodeDisplay: "Diabetes", ClinicalStatus: "active"},
		{FHIRResourceID: "cond-2", ICD10Code: "J45", CodeDisplay: "Asma", ClinicalStatus: "resolved"},
	}
	mockRepository.Medications[testPatientFHIRID] = []portal.PortalMedication{
		{FHIRResourceID: "med-1", MedicationName: "Metformina", Status: "active"},
		{FHIRResourceID: "med-2", MedicationName: "Salbutamol", Status: "stopped"},
	}
	mockRepository.Reports[testPatientFHIRID] = []portal.PortalReport{
		{FHIRResourceID: "report-1", ReportDisplay: "Hemograma", Status: "final", Version: "2"},
	}
	mockRepository.Imaging[testPatientFHIRID] = []portal.PortalImaging{
		{FHIRResourceID: "img-1", Title: "Raio-X tórax", Modality: "CR", Status: "available"},
	}
	return mockRepository
}

func TestPortalService_GetDashboard_AggregatesClinicalSummary(testingInstance *testing.T) {
	mockRepository := buildDashboardFixture()
	portalService := portal.NewService(mockRepository)
	contextParam := context.Background()

	dashboard, err := portalService.GetDashboard(contextParam, testPatientFHIRID)

	assert.NoError(testingInstance, err)
	assert.NotNil(testingInstance, dashboard)
	assert.Equal(testingInstance, "Mariana Costa", dashboard.PatientInfo.FullName)
	assert.Len(testingInstance, dashboard.UpcomingEncounters, 2)
	assert.Len(testingInstance, dashboard.ActiveConditions, 1)
	assert.Equal(testingInstance, "E11", dashboard.ActiveConditions[0].ICD10Code)
	assert.Len(testingInstance, dashboard.ActiveMedications, 1)
	assert.Equal(testingInstance, "Metformina", dashboard.ActiveMedications[0].MedicationName)
	assert.Len(testingInstance, dashboard.RecentObservations, 1)
	assert.Len(testingInstance, dashboard.RecentReports, 1)
	assert.Len(testingInstance, dashboard.RecentImaging, 1)
}

func TestPortalService_GetDashboard_LimitsLargeCollections(testingInstance *testing.T) {
	mockRepository := mocks.NewMockPortalRepository()
	mockRepository.Patients[testPatientFHIRID] = &portal.PatientInfo{FHIRResourceID: testPatientFHIRID, FullName: "Paciente Volume"}
	manyObservations := make([]portal.PortalObservation, 25)
	for index := range manyObservations {
		manyObservations[index] = portal.PortalObservation{FHIRResourceID: string(rune('a' + index))}
	}
	manyReports := make([]portal.PortalReport, 15)
	for index := range manyReports {
		manyReports[index] = portal.PortalReport{FHIRResourceID: string(rune('a' + index))}
	}
	manyImaging := make([]portal.PortalImaging, 12)
	for index := range manyImaging {
		manyImaging[index] = portal.PortalImaging{FHIRResourceID: string(rune('a' + index))}
	}
	mockRepository.Observations[testPatientFHIRID] = manyObservations
	mockRepository.Reports[testPatientFHIRID] = manyReports
	mockRepository.Imaging[testPatientFHIRID] = manyImaging
	portalService := portal.NewService(mockRepository)
	contextParam := context.Background()

	dashboard, err := portalService.GetDashboard(contextParam, testPatientFHIRID)

	assert.NoError(testingInstance, err)
	assert.Len(testingInstance, dashboard.RecentObservations, 20)
	assert.Len(testingInstance, dashboard.RecentReports, 10)
	assert.Len(testingInstance, dashboard.RecentImaging, 10)
}

func TestPortalService_GetDashboard_PatientNotFound(testingInstance *testing.T) {
	mockRepository := buildDashboardFixture()
	portalService := portal.NewService(mockRepository)
	contextParam := context.Background()

	dashboard, err := portalService.GetDashboard(contextParam, "patient-inexistente")

	assert.Nil(testingInstance, dashboard)
	assert.ErrorIs(testingInstance, err, portal.ErrPatientNotFound)
}

func TestPortalService_GetDashboard_RepositoryFailureSurfacesInternalSections(testingInstance *testing.T) {
	mockRepository := buildDashboardFixture()
	mockRepository.ObservationsErr = errors.New("healthcare api unavailable")
	portalService := portal.NewService(mockRepository)
	contextParam := context.Background()

	dashboard, err := portalService.GetDashboard(contextParam, testPatientFHIRID)

	assert.NoError(testingInstance, err)
	assert.NotNil(testingInstance, dashboard)
	assert.Empty(testingInstance, dashboard.RecentObservations)
	assert.Len(testingInstance, dashboard.ActiveConditions, 1)
}

func TestPortalService_GetEncounters_PassesThroughRepository(testingInstance *testing.T) {
	mockRepository := buildDashboardFixture()
	portalService := portal.NewService(mockRepository)
	contextParam := context.Background()

	encounters, err := portalService.GetEncounters(contextParam, testPatientFHIRID)

	assert.NoError(testingInstance, err)
	assert.Len(testingInstance, encounters, 3)
}

func TestPortalService_GetObservations_PassesThroughRepository(testingInstance *testing.T) {
	mockRepository := buildDashboardFixture()
	portalService := portal.NewService(mockRepository)
	contextParam := context.Background()

	observations, err := portalService.GetObservations(contextParam, testPatientFHIRID)

	assert.NoError(testingInstance, err)
	assert.Len(testingInstance, observations, 1)
	assert.Equal(testingInstance, "8867-4", observations[0].LoincCode)
}

func TestPortalService_GetConditions_PassesThroughRepository(testingInstance *testing.T) {
	mockRepository := buildDashboardFixture()
	portalService := portal.NewService(mockRepository)
	contextParam := context.Background()

	conditions, err := portalService.GetConditions(contextParam, testPatientFHIRID)

	assert.NoError(testingInstance, err)
	assert.Len(testingInstance, conditions, 2)
}

func TestPortalService_GetMedications_PassesThroughRepository(testingInstance *testing.T) {
	mockRepository := buildDashboardFixture()
	portalService := portal.NewService(mockRepository)
	contextParam := context.Background()

	medications, err := portalService.GetMedications(contextParam, testPatientFHIRID)

	assert.NoError(testingInstance, err)
	assert.Len(testingInstance, medications, 2)
}

func TestPortalService_GetReports_PassesThroughRepository(testingInstance *testing.T) {
	mockRepository := buildDashboardFixture()
	portalService := portal.NewService(mockRepository)
	contextParam := context.Background()

	reports, err := portalService.GetReports(contextParam, testPatientFHIRID)

	assert.NoError(testingInstance, err)
	assert.Len(testingInstance, reports, 1)
	assert.Equal(testingInstance, "2", reports[0].Version)
}

func TestPortalService_GetImaging_PassesThroughRepository(testingInstance *testing.T) {
	mockRepository := buildDashboardFixture()
	portalService := portal.NewService(mockRepository)
	contextParam := context.Background()

	imaging, err := portalService.GetImaging(contextParam, testPatientFHIRID)

	assert.NoError(testingInstance, err)
	assert.Len(testingInstance, imaging, 1)
	assert.Equal(testingInstance, "Raio-X tórax", imaging[0].Title)
}

func TestPortalService_GetEncounters_PropagatesRepositoryError(testingInstance *testing.T) {
	mockRepository := buildDashboardFixture()
	mockRepository.EncountersErr = errors.New("healthcare api unavailable")
	portalService := portal.NewService(mockRepository)
	contextParam := context.Background()

	encounters, err := portalService.GetEncounters(contextParam, testPatientFHIRID)

	assert.Error(testingInstance, err)
	assert.Nil(testingInstance, encounters)
}

func TestPortalService_GetDashboard_StartsAtZeroTimeWhenUndefined(testingInstance *testing.T) {
	mockRepository := buildDashboardFixture()
	mockRepository.Encounters[testPatientFHIRID] = []portal.PortalEncounter{
		{FHIRResourceID: "enc-4", Status: "planned"},
	}
	portalService := portal.NewService(mockRepository)
	contextParam := context.Background()

	dashboard, err := portalService.GetDashboard(contextParam, testPatientFHIRID)

	assert.NoError(testingInstance, err)
	assert.Len(testingInstance, dashboard.UpcomingEncounters, 1)
	assert.Equal(testingInstance, time.Time{}, dashboard.UpcomingEncounters[0].StartedAt)
}
