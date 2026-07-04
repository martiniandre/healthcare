package mocks

import (
	"context"

	"github.com/healthcare/backend/internal/modules/condition"
)

type MockConditionRepository struct {
	CreateConditionFn        func(ctx context.Context, entity *condition.Condition) (*condition.Condition, error)
	GetConditionsByPatientFn func(ctx context.Context, patientFHIRID string) ([]*condition.Condition, error)
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
