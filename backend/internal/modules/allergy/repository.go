package allergy

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
	CreateAllergyIntolerance(ctx context.Context, allergy *Allergy) (*Allergy, error)
	GetAllergyIntolerancesByPatient(ctx context.Context, patientFHIRID string) ([]*Allergy, error)
	UpdateAllergyIntolerance(ctx context.Context, fhirResourceID string, allergy *Allergy) (*Allergy, error)
	DeleteAllergyIntolerance(ctx context.Context, fhirResourceID string) error
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

func (allergyRepository *repository) UpdateAllergyIntolerance(ctx context.Context, fhirResourceID string, allergy *Allergy) (*Allergy, error) {
	fhirAllergy := fhir.NewAllergyIntoleranceResource(
		allergy.PatientFHIRID,
		allergy.AllergenCode,
		allergy.AllergenDisplay,
		allergy.ClinicalStatus,
		allergy.Reaction,
	)

	responseBody, updateErr := allergyRepository.fhirClient.UpdateResource(ctx, "AllergyIntolerance", fhirResourceID, fhirAllergy)
	if updateErr != nil {
		return nil, fmt.Errorf("failed to update allergy intolerance: %w", updateErr)
	}

	var updatedResource map[string]interface{}
	if parseErr := json.Unmarshal(responseBody, &updatedResource); parseErr != nil {
		return nil, fmt.Errorf("failed to parse allergy update response: %w", parseErr)
	}

	fhirID, _ := updatedResource["id"].(string)
	allergy.FHIRResourceID = fhirID
	return allergy, nil
}

func (allergyRepository *repository) DeleteAllergyIntolerance(ctx context.Context, fhirResourceID string) error {
	return allergyRepository.fhirClient.DeleteResource(ctx, "AllergyIntolerance/"+fhirResourceID)
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
		if clinicalStatus, ok := resource["clinicalStatus"].(map[string]interface{}); ok {
			if coding, ok := clinicalStatus["coding"].([]interface{}); ok && len(coding) > 0 {
				if firstCoding, ok := coding[0].(map[string]interface{}); ok {
					allergy.ClinicalStatus, _ = firstCoding["code"].(string)
				}
			}
		}
		if codes, ok := resource["code"].(map[string]interface{}); ok {
			allergy.AllergenDisplay, _ = codes["text"].(string)
			if coding, ok := codes["coding"].([]interface{}); ok && len(coding) > 0 {
				if firstCoding, ok := coding[0].(map[string]interface{}); ok {
					if code, ok := firstCoding["code"].(string); ok {
						allergy.AllergenCode = code
					}
					if display, ok := firstCoding["display"].(string); ok && allergy.AllergenDisplay == "" {
						allergy.AllergenDisplay = display
					}
				}
			}
		}
		if patient, ok := resource["patient"].(map[string]interface{}); ok {
			if ref, ok := patient["reference"].(string); ok {
				parts := strings.SplitN(ref, "/", 2)
				if len(parts) == 2 {
					allergy.PatientFHIRID = parts[1]
				}
			}
		}
		if recordedStr, ok := resource["recordedDate"].(string); ok {
			if parsedTime, parseErr := time.Parse(time.RFC3339, recordedStr); parseErr == nil {
				allergy.RecordedAt = parsedTime
			}
		}
		if reactions, ok := resource["reaction"].([]interface{}); ok && len(reactions) > 0 {
			if firstReaction, ok := reactions[0].(map[string]interface{}); ok {
				if manifestations, ok := firstReaction["manifestation"].([]interface{}); ok && len(manifestations) > 0 {
					if firstManifestation, ok := manifestations[0].(map[string]interface{}); ok {
						allergy.Reaction, _ = firstManifestation["text"].(string)
					}
				}
			}
		}
		allergies = append(allergies, allergy)
	}
	return allergies, nil
}
