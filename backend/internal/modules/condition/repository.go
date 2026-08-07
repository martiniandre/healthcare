package condition

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/healthcare/backend/internal/shared/apperrors"
	"github.com/healthcare/backend/internal/shared/fhir"
	"github.com/healthcare/backend/internal/shared/healthcare"
)

type Repository interface {
	CreateCondition(ctx context.Context, condition *Condition) (*Condition, error)
	GetConditionByID(ctx context.Context, fhirResourceID string) (*Condition, error)
	GetConditionsByPatient(ctx context.Context, patientFHIRID string) ([]*Condition, error)
	UpdateCondition(ctx context.Context, fhirResourceID string, condition *Condition) (*Condition, error)
	DeleteCondition(ctx context.Context, fhirResourceID string) error
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

func (conditionRepository *repository) UpdateCondition(ctx context.Context, fhirResourceID string, condition *Condition) (*Condition, error) {
	fhirCondition := fhir.NewConditionResource(
		condition.PatientFHIRID,
		condition.EncounterFHIRID,
		condition.ICD10Code,
		condition.CodeDisplay,
		condition.ClinicalStatus,
		condition.OnsetAt,
	)

	responseBody, err := conditionRepository.fhirClient.UpdateResource(ctx, "Condition", fhirResourceID, fhirCondition)
	if err != nil {
		return nil, fmt.Errorf("failed to update condition: %w", err)
	}

	var updatedResource map[string]interface{}
	if err := json.Unmarshal(responseBody, &updatedResource); err != nil {
		return nil, fmt.Errorf("failed to parse condition response: %w", err)
	}

	fhirID, _ := updatedResource["id"].(string)
	condition.FHIRResourceID = fhirID
	return condition, nil
}

func (conditionRepository *repository) GetConditionByID(ctx context.Context, fhirResourceID string) (*Condition, error) {
	responseBody, err := conditionRepository.fhirClient.GetResource(ctx, "Condition", fhirResourceID)
	if err != nil {
		if healthcare.IsNotFound(err) {
			return nil, fmt.Errorf("failed to get condition: %w", apperrors.ErrConditionNotFound)
		}
		return nil, fmt.Errorf("failed to get condition: %w", err)
	}

	decodedResource, err := fhir.DecodeResource[fhir.Condition](responseBody)
	if err != nil {
		return nil, fmt.Errorf("failed to parse condition response: %w", err)
	}
	return mapFHIRConditionToDomain(decodedResource), nil
}

func (conditionRepository *repository) DeleteCondition(ctx context.Context, fhirResourceID string) error {
	err := conditionRepository.fhirClient.DeleteResource(ctx, "Condition/"+fhirResourceID)
	if err != nil {
		if healthcare.IsNotFound(err) {
			return fmt.Errorf("failed to delete condition: %w", apperrors.ErrConditionNotFound)
		}
		return fmt.Errorf("failed to delete condition: %w", err)
	}
	return nil
}

func (conditionRepository *repository) GetConditionsByPatient(ctx context.Context, patientFHIRID string) ([]*Condition, error) {
	queryParams := url.Values{"subject": []string{fmt.Sprintf("Patient/%s", patientFHIRID)}}.Encode()
	responseBody, err := conditionRepository.fhirClient.SearchResources(ctx, "Condition", queryParams)
	if err != nil {
		return nil, fmt.Errorf("failed to search conditions: %w", err)
	}
	return parseConditionBundle(responseBody)
}

func parseConditionBundle(responseBody json.RawMessage) ([]*Condition, error) {
	decodedResources, err := fhir.DecodeBundle[fhir.Condition](responseBody)
	if err != nil {
		return nil, err
	}
	conditions := make([]*Condition, 0, len(decodedResources))
	for _, resource := range decodedResources {
		conditions = append(conditions, mapFHIRConditionToDomain(&resource))
	}
	return conditions, nil
}

func mapFHIRConditionToDomain(resource *fhir.Condition) *Condition {
	condition := &Condition{}
	condition.FHIRResourceID = resource.ID
	clinicalStatusCode, _, _ := fhir.CodeableConceptParts(resource.ClinicalStatus)
	condition.ClinicalStatus = clinicalStatusCode
	code, display, text := fhir.CodeableConceptParts(resource.Code)
	condition.CodeDisplay = text
	condition.ICD10Code = code
	if condition.CodeDisplay == "" {
		condition.CodeDisplay = display
	}
	condition.PatientFHIRID = fhir.SplitReferenceID(resource.Subject.Reference)
	condition.EncounterFHIRID = fhir.SplitReferenceID(resource.Encounter.Reference)
	if parsedTime, ok := fhir.ParseRFC3339(resource.OnsetDateTime); ok {
		condition.OnsetAt = parsedTime
	}
	if condition.OnsetAt.IsZero() {
		if parsedTime, ok := fhir.ParseRFC3339(resource.RecordedDate); ok {
			condition.OnsetAt = parsedTime
		}
	}
	return condition
}
