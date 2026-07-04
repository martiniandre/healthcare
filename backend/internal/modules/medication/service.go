package medication

import (
	"context"
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
	if medication.PatientFHIRID == "" || medication.MedicationCode == "" {
		return nil, ErrMedicationRequestNotFound
	}
	if medication.Status == "" {
		medication.Status = "active"
	}
	return medicationService.repo.CreateMedicationRequest(ctx, medication)
}

func (medicationService *service) GetMedicationRequestsByEncounter(ctx context.Context, encounterFHIRID string) ([]*Medication, error) {
	return medicationService.repo.GetMedicationRequestsByEncounter(ctx, encounterFHIRID)
}
