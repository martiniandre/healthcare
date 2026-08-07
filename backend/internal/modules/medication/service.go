package medication

import (
	"context"

	"github.com/healthcare/backend/internal/shared/apperrors"
)

type Service interface {
	CreateMedicationRequest(ctx context.Context, input CreateMedicationInput) (*Medication, error)
	GetMedicationRequestsByEncounter(ctx context.Context, encounterFHIRID string) ([]*Medication, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (medicationService *service) CreateMedicationRequest(ctx context.Context, input CreateMedicationInput) (*Medication, error) {
	if fieldViolations := validateMedicationFields(input); len(fieldViolations) > 0 {
		return nil, apperrors.InvalidArgument("invalid medication request input", fieldViolations)
	}

	newMedication := &Medication{
		EncounterFHIRID:    input.EncounterFHIRID,
		PatientFHIRID:      input.PatientFHIRID,
		PractitionerFHIRID: input.PractitionerFHIRID,
		MedicationCode:     input.MedicationCode,
		MedicationName:     input.MedicationName,
		DosageInstructions: input.DosageInstructions,
		Status:             "active",
	}

	return medicationService.repo.CreateMedicationRequest(ctx, newMedication)
}

func (medicationService *service) GetMedicationRequestsByEncounter(ctx context.Context, encounterFHIRID string) ([]*Medication, error) {
	return medicationService.repo.GetMedicationRequestsByEncounter(ctx, encounterFHIRID)
}

func validateMedicationFields(input CreateMedicationInput) map[string]string {
	fieldViolations := make(map[string]string)
	if input.EncounterFHIRID == "" {
		fieldViolations["encounter_fhir_id"] = "is required"
	}
	if input.PatientFHIRID == "" {
		fieldViolations["patient_fhir_id"] = "is required"
	}
	if input.MedicationCode == "" && input.MedicationName == "" {
		fieldViolations["medication"] = "either name or code is required"
	}
	return fieldViolations
}
