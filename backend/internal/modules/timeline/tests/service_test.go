package tests

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/healthcare/backend/internal/modules/timeline"
	"github.com/healthcare/backend/internal/modules/timeline/mocks"
	"github.com/healthcare/backend/internal/shared/apperrors"
	"github.com/healthcare/backend/internal/shared/healthcare"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func buildTimelineFixture() *mocks.MockRepository {
	mockRepository := mocks.NewMockRepository()

	mockRepository.Entries["Encounter"] = []timeline.TimelineEntry{
		{ResourceType: "Encounter", FHIRResourceID: "enc-1", Title: "Cardiology follow-up", Status: "finished", RecordedAt: time.Date(2026, 8, 12, 9, 0, 0, 0, time.UTC)},
	}
	mockRepository.Entries["Observation"] = []timeline.TimelineEntry{
		{ResourceType: "Observation", FHIRResourceID: "obs-1", Title: "Blood pressure", Code: "85354-9", RecordedAt: time.Date(2026, 8, 2, 10, 30, 0, 0, time.UTC)},
	}
	mockRepository.Entries["Condition"] = []timeline.TimelineEntry{
		{ResourceType: "Condition", FHIRResourceID: "cond-1", Title: "Hypertension", Code: "I10", Status: "active", RecordedAt: time.Date(2026, 7, 15, 8, 0, 0, 0, time.UTC)},
	}
	mockRepository.Entries["MedicationRequest"] = []timeline.TimelineEntry{
		{ResourceType: "MedicationRequest", FHIRResourceID: "med-1", Title: "Losartan", Status: "active", RecordedAt: time.Date(2026, 8, 10, 11, 0, 0, 0, time.UTC)},
	}
	mockRepository.Entries["DiagnosticReport"] = []timeline.TimelineEntry{
		{ResourceType: "DiagnosticReport", FHIRResourceID: "rep-1", Title: "Chest X-ray", Status: "final", RecordedAt: time.Date(2026, 7, 28, 14, 0, 0, 0, time.UTC)},
	}
	mockRepository.Entries["ImagingStudy"] = []timeline.TimelineEntry{
		{ResourceType: "ImagingStudy", FHIRResourceID: "img-1", Title: "Radiograph", Modality: "DX", RecordedAt: time.Date(2026, 6, 20, 16, 0, 0, 0, time.UTC)},
	}
	mockRepository.Entries["AllergyIntolerance"] = []timeline.TimelineEntry{
		{ResourceType: "AllergyIntolerance", FHIRResourceID: "alg-1", Title: "Penicillin", Criticality: "high", RecordedAt: time.Date(2026, 7, 1, 12, 0, 0, 0, time.UTC)},
	}

	return mockRepository
}

func TestTimelineService_GetTimeline_MergesAllTypesInDescendingOrder(t *testing.T) {
	mockRepository := buildTimelineFixture()
	timelineService := timeline.NewService(mockRepository)

	timelinePage, serviceErr := timelineService.GetTimeline(context.Background(), "fhir-pat-1", timeline.TimelineFilter{})

	assert.NoError(t, serviceErr)
	require.NotNil(t, timelinePage)
	assert.Len(t, timelinePage.Entries, 7)
	assert.Empty(t, timelinePage.UnavailableTypes)

	expectedOrder := []string{"Encounter", "MedicationRequest", "Observation", "DiagnosticReport", "Condition", "AllergyIntolerance", "ImagingStudy"}
	for entryIndex, expectedType := range expectedOrder {
		assert.Equal(t, expectedType, timelinePage.Entries[entryIndex].ResourceType)
	}

	assert.Nil(t, timelinePage.NextCursor)
}

func TestTimelineService_GetTimeline_ReturnsPartialSuccessWhenOneTypeFails(t *testing.T) {
	mockRepository := buildTimelineFixture()
	mockRepository.FetchErrors["DiagnosticReport"] = errors.New("upstream timeout")
	timelineService := timeline.NewService(mockRepository)

	timelinePage, serviceErr := timelineService.GetTimeline(context.Background(), "fhir-pat-1", timeline.TimelineFilter{})

	assert.NoError(t, serviceErr)
	require.NotNil(t, timelinePage)
	assert.Equal(t, []string{"DiagnosticReport"}, timelinePage.UnavailableTypes)
	assert.Len(t, timelinePage.Entries, 6)
}

