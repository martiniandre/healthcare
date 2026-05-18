package telemetry

import (
	"context"
	"errors"

	"github.com/google/uuid"
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
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (telemetryService *service) GetRooms(ctx context.Context) ([]*Room, error) {
	return telemetryService.repo.GetRooms(ctx)
}

func (telemetryService *service) UnlockRoom(ctx context.Context, roomID uuid.UUID, passcode string) (*Room, error) {
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
	bed, err := telemetryService.repo.GetBedByID(ctx, bedID)
	if err != nil {
		return ErrBedNotFound
	}

	bed.Bpm = bpm
	bed.Spo2 = spo2
	bed.Temperature = temperature
	bed.Status = status
	bed.Condition = condition

	return telemetryService.repo.UpdateBedCondition(ctx, bed)
}
