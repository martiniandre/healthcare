package timeline

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
	PatientExists(ctx context.Context, patientFHIRID string) error
	FetchEncounters(ctx context.Context, patientFHIRID string, before *time.Time, limit int) ([]TimelineEntry, error)
	FetchObservations(ctx context.Context, patientFHIRID string, before *time.Time, limit int) ([]TimelineEntry, error)
	FetchConditions(ctx context.Context, patientFHIRID string, onlyActive bool, before *time.Time, limit int) ([]TimelineEntry, error)
	FetchMedications(ctx context.Context, patientFHIRID string, onlyActive bool, before *time.Time, limit int) ([]TimelineEntry, error)
	FetchReports(ctx context.Context, patientFHIRID string, before *time.Time, limit int) ([]TimelineEntry, error)
	FetchImaging(ctx context.Context, patientFHIRID string, before *time.Time, limit int) ([]TimelineEntry, error)
	FetchAllergies(ctx context.Context, patientFHIRID string, before *time.Time, limit int) ([]TimelineEntry, error)
}

type repository struct {
	fhirClient healthcare.FHIRClient
}

func NewRepository(fhirClient healthcare.FHIRClient) Repository {
	return &repository{fhirClient: fhirClient}
}

func (timelineRepository *repository) PatientExists(ctx context.Context, patientFHIRID string) error {
	responseBody, err := timelineRepository.fhirClient.GetResource(ctx, "Patient", patientFHIRID)
	if err != nil {
		return fmt.Errorf("failed to get patient from healthcare api: %w", err)
	}

	if _, decodeErr := fhir.DecodeResource[fhir.Patient](responseBody); decodeErr != nil {
		return fmt.Errorf("failed to parse patient resource: %w", decodeErr)
	}

	return nil
}

func (timelineRepository *repository) FetchEncounters(ctx context.Context, patientFHIRID string, before *time.Time, limit int) ([]TimelineEntry, error) {
	queryParams := url.Values{"subject": {patientReference(patientFHIRID)}, "_sort": {"-date"}, "_count": {fmt.Sprintf("%d", limit)}}
	if before != nil {
		queryParams.Set("date", "lt"+before.Format(time.RFC3339))
	}

	responseBody, err := timelineRepository.fhirClient.SearchResources(ctx, "Encounter", queryParams.Encode())
	if err != nil {
		return nil, fmt.Errorf("failed to search encounters: %w", err)
	}

	return parseEncounterTimelineBundle(responseBody)
}

func (timelineRepository *repository) FetchObservations(ctx context.Context, patientFHIRID string, before *time.Time, limit int) ([]TimelineEntry, error) {
	queryParams := url.Values{"subject": {patientReference(patientFHIRID)}, "_sort": {"-date"}, "_count": {fmt.Sprintf("%d", limit)}}
	if before != nil {
		queryParams.Set("date", "lt"+before.Format(time.RFC3339))
	}

	responseBody, err := timelineRepository.fhirClient.SearchResources(ctx, "Observation", queryParams.Encode())
	if err != nil {
		return nil, fmt.Errorf("failed to search observations: %w", err)
	}

	return parseObservationTimelineBundle(responseBody)
}

func (timelineRepository *repository) FetchConditions(ctx context.Context, patientFHIRID string, onlyActive bool, before *time.Time, limit int) ([]TimelineEntry, error) {
	queryParams := url.Values{"subject": {patientReference(patientFHIRID)}, "_sort": {"-recorded-date"}, "_count": {fmt.Sprintf("%d", limit)}}
	if onlyActive {
		queryParams.Set("clinical-status", "active")
	}
	if before != nil {
		queryParams.Set("recorded-date", "lt"+before.Format(time.RFC3339))
	}

	responseBody, err := timelineRepository.fhirClient.SearchResources(ctx, "Condition", queryParams.Encode())
	if err != nil {
		return nil, fmt.Errorf("failed to search conditions: %w", err)
	}

	return parseConditionTimelineBundle(responseBody)
}

