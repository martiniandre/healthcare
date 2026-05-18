package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/healthcare/backend/internal/modules/telemetry"
)

type MockTelemetryRepository struct {
	Rooms map[uuid.UUID]*telemetry.Room
	Beds  map[uuid.UUID]*telemetry.Bed
	Err   error
}

func NewMockTelemetryRepository() *MockTelemetryRepository {
	return &MockTelemetryRepository{
		Rooms: make(map[uuid.UUID]*telemetry.Room),
		Beds:  make(map[uuid.UUID]*telemetry.Bed),
	}
}

func (mockRepo *MockTelemetryRepository) GetRooms(ctx context.Context) ([]*telemetry.Room, error) {
	if mockRepo.Err != nil {
		return nil, mockRepo.Err
	}

	rooms := make([]*telemetry.Room, 0, len(mockRepo.Rooms))
	for _, room := range mockRepo.Rooms {
		rooms = append(rooms, room)
	}

	return rooms, nil
}

func (mockRepo *MockTelemetryRepository) GetRoomByID(ctx context.Context, roomID uuid.UUID) (*telemetry.Room, error) {
	if mockRepo.Err != nil {
		return nil, mockRepo.Err
	}

	room, exists := mockRepo.Rooms[roomID]
	if !exists {
		return nil, telemetry.ErrRoomNotFound
	}

	return room, nil
}

func (mockRepo *MockTelemetryRepository) GetBedsByRoomID(ctx context.Context, roomID uuid.UUID) ([]*telemetry.Bed, error) {
	if mockRepo.Err != nil {
		return nil, mockRepo.Err
	}

	beds := make([]*telemetry.Bed, 0)
	for _, bed := range mockRepo.Beds {
		if bed.RoomID == roomID {
			beds = append(beds, bed)
		}
	}

	return beds, nil
}

func (mockRepo *MockTelemetryRepository) GetBedByID(ctx context.Context, bedID uuid.UUID) (*telemetry.Bed, error) {
	if mockRepo.Err != nil {
		return nil, mockRepo.Err
	}

	bed, exists := mockRepo.Beds[bedID]
	if !exists {
		return nil, telemetry.ErrBedNotFound
	}

	return bed, nil
}

func (mockRepo *MockTelemetryRepository) UpdateBedCondition(ctx context.Context, bed *telemetry.Bed) error {
	if mockRepo.Err != nil {
		return mockRepo.Err
	}

	mockRepo.Beds[bed.ID] = bed
	return nil
}
