package mocks

import (
	"context"

	"github.com/healthcare/backend/internal/modules/stats"
)

type MockStatsRepository struct {
	GetTotalPatientsCountFn   func(contextParameter context.Context) (int, error)
	GetEncountersFn           func(contextParameter context.Context) ([]stats.FHIREncounter, error)
	GetConditionsFn           func(contextParameter context.Context) ([]stats.FHIRCondition, error)
	GetExamModalitiesCountsFn func(contextParameter context.Context) (map[string]int, error)
}

func (mockRepo *MockStatsRepository) GetTotalPatientsCount(contextParameter context.Context) (int, error) {
	if mockRepo.GetTotalPatientsCountFn != nil {
		return mockRepo.GetTotalPatientsCountFn(contextParameter)
	}
	return 100, nil
}

func (mockRepo *MockStatsRepository) GetEncounters(contextParameter context.Context) ([]stats.FHIREncounter, error) {
	if mockRepo.GetEncountersFn != nil {
		return mockRepo.GetEncountersFn(contextParameter)
	}
	return []stats.FHIREncounter{}, nil
}

func (mockRepo *MockStatsRepository) GetConditions(contextParameter context.Context) ([]stats.FHIRCondition, error) {
	if mockRepo.GetConditionsFn != nil {
		return mockRepo.GetConditionsFn(contextParameter)
	}
	return []stats.FHIRCondition{}, nil
}

func (mockRepo *MockStatsRepository) GetExamModalitiesCounts(contextParameter context.Context) (map[string]int, error) {
	if mockRepo.GetExamModalitiesCountsFn != nil {
		return mockRepo.GetExamModalitiesCountsFn(contextParameter)
	}
	return map[string]int{}, nil
}
