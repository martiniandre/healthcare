package tests

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/healthcare/backend/internal/modules/portal"
	"github.com/healthcare/backend/internal/shared/ctxkeys"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakePortalService struct {
	dashboard    *portal.PortalDashboard
	encounters   []portal.PortalEncounter
	observations []portal.PortalObservation
	conditions   []portal.PortalCondition
	medications  []portal.PortalMedication
	reports      []portal.PortalReport
	imaging      []portal.PortalImaging
	serviceErr   error
}

func (fakeService *fakePortalService) GetDashboard(ctx context.Context, fhirResourceID string) (*portal.PortalDashboard, error) {
	return fakeService.dashboard, fakeService.serviceErr
}

func (fakeService *fakePortalService) GetEncounters(ctx context.Context, fhirResourceID string) ([]portal.PortalEncounter, error) {
	return fakeService.encounters, fakeService.serviceErr
}

func (fakeService *fakePortalService) GetObservations(ctx context.Context, fhirResourceID string) ([]portal.PortalObservation, error) {
	return fakeService.observations, fakeService.serviceErr
}

func (fakeService *fakePortalService) GetConditions(ctx context.Context, fhirResourceID string) ([]portal.PortalCondition, error) {
	return fakeService.conditions, fakeService.serviceErr
}

func (fakeService *fakePortalService) GetMedications(ctx context.Context, fhirResourceID string) ([]portal.PortalMedication, error) {
	return fakeService.medications, fakeService.serviceErr
}

func (fakeService *fakePortalService) GetReports(ctx context.Context, fhirResourceID string) ([]portal.PortalReport, error) {
	return fakeService.reports, fakeService.serviceErr
}

func (fakeService *fakePortalService) GetImaging(ctx context.Context, fhirResourceID string) ([]portal.PortalImaging, error) {
	return fakeService.imaging, fakeService.serviceErr
}

func buildAuthenticatedPortalRequest() *http.Request {
	httpRequest := httptest.NewRequest(http.MethodGet, "/api/v1/portal/dashboard", nil)
	contextWithUser := context.WithValue(httpRequest.Context(), ctxkeys.UserIDKey, testPatientFHIRID)
	contextWithRole := context.WithValue(contextWithUser, ctxkeys.RoleKey, "patient")
	return httpRequest.WithContext(contextWithRole)
}

func TestPortalHTTPHandler_GetDashboard_ReturnsDashboardForAuthenticatedPatient(testingInstance *testing.T) {
	fakeService := &fakePortalService{
		dashboard: &portal.PortalDashboard{
			PatientInfo: portal.PatientInfo{FHIRResourceID: testPatientFHIRID, FullName: "Mariana Costa"},
		},
	}
	handler := portal.NewHTTPHandler(fakeService)
	responseRecorder := httptest.NewRecorder()

	handler.GetDashboard(responseRecorder, buildAuthenticatedPortalRequest())

	require.Equal(testingInstance, http.StatusOK, responseRecorder.Code)
	assert.Contains(testingInstance, responseRecorder.Body.String(), "Mariana Costa")
}

func TestPortalHTTPHandler_GetDashboard_UnauthorizedWithoutUser(testingInstance *testing.T) {
	fakeService := &fakePortalService{dashboard: &portal.PortalDashboard{}}
	handler := portal.NewHTTPHandler(fakeService)
	responseRecorder := httptest.NewRecorder()
	httpRequest := httptest.NewRequest(http.MethodGet, "/api/v1/portal/dashboard", nil)

	handler.GetDashboard(responseRecorder, httpRequest)

	require.Equal(testingInstance, http.StatusUnauthorized, responseRecorder.Code)
}

func TestPortalHTTPHandler_GetDashboard_ServiceErrorReturnsInternalServerError(testingInstance *testing.T) {
	fakeService := &fakePortalService{serviceErr: portal.ErrPatientNotFound}
	handler := portal.NewHTTPHandler(fakeService)
	responseRecorder := httptest.NewRecorder()

	handler.GetDashboard(responseRecorder, buildAuthenticatedPortalRequest())

	require.Equal(testingInstance, http.StatusInternalServerError, responseRecorder.Code)
}

