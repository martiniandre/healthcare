package tests

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/healthcare/backend/internal/modules/observation"
	"github.com/healthcare/backend/internal/modules/observation/mocks"
	"github.com/healthcare/backend/internal/shared/apperrors"
	"github.com/stretchr/testify/assert"
)

var errRepositoryFailure = errors.New("repository failure")

func TestCreateObservation_ValidInputWithoutObservedAt_SetsCurrentTime(testingInstance *testing.T) {
	observationService := observation.NewService(&mocks.MockObservationRepository{})

	createdObservation, createError := observationService.CreateObservation(context.Background(), observation.CreateObservationInput{
		EncounterFHIRID: "encounter-123",
		PatientFHIRID:   "patient-456",
		LoincCode:       "55284-4",
		CodeDisplay:     "Blood pressure",
		ValueQuantity:   120,
		ValueUnit:       "mmHg",
	})

	assert.NoError(testingInstance, createError)
	assert.NotNil(testingInstance, createdObservation)
	assert.Equal(testingInstance, "encounter-123", createdObservation.EncounterFHIRID)
	assert.Equal(testingInstance, "patient-456", createdObservation.PatientFHIRID)
	assert.Equal(testingInstance, "55284-4", createdObservation.LoincCode)
	assert.Equal(testingInstance, 120.0, createdObservation.ValueQuantity)
	assert.Equal(testingInstance, "mmHg", createdObservation.ValueUnit)
	assert.False(testingInstance, createdObservation.ObservedAt.IsZero())
	assert.WithinDuration(testingInstance, time.Now(), createdObservation.ObservedAt, time.Second)
}

func TestCreateObservation_ValidInputWithExplicitObservedAt_PreservesProvidedTime(testingInstance *testing.T) {
	observationService := observation.NewService(&mocks.MockObservationRepository{})
	explicitObservedAt := time.Date(2024, time.May, 10, 9, 30, 0, 0, time.UTC)

	createdObservation, createError := observationService.CreateObservation(context.Background(), observation.CreateObservationInput{
		EncounterFHIRID: "encounter-123",
		PatientFHIRID:   "patient-456",
		LoincCode:       "55284-4",
		ValueQuantity:   120,
		ValueUnit:       "mmHg",
		ObservedAt:      &explicitObservedAt,
	})

	assert.NoError(testingInstance, createError)
	assert.NotNil(testingInstance, createdObservation)
	assert.Equal(testingInstance, explicitObservedAt, createdObservation.ObservedAt)
}

func TestCreateObservation_MissingEncounterFHIRID_ReturnsInvalidArgumentAndSkipsRepository(testingInstance *testing.T) {
	repositoryCalled := false
	observationService := observation.NewService(&mocks.MockObservationRepository{
		CreateObservationFn: func(ctx context.Context, entity *observation.Observation) (*observation.Observation, error) {
			repositoryCalled = true
			return entity, nil
		},
	})

	createdObservation, createError := observationService.CreateObservation(context.Background(), observation.CreateObservationInput{
		PatientFHIRID: "patient-456",
		LoincCode:     "55284-4",
	})

	assert.Error(testingInstance, createError)
	assert.Nil(testingInstance, createdObservation)
	assert.False(testingInstance, repositoryCalled)
	var appError apperrors.AppError
	assert.ErrorAs(testingInstance, createError, &appError)
}

func TestCreateObservation_MissingPatientFHIRID_ReturnsInvalidArgumentAndSkipsRepository(testingInstance *testing.T) {
	repositoryCalled := false
	observationService := observation.NewService(&mocks.MockObservationRepository{
		CreateObservationFn: func(ctx context.Context, entity *observation.Observation) (*observation.Observation, error) {
			repositoryCalled = true
			return entity, nil
		},
	})

	createdObservation, createError := observationService.CreateObservation(context.Background(), observation.CreateObservationInput{
		EncounterFHIRID: "encounter-123",
		LoincCode:       "55284-4",
	})

	assert.Error(testingInstance, createError)
	assert.Nil(testingInstance, createdObservation)
	assert.False(testingInstance, repositoryCalled)
}