func (timelineRepository *repository) FetchMedications(ctx context.Context, patientFHIRID string, onlyActive bool, before *time.Time, limit int) ([]TimelineEntry, error) {
	queryParams := url.Values{"subject": {patientReference(patientFHIRID)}, "_sort": {"-authoredon"}, "_count": {fmt.Sprintf("%d", limit)}}
	if onlyActive {
		queryParams.Set("status", "active")
	}
	if before != nil {
		queryParams.Set("authoredon", "lt"+before.Format(time.RFC3339))
	}

	responseBody, err := timelineRepository.fhirClient.SearchResources(ctx, "MedicationRequest", queryParams.Encode())
	if err != nil {
		return nil, fmt.Errorf("failed to search medication requests: %w", err)
	}

	return parseMedicationTimelineBundle(responseBody)
}

func (timelineRepository *repository) FetchReports(ctx context.Context, patientFHIRID string, before *time.Time, limit int) ([]TimelineEntry, error) {
	queryParams := url.Values{"subject": {patientReference(patientFHIRID)}, "_sort": {"-date"}, "_count": {fmt.Sprintf("%d", limit)}}
	if before != nil {
		queryParams.Set("date", "lt"+before.Format(time.RFC3339))
	}

	responseBody, err := timelineRepository.fhirClient.SearchResources(ctx, "DiagnosticReport", queryParams.Encode())
	if err != nil {
		return nil, fmt.Errorf("failed to search diagnostic reports: %w", err)
	}

	return parseReportTimelineBundle(responseBody)
}

func (timelineRepository *repository) FetchImaging(ctx context.Context, patientFHIRID string, before *time.Time, limit int) ([]TimelineEntry, error) {
	queryParams := url.Values{"subject": {patientReference(patientFHIRID)}, "_sort": {"-started"}, "_count": {fmt.Sprintf("%d", limit)}}
	if before != nil {
		queryParams.Set("started", "lt"+before.Format(time.RFC3339))
	}

	responseBody, err := timelineRepository.fhirClient.SearchResources(ctx, "ImagingStudy", queryParams.Encode())
	if err != nil {
		return nil, fmt.Errorf("failed to search imaging studies: %w", err)
	}

	return parseImagingTimelineBundle(responseBody)
}

func (timelineRepository *repository) FetchAllergies(ctx context.Context, patientFHIRID string, before *time.Time, limit int) ([]TimelineEntry, error) {
	queryParams := url.Values{"subject": {patientReference(patientFHIRID)}, "_sort": {"-date"}, "_count": {fmt.Sprintf("%d", limit)}}
	if before != nil {
		queryParams.Set("date", "lt"+before.Format(time.RFC3339))
	}

	responseBody, err := timelineRepository.fhirClient.SearchResources(ctx, "AllergyIntolerance", queryParams.Encode())
	if err != nil {
		return nil, fmt.Errorf("failed to search allergy intolerances: %w", err)
	}

	return parseAllergyTimelineBundle(responseBody)
}

func patientReference(patientFHIRID string) string {
	return fmt.Sprintf("Patient/%s", patientFHIRID)
}

func parseEncounterTimelineBundle(responseBody json.RawMessage) ([]TimelineEntry, error) {
	decodedResources, err := fhir.DecodeBundle[fhir.Encounter](responseBody)
	if err != nil {
		return nil, err
	}

	result := make([]TimelineEntry, 0, len(decodedResources))
	for _, resource := range decodedResources {
		entry := TimelineEntry{
			ResourceType:   "Encounter",
			FHIRResourceID: resource.ID,
			Status:         resource.Status,
		}

		if len(resource.ReasonCode) > 0 {
			_, display, text := fhir.CodeableConceptParts(resource.ReasonCode[0])
			entry.Title = display
			if entry.Title == "" {
				entry.Title = text
			}
		}
		if entry.Title == "" {
			entry.Title = "Encounter"
		}

		if resource.Period != nil {
			if parsed, ok := fhir.ParseRFC3339(resource.Period.Start); ok {
				parsedTime := parsed
				entry.RecordedAt = parsedTime
				entry.ClinicalDate = &parsedTime
			}
			if parsed, ok := fhir.ParseRFC3339(resource.Period.End); ok {
				entry.PeriodEnd = &parsed
			}
		}

		result = append(result, entry)
	}

	return result, nil
}

