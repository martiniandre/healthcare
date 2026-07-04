package allergy

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/healthcare/backend/internal/shared/fhir"
	"github.com/healthcare/backend/internal/shared/healthcare"
)

type Repository interface {
	CreateAllergyIntolerance(ctx context.Context, allergy *Allergy) (*Allergy, error)
	GetAllergyIntolerancesByPatient(ctx context.Context, patientFHIRID string) ([]*Allergy, error)
}

type repository struct {
	fhirClient healthcare.FHIRClient
}

func NewRepository(fhirClient healthcare.FHIRClient) Repository {
	return &repository{fhirClient: fhirClient}
}

func (allergyRepository *repository) CreateAllergyIntolerance(ctx context.Context, allergy *Allergy) (*Allergy, error) {
	fhirAllergy := fhir.NewAllergyIntoleranceResource(
		allergy.PatientFHIRID,
		allergy.AllergenCode,
		allergy.AllergenDisplay,
		allergy.ClinicalStatus,
		allergy.Reaction,
	)

	responseBody, err := allergyRepository.fhirClient.CreateResource(ctx, "AllergyIntolerance", fhirAllergy)
	if err != nil {
		return nil, fmt.Errorf("failed to create allergy intolerance: %w", err)
	}

	var createdResource map[string]interface{}
	if err := json.Unmarshal(responseBody, &createdResource); err != nil {
		return nil, fmt.Errorf("failed to parse allergy response: %w", err)
	}

	fhirID, _ := createdResource["id"].(string)
	allergy.FHIRResourceID = fhirID
	allergy.RecordedAt = time.Now()
	return allergy, nil
}

func (allergyRepository *repository) GetAllergyIntolerancesByPatient(ctx context.Context, patientFHIRID string) ([]*Allergy, error) {
	queryParams := url.Values{"patient": []string{fmt.Sprintf("Patient/%s", patientFHIRID)}}.Encode()
	responseBody, err := allergyRepository.fhirClient.SearchResources(ctx, "AllergyIntolerance", queryParams)
	if err != nil {
		return nil, fmt.Errorf("failed to search allergy intolerances: %w", err)
	}
	return parseAllergyBundle(responseBody)
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

func parseAllergyBundle(responseBody json.RawMessage) ([]*Allergy, error) {
	entries, err := extractBundleEntries(responseBody)
	if err != nil {
		return nil, err
	}
	allergies := make([]*Allergy, 0, len(entries))
	for _, resource := range entries {
		allergy := &Allergy{}
		allergy.FHIRResourceID, _ = resource["id"].(string)
		if codes, ok := resource["code"].(map[string]interface{}); ok {
			allergy.AllergenDisplay, _ = codes["text"].(string)
		}
		allergies = append(allergies, allergy)
	}
	return allergies, nil
}
