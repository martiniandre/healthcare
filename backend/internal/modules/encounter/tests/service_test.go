package tests

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/healthcare/backend/internal/modules/encounter"
	"github.com/healthcare/backend/internal/modules/encounter/mocks"
	"github.com/healthcare/backend/internal/shared/apperrors"
	"github.com/healthcare/backend/internal/shared/eventbus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockEncounterEventBus struct {
	PublishedEvents []eventbus.Event
}

func (mockBus *mockEncounterEventBus) Publish(ctx context.Context, event eventbus.Event) error {
	mockBus.PublishedEvents = append(mockBus.PublishedEvents, event)
	return nil
}

func (mockBus *mockEncounterEventBus) Subscribe(eventName string, handler eventbus.Handler) {}

func TestCreateEncounter_ValidInput_DefaultsToInProgressAndPublishesEvent(testingInstance *testing.T) {
	eventBus := &mockEncounterEventBus{}
	encounterService := encounter.NewService(&mocks.MockEncounterRepository{}, eventBus)

	input := encounter.CreateEncounterInput{
		PatientFHIRID:  "patient-fhir-123",
		PractitionerID: "practitioner-456",
		ReasonCode:     "Z00.0",
		ReasonDisplay:  "Routine check-up",
	}

	result, err := encounterService.CreateEncounter(context.Background(), input)

	require.NoError(testingInstance, err)
	require.NotNil(testingInstance, result)
	assert.Equal(testingInstance, "in-progress", result.Status)
	assert.False(testingInstance, result.StartedAt.IsZero())
	assert.WithinDuration(testingInstance, time.Now(), result.StartedAt, 5*time.Second)
	assert.Equal(testingInstance, "patient-fhir-123", result.PatientFHIRID)
	assert.Equal(testingInstance, "practitioner-456", result.PractitionerID)
	assert.Equal(testingInstance, "Z00.0", result.ReasonCode)
	assert.Equal(testingInstance, "Routine check-up", result.ReasonDisplay)
	require.Len(testingInstance, eventBus.PublishedEvents, 1)
	assert.Equal(testingInstance, "encounter.created", eventBus.PublishedEvents[0].Name)
	assert.Equal(testingInstance, "Novo Atendimento Criado", eventBus.PublishedEvents[0].Data["title"])
	assert.Equal(testingInstance, "encounter", eventBus.PublishedEvents[0].Data["resource_type"])
}

func TestCreateEncounter_MissingRequiredField_ReturnsErrorAndDoesNotCallRepository(testingInstance *testing.T) {
	testCases := []struct {
		name             string
		input            encounter.CreateEncounterInput
		expectedFieldKey string
	}{
		{
			name: "missing patient fhir id",
			input: encounter.CreateEncounterInput{
				PractitionerID: "practitioner-456",
				ReasonDisplay:  "Routine check-up",
			},
			expectedFieldKey: "patient_fhir_id",
		},
		{
			name: "missing practitioner id",
			input: encounter.CreateEncounterInput{
				PatientFHIRID: "patient-fhir-123",
				ReasonDisplay: "Routine check-up",
			},
			expectedFieldKey: "practitioner_id",
		},
		{
			name: "missing reason display",
			input: encounter.CreateEncounterInput{
				PatientFHIRID:  "patient-fhir-123",
				PractitionerID: "practitioner-456",
			},
			expectedFieldKey: "reason",
		},
	}

	for _, testCase := range testCases {
		testingInstance.Run(testCase.name, func(subTest *testing.T) {
			mockRepository := &mocks.MockEncounterRepository{
				CreateEncounterFn: func(ctx context.Context, entity *encounter.Encounter) (*encounter.Encounter, error) {
					subTest.Fatal("repository should not be called for invalid input")
					return nil, nil
				},
			}
			encounterService := encounter.NewService(mockRepository, nil)

			result, err := encounterService.CreateEncounter(context.Background(), testCase.input)

			require.Error(subTest, err)
			assert.Nil(subTest, result)
			var appError apperrors.AppError
			require.True(subTest, errors.As(err, &appError))
			assert.Contains(subTest, appError.Message, testCase.expectedFieldKey)
		})
	}
}

func TestCreateEncounter_RepositoryFailure_ReturnsError(testingInstance *testing.T) {
	expectedErr := errors.New("fhir provider unavailable")
	mockRepository := &mocks.MockEncounterRepository{
		CreateEncounterFn: func(ctx context.Context, entity *encounter.Encounter) (*encounter.Encounter, error) {
			return nil, expectedErr
		},
	}
	encounterService := encounter.NewService(mockRepository, nil)

	input := encounter.CreateEncounterInput{
		PatientFHIRID:  "patient-fhir-123",
		PractitionerID: "practitioner-456",
		ReasonDisplay:  "Routine check-up",
	}

	result, err := encounterService.CreateEncounter(context.Background(), input)

	assert.Nil(testingInstance, result)
	assert.ErrorIs(testingInstance, err, expectedErr)
}