func parseObservationTimelineBundle(responseBody json.RawMessage) ([]TimelineEntry, error) {
	decodedResources, err := fhir.DecodeBundle[fhir.Observation](responseBody)
	if err != nil {
		return nil, err
	}

	result := make([]TimelineEntry, 0, len(decodedResources))
	for _, resource := range decodedResources {
		entry := TimelineEntry{
			ResourceType:   "Observation",
			FHIRResourceID: resource.ID,
			Status:         resource.Status,
		}

		code, display, text := fhir.CodeableConceptParts(resource.Code)
		entry.Code = code
		entry.Title = display
		if entry.Title == "" {
			entry.Title = text
		}
		if entry.Title == "" {
			entry.Title = "Observation"
		}

		if resource.ValueQuantity != nil {
			valueQuantity := resource.ValueQuantity.Value
			entry.ValueQuantity = &valueQuantity
			entry.ValueUnit = resource.ValueQuantity.Unit
		}

		if len(resource.ReferenceRange) > 0 {
			referenceRange := resource.ReferenceRange[0]
			if referenceRange.Low != nil {
				referenceLow := referenceRange.Low.Value
				entry.ReferenceLow = &referenceLow
			}
			if referenceRange.High != nil {
				referenceHigh := referenceRange.High.Value
				entry.ReferenceHigh = &referenceHigh
			}
		}

		if parsed, ok := fhir.ParseRFC3339(resource.EffectiveDateTime); ok {
			parsedTime := parsed
			entry.RecordedAt = parsedTime
			entry.ClinicalDate = &parsedTime
		}

		result = append(result, entry)
	}

	return result, nil
}

func parseConditionTimelineBundle(responseBody json.RawMessage) ([]TimelineEntry, error) {
	decodedResources, err := fhir.DecodeBundle[fhir.Condition](responseBody)
	if err != nil {
		return nil, err
	}

	result := make([]TimelineEntry, 0, len(decodedResources))
	for _, resource := range decodedResources {
		clinicalStatusCode, _, _ := fhir.CodeableConceptParts(resource.ClinicalStatus)
		code, display, _ := fhir.CodeableConceptParts(resource.Code)

		entry := TimelineEntry{
			ResourceType:   "Condition",
			FHIRResourceID: resource.ID,
			Title:          display,
			Status:         clinicalStatusCode,
			Code:           code,
			OnsetDate:      resource.OnsetDateTime,
		}
		if entry.Title == "" {
			entry.Title = "Condition"
		}

		if recordedAt, ok := resolveRecordedAt(resource.RecordedDate, nil); ok {
			entry.RecordedAt = recordedAt
		} else {
			continue
		}

		result = append(result, entry)
	}

	return result, nil
}

func parseMedicationTimelineBundle(responseBody json.RawMessage) ([]TimelineEntry, error) {
	decodedResources, err := fhir.DecodeBundle[fhir.MedicationRequest](responseBody)
	if err != nil {
		return nil, err
	}

	result := make([]TimelineEntry, 0, len(decodedResources))
	for _, resource := range decodedResources {
		_, display, text := fhir.CodeableConceptParts(resource.MedicationCodeableConcept)

		entry := TimelineEntry{
			ResourceType:   "MedicationRequest",
			FHIRResourceID: resource.ID,
			Title:          display,
			Status:         resource.Status,
		}
		if entry.Title == "" {
			entry.Title = text
		}
		if entry.Title == "" {
			entry.Title = "MedicationRequest"
		}

		if len(resource.DosageInstruction) > 0 {
			entry.DosageInstructions = resource.DosageInstruction[0].Text
		}

		if recordedAt, ok := resolveRecordedAt(resource.AuthoredOn, nil); ok {
			entry.RecordedAt = recordedAt
		} else {
			continue
		}

		result = append(result, entry)
	}

	return result, nil
}

