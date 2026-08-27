package schedule

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/healthcare/backend/internal/shared/apperrors"
)

type Repository interface {
	CreateAppointment(ctx context.Context, appointment *Appointment) (*Appointment, error)
	CancelAppointment(ctx context.Context, appointmentID uuid.UUID) (*Appointment, error)
	RescheduleAppointment(ctx context.Context, appointmentID uuid.UUID, startsAt time.Time, endsAt time.Time) (*Appointment, error)
	GetAppointmentByID(ctx context.Context, appointmentID uuid.UUID) (*Appointment, error)
	ListAppointmentsByPatient(ctx context.Context, patientFHIRID string) ([]*Appointment, error)
	ListAppointmentsByStaffOnDate(ctx context.Context, staffID uuid.UUID, date time.Time) ([]*Appointment, error)
	ListAppointmentsByStaffInRange(ctx context.Context, staffID uuid.UUID, startDate time.Time, endDate time.Time) ([]*Appointment, error)
	ResolveActiveEmployeeIDByEmail(ctx context.Context, email string) (*uuid.UUID, error)
	ResolvePatientFHIRIDByUserID(ctx context.Context, userID string) (string, error)
	FindIdempotencyKey(ctx context.Context, idempotencyKey string) (*IdempotencyKey, error)
	SaveIdempotencyKey(ctx context.Context, idempotencyKey *IdempotencyKey) error
}

type repository struct {
	dbPool *pgxpool.Pool
}

func NewRepository(dbPool *pgxpool.Pool) Repository {
	return &repository{dbPool: dbPool}
}

