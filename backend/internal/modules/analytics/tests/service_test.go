package tests

import (
	"context"
	"errors"
	"testing"

	"github.com/healthcare/backend/internal/modules/analytics"
	"github.com/healthcare/backend/internal/modules/analytics/mocks"
	"github.com/stretchr/testify/assert"
)

var errorRepositoryFailure = errors.New("repository query failed")

func TestGetStats_Success(testingT *testing.T) {
	analyticsService := analytics.NewService(&mocks.MockStatsRepository{
		GetTotalPatientsCountFn: func(contextParameter context.Context) (int, error) {
			return 250, nil
		},
		GetExamModalitiesCountsFn: func(contextParameter context.Context) (map[string]int, error) {
			return map[string]int{
				"CT": 10,
				"MR": 20,
			}, nil
		},
	})

	analyticsData, errorInstance := analyticsService.GetStats(context.Background())

	assert.NoError(testingT, errorInstance)
	assert.Equal(testingT, 250, analyticsData.TotalPatients)
	assert.Equal(testingT, 0.0, analyticsData.AvgServiceDurationMinutes)
	assert.Len(testingT, analyticsData.WeeklyConsultations, 7)
	assert.Len(testingT, analyticsData.ExamModalities, 4)
	assert.Len(testingT, analyticsData.PathologyCases, 0)
}

func TestGetStats_TotalPatientsCountReturnsActualValue(testingT *testing.T) {
	analyticsService := analytics.NewService(&mocks.MockStatsRepository{
		GetTotalPatientsCountFn: func(contextParameter context.Context) (int, error) {
			return 1, nil
		},
	})

	analyticsData, errorInstance := analyticsService.GetStats(context.Background())

	assert.NoError(testingT, errorInstance)
	assert.Equal(testingT, 1, analyticsData.TotalPatients)
}

func TestGetStats_GetTotalPatientsError_ReturnsError(testingT *testing.T) {
	analyticsService := analytics.NewService(&mocks.MockStatsRepository{
		GetTotalPatientsCountFn: func(contextParameter context.Context) (int, error) {
			return 0, errorRepositoryFailure
		},
	})

	analyticsData, errorInstance := analyticsService.GetStats(context.Background())

	assert.Error(testingT, errorInstance)
	assert.Equal(testingT, 0, analyticsData.TotalPatients)
}

func TestGetStats_GetEncountersError_ReturnsError(testingT *testing.T) {
	analyticsService := analytics.NewService(&mocks.MockStatsRepository{
		GetEncountersFn: func(contextParameter context.Context) ([]analytics.FHIREncounter, error) {
			return nil, errorRepositoryFailure
		},
	})

	_, errorInstance := analyticsService.GetStats(context.Background())

	assert.Error(testingT, errorInstance)
}

func TestGetStats_CalculatesDurationCorrectly(testingT *testing.T) {
	analyticsService := analytics.NewService(&mocks.MockStatsRepository{
		GetEncountersFn: func(contextParameter context.Context) ([]analytics.FHIREncounter, error) {
			return []analytics.FHIREncounter{
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

	analyticsData, errorInstance := analyticsService.GetStats(context.Background())

	assert.NoError(testingT, errorInstance)
	assert.Equal(testingT, 25.0, analyticsData.AvgServiceDurationMinutes)
}

func TestGetStats_PathologyCasesWithConditions(testingT *testing.T) {
	analyticsService := analytics.NewService(&mocks.MockStatsRepository{
		GetConditionsFn: func(contextParameter context.Context) ([]analytics.FHIRCondition, error) {
			return []analytics.FHIRCondition{
				{ID: "1", ICD10Code: "J45.9"},
				{ID: "2", ICD10Code: "I10"},
				{ID: "3", ICD10Code: "I10"},
				{ID: "4", ICD10Code: "E11.9"},
			}, nil
		},
	})

	analyticsData, errorInstance := analyticsService.GetStats(context.Background())

	assert.NoError(testingT, errorInstance)
	assert.Len(testingT, analyticsData.PathologyCases, 3)
}

func TestGetStats_PathologyCasesEmptyWhenNoConditions(testingT *testing.T) {
	analyticsService := analytics.NewService(&mocks.MockStatsRepository{})

	analyticsData, errorInstance := analyticsService.GetStats(context.Background())

	assert.NoError(testingT, errorInstance)
	assert.Len(testingT, analyticsData.PathologyCases, 0)
}
