package clinical

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
	CreateEncounter(ctx context.Context, encounter *Encounter) (*Encounter, error)
	GetEncountersByPatient(ctx context.Context, patientFHIRID string) ([]*Encounter, error)

	CreateObservation(ctx context.Context, observation *Observation) (*Observation, error)
	GetObservationsByEncounter(ctx context.Context, encounterFHIRID string) ([]*Observation, error)
	GetObservationsByPatient(ctx context.Context, patientFHIRID string) ([]*Observation, error)

	CreateCondition(ctx context.Context, condition *Condition) (*Condition, error)
	GetConditionsByPatient(ctx context.Context, patientFHIRID string) ([]*Condition, error)

	CreateAllergyIntolerance(ctx context.Context, allergy *AllergyIntolerance) (*AllergyIntolerance, error)
	GetAllergyIntolerancesByPatient(ctx context.Context, patientFHIRID string) ([]*AllergyIntolerance, error)

	CreateMedicationRequest(ctx context.Context, medicationRequest *MedicationRequest) (*MedicationRequest, error)
	GetMedicationRequestsByEncounter(ctx context.Context, encounterFHIRID string) ([]*MedicationRequest, error)

	CreateDiagnosticReport(ctx context.Context, report *DiagnosticReport) (*DiagnosticReport, error)
	GetDiagnosticReportsByEncounter(ctx context.Context, encounterFHIRID string) ([]*DiagnosticReport, error)
}

type repository struct {
	fhirClient healthcare.FHIRClient
}

func NewRepository(fhirClient healthcare.FHIRClient) Repository {
	return &repository{fhirClient: fhirClient}
}

