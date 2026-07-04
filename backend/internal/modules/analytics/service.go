package analytics

import (
	"context"
	"fmt"
	"time"

	"github.com/healthcare/backend/internal/shared/apperrors"
)

type Service interface {
	GetStats(contextParameter context.Context) (Stats, error)
	GetDashboardData(contextParameter context.Context) (*DashboardData, error)
	GetConsultationsPerDoctor(contextParameter context.Context) ([]DoctorConsultation, error)
	GetOccupancyRate(contextParameter context.Context) (*OccupancyRate, error)
	GetAvgWaitTime(contextParameter context.Context) (*AvgWaitTime, error)
	GetTopDiagnoses(contextParameter context.Context) ([]DiagnosisCount, error)
}

type service struct {
	analyticsRepository Repository
}

func NewService(analyticsRepository Repository) Service {
	return &service{
		analyticsRepository: analyticsRepository,
	}
}

func (analyticsService *service) GetStats(contextParameter context.Context) (Stats, error) {
	totalPatientsCount, errorInstance := analyticsService.analyticsRepository.GetTotalPatientsCount(contextParameter)
	if errorInstance != nil {
		return Stats{}, fmt.Errorf("failed to get total patients count: %w", apperrors.ErrInternalServer)
	}

	encounters, errorInstance := analyticsService.analyticsRepository.GetEncounters(contextParameter)
	if errorInstance != nil {
		return Stats{}, fmt.Errorf("failed to get encounters: %w", apperrors.ErrInternalServer)
	}

	var totalDuration float64
	var countWithDuration int
	for _, encounterElement := range encounters {
		if encounterElement.StartedAt != "" && encounterElement.EndedAt != "" {
			startTime, errStart := time.Parse(time.RFC3339, encounterElement.StartedAt)
			endTime, errEnd := time.Parse(time.RFC3339, encounterElement.EndedAt)
			if errStart == nil && errEnd == nil {
				durationMinutes := endTime.Sub(startTime).Minutes()
				if durationMinutes > 0 {
					totalDuration += durationMinutes
					countWithDuration++
				}
			}
		}
	}

	var averageServiceDuration float64
	if countWithDuration > 0 {
		averageServiceDuration = totalDuration / float64(countWithDuration)
	}

	weeklyCounts := map[string]int{
		"Mon": 0, "Tue": 0, "Wed": 0, "Thu": 0, "Fri": 0, "Sat": 0, "Sun": 0,
	}
	var totalWeeklyCount int
	for _, encounterElement := range encounters {
		if encounterElement.StartedAt != "" {
			startTime, errStart := time.Parse(time.RFC3339, encounterElement.StartedAt)
			if errStart == nil {
				weekday := startTime.Weekday().String()
				shortWeekday := weekday[:3]
				weeklyCounts[shortWeekday]++
				totalWeeklyCount++
			}
		}
	}

	weeklyConsultationsList := make([]WeeklyConsultation, 0)
	for _, dayValue := range []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"} {
		weeklyConsultationsList = append(weeklyConsultationsList, WeeklyConsultation{
			DayName: dayValue,
			Count:   weeklyCounts[dayValue],
		})
	}

	modalityCounts, errorInstance := analyticsService.analyticsRepository.GetExamModalitiesCounts(contextParameter)
	if errorInstance != nil {
		return Stats{}, fmt.Errorf("failed to get modality counts: %w", apperrors.ErrInternalServer)
	}

	var totalExamsCount int
	for _, examCountValue := range modalityCounts {
		totalExamsCount += examCountValue
	}

	displayNameMap := map[string]string{
		"CT": "CT (Tomografia)",
		"MR": "MR (Ressonância)",
		"CR": "CR (Raio-X)",
		"US": "US (Ultrassom)",
	}
	colorMap := map[string]string{
		"CT": "#2563eb",
		"MR": "#0d9488",
		"CR": "#8b5cf6",
		"US": "#f59e0b",
	}

	modalitiesList := []string{"CT", "MR", "CR", "US"}

	examModalitiesList := make([]ExamModality, 0, len(modalitiesList))
	for _, modalityKey := range modalitiesList {
		countVal := modalityCounts[modalityKey]
		percentageVal := 0.0
		if totalExamsCount > 0 {
			percentageVal = (float64(countVal) / float64(totalExamsCount)) * 100.0
		}
		examModalitiesList = append(examModalitiesList, ExamModality{
			Modality:   displayNameMap[modalityKey],
			Percentage: percentageVal,
			Count:      countVal,
			Color:      colorMap[modalityKey],
		})
	}

	conditions, errorInstance := analyticsService.analyticsRepository.GetConditions(contextParameter)
	if errorInstance != nil {
		return Stats{}, fmt.Errorf("failed to get conditions: %w", apperrors.ErrInternalServer)
	}

	asthmaCount := 0
	hypertensionCount := 0
	diabetesCount := 0

	for _, conditionElement := range conditions {
		switch conditionElement.ICD10Code {
		case "J45.9":
			asthmaCount++
		case "I10":
			hypertensionCount++
		case "E11.9":
			diabetesCount++
		}
	}

	pathologyCasesList := []PathologyCase{}

	if asthmaCount > 0 {
		pathologyCasesList = append(pathologyCasesList, PathologyCase{
			Code:        "J45.9",
			Description: "Unspecified asthma",
			Category:    "Respiratory",
			ActiveCases: asthmaCount,
			Trend:       "+5%",
		})
	}

	if hypertensionCount > 0 {
		pathologyCasesList = append(pathologyCasesList, PathologyCase{
			Code:        "I10",
			Description: "Primary essential hypertension",
			Category:    "Cardiovascular",
			ActiveCases: hypertensionCount,
			Trend:       "Stable",
		})
	}

	if diabetesCount > 0 {
		pathologyCasesList = append(pathologyCasesList, PathologyCase{
			Code:        "E11.9",
			Description: "Type 2 diabetes mellitus",
			Category:    "Endocrine",
			ActiveCases: diabetesCount,
			Trend:       "+12%",
		})
	}

	analyticsData := Stats{
		TotalPatients:             totalPatientsCount,
		FHIRComplianceRate:        0,
		AvgServiceDurationMinutes: averageServiceDuration,
		WeeklyConsultations:       weeklyConsultationsList,
		ExamModalities:            examModalitiesList,
		PathologyCases:            pathologyCasesList,
	}

	return analyticsData, nil
}

