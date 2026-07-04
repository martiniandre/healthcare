package encounter

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/healthcare/backend/internal/shared/fhir"
	"github.com/healthcare/backend/internal/shared/healthcare"
)

type Repository interface {
	CreateEncounter(ctx context.Context, encounter *Encounter) (*Encounter, error)
	GetEncounterByID(ctx context.Context, fhirResourceID string) (*Encounter, error)
	GetEncountersByPatient(ctx context.Context, patientFHIRID string) ([]*Encounter, error)
}

type repository struct {
	fhirClient healthcare.FHIRClient
}

func NewRepository(fhirClient healthcare.FHIRClient) Repository {
	return &repository{fhirClient: fhirClient}
}

func (encounterRepository *repository) CreateEncounter(ctx context.Context, encounter *Encounter) (*Encounter, error) {
	fhirEncounter := fhir.NewEncounterResource(encounter.PatientFHIRID, encounter.PractitionerID, encounter.ReasonCode, encounter.ReasonDisplay)

	responseBody, err := encounterRepository.fhirClient.CreateResource(ctx, "Encounter", fhirEncounter)
	if err != nil {
		return nil, fmt.Errorf("failed to create encounter: %w", err)
	}

	var createdResource map[string]interface{}
	if err := json.Unmarshal(responseBody, &createdResource); err != nil {
		return nil, fmt.Errorf("failed to parse encounter response: %w", err)
	}

	fhirID, _ := createdResource["id"].(string)
	encounter.FHIRResourceID = fhirID
	return encounter, nil
}

func (encounterRepository *repository) GetEncounterByID(ctx context.Context, fhirResourceID string) (*Encounter, error) {
	responseBody, err := encounterRepository.fhirClient.GetResource(ctx, "Encounter", fhirResourceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get encounter: %w", err)
	}

	var resource map[string]interface{}
	if err := json.Unmarshal(responseBody, &resource); err != nil {
		return nil, fmt.Errorf("failed to parse encounter response: %w", err)
	}

	encounter := &Encounter{}
	encounter.FHIRResourceID, _ = resource["id"].(string)
	encounter.Status, _ = resource["status"].(string)
	if subject, ok := resource["subject"].(map[string]interface{}); ok {
		encounter.PatientFHIRID, _ = subject["reference"].(string)
	}

	return encounter, nil
}

func (encounterRepository *repository) GetEncountersByPatient(ctx context.Context, patientFHIRID string) ([]*Encounter, error) {
	queryParams := url.Values{"subject": []string{fmt.Sprintf("Patient/%s", patientFHIRID)}}.Encode()
	responseBody, err := encounterRepository.fhirClient.SearchResources(ctx, "Encounter", queryParams)
	if err != nil {
		return nil, fmt.Errorf("failed to search encounters: %w", err)
	}

	return parseEncounterBundle(responseBody)
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

func parseEncounterBundle(responseBody json.RawMessage) ([]*Encounter, error) {
	entries, err := extractBundleEntries(responseBody)
	if err != nil {
		return nil, err
	}
	encounters := make([]*Encounter, 0, len(entries))
	for _, resource := range entries {
		encounter := &Encounter{}
		encounter.FHIRResourceID, _ = resource["id"].(string)
		encounter.Status, _ = resource["status"].(string)
		if subject, ok := resource["subject"].(map[string]interface{}); ok {
			encounter.PatientFHIRID, _ = subject["reference"].(string)
		}
		encounters = append(encounters, encounter)
	}
	return encounters, nil
}
