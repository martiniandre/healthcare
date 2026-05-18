package telemetry

import (
	"time"

	"github.com/google/uuid"
)

type Room struct {
	ID          uuid.UUID `db:"id"`
	Name        string    `db:"name"`
	Passcode    string    `db:"passcode"`
	Description string    `db:"description"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

type Bed struct {
	ID          uuid.UUID `db:"id"`
	RoomID      uuid.UUID `db:"room_id"`
	BedNumber   string    `db:"bed_number"`
	PatientName string    `db:"patient_name"`
	Age         int32     `db:"age"`
	Gender      string    `db:"gender"`
	Bpm         int32     `db:"bpm"`
	Spo2        int32     `db:"spo2"`
	Temperature float64   `db:"temperature"`
	Status      string    `db:"status"`
	Condition   string    `db:"condition"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}
