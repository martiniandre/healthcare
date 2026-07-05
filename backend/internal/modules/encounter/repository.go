package encounter

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"

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

	return parseEncounterFromResource(resource), nil
}

func parseEncounterFromResource(resource map[string]interface{}) *Encounter {
	encounter := &Encounter{}
	encounter.FHIRResourceID, _ = resource["id"].(string)
	encounter.Status, _ = resource["status"].(string)
	if subject, ok := resource["subject"].(map[string]interface{}); ok {
		encounter.PatientFHIRID, _ = subject["reference"].(string)
	}
	if period, ok := resource["period"].(map[string]interface{}); ok {
		if start, ok := period["start"].(string); ok {
			if parsedTime, parseErr := time.Parse(time.RFC3339, start); parseErr == nil {
				encounter.StartedAt = parsedTime
			}
		}
	}
	if reasonCodes, ok := resource["reasonCode"].([]interface{}); ok && len(reasonCodes) > 0 {
		if firstReason, ok := reasonCodes[0].(map[string]interface{}); ok {
			if text, ok := firstReason["text"].(string); ok {
				encounter.ReasonCode = text
				encounter.ReasonDisplay = text
			} else if coding, ok := firstReason["coding"].([]interface{}); ok && len(coding) > 0 {
				if firstCoding, ok := coding[0].(map[string]interface{}); ok {
					encounter.ReasonCode, _ = firstCoding["code"].(string)
					encounter.ReasonDisplay, _ = firstCoding["display"].(string)
				}
			}
		}
	}
	if participants, ok := resource["participant"].([]interface{}); ok && len(participants) > 0 {
		if firstParticipant, ok := participants[0].(map[string]interface{}); ok {
			if individual, ok := firstParticipant["individual"].(map[string]interface{}); ok {
				if ref, ok := individual["reference"].(string); ok {
					parts := strings.SplitN(ref, "/", 2)
					if len(parts) == 2 {
						encounter.PractitionerID = parts[1]
					}
				}
			}
		}
	}

	return encounter
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
		encounter := parseEncounterFromResource(resource)
		encounters = append(encounters, encounter)
	}
	return encounters, nil
}
