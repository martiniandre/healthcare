package schedule

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/healthcare/backend/internal/shared/apperrors"
)

func TestCreateUnavailability_ValidInputCreatesWindow(t *testing.T) {
	repositoryMock := &MockUnavailabilityRepository{}
	unavailabilityService := NewUnavailabilityService(repositoryMock)
	staffID := uuid.New()
	futureTime := time.Now().Add(48 * time.Hour).Truncate(time.Second)

	createdUnavailability, createErr := unavailabilityService.CreateUnavailability(context.Background(), CreateUnavailabilityInput{
		StaffID:  staffID,
		StartsAt: futureTime,
		EndsAt:   futureTime.Add(24 * time.Hour),
		Reason:   "Congresso médico",
	})

	if createErr != nil {
		t.Fatalf("unexpected error: %v", createErr)
	}
	if createdUnavailability == nil {
		t.Fatal("expected created unavailability, got nil")
	}
	if createdUnavailability.StaffID != staffID {
		t.Errorf("expected staff %s, got %s", staffID, createdUnavailability.StaffID)
	}
	if createdUnavailability.Reason != "Congresso médico" {
		t.Errorf("expected reason %s, got %s", "Congresso médico", createdUnavailability.Reason)
	}
}

func TestCreateUnavailability_RejectsMissingStaff(t *testing.T) {
	repositoryMock := &MockUnavailabilityRepository{}
	unavailabilityService := NewUnavailabilityService(repositoryMock)
	futureTime := time.Now().Add(48 * time.Hour)

	_, createErr := unavailabilityService.CreateUnavailability(context.Background(), CreateUnavailabilityInput{
		StartsAt: futureTime,
		EndsAt:   futureTime.Add(24 * time.Hour),
	})

	if createErr == nil {
		t.Fatal("expected validation error, got nil")
	}
	var validationError apperrors.AppError
	if !errors.As(createErr, &validationError) {
		t.Errorf("expected app error, got %v", createErr)
	}
}

func TestCreateUnavailability_RejectsPastStartsAt(t *testing.T) {
	repositoryMock := &MockUnavailabilityRepository{}
	unavailabilityService := NewUnavailabilityService(repositoryMock)
	pastTime := time.Now().Add(-time.Hour)

	_, createErr := unavailabilityService.CreateUnavailability(context.Background(), CreateUnavailabilityInput{
		StaffID:  uuid.New(),
		StartsAt: pastTime,
		EndsAt:   pastTime.Add(24 * time.Hour),
	})

	if createErr == nil {
		t.Fatal("expected validation error, got nil")
	}
	var validationError apperrors.AppError
	if !errors.As(createErr, &validationError) {
		t.Errorf("expected app error, got %v", createErr)
	}
}

func TestCreateUnavailability_RejectsEndsBeforeStarts(t *testing.T) {
	repositoryMock := &MockUnavailabilityRepository{}
	unavailabilityService := NewUnavailabilityService(repositoryMock)
	futureTime := time.Now().Add(48 * time.Hour)

	_, createErr := unavailabilityService.CreateUnavailability(context.Background(), CreateUnavailabilityInput{
		StaffID:  uuid.New(),
		StartsAt: futureTime,
		EndsAt:   futureTime.Add(-time.Hour),
	})

	if createErr == nil {
		t.Fatal("expected validation error, got nil")
	}
	var validationError apperrors.AppError
	if !errors.As(createErr, &validationError) {
		t.Errorf("expected app error, got %v", createErr)
	}
}

func TestCreateUnavailability_RejectsReasonOverMaximumLength(t *testing.T) {
	repositoryMock := &MockUnavailabilityRepository{}
	unavailabilityService := NewUnavailabilityService(repositoryMock)
	futureTime := time.Now().Add(48 * time.Hour)

	_, createErr := unavailabilityService.CreateUnavailability(context.Background(), CreateUnavailabilityInput{
		StaffID:  uuid.New(),
		StartsAt: futureTime,
		EndsAt:   futureTime.Add(24 * time.Hour),
		Reason:   string(make([]byte, 501)),
	})

	if createErr == nil {
		t.Fatal("expected validation error, got nil")
	}
	var validationError apperrors.AppError
	if !errors.As(createErr, &validationError) {
		t.Errorf("expected app error, got %v", createErr)
	}
}

