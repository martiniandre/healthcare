package timeline

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/healthcare/backend/internal/shared/fhir"
	"github.com/healthcare/backend/internal/shared/healthcare"
)

type Repository interface {
	PatientExists(ctx context.Context, patientFHIRID string) error
	ListEntries(ctx context.Context, entryType EntryType, query FetchQuery) ([]Entry, error)
}

type entryMapper func(responseBody json.RawMessage) ([]Entry, error)

type searchSpecification struct {
	resourceType string
	patientParam string
	dateParam    string
	sortValue    string
	mapEntries   entryMapper
}

var searchSpecifications = map[EntryType]searchSpecification{
	EntryTypeEncounter:   {"Encounter", "subject", "date", "-date", mapEncounters},
	EntryTypeObservation: {"Observation", "subject", "date", "-date", mapObservations},
	EntryTypeCondition:   {"Condition", "subject", "recorded-date", "-recorded-date", mapConditions},
	EntryTypeAllergy:     {"AllergyIntolerance", "patient", "date", "-date", mapAllergies},
	EntryTypeMedication:  {"MedicationRequest", "subject", "authoredon", "-authoredon", mapMedications},
	EntryTypeReport:      {"DiagnosticReport", "subject", "issued", "-issued", mapReports},
	EntryTypeImaging:     {"ImagingStudy", "subject", "started", "-started", mapImagingStudies},
}

type repository struct {
	fhirClient healthcare.FHIRClient
}

func NewRepository(fhirClient healthcare.FHIRClient) Repository {
	return &repository{fhirClient: fhirClient}
}

func (timelineRepository *repository) PatientExists(ctx context.Context, patientFHIRID string) error {
	if _, getError := timelineRepository.fhirClient.GetResource(ctx, "Patient", patientFHIRID); getError != nil {
		return fmt.Errorf("failed to get patient from healthcare api: %w", getError)
	}
	return nil
}

func (timelineRepository *repository) ListEntries(ctx context.Context, entryType EntryType, query FetchQuery) ([]Entry, error) {
	specification, known := searchSpecifications[entryType]
	if !known {
		return nil, fmt.Errorf("unsupported timeline entry type %q", entryType)
	}

	queryParams := url.Values{
		specification.patientParam: {"Patient/" + query.PatientFHIRID},
		specification.dateParam: {
			"ge" + query.Window.From.UTC().Format(time.RFC3339Nano),
			"le" + query.Window.To.UTC().Format(time.RFC3339Nano),
		},
		"_sort":  {specification.sortValue},
		"_count": {strconv.Itoa(query.MaxEntries)},
	}.Encode()

	responseBody, searchError := timelineRepository.fhirClient.SearchResourcesPage(ctx, specification.resourceType, queryParams)
	if searchError != nil {
		return nil, fmt.Errorf("failed to search %s: %w", specification.resourceType, searchError)
	}
	return specification.mapEntries(responseBody)
}

var fhirTimeLayouts = []string{time.RFC3339Nano, "2006-01-02T15:04:05", "2006-01-02"}

func parseFHIRTime(rawValue string) (time.Time, bool) {
	for _, layout := range fhirTimeLayouts {
		if parsed, parseError := time.Parse(layout, rawValue); parseError == nil {
			return parsed.UTC(), true
		}
	}
	return time.Time{}, false
}

func optionalFHIRTime(rawValue string) *time.Time {
	parsed, ok := parseFHIRTime(rawValue)
	if !ok {
		return nil
	}
	return &parsed
}

func conceptLabel(concept fhir.CodeableConcept) string {
	code, display, text := fhir.CodeableConceptParts(concept)
	if display != "" {
		return display
	}
	if text != "" {
		return text
	}
	return code
}

func conceptCode(concept fhir.CodeableConcept) string {
	code, _, _ := fhir.CodeableConceptParts(concept)
	return code
}

func mapEncounters(responseBody json.RawMessage) ([]Entry, error) {
	resources, decodeError := fhir.DecodeBundle[fhir.Encounter](responseBody)
	if decodeError != nil {
		return nil, decodeError
	}
	entries := make([]Entry, 0, len(resources))
	for _, resource := range resources {
		if resource.Period == nil {
			continue
		}
		occurredAt, ok := parseFHIRTime(resource.Period.Start)
		if !ok {
			continue
		}
		title := ""
		if len(resource.ReasonCode) > 0 {
			title = conceptLabel(resource.ReasonCode[0])
		}
		entries = append(entries, Entry{
			Type:        EntryTypeEncounter,
			ResourceID:  resource.ID,
			OccurredAt:  occurredAt,
			Title:       title,
			Status:      resource.Status,
			EncounterID: resource.ID,
		})
	}
	return entries, nil
}

func mapObservations(responseBody json.RawMessage) ([]Entry, error) {
	resources, decodeError := fhir.DecodeBundle[fhir.Observation](responseBody)
	if decodeError != nil {
		return nil, decodeError
	}
	entries := make([]Entry, 0, len(resources))
	for _, resource := range resources {
		occurredAt, ok := parseFHIRTime(resource.EffectiveDateTime)
		if !ok {
			continue
		}
		summary := ""
		if resource.ValueQuantity != nil && resource.DataAbsentReason == nil {
			summary = strings.TrimSpace(strconv.FormatFloat(resource.ValueQuantity.Value, 'f', -1, 64) + " " + resource.ValueQuantity.Unit)
		}
		entries = append(entries, Entry{
			Type:        EntryTypeObservation,
			ResourceID:  resource.ID,
			OccurredAt:  occurredAt,
			Title:       conceptLabel(resource.Code),
			Summary:     summary,
			Status:      resource.Status,
			EncounterID: fhir.SplitReferenceID(resource.Encounter.Reference),
		})
	}
	return entries, nil
}

