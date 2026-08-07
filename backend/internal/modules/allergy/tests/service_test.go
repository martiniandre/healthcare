package tests

import (
	"context"
	"errors"
	"testing"

	"github.com/healthcare/backend/internal/modules/allergy"
	"github.com/healthcare/backend/internal/modules/allergy/mocks"
	"github.com/healthcare/backend/internal/shared/apperrors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateAllergyIntolerance_ValidInput_DefaultsToActive(testingInstance *testing.T) {
	allergyService := allergy.NewService(&mocks.MockAllergyRepository{})

	input := allergy.CreateAllergyInput{
		PatientFHIRID:   "patient-fhir-123",
		AllergenCode:    "7980",
		AllergenDisplay: "Penicillin",
		Reaction:        "Anaphylaxis",
	}

	result, err := allergyService.CreateAllergyIntolerance(context.Background(), input)

	require.NoError(testingInstance, err)
	assert.Equal(testingInstance, "active", result.ClinicalStatus)
	assert.Equal(testingInstance, "patient-fhir-123", result.PatientFHIRID)
	assert.Equal(testingInstance, "7980", result.AllergenCode)
	assert.Equal(testingInstance, "Penicillin", result.AllergenDisplay)
	assert.Equal(testingInstance, "Anaphylaxis", result.Reaction)
}

func TestCreateAllergyIntolerance_MissingRequiredField_ReturnsErrorAndDoesNotCallRepository(testingInstance *testing.T) {
	testCases := []struct {
		name             string
		input            allergy.CreateAllergyInput
		expectedFieldKey string
	}{
		{
			name: "missing patient fhir id",
			input: allergy.CreateAllergyInput{
				AllergenCode: "7980",
			},
			expectedFieldKey: "patient_fhir_id",
		},
		{
			name: "missing allergen code",
			input: allergy.CreateAllergyInput{
				PatientFHIRID: "patient-fhir-123",
			},
			expectedFieldKey: "allergen_code",
		},
		{
			name: "invalid clinical status",
			input: allergy.CreateAllergyInput{
				PatientFHIRID:  "patient-fhir-123",
				AllergenCode:   "7980",
				ClinicalStatus: "unknown-status",
			},
			expectedFieldKey: "clinical_status",
		},
	}

	for _, testCase := range testCases {
		testingInstance.Run(testCase.name, func(subTest *testing.T) {
			mockRepository := &mocks.MockAllergyRepository{
				CreateAllergyIntoleranceFn: func(ctx context.Context, entity *allergy.Allergy) (*allergy.Allergy, error) {
					subTest.Fatal("repository should not be called for invalid input")
					return nil, nil
				},
			}
			allergyService := allergy.NewService(mockRepository)

			result, err := allergyService.CreateAllergyIntolerance(context.Background(), testCase.input)

			require.Error(subTest, err)
			assert.Nil(subTest, result)
			var appError apperrors.AppError
			require.True(subTest, errors.As(err, &appError))
			assert.Contains(subTest, appError.Message, testCase.expectedFieldKey)
		})
	}
}

func TestCreateAllergyIntolerance_RepositoryFailure_ReturnsError(testingInstance *testing.T) {
	expectedErr := errors.New("fhir provider unavailable")
	mockRepository := &mocks.MockAllergyRepository{
		CreateAllergyIntoleranceFn: func(ctx context.Context, entity *allergy.Allergy) (*allergy.Allergy, error) {
			return nil, expectedErr
		},
	}
	allergyService := allergy.NewService(mockRepository)

	input := allergy.CreateAllergyInput{
		PatientFHIRID: "patient-fhir-123",
		AllergenCode:  "7980",
	}

	result, err := allergyService.CreateAllergyIntolerance(context.Background(), input)

	assert.Nil(testingInstance, result)
	assert.ErrorIs(testingInstance, err, expectedErr)
}

