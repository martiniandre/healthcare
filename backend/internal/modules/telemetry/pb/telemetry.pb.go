package pb

import "context"

type TelemetryRoom struct {
	Id          string
	Name        string
	Description string
}

type TelemetryBed struct {
	Id          string
	RoomId      string
	BedNumber   string
	PatientName string
	Age         int32
	Gender      string
	Bpm         int32
	Spo2        int32
	Temperature float64
	Status      string
	Condition   string
}

type GetRoomsRequest struct{}

type GetRoomsResponse struct {
	Rooms []*TelemetryRoom
}

type UnlockRoomRequest struct {
	RoomId   string
	Passcode string
}

type UnlockRoomResponse struct {
	Success  bool
	RoomName string
}

type GetBedsRequest struct {
	RoomId string
}

type GetBedsResponse struct {
	Beds []*TelemetryBed
}

type UpdateBedConditionRequest struct {
	BedId       string
	Bpm         int32
	Spo2        int32
	Temperature float64
	Status      string
	Condition   string
}

type UpdateBedConditionResponse struct {
	Success bool
}

type TelemetryServiceServer interface {
	GetRooms(ctx context.Context, req *GetRoomsRequest) (*GetRoomsResponse, error)
	UnlockRoom(ctx context.Context, req *UnlockRoomRequest) (*UnlockRoomResponse, error)
	GetBeds(ctx context.Context, req *GetBedsRequest) (*GetBedsResponse, error)
	UpdateBedCondition(ctx context.Context, req *UpdateBedConditionRequest) (*UpdateBedConditionResponse, error)
}

func RegisterTelemetryServiceServer(_ interface{}, server TelemetryServiceServer) {}
