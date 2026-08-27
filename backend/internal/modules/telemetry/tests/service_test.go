package tests

import (
	"context"
	"testing"

	"golang.org/x/crypto/bcrypt"

	"github.com/google/uuid"
	"github.com/healthcare/backend/internal/modules/telemetry"
	"github.com/healthcare/backend/internal/modules/telemetry/mocks"
	"github.com/healthcare/backend/internal/shared/apperrors"
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

	err := telemetryService.UpdateBedCondition(contextParam, telemetry.UpdateBedConditionInput{
		BedID:       bedID.String(),
		Bpm:         42,
		Spo2:        85,
		Temperature: 35.1,
		Status:      "danger",
		Condition:   "Hipotermia leve com bradicardia",
	})

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

	err := telemetryService.UpdateBedCondition(contextParam, telemetry.UpdateBedConditionInput{
		BedID:       bedID.String(),
		Bpm:         125,
		Spo2:        85,
		Temperature: 39.0,
		Status:      "danger",
		Condition:   "Taquicardia persistente",
	})

	assert.NoError(testingInstance, err)
	assert.Len(testingInstance, eventBus.PublishedEvents, 0, "Should not publish alert when status was already danger")
}

func TestTelemetryService_UpdateBedCondition_InvalidClinicalRange(testingInstance *testing.T) {
	mockRepository := mocks.NewMockTelemetryRepository()
	telemetryService := telemetry.NewService(mockRepository, nil)
	contextParam := context.Background()

	bedID := uuid.New()
	mockRepository.Beds[bedID] = &telemetry.Bed{
		ID:          bedID,
		BedNumber:   "Leito 05",
		PatientName: "Marcos Lima",
		Bpm:         80,
		Spo2:        98,
		Temperature: 36.5,
		Status:      "normal",
		Condition:   "Normal",
	}

	err := telemetryService.UpdateBedCondition(contextParam, telemetry.UpdateBedConditionInput{
		BedID:       bedID.String(),
		Bpm:         350,
		Spo2:        120,
		Temperature: 50.0,
		Status:      "danger",
		Condition:   "Critico",
	})

	var appError apperrors.AppError
	assert.ErrorAs(testingInstance, err, &appError)
	assert.Equal(testingInstance, 400, appError.HTTPCode)
	assert.Contains(testingInstance, appError.Message, "bpm")
	assert.Contains(testingInstance, appError.Message, "spo2")
	assert.Contains(testingInstance, appError.Message, "temperature")
}

func TestTelemetryService_UpdateBedCondition_BedNotFound(testingInstance *testing.T) {
	mockRepository := mocks.NewMockTelemetryRepository()
	telemetryService := telemetry.NewService(mockRepository, nil)
	contextParam := context.Background()

	err := telemetryService.UpdateBedCondition(contextParam, telemetry.UpdateBedConditionInput{
		BedID:       uuid.New().String(),
		Bpm:         80,
		Spo2:        98,
		Temperature: 36.5,
		Status:      "normal",
		Condition:   "Normal",
	})

	assert.ErrorIs(testingInstance, err, apperrors.ErrBedNotFound)
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
	hashedPasscode, hashError := bcrypt.GenerateFromPassword([]byte("4321"), bcrypt.MinCost)
	assert.NoError(testingInstance, hashError)
	mockRepository.Rooms[roomID] = &telemetry.Room{
		ID:          roomID,
		Name:        "Sala Vermelha",
		Passcode:    string(hashedPasscode),
		Description: "UTI",
	}

	unlockedRoom, err := telemetryService.UnlockRoom(contextParam, telemetry.UnlockRoomInput{
		RoomID:   roomID.String(),
		Passcode: "4321",
	})
	assert.NoError(testingInstance, err)
	assert.NotNil(testingInstance, unlockedRoom)
	assert.Equal(testingInstance, "Sala Vermelha", unlockedRoom.Name)

	_, errInvalid := telemetryService.UnlockRoom(contextParam, telemetry.UnlockRoomInput{
		RoomID:   roomID.String(),
		Passcode: "wrong_pin",
	})
	assert.ErrorIs(testingInstance, errInvalid, apperrors.ErrInvalidPasscode)

	_, errNotFound := telemetryService.UnlockRoom(contextParam, telemetry.UnlockRoomInput{
		RoomID:   uuid.New().String(),
		Passcode: "4321",
	})
	assert.ErrorIs(testingInstance, errNotFound, apperrors.ErrRoomNotFound)
}

func TestTelemetryService_UnlockRoom_InvalidInput(testingInstance *testing.T) {
	mockRepository := mocks.NewMockTelemetryRepository()
	telemetryService := telemetry.NewService(mockRepository, nil)
	contextParam := context.Background()

	_, err := telemetryService.UnlockRoom(contextParam, telemetry.UnlockRoomInput{
		RoomID:   "not-a-uuid",
		Passcode: "4321",
	})

	var appError apperrors.AppError
	assert.ErrorAs(testingInstance, err, &appError)
	assert.Equal(testingInstance, 400, appError.HTTPCode)
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

	beds, err := telemetryService.GetBeds(contextParam, telemetry.GetBedsInput{RoomID: roomID.String()})
	assert.NoError(testingInstance, err)
	assert.Len(testingInstance, beds, 1)
	assert.Equal(testingInstance, "Ana Silva", beds[0].PatientName)

	_, errNotFound := telemetryService.GetBeds(contextParam, telemetry.GetBedsInput{RoomID: uuid.New().String()})
	assert.ErrorIs(testingInstance, errNotFound, apperrors.ErrRoomNotFound)
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

	err := telemetryService.UpdateBedCondition(contextParam, telemetry.UpdateBedConditionInput{
		BedID:       bedID.String(),
		Bpm:         52,
		Spo2:        95,
		Temperature: 37.1,
		Status:      "warning",
		Condition:   "Bradicardia",
	})
	assert.NoError(testingInstance, err)

	updatedBed := mockRepository.Beds[bedID]
	assert.Equal(testingInstance, int32(52), updatedBed.Bpm)
	assert.Equal(testingInstance, "warning", updatedBed.Status)
	assert.Equal(testingInstance, "Bradicardia", updatedBed.Condition)
}
