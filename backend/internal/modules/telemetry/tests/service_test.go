package tests

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/healthcare/backend/internal/modules/telemetry"
	"github.com/healthcare/backend/internal/modules/telemetry/mocks"
	"github.com/healthcare/backend/internal/shared/eventbus"
	"github.com/stretchr/testify/assert"
)

type mockTelemetryEventBus struct {
	PublishedEvents []eventbus.Event
}

func (mockBus *mockTelemetryEventBus) Publish(ctx context.Context, event eventbus.Event) error {
	mockBus.PublishedEvents = append(mockBus.PublishedEvents, event)
	return nil
}

func (mockBus *mockTelemetryEventBus) Subscribe(eventName string, handler eventbus.Handler) {}

func TestTelemetryService_UpdateBedCondition_PublishesAlert(testingInstance *testing.T) {
	eventBus := &mockTelemetryEventBus{}
	mockRepository := mocks.NewMockTelemetryRepository()
	telemetryService := telemetry.NewService(mockRepository, eventBus)
	contextParam := context.Background()

	bedID := uuid.New()
	mockRepository.Beds[bedID] = &telemetry.Bed{
		ID:          bedID,
		BedNumber:   "Leito 03",
		PatientName: "Carlos Souza",
		Bpm:         72,
		Spo2:        98,
		Temperature: 36.5,
		Status:      "normal",
		Condition:   "Estavel",
	}

	err := telemetryService.UpdateBedCondition(contextParam, bedID, 42, 85, 35.1, "danger", "Hipotermia leve com bradicardia")

	assert.NoError(testingInstance, err)
	assert.Len(testingInstance, eventBus.PublishedEvents, 1)
	assert.Equal(testingInstance, "telemetry.alert", eventBus.PublishedEvents[0].Name)
	assert.Equal(testingInstance, "bed", eventBus.PublishedEvents[0].Data["resource_type"])
	assert.Contains(testingInstance, eventBus.PublishedEvents[0].Data["title"], "Leito 03")
	assert.Contains(testingInstance, eventBus.PublishedEvents[0].Data["body"], "Hipotermia")
}

func TestTelemetryService_UpdateBedCondition_NoAlertOnSameStatus(testingInstance *testing.T) {
	eventBus := &mockTelemetryEventBus{}
	mockRepository := mocks.NewMockTelemetryRepository()
	telemetryService := telemetry.NewService(mockRepository, eventBus)
	contextParam := context.Background()

	bedID := uuid.New()
	mockRepository.Beds[bedID] = &telemetry.Bed{
		ID:          bedID,
		BedNumber:   "Leito 04",
		PatientName: "Ana Paula",
		Bpm:         120,
		Spo2:        88,
		Temperature: 38.5,
		Status:      "danger",
		Condition:   "Taquicardia",
	}

	err := telemetryService.UpdateBedCondition(contextParam, bedID, 125, 85, 39.0, "danger", "Taquicardia persistente")

	assert.NoError(testingInstance, err)
	assert.Len(testingInstance, eventBus.PublishedEvents, 0, "Should not publish alert when status was already danger")
}

func TestTelemetryService_GetRooms(testingInstance *testing.T) {
	mockRepository := mocks.NewMockTelemetryRepository()
	telemetryService := telemetry.NewService(mockRepository, nil)
	contextParam := context.Background()

	roomID := uuid.New()
	mockRepository.Rooms[roomID] = &telemetry.Room{
		ID:          roomID,
		Name:        "Sala Verde",
		Passcode:    "1234",
		Description: "Estável",
	}

	rooms, err := telemetryService.GetRooms(contextParam)
	assert.NoError(testingInstance, err)
	assert.Len(testingInstance, rooms, 1)
	assert.Equal(testingInstance, "Sala Verde", rooms[0].Name)
}

func TestTelemetryService_UnlockRoom(testingInstance *testing.T) {
	mockRepository := mocks.NewMockTelemetryRepository()
	telemetryService := telemetry.NewService(mockRepository, nil)
	contextParam := context.Background()

	roomID := uuid.New()
	mockRepository.Rooms[roomID] = &telemetry.Room{
		ID:          roomID,
		Name:        "Sala Vermelha",
		Passcode:    "4321",
		Description: "UTI",
	}

	unlockedRoom, err := telemetryService.UnlockRoom(contextParam, roomID, "4321")
	assert.NoError(testingInstance, err)
	assert.NotNil(testingInstance, unlockedRoom)
	assert.Equal(testingInstance, "Sala Vermelha", unlockedRoom.Name)

	_, errInvalid := telemetryService.UnlockRoom(contextParam, roomID, "wrong_pin")
	assert.ErrorIs(testingInstance, errInvalid, telemetry.ErrInvalidPasscode)

	_, errNotFound := telemetryService.UnlockRoom(contextParam, uuid.New(), "4321")
	assert.ErrorIs(testingInstance, errNotFound, telemetry.ErrRoomNotFound)
}

func TestTelemetryService_GetBeds(testingInstance *testing.T) {
	mockRepository := mocks.NewMockTelemetryRepository()
	telemetryService := telemetry.NewService(mockRepository, nil)
	contextParam := context.Background()

	roomID := uuid.New()
	mockRepository.Rooms[roomID] = &telemetry.Room{
		ID:          roomID,
		Name:        "Sala Amarela",
		Passcode:    "9999",
		Description: "Semi",
	}

	bedID := uuid.New()
	mockRepository.Beds[bedID] = &telemetry.Bed{
		ID:          bedID,
		RoomID:      roomID,
		BedNumber:   "Leito 01",
		PatientName: "Ana Silva",
		Bpm:         80,
		Spo2:        98,
		Temperature: 36.5,
		Status:      "normal",
		Condition:   "Normal",
	}

	beds, err := telemetryService.GetBeds(contextParam, roomID)
	assert.NoError(testingInstance, err)
	assert.Len(testingInstance, beds, 1)
	assert.Equal(testingInstance, "Ana Silva", beds[0].PatientName)
}

func TestTelemetryService_UpdateBedCondition(testingInstance *testing.T) {
	mockRepository := mocks.NewMockTelemetryRepository()
	telemetryService := telemetry.NewService(mockRepository, nil)
	contextParam := context.Background()

	roomID := uuid.New()
	bedID := uuid.New()
	mockRepository.Beds[bedID] = &telemetry.Bed{
		ID:          bedID,
		RoomID:      roomID,
		BedNumber:   "Leito 02",
		PatientName: "Bruno Costa",
		Bpm:         80,
		Spo2:        98,
		Temperature: 36.5,
		Status:      "normal",
		Condition:   "Normal",
	}

	err := telemetryService.UpdateBedCondition(contextParam, bedID, 52, 95, 37.1, "warning", "Bradicardia")
	assert.NoError(testingInstance, err)

	updatedBed := mockRepository.Beds[bedID]
	assert.Equal(testingInstance, int32(52), updatedBed.Bpm)
	assert.Equal(testingInstance, "warning", updatedBed.Status)
	assert.Equal(testingInstance, "Bradicardia", updatedBed.Condition)
}
