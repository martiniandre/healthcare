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

func TestGetDashboardData_Success(testingT *testing.T) {
	analyticsService := analytics.NewService(&mocks.MockStatsRepository{
		GetConsultationsPerDoctorFn: func(contextParameter context.Context) ([]analytics.DoctorConsultation, error) {
			return []analytics.DoctorConsultation{
				{DoctorName: "Dr. Silva", Specialty: "Cardiology", Count: 5},
				{DoctorName: "Dr. Souza", Specialty: "Neurology", Count: 3},
			}, nil
		},
		GetOccupancyRateDataFn: func(contextParameter context.Context) (*analytics.OccupancyRate, error) {
			return &analytics.OccupancyRate{Rate: 75.0, TotalBeds: 20, OccupiedBeds: 15}, nil
		},
		GetAvgWaitTimeDataFn: func(contextParameter context.Context) (*analytics.AvgWaitTime, error) {
			return &analytics.AvgWaitTime{
				AverageMinutes: 12.5,
				ByDepartment: []analytics.DepartmentWaitTime{
					{Department: "Emergency", Minutes: 15.0},
					{Department: "Clinic", Minutes: 10.0},
				},
			}, nil
		},
		GetTopDiagnosesDataFn: func(contextParameter context.Context) ([]analytics.DiagnosisCount, error) {
			return []analytics.DiagnosisCount{
				{ICD10Code: "I10", Description: "Hypertension", Count: 10},
			}, nil
		},
	})

	dashboardData, errorInstance := analyticsService.GetDashboardData(context.Background())

	assert.NoError(testingT, errorInstance)
	assert.Equal(testingT, 8, dashboardData.ConsultationsToday)
	assert.Equal(testingT, "+8%", dashboardData.ConsultationsTrend)
	assert.Equal(testingT, 75.0, dashboardData.OccupancyRate)
	assert.Equal(testingT, 20, dashboardData.OccupancyTotalBeds)
	assert.Equal(testingT, 15, dashboardData.OccupancyOccupiedBeds)
	assert.Equal(testingT, 12.5, dashboardData.AvgWaitTimeMinutes)
	assert.Equal(testingT, 1, dashboardData.NewDiagnosesToday)
	assert.Len(testingT, dashboardData.ConsultationsPerDoctor, 2)
	assert.Len(testingT, dashboardData.WaitTimeByDepartment, 2)
	assert.Len(testingT, dashboardData.TopDiagnoses, 1)
}

func TestGetDashboardData_ConsultationsPerDoctorError_ReturnsError(testingT *testing.T) {
	analyticsService := analytics.NewService(&mocks.MockStatsRepository{
		GetConsultationsPerDoctorFn: func(contextParameter context.Context) ([]analytics.DoctorConsultation, error) {
			return nil, errorRepositoryFailure
		},
	})

	_, errorInstance := analyticsService.GetDashboardData(context.Background())

	assert.Error(testingT, errorInstance)
}

func TestGetDashboardData_OccupancyRateError_ReturnsError(testingT *testing.T) {
	analyticsService := analytics.NewService(&mocks.MockStatsRepository{
		GetConsultationsPerDoctorFn: func(contextParameter context.Context) ([]analytics.DoctorConsultation, error) {
			return []analytics.DoctorConsultation{}, nil
		},
		GetOccupancyRateDataFn: func(contextParameter context.Context) (*analytics.OccupancyRate, error) {
			return nil, errorRepositoryFailure
		},
	})

	_, errorInstance := analyticsService.GetDashboardData(context.Background())

	assert.Error(testingT, errorInstance)
}

func TestGetDashboardData_AvgWaitTimeError_ReturnsError(testingT *testing.T) {
	analyticsService := analytics.NewService(&mocks.MockStatsRepository{
		GetConsultationsPerDoctorFn: func(contextParameter context.Context) ([]analytics.DoctorConsultation, error) {
			return []analytics.DoctorConsultation{}, nil
		},
		GetOccupancyRateDataFn: func(contextParameter context.Context) (*analytics.OccupancyRate, error) {
			return &analytics.OccupancyRate{}, nil
		},
		GetAvgWaitTimeDataFn: func(contextParameter context.Context) (*analytics.AvgWaitTime, error) {
			return nil, errorRepositoryFailure
		},
	})

	_, errorInstance := analyticsService.GetDashboardData(context.Background())

	assert.Error(testingT, errorInstance)
}