func (analyticsService *service) GetDashboardData(contextParameter context.Context) (*DashboardData, error) {
	consultationsPerDoctor, errorInstance := analyticsService.analyticsRepository.GetConsultationsPerDoctor(contextParameter)
	if errorInstance != nil {
		return nil, fmt.Errorf("failed to get consultations per doctor: %w", apperrors.ErrInternalServer)
	}

	occupancyRate, errorInstance := analyticsService.analyticsRepository.GetOccupancyRateData(contextParameter)
	if errorInstance != nil {
		return nil, fmt.Errorf("failed to get occupancy rate: %w", apperrors.ErrInternalServer)
	}

	avgWaitTime, errorInstance := analyticsService.analyticsRepository.GetAvgWaitTimeData(contextParameter)
	if errorInstance != nil {
		return nil, fmt.Errorf("failed to get average wait time: %w", apperrors.ErrInternalServer)
	}

	topDiagnoses, errorInstance := analyticsService.analyticsRepository.GetTopDiagnosesData(contextParameter)
	if errorInstance != nil {
		return nil, fmt.Errorf("failed to get top diagnoses: %w", apperrors.ErrInternalServer)
	}

	var totalConsultations int
	for _, doctorConsultation := range consultationsPerDoctor {
		totalConsultations += doctorConsultation.Count
	}

	var totalWaitMinutes float64
	var departmentCount int
	for _, departmentWaitTime := range avgWaitTime.ByDepartment {
		totalWaitMinutes += departmentWaitTime.Minutes
		departmentCount++
	}

	var averageWaitTime float64
	if departmentCount > 0 {
		averageWaitTime = totalWaitMinutes / float64(departmentCount)
	}

	dashboardData := &DashboardData{
		ConsultationsToday:      totalConsultations,
		ConsultationsTrend:      "+8%",
		OccupancyRate:           occupancyRate.Rate,
		OccupancyTotalBeds:      occupancyRate.TotalBeds,
		OccupancyOccupiedBeds:   occupancyRate.OccupiedBeds,
		AvgWaitTimeMinutes:      averageWaitTime,
		ActivePatients:          0,
		ExamsToday:              0,
		NewDiagnosesToday:       len(topDiagnoses),
		ConsultationsPerDoctor:  consultationsPerDoctor,
		WaitTimeByDepartment:    avgWaitTime.ByDepartment,
		TopDiagnoses:            topDiagnoses,
	}

	return dashboardData, nil
}

func (analyticsService *service) GetConsultationsPerDoctor(contextParameter context.Context) ([]DoctorConsultation, error) {
	consultationsPerDoctor, errorInstance := analyticsService.analyticsRepository.GetConsultationsPerDoctor(contextParameter)
	if errorInstance != nil {
		return nil, fmt.Errorf("failed to get consultations per doctor: %w", apperrors.ErrInternalServer)
	}
	return consultationsPerDoctor, nil
}

func (analyticsService *service) GetOccupancyRate(contextParameter context.Context) (*OccupancyRate, error) {
	occupancyRate, errorInstance := analyticsService.analyticsRepository.GetOccupancyRateData(contextParameter)
	if errorInstance != nil {
		return nil, fmt.Errorf("failed to get occupancy rate: %w", apperrors.ErrInternalServer)
	}
	return occupancyRate, nil
}

func (analyticsService *service) GetAvgWaitTime(contextParameter context.Context) (*AvgWaitTime, error) {
	avgWaitTime, errorInstance := analyticsService.analyticsRepository.GetAvgWaitTimeData(contextParameter)
	if errorInstance != nil {
		return nil, fmt.Errorf("failed to get average wait time: %w", apperrors.ErrInternalServer)
	}
	return avgWaitTime, nil
}

func (analyticsService *service) GetTopDiagnoses(contextParameter context.Context) ([]DiagnosisCount, error) {
	topDiagnoses, errorInstance := analyticsService.analyticsRepository.GetTopDiagnosesData(contextParameter)
	if errorInstance != nil {
		return nil, fmt.Errorf("failed to get top diagnoses: %w", apperrors.ErrInternalServer)
	}
	return topDiagnoses, nil
}
