package tests

import (
	"context"
	"errors"
	"testing"

	"github.com/healthcare/backend/internal/modules/observation"
	"github.com/healthcare/backend/internal/modules/observation/mocks"
	"github.com/stretchr/testify/assert"
)

var errRepositoryFailure = errors.New("repository failure")

func TestCreateObservation_Success(t *testing.T) {
	observationService := observation.NewService(&mocks.MockObservationRepository{})

	entity := &observation.Observation{
		EncounterFHIRID: "encounter-123",
		PatientFHIRID:   "patient-456",
		LoincCode:       "55284-4",
		CodeDisplay:     "Blood pressure",
		ValueQuantity:   120,
		ValueUnit:       "mmHg",
	}

	result, err := observationService.CreateObservation(context.Background(), entity)

	assert.NoError(t, err)
	assert.Equal(t, "55284-4", result.LoincCode)
}

func TestCreateObservation_MissingEncounterID_ReturnsError(t *testing.T) {
	observationService := observation.NewService(&mocks.MockObservationRepository{})

	entity := &observation.Observation{EncounterFHIRID: "", PatientFHIRID: "patient-456", LoincCode: "55284-4"}

	result, err := observationService.CreateObservation(context.Background(), entity)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestGetObservationsByEncounter_RepositoryFailure_ReturnsError(t *testing.T) {
	observationService := observation.NewService(&mocks.MockObservationRepository{
		GetObservationsByEncounterFn: func(_ context.Context, _ string) ([]*observation.Observation, error) {
			return nil, errRepositoryFailure
		},
	})

	result, err := observationService.GetObservationsByEncounter(context.Background(), "encounter-123")

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestGetObservationsByPatient_Success(t *testing.T) {
	mockObservations := []*observation.Observation{
		{
			PatientFHIRID:   "patient-fhir-999",
			EncounterFHIRID: "encounter-fhir-888",
			LoincCode:       "8867-4",
			CodeDisplay:     "Heart Rate",
			ValueQuantity:   75,
			ValueUnit:       "bpm",
		},
	}

	observationService := observation.NewService(&mocks.MockObservationRepository{
		GetObservationsByPatientFn: func(ctx context.Context, patientFHIRID string) ([]*observation.Observation, error) {
			return mockObservations, nil
		},
	})

	retrievedObservations, observationError := observationService.GetObservationsByPatient(context.Background(), "patient-fhir-999")

	assert.NoError(t, observationError)
	assert.NotNil(t, retrievedObservations)
	assert.Len(t, retrievedObservations, 1)
	assert.Equal(t, "8867-4", retrievedObservations[0].LoincCode)
}

func TestGetObservationsByPatient_RepositoryFailure_ReturnsError(t *testing.T) {
	observationService := observation.NewService(&mocks.MockObservationRepository{
		GetObservationsByPatientFn: func(ctx context.Context, patientFHIRID string) ([]*observation.Observation, error) {
			return nil, errRepositoryFailure
		},
	})

	retrievedObservations, observationError := observationService.GetObservationsByPatient(context.Background(), "patient-fhir-999")

	assert.Error(t, observationError)
	assert.Nil(t, retrievedObservations)
	assert.True(t, errors.Is(observationError, errRepositoryFailure))
}
