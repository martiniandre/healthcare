package tests

import (
	"context"
	"errors"
	"testing"

	"github.com/healthcare/backend/internal/modules/condition"
	"github.com/healthcare/backend/internal/modules/condition/mocks"
	"github.com/healthcare/backend/internal/shared/apperrors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateCondition_ValidInput_DefaultsClinicalStatusToActive(testingInstance *testing.T) {
	conditionService := condition.NewService(&mocks.MockConditionRepository{})

	input := condition.CreateConditionInput{
		PatientFHIRID:   "patient-fhir-123",
		EncounterFHIRID: "encounter-456",
		ICD10Code:       "I10",
		CodeDisplay:     "Essential hypertension",
	}

	result, err := conditionService.CreateCondition(context.Background(), input)

	require.NoError(testingInstance, err)
	assert.Equal(testingInstance, "active", result.ClinicalStatus)
	assert.Equal(testingInstance, "patient-fhir-123", result.PatientFHIRID)
	assert.Equal(testingInstance, "I10", result.ICD10Code)
	assert.Equal(testingInstance, "Essential hypertension", result.CodeDisplay)
	assert.Equal(testingInstance, "encounter-456", result.EncounterFHIRID)
}

func TestCreateCondition_InvalidInput_ReturnsErrorAndDoesNotCallRepository(testingInstance *testing.T) {
	testCases := []struct {
		name             string
		input            condition.CreateConditionInput
		expectedFieldKey string
	}{
		{
			name: "missing patient fhir id",
			input: condition.CreateConditionInput{
				ICD10Code:   "I10",
				CodeDisplay: "Essential hypertension",
			},
			expectedFieldKey: "patient_fhir_id",
		},
		{
			name: "missing icd10 code",
			input: condition.CreateConditionInput{
				PatientFHIRID: "patient-fhir-123",
				CodeDisplay:   "Essential hypertension",
			},
			expectedFieldKey: "icd10_code",
		},
		{
			name: "missing code display",
			input: condition.CreateConditionInput{
				PatientFHIRID: "patient-fhir-123",
				ICD10Code:     "I10",
			},
			expectedFieldKey: "code_display",
		},
		{
			name: "invalid icd10 code",
			input: condition.CreateConditionInput{
				PatientFHIRID: "patient-fhir-123",
				ICD10Code:     "not-a-code",
				CodeDisplay:   "Essential hypertension",
			},
			expectedFieldKey: "icd10_code",
		},
		{
			name: "invalid clinical status",
			input: condition.CreateConditionInput{
				PatientFHIRID:  "patient-fhir-123",
				ICD10Code:      "I10",
				CodeDisplay:    "Essential hypertension",
				ClinicalStatus: "unknown-status",
			},
			expectedFieldKey: "clinical_status",
		},
	}

	for _, testCase := range testCases {
		testingInstance.Run(testCase.name, func(subTest *testing.T) {
			mockRepository := &mocks.MockConditionRepository{
				CreateConditionFn: func(ctx context.Context, entity *condition.Condition) (*condition.Condition, error) {
					subTest.Fatal("repository should not be called for invalid input")
					return nil, nil
				},
			}
			conditionService := condition.NewService(mockRepository)

			result, err := conditionService.CreateCondition(context.Background(), testCase.input)

			require.Error(subTest, err)
			assert.Nil(subTest, result)
			var appError apperrors.AppError
			require.True(subTest, errors.As(err, &appError))
			assert.Contains(subTest, appError.Message, testCase.expectedFieldKey)
		})
	}
}

func TestCreateCondition_RepositoryFailure_ReturnsError(testingInstance *testing.T) {
	expectedErr := errors.New("fhir provider unavailable")
	mockRepository := &mocks.MockConditionRepository{
		CreateConditionFn: func(ctx context.Context, entity *condition.Condition) (*condition.Condition, error) {
			return nil, expectedErr
		},
	}
	conditionService := condition.NewService(mockRepository)

	input := condition.CreateConditionInput{
		PatientFHIRID: "patient-fhir-123",
		ICD10Code:     "I10",
		CodeDisplay:   "Essential hypertension",
	}

	result, err := conditionService.CreateCondition(context.Background(), input)

	assert.Nil(testingInstance, result)
	assert.ErrorIs(testingInstance, err, expectedErr)
}