func parseReportTimelineBundle(responseBody json.RawMessage) ([]TimelineEntry, error) {
	decodedResources, err := fhir.DecodeBundle[fhir.DiagnosticReport](responseBody)
	if err != nil {
		return nil, err
	}

	result := make([]TimelineEntry, 0, len(decodedResources))
	for _, resource := range decodedResources {
		_, display, text := fhir.CodeableConceptParts(resource.Code)

		entry := TimelineEntry{
			ResourceType:   "DiagnosticReport",
			FHIRResourceID: resource.ID,
			Title:          display,
			Status:         resource.Status,
			Conclusion:     resource.Conclusion,
			Version:        fhir.ResourceVersionID(resource.Meta),
		}
		if entry.Title == "" {
			entry.Title = text
		}
		if entry.Title == "" {
			entry.Title = "DiagnosticReport"
		}

		if recordedAt, ok := resolveRecordedAt(resource.Issued, resource.Meta); ok {
			entry.RecordedAt = recordedAt
			entry.ClinicalDate = &recordedAt
		} else {
			continue
		}

		result = append(result, entry)
	}

	return result, nil
}

func parseImagingTimelineBundle(responseBody json.RawMessage) ([]TimelineEntry, error) {
	decodedResources, err := fhir.DecodeBundle[fhir.ImagingStudy](responseBody)
	if err != nil {
		return nil, err
	}

	result := make([]TimelineEntry, 0, len(decodedResources))
	for _, resource := range decodedResources {
		entry := TimelineEntry{
			ResourceType:   "ImagingStudy",
			FHIRResourceID: resource.ID,
			Status:         resource.Status,
			Title:          resource.Description,
		}

		if len(resource.Modality) > 0 {
			entry.Modality = resource.Modality[0].Display
			if entry.Modality == "" {
				entry.Modality = resource.Modality[0].Code
			}
		}
		if entry.Title == "" {
			entry.Title = entry.Modality
		}
		if entry.Title == "" {
			entry.Title = "ImagingStudy"
		}

		if recordedAt, ok := resolveRecordedAt(resource.Started, nil); ok {
			entry.RecordedAt = recordedAt
			entry.ClinicalDate = &recordedAt
		} else {
			continue
		}

		result = append(result, entry)
	}

	return result, nil
}

func parseAllergyTimelineBundle(responseBody json.RawMessage) ([]TimelineEntry, error) {
	decodedResources, err := fhir.DecodeBundle[fhir.AllergyIntolerance](responseBody)
	if err != nil {
		return nil, err
	}

	result := make([]TimelineEntry, 0, len(decodedResources))
	for _, resource := range decodedResources {
		clinicalStatusCode, _, _ := fhir.CodeableConceptParts(resource.ClinicalStatus)
		_, display, text := fhir.CodeableConceptParts(resource.Code)

		entry := TimelineEntry{
			ResourceType:   "AllergyIntolerance",
			FHIRResourceID: resource.ID,
			Title:          display,
			Status:         clinicalStatusCode,
			Criticality:    resource.Criticality,
		}
		if entry.Title == "" {
			entry.Title = text
		}
		if entry.Title == "" {
			entry.Title = "AllergyIntolerance"
		}

		if len(resource.Reaction) > 0 && len(resource.Reaction[0].Manifestation) > 0 {
			_, manifestationDisplay, manifestationText := fhir.CodeableConceptParts(resource.Reaction[0].Manifestation[0])
			entry.Reaction = manifestationDisplay
			if entry.Reaction == "" {
				entry.Reaction = manifestationText
			}
			if entry.Reaction == "" {
				entry.Reaction = resource.Reaction[0].Severity
			}
		}

		if recordedAt, ok := resolveRecordedAt(resource.RecordedDate, nil); ok {
			entry.RecordedAt = recordedAt
			entry.ClinicalDate = &recordedAt
		} else {
			continue
		}

		result = append(result, entry)
	}

	return result, nil
}

func resolveRecordedAt(primaryValue string, meta *fhir.ResourceMeta) (time.Time, bool) {
	if parsed, ok := fhir.ParseRFC3339(primaryValue); ok {
		return parsed, true
	}
	if meta != nil {
		if parsed, ok := fhir.ParseRFC3339(meta.LastUpdated); ok {
			return parsed, true
		}
	}
	return time.Time{}, false
}