func (appointmentRepository *repository) CreateAppointment(ctx context.Context, appointment *Appointment) (*Appointment, error) {
	transaction, beginErr := appointmentRepository.dbPool.Begin(ctx)
	if beginErr != nil {
		return nil, beginErr
	}
	defer transaction.Rollback(ctx)

	var activeEmployeeID uuid.UUID
	employeeCheckQuery := `SELECT id FROM employees WHERE id = $1 AND is_active = true`
	employeeCheckErr := transaction.QueryRow(ctx, employeeCheckQuery, appointment.StaffID).Scan(&activeEmployeeID)
	if employeeCheckErr != nil {
		if errors.Is(employeeCheckErr, pgx.ErrNoRows) {
			return nil, apperrors.ErrEmployeeNotFound
		}
		return nil, employeeCheckErr
	}

	overlapQuery := `SELECT id FROM appointments
		WHERE staff_id = $1 AND status IN ('scheduled', 'confirmed')
		AND starts_at < $3 AND ends_at > $2
		FOR UPDATE`
	overlappingRows, overlapQueryErr := transaction.Query(ctx, overlapQuery, appointment.StaffID, appointment.StartsAt, appointment.EndsAt)
	if overlapQueryErr != nil {
		return nil, overlapQueryErr
	}
	hasOverlap := overlappingRows.Next()
	overlappingRows.Close()
	if overlapQueryErr := overlappingRows.Err(); overlapQueryErr != nil {
		return nil, overlapQueryErr
	}
	if hasOverlap {
		return nil, apperrors.ErrAppointmentConflict
	}

	unavailabilityOverlapQuery := `SELECT id FROM staff_unavailability
		WHERE staff_id = $1 AND starts_at < $3 AND ends_at > $2
		FOR UPDATE`
	overlappingUnavailabilityRows, unavailabilityQueryErr := transaction.Query(ctx, unavailabilityOverlapQuery, appointment.StaffID, appointment.StartsAt, appointment.EndsAt)
	if unavailabilityQueryErr != nil {
		return nil, unavailabilityQueryErr
	}
	hasUnavailabilityOverlap := overlappingUnavailabilityRows.Next()
	overlappingUnavailabilityRows.Close()
	if unavailabilityQueryErr := overlappingUnavailabilityRows.Err(); unavailabilityQueryErr != nil {
		return nil, unavailabilityQueryErr
	}
	if hasUnavailabilityOverlap {
		return nil, apperrors.ErrAppointmentConflict
	}

	insertQuery := `INSERT INTO appointments (id, patient_fhir_id, staff_id, starts_at, ends_at, status, reason, version, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, insertErr := transaction.Exec(ctx, insertQuery,
		appointment.ID, appointment.PatientFHIRID, appointment.StaffID,
		appointment.StartsAt, appointment.EndsAt, appointment.Status,
		appointment.Reason, appointment.Version, appointment.CreatedBy,
	)
	if insertErr != nil {
		var postgresError *pgconn.PgError
		if errors.As(insertErr, &postgresError) && postgresError.Code == "23505" {
			return nil, apperrors.ErrAppointmentConflict
		}
		return nil, insertErr
	}

	if commitErr := transaction.Commit(ctx); commitErr != nil {
		return nil, commitErr
	}

	appointment.CreatedAt = time.Now()
	appointment.UpdatedAt = appointment.CreatedAt
	return appointment, nil
}

func (appointmentRepository *repository) CancelAppointment(ctx context.Context, appointmentID uuid.UUID) (*Appointment, error) {
	updateQuery := `UPDATE appointments
		SET status = 'cancelled', version = version + 1, updated_at = NOW()
		WHERE id = $1
		RETURNING id, patient_fhir_id, staff_id, starts_at, ends_at, status, reason, version, created_by, created_at, updated_at`
	var appointment Appointment
	scanErr := appointmentRepository.dbPool.QueryRow(ctx, updateQuery, appointmentID).Scan(
		&appointment.ID, &appointment.PatientFHIRID, &appointment.StaffID,
		&appointment.StartsAt, &appointment.EndsAt, &appointment.Status,
		&appointment.Reason, &appointment.Version, &appointment.CreatedBy,
		&appointment.CreatedAt, &appointment.UpdatedAt,
	)
	if scanErr != nil {
		if errors.Is(scanErr, pgx.ErrNoRows) {
			return nil, apperrors.ErrAppointmentNotFound
		}
		return nil, scanErr
	}
	return &appointment, nil
}

func (appointmentRepository *repository) GetAppointmentByID(ctx context.Context, appointmentID uuid.UUID) (*Appointment, error) {
	getQuery := `SELECT id, patient_fhir_id, staff_id, starts_at, ends_at, status, reason, version, created_by, created_at, updated_at
		FROM appointments WHERE id = $1`
	var appointment Appointment
	scanErr := appointmentRepository.dbPool.QueryRow(ctx, getQuery, appointmentID).Scan(
		&appointment.ID, &appointment.PatientFHIRID, &appointment.StaffID,
		&appointment.StartsAt, &appointment.EndsAt, &appointment.Status,
		&appointment.Reason, &appointment.Version, &appointment.CreatedBy,
		&appointment.CreatedAt, &appointment.UpdatedAt,
	)
	if scanErr != nil {
		if errors.Is(scanErr, pgx.ErrNoRows) {
			return nil, apperrors.ErrAppointmentNotFound
		}
		return nil, scanErr
	}
	return &appointment, nil
}

func (appointmentRepository *repository) ListAppointmentsByPatient(ctx context.Context, patientFHIRID string) ([]*Appointment, error) {
	listQuery := `SELECT id, patient_fhir_id, staff_id, starts_at, ends_at, status, reason, version, created_by, created_at, updated_at
		FROM appointments WHERE patient_fhir_id = $1 ORDER BY starts_at DESC`
	return appointmentRepository.queryAppointments(ctx, listQuery, patientFHIRID)
}

func (appointmentRepository *repository) ListAppointmentsByStaffOnDate(ctx context.Context, staffID uuid.UUID, date time.Time) ([]*Appointment, error) {
	listQuery := `SELECT id, patient_fhir_id, staff_id, starts_at, ends_at, status, reason, version, created_by, created_at, updated_at
		FROM appointments WHERE staff_id = $1 AND starts_at::date = $2::date ORDER BY starts_at ASC`
	return appointmentRepository.queryAppointments(ctx, listQuery, staffID, date)
}

func (appointmentRepository *repository) ListAppointmentsByStaffInRange(ctx context.Context, staffID uuid.UUID, startDate time.Time, endDate time.Time) ([]*Appointment, error) {
	listQuery := `SELECT id, patient_fhir_id, staff_id, starts_at, ends_at, status, reason, version, created_by, created_at, updated_at
		FROM appointments
		WHERE staff_id = $1 AND starts_at::date >= $2::date AND starts_at::date <= $3::date
		ORDER BY starts_at ASC`
	return appointmentRepository.queryAppointments(ctx, listQuery, staffID, startDate, endDate)
}

func (appointmentRepository *repository) RescheduleAppointment(ctx context.Context, appointmentID uuid.UUID, startsAt time.Time, endsAt time.Time) (*Appointment, error) {
	transaction, beginErr := appointmentRepository.dbPool.Begin(ctx)
	if beginErr != nil {
		return nil, beginErr
	}
	defer transaction.Rollback(ctx)

	var currentStatus string
	var currentStaffID uuid.UUID
	statusQuery := `SELECT status, staff_id FROM appointments WHERE id = $1 FOR UPDATE`
	statusScanErr := transaction.QueryRow(ctx, statusQuery, appointmentID).Scan(&currentStatus, &currentStaffID)
	if statusScanErr != nil {
		if errors.Is(statusScanErr, pgx.ErrNoRows) {
			return nil, apperrors.ErrAppointmentNotFound
		}
		return nil, statusScanErr
	}
	if currentStatus != string(AppointmentStatusScheduled) && currentStatus != string(AppointmentStatusConfirmed) {
		return nil, apperrors.ErrAppointmentInvalidTransition
	}

	overlapQuery := `SELECT id FROM appointments
		WHERE id <> $1 AND staff_id = $2 AND status IN ('scheduled', 'confirmed')
		AND starts_at < $4 AND ends_at > $3
		FOR UPDATE`
	overlappingRows, overlapQueryErr := transaction.Query(ctx, overlapQuery, appointmentID, currentStaffID, startsAt, endsAt)
	if overlapQueryErr != nil {
		return nil, overlapQueryErr
	}
	hasOverlap := overlappingRows.Next()
	overlappingRows.Close()
	if overlappingRowsErr := overlappingRows.Err(); overlappingRowsErr != nil {
		return nil, overlappingRowsErr
	}
	if hasOverlap {
		return nil, apperrors.ErrAppointmentConflict
	}

	unavailabilityOverlapQuery := `SELECT id FROM staff_unavailability
		WHERE staff_id = $1 AND starts_at < $3 AND ends_at > $2
		FOR UPDATE`
	overlappingUnavailabilityRows, unavailabilityQueryErr := transaction.Query(ctx, unavailabilityOverlapQuery, currentStaffID, startsAt, endsAt)
	if unavailabilityQueryErr != nil {
		return nil, unavailabilityQueryErr
	}
	hasUnavailabilityOverlap := overlappingUnavailabilityRows.Next()
	overlappingUnavailabilityRows.Close()
	if overlappingUnavailabilityRowsErr := overlappingUnavailabilityRows.Err(); overlappingUnavailabilityRowsErr != nil {
		return nil, overlappingUnavailabilityRowsErr
	}
	if hasUnavailabilityOverlap {
		return nil, apperrors.ErrAppointmentConflict
	}

	updateQuery := `UPDATE appointments
		SET starts_at = $2, ends_at = $3, version = version + 1, updated_at = NOW()
		WHERE id = $1
		RETURNING id, patient_fhir_id, staff_id, starts_at, ends_at, status, reason, version, created_by, created_at, updated_at`
	var rescheduledAppointment Appointment
	scanErr := transaction.QueryRow(ctx, updateQuery, appointmentID, startsAt, endsAt).Scan(
		&rescheduledAppointment.ID, &rescheduledAppointment.PatientFHIRID, &rescheduledAppointment.StaffID,
		&rescheduledAppointment.StartsAt, &rescheduledAppointment.EndsAt, &rescheduledAppointment.Status,
		&rescheduledAppointment.Reason, &rescheduledAppointment.Version, &rescheduledAppointment.CreatedBy,
		&rescheduledAppointment.CreatedAt, &rescheduledAppointment.UpdatedAt,
	)
	if scanErr != nil {
		return nil, scanErr
	}

	if commitErr := transaction.Commit(ctx); commitErr != nil {
		return nil, commitErr
	}
	return &rescheduledAppointment, nil
}

func (appointmentRepository *repository) ResolveActiveEmployeeIDByEmail(ctx context.Context, email string) (*uuid.UUID, error) {
	return resolveActiveEmployeeIDByEmail(ctx, appointmentRepository.dbPool, email)
}

func (appointmentRepository *repository) ResolvePatientFHIRIDByUserID(ctx context.Context, userID string) (string, error) {
	query := `SELECT patient_fhir_id FROM patient_user_links WHERE user_id = $1`
	var patientFHIRID string
	scanErr := appointmentRepository.dbPool.QueryRow(ctx, query, userID).Scan(&patientFHIRID)
	if scanErr != nil {
		return "", apperrors.ErrPatientNotFound
	}
	return patientFHIRID, nil
}

func (appointmentRepository *repository) queryAppointments(ctx context.Context, query string, queryArgs ...interface{}) ([]*Appointment, error) {
	rows, queryErr := appointmentRepository.dbPool.Query(ctx, query, queryArgs...)
	if queryErr != nil {
		return nil, queryErr
	}
	defer rows.Close()

	appointments := make([]*Appointment, 0)
	for rows.Next() {
		var appointment Appointment
		scanErr := rows.Scan(
			&appointment.ID, &appointment.PatientFHIRID, &appointment.StaffID,
			&appointment.StartsAt, &appointment.EndsAt, &appointment.Status,
			&appointment.Reason, &appointment.Version, &appointment.CreatedBy,
			&appointment.CreatedAt, &appointment.UpdatedAt,
		)
		if scanErr != nil {
			return nil, scanErr
		}
		appointments = append(appointments, &appointment)
	}
	return appointments, nil
}

func (appointmentRepository *repository) FindIdempotencyKey(ctx context.Context, idempotencyKey string) (*IdempotencyKey, error) {
	getQuery := `SELECT id, idempotency_key, request_hash, response_status, response_body, created_at
		FROM idempotency_keys WHERE idempotency_key = $1`
	var storedKey IdempotencyKey
	scanErr := appointmentRepository.dbPool.QueryRow(ctx, getQuery, idempotencyKey).Scan(
		&storedKey.ID, &storedKey.IdempotencyKey, &storedKey.RequestHash,
		&storedKey.ResponseStatus, &storedKey.ResponseBody, &storedKey.CreatedAt,
	)
	if scanErr != nil {
		if errors.Is(scanErr, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, scanErr
	}
	return &storedKey, nil
}

func (appointmentRepository *repository) SaveIdempotencyKey(ctx context.Context, idempotencyKey *IdempotencyKey) error {
	insertQuery := `INSERT INTO idempotency_keys (id, idempotency_key, request_hash, response_status, response_body)
		VALUES ($1, $2, $3, $4, $5)`
	_, insertErr := appointmentRepository.dbPool.Exec(ctx, insertQuery,
		idempotencyKey.ID, idempotencyKey.IdempotencyKey, idempotencyKey.RequestHash,
		idempotencyKey.ResponseStatus, idempotencyKey.ResponseBody,
	)
	if insertErr != nil {
		return insertErr
	}
	return nil
}