func TestPortalHTTPHandler_GetEncounters_ReturnsEncounters(testingInstance *testing.T) {
	fakeService := &fakePortalService{
		encounters: []portal.PortalEncounter{{FHIRResourceID: "enc-1", Status: "planned"}},
	}
	handler := portal.NewHTTPHandler(fakeService)
	responseRecorder := httptest.NewRecorder()

	handler.GetEncounters(responseRecorder, buildAuthenticatedPortalRequest())

	require.Equal(testingInstance, http.StatusOK, responseRecorder.Code)
	assert.Contains(testingInstance, responseRecorder.Body.String(), "enc-1")
}

func TestPortalHTTPHandler_GetObservations_ReturnsObservations(testingInstance *testing.T) {
	fakeService := &fakePortalService{
		observations: []portal.PortalObservation{{FHIRResourceID: "obs-1", CodeDisplay: "Pressão arterial"}},
	}
	handler := portal.NewHTTPHandler(fakeService)
	responseRecorder := httptest.NewRecorder()

	handler.GetObservations(responseRecorder, buildAuthenticatedPortalRequest())

	require.Equal(testingInstance, http.StatusOK, responseRecorder.Code)
	assert.Contains(testingInstance, responseRecorder.Body.String(), "Pressão arterial")
}

func TestPortalHTTPHandler_GetConditions_ReturnsConditions(testingInstance *testing.T) {
	fakeService := &fakePortalService{
		conditions: []portal.PortalCondition{{FHIRResourceID: "cond-1", CodeDisplay: "Diabetes"}},
	}
	handler := portal.NewHTTPHandler(fakeService)
	responseRecorder := httptest.NewRecorder()

	handler.GetConditions(responseRecorder, buildAuthenticatedPortalRequest())

	require.Equal(testingInstance, http.StatusOK, responseRecorder.Code)
	assert.Contains(testingInstance, responseRecorder.Body.String(), "Diabetes")
}

func TestPortalHTTPHandler_GetMedications_ReturnsMedications(testingInstance *testing.T) {
	fakeService := &fakePortalService{
		medications: []portal.PortalMedication{{FHIRResourceID: "med-1", MedicationName: "Metformina"}},
	}
	handler := portal.NewHTTPHandler(fakeService)
	responseRecorder := httptest.NewRecorder()

	handler.GetMedications(responseRecorder, buildAuthenticatedPortalRequest())

	require.Equal(testingInstance, http.StatusOK, responseRecorder.Code)
	assert.Contains(testingInstance, responseRecorder.Body.String(), "Metformina")
}

func TestPortalHTTPHandler_GetReports_ReturnsReports(testingInstance *testing.T) {
	fakeService := &fakePortalService{
		reports: []portal.PortalReport{{FHIRResourceID: "report-1", ReportDisplay: "Hemograma"}},
	}
	handler := portal.NewHTTPHandler(fakeService)
	responseRecorder := httptest.NewRecorder()

	handler.GetReports(responseRecorder, buildAuthenticatedPortalRequest())

	require.Equal(testingInstance, http.StatusOK, responseRecorder.Code)
	assert.Contains(testingInstance, responseRecorder.Body.String(), "Hemograma")
}

func TestPortalHTTPHandler_GetImaging_ReturnsImaging(testingInstance *testing.T) {
	fakeService := &fakePortalService{
		imaging: []portal.PortalImaging{{FHIRResourceID: "img-1", Title: "Raio-X tórax"}},
	}
	handler := portal.NewHTTPHandler(fakeService)
	responseRecorder := httptest.NewRecorder()

	handler.GetImaging(responseRecorder, buildAuthenticatedPortalRequest())

	require.Equal(testingInstance, http.StatusOK, responseRecorder.Code)
	assert.Contains(testingInstance, responseRecorder.Body.String(), "Raio-X tórax")
}