func TestUpdateAllergyIntolerance_StatusOnly_PreservesOtherFieldsAndMerges(testingInstance *testing.T) {
	currentAllergy := &allergy.Allergy{
		FHIRResourceID:  "allergy-fhir-789",
		PatientFHIRID:   "patient-fhir-123",
		AllergenCode:    "7980",
		AllergenDisplay: "Penicillin",
		ClinicalStatus:  "active",
		Reaction:        "Anaphylaxis",
	}
	completedStatus := "inactive"

	var mergedEntity *allergy.Allergy
	mockRepository := &mocks.MockAllergyRepository{
		GetAllergyIntoleranceByIDFn: func(ctx context.Context, fhirResourceID string) (*allergy.Allergy, error) {
			return currentAllergy, nil
		},
		UpdateAllergyIntoleranceFn: func(ctx context.Context, fhirResourceID string, entity *allergy.Allergy) (*allergy.Allergy, error) {
			mergedEntity = entity
			return entity, nil
		},
	}
	allergyService := allergy.NewService(mockRepository)

	input := allergy.UpdateAllergyInput{ClinicalStatus: &completedStatus}

	result, err := allergyService.UpdateAllergyIntolerance(context.Background(), "allergy-fhir-789", input)

	require.NoError(testingInstance, err)
	require.NotNil(testingInstance, result)
	require.NotNil(testingInstance, mergedEntity)
	assert.Equal(testingInstance, "inactive", mergedEntity.ClinicalStatus)
	assert.Equal(testingInstance, "patient-fhir-123", mergedEntity.PatientFHIRID)
	assert.Equal(testingInstance, "7980", mergedEntity.AllergenCode)
	assert.Equal(testingInstance, "Penicillin", mergedEntity.AllergenDisplay)
	assert.Equal(testingInstance, "Anaphylaxis", mergedEntity.Reaction)
	assert.Equal(testingInstance, "inactive", result.ClinicalStatus)
}

func TestUpdateAllergyIntolerance_AllergyNotFound_ReturnsError(testingInstance *testing.T) {
	mockRepository := &mocks.MockAllergyRepository{
		GetAllergyIntoleranceByIDFn: func(ctx context.Context, fhirResourceID string) (*allergy.Allergy, error) {
			return nil, apperrors.ErrAllergyIntoleranceNotFound
		},
		UpdateAllergyIntoleranceFn: func(ctx context.Context, fhirResourceID string, entity *allergy.Allergy) (*allergy.Allergy, error) {
			testingInstance.Fatal("repository update should not be called when allergy is missing")
			return nil, nil
		},
	}
	allergyService := allergy.NewService(mockRepository)

	clinicalStatus := "inactive"
	input := allergy.UpdateAllergyInput{ClinicalStatus: &clinicalStatus}

	result, err := allergyService.UpdateAllergyIntolerance(context.Background(), "allergy-fhir-404", input)

	require.Error(testingInstance, err)
	assert.Nil(testingInstance, result)
	var appError apperrors.AppError
	require.True(testingInstance, errors.As(err, &appError))
	assert.Equal(testingInstance, "allergy intolerance not found", appError.Message)
}

func TestUpdateAllergyIntolerance_InvalidClinicalStatus_ReturnsErrorAndDoesNotCallRepository(testingInstance *testing.T) {
	currentAllergy := &allergy.Allergy{
		FHIRResourceID: "allergy-fhir-789",
		PatientFHIRID:  "patient-fhir-123",
		AllergenCode:   "7980",
		ClinicalStatus: "active",
	}
	invalidStatus := "unknown-status"

	mockRepository := &mocks.MockAllergyRepository{
		GetAllergyIntoleranceByIDFn: func(ctx context.Context, fhirResourceID string) (*allergy.Allergy, error) {
			return currentAllergy, nil
		},
		UpdateAllergyIntoleranceFn: func(ctx context.Context, fhirResourceID string, entity *allergy.Allergy) (*allergy.Allergy, error) {
			testingInstance.Fatal("repository should not be called for invalid clinical status")
			return nil, nil
		},
	}
	allergyService := allergy.NewService(mockRepository)

	input := allergy.UpdateAllergyInput{ClinicalStatus: &invalidStatus}

	result, err := allergyService.UpdateAllergyIntolerance(context.Background(), "allergy-fhir-789", input)

	require.Error(testingInstance, err)
	assert.Nil(testingInstance, result)
	var appError apperrors.AppError
	require.True(testingInstance, errors.As(err, &appError))
	assert.Contains(testingInstance, appError.Message, "clinical_status")
}
