package condition

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/healthcare/backend/internal/shared/fhir"
	"github.com/healthcare/backend/internal/shared/healthcare"
)

type Repository interface {
	CreateCondition(ctx context.Context, condition *Condition) (*Condition, error)
	GetConditionsByPatient(ctx context.Context, patientFHIRID string) ([]*Condition, error)
}

type repository struct {
	fhirClient healthcare.FHIRClient
}

func NewRepository(fhirClient healthcare.FHIRClient) Repository {
	return &repository{fhirClient: fhirClient}
}

func (conditionRepository *repository) CreateCondition(ctx context.Context, condition *Condition) (*Condition, error) {
	fhirCondition := fhir.NewConditionResource(
		condition.PatientFHIRID,
		condition.EncounterFHIRID,
		condition.ICD10Code,
		condition.CodeDisplay,
		condition.ClinicalStatus,
		condition.OnsetAt,
	)

	responseBody, err := conditionRepository.fhirClient.CreateResource(ctx, "Condition", fhirCondition)
	if err != nil {
		return nil, fmt.Errorf("failed to create condition: %w", err)
	}

	var createdResource map[string]interface{}
	if err := json.Unmarshal(responseBody, &createdResource); err != nil {
		return nil, fmt.Errorf("failed to parse condition response: %w", err)
	}

	fhirID, _ := createdResource["id"].(string)
	condition.FHIRResourceID = fhirID
	return condition, nil
}

func (conditionRepository *repository) GetConditionsByPatient(ctx context.Context, patientFHIRID string) ([]*Condition, error) {
	queryParams := url.Values{"subject": []string{fmt.Sprintf("Patient/%s", patientFHIRID)}}.Encode()
	responseBody, err := conditionRepository.fhirClient.SearchResources(ctx, "Condition", queryParams)
	if err != nil {
		return nil, fmt.Errorf("failed to search conditions: %w", err)
	}
	return parseConditionBundle(responseBody)
}

func extractBundleEntries(responseBody json.RawMessage) ([]map[string]interface{}, error) {
	var bundle map[string]interface{}
	if err := json.Unmarshal(responseBody, &bundle); err != nil {
		return nil, err
	}
	rawEntries, ok := bundle["entry"].([]interface{})
	if !ok {
		return []map[string]interface{}{}, nil
	}
	entries := make([]map[string]interface{}, 0, len(rawEntries))
	for _, rawEntry := range rawEntries {
		entryMap, ok := rawEntry.(map[string]interface{})
		if !ok {
			continue
		}
		resource, ok := entryMap["resource"].(map[string]interface{})
		if !ok {
			continue
		}
		entries = append(entries, resource)
	}
	return entries, nil
}

func parseConditionBundle(responseBody json.RawMessage) ([]*Condition, error) {
	entries, err := extractBundleEntries(responseBody)
	if err != nil {
		return nil, err
	}
	conditions := make([]*Condition, 0, len(entries))
	for _, resource := range entries {
		condition := &Condition{}
		condition.FHIRResourceID, _ = resource["id"].(string)
		if codes, ok := resource["code"].(map[string]interface{}); ok {
			condition.CodeDisplay, _ = codes["text"].(string)
		}
		conditions = append(conditions, condition)
	}
	return conditions, nil
}
