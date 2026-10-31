package tests

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/healthcare/backend/internal/modules/timeline"
	"github.com/healthcare/backend/internal/modules/timeline/mocks"
	"github.com/healthcare/backend/internal/shared/apperrors"
	"github.com/healthcare/backend/internal/shared/healthcare"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var fixedNow = time.Date(2026, 10, 31, 12, 0, 0, 0, time.UTC)

var errFHIRUnavailable = errors.New("fhir unavailable")

func fixedClock() time.Time {
	return fixedNow
}

func newTestService(repository timeline.Repository) timeline.Service {
	return timeline.NewServiceWithClock(repository, fixedClock)
}

func entryAt(entryType timeline.EntryType, resourceID string, occurredAt time.Time) timeline.Entry {
	return timeline.Entry{Type: entryType, ResourceID: resourceID, OccurredAt: occurredAt, Title: resourceID}
}

func repositoryReturning(entriesByType map[timeline.EntryType][]timeline.Entry) *mocks.MockRepository {
	return &mocks.MockRepository{
		ListEntriesFn: func(ctx context.Context, entryType timeline.EntryType, query timeline.FetchQuery) ([]timeline.Entry, error) {
			return entriesByType[entryType], nil
		},
	}
}

func entryIdentifiers(entries []timeline.Entry) []string {
	identifiers := make([]string, 0, len(entries))
	for _, entry := range entries {
		identifiers = append(identifiers, entry.ID())
	}
	return identifiers
}

func TestGetTimeline_MergesTypesNewestFirstWithStableTiebreak(testingT *testing.T) {
	sharedMoment := fixedNow.Add(-48 * time.Hour)
	timelineService := newTestService(repositoryReturning(map[timeline.EntryType][]timeline.Entry{
		timeline.EntryTypeEncounter:   {entryAt(timeline.EntryTypeEncounter, "enc-1", fixedNow.Add(-72*time.Hour))},
		timeline.EntryTypeObservation: {entryAt(timeline.EntryTypeObservation, "obs-b", sharedMoment), entryAt(timeline.EntryTypeObservation, "obs-a", sharedMoment)},
		timeline.EntryTypeReport:      {entryAt(timeline.EntryTypeReport, "rep-1", fixedNow.Add(-1*time.Hour))},
		timeline.EntryTypeCondition:   {entryAt(timeline.EntryTypeCondition, "cond-1", sharedMoment)},
	}))

	page, serviceError := timelineService.GetTimeline(context.Background(), timeline.ListInput{PatientFHIRID: "patient-1"})

	require.NoError(testingT, serviceError)
	assert.Equal(testingT, []string{"report:rep-1", "condition:cond-1", "observation:obs-a", "observation:obs-b", "encounter:enc-1"}, entryIdentifiers(page.Entries))
	assert.Empty(testingT, page.FailedTypes)
	assert.Empty(testingT, page.NextCursor)
}

func TestGetTimeline_DefaultWindowIsLastSixMonths(testingT *testing.T) {
	var capturedQuery timeline.FetchQuery
	repository := &mocks.MockRepository{
		ListEntriesFn: func(ctx context.Context, entryType timeline.EntryType, query timeline.FetchQuery) ([]timeline.Entry, error) {
			if entryType == timeline.EntryTypeEncounter {
				capturedQuery = query
			}
			return nil, nil
		},
	}

	page, serviceError := newTestService(repository).GetTimeline(context.Background(), timeline.ListInput{PatientFHIRID: "patient-1"})

	require.NoError(testingT, serviceError)
	assert.Equal(testingT, time.Date(2026, 4, 30, 12, 0, 0, 0, time.UTC), page.Window.From)
	assert.Equal(testingT, fixedNow, page.Window.To)
	assert.Equal(testingT, page.Window, capturedQuery.Window)
	assert.Equal(testingT, timeline.DefaultLimit+1, capturedQuery.MaxEntries)
	assert.Equal(testingT, "patient-1", capturedQuery.PatientFHIRID)
}

