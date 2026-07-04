package telemetry

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/healthcare/backend/internal/shared/eventbus"
)

var ErrInvalidPasscode = errors.New("invalid passcode for telemetry room")
var ErrRoomNotFound = errors.New("room not found")
var ErrBedNotFound = errors.New("bed not found")

type Service interface {
	GetRooms(ctx context.Context) ([]*Room, error)
	UnlockRoom(ctx context.Context, roomID uuid.UUID, passcode string) (*Room, error)
	GetBeds(ctx context.Context, roomID uuid.UUID) ([]*Bed, error)
	UpdateBedCondition(ctx context.Context, bedID uuid.UUID, bpm, spo2 int32, temperature float64, status, condition string) error
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

func (telemetryService *service) UnlockRoom(ctx context.Context, roomID uuid.UUID, passcode string) (*Room, error) {
	if strings.TrimSpace(passcode) == "" {
		return nil, errors.New("passcode is required")
	}

	room, err := telemetryService.repo.GetRoomByID(ctx, roomID)
	if err != nil {
		return nil, ErrRoomNotFound
	}

	if room.Passcode != passcode {
		return nil, ErrInvalidPasscode
	}

	return room, nil
}

func (telemetryService *service) GetBeds(ctx context.Context, roomID uuid.UUID) ([]*Bed, error) {
	_, err := telemetryService.repo.GetRoomByID(ctx, roomID)
	if err != nil {
		return nil, ErrRoomNotFound
	}

	return telemetryService.repo.GetBedsByRoomID(ctx, roomID)
}

func (telemetryService *service) UpdateBedCondition(ctx context.Context, bedID uuid.UUID, bpm, spo2 int32, temperature float64, status, condition string) error {
	if bpm < 0 || bpm > 300 {
		return errors.New("BPM out of clinical range (0-300)")
	}
	if spo2 < 0 || spo2 > 100 {
		return errors.New("SpO2 out of clinical range (0-100)")
	}
	if temperature < 30.0 || temperature > 45.0 {
		return errors.New("temperature out of clinical range (30-45)")
	}

	bed, err := telemetryService.repo.GetBedByID(ctx, bedID)
	if err != nil {
		return ErrBedNotFound
	}

	previousStatus := bed.Status

	bed.Bpm = bpm
	bed.Spo2 = spo2
	bed.Temperature = temperature
	bed.Status = status
	bed.Condition = condition

	err = telemetryService.repo.UpdateBedCondition(ctx, bed)
	if err != nil {
		return err
	}

	if status == "danger" && previousStatus != "danger" && telemetryService.eventBus != nil {
		telemetryService.eventBus.Publish(ctx, eventbus.Event{
			Name: "telemetry.alert",
			Data: map[string]any{
				"title":         "Alerta Clínico - Leito " + bed.BedNumber,
				"body":          fmt.Sprintf("Paciente %s apresenta condição %s (BPM: %d, SpO2: %d%%).", bed.PatientName, condition, bpm, spo2),
				"resource_type": "bed",
				"resource_id":   bed.ID.String(),
			},
		})
	}

	return nil
}