func TestCreateObservation_MissingLoincCode_ReturnsInvalidArgumentAndSkipsRepository(testingInstance *testing.T) {
	repositoryCalled := false
	observationService := observation.NewService(&mocks.MockObservationRepository{
		CreateObservationFn: func(ctx context.Context, entity *observation.Observation) (*observation.Observation, error) {
			repositoryCalled = true
			return entity, nil
		},
	})

	createdObservation, createError := observationService.CreateObservation(context.Background(), observation.CreateObservationInput{
		EncounterFHIRID: "encounter-123",
		PatientFHIRID:   "patient-456",
	})

	assert.Error(testingInstance, createError)
	assert.Nil(testingInstance, createdObservation)
	assert.False(testingInstance, repositoryCalled)
}

func TestCreateObservation_InvalidLoincFormat_ReturnsInvalidArgumentAndSkipsRepository(testingInstance *testing.T) {
	repositoryCalled := false
	observationService := observation.NewService(&mocks.MockObservationRepository{
		CreateObservationFn: func(ctx context.Context, entity *observation.Observation) (*observation.Observation, error) {
			repositoryCalled = true
			return entity, nil
		},
	})

	createdObservation, createError := observationService.CreateObservation(context.Background(), observation.CreateObservationInput{
		EncounterFHIRID: "encounter-123",
		PatientFHIRID:   "patient-456",
		LoincCode:       "88674",
		ValueQuantity:   120,
	})

	assert.Error(testingInstance, createError)
	assert.Nil(testingInstance, createdObservation)
	assert.False(testingInstance, repositoryCalled)
	var appError apperrors.AppError
	assert.ErrorAs(testingInstance, createError, &appError)
}

func TestCreateObservation_ValueOutOfClinicalRange_ReturnsInvalidArgumentAndSkipsRepository(testingInstance *testing.T) {
	repositoryCalled := false
	observationService := observation.NewService(&mocks.MockObservationRepository{
		CreateObservationFn: func(ctx context.Context, entity *observation.Observation) (*observation.Observation, error) {
			repositoryCalled = true
			return entity, nil
		},
	})

	createdObservation, createError := observationService.CreateObservation(context.Background(), observation.CreateObservationInput{
		EncounterFHIRID: "encounter-123",
		PatientFHIRID:   "patient-456",
		LoincCode:       "8867-4",
		ValueQuantity:   350,
	})

	assert.Error(testingInstance, createError)
	assert.Nil(testingInstance, createdObservation)
	assert.False(testingInstance, repositoryCalled)
	var appError apperrors.AppError
	assert.ErrorAs(testingInstance, createError, &appError)
}

func TestCreateObservation_RepositoryFailure_ReturnsError(testingInstance *testing.T) {
	observationService := observation.NewService(&mocks.MockObservationRepository{
		CreateObservationFn: func(ctx context.Context, entity *observation.Observation) (*observation.Observation, error) {
			return nil, errRepositoryFailure
		},
	})

	createdObservation, createError := observationService.CreateObservation(context.Background(), observation.CreateObservationInput{
		EncounterFHIRID: "encounter-123",
		PatientFHIRID:   "patient-456",
		LoincCode:       "55284-4",
		ValueQuantity:   120,
		ValueUnit:       "mmHg",
	})

	assert.Error(testingInstance, createError)
	assert.Nil(testingInstance, createdObservation)
	assert.True(testingInstance, errors.Is(createError, errRepositoryFailure))
}

func TestGetObservationsByEncounter_Success(testingInstance *testing.T) {
	mockObservations := []*observation.Observation{
		{
			FHIRResourceID:  "observation-1",
			EncounterFHIRID: "encounter-123",
			PatientFHIRID:   "patient-456",
			LoincCode:       "8867-4",
			CodeDisplay:     "Heart Rate",
			ValueQuantity:   75,
			ValueUnit:       "bpm",
		},
	}

	observationService := observation.NewService(&mocks.MockObservationRepository{
		GetObservationsByEncounterFn: func(ctx context.Context, encounterFHIRID string) ([]*observation.Observation, error) {
			return mockObservations, nil
		},
	})

	retrievedObservations, encounterError := observationService.GetObservationsByEncounter(context.Background(), "encounter-123")

	assert.NoError(testingInstance, encounterError)
	assert.NotNil(testingInstance, retrievedObservations)
	assert.Len(testingInstance, retrievedObservations, 1)
	assert.Equal(testingInstance, "8867-4", retrievedObservations[0].LoincCode)
}

