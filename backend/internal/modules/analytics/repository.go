package analytics

import (
	"context"
	"encoding/json"

	"github.com/healthcare/backend/internal/shared/healthcare"
	"github.com/jackc/pgx/v5/pgxpool"
)

type FHIREncounter struct {
	ID        string
	Status    string
	StartedAt string
	EndedAt   string
}

type FHIRCondition struct {
	ID        string
	ICD10Code string
}

type Repository interface {
	GetTotalPatientsCount(contextParameter context.Context) (int, error)
	GetEncounters(contextParameter context.Context) ([]FHIREncounter, error)
	GetConditions(contextParameter context.Context) ([]FHIRCondition, error)
	GetExamModalitiesCounts(contextParameter context.Context) (map[string]int, error)
	GetConsultationsPerDoctor(contextParameter context.Context) ([]DoctorConsultation, error)
	GetOccupancyRateData(contextParameter context.Context) (*OccupancyRate, error)
	GetAvgWaitTimeData(contextParameter context.Context) (*AvgWaitTime, error)
	GetTopDiagnosesData(contextParameter context.Context) ([]DiagnosisCount, error)
}

type repository struct {
	dbPool     *pgxpool.Pool
	fhirClient healthcare.FHIRClient
}

func NewRepository(dbPool *pgxpool.Pool, fhirClient healthcare.FHIRClient) Repository {
	return &repository{
		dbPool:     dbPool,
		fhirClient: fhirClient,
	}
}

func (analyticsRepository *repository) GetTotalPatientsCount(contextParameter context.Context) (int, error) {
	responseBody, errorInstance := analyticsRepository.fhirClient.SearchResources(contextParameter, "Patient", "_count=100")
	if errorInstance != nil {
		return 0, errorInstance
	}
	var bundle map[string]interface{}
	if unmarshalError := json.Unmarshal(responseBody, &bundle); unmarshalError != nil {
		return 0, unmarshalError
	}
	totalValue, totalExists := bundle["total"].(float64)
	if totalExists {
		return int(totalValue), nil
	}
	entries, entriesExists := bundle["entry"].([]interface{})
	if !entriesExists {
		return 0, nil
	}
	return len(entries), nil
}

func (analyticsRepository *repository) GetEncounters(contextParameter context.Context) ([]FHIREncounter, error) {
	responseBody, errorInstance := analyticsRepository.fhirClient.SearchResources(contextParameter, "Encounter", "_count=100")
	if errorInstance != nil {
		return nil, errorInstance
	}
	var bundle map[string]interface{}
	if unmarshalError := json.Unmarshal(responseBody, &bundle); unmarshalError != nil {
		return nil, unmarshalError
	}
	entries, entriesExists := bundle["entry"].([]interface{})
	if !entriesExists {
		return []FHIREncounter{}, nil
	}
	encounters := make([]FHIREncounter, 0, len(entries))
	for _, entryElement := range entries {
		entryMap, entryOk := entryElement.(map[string]interface{})
		if !entryOk {
			continue
		}
		resourceMap, resourceOk := entryMap["resource"].(map[string]interface{})
		if !resourceOk {
			continue
		}
		var encounter FHIREncounter
		encounter.ID, _ = resourceMap["id"].(string)
		encounter.Status, _ = resourceMap["status"].(string)
		if periodMap, periodOk := resourceMap["period"].(map[string]interface{}); periodOk {
			encounter.StartedAt, _ = periodMap["start"].(string)
			encounter.EndedAt, _ = periodMap["end"].(string)
		}
		encounters = append(encounters, encounter)
	}
	return encounters, nil
}

func (analyticsRepository *repository) GetConditions(contextParameter context.Context) ([]FHIRCondition, error) {
	responseBody, errorInstance := analyticsRepository.fhirClient.SearchResources(contextParameter, "Condition", "_count=100")
	if errorInstance != nil {
		return nil, errorInstance
	}
	var bundle map[string]interface{}
	if unmarshalError := json.Unmarshal(responseBody, &bundle); unmarshalError != nil {
		return nil, unmarshalError
	}
	entries, entriesExists := bundle["entry"].([]interface{})
	if !entriesExists {
		return []FHIRCondition{}, nil
	}
	conditions := make([]FHIRCondition, 0, len(entries))
	for _, entryElement := range entries {
		entryMap, entryOk := entryElement.(map[string]interface{})
		if !entryOk {
			continue
		}
		resourceMap, resourceOk := entryMap["resource"].(map[string]interface{})
		if !resourceOk {
			continue
		}
		var condition FHIRCondition
		condition.ID, _ = resourceMap["id"].(string)
		if codeMap, codeOk := resourceMap["code"].(map[string]interface{}); codeOk {
			if codingList, codingOk := codeMap["coding"].([]interface{}); codingOk && len(codingList) > 0 {
				if firstCodingMap, firstCodingOk := codingList[0].(map[string]interface{}); firstCodingOk {
					condition.ICD10Code, _ = firstCodingMap["code"].(string)
				}
			}
		}
		conditions = append(conditions, condition)
	}
	return conditions, nil
}

