package schedule

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/healthcare/backend/internal/shared/apperrors"
)

func mustMarshalAppointment(t *testing.T, appointment *Appointment) []byte {
	t.Helper()
	marshaledBody, marshalErr := json.Marshal(appointment)
	if marshalErr != nil {
		t.Fatalf("failed to marshal appointment: %v", marshalErr)
	}
	return marshaledBody
}

func TestCreateAppointment_ValidInputCreatesScheduledAppointment(t *testing.T) {
	repositoryMock := &MockRepository{}
	appointmentService := NewService(repositoryMock, nil)
	staffID := uuid.New()
	futureTime := time.Now().Add(48 * time.Hour).Truncate(time.Second)

	createdAppointment, createErr := appointmentService.CreateAppointment(context.Background(), CreateAppointmentInput{
		PatientFHIRID: "patient-123",
		StaffID:       staffID,
		StartsAt:      futureTime,
		EndsAt:        futureTime.Add(time.Hour),
		Reason:        "Follow-up",
	})

	if createErr != nil {
		t.Fatalf("unexpected error: %v", createErr)
	}
	if createdAppointment == nil {
		t.Fatal("expected created appointment, got nil")
	}
	if createdAppointment.Status != AppointmentStatusScheduled {
		t.Errorf("expected status %s, got %s", AppointmentStatusScheduled, createdAppointment.Status)
	}
	if createdAppointment.Version != 1 {
		t.Errorf("expected version 1, got %d", createdAppointment.Version)
	}
}

func TestCreateAppointment_RejectsPastStartsAt(t *testing.T) {
	repositoryMock := &MockRepository{}
	appointmentService := NewService(repositoryMock, nil)
	pastTime := time.Now().Add(-time.Hour)

	_, createErr := appointmentService.CreateAppointment(context.Background(), CreateAppointmentInput{
		PatientFHIRID: "patient-123",
		StaffID:       uuid.New(),
		StartsAt:      pastTime,
		EndsAt:        pastTime.Add(time.Hour),
	})

	if createErr == nil {
		t.Fatal("expected validation error, got nil")
	}
	var validationError apperrors.AppError
	if !errors.As(createErr, &validationError) {
		t.Errorf("expected app error, got %v", createErr)
	}
}

func TestCreateAppointment_RejectsEndsBeforeStarts(t *testing.T) {
	repositoryMock := &MockRepository{}
	appointmentService := NewService(repositoryMock, nil)
	futureTime := time.Now().Add(48 * time.Hour)

	_, createErr := appointmentService.CreateAppointment(context.Background(), CreateAppointmentInput{
		PatientFHIRID: "patient-123",
		StaffID:       uuid.New(),
		StartsAt:      futureTime,
		EndsAt:        futureTime.Add(-time.Hour),
	})

	if createErr == nil {
		t.Fatal("expected validation error, got nil")
	}
	var validationError apperrors.AppError
	if !errors.As(createErr, &validationError) {
		t.Errorf("expected app error, got %v", createErr)
	}
}

func TestCreateAppointment_PropagatesConflictFromRepository(t *testing.T) {
	repositoryMock := &MockRepository{
		CreateAppointmentFunc: func(ctx context.Context, appointment *Appointment) (*Appointment, error) {
			return nil, apperrors.ErrAppointmentConflict
		},
	}
	appointmentService := NewService(repositoryMock, nil)
	futureTime := time.Now().Add(48 * time.Hour)

	_, createErr := appointmentService.CreateAppointment(context.Background(), CreateAppointmentInput{
		PatientFHIRID: "patient-123",
		StaffID:       uuid.New(),
		StartsAt:      futureTime,
		EndsAt:        futureTime.Add(time.Hour),
	})

	if createErr == nil {
		t.Fatal("expected conflict error, got nil")
	}
	if !errors.Is(createErr, apperrors.ErrAppointmentConflict) {
		t.Errorf("expected appointment conflict, got %v", createErr)
	}
}

