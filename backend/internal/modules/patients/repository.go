package patients

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/healthcare/backend/internal/shared/fhir"
	"github.com/healthcare/backend/internal/shared/healthcare"
)

const documentIdentifierSystem = "urn:oid:2.16.840.1.113883.13.237"

type Repository interface {
	CreatePatient(ctx context.Context, patient *Patient) (*Patient, error)
	GetPatientByID(ctx context.Context, fhirResourceID string) (*Patient, error)
	GetPatientByDocumentID(ctx context.Context, documentID string) (*Patient, error)
	ListPatients(ctx context.Context, search string, sortField string, sortDirection string, page int, limit int) ([]*Patient, error)
}

type repository struct {
	fhirClient healthcare.FHIRClient
}

func NewRepository(fhirClient healthcare.FHIRClient) Repository {
	return &repository{fhirClient: fhirClient}
}

func (patientRepository *repository) CreatePatient(ctx context.Context, patient *Patient) (*Patient, error) {
	nameParts := strings.SplitN(patient.FullName, " ", 2)
	givenName := nameParts[0]
	familyName := ""
	if len(nameParts) > 1 {
		familyName = nameParts[1]
	}

	fhirPatient := fhir.NewPatientResource(
		givenName,
		familyName,
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
	queryParams := url.Values{"identifier": []string{documentIdentifierSystem + "|" + documentID}}.Encode()
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

func (patientRepository *repository) ListPatients(ctx context.Context, search string, sortField string, sortDirection string, page int, limit int) ([]*Patient, error) {
	v := url.Values{}
	
	if limit <= 0 {
		limit = 100
	}
	v.Add("_count", fmt.Sprintf("%d", limit))
	
	// Note: Healthcare API might not support _offset directly for all resources, but we try standard FHIR
	if page > 1 {
		offset := (page - 1) * limit
		v.Add("_offset", fmt.Sprintf("%d", offset))
	}
	
	if search != "" {
		v.Add("name:contains", search)
	}
	
	if sortField != "" {
		fhirField := ""
		switch sortField {
		case "full_name":
			fhirField = "name"
		case "birth_date":
			fhirField = "birthdate"
		case "document_id":
			fhirField = "identifier"
		}
		if fhirField != "" {
			if sortDirection == "desc" {
				fhirField = "-" + fhirField
			}
			v.Add("_sort", fhirField)
		}
	}

	queryParams := v.Encode()
	responseBody, err := patientRepository.fhirClient.SearchResources(ctx, "Patient", queryParams)
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

	patientID := uuid.NewSHA1(uuid.NameSpaceDNS, []byte(fhirResourceID))

	patient := &Patient{
		ID:             patientID,
		FHIRResourceID: fhirResourceID,
		IsActive:       true,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if meta, hasMeta := fhirResource["meta"].(map[string]interface{}); hasMeta {
		if lastUpdated, hasLastUpdated := meta["lastUpdated"].(string); hasLastUpdated {
			parsedTime, parseErr := time.Parse(time.RFC3339, lastUpdated)
			if parseErr == nil {
				patient.UpdatedAt = parsedTime
				patient.CreatedAt = parsedTime
			}
		}
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
