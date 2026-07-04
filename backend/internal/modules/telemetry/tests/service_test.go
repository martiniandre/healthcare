package tests

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/healthcare/backend/internal/modules/telemetry"
	"github.com/healthcare/backend/internal/modules/telemetry/mocks"
	"github.com/stretchr/testify/assert"
)

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
