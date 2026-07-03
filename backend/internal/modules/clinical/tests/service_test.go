package tests

import (
	"context"
	"errors"
	"testing"

	"github.com/healthcare/backend/internal/modules/clinical"
	"github.com/healthcare/backend/internal/modules/clinical/mocks"
	"github.com/stretchr/testify/assert"
)

var errRepositoryFailure = errors.New("repository failure")

func TestCreateEncounter_Success(t *testing.T) {
	svc := clinical.NewService(&mocks.MockClinicalRepository{})

	encounter := &clinical.Encounter{
		PatientFHIRID:  "patient-fhir-123",
		PractitionerID: "practitioner-456",
		ReasonCode:     "Z00.0",
		ReasonDisplay:  "Routine check-up",
	}

	result, err := svc.CreateEncounter(context.Background(), encounter)

	assert.NoError(t, err)
	assert.Equal(t, "patient-fhir-123", result.PatientFHIRID)
}

func TestCreateEncounter_MissingPatientFHIRID_ReturnsError(t *testing.T) {
	svc := clinical.NewService(&mocks.MockClinicalRepository{})

	encounter := &clinical.Encounter{PatientFHIRID: "", PractitionerID: "practitioner-456"}

	result, err := svc.CreateEncounter(context.Background(), encounter)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestCreateObservation_Success(t *testing.T) {
	svc := clinical.NewService(&mocks.MockClinicalRepository{})

	observation := &clinical.Observation{
		EncounterFHIRID: "encounter-123",
		PatientFHIRID:   "patient-456",
		LoincCode:       "55284-4",
		CodeDisplay:     "Blood pressure",
		ValueQuantity:   120,
		ValueUnit:       "mmHg",
	}

	result, err := svc.CreateObservation(context.Background(), observation)

	assert.NoError(t, err)
	assert.Equal(t, "55284-4", result.LoincCode)
}

func TestCreateObservation_MissingEncounterID_ReturnsError(t *testing.T) {
	svc := clinical.NewService(&mocks.MockClinicalRepository{})

	observation := &clinical.Observation{EncounterFHIRID: "", PatientFHIRID: "patient-456", LoincCode: "55284-4"}

	result, err := svc.CreateObservation(context.Background(), observation)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestCreateCondition_DefaultsClinicalStatusToActive(t *testing.T) {
	svc := clinical.NewService(&mocks.MockClinicalRepository{})

	condition := &clinical.Condition{
		PatientFHIRID:   "patient-123",
		EncounterFHIRID: "encounter-456",
		ICD10Code:       "I10",
		CodeDisplay:     "Essential hypertension",
		ClinicalStatus:  "",
	}

	result, err := svc.CreateCondition(context.Background(), condition)

	assert.NoError(t, err)
	assert.Equal(t, "active", result.ClinicalStatus)
}

func TestCreateCondition_MissingICD10Code_ReturnsError(t *testing.T) {
	svc := clinical.NewService(&mocks.MockClinicalRepository{})

	condition := &clinical.Condition{PatientFHIRID: "patient-123", ICD10Code: ""}

	result, err := svc.CreateCondition(context.Background(), condition)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestCreateAllergyIntolerance_DefaultsToActive(t *testing.T) {
	svc := clinical.NewService(&mocks.MockClinicalRepository{})

	allergy := &clinical.AllergyIntolerance{
		PatientFHIRID:   "patient-123",
		AllergenCode:    "7980",
		AllergenDisplay: "Penicillin",
		Reaction:        "Anaphylaxis",
	}

	result, err := svc.CreateAllergyIntolerance(context.Background(), allergy)

	assert.NoError(t, err)
	assert.Equal(t, "active", result.ClinicalStatus)
}

func TestCreateMedicationRequest_DefaultsToActive(t *testing.T) {
	svc := clinical.NewService(&mocks.MockClinicalRepository{})

	medication := &clinical.MedicationRequest{
		PatientFHIRID:      "patient-123",
		EncounterFHIRID:    "encounter-456",
		MedicationCode:     "10582",
		MedicationName:     "Amoxicillin",
		DosageInstructions: "500mg every 8 hours",
	}

	result, err := svc.CreateMedicationRequest(context.Background(), medication)

	assert.NoError(t, err)
	assert.Equal(t, "active", result.Status)
}

func TestCreateDiagnosticReport_DefaultsToFinal(t *testing.T) {
	svc := clinical.NewService(&mocks.MockClinicalRepository{})

	report := &clinical.DiagnosticReport{
		PatientFHIRID:   "patient-123",
		EncounterFHIRID: "encounter-456",
		ReportCode:      "24323-8",
		ReportDisplay:   "Complete blood count",
		Conclusion:      "Normal values",
	}

	result, err := svc.CreateDiagnosticReport(context.Background(), report)

	assert.NoError(t, err)
	assert.Equal(t, "final", result.Status)
}

func TestGetObservationsByEncounter_RepositoryFailure_ReturnsError(t *testing.T) {
	svc := clinical.NewService(&mocks.MockClinicalRepository{
		GetObservationsByEncounterFn: func(_ context.Context, _ string) ([]*clinical.Observation, error) {
			return nil, errRepositoryFailure
		},
	})

	result, err := svc.GetObservationsByEncounter(context.Background(), "encounter-123")

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestGetObservationsByPatient_Success(t *testing.T) {
	mockClinicalObservations := []*clinical.Observation{
		{
			PatientFHIRID:   "patient-fhir-999",
			EncounterFHIRID: "encounter-fhir-888",
			LoincCode:       "8867-4",
			CodeDisplay:     "Heart Rate",
			ValueQuantity:   75,
			ValueUnit:       "bpm",
		},
	}

	clinicalService := clinical.NewService(&mocks.MockClinicalRepository{
		GetObservationsByPatientFn: func(ctx context.Context, patientFHIRID string) ([]*clinical.Observation, error) {
			return mockClinicalObservations, nil
		},
	})

	retrievedObservations, observationError := clinicalService.GetObservationsByPatient(context.Background(), "patient-fhir-999")

	assert.NoError(t, observationError)
	assert.NotNil(t, retrievedObservations)
	assert.Len(t, retrievedObservations, 1)
	assert.Equal(t, "8867-4", retrievedObservations[0].LoincCode)
}

func TestGetObservationsByPatient_RepositoryFailure_ReturnsError(t *testing.T) {
	clinicalService := clinical.NewService(&mocks.MockClinicalRepository{
		GetObservationsByPatientFn: func(ctx context.Context, patientFHIRID string) ([]*clinical.Observation, error) {
			return nil, errRepositoryFailure
		},
	})

	retrievedObservations, observationError := clinicalService.GetObservationsByPatient(context.Background(), "patient-fhir-999")

	assert.Error(t, observationError)
	assert.Nil(t, retrievedObservations)
	assert.True(t, errors.Is(observationError, errRepositoryFailure))
}
