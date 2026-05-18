package patients

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/healthcare/backend/internal/shared/fhir"
	"github.com/healthcare/backend/internal/shared/healthcare"
)

type Repository interface {
	CreatePatient(ctx context.Context, patient *Patient) (*Patient, error)
	GetPatientByID(ctx context.Context, fhirResourceID string) (*Patient, error)
	GetPatientByDocumentID(ctx context.Context, documentID string) (*Patient, error)
	ListPatients(ctx context.Context) ([]*Patient, error)
}

type repository struct {
	fhirClient healthcare.FHIRClient
}

func NewRepository(fhirClient healthcare.FHIRClient) Repository {
	return &repository{fhirClient: fhirClient}
}

func (patientRepository *repository) CreatePatient(ctx context.Context, patient *Patient) (*Patient, error) {
	fhirPatient := fhir.NewPatientResource(
		patient.FullName,
		patient.DocumentID,
		patient.PhoneNumber,
		patient.BirthDate.Format("2006-01-02"),
	)

	responseBody, err := patientRepository.fhirClient.CreateResource(ctx, "Patient", fhirPatient)
	if err != nil {
		return nil, fmt.Errorf("failed to create patient in healthcare api: %w", err)
	}

	var createdResource map[string]interface{}
	if err := json.Unmarshal(responseBody, &createdResource); err != nil {
		return nil, fmt.Errorf("failed to parse healthcare api response: %w", err)
	}

	fhirID, ok := createdResource["id"].(string)
	if !ok {
		return nil, fmt.Errorf("healthcare api did not return a valid resource id")
	}

	patient.FHIRResourceID = fhirID
	patient.ID = uuid.New()
	return patient, nil
}

func (patientRepository *repository) GetPatientByID(ctx context.Context, fhirResourceID string) (*Patient, error) {
	responseBody, err := patientRepository.fhirClient.GetResource(ctx, "Patient", fhirResourceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get patient from healthcare api: %w", err)
	}

	return parsePatientFromFHIR(responseBody, fhirResourceID)
}

func (patientRepository *repository) GetPatientByDocumentID(ctx context.Context, documentID string) (*Patient, error) {
	queryParams := fmt.Sprintf("identifier=%s", documentID)
	responseBody, err := patientRepository.fhirClient.SearchResources(ctx, "Patient", queryParams)
	if err != nil {
		return nil, fmt.Errorf("failed to search patient by document in healthcare api: %w", err)
	}

	var bundle map[string]interface{}
	if err := json.Unmarshal(responseBody, &bundle); err != nil {
		return nil, fmt.Errorf("failed to parse search response: %w", err)
	}

	entries, ok := bundle["entry"].([]interface{})
	if !ok || len(entries) == 0 {
		return nil, ErrPatientNotFound
	}

	firstEntry, ok := entries[0].(map[string]interface{})
	if !ok {
		return nil, ErrPatientNotFound
	}

	resource, ok := firstEntry["resource"].(map[string]interface{})
	if !ok {
		return nil, ErrPatientNotFound
	}

	fhirID, _ := resource["id"].(string)
	entryBytes, err := json.Marshal(resource)
	if err != nil {
		return nil, err
	}

	return parsePatientFromFHIR(entryBytes, fhirID)
}

func (patientRepository *repository) ListPatients(ctx context.Context) ([]*Patient, error) {
	responseBody, err := patientRepository.fhirClient.SearchResources(ctx, "Patient", "_count=100")
	if err != nil {
		return nil, fmt.Errorf("failed to list patients from healthcare api: %w", err)
	}

	var bundle map[string]interface{}
	if err := json.Unmarshal(responseBody, &bundle); err != nil {
		return nil, fmt.Errorf("failed to parse list response: %w", err)
	}

	entries, ok := bundle["entry"].([]interface{})
	if !ok {
		return []*Patient{}, nil
	}

	patientList := make([]*Patient, 0, len(entries))
	for _, entry := range entries {
		entryMap, ok := entry.(map[string]interface{})
		if !ok {
			continue
		}
		resource, ok := entryMap["resource"].(map[string]interface{})
		if !ok {
			continue
		}
		fhirID, _ := resource["id"].(string)
		entryBytes, err := json.Marshal(resource)
		if err != nil {
			continue
		}
		patient, err := parsePatientFromFHIR(entryBytes, fhirID)
		if err != nil {
			continue
		}
		patientList = append(patientList, patient)
	}

	return patientList, nil
}

func parsePatientFromFHIR(responseBody json.RawMessage, fhirResourceID string) (*Patient, error) {
	var fhirResource map[string]interface{}
	if err := json.Unmarshal(responseBody, &fhirResource); err != nil {
		return nil, fmt.Errorf("failed to parse fhir resource: %w", err)
	}

	patient := &Patient{
		ID:             uuid.New(),
		FHIRResourceID: fhirResourceID,
		IsActive:       true,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if names, ok := fhirResource["name"].([]interface{}); ok && len(names) > 0 {
		if firstNameMap, ok := names[0].(map[string]interface{}); ok {
			if family, ok := firstNameMap["family"].(string); ok {
				patient.FullName = family
			}
		}
	}

	if birthDateStr, ok := fhirResource["birthDate"].(string); ok {
		parsedBirthDate, err := time.Parse("2006-01-02", birthDateStr)
		if err == nil {
			patient.BirthDate = parsedBirthDate
		}
	}

	if identifiers, ok := fhirResource["identifier"].([]interface{}); ok {
		for _, identifier := range identifiers {
			if identifierMap, ok := identifier.(map[string]interface{}); ok {
				if value, ok := identifierMap["value"].(string); ok {
					patient.DocumentID = value
					break
				}
			}
		}
	}

	if telecom, ok := fhirResource["telecom"].([]interface{}); ok {
		for _, contact := range telecom {
			if contactMap, ok := contact.(map[string]interface{}); ok {
				if system, ok := contactMap["system"].(string); ok && system == "phone" {
					if value, ok := contactMap["value"].(string); ok {
						patient.PhoneNumber = value
					}
				}
			}
		}
	}

	return patient, nil
}