func (analyticsRepository *repository) GetConsultationsPerDoctor(contextParameter context.Context) ([]DoctorConsultation, error) {
	queryStatement := `
		SELECT s.name, s.specialty, COUNT(e.id) as consultation_count
		FROM encounters e
		JOIN staff s ON e.practitioner_id = s.id
		WHERE e.created_at >= NOW() - INTERVAL '30 days'
		GROUP BY s.id, s.name, s.specialty
		ORDER BY consultation_count DESC
	`
	rowsResult, queryError := analyticsRepository.dbPool.Query(contextParameter, queryStatement)
	if queryError != nil {
		return nil, queryError
	}
	defer rowsResult.Close()

	consultationsList := make([]DoctorConsultation, 0)
	for rowsResult.Next() {
		var doctorConsultation DoctorConsultation
		if scanError := rowsResult.Scan(&doctorConsultation.DoctorName, &doctorConsultation.Specialty, &doctorConsultation.Count); scanError != nil {
			return nil, scanError
		}
		consultationsList = append(consultationsList, doctorConsultation)
	}
	if rowsError := rowsResult.Err(); rowsError != nil {
		return nil, rowsError
	}
	return consultationsList, nil
}

func (analyticsRepository *repository) GetOccupancyRateData(contextParameter context.Context) (*OccupancyRate, error) {
	queryStatement := `
		SELECT
			COUNT(*) as total_beds,
			COUNT(*) FILTER (WHERE status = 'occupied') as occupied_beds
		FROM beds
	`
	rowResult := analyticsRepository.dbPool.QueryRow(contextParameter, queryStatement)
	var occupancyRate OccupancyRate
	if scanError := rowResult.Scan(&occupancyRate.TotalBeds, &occupancyRate.OccupiedBeds); scanError != nil {
		return nil, scanError
	}
	if occupancyRate.TotalBeds > 0 {
		occupancyRate.Rate = (float64(occupancyRate.OccupiedBeds) / float64(occupancyRate.TotalBeds)) * 100
	}
	return &occupancyRate, nil
}

func (analyticsRepository *repository) GetAvgWaitTimeData(contextParameter context.Context) (*AvgWaitTime, error) {
	queryStatement := `
		SELECT
			AVG(EXTRACT(EPOCH FROM (started_at - arrived_at)) / 60) as avg_minutes
		FROM encounters
		WHERE started_at IS NOT NULL AND arrived_at IS NOT NULL
			AND created_at >= NOW() - INTERVAL '30 days'
	`
	rowResult := analyticsRepository.dbPool.QueryRow(contextParameter, queryStatement)
	var avgWaitTime AvgWaitTime
	avgWaitTime.ByDepartment = make([]DepartmentWaitTime, 0)
	if scanError := rowResult.Scan(&avgWaitTime.AverageMinutes); scanError != nil {
		return nil, scanError
	}
	return &avgWaitTime, nil
}

func (analyticsRepository *repository) GetTopDiagnosesData(contextParameter context.Context) ([]DiagnosisCount, error) {
	queryStatement := `
		SELECT c.code, c.display, COUNT(*) as diagnosis_count
		FROM conditions c
		WHERE c.created_at >= NOW() - INTERVAL '30 days'
		GROUP BY c.code, c.display
		ORDER BY diagnosis_count DESC
		LIMIT 10
	`
	rowsResult, queryError := analyticsRepository.dbPool.Query(contextParameter, queryStatement)
	if queryError != nil {
		return nil, queryError
	}
	defer rowsResult.Close()

	diagnosesList := make([]DiagnosisCount, 0)
	for rowsResult.Next() {
		var diagnosisCount DiagnosisCount
		if scanError := rowsResult.Scan(&diagnosisCount.ICD10Code, &diagnosisCount.Description, &diagnosisCount.Count); scanError != nil {
			return nil, scanError
		}
		diagnosesList = append(diagnosesList, diagnosisCount)
	}
	if rowsError := rowsResult.Err(); rowsError != nil {
		return nil, rowsError
	}
	return diagnosesList, nil
}

func (analyticsRepository *repository) GetExamModalitiesCounts(contextParameter context.Context) (map[string]int, error) {
	queryStatement := `
		SELECT modality, COUNT(*) 
		FROM imaging_studies 
		GROUP BY modality
	`
	rowsResult, queryError := analyticsRepository.dbPool.Query(contextParameter, queryStatement)
	if queryError != nil {
		return nil, queryError
	}
	defer rowsResult.Close()

	countsMap := make(map[string]int)
	for rowsResult.Next() {
		var modality string
		var count int
		if scanError := rowsResult.Scan(&modality, &count); scanError != nil {
			return nil, scanError
		}
		countsMap[modality] = count
	}
	if rowsError := rowsResult.Err(); rowsError != nil {
		return nil, rowsError
	}
	return countsMap, nil
}
