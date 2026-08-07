package allergy

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

func parseAllergyBundle(responseBody json.RawMessage) ([]*Allergy, error) {
	decodedResources, err := fhir.DecodeBundle[fhir.AllergyIntolerance](responseBody)
	if err != nil {
		return nil, err
	}
	allergies := make([]*Allergy, 0, len(decodedResources))
	for _, resource := range decodedResources {
		allergy := &Allergy{}
		allergy.FHIRResourceID = resource.ID
		clinicalStatusCode, _, _ := fhir.CodeableConceptParts(resource.ClinicalStatus)
		allergy.ClinicalStatus = clinicalStatusCode
		code, display, text := fhir.CodeableConceptParts(resource.Code)
		allergy.AllergenDisplay = text
		allergy.AllergenCode = code
		if allergy.AllergenDisplay == "" {
			allergy.AllergenDisplay = display
		}
		allergy.PatientFHIRID = fhir.SplitReferenceID(resource.Patient.Reference)
		if parsedTime, ok := fhir.ParseRFC3339(resource.RecordedDate); ok {
			allergy.RecordedAt = parsedTime
		}
		if len(resource.Reaction) > 0 && len(resource.Reaction[0].Manifestation) > 0 {
			_, _, text := fhir.CodeableConceptParts(resource.Reaction[0].Manifestation[0])
			allergy.Reaction = text
		}
		allergies = append(allergies, allergy)
	}
	return allergies, nil
}