func TestGetObservationsByEncounter_RepositoryFailure_ReturnsError(testingInstance *testing.T) {
	observationService := observation.NewService(&mocks.MockObservationRepository{
		GetObservationsByEncounterFn: func(ctx context.Context, encounterFHIRID string) ([]*observation.Observation, error) {
			return nil, errRepositoryFailure
		},
	})

	retrievedObservations, encounterError := observationService.GetObservationsByEncounter(context.Background(), "encounter-123")

	assert.Error(testingInstance, encounterError)
	assert.Nil(testingInstance, retrievedObservations)
	assert.True(testingInstance, errors.Is(encounterError, errRepositoryFailure))
}

func TestGetObservationsByEncounter_RepositoryNotFound_ReturnsError(testingInstance *testing.T) {
	observationService := observation.NewService(&mocks.MockObservationRepository{
		GetObservationsByEncounterFn: func(ctx context.Context, encounterFHIRID string) ([]*observation.Observation, error) {
			return nil, apperrors.ErrObservationNotFound
		},
	})

	retrievedObservations, encounterError := observationService.GetObservationsByEncounter(context.Background(), "encounter-123")

	assert.Error(testingInstance, encounterError)
	assert.Nil(testingInstance, retrievedObservations)
	assert.True(testingInstance, errors.Is(encounterError, apperrors.ErrObservationNotFound))
}

func TestGetObservationsByPatient_Success(testingInstance *testing.T) {
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

	assert.NoError(testingInstance, observationError)
	assert.NotNil(testingInstance, retrievedObservations)
	assert.Len(testingInstance, retrievedObservations, 1)
	assert.Equal(testingInstance, "8867-4", retrievedObservations[0].LoincCode)
}

func TestGetObservationsByPatient_RepositoryFailure_ReturnsError(testingInstance *testing.T) {
	observationService := observation.NewService(&mocks.MockObservationRepository{
		GetObservationsByPatientFn: func(ctx context.Context, patientFHIRID string) ([]*observation.Observation, error) {
			return nil, errRepositoryFailure
		},
	})

	retrievedObservations, observationError := observationService.GetObservationsByPatient(context.Background(), "patient-fhir-999")

	assert.Error(testingInstance, observationError)
	assert.Nil(testingInstance, retrievedObservations)
	assert.True(testingInstance, errors.Is(observationError, errRepositoryFailure))
}

func TestGetObservationsByPatient_RepositoryNotFound_ReturnsError(testingInstance *testing.T) {
	observationService := observation.NewService(&mocks.MockObservationRepository{
		GetObservationsByPatientFn: func(ctx context.Context, patientFHIRID string) ([]*observation.Observation, error) {
			return nil, apperrors.ErrObservationNotFound
		},
	})

	retrievedObservations, observationError := observationService.GetObservationsByPatient(context.Background(), "patient-fhir-999")

	assert.Error(testingInstance, observationError)
	assert.Nil(testingInstance, retrievedObservations)
	assert.True(testingInstance, errors.Is(observationError, apperrors.ErrObservationNotFound))
}

func TestUpdateObservation_Success(testingInstance *testing.T) {
	updatedEntity := &observation.Observation{
		FHIRResourceID:  "observation-789",
		EncounterFHIRID: "encounter-123",
		PatientFHIRID:   "patient-456",
		LoincCode:       "8867-4",
		CodeDisplay:     "Heart Rate",
		ValueQuantity:   75,
		ValueUnit:       "bpm",
		ObservedAt:      time.Now(),
	}

	observationService := observation.NewService(&mocks.MockObservationRepository{
		UpdateObservationFn: func(ctx context.Context, fhirResourceID string, entity *observation.Observation) (*observation.Observation, error) {
			return updatedEntity, nil
		},
	})

	resultObservation, updateError := observationService.UpdateObservation(context.Background(), "observation-789", &observation.Observation{
		EncounterFHIRID: "encounter-123",
		PatientFHIRID:   "patient-456",
		LoincCode:       "8867-4",
		ValueQuantity:   75,
		ValueUnit:       "bpm",
	})

	assert.NoError(testingInstance, updateError)
	assert.NotNil(testingInstance, resultObservation)
	assert.Equal(testingInstance, "observation-789", resultObservation.FHIRResourceID)
	assert.Equal(testingInstance, "8867-4", resultObservation.LoincCode)
}

