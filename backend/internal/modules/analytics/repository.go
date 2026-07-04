package analytics

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/healthcare/backend/internal/shared/healthcare"
	"github.com/jackc/pgx/v5/pgxpool"
)

type FHIREncounter struct {
	ID             string
	Status         string
	StartedAt      string
	EndedAt        string
	PractitionerID string
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
		if participantList, participantOk := resourceMap["participant"].([]interface{}); participantOk && len(participantList) > 0 {
			if firstParticipant, firstOk := participantList[0].(map[string]interface{}); firstOk {
				if individualMap, individualOk := firstParticipant["individual"].(map[string]interface{}); individualOk {
					if referenceValue, refOk := individualMap["reference"].(string); refOk {
						referenceParts := strings.Split(referenceValue, "/")
						encounter.PractitionerID = referenceParts[len(referenceParts)-1]
					}
				}
			}
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
	encounters, errorInstance := analyticsRepository.GetEncounters(contextParameter)
	if errorInstance != nil {
		return nil, errorInstance
	}

	employeeRows, queryError := analyticsRepository.dbPool.Query(contextParameter,
		`SELECT id, full_name, COALESCE(crm_number, '') FROM employees WHERE role = 'doctor' AND is_active = true`)
	if queryError != nil {
		return nil, queryError
	}
	defer employeeRows.Close()

	type doctorInfo struct {
		Name     string
		CRM      string
	}
	doctorsByID := make(map[string]doctorInfo)
	for employeeRows.Next() {
		var id string
		var info doctorInfo
		if scanError := employeeRows.Scan(&id, &info.Name, &info.CRM); scanError != nil {
			return nil, scanError
		}
		doctorsByID[id] = info
	}
	if rowsError := employeeRows.Err(); rowsError != nil {
		return nil, rowsError
	}

	consultationCounts := make(map[string]*DoctorConsultation)
	for _, encounterElement := range encounters {
		if encounterElement.PractitionerID == "" {
			continue
		}
		doctorData, doctorExists := doctorsByID[encounterElement.PractitionerID]
		if !doctorExists {
			continue
		}
		entry, entryExists := consultationCounts[encounterElement.PractitionerID]
		if !entryExists {
			specialtyValue := ""
			if doctorData.CRM != "" {
				specialtyValue = doctorData.CRM
			}
			consultationCounts[encounterElement.PractitionerID] = &DoctorConsultation{
				DoctorName: doctorData.Name,
				Specialty:  specialtyValue,
				Count:      0,
			}
			entry = consultationCounts[encounterElement.PractitionerID]
		}
		entry.Count++
	}

	consultationsList := make([]DoctorConsultation, 0, len(consultationCounts))
	for _, consultationData := range consultationCounts {
		consultationsList = append(consultationsList, *consultationData)
	}
	return consultationsList, nil
}

func (analyticsRepository *repository) GetOccupancyRateData(contextParameter context.Context) (*OccupancyRate, error) {
	queryStatement := `
		SELECT
			COUNT(*) as total_beds,
			COUNT(*) FILTER (WHERE status = 'occupied') as occupied_beds
		FROM telemetry_beds
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
	return &AvgWaitTime{
		AverageMinutes: 0,
		ByDepartment:   make([]DepartmentWaitTime, 0),
	}, nil
}

func (analyticsRepository *repository) GetTopDiagnosesData(contextParameter context.Context) ([]DiagnosisCount, error) {
	conditions, errorInstance := analyticsRepository.GetConditions(contextParameter)
	if errorInstance != nil {
		return nil, errorInstance
	}

	diagnosisCounts := make(map[string]*DiagnosisCount)
	codeOrder := make([]string, 0)
	for _, conditionElement := range conditions {
		if conditionElement.ICD10Code == "" {
			continue
		}
		_, entryExists := diagnosisCounts[conditionElement.ICD10Code]
		if !entryExists {
			codeOrder = append(codeOrder, conditionElement.ICD10Code)
			diagnosisCounts[conditionElement.ICD10Code] = &DiagnosisCount{
				ICD10Code:   conditionElement.ICD10Code,
				Description: "",
				Count:       0,
			}
		}
		diagnosisCounts[conditionElement.ICD10Code].Count++
	}

	diagnosesList := make([]DiagnosisCount, 0, len(codeOrder))
	for _, icdCode := range codeOrder {
		diagnosesList = append(diagnosesList, *diagnosisCounts[icdCode])
	}

	const topDiagnosesLimit = 10
	if len(diagnosesList) > topDiagnosesLimit {
		diagnosesList = diagnosesList[:topDiagnosesLimit]
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
