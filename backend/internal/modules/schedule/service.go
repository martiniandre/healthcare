package schedule

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/healthcare/backend/internal/shared/apperrors"
	"github.com/healthcare/backend/internal/shared/ctxkeys"
	"github.com/healthcare/backend/internal/shared/eventbus"
)

type Service interface {
	CreateAppointment(ctx context.Context, input CreateAppointmentInput) (*Appointment, error)
	CancelAppointment(ctx context.Context, appointmentID uuid.UUID) (*Appointment, error)
	GetAppointment(ctx context.Context, appointmentID uuid.UUID) (*Appointment, error)
	ListAppointmentsByPatient(ctx context.Context, patientFHIRID string) ([]*Appointment, error)
	ListAppointmentsByStaffOnDate(ctx context.Context, staffID uuid.UUID, date time.Time) ([]*Appointment, error)
}

type service struct {
	repo     Repository
	eventBus eventbus.Bus
}

func NewService(repo Repository, eventBus eventbus.Bus) Service {
	return &service{repo: repo, eventBus: eventBus}
}

func (appointmentService *service) CreateAppointment(ctx context.Context, input CreateAppointmentInput) (*Appointment, error) {
	if input.IdempotencyKey != "" {
		cachedAppointment, cachedErr := appointmentService.resolveIdempotency(ctx, input)
		if cachedErr != nil {
			return nil, cachedErr
		}
		if cachedAppointment != nil {
			return cachedAppointment, nil
		}
	}

	fieldViolations := make(map[string]string)
	if input.PatientFHIRID == "" {
		fieldViolations["patient_fhir_id"] = "is required"
	}
	if input.StaffID == uuid.Nil {
		fieldViolations["staff_id"] = "is required"
	}
	if input.StartsAt.IsZero() {
		fieldViolations["starts_at"] = "is required"
	} else if !input.StartsAt.After(time.Now()) {
		fieldViolations["starts_at"] = "must be in the future"
	}
	if input.EndsAt.IsZero() {
		fieldViolations["ends_at"] = "is required"
	} else if !input.EndsAt.After(input.StartsAt) {
		fieldViolations["ends_at"] = "must be after starts_at"
	}
	if len(fieldViolations) > 0 {
		return nil, apperrors.InvalidArgument("invalid appointment input", fieldViolations)
	}

	newAppointment := &Appointment{
		ID:            uuid.New(),
		PatientFHIRID: input.PatientFHIRID,
		StaffID:       input.StaffID,
		StartsAt:      input.StartsAt,
		EndsAt:        input.EndsAt,
		Status:        AppointmentStatusScheduled,
		Reason:        input.Reason,
		Version:       1,
		CreatedBy:     actorFromContext(ctx),
	}

	createdAppointment, createErr := appointmentService.repo.CreateAppointment(ctx, newAppointment)
	if createErr != nil {
		return nil, createErr
	}

	appointmentService.publishAppointmentEvent(ctx, "appointment.created", createdAppointment)

	if input.IdempotencyKey != "" {
		appointmentService.storeIdempotencyKey(ctx, input, createdAppointment)
	}

	return createdAppointment, nil
}

func (appointmentService *service) CancelAppointment(ctx context.Context, appointmentID uuid.UUID) (*Appointment, error) {
	currentAppointment, getErr := appointmentService.repo.GetAppointmentByID(ctx, appointmentID)
	if getErr != nil {
		return nil, getErr
	}

	if currentAppointment.Status == AppointmentStatusCancelled {
		return currentAppointment, nil
	}
	if currentAppointment.Status == AppointmentStatusFinished {
		return nil, apperrors.ErrAppointmentInvalidTransition
	}

	cancelledAppointment, cancelErr := appointmentService.repo.CancelAppointment(ctx, appointmentID)
	if cancelErr != nil {
		return nil, cancelErr
	}

	appointmentService.publishAppointmentEvent(ctx, "appointment.cancelled", cancelledAppointment)

	return cancelledAppointment, nil
}

func (appointmentService *service) GetAppointment(ctx context.Context, appointmentID uuid.UUID) (*Appointment, error) {
	return appointmentService.repo.GetAppointmentByID(ctx, appointmentID)
}

func (appointmentService *service) ListAppointmentsByPatient(ctx context.Context, patientFHIRID string) ([]*Appointment, error) {
	if patientFHIRID == "" {
		return nil, apperrors.InvalidArgument("invalid appointment filter", map[string]string{"patient_fhir_id": "is required"})
	}
	return appointmentService.repo.ListAppointmentsByPatient(ctx, patientFHIRID)
}

func (appointmentService *service) ListAppointmentsByStaffOnDate(ctx context.Context, staffID uuid.UUID, date time.Time) ([]*Appointment, error) {
	if staffID == uuid.Nil {
		return nil, apperrors.InvalidArgument("invalid appointment filter", map[string]string{"staff_id": "is required"})
	}
	return appointmentService.repo.ListAppointmentsByStaffOnDate(ctx, staffID, date)
}

func (appointmentService *service) resolveIdempotency(ctx context.Context, input CreateAppointmentInput) (*Appointment, error) {
	storedKey, findErr := appointmentService.repo.FindIdempotencyKey(ctx, input.IdempotencyKey)
	if findErr != nil {
		return nil, findErr
	}
	if storedKey == nil {
		return nil, nil
	}

	if storedKey.RequestHash != input.RequestHash {
		return nil, apperrors.ErrBadRequest.WithFields(map[string]string{"idempotency_key": "reused with a different request payload"})
	}

	var cachedAppointment Appointment
	if unmarshalErr := json.Unmarshal(storedKey.ResponseBody, &cachedAppointment); unmarshalErr != nil {
		return nil, unmarshalErr
	}
	return &cachedAppointment, nil
}

func (appointmentService *service) storeIdempotencyKey(ctx context.Context, input CreateAppointmentInput, appointment *Appointment) {
	responseBody, marshalErr := json.Marshal(appointment)
	if marshalErr != nil {
		return
	}
	appointmentService.repo.SaveIdempotencyKey(ctx, &IdempotencyKey{
		ID:             uuid.New(),
		IdempotencyKey: input.IdempotencyKey,
		RequestHash:    input.RequestHash,
		ResponseStatus: 201,
		ResponseBody:   responseBody,
	})
}

func (appointmentService *service) publishAppointmentEvent(ctx context.Context, eventName string, appointment *Appointment) {
	if appointmentService.eventBus == nil {
		return
	}
	appointmentService.eventBus.Publish(ctx, eventbus.Event{
		Name: eventName,
		Data: map[string]any{
			"appointment_id": appointment.ID.String(),
			"patient_id":     appointment.PatientFHIRID,
			"staff_id":       appointment.StaffID.String(),
			"status":         string(appointment.Status),
		},
	})
}

func actorFromContext(ctx context.Context) *uuid.UUID {
	userIDString, exists := ctx.Value(ctxkeys.UserIDKey).(string)
	if !exists || userIDString == "" {
		return nil
	}
	parsedUserID, parseErr := uuid.Parse(userIDString)
	if parseErr != nil {
		return nil
	}
	return &parsedUserID
}