func TestUpdateCondition_ClinicalStatusOnly_PreservesOtherFieldsAndMerges(testingInstance *testing.T) {
	currentCondition := &condition.Condition{
		FHIRResourceID:  "condition-fhir-789",
		PatientFHIRID:   "patient-fhir-123",
		EncounterFHIRID: "encounter-456",
		ICD10Code:       "I10",
		CodeDisplay:     "Essential hypertension",
		ClinicalStatus:  "active",
	}
	completedStatus := "resolved"

	var mergedEntity *condition.Condition
	mockRepository := &mocks.MockConditionRepository{
		GetConditionByIDFn: func(ctx context.Context, fhirResourceID string) (*condition.Condition, error) {
			return currentCondition, nil
		},
		UpdateConditionFn: func(ctx context.Context, fhirResourceID string, entity *condition.Condition) (*condition.Condition, error) {
			mergedEntity = entity
			return entity, nil
		},
	}
	conditionService := condition.NewService(mockRepository)

	input := condition.UpdateConditionInput{ClinicalStatus: &completedStatus}

	result, err := conditionService.UpdateCondition(context.Background(), "condition-fhir-789", input)

	require.NoError(testingInstance, err)
	require.NotNil(testingInstance, result)
	require.NotNil(testingInstance, mergedEntity)
	assert.Equal(testingInstance, "resolved", mergedEntity.ClinicalStatus)
	assert.Equal(testingInstance, "patient-fhir-123", mergedEntity.PatientFHIRID)
	assert.Equal(testingInstance, "encounter-456", mergedEntity.EncounterFHIRID)
	assert.Equal(testingInstance, "I10", mergedEntity.ICD10Code)
	assert.Equal(testingInstance, "Essential hypertension", mergedEntity.CodeDisplay)
	assert.Equal(testingInstance, "resolved", result.ClinicalStatus)
}

func TestUpdateCondition_ConditionNotFound_ReturnsError(testingInstance *testing.T) {
	mockRepository := &mocks.MockConditionRepository{
		GetConditionByIDFn: func(ctx context.Context, fhirResourceID string) (*condition.Condition, error) {
			return nil, apperrors.ErrConditionNotFound
		},
		UpdateConditionFn: func(ctx context.Context, fhirResourceID string, entity *condition.Condition) (*condition.Condition, error) {
			testingInstance.Fatal("repository update should not be called when condition is missing")
			return nil, nil
		},
	}
	conditionService := condition.NewService(mockRepository)

	clinicalStatus := "resolved"
	input := condition.UpdateConditionInput{ClinicalStatus: &clinicalStatus}

	result, err := conditionService.UpdateCondition(context.Background(), "condition-fhir-404", input)

	require.Error(testingInstance, err)
	assert.Nil(testingInstance, result)
	var appError apperrors.AppError
	require.True(testingInstance, errors.As(err, &appError))
	assert.Equal(testingInstance, "condition not found", appError.Message)
}

func TestUpdateCondition_InvalidICD10Code_ReturnsErrorAndDoesNotCallRepository(testingInstance *testing.T) {
	currentCondition := &condition.Condition{
		FHIRResourceID: "condition-fhir-789",
		PatientFHIRID:  "patient-fhir-123",
		ICD10Code:      "I10",
		CodeDisplay:    "Essential hypertension",
		ClinicalStatus: "active",
	}
	invalidCode := "not-a-code"

	mockRepository := &mocks.MockConditionRepository{
		GetConditionByIDFn: func(ctx context.Context, fhirResourceID string) (*condition.Condition, error) {
			return currentCondition, nil
		},
		UpdateConditionFn: func(ctx context.Context, fhirResourceID string, entity *condition.Condition) (*condition.Condition, error) {
			testingInstance.Fatal("repository should not be called for invalid icd10 code")
			return nil, nil
		},
	}
	conditionService := condition.NewService(mockRepository)

	input := condition.UpdateConditionInput{ICD10Code: &invalidCode}

	result, err := conditionService.UpdateCondition(context.Background(), "condition-fhir-789", input)

	require.Error(testingInstance, err)
	assert.Nil(testingInstance, result)
	var appError apperrors.AppError
	require.True(testingInstance, errors.As(err, &appError))
	assert.Contains(testingInstance, appError.Message, "icd10_code")
}