func TestGetDashboardData_TopDiagnosesError_ReturnsError(testingT *testing.T) {
	analyticsService := analytics.NewService(&mocks.MockStatsRepository{
		GetConsultationsPerDoctorFn: func(contextParameter context.Context) ([]analytics.DoctorConsultation, error) {
			return []analytics.DoctorConsultation{}, nil
		},
		GetOccupancyRateDataFn: func(contextParameter context.Context) (*analytics.OccupancyRate, error) {
			return &analytics.OccupancyRate{}, nil
		},
		GetAvgWaitTimeDataFn: func(contextParameter context.Context) (*analytics.AvgWaitTime, error) {
			return &analytics.AvgWaitTime{}, nil
		},
		GetTopDiagnosesDataFn: func(contextParameter context.Context) ([]analytics.DiagnosisCount, error) {
			return nil, errorRepositoryFailure
		},
	})

	_, errorInstance := analyticsService.GetDashboardData(context.Background())

	assert.Error(testingT, errorInstance)
}

func TestGetConsultationsPerDoctor_Success(testingT *testing.T) {
	analyticsService := analytics.NewService(&mocks.MockStatsRepository{
		GetConsultationsPerDoctorFn: func(contextParameter context.Context) ([]analytics.DoctorConsultation, error) {
			return []analytics.DoctorConsultation{
				{DoctorName: "Dr. Silva", Specialty: "Cardiology", Count: 5},
			}, nil
		},
	})

	consultations, errorInstance := analyticsService.GetConsultationsPerDoctor(context.Background())

	assert.NoError(testingT, errorInstance)
	assert.Len(testingT, consultations, 1)
	assert.Equal(testingT, "Dr. Silva", consultations[0].DoctorName)
}

func TestGetConsultationsPerDoctor_Error_ReturnsError(testingT *testing.T) {
	analyticsService := analytics.NewService(&mocks.MockStatsRepository{
		GetConsultationsPerDoctorFn: func(contextParameter context.Context) ([]analytics.DoctorConsultation, error) {
			return nil, errorRepositoryFailure
		},
	})

	_, errorInstance := analyticsService.GetConsultationsPerDoctor(context.Background())

	assert.Error(testingT, errorInstance)
}

func TestGetOccupancyRate_Success(testingT *testing.T) {
	analyticsService := analytics.NewService(&mocks.MockStatsRepository{
		GetOccupancyRateDataFn: func(contextParameter context.Context) (*analytics.OccupancyRate, error) {
			return &analytics.OccupancyRate{Rate: 80.0, TotalBeds: 10, OccupiedBeds: 8}, nil
		},
	})

	occupancyRate, errorInstance := analyticsService.GetOccupancyRate(context.Background())

	assert.NoError(testingT, errorInstance)
	assert.Equal(testingT, 80.0, occupancyRate.Rate)
	assert.Equal(testingT, 10, occupancyRate.TotalBeds)
	assert.Equal(testingT, 8, occupancyRate.OccupiedBeds)
}

func TestGetAvgWaitTime_Success(testingT *testing.T) {
	analyticsService := analytics.NewService(&mocks.MockStatsRepository{
		GetAvgWaitTimeDataFn: func(contextParameter context.Context) (*analytics.AvgWaitTime, error) {
			return &analytics.AvgWaitTime{
				AverageMinutes: 15.0,
				ByDepartment: []analytics.DepartmentWaitTime{
					{Department: "Emergency", Minutes: 20.0},
				},
			}, nil
		},
	})

	avgWaitTime, errorInstance := analyticsService.GetAvgWaitTime(context.Background())

	assert.NoError(testingT, errorInstance)
	assert.Equal(testingT, 15.0, avgWaitTime.AverageMinutes)
	assert.Len(testingT, avgWaitTime.ByDepartment, 1)
}

func TestGetTopDiagnoses_Success(testingT *testing.T) {
	analyticsService := analytics.NewService(&mocks.MockStatsRepository{
		GetTopDiagnosesDataFn: func(contextParameter context.Context) ([]analytics.DiagnosisCount, error) {
			return []analytics.DiagnosisCount{
				{ICD10Code: "I10", Description: "Hypertension", Count: 10},
				{ICD10Code: "E11.9", Description: "Diabetes", Count: 5},
			}, nil
		},
	})

	topDiagnoses, errorInstance := analyticsService.GetTopDiagnoses(context.Background())

	assert.NoError(testingT, errorInstance)
	assert.Len(testingT, topDiagnoses, 2)
	assert.Equal(testingT, "I10", topDiagnoses[0].ICD10Code)
}
