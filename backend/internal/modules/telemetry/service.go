package telemetry

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/healthcare/backend/internal/shared/apperrors"
	"github.com/healthcare/backend/internal/shared/eventbus"
)

type Service interface {
	GetRooms(ctx context.Context) ([]*Room, error)
	UnlockRoom(ctx context.Context, input UnlockRoomInput) (*Room, error)
	GetBeds(ctx context.Context, input GetBedsInput) ([]*Bed, error)
	UpdateBedCondition(ctx context.Context, input UpdateBedConditionInput) error
}

type service struct {
	repo     Repository
	eventBus eventbus.Bus
}

func NewService(repo Repository, eventBus eventbus.Bus) Service {
	return &service{repo: repo, eventBus: eventBus}
}

func (telemetryService *service) GetRooms(ctx context.Context) ([]*Room, error) {
	return telemetryService.repo.GetRooms(ctx)
}

func (telemetryService *service) UnlockRoom(ctx context.Context, input UnlockRoomInput) (*Room, error) {
	if fieldViolations := validateUnlockRoomInput(input); len(fieldViolations) > 0 {
		return nil, apperrors.InvalidArgument("invalid unlock room input", fieldViolations)
	}

	roomID, err := uuid.Parse(input.RoomID)
	if err != nil {
		return nil, apperrors.InvalidArgument("invalid unlock room input", map[string]string{"room_id": "must be a valid UUID"})
	}

	room, err := telemetryService.repo.GetRoomByID(ctx, roomID)
	if err != nil {
		return nil, err
	}

	if room.Passcode != input.Passcode {
		return nil, apperrors.ErrInvalidPasscode
	}

	return room, nil
}

func (telemetryService *service) GetBeds(ctx context.Context, input GetBedsInput) ([]*Bed, error) {
	roomID, err := uuid.Parse(input.RoomID)
	if err != nil {
		return nil, apperrors.InvalidArgument("invalid get beds input", map[string]string{"room_id": "must be a valid UUID"})
	}

	_, err = telemetryService.repo.GetRoomByID(ctx, roomID)
	if err != nil {
		return nil, err
	}

	return telemetryService.repo.GetBedsByRoomID(ctx, roomID)
}

func (telemetryService *service) UpdateBedCondition(ctx context.Context, input UpdateBedConditionInput) error {
	if fieldViolations := validateBedConditionInput(input); len(fieldViolations) > 0 {
		return apperrors.InvalidArgument("invalid bed condition input", fieldViolations)
	}

	bedID, err := uuid.Parse(input.BedID)
	if err != nil {
		return apperrors.InvalidArgument("invalid bed condition input", map[string]string{"bed_id": "must be a valid UUID"})
	}

	bed, err := telemetryService.repo.GetBedByID(ctx, bedID)
	if err != nil {
		return err
	}

	previousStatus := bed.Status

	bed.Bpm = input.Bpm
	bed.Spo2 = input.Spo2
	bed.Temperature = input.Temperature
	bed.Status = input.Status
	bed.Condition = input.Condition

	err = telemetryService.repo.UpdateBedCondition(ctx, bed)
	if err != nil {
		return err
	}

	if input.Status == "danger" && previousStatus != "danger" && telemetryService.eventBus != nil {
		telemetryService.eventBus.Publish(ctx, eventbus.Event{
			Name: "telemetry.alert",
			Data: map[string]any{
				"title":         "Alerta Clínico - Leito " + bed.BedNumber,
				"body":          fmt.Sprintf("Paciente %s apresenta condição %s (BPM: %d, SpO2: %d%%).", bed.PatientName, input.Condition, input.Bpm, input.Spo2),
				"resource_type": "bed",
				"resource_id":   bed.ID.String(),
			},
		})
	}

	return nil
}

func validateUnlockRoomInput(input UnlockRoomInput) map[string]string {
	fieldViolations := make(map[string]string)
	if strings.TrimSpace(input.Passcode) == "" {
		fieldViolations["passcode"] = "is required"
	}
	if strings.TrimSpace(input.RoomID) == "" {
		fieldViolations["room_id"] = "is required"
	}
	return fieldViolations
}

func validateBedConditionInput(input UpdateBedConditionInput) map[string]string {
	fieldViolations := make(map[string]string)
	if strings.TrimSpace(input.BedID) == "" {
		fieldViolations["bed_id"] = "is required"
	}
	if input.Bpm < 0 || input.Bpm > 300 {
		fieldViolations["bpm"] = "out of clinical range (0-300)"
	}
	if input.Spo2 < 0 || input.Spo2 > 100 {
		fieldViolations["spo2"] = "out of clinical range (0-100)"
	}
	if input.Temperature < 30.0 || input.Temperature > 45.0 {
		fieldViolations["temperature"] = "out of clinical range (30-45)"
	}
	if strings.TrimSpace(input.Status) == "" {
		fieldViolations["status"] = "is required"
	}
	if strings.TrimSpace(input.Condition) == "" {
		fieldViolations["condition"] = "is required"
	}
	return fieldViolations
}
