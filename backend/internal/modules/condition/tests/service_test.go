package tests

import (
	"context"
	"testing"

	"github.com/healthcare/backend/internal/modules/condition"
	"github.com/healthcare/backend/internal/modules/condition/mocks"
	"github.com/stretchr/testify/assert"
)

func TestCreateCondition_DefaultsClinicalStatusToActive(t *testing.T) {
	conditionService := condition.NewService(&mocks.MockConditionRepository{})

	entity := &condition.Condition{
		PatientFHIRID:   "patient-123",
		EncounterFHIRID: "encounter-456",
		ICD10Code:       "I10",
		CodeDisplay:     "Essential hypertension",
		ClinicalStatus:  "",
	}

	result, err := conditionService.CreateCondition(context.Background(), entity)

	assert.NoError(t, err)
	assert.Equal(t, "active", result.ClinicalStatus)
}

func TestCreateCondition_MissingICD10Code_ReturnsError(t *testing.T) {
	conditionService := condition.NewService(&mocks.MockConditionRepository{})

	entity := &condition.Condition{PatientFHIRID: "patient-123", ICD10Code: ""}

	result, err := conditionService.CreateCondition(context.Background(), entity)

	assert.Error(t, err)
	assert.Nil(t, result)
}