func TestCreateAppointment_ReplaysCachedResponseForSameIdempotencyKey(t *testing.T) {
	futureTime := time.Now().Add(48 * time.Hour).Truncate(time.Second)
	cachedAppointment := &Appointment{
		ID:            uuid.New(),
		PatientFHIRID: "patient-123",
		Status:        AppointmentStatusScheduled,
		Version:       1,
	}
	repositoryMock := &MockRepository{
		FindIdempotencyKeyFunc: func(ctx context.Context, idempotencyKey string) (*IdempotencyKey, error) {
			return &IdempotencyKey{
				IdempotencyKey: "key-abc",
				RequestHash:    "request-hash-value",
				ResponseStatus: 201,
				ResponseBody:   mustMarshalAppointment(t, cachedAppointment),
			}, nil
		},
	}
	appointmentService := NewService(repositoryMock, nil)

	createdAppointment, createErr := appointmentService.CreateAppointment(context.Background(), CreateAppointmentInput{
		PatientFHIRID:  "patient-123",
		StaffID:        uuid.New(),
		StartsAt:       futureTime,
		EndsAt:         futureTime.Add(time.Hour),
		IdempotencyKey: "key-abc",
		RequestHash:    "request-hash-value",
	})

	if createErr != nil {
		t.Fatalf("unexpected error: %v", createErr)
	}
	if createdAppointment.ID != cachedAppointment.ID {
		t.Errorf("expected cached appointment %s, got %s", cachedAppointment.ID, createdAppointment.ID)
	}
	if repositoryMock.CreateAppointmentFunc != nil {
		t.Error("repository CreateAppointment should not be called on idempotent replay")
	}
}

func TestCreateAppointment_RejectsIdempotencyKeyReusedWithDifferentPayload(t *testing.T) {
	futureTime := time.Now().Add(48 * time.Hour)
	repositoryMock := &MockRepository{
		FindIdempotencyKeyFunc: func(ctx context.Context, idempotencyKey string) (*IdempotencyKey, error) {
			return &IdempotencyKey{
				IdempotencyKey: "key-abc",
				RequestHash:    "different-hash",
				ResponseStatus: 201,
				ResponseBody:   []byte("{}"),
			}, nil
		},
	}
	appointmentService := NewService(repositoryMock, nil)

	_, createErr := appointmentService.CreateAppointment(context.Background(), CreateAppointmentInput{
		PatientFHIRID:  "patient-123",
		StaffID:        uuid.New(),
		StartsAt:       futureTime,
		EndsAt:         futureTime.Add(time.Hour),
		IdempotencyKey: "key-abc",
		RequestHash:    "request-hash-value",
	})

	if createErr == nil {
		t.Fatal("expected idempotency conflict error, got nil")
	}
	if !strings.Contains(createErr.Error(), "idempotency") {
		t.Errorf("expected idempotency conflict message, got %v", createErr)
	}
}

func TestCancelAppointment_CancelsScheduledAppointment(t *testing.T) {
	existingAppointment := &Appointment{
		ID:      uuid.New(),
		Status:  AppointmentStatusScheduled,
		Version: 2,
	}
	repositoryMock := &MockRepository{
		GetAppointmentByIDFunc: func(ctx context.Context, appointmentID uuid.UUID) (*Appointment, error) {
			return existingAppointment, nil
		},
		CancelAppointmentFunc: func(ctx context.Context, appointmentID uuid.UUID) (*Appointment, error) {
			cancelledAppointment := *existingAppointment
			cancelledAppointment.Status = AppointmentStatusCancelled
			cancelledAppointment.Version++
			return &cancelledAppointment, nil
		},
	}
	appointmentService := NewService(repositoryMock, nil)

	cancelledAppointment, cancelErr := appointmentService.CancelAppointment(context.Background(), existingAppointment.ID)
	if cancelErr != nil {
		t.Fatalf("unexpected error: %v", cancelErr)
	}
	if cancelledAppointment.Status != AppointmentStatusCancelled {
		t.Errorf("expected status %s, got %s", AppointmentStatusCancelled, cancelledAppointment.Status)
	}
	if cancelledAppointment.Version != 3 {
		t.Errorf("expected version 3, got %d", cancelledAppointment.Version)
	}
}

