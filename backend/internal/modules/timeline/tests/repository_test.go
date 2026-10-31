package tests

import (
	"context"
	"encoding/json"
	"net/url"
	"testing"
	"time"

	"github.com/healthcare/backend/internal/modules/timeline"
	"github.com/healthcare/backend/internal/modules/timeline/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func bundleOf(resources ...string) json.RawMessage {
	joined := ""
	for resourceIndex, resource := range resources {
		if resourceIndex > 0 {
			joined += ","
		}
		joined += `{"resource":` + resource + `}`
	}
	return json.RawMessage(`{"resourceType":"Bundle","entry":[` + joined + `]}`)
}

func listEntriesFor(testingT *testing.T, entryType timeline.EntryType, resource string) []timeline.Entry {
	testingT.Helper()
	fhirClient := &mocks.MockFHIRClient{
		SearchResourcesPageFn: func(ctx context.Context, resourceType string, queryParams string) (json.RawMessage, error) {
			return bundleOf(resource), nil
		},
	}
	entries, listError := timeline.NewRepository(fhirClient).ListEntries(context.Background(), entryType, timeline.FetchQuery{
		PatientFHIRID: "patient-1",
		Window:        timeline.Window{From: fixedNow.AddDate(0, -6, 0), To: fixedNow},
		MaxEntries:    31,
	})
	require.NoError(testingT, listError)
	return entries
}