func TestCreateUnavailability_ReturnsStaffNotFoundFromRepository(t *testing.T) {
	repositoryMock := &MockUnavailabilityRepository{
		CreateUnavailabilityFunc: func(ctx context.Context, unavailability *StaffUnavailability) (*StaffUnavailability, error) {
			return nil, apperrors.ErrEmployeeNotFound
		},
	}
	unavailabilityService := NewUnavailabilityService(repositoryMock)
	futureTime := time.Now().Add(48 * time.Hour)

	_, createErr := unavailabilityService.CreateUnavailability(context.Background(), CreateUnavailabilityInput{
		StaffID:  uuid.New(),
		StartsAt: futureTime,
		EndsAt:   futureTime.Add(24 * time.Hour),
	})

	if createErr == nil {
		t.Fatal("expected employee not found error, got nil")
	}
	if !errors.Is(createErr, apperrors.ErrEmployeeNotFound) {
		t.Errorf("expected employee not found, got %v", createErr)
	}
}

func TestCreateUnavailability_RejectsOverlappingWindow(t *testing.T) {
	repositoryMock := &MockUnavailabilityRepository{
		CreateUnavailabilityFunc: func(ctx context.Context, unavailability *StaffUnavailability) (*StaffUnavailability, error) {
			return nil, apperrors.ErrUnavailabilityConflict
		},
	}
	unavailabilityService := NewUnavailabilityService(repositoryMock)
	futureTime := time.Now().Add(48 * time.Hour)

	_, createErr := unavailabilityService.CreateUnavailability(context.Background(), CreateUnavailabilityInput{
		StaffID:  uuid.New(),
		StartsAt: futureTime,
		EndsAt:   futureTime.Add(24 * time.Hour),
	})

	if createErr == nil {
		t.Fatal("expected overlap conflict error, got nil")
	}
	if !errors.Is(createErr, apperrors.ErrUnavailabilityConflict) {
		t.Errorf("expected unavailability conflict, got %v", createErr)
	}
}

func TestDeleteUnavailability_DeletesExistingWindow(t *testing.T) {
	existingUnavailability := &StaffUnavailability{
		ID:       uuid.New(),
		StaffID:  uuid.New(),
		StartsAt: time.Now().Add(48 * time.Hour),
		EndsAt:   time.Now().Add(72 * time.Hour),
		Reason:   "Férias",
	}
	repositoryMock := &MockUnavailabilityRepository{
		DeleteUnavailabilityFunc: func(ctx context.Context, unavailabilityID uuid.UUID) (*StaffUnavailability, error) {
			return existingUnavailability, nil
		},
	}
	unavailabilityService := NewUnavailabilityService(repositoryMock)

	deletedUnavailability, deleteErr := unavailabilityService.DeleteUnavailability(context.Background(), existingUnavailability.ID)
	if deleteErr != nil {
		t.Fatalf("unexpected error: %v", deleteErr)
	}
	if deletedUnavailability.ID != existingUnavailability.ID {
		t.Errorf("expected deleted window %s, got %s", existingUnavailability.ID, deletedUnavailability.ID)
	}
}

func TestDeleteUnavailability_ReturnsNotFoundWhenMissing(t *testing.T) {
	repositoryMock := &MockUnavailabilityRepository{}
	unavailabilityService := NewUnavailabilityService(repositoryMock)

	_, deleteErr := unavailabilityService.DeleteUnavailability(context.Background(), uuid.New())
	if deleteErr == nil {
		t.Fatal("expected not found error, got nil")
	}
	if !errors.Is(deleteErr, apperrors.ErrUnavailabilityNotFound) {
		t.Errorf("expected unavailability not found, got %v", deleteErr)
	}
}

