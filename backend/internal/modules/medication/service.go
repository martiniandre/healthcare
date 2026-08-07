package medication

import (
	"context"

	"github.com/healthcare/backend/internal/shared/apperrors"
)

type Service interface {
	CreateMedicationRequest(ctx context.Context, medication *Medication) (*Medication, error)
	GetMedicationRequestsByEncounter(ctx context.Context, encounterFHIRID string) ([]*Medication, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (medicationService *service) CreateMedicationRequest(ctx context.Context, medication *Medication) (*Medication, error) {
	if medication.PatientFHIRID == "" {
		return nil, apperrors.ErrBadRequest.WithFields(map[string]string{"patient_fhir_id": "required"})
	}
	if medication.MedicationName == "" && medication.MedicationCode == "" {
		return nil, apperrors.ErrBadRequest.WithFields(map[string]string{"medication": "either name or code is required"})
	}
	if medication.Status == "" {
		medication.Status = "active"
	}
	return medicationService.repo.CreateMedicationRequest(ctx, medication)
}

func (medicationService *service) GetMedicationRequestsByEncounter(ctx context.Context, encounterFHIRID string) ([]*Medication, error) {
	return medicationService.repo.GetMedicationRequestsByEncounter(ctx, encounterFHIRID)
}
