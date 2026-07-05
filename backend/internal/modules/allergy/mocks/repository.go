package mocks

import (
	"context"

	"github.com/healthcare/backend/internal/modules/allergy"
)

type MockAllergyRepository struct {
	CreateAllergyIntoleranceFn        func(ctx context.Context, entity *allergy.Allergy) (*allergy.Allergy, error)
	GetAllergyIntolerancesByPatientFn func(ctx context.Context, patientFHIRID string) ([]*allergy.Allergy, error)
	UpdateAllergyIntoleranceFn        func(ctx context.Context, fhirResourceID string, entity *allergy.Allergy) (*allergy.Allergy, error)
	DeleteAllergyIntoleranceFn        func(ctx context.Context, fhirResourceID string) error
}

func (mockRepo *MockAllergyRepository) CreateAllergyIntolerance(ctx context.Context, entity *allergy.Allergy) (*allergy.Allergy, error) {
	if mockRepo.CreateAllergyIntoleranceFn != nil {
		return mockRepo.CreateAllergyIntoleranceFn(ctx, entity)
	}
	return entity, nil
}

func (mockRepo *MockAllergyRepository) GetAllergyIntolerancesByPatient(ctx context.Context, patientFHIRID string) ([]*allergy.Allergy, error) {
	if mockRepo.GetAllergyIntolerancesByPatientFn != nil {
		return mockRepo.GetAllergyIntolerancesByPatientFn(ctx, patientFHIRID)
	}
	return []*allergy.Allergy{}, nil
}

func (mockRepo *MockAllergyRepository) UpdateAllergyIntolerance(ctx context.Context, fhirResourceID string, entity *allergy.Allergy) (*allergy.Allergy, error) {
	if mockRepo.UpdateAllergyIntoleranceFn != nil {
		return mockRepo.UpdateAllergyIntoleranceFn(ctx, fhirResourceID, entity)
	}
	return entity, nil
}

func (mockRepo *MockAllergyRepository) DeleteAllergyIntolerance(ctx context.Context, fhirResourceID string) error {
	if mockRepo.DeleteAllergyIntoleranceFn != nil {
		return mockRepo.DeleteAllergyIntoleranceFn(ctx, fhirResourceID)
	}
	return nil
}
