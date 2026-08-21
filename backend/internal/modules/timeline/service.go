package timeline

import (
	"context"
	"errors"
	"sort"
	"sync"
	"time"

	"github.com/healthcare/backend/internal/shared/apperrors"
	"github.com/healthcare/backend/internal/shared/healthcare"
)

var ErrPatientNotFound = errors.New("patient not found")

const (
	DefaultPageSize = 20
	MaxPageSize     = 100
)

var SupportedResourceTypes = []string{
	"Encounter",
	"Observation",
	"Condition",
	"MedicationRequest",
	"DiagnosticReport",
	"ImagingStudy",
	"AllergyIntolerance",
}

type TimelineFilter struct {
	Types  []string
	Status string
	Before *time.Time
	Limit  int
}

type Service interface {
	GetTimeline(ctx context.Context, patientFHIRID string, filter TimelineFilter) (*TimelinePage, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (timelineService *service) GetTimeline(ctx context.Context, patientFHIRID string, filter TimelineFilter) (*TimelinePage, error) {
	selectedTypes, filterErr := normalizeFilter(&filter)
	if filterErr != nil {
		return nil, filterErr
	}

	if patientErr := timelineService.repo.PatientExists(ctx, patientFHIRID); patientErr != nil {
		if healthcare.IsNotFound(patientErr) {
			return nil, apperrors.ErrPatientNotFound
		}
		return nil, apperrors.ErrInternalServer
	}

	results := make([][]TimelineEntry, len(selectedTypes))
	failures := make([]error, len(selectedTypes))

	var fetchGroup sync.WaitGroup
	for typeIndex, resourceType := range selectedTypes {
		fetchGroup.Add(1)
		go func(typeIndex int, resourceType string) {
			defer fetchGroup.Done()
			entries, fetchErr := timelineService.fetchType(ctx, resourceType, patientFHIRID, filter)
			if fetchErr != nil {
				failures[typeIndex] = fetchErr
				return
			}
			results[typeIndex] = entries
		}(typeIndex, resourceType)
	}
	fetchGroup.Wait()

	unavailableTypes := make([]string, 0)
	for typeIndex, resourceType := range selectedTypes {
		if failures[typeIndex] != nil {
			unavailableTypes = append(unavailableTypes, resourceType)
		}
	}

	if len(unavailableTypes) == len(selectedTypes) {
		return nil, apperrors.ErrServiceUnavailable
	}

	mergedEntries := mergeEntries(results)
	sort.SliceStable(mergedEntries, func(firstIndex, secondIndex int) bool {
		firstEntry := mergedEntries[firstIndex]
		secondEntry := mergedEntries[secondIndex]
		if !firstEntry.RecordedAt.Equal(secondEntry.RecordedAt) {
			return firstEntry.RecordedAt.After(secondEntry.RecordedAt)
		}
		if firstEntry.ResourceType != secondEntry.ResourceType {
			return firstEntry.ResourceType < secondEntry.ResourceType
		}
		return firstEntry.FHIRResourceID < secondEntry.FHIRResourceID
	})

	page := &TimelinePage{
		Entries:          mergedEntries,
		NextCursor:       nil,
		UnavailableTypes: unavailableTypes,
	}

	if len(mergedEntries) > filter.Limit {
		page.Entries = mergedEntries[:filter.Limit]
		lastEntry := page.Entries[len(page.Entries)-1]
		nextCursor := lastEntry.RecordedAt
		page.NextCursor = &nextCursor
	} else if len(mergedEntries) == filter.Limit && len(mergedEntries) > 0 {
		lastEntry := page.Entries[len(page.Entries)-1]
		nextCursor := lastEntry.RecordedAt
		page.NextCursor = &nextCursor
	}

	if page.Entries == nil {
		page.Entries = []TimelineEntry{}
	}

	return page, nil
}

func (timelineService *service) fetchType(ctx context.Context, resourceType string, patientFHIRID string, filter TimelineFilter) ([]TimelineEntry, error) {
	switch resourceType {
	case "Encounter":
		return timelineService.repo.FetchEncounters(ctx, patientFHIRID, filter.Before, filter.Limit)
	case "Observation":
		return timelineService.repo.FetchObservations(ctx, patientFHIRID, filter.Before, filter.Limit)
	case "Condition":
		return timelineService.repo.FetchConditions(ctx, patientFHIRID, filter.Status == "active", filter.Before, filter.Limit)
	case "MedicationRequest":
		return timelineService.repo.FetchMedications(ctx, patientFHIRID, filter.Status == "active", filter.Before, filter.Limit)
	case "DiagnosticReport":
		return timelineService.repo.FetchReports(ctx, patientFHIRID, filter.Before, filter.Limit)
	case "ImagingStudy":
		return timelineService.repo.FetchImaging(ctx, patientFHIRID, filter.Before, filter.Limit)
	case "AllergyIntolerance":
		return timelineService.repo.FetchAllergies(ctx, patientFHIRID, filter.Before, filter.Limit)
	default:
		return nil, apperrors.InvalidArgument("unsupported resource type", map[string]string{"types": resourceType})
	}
}

func normalizeFilter(filter *TimelineFilter) ([]string, error) {
	if filter.Limit <= 0 {
		filter.Limit = DefaultPageSize
	}
	if filter.Limit > MaxPageSize {
		return nil, apperrors.InvalidArgument("limit exceeds maximum page size", map[string]string{"limit": "must be between 1 and 100"})
	}

	if filter.Status != "" && filter.Status != "active" {
		return nil, apperrors.InvalidArgument("unsupported status filter", map[string]string{"status": "only 'active' is supported"})
	}

	supportedSet := make(map[string]bool, len(SupportedResourceTypes))
	for _, resourceType := range SupportedResourceTypes {
		supportedSet[resourceType] = true
	}

	if len(filter.Types) == 0 {
		return SupportedResourceTypes, nil
	}

	selectedTypes := make([]string, 0, len(filter.Types))
	for _, requestedType := range filter.Types {
		if !supportedSet[requestedType] {
			return nil, apperrors.InvalidArgument("unsupported resource type", map[string]string{"types": requestedType})
		}
		if !containsType(selectedTypes, requestedType) {
			selectedTypes = append(selectedTypes, requestedType)
		}
	}

	return selectedTypes, nil
}

func containsType(resourceTypes []string, resourceType string) bool {
	for _, existingType := range resourceTypes {
		if existingType == resourceType {
			return true
		}
	}
	return false
}

func mergeEntries(results [][]TimelineEntry) []TimelineEntry {
	totalLength := 0
	for _, typeEntries := range results {
		totalLength += len(typeEntries)
	}

	mergedEntries := make([]TimelineEntry, 0, totalLength)
	for _, typeEntries := range results {
		mergedEntries = append(mergedEntries, typeEntries...)
	}

	return mergedEntries
}
