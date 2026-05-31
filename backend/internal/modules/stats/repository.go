package stats

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

func (statsRepository *repository) GetTotalPatientsCount(contextParameter context.Context) (int, error) {
	responseBody, errorInstance := statsRepository.fhirClient.SearchResources(contextParameter, "Patient", "_count=100")
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

func (statsRepository *repository) GetEncounters(contextParameter context.Context) ([]FHIREncounter, error) {
	responseBody, errorInstance := statsRepository.fhirClient.SearchResources(contextParameter, "Encounter", "_count=100")
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

func (statsRepository *repository) GetConditions(contextParameter context.Context) ([]FHIRCondition, error) {
	responseBody, errorInstance := statsRepository.fhirClient.SearchResources(contextParameter, "Condition", "_count=100")
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

func (statsRepository *repository) GetExamModalitiesCounts(contextParameter context.Context) (map[string]int, error) {
	queryStatement := `
		SELECT modality, COUNT(*) 
		FROM imaging_studies 
		GROUP BY modality
	`
	rowsResult, queryError := statsRepository.dbPool.Query(contextParameter, queryStatement)
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
