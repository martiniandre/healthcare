package tests

import (
	"context"
	"errors"
	"testing"

	"github.com/healthcare/backend/internal/modules/stats"
	"github.com/healthcare/backend/internal/modules/stats/mocks"
	"github.com/stretchr/testify/assert"
)

var errorRepositoryFailure = errors.New("repository query failed")

func TestGetStats_Success(testingT *testing.T) {
	statsService := stats.NewService(&mocks.MockStatsRepository{
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
	statsService := stats.NewService(&mocks.MockStatsRepository{
		GetTotalPatientsCountFn: func(contextParameter context.Context) (int, error) {
			return 1, nil
		},
	})

	statsData, errorInstance := statsService.GetStats(context.Background())

	assert.NoError(testingT, errorInstance)
	assert.Equal(testingT, 340, statsData.TotalPatients)
}

func TestGetStats_GetTotalPatientsError_ReturnsError(testingT *testing.T) {
	statsService := stats.NewService(&mocks.MockStatsRepository{
		GetTotalPatientsCountFn: func(contextParameter context.Context) (int, error) {
			return 0, errorRepositoryFailure
		},
	})

	statsData, errorInstance := statsService.GetStats(context.Background())

	assert.Error(testingT, errorInstance)
	assert.Equal(testingT, 0, statsData.TotalPatients)
}

func TestGetStats_GetEncountersError_ReturnsError(testingT *testing.T) {
	statsService := stats.NewService(&mocks.MockStatsRepository{
		GetEncountersFn: func(contextParameter context.Context) ([]stats.FHIREncounter, error) {
			return nil, errorRepositoryFailure
		},
	})

	_, errorInstance := statsService.GetStats(context.Background())

	assert.Error(testingT, errorInstance)
}

func TestGetStats_CalculatesDurationCorrectly(testingT *testing.T) {
	statsService := stats.NewService(&mocks.MockStatsRepository{
		GetEncountersFn: func(contextParameter context.Context) ([]stats.FHIREncounter, error) {
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
