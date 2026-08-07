package allergy

import (
	"context"

	"github.com/healthcare/backend/internal/shared/apperrors"
	"github.com/healthcare/backend/internal/shared/validator"
)

type Service interface {
	CreateAllergyIntolerance(ctx context.Context, input CreateAllergyInput) (*Allergy, error)
	GetAllergyIntolerancesByPatient(ctx context.Context, patientFHIRID string) ([]*Allergy, error)
	UpdateAllergyIntolerance(ctx context.Context, fhirResourceID string, input UpdateAllergyInput) (*Allergy, error)
	DeleteAllergyIntolerance(ctx context.Context, fhirResourceID string) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (allergyService *service) CreateAllergyIntolerance(ctx context.Context, input CreateAllergyInput) (*Allergy, error) {
	if fieldViolations := validateAllergyFields(input.PatientFHIRID, input.AllergenCode, input.ClinicalStatus); len(fieldViolations) > 0 {
		return nil, apperrors.InvalidArgument("invalid allergy intolerance input", fieldViolations)
	}

	newAllergy := &Allergy{
		PatientFHIRID:   input.PatientFHIRID,
		AllergenCode:    input.AllergenCode,
		AllergenDisplay: input.AllergenDisplay,
		ClinicalStatus:  normalizeClinicalStatus(input.ClinicalStatus),
		Reaction:        input.Reaction,
	}

	return allergyService.repo.CreateAllergyIntolerance(ctx, newAllergy)
}

func (allergyService *service) UpdateAllergyIntolerance(ctx context.Context, fhirResourceID string, input UpdateAllergyInput) (*Allergy, error) {
	currentAllergy, fetchErr := allergyService.repo.GetAllergyIntoleranceByID(ctx, fhirResourceID)
	if fetchErr != nil {
		return nil, fetchErr
	}

	mergedAllergy := mergeAllergyInput(currentAllergy, input)

	if fieldViolations := validateAllergyFields(mergedAllergy.PatientFHIRID, mergedAllergy.AllergenCode, mergedAllergy.ClinicalStatus); len(fieldViolations) > 0 {
		return nil, apperrors.InvalidArgument("invalid allergy intolerance input", fieldViolations)
	}
	mergedAllergy.ClinicalStatus = normalizeClinicalStatus(mergedAllergy.ClinicalStatus)

	return allergyService.repo.UpdateAllergyIntolerance(ctx, fhirResourceID, mergedAllergy)
}

func (allergyService *service) DeleteAllergyIntolerance(ctx context.Context, fhirResourceID string) error {
	return allergyService.repo.DeleteAllergyIntolerance(ctx, fhirResourceID)
}

func (allergyService *service) GetAllergyIntolerancesByPatient(ctx context.Context, patientFHIRID string) ([]*Allergy, error) {
	return allergyService.repo.GetAllergyIntolerancesByPatient(ctx, patientFHIRID)
}

func validateAllergyFields(patientFHIRID string, allergenCode string, clinicalStatus string) map[string]string {
	fieldViolations := make(map[string]string)
	if patientFHIRID == "" {
		fieldViolations["patient_fhir_id"] = "is required"
	}
	if allergenCode == "" {
		fieldViolations["allergen_code"] = "is required"
	}
	if clinicalStatus != "" && !validator.IsValidClinicalStatus(clinicalStatus) {
		fieldViolations["clinical_status"] = "is not a valid clinical status"
	}
	return fieldViolations
}

func normalizeClinicalStatus(clinicalStatus string) string {
	if clinicalStatus == "" {
		return "active"
	}
	return clinicalStatus
}

func mergeAllergyInput(currentAllergy *Allergy, input UpdateAllergyInput) *Allergy {
	mergedAllergy := *currentAllergy
	if input.AllergenCode != nil {
		mergedAllergy.AllergenCode = *input.AllergenCode
	}
	if input.AllergenDisplay != nil {
		mergedAllergy.AllergenDisplay = *input.AllergenDisplay
	}
	if input.ClinicalStatus != nil {
		mergedAllergy.ClinicalStatus = *input.ClinicalStatus
	}
	if input.Reaction != nil {
		mergedAllergy.Reaction = *input.Reaction
	}
	return &mergedAllergy
}
