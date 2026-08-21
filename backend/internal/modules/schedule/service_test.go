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

func futureAlignedSlotStart(slotMinute int) time.Time {
	currentTime := time.Now()
	alignedStart := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), currentTime.Hour(), 0, 0, 0, currentTime.Location()).Add(48 * time.Hour)
	return alignedStart.Add(time.Duration(slotMinute) * time.Minute)
}

func assertFieldViolation(t *testing.T, validationErr error, fieldName string) {
	t.Helper()
	var validationError apperrors.AppError
	if !errors.As(validationErr, &validationError) {
		t.Fatalf("expected app error, got %v", validationErr)
	}
	if !strings.Contains(validationError.Message, fieldName) {
		t.Errorf("expected field violation %s in message, got %v", fieldName, validationError.Message)
	}
}

func TestCreateAppointment_ValidInputCreatesScheduledAppointment(t *testing.T) {
	repositoryMock := &MockRepository{}
	appointmentService := NewService(repositoryMock, nil, nil)
	staffID := uuid.New()
	slotStart := futureAlignedSlotStart(0)

	createdAppointment, createErr := appointmentService.CreateAppointment(context.Background(), CreateAppointmentInput{
		PatientFHIRID: "patient-123",
		StaffID:       staffID,
		StartsAt:      slotStart,
		EndsAt:        slotStart.Add(30 * time.Minute),
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

func TestCreateAppointment_AcceptsFortyFiveMinuteSlot(t *testing.T) {
	repositoryMock := &MockRepository{}
	appointmentService := NewService(repositoryMock, nil, nil)
	slotStart := futureAlignedSlotStart(30)

	createdAppointment, createErr := appointmentService.CreateAppointment(context.Background(), CreateAppointmentInput{
		PatientFHIRID: "patient-123",
		StaffID:       uuid.New(),
		StartsAt:      slotStart,
		EndsAt:        slotStart.Add(45 * time.Minute),
	})

	if createErr != nil {
		t.Fatalf("unexpected error: %v", createErr)
	}
	if createdAppointment == nil {
		t.Fatal("expected created appointment, got nil")
	}
	if createdAppointment.EndsAt.Sub(createdAppointment.StartsAt) != 45*time.Minute {
		t.Errorf("expected 45 minute slot, got %s", createdAppointment.EndsAt.Sub(createdAppointment.StartsAt))
	}
}

func TestCreateAppointment_RejectsOneHourDuration(t *testing.T) {
	repositoryMock := &MockRepository{}
	appointmentService := NewService(repositoryMock, nil, nil)
	slotStart := futureAlignedSlotStart(0)

	_, createErr := appointmentService.CreateAppointment(context.Background(), CreateAppointmentInput{
		PatientFHIRID: "patient-123",
		StaffID:       uuid.New(),
		StartsAt:      slotStart,
		EndsAt:        slotStart.Add(time.Hour),
	})

	if createErr == nil {
		t.Fatal("expected validation error, got nil")
	}
	assertFieldViolation(t, createErr, "ends_at")
}

func TestCreateAppointment_RejectsArbitraryDuration(t *testing.T) {
	repositoryMock := &MockRepository{}
	appointmentService := NewService(repositoryMock, nil, nil)
	slotStart := futureAlignedSlotStart(0)

	_, createErr := appointmentService.CreateAppointment(context.Background(), CreateAppointmentInput{
		PatientFHIRID: "patient-123",
		StaffID:       uuid.New(),
		StartsAt:      slotStart,
		EndsAt:        slotStart.Add(20 * time.Minute),
	})

	if createErr == nil {
		t.Fatal("expected validation error, got nil")
	}
	assertFieldViolation(t, createErr, "ends_at")
}

func TestCreateAppointment_RejectsUnalignedStart(t *testing.T) {
	repositoryMock := &MockRepository{}
	appointmentService := NewService(repositoryMock, nil, nil)
	unalignedStart := futureAlignedSlotStart(15)

	_, createErr := appointmentService.CreateAppointment(context.Background(), CreateAppointmentInput{
		PatientFHIRID: "patient-123",
		StaffID:       uuid.New(),
		StartsAt:      unalignedStart,
		EndsAt:        unalignedStart.Add(30 * time.Minute),
	})

	if createErr == nil {
		t.Fatal("expected validation error, got nil")
	}
	assertFieldViolation(t, createErr, "starts_at")
}

func TestCreateAppointment_RejectsPastStartsAt(t *testing.T) {
	repositoryMock := &MockRepository{}
	appointmentService := NewService(repositoryMock, nil, nil)
	pastTime := time.Now().Add(-time.Hour)

	_, createErr := appointmentService.CreateAppointment(context.Background(), CreateAppointmentInput{
		PatientFHIRID: "patient-123",
		StaffID:       uuid.New(),
		StartsAt:      pastTime,
		EndsAt:        pastTime.Add(30 * time.Minute),
	})

	if createErr == nil {
		t.Fatal("expected validation error, got nil")
	}
	assertFieldViolation(t, createErr, "starts_at")
}

func TestCreateAppointment_RejectsEndsBeforeStarts(t *testing.T) {
	repositoryMock := &MockRepository{}
	appointmentService := NewService(repositoryMock, nil, nil)
	slotStart := futureAlignedSlotStart(0)

	_, createErr := appointmentService.CreateAppointment(context.Background(), CreateAppointmentInput{
		PatientFHIRID: "patient-123",
		StaffID:       uuid.New(),
		StartsAt:      slotStart,
		EndsAt:        slotStart.Add(-time.Minute),
	})

	if createErr == nil {
		t.Fatal("expected validation error, got nil")
	}
	assertFieldViolation(t, createErr, "ends_at")
}

func TestCreateAppointment_PropagatesConflictFromRepository(t *testing.T) {
	repositoryMock := &MockRepository{
		CreateAppointmentFunc: func(ctx context.Context, appointment *Appointment) (*Appointment, error) {
			return nil, apperrors.ErrAppointmentConflict
		},
	}
	appointmentService := NewService(repositoryMock, nil, nil)
	slotStart := futureAlignedSlotStart(0)

	_, createErr := appointmentService.CreateAppointment(context.Background(), CreateAppointmentInput{
		PatientFHIRID: "patient-123",
		StaffID:       uuid.New(),
		StartsAt:      slotStart,
		EndsAt:        slotStart.Add(30 * time.Minute),
	})

	if createErr == nil {
		t.Fatal("expected conflict error, got nil")
	}
	if !errors.Is(createErr, apperrors.ErrAppointmentConflict) {
		t.Errorf("expected appointment conflict, got %v", createErr)
	}
}

func TestCreateAppointment_ReplaysCachedResponseForSameIdempotencyKey(t *testing.T) {
	slotStart := futureAlignedSlotStart(0)
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
	appointmentService := NewService(repositoryMock, nil, nil)

	createdAppointment, createErr := appointmentService.CreateAppointment(context.Background(), CreateAppointmentInput{
		PatientFHIRID:  "patient-123",
		StaffID:        uuid.New(),
		StartsAt:       slotStart,
		EndsAt:         slotStart.Add(30 * time.Minute),
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
	slotStart := futureAlignedSlotStart(0)
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
	appointmentService := NewService(repositoryMock, nil, nil)

	_, createErr := appointmentService.CreateAppointment(context.Background(), CreateAppointmentInput{
		PatientFHIRID:  "patient-123",
		StaffID:        uuid.New(),
		StartsAt:       slotStart,
		EndsAt:         slotStart.Add(30 * time.Minute),
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
	appointmentService := NewService(repositoryMock, nil, nil)

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
	appointmentService := NewService(repositoryMock, nil, nil)

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
	appointmentService := NewService(repositoryMock, nil, nil)

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
	appointmentService := NewService(repositoryMock, nil, nil)

	_, cancelErr := appointmentService.CancelAppointment(context.Background(), uuid.New())
	if cancelErr == nil {
		t.Fatal("expected not found error, got nil")
	}
	if !errors.Is(cancelErr, apperrors.ErrAppointmentNotFound) {
		t.Errorf("expected appointment not found, got %v", cancelErr)
	}
}

func TestListAppointmentsByPatient_RejectsEmptyFilter(t *testing.T) {
	appointmentService := NewService(&MockRepository{}, nil, nil)

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
	appointmentService := NewService(&MockRepository{}, nil, nil)

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
		EndsAt:        "2026-08-03T09:30:00Z",
		Reason:        "Follow-up",
	}

	firstHash := computeRequestHash(payload)
	secondHash := computeRequestHash(payload)
	if firstHash != secondHash {
		t.Errorf("expected deterministic hash, got %s and %s", firstHash, secondHash)
	}
}

func TestComputeRequestHash_NormalizesEmptyReason(t *testing.T) {
	staffID := uuid.New().String()
	emptyReasonPayload := CreateAppointmentRequest{
		PatientFhirID: "patient-123",
		StaffID:       staffID,
		StartsAt:      "2026-08-03T09:00:00Z",
		EndsAt:        "2026-08-03T09:30:00Z",
		Reason:        "",
	}
	absentReasonPayload := CreateAppointmentRequest{
		PatientFhirID: "patient-123",
		StaffID:       staffID,
		StartsAt:      "2026-08-03T09:00:00Z",
		EndsAt:        "2026-08-03T09:30:00Z",
	}
	withReasonPayload := CreateAppointmentRequest{
		PatientFhirID: "patient-123",
		StaffID:       staffID,
		StartsAt:      "2026-08-03T09:00:00Z",
		EndsAt:        "2026-08-03T09:30:00Z",
		Reason:        "Follow-up",
	}

	emptyReasonHash := computeRequestHash(emptyReasonPayload)
	absentReasonHash := computeRequestHash(absentReasonPayload)
	withReasonHash := computeRequestHash(withReasonPayload)

	if emptyReasonHash != absentReasonHash {
		t.Errorf("expected empty reason and absent reason to produce the same hash, got %s and %s", emptyReasonHash, absentReasonHash)
	}
	if emptyReasonHash == withReasonHash {
		t.Error("expected hash with a reason to differ from hash without a reason")
	}
}

func TestListAppointmentsByStaffInRange_ReturnsAppointments(t *testing.T) {
	staffID := uuid.New()
	rangeStart := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	rangeEnd := time.Date(2026, 8, 31, 0, 0, 0, 0, time.UTC)
	expectedAppointments := []*Appointment{{ID: uuid.New(), StaffID: staffID}}

	repositoryMock := &MockRepository{
		ListAppointmentsByStaffInRangeFunc: func(ctx context.Context, receivedStaffID uuid.UUID, receivedStart time.Time, receivedEnd time.Time) ([]*Appointment, error) {
			if receivedStaffID != staffID {
				t.Errorf("expected staff %s, got %s", staffID, receivedStaffID)
			}
			if !receivedStart.Equal(rangeStart) || !receivedEnd.Equal(rangeEnd) {
				t.Errorf("unexpected range bounds %s..%s", receivedStart, receivedEnd)
			}
			return expectedAppointments, nil
		},
	}
	appointmentService := NewService(repositoryMock, nil, nil)

	appointments, listErr := appointmentService.ListAppointmentsByStaffInRange(context.Background(), staffID, rangeStart, rangeEnd)
	if listErr != nil {
		t.Fatalf("unexpected error: %v", listErr)
	}
	if len(appointments) != 1 {
		t.Errorf("expected 1 appointment, got %d", len(appointments))
	}
}

func TestListAppointmentsByStaffInRange_RejectsInvertedRange(t *testing.T) {
	repositoryMock := &MockRepository{}
	appointmentService := NewService(repositoryMock, nil, nil)

	_, listErr := appointmentService.ListAppointmentsByStaffInRange(
		context.Background(),
		uuid.New(),
		time.Date(2026, 8, 31, 0, 0, 0, 0, time.UTC),
		time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC),
	)
	if listErr == nil {
		t.Fatal("expected validation error, got nil")
	}
	assertFieldViolation(t, listErr, "end_date")
}

func TestListAppointmentsByStaffInRange_RejectsMissingBounds(t *testing.T) {
	repositoryMock := &MockRepository{}
	appointmentService := NewService(repositoryMock, nil, nil)

	_, missingStartErr := appointmentService.ListAppointmentsByStaffInRange(context.Background(), uuid.New(), time.Time{}, time.Now())
	if missingStartErr == nil {
		t.Fatal("expected validation error for missing start_date, got nil")
	}
	assertFieldViolation(t, missingStartErr, "start_date")

	_, missingEndErr := appointmentService.ListAppointmentsByStaffInRange(context.Background(), uuid.New(), time.Now(), time.Time{})
	if missingEndErr == nil {
		t.Fatal("expected validation error for missing end_date, got nil")
	}
	assertFieldViolation(t, missingEndErr, "end_date")
}
