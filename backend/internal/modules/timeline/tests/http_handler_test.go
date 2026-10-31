package tests

import (
	"context"
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

func performTimelineRequest(testingT *testing.T, repository timeline.Repository, rawQuery string) *httptest.ResponseRecorder {
	testingT.Helper()
	handler := timeline.NewHTTPHandler(newTestService(repository))
	httpRequest := httptest.NewRequest(http.MethodGet, "/api/v1/patients/patient-1/timeline?"+rawQuery, nil)
	httpRequest.SetPathValue("patientFhirId", "patient-1")
	responseRecorder := httptest.NewRecorder()
	handler.GetTimeline(responseRecorder, httpRequest)
	return responseRecorder
}

func TestGetTimelineHandler_ReturnsContractShape(testingT *testing.T) {
	clinicalDate := time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)
	repository := &mocks.MockRepository{
		ListEntriesFn: func(ctx context.Context, entryType timeline.EntryType, query timeline.FetchQuery) ([]timeline.Entry, error) {
			switch entryType {
			case timeline.EntryTypeCondition:
				return []timeline.Entry{{Type: timeline.EntryTypeCondition, ResourceID: "cond-1", OccurredAt: fixedNow.Add(-time.Hour), ClinicalDate: &clinicalDate, Title: "Hypertension", Status: "active", EncounterID: "enc-9"}}, nil
			case timeline.EntryTypeAllergy:
				return nil, errFHIRUnavailable
			}
			return nil, nil
		},
	}

	responseRecorder := performTimelineRequest(testingT, repository, "")

	require.Equal(testingT, http.StatusOK, responseRecorder.Code)
	var body map[string]interface{}
	require.NoError(testingT, json.Unmarshal(responseRecorder.Body.Bytes(), &body))
	items := body["items"].([]interface{})
	require.Len(testingT, items, 1)
	item := items[0].(map[string]interface{})
	assert.Equal(testingT, "condition:cond-1", item["id"])
	assert.Equal(testingT, "condition", item["type"])
	assert.Equal(testingT, "2026-10-31T11:00:00Z", item["occurredAt"])
	assert.Equal(testingT, "2026-09-01T00:00:00Z", item["clinicalDate"])
	assert.Equal(testingT, "Hypertension", item["title"])
	assert.Equal(testingT, "", item["summary"])
	assert.Equal(testingT, "active", item["status"])
	assert.Equal(testingT, "cond-1", item["resourceId"])
	assert.Equal(testingT, "enc-9", item["encounterId"])
	assert.Nil(testingT, body["nextCursor"])
	assert.Equal(testingT, map[string]interface{}{"failedTypes": []interface{}{"allergy"}}, body["partial"])
	assert.Equal(testingT, map[string]interface{}{"from": "2026-04-30T12:00:00Z", "to": "2026-10-31T12:00:00Z"}, body["window"])
}

func TestGetTimelineHandler_NullableFieldsSerializeAsNull(testingT *testing.T) {
	repository := &mocks.MockRepository{
		ListEntriesFn: func(ctx context.Context, entryType timeline.EntryType, query timeline.FetchQuery) ([]timeline.Entry, error) {
			if entryType == timeline.EntryTypeImaging {
				return []timeline.Entry{entryAt(timeline.EntryTypeImaging, "img-1", fixedNow.Add(-time.Hour))}, nil
			}
			return nil, nil
		},
	}

	responseRecorder := performTimelineRequest(testingT, repository, "")

	var body struct {
		Items []map[string]interface{} `json:"items"`
	}
	require.NoError(testingT, json.Unmarshal(responseRecorder.Body.Bytes(), &body))
	require.Len(testingT, body.Items, 1)
	clinicalDate, hasClinicalDate := body.Items[0]["clinicalDate"]
	encounterID, hasEncounterID := body.Items[0]["encounterId"]
	assert.True(testingT, hasClinicalDate)
	assert.Nil(testingT, clinicalDate)
	assert.True(testingT, hasEncounterID)
	assert.Nil(testingT, encounterID)
}

func TestGetTimelineHandler_ReturnsNextCursorWhenMoreEntriesExist(testingT *testing.T) {
	repository := &mocks.MockRepository{
		ListEntriesFn: func(ctx context.Context, entryType timeline.EntryType, query timeline.FetchQuery) ([]timeline.Entry, error) {
			if entryType == timeline.EntryTypeObservation {
				return []timeline.Entry{
					entryAt(timeline.EntryTypeObservation, "obs-1", fixedNow.Add(-1*time.Hour)),
					entryAt(timeline.EntryTypeObservation, "obs-2", fixedNow.Add(-2*time.Hour)),
				}, nil
			}
			return nil, nil
		},
	}

	responseRecorder := performTimelineRequest(testingT, repository, "limit=1")

	var body struct {
		Items      []map[string]interface{} `json:"items"`
		NextCursor *string                  `json:"nextCursor"`
	}
	require.NoError(testingT, json.Unmarshal(responseRecorder.Body.Bytes(), &body))
	assert.Len(testingT, body.Items, 1)
	require.NotNil(testingT, body.NextCursor)
	assert.NotEmpty(testingT, *body.NextCursor)
}

func TestGetTimelineHandler_InvalidParametersReturnBadRequest(testingT *testing.T) {
	for _, rawQuery := range []string{"from=abc", "types=appointment", "limit=0", "cursor=bogus"} {
		responseRecorder := performTimelineRequest(testingT, &mocks.MockRepository{}, rawQuery)
		assert.Equal(testingT, http.StatusBadRequest, responseRecorder.Code, rawQuery)
	}
}

func TestGetTimelineHandler_UnknownPatientReturnsNotFound(testingT *testing.T) {
	repository := &mocks.MockRepository{
		PatientExistsFn: func(ctx context.Context, patientFHIRID string) error {
			return &healthcare.NotFoundError{ResourceType: "Patient", ResourceID: patientFHIRID}
		},
	}

	responseRecorder := performTimelineRequest(testingT, repository, "")

	assert.Equal(testingT, http.StatusNotFound, responseRecorder.Code)
}

func TestGetTimelineHandler_TotalFailureReturnsServiceUnavailable(testingT *testing.T) {
	repository := &mocks.MockRepository{
		ListEntriesFn: func(ctx context.Context, entryType timeline.EntryType, query timeline.FetchQuery) ([]timeline.Entry, error) {
			return nil, errFHIRUnavailable
		},
	}

	responseRecorder := performTimelineRequest(testingT, repository, "")

	assert.Equal(testingT, http.StatusServiceUnavailable, responseRecorder.Code)
}
