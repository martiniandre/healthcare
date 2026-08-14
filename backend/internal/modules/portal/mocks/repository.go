package mocks

import (
	"context"

	"github.com/healthcare/backend/internal/modules/portal"
)

type MockPortalRepository struct {
	Patients        map[string]*portal.PatientInfo
	Encounters      map[string][]portal.PortalEncounter
	Observations    map[string][]portal.PortalObservation
	Conditions      map[string][]portal.PortalCondition
	Medications     map[string][]portal.PortalMedication
	Reports         map[string][]portal.PortalReport
	Imaging         map[string][]portal.PortalImaging
	PatientErr      error
	EncountersErr   error
	ObservationsErr error
	ConditionsErr   error
	MedicationsErr  error
	ReportsErr      error
	ImagingErr      error
}

func NewMockPortalRepository() *MockPortalRepository {
	return &MockPortalRepository{
		Patients:     make(map[string]*portal.PatientInfo),
		Encounters:   make(map[string][]portal.PortalEncounter),
		Observations: make(map[string][]portal.PortalObservation),
		Conditions:   make(map[string][]portal.PortalCondition),
		Medications:  make(map[string][]portal.PortalMedication),
		Reports:      make(map[string][]portal.PortalReport),
		Imaging:      make(map[string][]portal.PortalImaging),
	}
}

func (mockRepository *MockPortalRepository) GetPatient(ctx context.Context, fhirResourceID string) (*portal.PatientInfo, error) {
	if mockRepository.PatientErr != nil {
		return nil, mockRepository.PatientErr
	}
	patientInfo, exists := mockRepository.Patients[fhirResourceID]
	if !exists {
		return nil, portal.ErrPatientNotFound
	}
	return patientInfo, nil
}

func (mockRepository *MockPortalRepository) GetEncountersByPatient(ctx context.Context, patientFHIRID string) ([]portal.PortalEncounter, error) {
	if mockRepository.EncountersErr != nil {
		return nil, mockRepository.EncountersErr
	}
	return mockRepository.Encounters[patientFHIRID], nil
}

func (mockRepository *MockPortalRepository) GetObservationsByPatient(ctx context.Context, patientFHIRID string) ([]portal.PortalObservation, error) {
	if mockRepository.ObservationsErr != nil {
		return nil, mockRepository.ObservationsErr
	}
	return mockRepository.Observations[patientFHIRID], nil
}

func (mockRepository *MockPortalRepository) GetConditionsByPatient(ctx context.Context, patientFHIRID string) ([]portal.PortalCondition, error) {
	if mockRepository.ConditionsErr != nil {
		return nil, mockRepository.ConditionsErr
	}
	return mockRepository.Conditions[patientFHIRID], nil
}

func (mockRepository *MockPortalRepository) GetMedicationsByPatient(ctx context.Context, patientFHIRID string) ([]portal.PortalMedication, error) {
	if mockRepository.MedicationsErr != nil {
		return nil, mockRepository.MedicationsErr
	}
	return mockRepository.Medications[patientFHIRID], nil
}

func (mockRepository *MockPortalRepository) GetReportsByPatient(ctx context.Context, patientFHIRID string) ([]portal.PortalReport, error) {
	if mockRepository.ReportsErr != nil {
		return nil, mockRepository.ReportsErr
	}
	return mockRepository.Reports[patientFHIRID], nil
}

func (mockRepository *MockPortalRepository) GetImagingByPatient(ctx context.Context, patientFHIRID string) ([]portal.PortalImaging, error) {
	if mockRepository.ImagingErr != nil {
		return nil, mockRepository.ImagingErr
	}
	return mockRepository.Imaging[patientFHIRID], nil
}
