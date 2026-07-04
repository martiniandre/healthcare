package mocks

import (
	"context"

	"github.com/healthcare/backend/internal/modules/medication"
)

type MockMedicationRepository struct {
	CreateMedicationRequestFn          func(ctx context.Context, entity *medication.Medication) (*medication.Medication, error)
	GetMedicationRequestsByEncounterFn func(ctx context.Context, encounterFHIRID string) ([]*medication.Medication, error)
}

func (mockRepo *MockMedicationRepository) CreateMedicationRequest(ctx context.Context, entity *medication.Medication) (*medication.Medication, error) {
	if mockRepo.CreateMedicationRequestFn != nil {
		return mockRepo.CreateMedicationRequestFn(ctx, entity)
	}
	return entity, nil
}

func (mockRepo *MockMedicationRepository) GetMedicationRequestsByEncounter(ctx context.Context, encounterFHIRID string) ([]*medication.Medication, error) {
	if mockRepo.GetMedicationRequestsByEncounterFn != nil {
		return mockRepo.GetMedicationRequestsByEncounterFn(ctx, encounterFHIRID)
	}
	return []*medication.Medication{}, nil
}
