package portal

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/healthcare/backend/internal/shared/healthcare"
)

type Repository interface {
	GetPatient(ctx context.Context, fhirResourceID string) (*PatientInfo, error)
	GetEncountersByPatient(ctx context.Context, patientFHIRID string) ([]PortalEncounter, error)
	GetObservationsByPatient(ctx context.Context, patientFHIRID string) ([]PortalObservation, error)
	GetConditionsByPatient(ctx context.Context, patientFHIRID string) ([]PortalCondition, error)
	GetMedicationsByPatient(ctx context.Context, patientFHIRID string) ([]PortalMedication, error)
	GetReportsByPatient(ctx context.Context, patientFHIRID string) ([]PortalReport, error)
	GetImagingByPatient(ctx context.Context, patientFHIRID string) ([]PortalImaging, error)
}

type repository struct {
	fhirClient healthcare.FHIRClient
}

func NewRepository(fhirClient healthcare.FHIRClient) Repository {
	return &repository{fhirClient: fhirClient}
}

func (portalRepository *repository) GetPatient(ctx context.Context, fhirResourceID string) (*PatientInfo, error) {
	responseBody, err := portalRepository.fhirClient.GetResource(ctx, "Patient", fhirResourceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get patient from healthcare api: %w", err)
	}

	var resource map[string]interface{}
	if err := json.Unmarshal(responseBody, &resource); err != nil {
		return nil, fmt.Errorf("failed to parse patient resource: %w", err)
	}

	patientInfo := &PatientInfo{
		FHIRResourceID: fhirResourceID,
	}

	if names, ok := resource["name"].([]interface{}); ok && len(names) > 0 {
		if nameMap, ok := names[0].(map[string]interface{}); ok {
			family, _ := nameMap["family"].(string)
			givenRaw, _ := nameMap["given"].([]interface{})
			given := ""
			if len(givenRaw) > 0 {
				given, _ = givenRaw[0].(string)
			}
			patientInfo.FullName = strings.TrimSpace(given + " " + family)
		}
	}

	if birthDate, ok := resource["birthDate"].(string); ok {
		patientInfo.BirthDate = birthDate
	}

	if identifiers, ok := resource["identifier"].([]interface{}); ok {
		for _, identifier := range identifiers {
			if identifierMap, ok := identifier.(map[string]interface{}); ok {
				if value, ok := identifierMap["value"].(string); ok {
					patientInfo.DocumentID = value
					break
				}
			}
		}
	}

	return patientInfo, nil
}

func (portalRepository *repository) GetEncountersByPatient(ctx context.Context, patientFHIRID string) ([]PortalEncounter, error) {
	queryParams := url.Values{"subject": {fmt.Sprintf("Patient/%s", patientFHIRID)}, "_sort": {"-date"}}.Encode()
	responseBody, err := portalRepository.fhirClient.SearchResources(ctx, "Encounter", queryParams)
	if err != nil {
		return nil, fmt.Errorf("failed to search encounters: %w", err)
	}

	return parseEncounterPortalBundle(responseBody)
}

func parseEncounterPortalBundle(responseBody json.RawMessage) ([]PortalEncounter, error) {
	entries, err := extractBundleEntries(responseBody)
	if err != nil {
		return nil, err
	}

	result := make([]PortalEncounter, 0, len(entries))
	for _, resource := range entries {
		encounter := PortalEncounter{}
		encounter.FHIRResourceID, _ = resource["id"].(string)
		encounter.Status, _ = resource["status"].(string)

		if period, ok := resource["period"].(map[string]interface{}); ok {
			if startStr, ok := period["start"].(string); ok {
				if parsed, parseErr := time.Parse(time.RFC3339, startStr); parseErr == nil {
					encounter.StartedAt = parsed
				}
			}
			if endStr, ok := period["end"].(string); ok {
				if parsed, parseErr := time.Parse(time.RFC3339, endStr); parseErr == nil {
					encounter.EndedAt = &parsed
				}
			}
		}

		if reasonCode, ok := resource["reasonCode"].([]interface{}); ok && len(reasonCode) > 0 {
			if firstReason, ok := reasonCode[0].(map[string]interface{}); ok {
				if coding, ok := firstReason["coding"].([]interface{}); ok && len(coding) > 0 {
					if firstCoding, ok := coding[0].(map[string]interface{}); ok {
						if display, ok := firstCoding["display"].(string); ok {
							encounter.ReasonDisplay = display
						}
					}
				}
			}
		}

		result = append(result, encounter)
	}

	return result, nil
}

