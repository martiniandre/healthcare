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
	ID          uuid.UUID `db:"id" json:"id"`
	RoomID      uuid.UUID `db:"room_id" json:"roomId"`
	BedNumber   string    `db:"bed_number" json:"bedNumber"`
	PatientName string    `db:"patient_name" json:"patientName"`
	Age         int32     `db:"age" json:"age"`
	Gender      string    `db:"gender" json:"gender"`
	Bpm         int32     `db:"bpm" json:"bpm"`
	Spo2        int32     `db:"spo2" json:"spo2"`
	Temperature float64   `db:"temperature" json:"temperature"`
	Status      string    `db:"status" json:"status"`
	Condition   string    `db:"condition" json:"condition"`
	CreatedAt   time.Time `db:"created_at" json:"-"`
	UpdatedAt   time.Time `db:"updated_at" json:"-"`
}
