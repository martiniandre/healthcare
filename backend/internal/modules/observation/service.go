package observation

import (
	"context"
	"errors"

	"github.com/healthcare/backend/internal/shared/validator"
)

type Service interface {
	CreateObservation(ctx context.Context, observation *Observation) (*Observation, error)
	GetObservationsByEncounter(ctx context.Context, encounterFHIRID string) ([]*Observation, error)
	GetObservationsByPatient(ctx context.Context, patientFHIRID string) ([]*Observation, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (observationService *service) CreateObservation(ctx context.Context, observation *Observation) (*Observation, error) {
	if observation.EncounterFHIRID == "" || observation.PatientFHIRID == "" {
		return nil, ErrObservationNotFound
	}
	if !validator.IsValidLOINC(observation.LoincCode) {
		return nil, errors.New("invalid LOINC format")
	}
	if !validator.IsValidObservationRange(observation.LoincCode, observation.ValueQuantity) {
		return nil, errors.New("LOINC value quantity out of clinical range")
	}
	return observationService.repo.CreateObservation(ctx, observation)
}

func (observationService *service) GetObservationsByEncounter(ctx context.Context, encounterFHIRID string) ([]*Observation, error) {
	return observationService.repo.GetObservationsByEncounter(ctx, encounterFHIRID)
}

func (observationService *service) GetObservationsByPatient(ctx context.Context, patientFHIRID string) ([]*Observation, error) {
	return observationService.repo.GetObservationsByPatient(ctx, patientFHIRID)
}