func mapConditions(responseBody json.RawMessage) ([]Entry, error) {
	resources, decodeError := fhir.DecodeBundle[fhir.Condition](responseBody)
	if decodeError != nil {
		return nil, decodeError
	}
	entries := make([]Entry, 0, len(resources))
	for _, resource := range resources {
		occurredAt, ok := parseFHIRTime(resource.RecordedDate)
		if !ok {
			continue
		}
		summary := ""
		if conceptLabel(resource.Code) != conceptCode(resource.Code) {
			summary = conceptCode(resource.Code)
		}
		entries = append(entries, Entry{
			Type:         EntryTypeCondition,
			ResourceID:   resource.ID,
			OccurredAt:   occurredAt,
			ClinicalDate: optionalFHIRTime(resource.OnsetDateTime),
			Title:        conceptLabel(resource.Code),
			Summary:      summary,
			Status:       conceptCode(resource.ClinicalStatus),
			EncounterID:  fhir.SplitReferenceID(resource.Encounter.Reference),
		})
	}
	return entries, nil
}

func mapAllergies(responseBody json.RawMessage) ([]Entry, error) {
	resources, decodeError := fhir.DecodeBundle[fhir.AllergyIntolerance](responseBody)
	if decodeError != nil {
		return nil, decodeError
	}
	entries := make([]Entry, 0, len(resources))
	for _, resource := range resources {
		occurredAt, ok := parseFHIRTime(resource.RecordedDate)
		if !ok {
			continue
		}
		summary := ""
		if len(resource.Reaction) > 0 && len(resource.Reaction[0].Manifestation) > 0 {
			summary = conceptLabel(resource.Reaction[0].Manifestation[0])
		}
		entries = append(entries, Entry{
			Type:       EntryTypeAllergy,
			ResourceID: resource.ID,
			OccurredAt: occurredAt,
			Title:      conceptLabel(resource.Code),
			Summary:    summary,
			Status:     conceptCode(resource.ClinicalStatus),
		})
	}
	return entries, nil
}

func mapMedications(responseBody json.RawMessage) ([]Entry, error) {
	resources, decodeError := fhir.DecodeBundle[fhir.MedicationRequest](responseBody)
	if decodeError != nil {
		return nil, decodeError
	}
	entries := make([]Entry, 0, len(resources))
	for _, resource := range resources {
		occurredAt, ok := parseFHIRTime(resource.AuthoredOn)
		if !ok {
			continue
		}
		summary := ""
		if len(resource.DosageInstruction) > 0 {
			summary = resource.DosageInstruction[0].Text
		}
		entries = append(entries, Entry{
			Type:        EntryTypeMedication,
			ResourceID:  resource.ID,
			OccurredAt:  occurredAt,
			Title:       conceptLabel(resource.MedicationCodeableConcept),
			Summary:     summary,
			Status:      resource.Status,
			EncounterID: fhir.SplitReferenceID(resource.Encounter.Reference),
		})
	}
	return entries, nil
}

func mapReports(responseBody json.RawMessage) ([]Entry, error) {
	resources, decodeError := fhir.DecodeBundle[fhir.DiagnosticReport](responseBody)
	if decodeError != nil {
		return nil, decodeError
	}
	entries := make([]Entry, 0, len(resources))
	for _, resource := range resources {
		occurredAt, ok := parseFHIRTime(resource.Issued)
		if !ok {
			continue
		}
		entries = append(entries, Entry{
			Type:        EntryTypeReport,
			ResourceID:  resource.ID,
			OccurredAt:  occurredAt,
			Title:       conceptLabel(resource.Code),
			Summary:     resource.Conclusion,
			Status:      resource.Status,
			EncounterID: fhir.SplitReferenceID(resource.Encounter.Reference),
		})
	}
	return entries, nil
}

func mapImagingStudies(responseBody json.RawMessage) ([]Entry, error) {
	resources, decodeError := fhir.DecodeBundle[fhir.ImagingStudyResource](responseBody)
	if decodeError != nil {
		return nil, decodeError
	}
	entries := make([]Entry, 0, len(resources))
	for _, resource := range resources {
		occurredAt, ok := parseFHIRTime(resource.Started)
		if !ok {
			continue
		}
		modalities := make([]string, 0, len(resource.Series))
		for _, series := range resource.Series {
			if series.Modality.Code != "" {
				modalities = append(modalities, series.Modality.Code)
			}
		}
		title := resource.Description
		if title == "" {
			title = strings.Join(modalities, ", ")
		}
		entries = append(entries, Entry{
			Type:       EntryTypeImaging,
			ResourceID: resource.ID,
			OccurredAt: occurredAt,
			Title:      title,
			Summary:    strings.Join(modalities, ", "),
			Status:     resource.Status,
		})
	}
	return entries, nil
}
