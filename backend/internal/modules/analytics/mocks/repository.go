package mocks

import (
	"context"

	"github.com/healthcare/backend/internal/modules/analytics"
)

type MockStatsRepository struct {
	GetTotalPatientsCountFn    func(contextParameter context.Context) (int, error)
	GetEncountersFn            func(contextParameter context.Context) ([]analytics.FHIREncounter, error)
	GetConditionsFn            func(contextParameter context.Context) ([]analytics.FHIRCondition, error)
	GetExamModalitiesCountsFn  func(contextParameter context.Context) (map[string]int, error)
	GetConsultationsPerDoctorFn func(contextParameter context.Context) ([]analytics.DoctorConsultation, error)
	GetOccupancyRateDataFn     func(contextParameter context.Context) (*analytics.OccupancyRate, error)
	GetAvgWaitTimeDataFn       func(contextParameter context.Context) (*analytics.AvgWaitTime, error)
	GetTopDiagnosesDataFn      func(contextParameter context.Context) ([]analytics.DiagnosisCount, error)
}

func (mockRepo *MockStatsRepository) GetTotalPatientsCount(contextParameter context.Context) (int, error) {
	if mockRepo.GetTotalPatientsCountFn != nil {
		return mockRepo.GetTotalPatientsCountFn(contextParameter)
	}
	return 100, nil
}

func (mockRepo *MockStatsRepository) GetEncounters(contextParameter context.Context) ([]analytics.FHIREncounter, error) {
	if mockRepo.GetEncountersFn != nil {
		return mockRepo.GetEncountersFn(contextParameter)
	}
	return []analytics.FHIREncounter{}, nil
}

func (mockRepo *MockStatsRepository) GetConditions(contextParameter context.Context) ([]analytics.FHIRCondition, error) {
	if mockRepo.GetConditionsFn != nil {
		return mockRepo.GetConditionsFn(contextParameter)
	}
	return []analytics.FHIRCondition{}, nil
}

func (mockRepo *MockStatsRepository) GetExamModalitiesCounts(contextParameter context.Context) (map[string]int, error) {
	if mockRepo.GetExamModalitiesCountsFn != nil {
		return mockRepo.GetExamModalitiesCountsFn(contextParameter)
	}
	return map[string]int{}, nil
}

func (mockRepo *MockStatsRepository) GetConsultationsPerDoctor(contextParameter context.Context) ([]analytics.DoctorConsultation, error) {
	if mockRepo.GetConsultationsPerDoctorFn != nil {
		return mockRepo.GetConsultationsPerDoctorFn(contextParameter)
	}
	return []analytics.DoctorConsultation{}, nil
}

func (mockRepo *MockStatsRepository) GetOccupancyRateData(contextParameter context.Context) (*analytics.OccupancyRate, error) {
	if mockRepo.GetOccupancyRateDataFn != nil {
		return mockRepo.GetOccupancyRateDataFn(contextParameter)
	}
	return &analytics.OccupancyRate{Rate: 0, TotalBeds: 0, OccupiedBeds: 0}, nil
}

func (mockRepo *MockStatsRepository) GetAvgWaitTimeData(contextParameter context.Context) (*analytics.AvgWaitTime, error) {
	if mockRepo.GetAvgWaitTimeDataFn != nil {
		return mockRepo.GetAvgWaitTimeDataFn(contextParameter)
	}
	return &analytics.AvgWaitTime{AverageMinutes: 0, ByDepartment: []analytics.DepartmentWaitTime{}}, nil
}

func (mockRepo *MockStatsRepository) GetTopDiagnosesData(contextParameter context.Context) ([]analytics.DiagnosisCount, error) {
	if mockRepo.GetTopDiagnosesDataFn != nil {
		return mockRepo.GetTopDiagnosesDataFn(contextParameter)
	}
	return []analytics.DiagnosisCount{}, nil
}
