package schedule

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/healthcare/backend/internal/shared/apperrors"
)

type UnavailabilityRepository interface {
	CreateUnavailability(ctx context.Context, unavailability *StaffUnavailability) (*StaffUnavailability, error)
	ListUnavailabilityByStaff(ctx context.Context, staffID uuid.UUID, from time.Time, to time.Time) ([]*StaffUnavailability, error)
	DeleteUnavailability(ctx context.Context, unavailabilityID uuid.UUID) (*StaffUnavailability, error)
}

type unavailabilityRepository struct {
	dbPool *pgxpool.Pool
}

func NewUnavailabilityRepository(dbPool *pgxpool.Pool) UnavailabilityRepository {
	return &unavailabilityRepository{dbPool: dbPool}
}

func (unavailabilityRepository *unavailabilityRepository) CreateUnavailability(ctx context.Context, unavailability *StaffUnavailability) (*StaffUnavailability, error) {
	transaction, beginErr := unavailabilityRepository.dbPool.Begin(ctx)
	if beginErr != nil {
		return nil, beginErr
	}
	defer transaction.Rollback(ctx)

	var activeEmployeeID uuid.UUID
	employeeCheckQuery := `SELECT id FROM employees WHERE id = $1 AND is_active = true`
	employeeCheckErr := transaction.QueryRow(ctx, employeeCheckQuery, unavailability.StaffID).Scan(&activeEmployeeID)
	if employeeCheckErr != nil {
		if errors.Is(employeeCheckErr, pgx.ErrNoRows) {
			return nil, apperrors.ErrEmployeeNotFound
		}
		return nil, employeeCheckErr
	}

	overlapQuery := `SELECT id FROM staff_unavailability
		WHERE staff_id = $1 AND starts_at < $3 AND ends_at > $2
		FOR UPDATE`
	overlappingRows, overlapQueryErr := transaction.Query(ctx, overlapQuery, unavailability.StaffID, unavailability.StartsAt, unavailability.EndsAt)
	if overlapQueryErr != nil {
		return nil, overlapQueryErr
	}
	hasOverlap := overlappingRows.Next()
	overlappingRows.Close()
	if overlapQueryErr := overlappingRows.Err(); overlapQueryErr != nil {
		return nil, overlapQueryErr
	}
	if hasOverlap {
		return nil, apperrors.ErrUnavailabilityConflict
	}

	insertQuery := `INSERT INTO staff_unavailability (id, staff_id, starts_at, ends_at, reason, created_by)
		VALUES ($1, $2, $3, $4, $5, $6)`
	_, insertErr := transaction.Exec(ctx, insertQuery,
		unavailability.ID, unavailability.StaffID,
		unavailability.StartsAt, unavailability.EndsAt,
		unavailability.Reason, unavailability.CreatedBy,
	)
	if insertErr != nil {
		return nil, insertErr
	}

	if commitErr := transaction.Commit(ctx); commitErr != nil {
		return nil, commitErr
	}

	unavailability.CreatedAt = time.Now()
	unavailability.UpdatedAt = unavailability.CreatedAt
	return unavailability, nil
}

func (unavailabilityRepository *unavailabilityRepository) ListUnavailabilityByStaff(ctx context.Context, staffID uuid.UUID, from time.Time, to time.Time) ([]*StaffUnavailability, error) {
	listQuery := `SELECT id, staff_id, starts_at, ends_at, reason, created_by, created_at, updated_at
		FROM staff_unavailability
		WHERE staff_id = $1
		AND ($2::timestamptz IS NULL OR starts_at >= $2)
		AND ($3::timestamptz IS NULL OR ends_at <= $3)
		ORDER BY starts_at ASC`
	rows, queryErr := unavailabilityRepository.dbPool.Query(ctx, listQuery, staffID, nullableTime(from), nullableTime(to))
	if queryErr != nil {
		return nil, queryErr
	}
	defer rows.Close()

	unavailabilityWindows := make([]*StaffUnavailability, 0)
	for rows.Next() {
		var unavailabilityWindow StaffUnavailability
		scanErr := rows.Scan(
			&unavailabilityWindow.ID, &unavailabilityWindow.StaffID,
			&unavailabilityWindow.StartsAt, &unavailabilityWindow.EndsAt,
			&unavailabilityWindow.Reason, &unavailabilityWindow.CreatedBy,
			&unavailabilityWindow.CreatedAt, &unavailabilityWindow.UpdatedAt,
		)
		if scanErr != nil {
			return nil, scanErr
		}
		unavailabilityWindows = append(unavailabilityWindows, &unavailabilityWindow)
	}
	return unavailabilityWindows, nil
}

func (unavailabilityRepository *unavailabilityRepository) DeleteUnavailability(ctx context.Context, unavailabilityID uuid.UUID) (*StaffUnavailability, error) {
	deleteQuery := `DELETE FROM staff_unavailability
		WHERE id = $1
		RETURNING id, staff_id, starts_at, ends_at, reason, created_by, created_at, updated_at`
	var deletedUnavailability StaffUnavailability
	scanErr := unavailabilityRepository.dbPool.QueryRow(ctx, deleteQuery, unavailabilityID).Scan(
		&deletedUnavailability.ID, &deletedUnavailability.StaffID,
		&deletedUnavailability.StartsAt, &deletedUnavailability.EndsAt,
		&deletedUnavailability.Reason, &deletedUnavailability.CreatedBy,
		&deletedUnavailability.CreatedAt, &deletedUnavailability.UpdatedAt,
	)
	if scanErr != nil {
		if errors.Is(scanErr, pgx.ErrNoRows) {
			return nil, apperrors.ErrUnavailabilityNotFound
		}
		return nil, scanErr
	}
	return &deletedUnavailability, nil
}

func nullableTime(value time.Time) *time.Time {
	if value.IsZero() {
		return nil
	}
	return &value
}
