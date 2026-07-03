package mocks

import (
	"context"

	"github.com/healthcare/backend/internal/modules/clinical"
)

type MockClinicalRepository struct {
	CreateEncounterFn            func(ctx context.Context, encounter *clinical.Encounter) (*clinical.Encounter, error)
	GetEncountersByPatientFn     func(ctx context.Context, patientFHIRID string) ([]*clinical.Encounter, error)
	CreateObservationFn          func(ctx context.Context, observation *clinical.Observation) (*clinical.Observation, error)
	GetObservationsByEncounterFn func(ctx context.Context, encounterFHIRID string) ([]*clinical.Observation, error)
	GetObservationsByPatientFn   func(ctx context.Context, patientFHIRID string) ([]*clinical.Observation, error)
	CreateConditionFn            func(ctx context.Context, condition *clinical.Condition) (*clinical.Condition, error)
	GetConditionsByPatientFn     func(ctx context.Context, patientFHIRID string) ([]*clinical.Condition, error)
	CreateAllergyIntoleranceFn          func(ctx context.Context, allergy *clinical.AllergyIntolerance) (*clinical.AllergyIntolerance, error)
	GetAllergyIntolerancesByPatientFn    func(ctx context.Context, patientFHIRID string) ([]*clinical.AllergyIntolerance, error)
	CreateMedicationRequestFn           func(ctx context.Context, medicationRequest *clinical.MedicationRequest) (*clinical.MedicationRequest, error)
	GetMedicationRequestsByEncounterFn  func(ctx context.Context, encounterFHIRID string) ([]*clinical.MedicationRequest, error)
	CreateDiagnosticReportFn            func(ctx context.Context, report *clinical.DiagnosticReport) (*clinical.DiagnosticReport, error)
	GetDiagnosticReportsByEncounterFn   func(ctx context.Context, encounterFHIRID string) ([]*clinical.DiagnosticReport, error)
}

func (mockRepo *MockClinicalRepository) CreateEncounter(ctx context.Context, encounter *clinical.Encounter) (*clinical.Encounter, error) {
	if mockRepo.CreateEncounterFn != nil {
		return mockRepo.CreateEncounterFn(ctx, encounter)
	}
	return encounter, nil
}

func (mockRepo *MockClinicalRepository) GetEncountersByPatient(ctx context.Context, patientFHIRID string) ([]*clinical.Encounter, error) {
	if mockRepo.GetEncountersByPatientFn != nil {
		return mockRepo.GetEncountersByPatientFn(ctx, patientFHIRID)
	}
	return []*clinical.Encounter{}, nil
}

func (mockRepo *MockClinicalRepository) CreateObservation(ctx context.Context, observation *clinical.Observation) (*clinical.Observation, error) {
	if mockRepo.CreateObservationFn != nil {
		return mockRepo.CreateObservationFn(ctx, observation)
	}
	return observation, nil
}

func (mockRepo *MockClinicalRepository) GetObservationsByEncounter(ctx context.Context, encounterFHIRID string) ([]*clinical.Observation, error) {
	if mockRepo.GetObservationsByEncounterFn != nil {
		return mockRepo.GetObservationsByEncounterFn(ctx, encounterFHIRID)
	}
	return []*clinical.Observation{}, nil
}

func (mockRepo *MockClinicalRepository) GetObservationsByPatient(ctx context.Context, patientFHIRID string) ([]*clinical.Observation, error) {
	if mockRepo.GetObservationsByPatientFn != nil {
		return mockRepo.GetObservationsByPatientFn(ctx, patientFHIRID)
	}
	return []*clinical.Observation{}, nil
}

func (mockRepo *MockClinicalRepository) CreateCondition(ctx context.Context, condition *clinical.Condition) (*clinical.Condition, error) {
	if mockRepo.CreateConditionFn != nil {
		return mockRepo.CreateConditionFn(ctx, condition)
	}
	return condition, nil
}

func (mockRepo *MockClinicalRepository) GetConditionsByPatient(ctx context.Context, patientFHIRID string) ([]*clinical.Condition, error) {
	if mockRepo.GetConditionsByPatientFn != nil {
		return mockRepo.GetConditionsByPatientFn(ctx, patientFHIRID)
	}
	return []*clinical.Condition{}, nil
}

func (mockRepo *MockClinicalRepository) CreateAllergyIntolerance(ctx context.Context, allergy *clinical.AllergyIntolerance) (*clinical.AllergyIntolerance, error) {
	if mockRepo.CreateAllergyIntoleranceFn != nil {
		return mockRepo.CreateAllergyIntoleranceFn(ctx, allergy)
	}
	return allergy, nil
}

func (mockRepo *MockClinicalRepository) GetAllergyIntolerancesByPatient(ctx context.Context, patientFHIRID string) ([]*clinical.AllergyIntolerance, error) {
	if mockRepo.GetAllergyIntolerancesByPatientFn != nil {
		return mockRepo.GetAllergyIntolerancesByPatientFn(ctx, patientFHIRID)
	}
	return []*clinical.AllergyIntolerance{}, nil
}

func (mockRepo *MockClinicalRepository) CreateMedicationRequest(ctx context.Context, medicationRequest *clinical.MedicationRequest) (*clinical.MedicationRequest, error) {
	if mockRepo.CreateMedicationRequestFn != nil {
		return mockRepo.CreateMedicationRequestFn(ctx, medicationRequest)
	}
	return medicationRequest, nil
}

func (mockRepo *MockClinicalRepository) GetMedicationRequestsByEncounter(ctx context.Context, encounterFHIRID string) ([]*clinical.MedicationRequest, error) {
	if mockRepo.GetMedicationRequestsByEncounterFn != nil {
		return mockRepo.GetMedicationRequestsByEncounterFn(ctx, encounterFHIRID)
	}
	return []*clinical.MedicationRequest{}, nil
}

func (mockRepo *MockClinicalRepository) CreateDiagnosticReport(ctx context.Context, report *clinical.DiagnosticReport) (*clinical.DiagnosticReport, error) {
	if mockRepo.CreateDiagnosticReportFn != nil {
		return mockRepo.CreateDiagnosticReportFn(ctx, report)
	}
	return report, nil
}

func (mockRepo *MockClinicalRepository) GetDiagnosticReportsByEncounter(ctx context.Context, encounterFHIRID string) ([]*clinical.DiagnosticReport, error) {
	if mockRepo.GetDiagnosticReportsByEncounterFn != nil {
		return mockRepo.GetDiagnosticReportsByEncounterFn(ctx, encounterFHIRID)
	}
	return []*clinical.DiagnosticReport{}, nil
}
