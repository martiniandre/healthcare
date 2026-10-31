package timeline

import (
	"context"
	"log/slog"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/healthcare/backend/internal/shared/apperrors"
	"github.com/healthcare/backend/internal/shared/healthcare"
)

const fetchTimeout = 15 * time.Second

type Service interface {
	GetTimeline(ctx context.Context, input ListInput) (*Page, error)
}

type Clock func() time.Time

type service struct {
	repository Repository
	clock      Clock
}

func NewService(repository Repository) Service {
	return NewServiceWithClock(repository, time.Now)
}

func NewServiceWithClock(repository Repository, clock Clock) Service {
	return &service{repository: repository, clock: clock}
}

type listParameters struct {
	patientFHIRID string
	window        Window
	types         []EntryType
	cursor        *cursorPosition
	limit         int
}

type fetchOutcome struct {
	entries []Entry
	failed  bool
}

func (timelineService *service) GetTimeline(ctx context.Context, input ListInput) (*Page, error) {
	parameters, validationError := timelineService.parseInput(input)
	if validationError != nil {
		return nil, validationError
	}

	fetchUpperBound := parameters.window.To
	if parameters.cursor != nil && parameters.cursor.OccurredAt.Before(fetchUpperBound) {
		fetchUpperBound = parameters.cursor.OccurredAt
	}
	fetchQuery := FetchQuery{
		PatientFHIRID: parameters.patientFHIRID,
		Window:        Window{From: parameters.window.From, To: fetchUpperBound},
		MaxEntries:    parameters.limit + 1,
	}

	var waitGroup sync.WaitGroup
	var patientError error
	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()
		patientContext, cancel := context.WithTimeout(ctx, fetchTimeout)
		defer cancel()
		patientError = timelineService.repository.PatientExists(patientContext, parameters.patientFHIRID)
	}()

	outcomes := make([]fetchOutcome, len(parameters.types))
	for typeIndex, entryType := range parameters.types {
		waitGroup.Add(1)
		go func(typeIndex int, entryType EntryType) {
			defer waitGroup.Done()
			fetchContext, cancel := context.WithTimeout(ctx, fetchTimeout)
			defer cancel()
			entries, fetchError := timelineService.repository.ListEntries(fetchContext, entryType, fetchQuery)
			if fetchError != nil {
				slog.Error("failed to fetch timeline entries", "error", fetchError, "entry_type", entryType, "patient_fhir_id", parameters.patientFHIRID)
				outcomes[typeIndex] = fetchOutcome{failed: true}
				return
			}
			outcomes[typeIndex] = fetchOutcome{entries: entries}
		}(typeIndex, entryType)
	}
	waitGroup.Wait()

	if patientError != nil {
		if healthcare.IsNotFound(patientError) {
			return nil, apperrors.ErrPatientNotFound
		}
		slog.Error("failed to verify patient for timeline", "error", patientError, "patient_fhir_id", parameters.patientFHIRID)
		return nil, apperrors.ErrServiceUnavailable
	}

	failedTypes := make([]EntryType, 0)
	merged := make([]Entry, 0)
	for typeIndex, outcome := range outcomes {
		if outcome.failed {
			failedTypes = append(failedTypes, parameters.types[typeIndex])
			continue
		}
		for _, entry := range outcome.entries {
			if !entryIsWithinWindow(entry, parameters.window) {
				continue
			}
			if parameters.cursor != nil && !entryComesAfterCursor(entry, *parameters.cursor) {
				continue
			}
			merged = append(merged, entry)
		}
	}

	if len(failedTypes) == len(parameters.types) {
		return nil, apperrors.ErrServiceUnavailable
	}

	sort.SliceStable(merged, func(firstIndex int, secondIndex int) bool {
		return entryComesBefore(merged[firstIndex], merged[secondIndex])
	})

	page := &Page{Entries: merged, FailedTypes: failedTypes, Window: parameters.window}
	if len(merged) > parameters.limit {
		page.Entries = merged[:parameters.limit]
		page.NextCursor = encodeCursor(page.Entries[parameters.limit-1])
	}
	return page, nil
}

func entryIsWithinWindow(entry Entry, window Window) bool {
	return !entry.OccurredAt.Before(window.From) && !entry.OccurredAt.After(window.To)
}