func TestUpdateEncounter_StatusOnly_PreservesOtherFieldsAndMerges(testingInstance *testing.T) {
	currentEncounter := &encounter.Encounter{
		FHIRResourceID: "encounter-fhir-789",
		PatientFHIRID:  "patient-fhir-123",
		PractitionerID: "practitioner-456",
		Status:         "in-progress",
		ReasonCode:     "Z00.0",
		ReasonDisplay:  "Routine check-up",
		StartedAt:      time.Now().Add(-30 * time.Minute),
	}
	finishedStatus := "finished"

	var mergedEntity *encounter.Encounter
	mockRepository := &mocks.MockEncounterRepository{
		GetEncounterByIDFn: func(ctx context.Context, fhirResourceID string) (*encounter.Encounter, error) {
			return currentEncounter, nil
		},
		UpdateEncounterFn: func(ctx context.Context, fhirResourceID string, entity *encounter.Encounter) (*encounter.Encounter, error) {
			mergedEntity = entity
			return entity, nil
		},
	}
	encounterService := encounter.NewService(mockRepository, nil)

	input := encounter.UpdateEncounterInput{Status: &finishedStatus}

	result, err := encounterService.UpdateEncounter(context.Background(), "encounter-fhir-789", input)

	require.NoError(testingInstance, err)
	require.NotNil(testingInstance, result)
	require.NotNil(testingInstance, mergedEntity)
	assert.Equal(testingInstance, "finished", mergedEntity.Status)
	assert.Equal(testingInstance, "patient-fhir-123", mergedEntity.PatientFHIRID)
	assert.Equal(testingInstance, "practitioner-456", mergedEntity.PractitionerID)
	assert.Equal(testingInstance, "Z00.0", mergedEntity.ReasonCode)
	assert.Equal(testingInstance, "Routine check-up", mergedEntity.ReasonDisplay)
	assert.Equal(testingInstance, currentEncounter.StartedAt, mergedEntity.StartedAt)
	assert.Equal(testingInstance, "finished", result.Status)
}

func TestUpdateEncounter_InvalidStatusTransition_ReturnsErrorAndDoesNotCallRepository(testingInstance *testing.T) {
	currentEncounter := &encounter.Encounter{
		FHIRResourceID: "encounter-fhir-789",
		PatientFHIRID:  "patient-fhir-123",
		Status:         "finished",
	}
	cancelledStatus := "cancelled"

	mockRepository := &mocks.MockEncounterRepository{
		GetEncounterByIDFn: func(ctx context.Context, fhirResourceID string) (*encounter.Encounter, error) {
			return currentEncounter, nil
		},
		UpdateEncounterFn: func(ctx context.Context, fhirResourceID string, entity *encounter.Encounter) (*encounter.Encounter, error) {
			testingInstance.Fatal("repository should not be called for invalid status transition")
			return nil, nil
		},
	}
	encounterService := encounter.NewService(mockRepository, nil)

	input := encounter.UpdateEncounterInput{Status: &cancelledStatus}

	result, err := encounterService.UpdateEncounter(context.Background(), "encounter-fhir-789", input)

	require.Error(testingInstance, err)
	assert.Nil(testingInstance, result)
	var appError apperrors.AppError
	require.True(testingInstance, errors.As(err, &appError))
	assert.Contains(testingInstance, appError.Message, "invalid encounter status transition")
}

func TestUpdateEncounter_SameStatusIsAllowedNoOp(testingInstance *testing.T) {
	currentEncounter := &encounter.Encounter{
		FHIRResourceID: "encounter-fhir-789",
		PatientFHIRID:  "patient-fhir-123",
		Status:         "finished",
	}
	finishedStatus := "finished"

	mockRepository := &mocks.MockEncounterRepository{
		GetEncounterByIDFn: func(ctx context.Context, fhirResourceID string) (*encounter.Encounter, error) {
			return currentEncounter, nil
		},
		UpdateEncounterFn: func(ctx context.Context, fhirResourceID string, entity *encounter.Encounter) (*encounter.Encounter, error) {
			return entity, nil
		},
	}
	encounterService := encounter.NewService(mockRepository, nil)

	input := encounter.UpdateEncounterInput{Status: &finishedStatus}

	result, err := encounterService.UpdateEncounter(context.Background(), "encounter-fhir-789", input)

	require.NoError(testingInstance, err)
	assert.NotNil(testingInstance, result)
	assert.Equal(testingInstance, "finished", result.Status)
}

func TestUpdateEncounter_EncounterNotFound_ReturnsError(testingInstance *testing.T) {
	mockRepository := &mocks.MockEncounterRepository{
		GetEncounterByIDFn: func(ctx context.Context, fhirResourceID string) (*encounter.Encounter, error) {
			return nil, apperrors.ErrEncounterNotFound
		},
		UpdateEncounterFn: func(ctx context.Context, fhirResourceID string, entity *encounter.Encounter) (*encounter.Encounter, error) {
			testingInstance.Fatal("repository update should not be called when encounter is missing")
			return nil, nil
		},
	}
	encounterService := encounter.NewService(mockRepository, nil)

	reasonDisplay := "New reason"
	input := encounter.UpdateEncounterInput{ReasonDisplay: &reasonDisplay}

	result, err := encounterService.UpdateEncounter(context.Background(), "encounter-fhir-404", input)

	require.Error(testingInstance, err)
	assert.Nil(testingInstance, result)
	var appError apperrors.AppError
	require.True(testingInstance, errors.As(err, &appError))
	assert.Equal(testingInstance, "encounter not found", appError.Message)
}
