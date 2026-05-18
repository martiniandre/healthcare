package telemetry

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/healthcare/backend/internal/modules/telemetry/pb"
	"github.com/healthcare/backend/internal/shared/apperrors"
)

func mapTelemetryError(err error) error {
	if errors.Is(err, ErrInvalidPasscode) {
		return apperrors.ErrInvalidPasscode.ToGRPC()
	}
	if errors.Is(err, ErrRoomNotFound) {
		return apperrors.ErrRoomNotFound.ToGRPC()
	}
	if errors.Is(err, ErrBedNotFound) {
		return apperrors.ErrBedNotFound.ToGRPC()
	}
	return apperrors.ToGRPCStatus(err)
}

type GRPCHandler struct {
	service Service
}

func NewGRPCHandler(service Service) *GRPCHandler {
	return &GRPCHandler{service: service}
}

func (handler *GRPCHandler) GetRooms(ctx context.Context, req *pb.GetRoomsRequest) (*pb.GetRoomsResponse, error) {
	rooms, err := handler.service.GetRooms(ctx)
	if err != nil {
		return nil, mapTelemetryError(err)
	}

	roomResponses := make([]*pb.TelemetryRoom, 0, len(rooms))
	for _, room := range rooms {
		roomResponses = append(roomResponses, &pb.TelemetryRoom{
			Id:          room.ID.String(),
			Name:        room.Name,
			Description: room.Description,
		})
	}

	return &pb.GetRoomsResponse{Rooms: roomResponses}, nil
}

func (handler *GRPCHandler) UnlockRoom(ctx context.Context, req *pb.UnlockRoomRequest) (*pb.UnlockRoomResponse, error) {
	roomID, err := uuid.Parse(req.RoomId)
	if err != nil {
		return nil, apperrors.ErrBadRequest.ToGRPC()
	}

	room, err := handler.service.UnlockRoom(ctx, roomID, req.Passcode)
	if err != nil {
		return nil, mapTelemetryError(err)
	}

	return &pb.UnlockRoomResponse{
		Success:  true,
		RoomName: room.Name,
	}, nil
}

func (handler *GRPCHandler) GetBeds(ctx context.Context, req *pb.GetBedsRequest) (*pb.GetBedsResponse, error) {
	roomID, err := uuid.Parse(req.RoomId)
	if err != nil {
		return nil, apperrors.ErrBadRequest.ToGRPC()
	}

	beds, err := handler.service.GetBeds(ctx, roomID)
	if err != nil {
		return nil, mapTelemetryError(err)
	}

	bedResponses := make([]*pb.TelemetryBed, 0, len(beds))
	for _, bed := range beds {
		bedResponses = append(bedResponses, &pb.TelemetryBed{
			Id:          bed.ID.String(),
			RoomId:      bed.RoomID.String(),
			BedNumber:   bed.BedNumber,
			PatientName: bed.PatientName,
			Age:         bed.Age,
			Gender:      bed.Gender,
			Bpm:         bed.Bpm,
			Spo2:        bed.Spo2,
			Temperature: bed.Temperature,
			Status:      bed.Status,
			Condition:   bed.Condition,
		})
	}

	return &pb.GetBedsResponse{Beds: bedResponses}, nil
}

func (handler *GRPCHandler) UpdateBedCondition(ctx context.Context, req *pb.UpdateBedConditionRequest) (*pb.UpdateBedConditionResponse, error) {
	bedID, err := uuid.Parse(req.BedId)
	if err != nil {
		return nil, apperrors.ErrBadRequest.ToGRPC()
	}

	err = handler.service.UpdateBedCondition(ctx, bedID, req.Bpm, req.Spo2, req.Temperature, req.Status, req.Condition)
	if err != nil {
		return nil, mapTelemetryError(err)
	}

	return &pb.UpdateBedConditionResponse{Success: true}, nil
}
