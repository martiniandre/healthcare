package portal

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/healthcare/backend/internal/shared/fhir"
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

	decodedResource, err := fhir.DecodeResource[fhir.Patient](responseBody)
	if err != nil {
		return nil, fmt.Errorf("failed to parse patient resource: %w", err)
	}

	patientInfo := &PatientInfo{
		FHIRResourceID: fhirResourceID,
	}

	if len(decodedResource.Name) > 0 {
		name := decodedResource.Name[0]
		given := ""
		if len(name.Given) > 0 {
			given = name.Given[0]
		}
		patientInfo.FullName = strings.TrimSpace(given + " " + name.Family)
	}

	patientInfo.BirthDate = decodedResource.BirthDate

	if len(decodedResource.Identifier) > 0 {
		patientInfo.DocumentID = decodedResource.Identifier[0].Value
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
	decodedResources, err := fhir.DecodeBundle[fhir.Encounter](responseBody)
	if err != nil {
		return nil, err
	}

	result := make([]PortalEncounter, 0, len(decodedResources))
	for _, resource := range decodedResources {
		encounter := PortalEncounter{}
		encounter.FHIRResourceID = resource.ID
		encounter.Status = resource.Status

		if resource.Period != nil {
			if parsed, ok := fhir.ParseRFC3339(resource.Period.Start); ok {
				encounter.StartedAt = parsed
			}
			if parsed, ok := fhir.ParseRFC3339(resource.Period.End); ok {
				encounter.EndedAt = &parsed
			}
		}

		if len(resource.ReasonCode) > 0 {
			_, display, text := fhir.CodeableConceptParts(resource.ReasonCode[0])
			encounter.ReasonDisplay = display
			if encounter.ReasonDisplay == "" {
				encounter.ReasonDisplay = text
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
	decodedResources, err := fhir.DecodeBundle[fhir.Observation](responseBody)
	if err != nil {
		return nil, err
	}

	result := make([]PortalObservation, 0, len(decodedResources))
	for _, resource := range decodedResources {
		observation := PortalObservation{}
		observation.FHIRResourceID = resource.ID

		code, display, text := fhir.CodeableConceptParts(resource.Code)
		observation.LoincCode = code
		observation.CodeDisplay = display
		if observation.CodeDisplay == "" {
			observation.CodeDisplay = text
		}

		if resource.ValueQuantity != nil {
			observation.ValueQuantity = resource.ValueQuantity.Value
			observation.ValueUnit = resource.ValueQuantity.Unit
		}
		observation.NotPerformed = resource.DataAbsentReason != nil

		if parsed, ok := fhir.ParseRFC3339(resource.EffectiveDateTime); ok {
			observation.ObservedAt = parsed
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
	decodedResources, err := fhir.DecodeBundle[fhir.Condition](responseBody)
	if err != nil {
		return nil, err
	}

	result := make([]PortalCondition, 0, len(decodedResources))
	for _, resource := range decodedResources {
		condition := PortalCondition{}
		condition.FHIRResourceID = resource.ID
		clinicalStatusCode, _, _ := fhir.CodeableConceptParts(resource.ClinicalStatus)
		condition.ClinicalStatus = clinicalStatusCode

		code, display, _ := fhir.CodeableConceptParts(resource.Code)
		condition.ICD10Code = code
		condition.CodeDisplay = display

		condition.OnsetAt = resource.OnsetDateTime

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
	decodedResources, err := fhir.DecodeBundle[fhir.MedicationRequest](responseBody)
	if err != nil {
		return nil, err
	}

	result := make([]PortalMedication, 0, len(decodedResources))
	for _, resource := range decodedResources {
		medication := PortalMedication{}
		medication.FHIRResourceID = resource.ID
		medication.Status = resource.Status

		_, display, text := fhir.CodeableConceptParts(resource.MedicationCodeableConcept)
		medication.MedicationName = display
		if medication.MedicationName == "" {
			medication.MedicationName = text
		}

		if len(resource.DosageInstruction) > 0 {
			medication.DosageInstructions = resource.DosageInstruction[0].Text
		}

		medication.IssuedAt = resource.AuthoredOn

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
	decodedResources, err := fhir.DecodeBundle[fhir.DiagnosticReport](responseBody)
	if err != nil {
		return nil, err
	}

	result := make([]PortalReport, 0, len(decodedResources))
	for _, resource := range decodedResources {
		report := PortalReport{}
		report.FHIRResourceID = resource.ID
		report.Status = resource.Status
		report.Conclusion = resource.Conclusion

		_, display, text := fhir.CodeableConceptParts(resource.Code)
		report.ReportDisplay = display
		if report.ReportDisplay == "" {
			report.ReportDisplay = text
		}

		report.IssuedAt = resource.Issued
		report.Version = fhir.ResourceVersionID(resource.Meta)

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
	decodedResources, err := fhir.DecodeBundle[fhir.ImagingStudy](responseBody)
	if err != nil {
		return nil, err
	}

	result := make([]PortalImaging, 0, len(decodedResources))
	for _, resource := range decodedResources {
		imaging := PortalImaging{}
		imaging.FHIRResourceID = resource.ID
		imaging.Status = resource.Status
		imaging.CreatedAt = resource.Started
		imaging.Title = resource.Description

		if len(resource.Modality) > 0 {
			imaging.Modality = resource.Modality[0].Display
			if imaging.Modality == "" {
				imaging.Modality = resource.Modality[0].Code
			}
		}

		result = append(result, imaging)
	}

	return result, nil
}
