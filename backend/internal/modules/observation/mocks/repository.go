package mocks

import (
	"context"

	"github.com/healthcare/backend/internal/modules/observation"
)

type MockObservationRepository struct {
	CreateObservationFn          func(ctx context.Context, entity *observation.Observation) (*observation.Observation, error)
	GetObservationsByEncounterFn func(ctx context.Context, encounterFHIRID string) ([]*observation.Observation, error)
	GetObservationsByPatientFn   func(ctx context.Context, patientFHIRID string) ([]*observation.Observation, error)
	UpdateObservationFn          func(ctx context.Context, fhirResourceID string, entity *observation.Observation) (*observation.Observation, error)
	DeleteObservationFn          func(ctx context.Context, fhirResourceID string) error
}

func (mockRepo *MockObservationRepository) CreateObservation(ctx context.Context, entity *observation.Observation) (*observation.Observation, error) {
	if mockRepo.CreateObservationFn != nil {
		return mockRepo.CreateObservationFn(ctx, entity)
	}
	return entity, nil
}

func (mockRepo *MockObservationRepository) GetObservationsByEncounter(ctx context.Context, encounterFHIRID string) ([]*observation.Observation, error) {
	if mockRepo.GetObservationsByEncounterFn != nil {
		return mockRepo.GetObservationsByEncounterFn(ctx, encounterFHIRID)
	}
	return []*observation.Observation{}, nil
}

func (mockRepo *MockObservationRepository) GetObservationsByPatient(ctx context.Context, patientFHIRID string) ([]*observation.Observation, error) {
	if mockRepo.GetObservationsByPatientFn != nil {
		return mockRepo.GetObservationsByPatientFn(ctx, patientFHIRID)
	}
	return []*observation.Observation{}, nil
}

func (mockRepo *MockObservationRepository) UpdateObservation(ctx context.Context, fhirResourceID string, entity *observation.Observation) (*observation.Observation, error) {
	if mockRepo.UpdateObservationFn != nil {
		return mockRepo.UpdateObservationFn(ctx, fhirResourceID, entity)
	}
	return entity, nil
}

func (mockRepo *MockObservationRepository) DeleteObservation(ctx context.Context, fhirResourceID string) error {
	if mockRepo.DeleteObservationFn != nil {
		return mockRepo.DeleteObservationFn(ctx, fhirResourceID)
	}
	return nil
}
