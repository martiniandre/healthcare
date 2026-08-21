package schedule

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/healthcare/backend/internal/api/middleware"
	"github.com/healthcare/backend/internal/api/render"
	"github.com/healthcare/backend/internal/shared/ctxkeys"
)

type HTTPHandler struct {
	service Service
}

func NewHTTPHandler(service Service) *HTTPHandler {
	return &HTTPHandler{
		service: service,
	}
}

func (handler *HTTPHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.Handle("POST /api/v1/appointments", middleware.RequirePolicy("POST /api/v1/appointments")(http.HandlerFunc(handler.CreateAppointment)))
	mux.Handle("GET /api/v1/appointments", middleware.RequirePolicy("GET /api/v1/appointments")(http.HandlerFunc(handler.ListAppointments)))
	mux.Handle("GET /api/v1/appointments/my", middleware.RequirePolicy("GET /api/v1/appointments/my")(http.HandlerFunc(handler.ListMyAppointments)))
	mux.Handle("GET /api/v1/appointments/{appointmentId}", middleware.RequirePolicy("GET /api/v1/appointments/{appointmentId}")(http.HandlerFunc(handler.GetAppointment)))
	mux.Handle("POST /api/v1/appointments/{appointmentId}/cancel", middleware.RequirePolicy("POST /api/v1/appointments/{appointmentId}/cancel")(http.HandlerFunc(handler.CancelAppointment)))
}

// CreateAppointment godoc
//
//	@Summary		Create an appointment
//	@Description	Books a time slot for a patient with a staff member (idempotent)
//	@Tags			appointments
//	@Accept			json
//	@Produce		json
//	@Param			body	body	CreateAppointmentRequest	true	"Appointment data"
//	@Success		201		{object}	AppointmentResponse
//	@Failure		400		{object}	map[string]string
//	@Failure		409		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Router			/appointments [post]
func (handler *HTTPHandler) CreateAppointment(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
	var payload CreateAppointmentRequest
	if payloadDecodeErr := json.NewDecoder(httpRequest.Body).Decode(&payload); payloadDecodeErr != nil {
		render.Error(httpResponseWriter, http.StatusBadRequest, "Payload inválido.")
		return
	}

	staffID, staffParseErr := uuid.Parse(payload.StaffID)
	if staffParseErr != nil {
		render.Error(httpResponseWriter, http.StatusBadRequest, "staff_id inválido.")
		return
	}
	startsAt, startsParseErr := time.Parse(time.RFC3339, payload.StartsAt)
	if startsParseErr != nil {
		render.Error(httpResponseWriter, http.StatusBadRequest, "starts_at inválido.")
		return
	}
	endsAt, endsParseErr := time.Parse(time.RFC3339, payload.EndsAt)
	if endsParseErr != nil {
		render.Error(httpResponseWriter, http.StatusBadRequest, "ends_at inválido.")
		return
	}

	createdAppointment, createErr := handler.service.CreateAppointment(httpRequest.Context(), CreateAppointmentInput{
		PatientFHIRID:  payload.PatientFhirID,
		StaffID:        staffID,
		StartsAt:       startsAt,
		EndsAt:         endsAt,
		Reason:         payload.Reason,
		IdempotencyKey: payload.IdempotencyKey,
		RequestHash:    computeRequestHash(payload),
	})
	if createErr != nil {
		slog.Error("failed to create appointment", "error", createErr, "request_id", middleware.GetRequestID(httpRequest.Context()))
		render.ErrorFromAppError(httpResponseWriter, createErr)
		return
	}

	render.JSON(httpResponseWriter, http.StatusCreated, toAppointmentResponse(createdAppointment))
}

