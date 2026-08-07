package telemetry

type UnlockRoomInput struct {
	RoomID   string
	Passcode string
}

type GetBedsInput struct {
	RoomID string
}

type UpdateBedConditionInput struct {
	BedID       string
	Bpm         int32
	Spo2        int32
	Temperature float64
	Status      string
	Condition   string
}