func (clinicalRepo *repository) CreateEncounter(ctx context.Context, encounter *Encounter) (*Encounter, error) {
	fhirEncounter := fhir.NewEncounterResource(encounter.PatientFHIRID, encounter.PractitionerID, encounter.ReasonCode, encounter.ReasonDisplay)

	responseBody, err := clinicalRepo.fhirClient.CreateResource(ctx, "Encounter", fhirEncounter)
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

func (clinicalRepo *repository) GetEncountersByPatient(ctx context.Context, patientFHIRID string) ([]*Encounter, error) {
	queryParams := url.Values{"subject": []string{fmt.Sprintf("Patient/%s", patientFHIRID)}}.Encode()
	responseBody, err := clinicalRepo.fhirClient.SearchResources(ctx, "Encounter", queryParams)
	if err != nil {
		return nil, fmt.Errorf("failed to search encounters: %w", err)
	}
	return parseEncounterBundle(responseBody)
}

func (clinicalRepo *repository) CreateObservation(ctx context.Context, observation *Observation) (*Observation, error) {
	fhirObservation := fhir.NewObservationResource(
		observation.PatientFHIRID,
		observation.EncounterFHIRID,
		observation.LoincCode,
		observation.CodeDisplay,
		observation.ValueQuantity,
		observation.ValueUnit,
	)

	responseBody, err := clinicalRepo.fhirClient.CreateResource(ctx, "Observation", fhirObservation)
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

func (clinicalRepo *repository) GetObservationsByEncounter(ctx context.Context, encounterFHIRID string) ([]*Observation, error) {
	queryParams := url.Values{"encounter": []string{fmt.Sprintf("Encounter/%s", encounterFHIRID)}}.Encode()
	responseBody, err := clinicalRepo.fhirClient.SearchResources(ctx, "Observation", queryParams)
	if err != nil {
		return nil, fmt.Errorf("failed to search observations: %w", err)
	}
	return parseObservationBundle(responseBody)
}

func (clinicalRepo *repository) GetObservationsByPatient(ctx context.Context, patientFHIRID string) ([]*Observation, error) {
	queryParams := url.Values{"subject": []string{fmt.Sprintf("Patient/%s", patientFHIRID)}}.Encode()
	responseBody, err := clinicalRepo.fhirClient.SearchResources(ctx, "Observation", queryParams)
	if err != nil {
		return nil, fmt.Errorf("failed to search observations: %w", err)
	}
	return parseObservationBundle(responseBody)
}

func (clinicalRepo *repository) CreateCondition(ctx context.Context, condition *Condition) (*Condition, error) {
	fhirCondition := fhir.NewConditionResource(
		condition.PatientFHIRID,
		condition.EncounterFHIRID,
		condition.ICD10Code,
		condition.CodeDisplay,
		condition.ClinicalStatus,
		condition.OnsetAt,
	)

	responseBody, err := clinicalRepo.fhirClient.CreateResource(ctx, "Condition", fhirCondition)
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

func (clinicalRepo *repository) GetConditionsByPatient(ctx context.Context, patientFHIRID string) ([]*Condition, error) {
	queryParams := url.Values{"subject": []string{fmt.Sprintf("Patient/%s", patientFHIRID)}}.Encode()
	responseBody, err := clinicalRepo.fhirClient.SearchResources(ctx, "Condition", queryParams)
	if err != nil {
		return nil, fmt.Errorf("failed to search conditions: %w", err)
	}
	return parseConditionBundle(responseBody)
}

func (clinicalRepo *repository) CreateAllergyIntolerance(ctx context.Context, allergy *AllergyIntolerance) (*AllergyIntolerance, error) {
	fhirAllergy := fhir.NewAllergyIntoleranceResource(
		allergy.PatientFHIRID,
		allergy.AllergenCode,
		allergy.AllergenDisplay,
		allergy.ClinicalStatus,
		allergy.Reaction,
	)

	responseBody, err := clinicalRepo.fhirClient.CreateResource(ctx, "AllergyIntolerance", fhirAllergy)
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

func (clinicalRepo *repository) GetAllergyIntolerancesByPatient(ctx context.Context, patientFHIRID string) ([]*AllergyIntolerance, error) {
	queryParams := url.Values{"patient": []string{fmt.Sprintf("Patient/%s", patientFHIRID)}}.Encode()
	responseBody, err := clinicalRepo.fhirClient.SearchResources(ctx, "AllergyIntolerance", queryParams)
	if err != nil {
		return nil, fmt.Errorf("failed to search allergy intolerances: %w", err)
	}
	return parseAllergyBundle(responseBody)
}

func (clinicalRepo *repository) CreateMedicationRequest(ctx context.Context, medicationRequest *MedicationRequest) (*MedicationRequest, error) {
	fhirMedication := fhir.NewMedicationRequestResource(
		medicationRequest.PatientFHIRID,
		medicationRequest.EncounterFHIRID,
		medicationRequest.PractitionerFHIRID,
		medicationRequest.MedicationCode,
		medicationRequest.MedicationName,
		medicationRequest.DosageInstructions,
	)

	responseBody, err := clinicalRepo.fhirClient.CreateResource(ctx, "MedicationRequest", fhirMedication)
	if err != nil {
		return nil, fmt.Errorf("failed to create medication request: %w", err)
	}

	var createdResource map[string]interface{}
	if err := json.Unmarshal(responseBody, &createdResource); err != nil {
		return nil, fmt.Errorf("failed to parse medication request response: %w", err)
	}

	fhirID, _ := createdResource["id"].(string)
	medicationRequest.FHIRResourceID = fhirID
	medicationRequest.IssuedAt = time.Now()
	return medicationRequest, nil
}

func (clinicalRepo *repository) GetMedicationRequestsByEncounter(ctx context.Context, encounterFHIRID string) ([]*MedicationRequest, error) {
	queryParams := url.Values{"encounter": []string{fmt.Sprintf("Encounter/%s", encounterFHIRID)}}.Encode()
	responseBody, err := clinicalRepo.fhirClient.SearchResources(ctx, "MedicationRequest", queryParams)
	if err != nil {
		return nil, fmt.Errorf("failed to search medication requests: %w", err)
	}
	return parseMedicationRequestBundle(responseBody)
}

func (clinicalRepo *repository) CreateDiagnosticReport(ctx context.Context, report *DiagnosticReport) (*DiagnosticReport, error) {
	fhirReport := fhir.NewDiagnosticReportResource(
		report.PatientFHIRID,
		report.EncounterFHIRID,
		report.ReportCode,
		report.ReportDisplay,
		report.Conclusion,
	)

	responseBody, err := clinicalRepo.fhirClient.CreateResource(ctx, "DiagnosticReport", fhirReport)
	if err != nil {
		return nil, fmt.Errorf("failed to create diagnostic report: %w", err)
	}

	var createdResource map[string]interface{}
	if err := json.Unmarshal(responseBody, &createdResource); err != nil {
		return nil, fmt.Errorf("failed to parse diagnostic report response: %w", err)
	}

	fhirID, _ := createdResource["id"].(string)
	report.FHIRResourceID = fhirID
	report.IssuedAt = time.Now()
	return report, nil
}

func (clinicalRepo *repository) GetDiagnosticReportsByEncounter(ctx context.Context, encounterFHIRID string) ([]*DiagnosticReport, error) {
	queryParams := url.Values{"encounter": []string{fmt.Sprintf("Encounter/%s", encounterFHIRID)}}.Encode()
	responseBody, err := clinicalRepo.fhirClient.SearchResources(ctx, "DiagnosticReport", queryParams)
	if err != nil {
		return nil, fmt.Errorf("failed to search diagnostic reports: %w", err)
	}
	return parseDiagnosticReportBundle(responseBody)
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

func parseAllergyBundle(responseBody json.RawMessage) ([]*AllergyIntolerance, error) {
	entries, err := extractBundleEntries(responseBody)
	if err != nil {
		return nil, err
	}
	allergies := make([]*AllergyIntolerance, 0, len(entries))
	for _, resource := range entries {
		allergy := &AllergyIntolerance{}
		allergy.FHIRResourceID, _ = resource["id"].(string)
		if codes, ok := resource["code"].(map[string]interface{}); ok {
			allergy.AllergenDisplay, _ = codes["text"].(string)
		}
		allergies = append(allergies, allergy)
	}
	return allergies, nil
}

func parseMedicationRequestBundle(responseBody json.RawMessage) ([]*MedicationRequest, error) {
	entries, err := extractBundleEntries(responseBody)
	if err != nil {
		return nil, err
	}
	medications := make([]*MedicationRequest, 0, len(entries))
	for _, resource := range entries {
		medication := &MedicationRequest{}
		medication.FHIRResourceID, _ = resource["id"].(string)
		medication.Status, _ = resource["status"].(string)
		if med, ok := resource["medicationCodeableConcept"].(map[string]interface{}); ok {
			medication.MedicationName, _ = med["text"].(string)
		}
		medications = append(medications, medication)
	}
	return medications, nil
}

func parseDiagnosticReportBundle(responseBody json.RawMessage) ([]*DiagnosticReport, error) {
	entries, err := extractBundleEntries(responseBody)
	if err != nil {
		return nil, err
	}
	reports := make([]*DiagnosticReport, 0, len(entries))
	for _, resource := range entries {
		report := &DiagnosticReport{}
		report.FHIRResourceID, _ = resource["id"].(string)
		report.Status, _ = resource["status"].(string)
		report.Conclusion, _ = resource["conclusion"].(string)
		if codes, ok := resource["code"].(map[string]interface{}); ok {
			report.ReportDisplay, _ = codes["text"].(string)
		}
		reports = append(reports, report)
	}
	return reports, nil
}