func (portalRepository *repository) GetObservationsByPatient(ctx context.Context, patientFHIRID string) ([]PortalObservation, error) {
	queryParams := url.Values{"subject": {fmt.Sprintf("Patient/%s", patientFHIRID)}, "_sort": {"-date"}, "_count": {"50"}}.Encode()
	responseBody, err := portalRepository.fhirClient.SearchResources(ctx, "Observation", queryParams)
	if err != nil {
		return nil, fmt.Errorf("failed to search observations: %w", err)
	}

	return parseObservationPortalBundle(responseBody)
}

func parseObservationPortalBundle(responseBody json.RawMessage) ([]PortalObservation, error) {
	entries, err := extractBundleEntries(responseBody)
	if err != nil {
		return nil, err
	}

	result := make([]PortalObservation, 0, len(entries))
	for _, resource := range entries {
		observation := PortalObservation{}
		observation.FHIRResourceID, _ = resource["id"].(string)

		if code, ok := resource["code"].(map[string]interface{}); ok {
			if coding, ok := code["coding"].([]interface{}); ok && len(coding) > 0 {
				if firstCoding, ok := coding[0].(map[string]interface{}); ok {
					observation.LoincCode, _ = firstCoding["code"].(string)
					observation.CodeDisplay, _ = firstCoding["display"].(string)
				}
			}
			if text, ok := code["text"].(string); ok && observation.CodeDisplay == "" {
				observation.CodeDisplay = text
			}
		}

		if valueQuantity, ok := resource["valueQuantity"].(map[string]interface{}); ok {
			if value, ok := valueQuantity["value"].(float64); ok {
				observation.ValueQuantity = value
			}
			if unit, ok := valueQuantity["unit"].(string); ok {
				observation.ValueUnit = unit
			}
		}

		if effectiveDateTime, ok := resource["effectiveDateTime"].(string); ok {
			if parsed, parseErr := time.Parse(time.RFC3339, effectiveDateTime); parseErr == nil {
				observation.ObservedAt = parsed
			}
		}

		result = append(result, observation)
	}

	return result, nil
}

func (portalRepository *repository) GetConditionsByPatient(ctx context.Context, patientFHIRID string) ([]PortalCondition, error) {
	queryParams := url.Values{"subject": {fmt.Sprintf("Patient/%s", patientFHIRID)}}.Encode()
	responseBody, err := portalRepository.fhirClient.SearchResources(ctx, "Condition", queryParams)
	if err != nil {
		return nil, fmt.Errorf("failed to search conditions: %w", err)
	}

	return parseConditionPortalBundle(responseBody)
}

func parseConditionPortalBundle(responseBody json.RawMessage) ([]PortalCondition, error) {
	entries, err := extractBundleEntries(responseBody)
	if err != nil {
		return nil, err
	}

	result := make([]PortalCondition, 0, len(entries))
	for _, resource := range entries {
		condition := PortalCondition{}
		condition.FHIRResourceID, _ = resource["id"].(string)
		condition.ClinicalStatus, _ = resource["clinicalStatus"].(string)

		if code, ok := resource["code"].(map[string]interface{}); ok {
			if coding, ok := code["coding"].([]interface{}); ok && len(coding) > 0 {
				if firstCoding, ok := coding[0].(map[string]interface{}); ok {
					condition.ICD10Code, _ = firstCoding["code"].(string)
					condition.CodeDisplay, _ = firstCoding["display"].(string)
				}
			}
		}

		if onsetDateTime, ok := resource["onsetDateTime"].(string); ok {
			condition.OnsetAt = onsetDateTime
		}

		result = append(result, condition)
	}

	return result, nil
}

func (portalRepository *repository) GetMedicationsByPatient(ctx context.Context, patientFHIRID string) ([]PortalMedication, error) {
	queryParams := url.Values{"subject": {fmt.Sprintf("Patient/%s", patientFHIRID)}, "_sort": {"-authoredon"}}.Encode()
	responseBody, err := portalRepository.fhirClient.SearchResources(ctx, "MedicationRequest", queryParams)
	if err != nil {
		return nil, fmt.Errorf("failed to search medication requests: %w", err)
	}

	return parseMedicationPortalBundle(responseBody)
}

