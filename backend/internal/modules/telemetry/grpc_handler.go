package telemetry

import (
	"context"

	"github.com/healthcare/backend/internal/modules/telemetry/pb"
	"github.com/healthcare/backend/internal/shared/apperrors"
)

type GRPCHandler struct {
	service Service
}

func NewGRPCHandler(service Service) *GRPCHandler {
	return &GRPCHandler{service: service}
}

func (handler *GRPCHandler) GetRooms(ctx context.Context, req *pb.GetRoomsRequest) (*pb.GetRoomsResponse, error) {
	rooms, err := handler.service.GetRooms(ctx)
	if err != nil {
		return nil, apperrors.ToGRPCStatus(err)
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
	room, err := handler.service.UnlockRoom(ctx, UnlockRoomInput{
		RoomID:   req.RoomId,
		Passcode: req.Passcode,
	})
	if err != nil {
		return nil, apperrors.ToGRPCStatus(err)
	}

	return &pb.UnlockRoomResponse{
		Success:  true,
		RoomName: room.Name,
	}, nil
}

func (handler *GRPCHandler) GetBeds(ctx context.Context, req *pb.GetBedsRequest) (*pb.GetBedsResponse, error) {
	beds, err := handler.service.GetBeds(ctx, GetBedsInput{RoomID: req.RoomId})
	if err != nil {
		return nil, apperrors.ToGRPCStatus(err)
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
	err := handler.service.UpdateBedCondition(ctx, UpdateBedConditionInput{
		BedID:       req.BedId,
		Bpm:         req.Bpm,
		Spo2:        req.Spo2,
		Temperature: req.Temperature,
		Status:      req.Status,
		Condition:   req.Condition,
	})
	if err != nil {
		return nil, apperrors.ToGRPCStatus(err)
	}

	return &pb.UpdateBedConditionResponse{Success: true}, nil
}
