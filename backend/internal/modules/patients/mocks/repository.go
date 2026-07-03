package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/healthcare/backend/internal/modules/patients"
)

type MockPatientRepository struct {
	Patients map[string]*patients.Patient
	ByDoc    map[string]*patients.Patient
	Err      error
}

func NewMockPatientRepository() *MockPatientRepository {
	return &MockPatientRepository{
		Patients: make(map[string]*patients.Patient),
		ByDoc:    make(map[string]*patients.Patient),
	}
}

func (mockRepo *MockPatientRepository) CreatePatient(contextParam context.Context, patient *patients.Patient) (*patients.Patient, error) {
	if mockRepo.Err != nil {
		return nil, mockRepo.Err
	}
	if patient.FHIRResourceID == "" {
		patient.FHIRResourceID = uuid.New().String()
	}
	mockRepo.Patients[patient.FHIRResourceID] = patient
	mockRepo.ByDoc[patient.DocumentID] = patient
	return patient, nil
}

func (mockRepo *MockPatientRepository) GetPatientByID(contextParam context.Context, fhirResourceID string) (*patients.Patient, error) {
	if mockRepo.Err != nil {
		return nil, mockRepo.Err
	}
	patient, exists := mockRepo.Patients[fhirResourceID]
	if !exists {
		return nil, patients.ErrPatientNotFound
	}
	return patient, nil
}

func (mockRepo *MockPatientRepository) GetPatientByDocumentID(contextParam context.Context, documentID string) (*patients.Patient, error) {
	if mockRepo.Err != nil {
		return nil, mockRepo.Err
	}
	patient, exists := mockRepo.ByDoc[documentID]
	if !exists {
		return nil, patients.ErrPatientNotFound
	}
	return patient, nil
}

func (mockRepo *MockPatientRepository) ListPatients(contextParam context.Context, search string, sortField string, sortDirection string, page int, limit int) ([]*patients.Patient, error) {
	if mockRepo.Err != nil {
		return nil, mockRepo.Err
	}
	result := make([]*patients.Patient, 0, len(mockRepo.Patients))
	for _, patient := range mockRepo.Patients {
		result = append(result, patient)
	}
	return result, nil
}
