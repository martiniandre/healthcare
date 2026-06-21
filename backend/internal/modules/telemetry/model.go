package telemetry

import (
	"time"

	"github.com/google/uuid"
)

type Room struct {
	ID          uuid.UUID `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	Passcode    string    `db:"passcode" json:"-"`
	Description string    `db:"description" json:"description"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
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
