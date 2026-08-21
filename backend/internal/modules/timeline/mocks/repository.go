package mocks

import (
	"context"
	"time"

	"github.com/healthcare/backend/internal/modules/timeline"
)

type MockRepository struct {
	Entries     map[string][]timeline.TimelineEntry
	FetchErrors map[string]error
	PatientErr  error

	RequestedTypes    []string
	ReceivedOnlyActive bool
	ReceivedBefore    *time.Time
	ReceivedLimit     int
}

func NewMockRepository() *MockRepository {
	return &MockRepository{
		Entries:     make(map[string][]timeline.TimelineEntry),
		FetchErrors: make(map[string]error),
	}
}

func (mockRepository *MockRepository) PatientExists(ctx context.Context, patientFHIRID string) error {
	return mockRepository.PatientErr
}

func (mockRepository *MockRepository) FetchEncounters(ctx context.Context, patientFHIRID string, before *time.Time, limit int) ([]timeline.TimelineEntry, error) {
	return mockRepository.fetch("Encounter", before, limit)
}

func (mockRepository *MockRepository) FetchObservations(ctx context.Context, patientFHIRID string, before *time.Time, limit int) ([]timeline.TimelineEntry, error) {
	return mockRepository.fetch("Observation", before, limit)
}

func (mockRepository *MockRepository) FetchConditions(ctx context.Context, patientFHIRID string, onlyActive bool, before *time.Time, limit int) ([]timeline.TimelineEntry, error) {
	mockRepository.ReceivedOnlyActive = onlyActive
	return mockRepository.fetch("Condition", before, limit)
}

func (mockRepository *MockRepository) FetchMedications(ctx context.Context, patientFHIRID string, onlyActive bool, before *time.Time, limit int) ([]timeline.TimelineEntry, error) {
	mockRepository.ReceivedOnlyActive = onlyActive
	return mockRepository.fetch("MedicationRequest", before, limit)
}

func (mockRepository *MockRepository) FetchReports(ctx context.Context, patientFHIRID string, before *time.Time, limit int) ([]timeline.TimelineEntry, error) {
	return mockRepository.fetch("DiagnosticReport", before, limit)
}

func (mockRepository *MockRepository) FetchImaging(ctx context.Context, patientFHIRID string, before *time.Time, limit int) ([]timeline.TimelineEntry, error) {
	return mockRepository.fetch("ImagingStudy", before, limit)
}

func (mockRepository *MockRepository) FetchAllergies(ctx context.Context, patientFHIRID string, before *time.Time, limit int) ([]timeline.TimelineEntry, error) {
	return mockRepository.fetch("AllergyIntolerance", before, limit)
}

func (mockRepository *MockRepository) fetch(resourceType string, before *time.Time, limit int) ([]timeline.TimelineEntry, error) {
	mockRepository.RequestedTypes = append(mockRepository.RequestedTypes, resourceType)
	mockRepository.ReceivedBefore = before
	mockRepository.ReceivedLimit = limit

	if fetchError, hasError := mockRepository.FetchErrors[resourceType]; hasError {
		return nil, fetchError
	}

	return mockRepository.Entries[resourceType], nil
}
