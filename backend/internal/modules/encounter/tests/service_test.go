package tests

import (
	"context"
	"testing"

	"github.com/healthcare/backend/internal/modules/encounter"
	"github.com/healthcare/backend/internal/modules/encounter/mocks"
	"github.com/stretchr/testify/assert"
)

func TestCreateEncounter_Success(t *testing.T) {
	encounterService := encounter.NewService(&mocks.MockEncounterRepository{}, nil)

	entity := &encounter.Encounter{
		PatientFHIRID:  "patient-fhir-123",
		PractitionerID: "practitioner-456",
		ReasonCode:     "Z00.0",
		ReasonDisplay:  "Routine check-up",
	}

	result, err := encounterService.CreateEncounter(context.Background(), entity)

	assert.NoError(t, err)
	assert.Equal(t, "patient-fhir-123", result.PatientFHIRID)
}

func TestCreateEncounter_MissingPatientFHIRID_ReturnsError(t *testing.T) {
	encounterService := encounter.NewService(&mocks.MockEncounterRepository{}, nil)

	entity := &encounter.Encounter{PatientFHIRID: "", PractitionerID: "practitioner-456"}

	result, err := encounterService.CreateEncounter(context.Background(), entity)

	assert.Error(t, err)
	assert.Nil(t, result)
}
