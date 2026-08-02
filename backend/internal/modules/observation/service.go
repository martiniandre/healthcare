package observation

import (
	"context"
	"time"

	"github.com/healthcare/backend/internal/shared/apperrors"
	"github.com/healthcare/backend/internal/shared/validator"
)

type Service interface {
	CreateObservation(ctx context.Context, input CreateObservationInput) (*Observation, error)
	GetObservationsByEncounter(ctx context.Context, encounterFHIRID string) ([]*Observation, error)
	GetObservationsByPatient(ctx context.Context, patientFHIRID string) ([]*Observation, error)
	UpdateObservation(ctx context.Context, fhirResourceID string, observation *Observation) (*Observation, error)
	DeleteObservation(ctx context.Context, fhirResourceID string) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (observationService *service) CreateObservation(ctx context.Context, input CreateObservationInput) (*Observation, error) {
	violations := make(map[string]string)
	if input.EncounterFHIRID == "" {
		violations["encounter_fhir_id"] = "is required"
	}
	if input.PatientFHIRID == "" {
		violations["patient_fhir_id"] = "is required"
	}
	if input.LoincCode == "" {
		violations["loinc_code"] = "is required"
	}
	if len(violations) > 0 {
		return nil, apperrors.InvalidArgument("invalid observation input", violations)
	}
	if !validator.IsValidLOINC(input.LoincCode) {
		return nil, apperrors.InvalidArgument("invalid observation input", map[string]string{"loinc_code": "invalid LOINC format"})
	}
	if !validator.IsValidObservationRange(input.LoincCode, input.ValueQuantity) {
		return nil, apperrors.InvalidArgument("invalid observation input", map[string]string{"value_quantity": "value quantity out of clinical range"})
	}

	observedAt := time.Now()
	if input.ObservedAt != nil && !input.ObservedAt.IsZero() {
		observedAt = *input.ObservedAt
	}

	observation := &Observation{
		EncounterFHIRID: input.EncounterFHIRID,
		PatientFHIRID:   input.PatientFHIRID,
		LoincCode:       input.LoincCode,
		CodeDisplay:     input.CodeDisplay,
		ValueQuantity:   input.ValueQuantity,
		ValueUnit:       input.ValueUnit,
		ObservedAt:      observedAt,
	}

	return observationService.repo.CreateObservation(ctx, observation)
}

func (observationService *service) GetObservationsByEncounter(ctx context.Context, encounterFHIRID string) ([]*Observation, error) {
	return observationService.repo.GetObservationsByEncounter(ctx, encounterFHIRID)
}

func (observationService *service) GetObservationsByPatient(ctx context.Context, patientFHIRID string) ([]*Observation, error) {
	return observationService.repo.GetObservationsByPatient(ctx, patientFHIRID)
}

func (observationService *service) UpdateObservation(ctx context.Context, fhirResourceID string, observation *Observation) (*Observation, error) {
	violations := make(map[string]string)
	if observation.PatientFHIRID == "" {
		violations["patient_fhir_id"] = "is required"
	}
	if observation.EncounterFHIRID == "" {
		violations["encounter_fhir_id"] = "is required"
	}
	if observation.LoincCode == "" {
		violations["loinc_code"] = "is required"
	}
	if len(violations) > 0 {
		return nil, apperrors.InvalidArgument("invalid observation input", violations)
	}
	return observationService.repo.UpdateObservation(ctx, fhirResourceID, observation)
}

func (observationService *service) DeleteObservation(ctx context.Context, fhirResourceID string) error {
	return observationService.repo.DeleteObservation(ctx, fhirResourceID)
}
