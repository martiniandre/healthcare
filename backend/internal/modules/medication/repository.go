package medication

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

func parseMedicationRequestBundle(responseBody json.RawMessage) ([]*Medication, error) {
	decodedResources, err := fhir.DecodeBundle[fhir.MedicationRequest](responseBody)
	if err != nil {
		return nil, err
	}
	medications := make([]*Medication, 0, len(decodedResources))
	for _, resource := range decodedResources {
		medication := &Medication{}
		medication.FHIRResourceID = resource.ID
		medication.Status = resource.Status
		code, display, text := fhir.CodeableConceptParts(resource.MedicationCodeableConcept)
		medication.MedicationName = text
		medication.MedicationCode = code
		if medication.MedicationName == "" {
			medication.MedicationName = display
		}
		medication.EncounterFHIRID = fhir.SplitReferenceID(resource.Encounter.Reference)
		medication.PatientFHIRID = fhir.SplitReferenceID(resource.Subject.Reference)
		medication.PractitionerFHIRID = fhir.SplitReferenceID(resource.Requester.Agent.Reference)
		if parsedTime, ok := fhir.ParseRFC3339(resource.AuthoredOn); ok {
			medication.IssuedAt = parsedTime
		}
		medications = append(medications, medication)
	}
	return medications, nil
}
