package mocks

import (
	"context"

	"github.com/healthcare/backend/internal/modules/timeline"
)

type MockRepository struct {
	PatientExistsFn func(ctx context.Context, patientFHIRID string) error
	ListEntriesFn   func(ctx context.Context, entryType timeline.EntryType, query timeline.FetchQuery) ([]timeline.Entry, error)
}

func (mockRepository *MockRepository) PatientExists(ctx context.Context, patientFHIRID string) error {
	if mockRepository.PatientExistsFn != nil {
		return mockRepository.PatientExistsFn(ctx, patientFHIRID)
	}
	return nil
}

func (mockRepository *MockRepository) ListEntries(ctx context.Context, entryType timeline.EntryType, query timeline.FetchQuery) ([]timeline.Entry, error) {
	if mockRepository.ListEntriesFn != nil {
		return mockRepository.ListEntriesFn(ctx, entryType, query)
	}
	return []timeline.Entry{}, nil
}
