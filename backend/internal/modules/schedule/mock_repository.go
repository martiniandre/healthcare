package schedule

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/healthcare/backend/internal/shared/apperrors"
)

type MockRepository struct {
	CreateAppointmentFunc             func(ctx context.Context, appointment *Appointment) (*Appointment, error)
	GetAppointmentByIDFunc            func(ctx context.Context, appointmentID uuid.UUID) (*Appointment, error)
	CancelAppointmentFunc             func(ctx context.Context, appointmentID uuid.UUID) (*Appointment, error)
	RescheduleAppointmentFunc         func(ctx context.Context, appointmentID uuid.UUID, startsAt time.Time, endsAt time.Time) (*Appointment, error)
	ListAppointmentsByPatientFunc     func(ctx context.Context, patientFHIRID string) ([]*Appointment, error)
	ListAppointmentsByStaffOnDateFunc func(ctx context.Context, staffID uuid.UUID, date time.Time) ([]*Appointment, error)
	ListAppointmentsByStaffInRangeFunc func(ctx context.Context, staffID uuid.UUID, startDate time.Time, endDate time.Time) ([]*Appointment, error)
	ResolveActiveEmployeeIDByEmailFunc func(ctx context.Context, email string) (*uuid.UUID, error)
	FindIdempotencyKeyFunc            func(ctx context.Context, idempotencyKey string) (*IdempotencyKey, error)
	SaveIdempotencyKeyFunc            func(ctx context.Context, key *IdempotencyKey) error
}

func (mock *MockRepository) CreateAppointment(ctx context.Context, appointment *Appointment) (*Appointment, error) {
	if mock.CreateAppointmentFunc != nil {
		return mock.CreateAppointmentFunc(ctx, appointment)
	}
	appointment.CreatedAt = time.Now()
	return appointment, nil
}

func (mock *MockRepository) GetAppointmentByID(ctx context.Context, appointmentID uuid.UUID) (*Appointment, error) {
	if mock.GetAppointmentByIDFunc != nil {
		return mock.GetAppointmentByIDFunc(ctx, appointmentID)
	}
	return nil, apperrors.ErrAppointmentNotFound
}

func (mock *MockRepository) CancelAppointment(ctx context.Context, appointmentID uuid.UUID) (*Appointment, error) {
	if mock.CancelAppointmentFunc != nil {
		return mock.CancelAppointmentFunc(ctx, appointmentID)
	}
	return nil, apperrors.ErrAppointmentNotFound
}

func (mock *MockRepository) ListAppointmentsByPatient(ctx context.Context, patientFHIRID string) ([]*Appointment, error) {
	if mock.ListAppointmentsByPatientFunc != nil {
		return mock.ListAppointmentsByPatientFunc(ctx, patientFHIRID)
	}
	return []*Appointment{}, nil
}

func (mock *MockRepository) ListAppointmentsByStaffOnDate(ctx context.Context, staffID uuid.UUID, date time.Time) ([]*Appointment, error) {
	if mock.ListAppointmentsByStaffOnDateFunc != nil {
		return mock.ListAppointmentsByStaffOnDateFunc(ctx, staffID, date)
	}
	return []*Appointment{}, nil
}

func (mock *MockRepository) ListAppointmentsByStaffInRange(ctx context.Context, staffID uuid.UUID, startDate time.Time, endDate time.Time) ([]*Appointment, error) {
	if mock.ListAppointmentsByStaffInRangeFunc != nil {
		return mock.ListAppointmentsByStaffInRangeFunc(ctx, staffID, startDate, endDate)
	}
	return []*Appointment{}, nil
}

func (mock *MockRepository) RescheduleAppointment(ctx context.Context, appointmentID uuid.UUID, startsAt time.Time, endsAt time.Time) (*Appointment, error) {
	if mock.RescheduleAppointmentFunc != nil {
		return mock.RescheduleAppointmentFunc(ctx, appointmentID, startsAt, endsAt)
	}
	return nil, apperrors.ErrAppointmentNotFound
}

func (mock *MockRepository) ResolveActiveEmployeeIDByEmail(ctx context.Context, email string) (*uuid.UUID, error) {
	if mock.ResolveActiveEmployeeIDByEmailFunc != nil {
		return mock.ResolveActiveEmployeeIDByEmailFunc(ctx, email)
	}
	return nil, apperrors.ErrPermissionDenied
}

func (mock *MockRepository) FindIdempotencyKey(ctx context.Context, idempotencyKey string) (*IdempotencyKey, error) {
	if mock.FindIdempotencyKeyFunc != nil {
		return mock.FindIdempotencyKeyFunc(ctx, idempotencyKey)
	}
	return nil, nil
}

func (mock *MockRepository) SaveIdempotencyKey(ctx context.Context, key *IdempotencyKey) error {
	if mock.SaveIdempotencyKeyFunc != nil {
		return mock.SaveIdempotencyKeyFunc(ctx, key)
	}
	return nil
}
