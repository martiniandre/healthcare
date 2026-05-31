package stats_test

import (
	"context"
	"errors"
	"testing"

	"github.com/healthcare/backend/internal/modules/stats"
	"github.com/stretchr/testify/assert"
)

var errorRepositoryFailure = errors.New("repository query failed")

type mockRepository struct {
	getTotalPatientsCountFn   func(contextParameter context.Context) (int, error)
	getEncountersFn           func(contextParameter context.Context) ([]stats.FHIREncounter, error)
	getConditionsFn           func(contextParameter context.Context) ([]stats.FHIRCondition, error)
	getExamModalitiesCountsFn func(contextParameter context.Context) (map[string]int, error)
}

func (mockRepo *mockRepository) GetTotalPatientsCount(contextParameter context.Context) (int, error) {
	if mockRepo.getTotalPatientsCountFn != nil {
		return mockRepo.getTotalPatientsCountFn(contextParameter)
	}
	return 100, nil
}

func (mockRepo *mockRepository) GetEncounters(contextParameter context.Context) ([]stats.FHIREncounter, error) {
	if mockRepo.getEncountersFn != nil {
		return mockRepo.getEncountersFn(contextParameter)
	}
	return []stats.FHIREncounter{}, nil
}

func (mockRepo *mockRepository) GetConditions(contextParameter context.Context) ([]stats.FHIRCondition, error) {
	if mockRepo.getConditionsFn != nil {
		return mockRepo.getConditionsFn(contextParameter)
	}
	return []stats.FHIRCondition{}, nil
}

func (mockRepo *mockRepository) GetExamModalitiesCounts(contextParameter context.Context) (map[string]int, error) {
	if mockRepo.getExamModalitiesCountsFn != nil {
		return mockRepo.getExamModalitiesCountsFn(contextParameter)
	}
	return map[string]int{}, nil
}

func TestGetStats_Success(testingT *testing.T) {
	statsService := stats.NewService(&mockRepository{
		getTotalPatientsCountFn: func(contextParameter context.Context) (int, error) {
			return 250, nil
		},
		getExamModalitiesCountsFn: func(contextParameter context.Context) (map[string]int, error) {
			return map[string]int{
				"CT": 10,
				"MR": 20,
			}, nil
		},
	})

	statsData, errorInstance := statsService.GetStats(context.Background())

	assert.NoError(testingT, errorInstance)
	assert.Equal(testingT, 250, statsData.TotalPatients)
	assert.Equal(testingT, 99.4, statsData.FHIRComplianceRate)
	assert.Equal(testingT, 14.5, statsData.AvgServiceDurationMinutes)
	assert.Len(testingT, statsData.WeeklyConsultations, 7)
	assert.Len(testingT, statsData.ExamModalities, 4)
	assert.Len(testingT, statsData.PathologyCases, 3)
}

func TestGetStats_TotalPatientsCountFallback(testingT *testing.T) {
	statsService := stats.NewService(&mockRepository{
		getTotalPatientsCountFn: func(contextParameter context.Context) (int, error) {
			return 1, nil
		},
	})

	statsData, errorInstance := statsService.GetStats(context.Background())

	assert.NoError(testingT, errorInstance)
	assert.Equal(testingT, 340, statsData.TotalPatients)
}

func TestGetStats_GetTotalPatientsError_ReturnsError(testingT *testing.T) {
	statsService := stats.NewService(&mockRepository{
		getTotalPatientsCountFn: func(contextParameter context.Context) (int, error) {
			return 0, errorRepositoryFailure
		},
	})

	statsData, errorInstance := statsService.GetStats(context.Background())

	assert.Error(testingT, errorInstance)
	assert.True(testingT, errors.Is(errorInstance, errorRepositoryFailure))
	assert.Equal(testingT, 0, statsData.TotalPatients)
}

func TestGetStats_GetEncountersError_ReturnsError(testingT *testing.T) {
	statsService := stats.NewService(&mockRepository{
		getEncountersFn: func(contextParameter context.Context) ([]stats.FHIREncounter, error) {
			return nil, errorRepositoryFailure
		},
	})

	_, errorInstance := statsService.GetStats(context.Background())

	assert.Error(testingT, errorInstance)
	assert.True(testingT, errors.Is(errorInstance, errorRepositoryFailure))
}

func TestGetStats_CalculatesDurationCorrectly(testingT *testing.T) {
	statsService := stats.NewService(&mockRepository{
		getEncountersFn: func(contextParameter context.Context) ([]stats.FHIREncounter, error) {
			return []stats.FHIREncounter{
				{
					StartedAt: "2026-05-30T10:00:00Z",
					EndedAt:   "2026-05-30T10:20:00Z",
				},
				{
					StartedAt: "2026-05-30T11:00:00Z",
					EndedAt:   "2026-05-30T11:30:00Z",
				},
			}, nil
		},
	})

	statsData, errorInstance := statsService.GetStats(context.Background())

	assert.NoError(testingT, errorInstance)
	assert.Equal(testingT, 25.0, statsData.AvgServiceDurationMinutes)
}