func (timelineService *service) parseInput(input ListInput) (listParameters, error) {
	violations := make(map[string]string)
	parameters := listParameters{patientFHIRID: strings.TrimSpace(input.PatientFHIRID), limit: DefaultLimit}

	if parameters.patientFHIRID == "" {
		violations["patientFhirId"] = "is required"
	}

	now := timelineService.clock().UTC()
	parameters.window = Window{From: subtractMonths(now, DefaultWindowMonths), To: now}

	if input.From != "" {
		parsedFrom, parsed := parseBoundary(input.From, false)
		if !parsed {
			violations["from"] = "must be an ISO date or RFC3339 timestamp"
		} else {
			parameters.window.From = parsedFrom
		}
	}
	if input.To != "" {
		parsedTo, parsed := parseBoundary(input.To, true)
		if !parsed {
			violations["to"] = "must be an ISO date or RFC3339 timestamp"
		} else {
			parameters.window.To = parsedTo
		}
	}
	if _, hasFromViolation := violations["from"]; !hasFromViolation {
		if _, hasToViolation := violations["to"]; !hasToViolation && parameters.window.From.After(parameters.window.To) {
			violations["from"] = "must not be after to"
		}
	}

	requestedTypes, typesError := parseTypes(input.Types)
	if typesError != "" {
		violations["types"] = typesError
	}
	parameters.types = requestedTypes

	if input.Limit != "" {
		parsedLimit, limitError := strconv.Atoi(input.Limit)
		if limitError != nil || parsedLimit < 1 || parsedLimit > MaxLimit {
			violations["limit"] = "must be an integer between 1 and " + strconv.Itoa(MaxLimit)
		} else {
			parameters.limit = parsedLimit
		}
	}

	if input.Cursor != "" {
		position, cursorError := decodeCursor(input.Cursor)
		if cursorError != nil {
			violations["cursor"] = "is invalid"
		} else {
			parameters.cursor = &position
		}
	}

	if len(violations) > 0 {
		return listParameters{}, apperrors.InvalidArgument("invalid timeline parameters", violations)
	}
	return parameters, nil
}

func parseTypes(rawTypes string) ([]EntryType, string) {
	if strings.TrimSpace(rawTypes) == "" {
		return append([]EntryType(nil), AllEntryTypes...), ""
	}
	requested := make(map[EntryType]bool)
	for _, rawType := range strings.Split(rawTypes, ",") {
		entryType := EntryType(strings.ToLower(strings.TrimSpace(rawType)))
		if !isKnownEntryType(entryType) {
			return nil, "contains unknown type " + strconv.Quote(strings.TrimSpace(rawType))
		}
		requested[entryType] = true
	}
	ordered := make([]EntryType, 0, len(requested))
	for _, knownType := range AllEntryTypes {
		if requested[knownType] {
			ordered = append(ordered, knownType)
		}
	}
	return ordered, ""
}

var boundaryLayouts = []string{time.RFC3339Nano, "2006-01-02T15:04:05", "2006-01-02T15:04"}

func parseBoundary(rawValue string, isUpperBound bool) (time.Time, bool) {
	trimmed := strings.TrimSpace(rawValue)
	if dateOnly, dateError := time.Parse("2006-01-02", trimmed); dateError == nil {
		if isUpperBound {
			return dateOnly.AddDate(0, 0, 1).Add(-time.Nanosecond).UTC(), true
		}
		return dateOnly.UTC(), true
	}
	for _, layout := range boundaryLayouts {
		if parsed, parseError := time.Parse(layout, trimmed); parseError == nil {
			return parsed.UTC(), true
		}
	}
	return time.Time{}, false
}

func subtractMonths(reference time.Time, months int) time.Time {
	firstOfTargetMonth := time.Date(reference.Year(), reference.Month()-time.Month(months), 1, reference.Hour(), reference.Minute(), reference.Second(), reference.Nanosecond(), reference.Location())
	lastDayOfTargetMonth := firstOfTargetMonth.AddDate(0, 1, -1).Day()
	targetDay := reference.Day()
	if targetDay > lastDayOfTargetMonth {
		targetDay = lastDayOfTargetMonth
	}
	return firstOfTargetMonth.AddDate(0, 0, targetDay-1)
}