func parseMedicationPortalBundle(responseBody json.RawMessage) ([]PortalMedication, error) {
	entries, err := extractBundleEntries(responseBody)
	if err != nil {
		return nil, err
	}

	result := make([]PortalMedication, 0, len(entries))
	for _, resource := range entries {
		medication := PortalMedication{}
		medication.FHIRResourceID, _ = resource["id"].(string)
		medication.Status, _ = resource["status"].(string)

		if medicationCodeable, ok := resource["medicationCodeableConcept"].(map[string]interface{}); ok {
			if coding, ok := medicationCodeable["coding"].([]interface{}); ok && len(coding) > 0 {
				if firstCoding, ok := coding[0].(map[string]interface{}); ok {
					medication.MedicationName, _ = firstCoding["display"].(string)
				}
			}
			if text, ok := medicationCodeable["text"].(string); ok && medication.MedicationName == "" {
				medication.MedicationName = text
			}
		}

		if dosageInstruction, ok := resource["dosageInstruction"].([]interface{}); ok && len(dosageInstruction) > 0 {
			if firstDosage, ok := dosageInstruction[0].(map[string]interface{}); ok {
				if text, ok := firstDosage["text"].(string); ok {
					medication.DosageInstructions = text
				}
			}
		}

		if authoredOn, ok := resource["authoredOn"].(string); ok {
			medication.IssuedAt = authoredOn
		}

		result = append(result, medication)
	}

	return result, nil
}

func (portalRepository *repository) GetReportsByPatient(ctx context.Context, patientFHIRID string) ([]PortalReport, error) {
	queryParams := url.Values{"subject": {fmt.Sprintf("Patient/%s", patientFHIRID)}, "_sort": {"-date"}}.Encode()
	responseBody, err := portalRepository.fhirClient.SearchResources(ctx, "DiagnosticReport", queryParams)
	if err != nil {
		return nil, fmt.Errorf("failed to search diagnostic reports: %w", err)
	}

	return parseReportPortalBundle(responseBody)
}

func parseReportPortalBundle(responseBody json.RawMessage) ([]PortalReport, error) {
	entries, err := extractBundleEntries(responseBody)
	if err != nil {
		return nil, err
	}

	result := make([]PortalReport, 0, len(entries))
	for _, resource := range entries {
		report := PortalReport{}
		report.FHIRResourceID, _ = resource["id"].(string)
		report.Status, _ = resource["status"].(string)

		if conclusion, ok := resource["conclusion"].(string); ok {
			report.Conclusion = conclusion
		}

		if code, ok := resource["code"].(map[string]interface{}); ok {
			if coding, ok := code["coding"].([]interface{}); ok && len(coding) > 0 {
				if firstCoding, ok := coding[0].(map[string]interface{}); ok {
					report.ReportDisplay, _ = firstCoding["display"].(string)
				}
			}
			if text, ok := code["text"].(string); ok && report.ReportDisplay == "" {
				report.ReportDisplay = text
			}
		}

		if issued, ok := resource["issued"].(string); ok {
			report.IssuedAt = issued
		}

		result = append(result, report)
	}

	return result, nil
}

func (portalRepository *repository) GetImagingByPatient(ctx context.Context, patientFHIRID string) ([]PortalImaging, error) {
	queryParams := url.Values{"subject": {fmt.Sprintf("Patient/%s", patientFHIRID)}, "_sort": {"-started"}}.Encode()
	responseBody, err := portalRepository.fhirClient.SearchResources(ctx, "ImagingStudy", queryParams)
	if err != nil {
		return nil, fmt.Errorf("failed to search imaging studies: %w", err)
	}

	return parseImagingPortalBundle(responseBody)
}

func parseImagingPortalBundle(responseBody json.RawMessage) ([]PortalImaging, error) {
	entries, err := extractBundleEntries(responseBody)
	if err != nil {
		return nil, err
	}

	result := make([]PortalImaging, 0, len(entries))
	for _, resource := range entries {
		imaging := PortalImaging{}
		imaging.FHIRResourceID, _ = resource["id"].(string)
		imaging.Status, _ = resource["status"].(string)

		if started, ok := resource["started"].(string); ok {
			imaging.CreatedAt = started
		}

		imaging.Title, _ = resource["description"].(string)

		if modality, ok := resource["modality"].([]interface{}); ok && len(modality) > 0 {
			if firstModality, ok := modality[0].(map[string]interface{}); ok {
				imaging.Modality, _ = firstModality["display"].(string)
				if imaging.Modality == "" {
					imaging.Modality, _ = firstModality["code"].(string)
				}
			}
		}

		result = append(result, imaging)
	}

	return result, nil
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
