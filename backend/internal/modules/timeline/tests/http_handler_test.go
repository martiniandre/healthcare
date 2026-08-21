package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/healthcare/backend/internal/modules/timeline"
	"github.com/healthcare/backend/internal/modules/timeline/mocks"
	"github.com/healthcare/backend/internal/shared/healthcare"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func buildTimelineHandler() (*timeline.HTTPHandler, *mocks.MockRepository) {
	mockRepository := buildTimelineFixture()
	return timeline.NewHTTPHandler(timeline.NewService(mockRepository)), mockRepository
}

func TestTimelineHTTPHandler_GetTimeline_ReturnsMergedFeed(t *testing.T) {
	handler, _ := buildTimelineHandler()

	request := httptest.NewRequest(http.MethodGet, "/api/v1/patients/fhir-pat-1/timeline", nil)
	responseRecorder := httptest.NewRecorder()

	handler.GetTimeline(responseRecorder, request)

	assert.Equal(t, http.StatusOK, responseRecorder.Code)

	var timelinePage timeline.TimelinePage
	unmarshalErr := json.Unmarshal(responseRecorder.Body.Bytes(), &timelinePage)
	require.NoError(t, unmarshalErr)
	assert.Len(t, timelinePage.Entries, 7)
	assert.Empty(t, timelinePage.UnavailableTypes)
}

func TestTimelineHTTPHandler_GetTimeline_ParsesQueryParameters(t *testing.T) {
	handler, mockRepository := buildTimelineHandler()
	mockRepository.Entries["Observation"] = []timeline.TimelineEntry{
		{ResourceType: "Observation", FHIRResourceID: "obs-1", RecordedAt: time.Date(2026, 8, 2, 10, 30, 0, 0, time.UTC)},
	}
	for _, resourceType := range []string{"Encounter", "Condition", "MedicationRequest", "DiagnosticReport", "ImagingStudy", "AllergyIntolerance"} {
		delete(mockRepository.Entries, resourceType)
	}

	request := httptest.NewRequest(http.MethodGet, "/api/v1/patients/fhir-pat-1/timeline?limit=5&before=2026-08-01T00%3A00%3A00Z&types=Observation&status=active", nil)
	responseRecorder := httptest.NewRecorder()

	handler.GetTimeline(responseRecorder, request)

	assert.Equal(t, http.StatusOK, responseRecorder.Code)
	assert.Equal(t, 5, mockRepository.ReceivedLimit)
	require.NotNil(t, mockRepository.ReceivedBefore)
	assert.True(t, mockRepository.ReceivedBefore.Equal(time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)))
}

func TestTimelineHTTPHandler_GetTimeline_RejectsInvalidLimit(t *testing.T) {
	handler, _ := buildTimelineHandler()

	request := httptest.NewRequest(http.MethodGet, "/api/v1/patients/fhir-pat-1/timeline?limit=9999", nil)
	responseRecorder := httptest.NewRecorder()

	handler.GetTimeline(responseRecorder, request)

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
}

func TestTimelineHTTPHandler_GetTimeline_RejectsInvalidBefore(t *testing.T) {
	handler, _ := buildTimelineHandler()

	request := httptest.NewRequest(http.MethodGet, "/api/v1/patients/fhir-pat-1/timeline?before=yesterday", nil)
	responseRecorder := httptest.NewRecorder()

	handler.GetTimeline(responseRecorder, request)

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
}

func TestTimelineHTTPHandler_GetTimeline_RejectsInvalidTypes(t *testing.T) {
	handler, _ := buildTimelineHandler()

	request := httptest.NewRequest(http.MethodGet, "/api/v1/patients/fhir-pat-1/timeline?types=Appointments", nil)
	responseRecorder := httptest.NewRecorder()

	handler.GetTimeline(responseRecorder, request)

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
}

func TestTimelineHTTPHandler_GetTimeline_MapsPatientNotFoundTo404(t *testing.T) {
	handler, mockRepository := buildTimelineHandler()
	mockRepository.PatientErr = &healthcare.NotFoundError{ResourceType: "Patient", ResourceID: "fhir-missing"}

	request := httptest.NewRequest(http.MethodGet, "/api/v1/patients/fhir-missing/timeline", nil)
	responseRecorder := httptest.NewRecorder()

	handler.GetTimeline(responseRecorder, request)

	assert.Equal(t, http.StatusNotFound, responseRecorder.Code)
}
