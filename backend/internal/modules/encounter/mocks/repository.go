package mocks

import (
	"context"

	"github.com/healthcare/backend/internal/modules/encounter"
)

type MockEncounterRepository struct {
	CreateEncounterFn        func(ctx context.Context, entity *encounter.Encounter) (*encounter.Encounter, error)
	GetEncounterByIDFn       func(ctx context.Context, fhirResourceID string) (*encounter.Encounter, error)
	GetEncountersByPatientFn func(ctx context.Context, patientFHIRID string) ([]*encounter.Encounter, error)
}

func (mockRepo *MockEncounterRepository) CreateEncounter(ctx context.Context, entity *encounter.Encounter) (*encounter.Encounter, error) {
	if mockRepo.CreateEncounterFn != nil {
		return mockRepo.CreateEncounterFn(ctx, entity)
	}
	return entity, nil
}

func (mockRepo *MockEncounterRepository) GetEncounterByID(ctx context.Context, fhirResourceID string) (*encounter.Encounter, error) {
	if mockRepo.GetEncounterByIDFn != nil {
		return mockRepo.GetEncounterByIDFn(ctx, fhirResourceID)
	}
	return &encounter.Encounter{}, nil
}

func (mockRepo *MockEncounterRepository) GetEncountersByPatient(ctx context.Context, patientFHIRID string) ([]*encounter.Encounter, error) {
	if mockRepo.GetEncountersByPatientFn != nil {
		return mockRepo.GetEncountersByPatientFn(ctx, patientFHIRID)
	}
	return []*encounter.Encounter{}, nil
}