func TestGetTimeline_ClientProvidedWindowAndDateOnlyUpperBound(testingT *testing.T) {
	timelineService := newTestService(repositoryReturning(map[timeline.EntryType][]timeline.Entry{
		timeline.EntryTypeEncounter: {
			entryAt(timeline.EntryTypeEncounter, "inside-end-of-day", time.Date(2026, 9, 30, 23, 0, 0, 0, time.UTC)),
			entryAt(timeline.EntryTypeEncounter, "after-window", time.Date(2026, 10, 1, 0, 0, 1, 0, time.UTC)),
			entryAt(timeline.EntryTypeEncounter, "before-window", time.Date(2026, 8, 31, 23, 59, 0, 0, time.UTC)),
		},
	}))

	page, serviceError := timelineService.GetTimeline(context.Background(), timeline.ListInput{
		PatientFHIRID: "patient-1",
		From:          "2026-09-01",
		To:            "2026-09-30",
	})

	require.NoError(testingT, serviceError)
	assert.Equal(testingT, []string{"encounter:inside-end-of-day"}, entryIdentifiers(page.Entries))
	assert.Equal(testingT, time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC), page.Window.From)
	assert.Equal(testingT, time.Date(2026, 9, 30, 23, 59, 59, 999999999, time.UTC), page.Window.To)
}

func TestGetTimeline_TypesFilterOnlyQueriesRequestedTypes(testingT *testing.T) {
	queriedTypes := make(chan timeline.EntryType, len(timeline.AllEntryTypes))
	repository := &mocks.MockRepository{
		ListEntriesFn: func(ctx context.Context, entryType timeline.EntryType, query timeline.FetchQuery) ([]timeline.Entry, error) {
			queriedTypes <- entryType
			return nil, nil
		},
	}

	_, serviceError := newTestService(repository).GetTimeline(context.Background(), timeline.ListInput{
		PatientFHIRID: "patient-1",
		Types:         "report, Encounter,report",
	})
	close(queriedTypes)

	require.NoError(testingT, serviceError)
	collected := make([]timeline.EntryType, 0)
	for queriedType := range queriedTypes {
		collected = append(collected, queriedType)
	}
	assert.ElementsMatch(testingT, []timeline.EntryType{timeline.EntryTypeEncounter, timeline.EntryTypeReport}, collected)
}

func TestGetTimeline_CursorPagesThroughAllEntriesWithoutDuplicates(testingT *testing.T) {
	entriesByType := map[timeline.EntryType][]timeline.Entry{}
	expectedOrder := make([]string, 0)
	for hoursAgo := 1; hoursAgo <= 7; hoursAgo++ {
		entryType := timeline.AllEntryTypes[hoursAgo%len(timeline.AllEntryTypes)]
		resourceID := fmt.Sprintf("res-%d", hoursAgo)
		entriesByType[entryType] = append(entriesByType[entryType], entryAt(entryType, resourceID, fixedNow.Add(-time.Duration(hoursAgo)*time.Hour)))
		expectedOrder = append(expectedOrder, string(entryType)+":"+resourceID)
	}
	timelineService := newTestService(repositoryReturning(entriesByType))

	collected := make([]string, 0)
	cursor := ""
	pageCount := 0
	for {
		page, serviceError := timelineService.GetTimeline(context.Background(), timeline.ListInput{PatientFHIRID: "patient-1", Limit: "3", Cursor: cursor})
		require.NoError(testingT, serviceError)
		collected = append(collected, entryIdentifiers(page.Entries)...)
		pageCount++
		if page.NextCursor == "" {
			break
		}
		cursor = page.NextCursor
		require.Less(testingT, pageCount, 10)
	}

	assert.Equal(testingT, expectedOrder, collected)
	assert.Equal(testingT, 3, pageCount)
}

func TestGetTimeline_CursorNarrowsFetchUpperBoundAndSkipsTies(testingT *testing.T) {
	tiedMoment := fixedNow.Add(-2 * time.Hour)
	var capturedQuery timeline.FetchQuery
	repository := &mocks.MockRepository{
		ListEntriesFn: func(ctx context.Context, entryType timeline.EntryType, query timeline.FetchQuery) ([]timeline.Entry, error) {
			if entryType != timeline.EntryTypeObservation {
				return nil, nil
			}
			capturedQuery = query
			return []timeline.Entry{
				entryAt(timeline.EntryTypeObservation, "a", tiedMoment),
				entryAt(timeline.EntryTypeObservation, "b", tiedMoment),
				entryAt(timeline.EntryTypeObservation, "c", tiedMoment),
			}, nil
		},
	}
	timelineService := newTestService(repository)

	firstPage, firstError := timelineService.GetTimeline(context.Background(), timeline.ListInput{PatientFHIRID: "patient-1", Limit: "1"})
	require.NoError(testingT, firstError)
	require.Equal(testingT, []string{"observation:a"}, entryIdentifiers(firstPage.Entries))

	secondPage, secondError := timelineService.GetTimeline(context.Background(), timeline.ListInput{PatientFHIRID: "patient-1", Limit: "1", Cursor: firstPage.NextCursor})
	require.NoError(testingT, secondError)

	assert.Equal(testingT, []string{"observation:b"}, entryIdentifiers(secondPage.Entries))
	assert.Equal(testingT, tiedMoment, capturedQuery.Window.To)
}

