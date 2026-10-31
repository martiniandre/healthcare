package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/healthcare/backend/internal/modules/medication"
)

type MockAllergyReader struct {
	Allergies []medication.PatientAllergy
	Err       error
}

func (mockReader *MockAllergyReader) ListPatientAllergies(ctx context.Context, patientFHIRID string) ([]medication.PatientAllergy, error) {
	return mockReader.Allergies, mockReader.Err
}

type MockSafetyCatalog struct {
	Snapshot *medication.CatalogSnapshot
	Err      error
}

func (mockCatalog *MockSafetyCatalog) LoadSnapshot(ctx context.Context) (*medication.CatalogSnapshot, error) {
	if mockCatalog.Err != nil {
		return nil, mockCatalog.Err
	}
	if mockCatalog.Snapshot != nil {
		return mockCatalog.Snapshot, nil
	}
	return medication.NewCatalogSnapshot(nil, nil, nil), nil
}

type MockOverrideRecorder struct {
	RecordErr error
	Recorded  []medication.SafetyOverride
	Linked    map[uuid.UUID]string
	Discarded []uuid.UUID
}

func (mockRecorder *MockOverrideRecorder) RecordOverride(ctx context.Context, override medication.SafetyOverride) error {
	if mockRecorder.RecordErr != nil {
		return mockRecorder.RecordErr
	}
	mockRecorder.Recorded = append(mockRecorder.Recorded, override)
	return nil
}

func (mockRecorder *MockOverrideRecorder) LinkMedicationRequest(ctx context.Context, overrideID uuid.UUID, medicationRequestFHIRID string) error {
	if mockRecorder.Linked == nil {
		mockRecorder.Linked = make(map[uuid.UUID]string)
	}
	mockRecorder.Linked[overrideID] = medicationRequestFHIRID
	return nil
}

func (mockRecorder *MockOverrideRecorder) DiscardOverride(ctx context.Context, overrideID uuid.UUID) error {
	mockRecorder.Discarded = append(mockRecorder.Discarded, overrideID)
	return nil
}

type MockAuditRecorder struct {
	Events []medication.OverrideAuditEvent
	Err    error
}

func (mockRecorder *MockAuditRecorder) RecordOverrideAudit(ctx context.Context, event medication.OverrideAuditEvent) error {
	if mockRecorder.Err != nil {
		return mockRecorder.Err
	}
	mockRecorder.Events = append(mockRecorder.Events, event)
	return nil
}