// ListAppointments godoc
//
//	@Summary		List appointments
//	@Description	Lists appointments filtered by patient or by staff on a date
//	@Tags			appointments
//	@Accept			json
//	@Produce		json
//	@Param			patient_fhir_id	query	string	false	"Patient FHIR ID"
//	@Param			staff_id		query	string	false	"Staff UUID"
//	@Param			date			query	string	false	"Date (YYYY-MM-DD) for staff listing"
//	@Param			start_date		query	string	false	"Range start date (YYYY-MM-DD), requires end_date"
//	@Param			end_date		query	string	false	"Range end date (YYYY-MM-DD), requires start_date"
//	@Success		200				{array}	AppointmentResponse
//	@Failure		500				{object}	map[string]string
//	@Router			/appointments [get]
func (handler *HTTPHandler) ListAppointments(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
	patientFhirID := httpRequest.URL.Query().Get("patient_fhir_id")
	staffIDValue := httpRequest.URL.Query().Get("staff_id")
	dateValue := httpRequest.URL.Query().Get("date")

	if patientFhirID != "" {
		appointments, listErr := handler.service.ListAppointmentsByPatient(httpRequest.Context(), patientFhirID)
		if listErr != nil {
			render.ErrorFromAppError(httpResponseWriter, listErr)
			return
		}
		render.JSON(httpResponseWriter, http.StatusOK, toAppointmentResponseList(appointments))
		return
	}

	if staffIDValue != "" {
		staffID, staffParseErr := uuid.Parse(staffIDValue)
		if staffParseErr != nil {
			render.Error(httpResponseWriter, http.StatusBadRequest, "staff_id inválido.")
			return
		}
		startDateValue := httpRequest.URL.Query().Get("start_date")
		endDateValue := httpRequest.URL.Query().Get("end_date")

		if startDateValue != "" && endDateValue != "" {
			startDate, startDateParseErr := time.Parse("2006-01-02", startDateValue)
			if startDateParseErr != nil {
				render.Error(httpResponseWriter, http.StatusBadRequest, "start_date inválido.")
				return
			}
			endDate, endDateParseErr := time.Parse("2006-01-02", endDateValue)
			if endDateParseErr != nil {
				render.Error(httpResponseWriter, http.StatusBadRequest, "end_date inválido.")
				return
			}
			appointments, listErr := handler.service.ListAppointmentsByStaffInRange(httpRequest.Context(), staffID, startDate, endDate)
			if listErr != nil {
				render.ErrorFromAppError(httpResponseWriter, listErr)
				return
			}
			render.JSON(httpResponseWriter, http.StatusOK, toAppointmentResponseList(appointments))
			return
		}

		date, dateParseErr := time.Parse("2006-01-02", dateValue)
		if dateParseErr != nil {
			render.Error(httpResponseWriter, http.StatusBadRequest, "date inválido.")
			return
		}
		appointments, listErr := handler.service.ListAppointmentsByStaffOnDate(httpRequest.Context(), staffID, date)
		if listErr != nil {
			render.ErrorFromAppError(httpResponseWriter, listErr)
			return
		}
		render.JSON(httpResponseWriter, http.StatusOK, toAppointmentResponseList(appointments))
		return
	}

	render.Error(httpResponseWriter, http.StatusBadRequest, "Informe patient_fhir_id ou staff_id.")
}

// ListMyAppointments godoc
//
//	@Summary		List the authenticated patient's appointments
//	@Description	Returns appointments for the authenticated patient
//	@Tags			appointments
//	@Produce		json
//	@Success		200	{array}	AppointmentResponse
//	@Failure		401	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/appointments/my [get]
func (handler *HTTPHandler) ListMyAppointments(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
	patientFHIRID, ok := httpRequest.Context().Value(ctxkeys.UserIDKey).(string)
	if !ok || patientFHIRID == "" {
		render.Error(httpResponseWriter, http.StatusUnauthorized, "Usuário não autenticado.")
		return
	}

	appointments, listErr := handler.service.ListAppointmentsByPatient(httpRequest.Context(), patientFHIRID)
	if listErr != nil {
		render.ErrorFromAppError(httpResponseWriter, listErr)
		return
	}

	render.JSON(httpResponseWriter, http.StatusOK, toAppointmentResponseList(appointments))
}

