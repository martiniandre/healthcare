package allergy

import (
	"context"

	"github.com/healthcare/backend/internal/shared/validator"
)

type Service interface {
	CreateAllergyIntolerance(ctx context.Context, allergy *Allergy) (*Allergy, error)
	GetAllergyIntolerancesByPatient(ctx context.Context, patientFHIRID string) ([]*Allergy, error)
	UpdateAllergyIntolerance(ctx context.Context, fhirResourceID string, allergy *Allergy) (*Allergy, error)
	DeleteAllergyIntolerance(ctx context.Context, fhirResourceID string) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (allergyService *service) CreateAllergyIntolerance(ctx context.Context, allergy *Allergy) (*Allergy, error) {
	if allergy.PatientFHIRID == "" || allergy.AllergenCode == "" {
		return nil, ErrAllergyNotFound
	}
	if allergy.ClinicalStatus == "" {
		allergy.ClinicalStatus = "active"
	} else if !validator.IsValidClinicalStatus(allergy.ClinicalStatus) {
		return nil, ErrAllergyNotFound
	}
	return allergyService.repo.CreateAllergyIntolerance(ctx, allergy)
}

func (allergyService *service) UpdateAllergyIntolerance(ctx context.Context, fhirResourceID string, allergy *Allergy) (*Allergy, error) {
	if allergy.PatientFHIRID == "" || allergy.AllergenCode == "" {
		return nil, ErrAllergyNotFound
	}
	if allergy.ClinicalStatus == "" {
		allergy.ClinicalStatus = "active"
	} else if !validator.IsValidClinicalStatus(allergy.ClinicalStatus) {
		return nil, ErrAllergyNotFound
	}
	return allergyService.repo.UpdateAllergyIntolerance(ctx, fhirResourceID, allergy)
}

func (allergyService *service) DeleteAllergyIntolerance(ctx context.Context, fhirResourceID string) error {
	return allergyService.repo.DeleteAllergyIntolerance(ctx, fhirResourceID)
}

func (allergyService *service) GetAllergyIntolerancesByPatient(ctx context.Context, patientFHIRID string) ([]*Allergy, error) {
	return allergyService.repo.GetAllergyIntolerancesByPatient(ctx, patientFHIRID)
}