func TestGetTimeline_PartialFailureReturnsRemainingTypesAndFailedList(testingT *testing.T) {
	repository := &mocks.MockRepository{
		ListEntriesFn: func(ctx context.Context, entryType timeline.EntryType, query timeline.FetchQuery) ([]timeline.Entry, error) {
			switch entryType {
			case timeline.EntryTypeReport, timeline.EntryTypeImaging:
				return nil, errFHIRUnavailable
			case timeline.EntryTypeEncounter:
				return []timeline.Entry{entryAt(timeline.EntryTypeEncounter, "enc-1", fixedNow.Add(-time.Hour))}, nil
			}
			return nil, nil
		},
	}

	page, serviceError := newTestService(repository).GetTimeline(context.Background(), timeline.ListInput{PatientFHIRID: "patient-1"})

	require.NoError(testingT, serviceError)
	assert.Equal(testingT, []string{"encounter:enc-1"}, entryIdentifiers(page.Entries))
	assert.Equal(testingT, []timeline.EntryType{timeline.EntryTypeReport, timeline.EntryTypeImaging}, page.FailedTypes)
}

func TestGetTimeline_AllTypesFailingReturnsServiceUnavailable(testingT *testing.T) {
	repository := &mocks.MockRepository{
		ListEntriesFn: func(ctx context.Context, entryType timeline.EntryType, query timeline.FetchQuery) ([]timeline.Entry, error) {
			return nil, errFHIRUnavailable
		},
	}

	page, serviceError := newTestService(repository).GetTimeline(context.Background(), timeline.ListInput{PatientFHIRID: "patient-1"})

	assert.Nil(testingT, page)
	assert.ErrorIs(testingT, serviceError, apperrors.ErrServiceUnavailable)
}

func TestGetTimeline_PatientNotFound(testingT *testing.T) {
	repository := &mocks.MockRepository{
		PatientExistsFn: func(ctx context.Context, patientFHIRID string) error {
			return fmt.Errorf("wrapped: %w", &healthcare.NotFoundError{ResourceType: "Patient", ResourceID: patientFHIRID})
		},
	}

	page, serviceError := newTestService(repository).GetTimeline(context.Background(), timeline.ListInput{PatientFHIRID: "missing"})

	assert.Nil(testingT, page)
	assert.ErrorIs(testingT, serviceError, apperrors.ErrPatientNotFound)
}

func TestGetTimeline_PatientLookupFailureReturnsServiceUnavailable(testingT *testing.T) {
	repository := &mocks.MockRepository{
		PatientExistsFn: func(ctx context.Context, patientFHIRID string) error {
			return errFHIRUnavailable
		},
	}

	_, serviceError := newTestService(repository).GetTimeline(context.Background(), timeline.ListInput{PatientFHIRID: "patient-1"})

	assert.ErrorIs(testingT, serviceError, apperrors.ErrServiceUnavailable)
}

func TestGetTimeline_InvalidInputsAreRejected(testingT *testing.T) {
	invalidInputs := map[string]timeline.ListInput{
		"invalid from":        {PatientFHIRID: "patient-1", From: "yesterday"},
		"invalid to":          {PatientFHIRID: "patient-1", To: "2026-13-45"},
		"from after to":       {PatientFHIRID: "patient-1", From: "2026-10-02", To: "2026-10-01"},
		"unknown type":        {PatientFHIRID: "patient-1", Types: "encounter,appointment"},
		"zero limit":          {PatientFHIRID: "patient-1", Limit: "0"},
		"limit above maximum": {PatientFHIRID: "patient-1", Limit: "101"},
		"non numeric limit":   {PatientFHIRID: "patient-1", Limit: "many"},
		"garbage cursor":      {PatientFHIRID: "patient-1", Cursor: "not-a-cursor"},
		"missing patient":     {},
	}

	for caseName, input := range invalidInputs {
		testingT.Run(caseName, func(subTest *testing.T) {
			repository := &mocks.MockRepository{
				ListEntriesFn: func(ctx context.Context, entryType timeline.EntryType, query timeline.FetchQuery) ([]timeline.Entry, error) {
					subTest.Fatal("repository must not be called for invalid input")
					return nil, nil
				},
			}

			page, serviceError := newTestService(repository).GetTimeline(context.Background(), input)

			assert.Nil(subTest, page)
			var appError apperrors.AppError
			require.ErrorAs(subTest, serviceError, &appError)
			assert.Equal(subTest, apperrors.ErrBadRequest.HTTPCode, appError.HTTPCode)
		})
	}
}