func TestListUnavailability_RejectsNilStaff(t *testing.T) {
	unavailabilityService := NewUnavailabilityService(&MockUnavailabilityRepository{})

	_, listErr := unavailabilityService.ListUnavailabilityByStaff(context.Background(), ListUnavailabilityInput{
		StaffID: uuid.Nil,
	})
	if listErr == nil {
		t.Fatal("expected validation error, got nil")
	}
	var validationError apperrors.AppError
	if !errors.As(listErr, &validationError) {
		t.Errorf("expected app error, got %v", listErr)
	}
}

func TestListUnavailability_RejectsRangeWithToBeforeFrom(t *testing.T) {
	unavailabilityService := NewUnavailabilityService(&MockUnavailabilityRepository{})
	baseTime := time.Now().Add(48 * time.Hour)

	_, listErr := unavailabilityService.ListUnavailabilityByStaff(context.Background(), ListUnavailabilityInput{
		StaffID: uuid.New(),
		From:    baseTime,
		To:      baseTime.Add(-time.Hour),
	})
	if listErr == nil {
		t.Fatal("expected validation error, got nil")
	}
	var validationError apperrors.AppError
	if !errors.As(listErr, &validationError) {
		t.Errorf("expected app error, got %v", listErr)
	}
}

func TestCreateAppointment_RejectsOverlappingUnavailabilityWindow(t *testing.T) {
	futureTime := time.Now().Add(48 * time.Hour).Truncate(time.Hour)
	appointmentRepositoryMock := &MockRepository{
		CreateAppointmentFunc: func(ctx context.Context, appointment *Appointment) (*Appointment, error) {
			unavailabilityOverlap := appointment.StartsAt.Before(futureTime.Add(24*time.Hour)) && appointment.EndsAt.After(futureTime)
			if unavailabilityOverlap {
				return nil, apperrors.ErrAppointmentConflict
			}
			appointment.CreatedAt = time.Now()
			return appointment, nil
		},
	}
	appointmentService := NewService(appointmentRepositoryMock, nil, nil)

	overlappingAppointment, createErr := appointmentService.CreateAppointment(context.Background(), CreateAppointmentInput{
		PatientFHIRID: "patient-123",
		StaffID:       uuid.New(),
		StartsAt:      futureTime.Add(2 * time.Hour),
		EndsAt:        futureTime.Add(2*time.Hour + 30*time.Minute),
		Reason:        "Consulta dentro de indisponibilidade",
	})

	if createErr == nil {
		t.Fatal("expected appointment conflict error, got nil")
	}
	if !errors.Is(createErr, apperrors.ErrAppointmentConflict) {
		t.Errorf("expected appointment conflict, got %v", createErr)
	}
	if overlappingAppointment != nil {
		t.Error("expected nil appointment on conflict")
	}
}

func TestCreateAppointment_AllowsWindowOutsideUnavailability(t *testing.T) {
	futureTime := time.Now().Add(48 * time.Hour).Truncate(time.Hour)
	appointmentRepositoryMock := &MockRepository{
		CreateAppointmentFunc: func(ctx context.Context, appointment *Appointment) (*Appointment, error) {
			unavailabilityOverlap := appointment.StartsAt.Before(futureTime.Add(24*time.Hour)) && appointment.EndsAt.After(futureTime)
			if unavailabilityOverlap {
				return nil, apperrors.ErrAppointmentConflict
			}
			appointment.CreatedAt = time.Now()
			return appointment, nil
		},
	}
	appointmentService := NewService(appointmentRepositoryMock, nil, nil)

	createdAppointment, createErr := appointmentService.CreateAppointment(context.Background(), CreateAppointmentInput{
		PatientFHIRID: "patient-123",
		StaffID:       uuid.New(),
		StartsAt:      futureTime.Add(48 * time.Hour),
		EndsAt:        futureTime.Add(48*time.Hour + 30*time.Minute),
		Reason:        "Consulta fora da indisponibilidade",
	})

	if createErr != nil {
		t.Fatalf("unexpected error: %v", createErr)
	}
	if createdAppointment == nil {
		t.Fatal("expected created appointment, got nil")
	}
}
