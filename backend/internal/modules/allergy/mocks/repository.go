package mocks

import (
	"context"

	"github.com/healthcare/backend/internal/modules/allergy"
)

type MockAllergyRepository struct {
	CreateAllergyIntoleranceFn        func(ctx context.Context, entity *allergy.Allergy) (*allergy.Allergy, error)
	GetAllergyIntolerancesByPatientFn func(ctx context.Context, patientFHIRID string) ([]*allergy.Allergy, error)
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
