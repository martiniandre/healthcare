package observation

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/healthcare/backend/internal/shared/fhir"
	"github.com/healthcare/backend/internal/shared/healthcare"
)

type Repository interface {
	CreateObservation(ctx context.Context, observation *Observation) (*Observation, error)
	GetObservationsByEncounter(ctx context.Context, encounterFHIRID string) ([]*Observation, error)
	GetObservationsByPatient(ctx context.Context, patientFHIRID string) ([]*Observation, error)
}

type repository struct {
	fhirClient healthcare.FHIRClient
}

func NewRepository(fhirClient healthcare.FHIRClient) Repository {
	return &repository{fhirClient: fhirClient}
}

func (observationRepository *repository) CreateObservation(ctx context.Context, observation *Observation) (*Observation, error) {
	fhirObservation := fhir.NewObservationResource(
		observation.PatientFHIRID,
		observation.EncounterFHIRID,
		observation.LoincCode,
		observation.CodeDisplay,
		observation.ValueQuantity,
		observation.ValueUnit,
	)

	responseBody, err := observationRepository.fhirClient.CreateResource(ctx, "Observation", fhirObservation)
	if err != nil {
		return nil, fmt.Errorf("failed to create observation: %w", err)
	}

	var createdResource map[string]interface{}
	if err := json.Unmarshal(responseBody, &createdResource); err != nil {
		return nil, fmt.Errorf("failed to parse observation response: %w", err)
	}

	fhirID, _ := createdResource["id"].(string)
	observation.FHIRResourceID = fhirID
	return observation, nil
}

func (observationRepository *repository) GetObservationsByEncounter(ctx context.Context, encounterFHIRID string) ([]*Observation, error) {
	queryParams := url.Values{"encounter": []string{fmt.Sprintf("Encounter/%s", encounterFHIRID)}}.Encode()
	responseBody, err := observationRepository.fhirClient.SearchResources(ctx, "Observation", queryParams)
	if err != nil {
		return nil, fmt.Errorf("failed to search observations: %w", err)
	}
	return parseObservationBundle(responseBody)
}

func (observationRepository *repository) GetObservationsByPatient(ctx context.Context, patientFHIRID string) ([]*Observation, error) {
	queryParams := url.Values{"subject": []string{fmt.Sprintf("Patient/%s", patientFHIRID)}}.Encode()
	responseBody, err := observationRepository.fhirClient.SearchResources(ctx, "Observation", queryParams)
	if err != nil {
		return nil, fmt.Errorf("failed to search observations: %w", err)
	}
	return parseObservationBundle(responseBody)
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

func parseObservationBundle(responseBody json.RawMessage) ([]*Observation, error) {
	entries, err := extractBundleEntries(responseBody)
	if err != nil {
		return nil, err
	}
	observations := make([]*Observation, 0, len(entries))
	for _, resource := range entries {
		observation := &Observation{}
		observation.FHIRResourceID, _ = resource["id"].(string)
		if codes, ok := resource["code"].(map[string]interface{}); ok {
			observation.CodeDisplay, _ = codes["text"].(string)
		}
		if valueQuantity, ok := resource["valueQuantity"].(map[string]interface{}); ok {
			observation.ValueQuantity, _ = valueQuantity["value"].(float64)
			observation.ValueUnit, _ = valueQuantity["unit"].(string)
		}
		observations = append(observations, observation)
	}
	return observations, nil
}
