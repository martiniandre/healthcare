package tests

import (
	"context"
	"testing"

	"github.com/healthcare/backend/internal/modules/medication"
	"github.com/healthcare/backend/internal/modules/medication/mocks"
	"github.com/stretchr/testify/assert"
)

func TestCreateMedicationRequest_DefaultsToActive(t *testing.T) {
	medicationService := medication.NewService(&mocks.MockMedicationRepository{})

	entity := &medication.Medication{
		PatientFHIRID:      "patient-123",
		EncounterFHIRID:    "encounter-456",
		MedicationCode:     "10582",
		MedicationName:     "Amoxicillin",
		DosageInstructions: "500mg every 8 hours",
	}

	result, err := medicationService.CreateMedicationRequest(context.Background(), entity)

	assert.NoError(t, err)
	assert.Equal(t, "active", result.Status)
}
