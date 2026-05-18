package telemetry

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	GetRooms(ctx context.Context) ([]*Room, error)
	GetRoomByID(ctx context.Context, roomID uuid.UUID) (*Room, error)
	GetBedsByRoomID(ctx context.Context, roomID uuid.UUID) ([]*Bed, error)
	GetBedByID(ctx context.Context, bedID uuid.UUID) (*Bed, error)
	UpdateBedCondition(ctx context.Context, bed *Bed) error
}

type repository struct {
	dbPool *pgxpool.Pool
}

func NewRepository(dbPool *pgxpool.Pool) Repository {
	return &repository{dbPool: dbPool}
}

func (telemetryRepository *repository) GetRooms(ctx context.Context) ([]*Room, error) {
	query := `SELECT id, name, passcode, description, created_at, updated_at 
			  FROM telemetry_rooms ORDER BY name ASC`

	rows, err := telemetryRepository.dbPool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rooms := make([]*Room, 0)
	for rows.Next() {
		room := &Room{}
		err := rows.Scan(
			&room.ID, &room.Name, &room.Passcode, &room.Description,
			&room.CreatedAt, &room.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		rooms = append(rooms, room)
	}

	return rooms, nil
}

func (telemetryRepository *repository) GetRoomByID(ctx context.Context, roomID uuid.UUID) (*Room, error) {
	query := `SELECT id, name, passcode, description, created_at, updated_at 
			  FROM telemetry_rooms WHERE id = $1`

	room := &Room{}
	err := telemetryRepository.dbPool.QueryRow(ctx, query, roomID).Scan(
		&room.ID, &room.Name, &room.Passcode, &room.Description,
		&room.CreatedAt, &room.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return room, nil
}

func (telemetryRepository *repository) GetBedsByRoomID(ctx context.Context, roomID uuid.UUID) ([]*Bed, error) {
	query := `SELECT id, room_id, bed_number, patient_name, age, gender, bpm, spo2, temperature, status, condition, created_at, updated_at 
			  FROM telemetry_beds WHERE room_id = $1 ORDER BY bed_number ASC`

	rows, err := telemetryRepository.dbPool.Query(ctx, query, roomID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	beds := make([]*Bed, 0)
	for rows.Next() {
		bed := &Bed{}
		err := rows.Scan(
			&bed.ID, &bed.RoomID, &bed.BedNumber, &bed.PatientName,
			&bed.Age, &bed.Gender, &bed.Bpm, &bed.Spo2, &bed.Temperature,
			&bed.Status, &bed.Condition, &bed.CreatedAt, &bed.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		beds = append(beds, bed)
	}

	return beds, nil
}

func (telemetryRepository *repository) GetBedByID(ctx context.Context, bedID uuid.UUID) (*Bed, error) {
	query := `SELECT id, room_id, bed_number, patient_name, age, gender, bpm, spo2, temperature, status, condition, created_at, updated_at 
			  FROM telemetry_beds WHERE id = $1`

	bed := &Bed{}
	err := telemetryRepository.dbPool.QueryRow(ctx, query, bedID).Scan(
		&bed.ID, &bed.RoomID, &bed.BedNumber, &bed.PatientName,
		&bed.Age, &bed.Gender, &bed.Bpm, &bed.Spo2, &bed.Temperature,
		&bed.Status, &bed.Condition, &bed.CreatedAt, &bed.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return bed, nil
}

func (telemetryRepository *repository) UpdateBedCondition(ctx context.Context, bed *Bed) error {
	query := `UPDATE telemetry_beds 
			  SET bpm = $1, spo2 = $2, temperature = $3, status = $4, condition = $5, updated_at = NOW() 
			  WHERE id = $6`

	_, err := telemetryRepository.dbPool.Exec(ctx, query,
		bed.Bpm, bed.Spo2, bed.Temperature, bed.Status, bed.Condition, bed.ID,
	)
	return err
}