func TestRepositoryMapping_OrderingKeyPerResourceType(testingT *testing.T) {
	testCases := []struct {
		name         string
		entryType    timeline.EntryType
		resource     string
		expectedTime time.Time
	}{
		{
			name:         "encounter uses period start",
			entryType:    timeline.EntryTypeEncounter,
			resource:     `{"resourceType":"Encounter","id":"e1","status":"finished","period":{"start":"2026-09-01T09:00:00Z","end":"2026-09-03T09:00:00Z"}}`,
			expectedTime: time.Date(2026, 9, 1, 9, 0, 0, 0, time.UTC),
		},
		{
			name:         "observation uses effectiveDateTime",
			entryType:    timeline.EntryTypeObservation,
			resource:     `{"resourceType":"Observation","id":"o1","status":"final","effectiveDateTime":"2026-09-02T10:00:00Z","issued":"2026-09-09T10:00:00Z"}`,
			expectedTime: time.Date(2026, 9, 2, 10, 0, 0, 0, time.UTC),
		},
		{
			name:         "condition uses recordedDate not onset",
			entryType:    timeline.EntryTypeCondition,
			resource:     `{"resourceType":"Condition","id":"c1","onsetDateTime":"2020-01-01T00:00:00Z","recordedDate":"2026-09-03T11:00:00Z"}`,
			expectedTime: time.Date(2026, 9, 3, 11, 0, 0, 0, time.UTC),
		},
		{
			name:         "allergy uses recordedDate",
			entryType:    timeline.EntryTypeAllergy,
			resource:     `{"resourceType":"AllergyIntolerance","id":"a1","recordedDate":"2026-09-04T12:00:00Z"}`,
			expectedTime: time.Date(2026, 9, 4, 12, 0, 0, 0, time.UTC),
		},
		{
			name:         "medication uses authoredOn",
			entryType:    timeline.EntryTypeMedication,
			resource:     `{"resourceType":"MedicationRequest","id":"m1","status":"active","authoredOn":"2026-09-05T13:00:00Z"}`,
			expectedTime: time.Date(2026, 9, 5, 13, 0, 0, 0, time.UTC),
		},
		{
			name:         "report uses issued",
			entryType:    timeline.EntryTypeReport,
			resource:     `{"resourceType":"DiagnosticReport","id":"r1","status":"final","issued":"2026-09-06T14:00:00Z","meta":{"lastUpdated":"2026-09-20T14:00:00Z"}}`,
			expectedTime: time.Date(2026, 9, 6, 14, 0, 0, 0, time.UTC),
		},
		{
			name:         "imaging uses started",
			entryType:    timeline.EntryTypeImaging,
			resource:     `{"resourceType":"ImagingStudy","id":"i1","status":"available","started":"2026-09-07T15:00:00Z"}`,
			expectedTime: time.Date(2026, 9, 7, 15, 0, 0, 0, time.UTC),
		},
		{
			name:         "date only value is read as midnight UTC",
			entryType:    timeline.EntryTypeCondition,
			resource:     `{"resourceType":"Condition","id":"c2","recordedDate":"2026-09-08"}`,
			expectedTime: time.Date(2026, 9, 8, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, testCase := range testCases {
		testingT.Run(testCase.name, func(subTest *testing.T) {
			entries := listEntriesFor(subTest, testCase.entryType, testCase.resource)

			require.Len(subTest, entries, 1)
			assert.Equal(subTest, testCase.entryType, entries[0].Type)
			assert.True(subTest, testCase.expectedTime.Equal(entries[0].OccurredAt), "got %s", entries[0].OccurredAt)
		})
	}
}

func TestRepositoryMapping_SkipsResourcesWithoutRecordedTime(testingT *testing.T) {
	entries := listEntriesFor(testingT, timeline.EntryTypeObservation, `{"resourceType":"Observation","id":"o1","status":"final"}`)

	assert.Empty(testingT, entries)
}

func TestRepositoryMapping_DisplayFields(testingT *testing.T) {
	observationEntries := listEntriesFor(testingT, timeline.EntryTypeObservation, `{"resourceType":"Observation","id":"o1","status":"final","code":{"coding":[{"code":"8867-4","display":"Heart rate"}]},"encounter":{"reference":"Encounter/enc-3"},"valueQuantity":{"value":72,"unit":"bpm"},"effectiveDateTime":"2026-09-02T10:00:00Z"}`)
	require.Len(testingT, observationEntries, 1)
	assert.Equal(testingT, "Heart rate", observationEntries[0].Title)
	assert.Equal(testingT, "72 bpm", observationEntries[0].Summary)
	assert.Equal(testingT, "final", observationEntries[0].Status)
	assert.Equal(testingT, "enc-3", observationEntries[0].EncounterID)
	assert.Equal(testingT, "observation:o1", observationEntries[0].ID())

	conditionEntries := listEntriesFor(testingT, timeline.EntryTypeCondition, `{"resourceType":"Condition","id":"c1","clinicalStatus":{"coding":[{"code":"active"}]},"code":{"coding":[{"code":"I10","display":"Hypertension"}]},"onsetDateTime":"2020-01-01T00:00:00Z","recordedDate":"2026-09-03T11:00:00Z"}`)
	require.Len(testingT, conditionEntries, 1)
	assert.Equal(testingT, "Hypertension", conditionEntries[0].Title)
	assert.Equal(testingT, "I10", conditionEntries[0].Summary)
	assert.Equal(testingT, "active", conditionEntries[0].Status)
	require.NotNil(testingT, conditionEntries[0].ClinicalDate)
	assert.True(testingT, time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC).Equal(*conditionEntries[0].ClinicalDate))

	encounterEntries := listEntriesFor(testingT, timeline.EntryTypeEncounter, `{"resourceType":"Encounter","id":"e1","status":"in-progress","period":{"start":"2026-09-01T09:00:00Z"},"reasonCode":[{"coding":[{"code":"R50","display":"Fever"}]}]}`)
	require.Len(testingT, encounterEntries, 1)
	assert.Equal(testingT, "Fever", encounterEntries[0].Title)
	assert.Equal(testingT, "e1", encounterEntries[0].EncounterID)
}

func TestRepositorySearch_BuildsDateFilteredPagedQueryPerType(testingT *testing.T) {
	fhirClient := &mocks.MockFHIRClient{}
	repository := timeline.NewRepository(fhirClient)
	query := timeline.FetchQuery{
		PatientFHIRID: "patient-1",
		Window:        timeline.Window{From: time.Date(2026, 4, 30, 12, 0, 0, 0, time.UTC), To: time.Date(2026, 10, 31, 12, 0, 0, 0, time.UTC)},
		MaxEntries:    31,
	}

	for _, entryType := range timeline.AllEntryTypes {
		_, listError := repository.ListEntries(context.Background(), entryType, query)
		require.NoError(testingT, listError)
	}

	searchCalls := fhirClient.SearchCalls()
	require.Len(testingT, searchCalls, len(timeline.AllEntryTypes))
	expectations := map[string][3]string{
		"Encounter":          {"subject", "date", "-date"},
		"Observation":        {"subject", "date", "-date"},
		"Condition":          {"subject", "recorded-date", "-recorded-date"},
		"AllergyIntolerance": {"patient", "date", "-date"},
		"MedicationRequest":  {"subject", "authoredon", "-authoredon"},
		"DiagnosticReport":   {"subject", "issued", "-issued"},
		"ImagingStudy":       {"subject", "started", "-started"},
	}
	for _, searchCall := range searchCalls {
		expected, known := expectations[searchCall.ResourceType]
		require.True(testingT, known, searchCall.ResourceType)
		parsedQuery, parseError := url.ParseQuery(searchCall.QueryParams)
		require.NoError(testingT, parseError)
		assert.Equal(testingT, []string{"Patient/patient-1"}, parsedQuery[expected[0]], searchCall.ResourceType)
		assert.Equal(testingT, []string{"ge2026-04-30T12:00:00Z", "le2026-10-31T12:00:00Z"}, parsedQuery[expected[1]], searchCall.ResourceType)
		assert.Equal(testingT, []string{expected[2]}, parsedQuery["_sort"], searchCall.ResourceType)
		assert.Equal(testingT, []string{"31"}, parsedQuery["_count"], searchCall.ResourceType)
	}
}

func TestRepositorySearch_PropagatesSearchFailure(testingT *testing.T) {
	fhirClient := &mocks.MockFHIRClient{
		SearchResourcesPageFn: func(ctx context.Context, resourceType string, queryParams string) (json.RawMessage, error) {
			return nil, errFHIRUnavailable
		},
	}

	_, listError := timeline.NewRepository(fhirClient).ListEntries(context.Background(), timeline.EntryTypeReport, timeline.FetchQuery{PatientFHIRID: "patient-1"})

	assert.ErrorIs(testingT, listError, errFHIRUnavailable)
}

func TestRepositoryPatientExists_WrapsLookupError(testingT *testing.T) {
	fhirClient := &mocks.MockFHIRClient{
		GetResourceFn: func(ctx context.Context, resourceType string, resourceID string) (json.RawMessage, error) {
			return nil, errFHIRUnavailable
		},
	}

	patientError := timeline.NewRepository(fhirClient).PatientExists(context.Background(), "patient-1")

	assert.ErrorIs(testingT, patientError, errFHIRUnavailable)
}