func TestCancelAppointment_IsIdempotentWhenAlreadyCancelled(t *testing.T) {
	alreadyCancelled := &Appointment{
		ID:      uuid.New(),
		Status:  AppointmentStatusCancelled,
		Version: 3,
	}
	repositoryMock := &MockRepository{
		GetAppointmentByIDFunc: func(ctx context.Context, appointmentID uuid.UUID) (*Appointment, error) {
			return alreadyCancelled, nil
		},
	}
	appointmentService := NewService(repositoryMock, nil)

	cancelledAppointment, cancelErr := appointmentService.CancelAppointment(context.Background(), alreadyCancelled.ID)
	if cancelErr != nil {
		t.Fatalf("unexpected error: %v", cancelErr)
	}
	if cancelledAppointment.Version != 3 {
		t.Errorf("expected no version change, got %d", cancelledAppointment.Version)
	}
	if repositoryMock.CancelAppointmentFunc != nil {
		t.Error("repository CancelAppointment should not be called for already cancelled appointment")
	}
}

func TestCancelAppointment_RejectsFinishedAppointment(t *testing.T) {
	finishedAppointment := &Appointment{
		ID:      uuid.New(),
		Status:  AppointmentStatusFinished,
		Version: 1,
	}
	repositoryMock := &MockRepository{
		GetAppointmentByIDFunc: func(ctx context.Context, appointmentID uuid.UUID) (*Appointment, error) {
			return finishedAppointment, nil
		},
	}
	appointmentService := NewService(repositoryMock, nil)

	_, cancelErr := appointmentService.CancelAppointment(context.Background(), finishedAppointment.ID)
	if cancelErr == nil {
		t.Fatal("expected invalid transition error, got nil")
	}
	if !errors.Is(cancelErr, apperrors.ErrAppointmentInvalidTransition) {
		t.Errorf("expected invalid transition, got %v", cancelErr)
	}
}

func TestCancelAppointment_ReturnsNotFoundWhenMissing(t *testing.T) {
	repositoryMock := &MockRepository{}
	appointmentService := NewService(repositoryMock, nil)

	_, cancelErr := appointmentService.CancelAppointment(context.Background(), uuid.New())
	if cancelErr == nil {
		t.Fatal("expected not found error, got nil")
	}
	if !errors.Is(cancelErr, apperrors.ErrAppointmentNotFound) {
		t.Errorf("expected appointment not found, got %v", cancelErr)
	}
}

func TestListAppointmentsByPatient_RejectsEmptyFilter(t *testing.T) {
	appointmentService := NewService(&MockRepository{}, nil)

	_, listErr := appointmentService.ListAppointmentsByPatient(context.Background(), "")
	if listErr == nil {
		t.Fatal("expected validation error, got nil")
	}
	var validationError apperrors.AppError
	if !errors.As(listErr, &validationError) {
		t.Errorf("expected app error, got %v", listErr)
	}
}

func TestListAppointmentsByStaffOnDate_RejectsNilStaff(t *testing.T) {
	appointmentService := NewService(&MockRepository{}, nil)

	_, listErr := appointmentService.ListAppointmentsByStaffOnDate(context.Background(), uuid.Nil, time.Now())
	if listErr == nil {
		t.Fatal("expected validation error, got nil")
	}
	var validationError apperrors.AppError
	if !errors.As(listErr, &validationError) {
		t.Errorf("expected app error, got %v", listErr)
	}
}

func TestComputeRequestHash_IsDeterministic(t *testing.T) {
	payload := CreateAppointmentRequest{
		PatientFhirID: "patient-123",
		StaffID:       uuid.New().String(),
		StartsAt:      "2026-08-03T09:00:00Z",
		EndsAt:        "2026-08-03T10:00:00Z",
		Reason:        "Follow-up",
	}

	firstHash := computeRequestHash(payload)
	secondHash := computeRequestHash(payload)
	if firstHash != secondHash {
		t.Errorf("expected deterministic hash, got %s and %s", firstHash, secondHash)
	}
}
