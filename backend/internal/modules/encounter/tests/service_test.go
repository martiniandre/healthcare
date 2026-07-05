package tests

import (
	"context"
	"testing"

	"github.com/healthcare/backend/internal/modules/encounter"
	"github.com/healthcare/backend/internal/modules/encounter/mocks"
	"github.com/healthcare/backend/internal/shared/eventbus"
	"github.com/stretchr/testify/assert"
)

type mockEncounterEventBus struct {
	PublishedEvents []eventbus.Event
}

func (mockBus *mockEncounterEventBus) Publish(ctx context.Context, event eventbus.Event) error {
	mockBus.PublishedEvents = append(mockBus.PublishedEvents, event)
	return nil
}

func (mockBus *mockEncounterEventBus) Subscribe(eventName string, handler eventbus.Handler) {}

func TestCreateEncounter_PublishesEvent(testingInstance *testing.T) {
	eventBus := &mockEncounterEventBus{}
	encounterService := encounter.NewService(&mocks.MockEncounterRepository{}, eventBus)

	entity := &encounter.Encounter{
		PatientFHIRID:  "patient-fhir-123",
		PractitionerID: "practitioner-456",
		ReasonCode:     "Z00.0",
		ReasonDisplay:  "Routine check-up",
	}

	result, err := encounterService.CreateEncounter(context.Background(), entity)

	assert.NoError(testingInstance, err)
	assert.NotNil(testingInstance, result)
	assert.Len(testingInstance, eventBus.PublishedEvents, 1)
	assert.Equal(testingInstance, "encounter.created", eventBus.PublishedEvents[0].Name)
	assert.Equal(testingInstance, "Novo Atendimento Criado", eventBus.PublishedEvents[0].Data["title"])
	assert.Equal(testingInstance, "encounter", eventBus.PublishedEvents[0].Data["resource_type"])
}

func TestCreateEncounter_Success(testingInstance *testing.T) {
	encounterService := encounter.NewService(&mocks.MockEncounterRepository{}, nil)

	entity := &encounter.Encounter{
		PatientFHIRID:  "patient-fhir-123",
		PractitionerID: "practitioner-456",
		ReasonCode:     "Z00.0",
		ReasonDisplay:  "Routine check-up",
	}

	result, err := encounterService.CreateEncounter(context.Background(), entity)

	assert.NoError(testingInstance, err)
	assert.Equal(testingInstance, "patient-fhir-123", result.PatientFHIRID)
}

func TestCreateEncounter_MissingPatientFHIRID_ReturnsError(testingInstance *testing.T) {
	encounterService := encounter.NewService(&mocks.MockEncounterRepository{}, nil)

	entity := &encounter.Encounter{PatientFHIRID: "", PractitionerID: "practitioner-456"}

	result, err := encounterService.CreateEncounter(context.Background(), entity)

	assert.Error(testingInstance, err)
	assert.Nil(testingInstance, result)
}