func TestUpdateObservation_MissingRequiredFields_ReturnsInvalidArgument(testingInstance *testing.T) {
	observationService := observation.NewService(&mocks.MockObservationRepository{})

	resultObservation, updateError := observationService.UpdateObservation(context.Background(), "observation-789", &observation.Observation{
		LoincCode: "8867-4",
	})

	assert.Error(testingInstance, updateError)
	assert.Nil(testingInstance, resultObservation)
	var appError apperrors.AppError
	assert.ErrorAs(testingInstance, updateError, &appError)
}

func TestUpdateObservation_RepositoryFailure_ReturnsError(testingInstance *testing.T) {
	observationService := observation.NewService(&mocks.MockObservationRepository{
		UpdateObservationFn: func(ctx context.Context, fhirResourceID string, entity *observation.Observation) (*observation.Observation, error) {
			return nil, errRepositoryFailure
		},
	})

	resultObservation, updateError := observationService.UpdateObservation(context.Background(), "observation-789", &observation.Observation{
		EncounterFHIRID: "encounter-123",
		PatientFHIRID:   "patient-456",
		LoincCode:       "8867-4",
	})

	assert.Error(testingInstance, updateError)
	assert.Nil(testingInstance, resultObservation)
	assert.True(testingInstance, errors.Is(updateError, errRepositoryFailure))
}

func TestUpdateObservation_RepositoryNotFound_ReturnsError(testingInstance *testing.T) {
	observationService := observation.NewService(&mocks.MockObservationRepository{
		UpdateObservationFn: func(ctx context.Context, fhirResourceID string, entity *observation.Observation) (*observation.Observation, error) {
			return nil, apperrors.ErrObservationNotFound
		},
	})

	resultObservation, updateError := observationService.UpdateObservation(context.Background(), "observation-789", &observation.Observation{
		EncounterFHIRID: "encounter-123",
		PatientFHIRID:   "patient-456",
		LoincCode:       "8867-4",
	})

	assert.Error(testingInstance, updateError)
	assert.Nil(testingInstance, resultObservation)
	assert.True(testingInstance, errors.Is(updateError, apperrors.ErrObservationNotFound))
}

func TestDeleteObservation_Success(testingInstance *testing.T) {
	observationService := observation.NewService(&mocks.MockObservationRepository{})

	deleteError := observationService.DeleteObservation(context.Background(), "observation-789")

	assert.NoError(testingInstance, deleteError)
}

func TestDeleteObservation_RepositoryFailure_ReturnsError(testingInstance *testing.T) {
	observationService := observation.NewService(&mocks.MockObservationRepository{
		DeleteObservationFn: func(ctx context.Context, fhirResourceID string) error {
			return errRepositoryFailure
		},
	})

	deleteError := observationService.DeleteObservation(context.Background(), "observation-789")

	assert.Error(testingInstance, deleteError)
	assert.True(testingInstance, errors.Is(deleteError, errRepositoryFailure))
}

func TestDeleteObservation_RepositoryNotFound_ReturnsError(testingInstance *testing.T) {
	observationService := observation.NewService(&mocks.MockObservationRepository{
		DeleteObservationFn: func(ctx context.Context, fhirResourceID string) error {
			return apperrors.ErrObservationNotFound
		},
	})

	deleteError := observationService.DeleteObservation(context.Background(), "observation-789")

	assert.Error(testingInstance, deleteError)
	assert.True(testingInstance, errors.Is(deleteError, apperrors.ErrObservationNotFound))
}