// GetAppointment godoc
//
//	@Summary		Get an appointment
//	@Description	Returns a single appointment by ID
//	@Tags			appointments
//	@Accept			json
//	@Produce		json
//	@Param			appointmentId	path	string	true	"Appointment UUID"
//	@Success		200				{object}	AppointmentResponse
//	@Failure		404				{object}	map[string]string
//	@Router			/appointments/{appointmentId} [get]
func (handler *HTTPHandler) GetAppointment(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
	appointmentID, idParseErr := uuid.Parse(httpRequest.PathValue("appointmentId"))
	if idParseErr != nil {
		render.Error(httpResponseWriter, http.StatusBadRequest, "appointmentId inválido.")
		return
	}

	appointment, getErr := handler.service.GetAppointment(httpRequest.Context(), appointmentID)
	if getErr != nil {
		render.ErrorFromAppError(httpResponseWriter, getErr)
		return
	}

	render.JSON(httpResponseWriter, http.StatusOK, toAppointmentResponse(appointment))
}

// CancelAppointment godoc
//
//	@Summary		Cancel an appointment
//	@Description	Cancels an existing appointment (idempotent)
//	@Tags			appointments
//	@Accept			json
//	@Produce		json
//	@Param			appointmentId	path	string	true	"Appointment UUID"
//	@Success		200				{object}	AppointmentResponse
//	@Failure		400				{object}	map[string]string
//	@Failure		404				{object}	map[string]string
//	@Router			/appointments/{appointmentId}/cancel [post]
func (handler *HTTPHandler) CancelAppointment(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
	appointmentID, idParseErr := uuid.Parse(httpRequest.PathValue("appointmentId"))
	if idParseErr != nil {
		render.Error(httpResponseWriter, http.StatusBadRequest, "appointmentId inválido.")
		return
	}

	cancelledAppointment, cancelErr := handler.service.CancelAppointment(httpRequest.Context(), appointmentID)
	if cancelErr != nil {
		render.ErrorFromAppError(httpResponseWriter, cancelErr)
		return
	}

	render.JSON(httpResponseWriter, http.StatusOK, toAppointmentResponse(cancelledAppointment))
}

type CreateAppointmentRequest struct {
	PatientFhirID  string `json:"patient_fhir_id"`
	StaffID        string `json:"staff_id"`
	StartsAt       string `json:"starts_at"`
	EndsAt         string `json:"ends_at"`
	Reason         string `json:"reason"`
	IdempotencyKey string `json:"idempotency_key"`
}

type AppointmentResponse struct {
	ID            string `json:"id"`
	PatientFhirID string `json:"patient_fhir_id"`
	StaffID       string `json:"staff_id"`
	StartsAt      string `json:"starts_at"`
	EndsAt        string `json:"ends_at"`
	Status        string `json:"status"`
	Reason        string `json:"reason"`
	Version       int    `json:"version"`
	CreatedAt     string `json:"created_at"`
}

func toAppointmentResponse(appointment *Appointment) AppointmentResponse {
	return AppointmentResponse{
		ID:            appointment.ID.String(),
		PatientFhirID: appointment.PatientFHIRID,
		StaffID:       appointment.StaffID.String(),
		StartsAt:      appointment.StartsAt.Format(time.RFC3339),
		EndsAt:        appointment.EndsAt.Format(time.RFC3339),
		Status:        string(appointment.Status),
		Reason:        appointment.Reason,
		Version:       appointment.Version,
		CreatedAt:     appointment.CreatedAt.Format(time.RFC3339),
	}
}

func toAppointmentResponseList(appointments []*Appointment) []AppointmentResponse {
	responseList := make([]AppointmentResponse, 0, len(appointments))
	for _, appointment := range appointments {
		responseList = append(responseList, toAppointmentResponse(appointment))
	}
	return responseList
}

func computeRequestHash(payload CreateAppointmentRequest) string {
	normalizedReason := payload.Reason
	if normalizedReason == "" {
		normalizedReason = "no-reason"
	}
	canonicalPayload := payload.PatientFhirID + "|" + payload.StaffID + "|" + payload.StartsAt + "|" + payload.EndsAt + "|" + normalizedReason
	payloadHash := sha256.Sum256([]byte(canonicalPayload))
	return hex.EncodeToString(payloadHash[:])
}
