package observation

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
	CreateObservation(ctx context.Context, observation *Observation) (*Observation, error)
	CreateObservationBatch(ctx context.Context, batch []*Observation) ([]*Observation, error)
	GetObservationsByEncounter(ctx context.Context, encounterFHIRID string) ([]*Observation, error)
	GetObservationsByPatient(ctx context.Context, patientFHIRID string) ([]*Observation, error)
	UpdateObservation(ctx context.Context, fhirResourceID string, observation *Observation) (*Observation, error)
	DeleteObservation(ctx context.Context, fhirResourceID string) error
}

type repository struct {
	fhirClient healthcare.FHIRClient
}

func NewRepository(fhirClient healthcare.FHIRClient) Repository {
	return &repository{fhirClient: fhirClient}
}

func (observationRepository *repository) CreateObservation(ctx context.Context, observation *Observation) (*Observation, error) {
	fhirObservation := fhir.NewObservationResource(
		observation.PatientFHIRID,
		observation.EncounterFHIRID,
		observation.LoincCode,
		observation.CodeDisplay,
		observation.ValueQuantity,
		observation.ValueUnit,
	)

	responseBody, err := observationRepository.fhirClient.CreateResource(ctx, "Observation", fhirObservation)
	if err != nil {
		if healthcare.IsNotFound(err) {
			return nil, apperrors.ErrObservationNotFound
		}
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

func (observationRepository *repository) CreateObservationBatch(ctx context.Context, batch []*Observation) ([]*Observation, error) {
	createdObservations := make([]*Observation, 0, len(batch))
	for _, observationEntity := range batch {
		var fhirObservation *fhir.ObservationResource
		if observationEntity.NotPerformed {
			fhirObservation = fhir.NewNotPerformedObservationResource(
				observationEntity.PatientFHIRID,
				observationEntity.EncounterFHIRID,
				observationEntity.LoincCode,
				observationEntity.CodeDisplay,
			)
		} else {
			fhirObservation = fhir.NewObservationResource(
				observationEntity.PatientFHIRID,
				observationEntity.EncounterFHIRID,
				observationEntity.LoincCode,
				observationEntity.CodeDisplay,
				observationEntity.ValueQuantity,
				observationEntity.ValueUnit,
			)
		}

		responseBody, createErr := observationRepository.fhirClient.CreateResource(ctx, "Observation", fhirObservation)
		if createErr != nil {
			if healthcare.IsNotFound(createErr) {
				return nil, apperrors.ErrObservationNotFound
			}
			return nil, fmt.Errorf("failed to create observation batch: %w", createErr)
		}

		var createdResource map[string]interface{}
		if unmarshalErr := json.Unmarshal(responseBody, &createdResource); unmarshalErr != nil {
			return nil, fmt.Errorf("failed to parse observation response: %w", unmarshalErr)
		}

		fhirID, _ := createdResource["id"].(string)
		observationEntity.FHIRResourceID = fhirID
		createdObservations = append(createdObservations, observationEntity)
	}
	return createdObservations, nil
}

func (observationRepository *repository) GetObservationsByEncounter(ctx context.Context, encounterFHIRID string) ([]*Observation, error) {
	queryParams := url.Values{"encounter": []string{fmt.Sprintf("Encounter/%s", encounterFHIRID)}}.Encode()
	responseBody, err := observationRepository.fhirClient.SearchResources(ctx, "Observation", queryParams)
	if err != nil {
		if healthcare.IsNotFound(err) {
			return nil, apperrors.ErrObservationNotFound
		}
		return nil, fmt.Errorf("failed to search observations: %w", err)
	}
	return parseObservationBundle(responseBody)
}

func (observationRepository *repository) GetObservationsByPatient(ctx context.Context, patientFHIRID string) ([]*Observation, error) {
	queryParams := url.Values{"subject": []string{fmt.Sprintf("Patient/%s", patientFHIRID)}}.Encode()
	responseBody, err := observationRepository.fhirClient.SearchResources(ctx, "Observation", queryParams)
	if err != nil {
		if healthcare.IsNotFound(err) {
			return nil, apperrors.ErrObservationNotFound
		}
		return nil, fmt.Errorf("failed to search observations: %w", err)
	}
	return parseObservationBundle(responseBody)
}

func (observationRepository *repository) UpdateObservation(ctx context.Context, fhirResourceID string, observation *Observation) (*Observation, error) {
	fhirObservation := fhir.NewObservationResource(
		observation.PatientFHIRID,
		observation.EncounterFHIRID,
		observation.LoincCode,
		observation.CodeDisplay,
		observation.ValueQuantity,
		observation.ValueUnit,
	)

	responseBody, err := observationRepository.fhirClient.UpdateResource(ctx, "Observation", fhirResourceID, fhirObservation)
	if err != nil {
		if healthcare.IsNotFound(err) {
			return nil, apperrors.ErrObservationNotFound
		}
		return nil, fmt.Errorf("failed to update observation: %w", err)
	}

	var updatedResource map[string]interface{}
	if err := json.Unmarshal(responseBody, &updatedResource); err != nil {
		return nil, fmt.Errorf("failed to parse observation response: %w", err)
	}

	fhirID, _ := updatedResource["id"].(string)
	observation.FHIRResourceID = fhirID
	return observation, nil
}

func (observationRepository *repository) DeleteObservation(ctx context.Context, fhirResourceID string) error {
	if err := observationRepository.fhirClient.DeleteResource(ctx, "Observation/"+fhirResourceID); err != nil {
		if healthcare.IsNotFound(err) {
			return apperrors.ErrObservationNotFound
		}
		return fmt.Errorf("failed to delete observation: %w", err)
	}
	return nil
}

func parseObservationBundle(responseBody json.RawMessage) ([]*Observation, error) {
	decodedResources, err := fhir.DecodeBundle[fhir.Observation](responseBody)
	if err != nil {
		return nil, err
	}
	observations := make([]*Observation, 0, len(decodedResources))
	for _, resource := range decodedResources {
		observation := &Observation{}
		observation.FHIRResourceID = resource.ID
		code, display, text := fhir.CodeableConceptParts(resource.Code)
		observation.CodeDisplay = text
		observation.LoincCode = code
		if observation.CodeDisplay == "" {
			observation.CodeDisplay = display
		}
		if resource.ValueQuantity != nil {
			observation.ValueQuantity = resource.ValueQuantity.Value
			observation.ValueUnit = resource.ValueQuantity.Unit
		} else if resource.DataAbsentReason != nil {
			observation.NotPerformed = true
		}
		observation.EncounterFHIRID = fhir.SplitReferenceID(resource.Encounter.Reference)
		observation.PatientFHIRID = fhir.SplitReferenceID(resource.Subject.Reference)
		if parsedTime, ok := fhir.ParseRFC3339(resource.EffectiveDateTime); ok {
			observation.ObservedAt = parsedTime
		}
		if observation.ObservedAt.IsZero() {
			if parsedTime, ok := fhir.ParseRFC3339(resource.Issued); ok {
				observation.ObservedAt = parsedTime
			}
		}
		observations = append(observations, observation)
	}
	return observations, nil
}
