package schedule

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/healthcare/backend/internal/modules/audit_logs"
	"github.com/healthcare/backend/internal/shared/apperrors"
	"github.com/healthcare/backend/internal/shared/ctxkeys"
	"github.com/healthcare/backend/internal/shared/eventbus"
	"github.com/healthcare/backend/internal/shared/payloaddiff"
)

type Service interface {
	CreateAppointment(ctx context.Context, input CreateAppointmentInput) (*Appointment, error)
	CancelAppointment(ctx context.Context, appointmentID uuid.UUID) (*Appointment, error)
	RescheduleAppointment(ctx context.Context, appointmentID uuid.UUID, startsAt time.Time, endsAt time.Time) (*Appointment, error)
	GetAppointment(ctx context.Context, appointmentID uuid.UUID) (*Appointment, error)
	ListAppointmentsByPatient(ctx context.Context, patientFHIRID string) ([]*Appointment, error)
	ListMyAppointments(ctx context.Context, authenticatedUserID string) ([]*Appointment, error)
	ListAppointmentsByStaffOnDate(ctx context.Context, staffID uuid.UUID, date time.Time) ([]*Appointment, error)
	ListAppointmentsByStaffInRange(ctx context.Context, staffID uuid.UUID, startDate time.Time, endDate time.Time) ([]*Appointment, error)
}

type service struct {
	repo         Repository
	eventBus     eventbus.Bus
	auditService audit_logs.Service
}

func NewService(repo Repository, eventBus eventbus.Bus, auditService audit_logs.Service) Service {
	return &service{repo: repo, eventBus: eventBus, auditService: auditService}
}

var allowedAppointmentDurations = []int{30, 45}

var allowedAppointmentStartMinutes = []int{0, 30, 45}

func isAllowedAppointmentDuration(appointmentStart time.Time, appointmentEnd time.Time) bool {
	durationMinutes := int(appointmentEnd.Sub(appointmentStart).Minutes())
	for _, allowedDuration := range allowedAppointmentDurations {
		if durationMinutes == allowedDuration {
			return true
		}
	}
	return false
}

func appointmentDurationViolationDescription() string {
	durationLabels := make([]string, 0, len(allowedAppointmentDurations))
	for _, allowedDuration := range allowedAppointmentDurations {
		durationLabels = append(durationLabels, strconv.Itoa(allowedDuration))
	}
	return "must be exactly " + strings.Join(durationLabels, " or ") + " minutes after starts_at"
}

func isAlignedToAllowedSlotStart(appointmentStart time.Time) bool {
	if appointmentStart.Second() != 0 {
		return false
	}
	startMinute := appointmentStart.Minute()
	for _, allowedMinute := range allowedAppointmentStartMinutes {
		if startMinute == allowedMinute {
			return true
		}
	}
	return false
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
	} else if !isAlignedToAllowedSlotStart(input.StartsAt) {
		fieldViolations["starts_at"] = "must start at an allowed slot time (:00, :30 or :45)"
	}
	if input.EndsAt.IsZero() {
		fieldViolations["ends_at"] = "is required"
	} else if !input.EndsAt.After(input.StartsAt) {
		fieldViolations["ends_at"] = "must be after starts_at"
	} else if !input.StartsAt.IsZero() && !isAllowedAppointmentDuration(input.StartsAt, input.EndsAt) {
		fieldViolations["ends_at"] = appointmentDurationViolationDescription()
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
	appointmentService.recordAppointmentAudit(ctx, "create", nil, createdAppointment)

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
	appointmentService.recordAppointmentAudit(ctx, "cancel", currentAppointment, cancelledAppointment)

	return cancelledAppointment, nil
}

