package medication

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
	CreateMedicationRequest(ctx context.Context, medication *Medication) (*Medication, error)
	GetMedicationRequestsByEncounter(ctx context.Context, encounterFHIRID string) ([]*Medication, error)
}

type repository struct {
	fhirClient healthcare.FHIRClient
}

func NewRepository(fhirClient healthcare.FHIRClient) Repository {
	return &repository{fhirClient: fhirClient}
}

func (medicationRepository *repository) CreateMedicationRequest(ctx context.Context, medication *Medication) (*Medication, error) {
	fhirMedication := fhir.NewMedicationRequestResource(
		medication.PatientFHIRID,
		medication.EncounterFHIRID,
		medication.PractitionerFHIRID,
		medication.MedicationCode,
		medication.MedicationName,
		medication.DosageInstructions,
	)

	responseBody, err := medicationRepository.fhirClient.CreateResource(ctx, "MedicationRequest", fhirMedication)
	if err != nil {
		return nil, fmt.Errorf("failed to create medication request: %w", err)
	}

	var createdResource map[string]interface{}
	if err := json.Unmarshal(responseBody, &createdResource); err != nil {
		return nil, fmt.Errorf("failed to parse medication request response: %w", err)
	}

	fhirID, _ := createdResource["id"].(string)
	medication.FHIRResourceID = fhirID
	medication.IssuedAt = time.Now()
	return medication, nil
}

func (medicationRepository *repository) GetMedicationRequestsByEncounter(ctx context.Context, encounterFHIRID string) ([]*Medication, error) {
	queryParams := url.Values{"encounter": []string{fmt.Sprintf("Encounter/%s", encounterFHIRID)}}.Encode()
	responseBody, err := medicationRepository.fhirClient.SearchResources(ctx, "MedicationRequest", queryParams)
	if err != nil {
		return nil, fmt.Errorf("failed to search medication requests: %w", err)
	}
	return parseMedicationRequestBundle(responseBody)
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

func parseMedicationRequestBundle(responseBody json.RawMessage) ([]*Medication, error) {
	entries, err := extractBundleEntries(responseBody)
	if err != nil {
		return nil, err
	}
	medications := make([]*Medication, 0, len(entries))
	for _, resource := range entries {
		medication := &Medication{}
		medication.FHIRResourceID, _ = resource["id"].(string)
		medication.Status, _ = resource["status"].(string)
		if med, ok := resource["medicationCodeableConcept"].(map[string]interface{}); ok {
			medication.MedicationName, _ = med["text"].(string)
			if coding, ok := med["coding"].([]interface{}); ok && len(coding) > 0 {
				if firstCoding, ok := coding[0].(map[string]interface{}); ok {
					medication.MedicationCode, _ = firstCoding["code"].(string)
					if display, ok := firstCoding["display"].(string); ok && medication.MedicationName == "" {
						medication.MedicationName = display
					}
				}
			}
		}
		if encounter, ok := resource["encounter"].(map[string]interface{}); ok {
			if ref, ok := encounter["reference"].(string); ok {
				parts := strings.SplitN(ref, "/", 2)
				if len(parts) == 2 {
					medication.EncounterFHIRID = parts[1]
				}
			}
		}
		if subject, ok := resource["subject"].(map[string]interface{}); ok {
			if ref, ok := subject["reference"].(string); ok {
				parts := strings.SplitN(ref, "/", 2)
				if len(parts) == 2 {
					medication.PatientFHIRID = parts[1]
				}
			}
		}
		if requester, ok := resource["requester"].(map[string]interface{}); ok {
			if agent, ok := requester["agent"].(map[string]interface{}); ok {
				if ref, ok := agent["reference"].(string); ok {
					parts := strings.SplitN(ref, "/", 2)
					if len(parts) == 2 {
						medication.PractitionerFHIRID = parts[1]
					}
				}
			}
		}
		if authoredStr, ok := resource["authoredOn"].(string); ok {
			if parsedTime, parseErr := time.Parse(time.RFC3339, authoredStr); parseErr == nil {
				medication.IssuedAt = parsedTime
			}
		}
		medications = append(medications, medication)
	}
	return medications, nil
}