func TestTimelineService_GetTimeline_FailsWhenAllRequestedTypesFail(t *testing.T) {
	mockRepository := mocks.NewMockRepository()
	mockRepository.FetchErrors["Encounter"] = errors.New("boom")
	mockRepository.FetchErrors["Observation"] = errors.New("boom")
	timelineService := timeline.NewService(mockRepository)

	timelinePage, serviceErr := timelineService.GetTimeline(context.Background(), "fhir-pat-1", timeline.TimelineFilter{Types: []string{"Encounter", "Observation"}})

	assert.Error(t, serviceErr)
	assert.Nil(t, timelinePage)

	var appError apperrors.AppError
	assert.ErrorAs(t, serviceErr, &appError)
	assert.Equal(t, http.StatusServiceUnavailable, appError.HTTPCode)
}

func TestTimelineService_GetTimeline_ReturnsNotFoundForMissingPatient(t *testing.T) {
	mockRepository := mocks.NewMockRepository()
	mockRepository.PatientErr = &healthcare.NotFoundError{ResourceType: "Patient", ResourceID: "fhir-missing"}
	timelineService := timeline.NewService(mockRepository)

	timelinePage, serviceErr := timelineService.GetTimeline(context.Background(), "fhir-missing", timeline.TimelineFilter{})

	assert.Error(t, serviceErr)
	assert.Nil(t, timelinePage)

	var appError apperrors.AppError
	assert.ErrorAs(t, serviceErr, &appError)
	assert.Equal(t, http.StatusNotFound, appError.HTTPCode)
}

func TestTimelineService_GetTimeline_RejectsUnsupportedResourceType(t *testing.T) {
	mockRepository := buildTimelineFixture()
	timelineService := timeline.NewService(mockRepository)

	timelinePage, serviceErr := timelineService.GetTimeline(context.Background(), "fhir-pat-1", timeline.TimelineFilter{Types: []string{"Nonsense"}})

	assert.Error(t, serviceErr)
	assert.Nil(t, timelinePage)

	var appError apperrors.AppError
	assert.ErrorAs(t, serviceErr, &appError)
	assert.Equal(t, http.StatusBadRequest, appError.HTTPCode)
}

func TestTimelineService_GetTimeline_RejectsUnsupportedStatusFilter(t *testing.T) {
	mockRepository := buildTimelineFixture()
	timelineService := timeline.NewService(mockRepository)

	timelinePage, serviceErr := timelineService.GetTimeline(context.Background(), "fhir-pat-1", timeline.TimelineFilter{Status: "cancelled"})

	assert.Error(t, serviceErr)
	assert.Nil(t, timelinePage)
}

func TestTimelineService_GetTimeline_AppliesDefaultPageSizeAndActiveFilter(t *testing.T) {
	mockRepository := buildTimelineFixture()
	timelineService := timeline.NewService(mockRepository)

	_, serviceErr := timelineService.GetTimeline(context.Background(), "fhir-pat-1", timeline.TimelineFilter{Status: "active"})

	assert.NoError(t, serviceErr)
	assert.Equal(t, timeline.DefaultPageSize, mockRepository.ReceivedLimit)
	assert.True(t, mockRepository.ReceivedOnlyActive)
}

func TestTimelineService_GetTimeline_ForwardsCursorToRepositories(t *testing.T) {
	mockRepository := buildTimelineFixture()
	timelineService := timeline.NewService(mockRepository)

	cursorDate := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	_, serviceErr := timelineService.GetTimeline(context.Background(), "fhir-pat-1", timeline.TimelineFilter{Before: &cursorDate})

	assert.NoError(t, serviceErr)
	require.NotNil(t, mockRepository.ReceivedBefore)
	assert.True(t, mockRepository.ReceivedBefore.Equal(cursorDate))
}

func TestTimelineService_GetTimeline_SetsNextCursorWhenPageIsFull(t *testing.T) {
	mockRepository := mocks.NewMockRepository()
	overflowEntries := make([]timeline.TimelineEntry, 0, timeline.DefaultPageSize+5)
	for entryIndex := 0; entryIndex < timeline.DefaultPageSize+5; entryIndex++ {
		overflowEntries = append(overflowEntries, timeline.TimelineEntry{
			ResourceType:   "Observation",
			FHIRResourceID: string(rune('a' + entryIndex)),
			RecordedAt:     time.Date(2026, 8, 1, entryIndex%24, 0, 0, 0, time.UTC).Add(time.Duration(entryIndex/24) * time.Hour),
		})
	}
	mockRepository.Entries["Observation"] = overflowEntries
	timelineService := timeline.NewService(mockRepository)

	timelinePage, serviceErr := timelineService.GetTimeline(context.Background(), "fhir-pat-1", timeline.TimelineFilter{})

	assert.NoError(t, serviceErr)
	require.NotNil(t, timelinePage)
	assert.Len(t, timelinePage.Entries, timeline.DefaultPageSize)
	require.NotNil(t, timelinePage.NextCursor)
	assert.True(t, timelinePage.NextCursor.Equal(timelinePage.Entries[timeline.DefaultPageSize-1].RecordedAt))
}