func (appointmentService *service) RescheduleAppointment(ctx context.Context, appointmentID uuid.UUID, startsAt time.Time, endsAt time.Time) (*Appointment, error) {
	if appointmentID == uuid.Nil {
		return nil, apperrors.InvalidArgument("invalid appointment input", map[string]string{"appointment_id": "is required"})
	}

	currentAppointment, getErr := appointmentService.repo.GetAppointmentByID(ctx, appointmentID)
	if getErr != nil {
		return nil, getErr
	}
	if currentAppointment.Status == AppointmentStatusCancelled {
		return nil, apperrors.ErrAppointmentInvalidTransition
	}
	if currentAppointment.Status == AppointmentStatusFinished {
		return nil, apperrors.ErrAppointmentInvalidTransition
	}

	fieldViolations := make(map[string]string)
	if startsAt.IsZero() {
		fieldViolations["starts_at"] = "is required"
	} else if !startsAt.After(time.Now()) {
		fieldViolations["starts_at"] = "must be in the future"
	} else if !isAlignedToAllowedSlotStart(startsAt) {
		fieldViolations["starts_at"] = "must start at an allowed slot time (:00, :30 or :45)"
	}
	if endsAt.IsZero() {
		fieldViolations["ends_at"] = "is required"
	} else if !endsAt.After(startsAt) {
		fieldViolations["ends_at"] = "must be after starts_at"
	} else if !startsAt.IsZero() && !isAllowedAppointmentDuration(startsAt, endsAt) {
		fieldViolations["ends_at"] = appointmentDurationViolationDescription()
	}
	if len(fieldViolations) > 0 {
		return nil, apperrors.InvalidArgument("invalid appointment input", fieldViolations)
	}

	rescheduledAppointment, rescheduleErr := appointmentService.repo.RescheduleAppointment(ctx, appointmentID, startsAt, endsAt)
	if rescheduleErr != nil {
		return nil, rescheduleErr
	}

	appointmentService.publishAppointmentEvent(ctx, "appointment.rescheduled", rescheduledAppointment)
	appointmentService.recordAppointmentAudit(ctx, "reschedule", currentAppointment, rescheduledAppointment)

	return rescheduledAppointment, nil
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

func (appointmentService *service) ListMyAppointments(ctx context.Context, authenticatedUserID string) ([]*Appointment, error) {
	if authenticatedUserID == "" {
		return nil, apperrors.ErrUserNotFound
	}

	patientFHIRID, resolveErr := appointmentService.repo.ResolvePatientFHIRIDByUserID(ctx, authenticatedUserID)
	if resolveErr != nil {
		return nil, resolveErr
	}

	return appointmentService.repo.ListAppointmentsByPatient(ctx, patientFHIRID)
}

func (appointmentService *service) ListAppointmentsByStaffOnDate(ctx context.Context, staffID uuid.UUID, date time.Time) ([]*Appointment, error) {
	if staffID == uuid.Nil {
		return nil, apperrors.InvalidArgument("invalid appointment filter", map[string]string{"staff_id": "is required"})
	}
	return appointmentService.repo.ListAppointmentsByStaffOnDate(ctx, staffID, date)
}

func (appointmentService *service) ListAppointmentsByStaffInRange(ctx context.Context, staffID uuid.UUID, startDate time.Time, endDate time.Time) ([]*Appointment, error) {
	fieldViolations := make(map[string]string)
	if staffID == uuid.Nil {
		fieldViolations["staff_id"] = "is required"
	}
	if startDate.IsZero() {
		fieldViolations["start_date"] = "is required"
	}
	if endDate.IsZero() {
		fieldViolations["end_date"] = "is required"
	} else if !startDate.IsZero() && endDate.Before(startDate) {
		fieldViolations["end_date"] = "must not be before start_date"
	}
	if len(fieldViolations) > 0 {
		return nil, apperrors.InvalidArgument("invalid appointment filter", fieldViolations)
	}
	return appointmentService.repo.ListAppointmentsByStaffInRange(ctx, staffID, startDate, endDate)
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

func (appointmentService *service) recordAppointmentAudit(ctx context.Context, action string, beforeAppointment *Appointment, afterAppointment *Appointment) {
	if appointmentService.auditService == nil {
		return
	}
	payloadChanges, diffErr := payloaddiff.Compute(beforeAppointment, afterAppointment)
	if diffErr != nil {
		return
	}
	_, auditErr := appointmentService.auditService.CreateResourceAuditLog(ctx, audit_logs.ResourceAuditLog{
		CorrelationID: requestIDFromScheduleContext(ctx),
		CallerUserID:  scheduleActorIDFromContext(ctx),
		CallerRole:    scheduleRoleFromContext(ctx),
		Method:        "Appointment" + action,
		AccessGranted: true,
		ResourceType:  "appointment",
		ResourceID:    afterAppointment.ID.String(),
		Action:        action,
		PayloadDiff:   payloadChanges,
	})
	if auditErr != nil {
		return
	}
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

func requestIDFromScheduleContext(ctx context.Context) string {
	requestID, exists := ctx.Value(ctxkeys.RequestIDKey).(string)
	if exists && requestID != "" {
		return requestID
	}
	correlationID, correlationExists := ctx.Value(ctxkeys.CorrelationIDKey).(string)
	if correlationExists {
		return correlationID
	}
	return ""
}

func scheduleActorIDFromContext(ctx context.Context) string {
	userIDString, exists := ctx.Value(ctxkeys.UserIDKey).(string)
	if exists {
		return userIDString
	}
	return ""
}

func scheduleRoleFromContext(ctx context.Context) string {
	roleString, exists := ctx.Value(ctxkeys.RoleKey).(string)
	if exists {
		return roleString
	}
	return ""
}
