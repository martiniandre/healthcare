package mocks

import (
	"context"

	"github.com/healthcare/backend/internal/modules/condition"
)

type MockConditionRepository struct {
	CreateConditionFn        func(ctx context.Context, entity *condition.Condition) (*condition.Condition, error)
	GetConditionsByPatientFn func(ctx context.Context, patientFHIRID string) ([]*condition.Condition, error)
	UpdateConditionFn        func(ctx context.Context, fhirResourceID string, entity *condition.Condition) (*condition.Condition, error)
	DeleteConditionFn        func(ctx context.Context, fhirResourceID string) error
}

func (mockRepo *MockConditionRepository) CreateCondition(ctx context.Context, entity *condition.Condition) (*condition.Condition, error) {
	if mockRepo.CreateConditionFn != nil {
		return mockRepo.CreateConditionFn(ctx, entity)
	}
	return entity, nil
}

func (mockRepo *MockConditionRepository) GetConditionsByPatient(ctx context.Context, patientFHIRID string) ([]*condition.Condition, error) {
	if mockRepo.GetConditionsByPatientFn != nil {
		return mockRepo.GetConditionsByPatientFn(ctx, patientFHIRID)
	}
	return []*condition.Condition{}, nil
}

func (mockRepo *MockConditionRepository) UpdateCondition(ctx context.Context, fhirResourceID string, entity *condition.Condition) (*condition.Condition, error) {
	if mockRepo.UpdateConditionFn != nil {
		return mockRepo.UpdateConditionFn(ctx, fhirResourceID, entity)
	}
	return entity, nil
}

func (mockRepo *MockConditionRepository) DeleteCondition(ctx context.Context, fhirResourceID string) error {
	if mockRepo.DeleteConditionFn != nil {
		return mockRepo.DeleteConditionFn(ctx, fhirResourceID)
	}
	return nil
}
