package encounter

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/healthcare/backend/internal/shared/apperrors"
	"github.com/healthcare/backend/internal/shared/fhir"
	"github.com/healthcare/backend/internal/shared/healthcare"
)

type Repository interface {
	CreateEncounter(ctx context.Context, encounter *Encounter) (*Encounter, error)
	GetEncounterByID(ctx context.Context, fhirResourceID string) (*Encounter, error)
	GetEncountersByPatient(ctx context.Context, patientFHIRID string) ([]*Encounter, error)
	UpdateEncounter(ctx context.Context, fhirResourceID string, encounter *Encounter) (*Encounter, error)
	DeleteEncounter(ctx context.Context, fhirResourceID string) error
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
		if healthcare.IsNotFound(err) {
			return nil, fmt.Errorf("failed to create encounter: %w", apperrors.ErrEncounterNotFound)
		}
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
		if healthcare.IsNotFound(err) {
			return nil, fmt.Errorf("failed to get encounter: %w", apperrors.ErrEncounterNotFound)
		}
		return nil, fmt.Errorf("failed to get encounter: %w", err)
	}

	decodedResource, err := fhir.DecodeResource[fhir.Encounter](responseBody)
	if err != nil {
		return nil, fmt.Errorf("failed to parse encounter response: %w", err)
	}
	return mapFHIREncounterToDomain(decodedResource), nil
}

func mapFHIREncounterToDomain(resource *fhir.Encounter) *Encounter {
	encounter := &Encounter{}
	encounter.FHIRResourceID = resource.ID
	encounter.Status = resource.Status
	encounter.PatientFHIRID = resource.Subject.Reference
	if resource.Period != nil {
		if parsedTime, ok := fhir.ParseRFC3339(resource.Period.Start); ok {
			encounter.StartedAt = parsedTime
		}
	}
	if len(resource.ReasonCode) > 0 {
		code, display, text := fhir.CodeableConceptParts(resource.ReasonCode[0])
		if text != "" {
			encounter.ReasonCode = text
			encounter.ReasonDisplay = text
		} else {
			encounter.ReasonCode = code
			encounter.ReasonDisplay = display
		}
	}
	if len(resource.Participant) > 0 {
		encounter.PractitionerID = fhir.SplitReferenceID(resource.Participant[0].Individual.Reference)
	}
	return encounter
}

func (encounterRepository *repository) UpdateEncounter(ctx context.Context, fhirResourceID string, encounter *Encounter) (*Encounter, error) {
	fhirEncounter := fhir.NewEncounterResource(encounter.PatientFHIRID, encounter.PractitionerID, encounter.ReasonCode, encounter.ReasonDisplay)
	fhirEncounter.ID = fhirResourceID
	fhirEncounter.Status = encounter.Status

	responseBody, err := encounterRepository.fhirClient.UpdateResource(ctx, "Encounter", fhirResourceID, fhirEncounter)
	if err != nil {
		if healthcare.IsNotFound(err) {
			return nil, fmt.Errorf("failed to update encounter: %w", apperrors.ErrEncounterNotFound)
		}
		return nil, fmt.Errorf("failed to update encounter: %w", err)
	}

	decodedResource, err := fhir.DecodeResource[fhir.Encounter](responseBody)
	if err != nil {
		return nil, fmt.Errorf("failed to parse updated encounter: %w", err)
	}
	return mapFHIREncounterToDomain(decodedResource), nil
}

func (encounterRepository *repository) DeleteEncounter(ctx context.Context, fhirResourceID string) error {
	err := encounterRepository.fhirClient.DeleteResource(ctx, "Encounter/"+fhirResourceID)
	if err != nil {
		if healthcare.IsNotFound(err) {
			return fmt.Errorf("failed to delete encounter: %w", apperrors.ErrEncounterNotFound)
		}
		return fmt.Errorf("failed to delete encounter: %w", err)
	}

	return nil
}

func (encounterRepository *repository) GetEncountersByPatient(ctx context.Context, patientFHIRID string) ([]*Encounter, error) {
	queryParams := url.Values{"subject": []string{fmt.Sprintf("Patient/%s", patientFHIRID)}}.Encode()
	responseBody, err := encounterRepository.fhirClient.SearchResources(ctx, "Encounter", queryParams)
	if err != nil {
		return nil, fmt.Errorf("failed to search encounters: %w", err)
	}

	return parseEncounterBundle(responseBody)
}

func parseEncounterBundle(responseBody json.RawMessage) ([]*Encounter, error) {
	decodedResources, err := fhir.DecodeBundle[fhir.Encounter](responseBody)
	if err != nil {
		return nil, err
	}
	encounters := make([]*Encounter, 0, len(decodedResources))
	for _, resource := range decodedResources {
		encounter := mapFHIREncounterToDomain(&resource)
		encounters = append(encounters, encounter)
	}
	return encounters, nil
}
