package schedule

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/healthcare/backend/internal/shared/apperrors"
)

type MockUnavailabilityRepository struct {
	CreateUnavailabilityFunc         func(ctx context.Context, unavailability *StaffUnavailability) (*StaffUnavailability, error)
	ListUnavailabilityByStaffFunc    func(ctx context.Context, staffID uuid.UUID, from time.Time, to time.Time) ([]*StaffUnavailability, error)
	GetUnavailabilityByIDFunc        func(ctx context.Context, unavailabilityID uuid.UUID) (*StaffUnavailability, error)
	DeleteUnavailabilityFunc         func(ctx context.Context, unavailabilityID uuid.UUID) (*StaffUnavailability, error)
	ResolveActiveEmployeeIDByEmailFunc func(ctx context.Context, email string) (*uuid.UUID, error)
}

func (mock *MockUnavailabilityRepository) CreateUnavailability(ctx context.Context, unavailability *StaffUnavailability) (*StaffUnavailability, error) {
	if mock.CreateUnavailabilityFunc != nil {
		return mock.CreateUnavailabilityFunc(ctx, unavailability)
	}
	unavailability.CreatedAt = time.Now()
	unavailability.UpdatedAt = unavailability.CreatedAt
	return unavailability, nil
}

func (mock *MockUnavailabilityRepository) ListUnavailabilityByStaff(ctx context.Context, staffID uuid.UUID, from time.Time, to time.Time) ([]*StaffUnavailability, error) {
	if mock.ListUnavailabilityByStaffFunc != nil {
		return mock.ListUnavailabilityByStaffFunc(ctx, staffID, from, to)
	}
	return []*StaffUnavailability{}, nil
}

func (mock *MockUnavailabilityRepository) GetUnavailabilityByID(ctx context.Context, unavailabilityID uuid.UUID) (*StaffUnavailability, error) {
	if mock.GetUnavailabilityByIDFunc != nil {
		return mock.GetUnavailabilityByIDFunc(ctx, unavailabilityID)
	}
	return nil, apperrors.ErrUnavailabilityNotFound
}

func (mock *MockUnavailabilityRepository) ResolveActiveEmployeeIDByEmail(ctx context.Context, email string) (*uuid.UUID, error) {
	if mock.ResolveActiveEmployeeIDByEmailFunc != nil {
		return mock.ResolveActiveEmployeeIDByEmailFunc(ctx, email)
	}
	return nil, apperrors.ErrPermissionDenied
}

func (mock *MockUnavailabilityRepository) DeleteUnavailability(ctx context.Context, unavailabilityID uuid.UUID) (*StaffUnavailability, error) {
	if mock.DeleteUnavailabilityFunc != nil {
		return mock.DeleteUnavailabilityFunc(ctx, unavailabilityID)
	}
	return nil, apperrors.ErrUnavailabilityNotFound
}
