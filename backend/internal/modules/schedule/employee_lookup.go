package schedule

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/healthcare/backend/internal/shared/apperrors"
)

func resolveActiveEmployeeIDByEmail(ctx context.Context, dbPool *pgxpool.Pool, email string) (*uuid.UUID, error) {
	if email == "" {
		return nil, apperrors.ErrPermissionDenied
	}
	var employeeID uuid.UUID
	lookupQuery := `SELECT id FROM employees WHERE email = $1 AND is_active = true`
	scanErr := dbPool.QueryRow(ctx, lookupQuery, email).Scan(&employeeID)
	if scanErr != nil {
		if errors.Is(scanErr, pgx.ErrNoRows) {
			return nil, apperrors.ErrPermissionDenied
		}
		return nil, scanErr
	}
	return &employeeID, nil
}
