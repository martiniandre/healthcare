package tests

import (
	"context"
	"errors"
	"testing"

	"github.com/healthcare/backend/internal/modules/medication"
	"github.com/healthcare/backend/internal/modules/medication/mocks"
	"github.com/healthcare/backend/internal/shared/apperrors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateMedicationRequest_ValidInput_DefaultsToActive(testingInstance *testing.T) {
	medicationService := medication.NewService(&mocks.MockMedicationRepository{})

	input := medication.CreateMedicationInput{
		EncounterFHIRID:    "encounter-fhir-456",
		PatientFHIRID:      "patient-fhir-123",
		PractitionerFHIRID: "practitioner-fhir-789",
		MedicationCode:     "10582",
		MedicationName:     "Amoxicillin",
		DosageInstructions: "500mg every 8 hours",
	}

	result, err := medicationService.CreateMedicationRequest(context.Background(), input)

	require.NoError(testingInstance, err)
	assert.Equal(testingInstance, "active", result.Status)
	assert.Equal(testingInstance, "encounter-fhir-456", result.EncounterFHIRID)
	assert.Equal(testingInstance, "patient-fhir-123", result.PatientFHIRID)
	assert.Equal(testingInstance, "10582", result.MedicationCode)
	assert.Equal(testingInstance, "Amoxicillin", result.MedicationName)
}

func TestCreateMedicationRequest_MissingField_ReturnsErrorAndDoesNotCallRepository(testingInstance *testing.T) {
	testCases := []struct {
		name             string
		input            medication.CreateMedicationInput
		expectedFieldKey string
	}{
		{
			name: "missing encounter fhir id",
			input: medication.CreateMedicationInput{
				PatientFHIRID: "patient-fhir-123",
				MedicationCode: "10582",
			},
			expectedFieldKey: "encounter_fhir_id",
		},
		{
			name: "missing patient fhir id",
			input: medication.CreateMedicationInput{
				EncounterFHIRID: "encounter-fhir-456",
				MedicationCode:  "10582",
			},
			expectedFieldKey: "patient_fhir_id",
		},
		{
			name: "missing medication name and code",
			input: medication.CreateMedicationInput{
				EncounterFHIRID:    "encounter-fhir-456",
				PatientFHIRID:      "patient-fhir-123",
				DosageInstructions: "500mg every 8 hours",
			},
			expectedFieldKey: "medication",
		},
	}

	for _, testCase := range testCases {
		testingInstance.Run(testCase.name, func(subTest *testing.T) {
			mockRepository := &mocks.MockMedicationRepository{
				CreateMedicationRequestFn: func(ctx context.Context, entity *medication.Medication) (*medication.Medication, error) {
					subTest.Fatal("repository should not be called for invalid input")
					return nil, nil
				},
			}
			medicationService := medication.NewService(mockRepository)

			result, err := medicationService.CreateMedicationRequest(context.Background(), testCase.input)

			require.Error(subTest, err)
			assert.Nil(subTest, result)
			var appError apperrors.AppError
			require.True(subTest, errors.As(err, &appError))
			assert.Contains(subTest, appError.Message, testCase.expectedFieldKey)
		})
	}
}

func TestCreateMedicationRequest_RepositoryFailure_ReturnsError(testingInstance *testing.T) {
	expectedErr := errors.New("fhir provider unavailable")
	mockRepository := &mocks.MockMedicationRepository{
		CreateMedicationRequestFn: func(ctx context.Context, entity *medication.Medication) (*medication.Medication, error) {
			return nil, expectedErr
		},
	}
	medicationService := medication.NewService(mockRepository)

	input := medication.CreateMedicationInput{
		EncounterFHIRID: "encounter-fhir-456",
		PatientFHIRID:   "patient-fhir-123",
		MedicationCode:  "10582",
	}

	result, err := medicationService.CreateMedicationRequest(context.Background(), input)

	assert.Nil(testingInstance, result)
	assert.ErrorIs(testingInstance, err, expectedErr)
}
