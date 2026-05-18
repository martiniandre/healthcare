package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/healthcare/backend/internal/modules/patients"
)

type MockPatientRepository struct {
	Patients map[uuid.UUID]*patients.Patient
	ByDoc    map[string]*patients.Patient
	Err      error
}

func NewMockPatientRepository() *MockPatientRepository {
	return &MockPatientRepository{
		Patients: make(map[uuid.UUID]*patients.Patient),
		ByDoc:    make(map[string]*patients.Patient),
	}
}

func (mockRepo *MockPatientRepository) CreatePatient(contextParam context.Context, patient *patients.Patient) error {
	if mockRepo.Err != nil {
		return mockRepo.Err
	}
	mockRepo.Patients[patient.ID] = patient
	mockRepo.ByDoc[patient.DocumentID] = patient
	return nil
}

func (mockRepo *MockPatientRepository) GetPatientByID(contextParam context.Context, patientID uuid.UUID) (*patients.Patient, error) {
	if mockRepo.Err != nil {
		return nil, mockRepo.Err
	}
	patient, exists := mockRepo.Patients[patientID]
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

func (mockRepo *MockPatientRepository) ListPatients(contextParam context.Context) ([]*patients.Patient, error) {
	if mockRepo.Err != nil {
		return nil, mockRepo.Err
	}
	result := make([]*patients.Patient, 0, len(mockRepo.Patients))
	for _, patient := range mockRepo.Patients {
		result = append(result, patient)
	}
	return result, nil
}
