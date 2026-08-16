package schedule

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/healthcare/backend/internal/api/middleware"
	"github.com/healthcare/backend/internal/api/render"
)

type UnavailabilityHTTPHandler struct {
	service UnavailabilityService
}

func NewUnavailabilityHTTPHandler(service UnavailabilityService) *UnavailabilityHTTPHandler {
	return &UnavailabilityHTTPHandler{
		service: service,
	}
}

func (handler *UnavailabilityHTTPHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.Handle("POST /api/v1/schedule/unavailability", middleware.RequirePolicy("POST /api/v1/schedule/unavailability")(http.HandlerFunc(handler.CreateUnavailability)))
	mux.Handle("GET /api/v1/schedule/unavailability", middleware.RequirePolicy("GET /api/v1/schedule/unavailability")(http.HandlerFunc(handler.ListUnavailability)))
	mux.Handle("DELETE /api/v1/schedule/unavailability/{unavailabilityId}", middleware.RequirePolicy("DELETE /api/v1/schedule/unavailability/{unavailabilityId}")(http.HandlerFunc(handler.DeleteUnavailability)))
}

func (handler *UnavailabilityHTTPHandler) CreateUnavailability(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
	var payload CreateUnavailabilityRequest
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

	createdUnavailability, createErr := handler.service.CreateUnavailability(httpRequest.Context(), CreateUnavailabilityInput{
		StaffID:  staffID,
		StartsAt: startsAt,
		EndsAt:   endsAt,
		Reason:   payload.Reason,
	})
	if createErr != nil {
		slog.Error("failed to create staff unavailability", "error", createErr, "request_id", middleware.GetRequestID(httpRequest.Context()))
		render.ErrorFromAppError(httpResponseWriter, createErr)
		return
	}

	render.JSON(httpResponseWriter, http.StatusCreated, toUnavailabilityResponse(createdUnavailability))
}

func (handler *UnavailabilityHTTPHandler) ListUnavailability(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
	staffIDValue := httpRequest.URL.Query().Get("staff_id")
	if staffIDValue == "" {
		render.Error(httpResponseWriter, http.StatusBadRequest, "Informe staff_id.")
		return
	}
	staffID, staffParseErr := uuid.Parse(staffIDValue)
	if staffParseErr != nil {
		render.Error(httpResponseWriter, http.StatusBadRequest, "staff_id inválido.")
		return
	}

	var fromTime time.Time
	var toTime time.Time
	fromValue := httpRequest.URL.Query().Get("from")
	if fromValue != "" {
		parsedFrom, fromParseErr := time.Parse(time.RFC3339, fromValue)
		if fromParseErr != nil {
			render.Error(httpResponseWriter, http.StatusBadRequest, "from inválido.")
			return
		}
		fromTime = parsedFrom
	}
	toValue := httpRequest.URL.Query().Get("to")
	if toValue != "" {
		parsedTo, toParseErr := time.Parse(time.RFC3339, toValue)
		if toParseErr != nil {
			render.Error(httpResponseWriter, http.StatusBadRequest, "to inválido.")
			return
		}
		toTime = parsedTo
	}

	unavailabilityWindows, listErr := handler.service.ListUnavailabilityByStaff(httpRequest.Context(), ListUnavailabilityInput{
		StaffID: staffID,
		From:    fromTime,
		To:      toTime,
	})
	if listErr != nil {
		render.ErrorFromAppError(httpResponseWriter, listErr)
		return
	}

	render.JSON(httpResponseWriter, http.StatusOK, toUnavailabilityResponseList(unavailabilityWindows))
}

func (handler *UnavailabilityHTTPHandler) DeleteUnavailability(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
	unavailabilityID, idParseErr := uuid.Parse(httpRequest.PathValue("unavailabilityId"))
	if idParseErr != nil {
		render.Error(httpResponseWriter, http.StatusBadRequest, "unavailabilityId inválido.")
		return
	}

	deletedUnavailability, deleteErr := handler.service.DeleteUnavailability(httpRequest.Context(), unavailabilityID)
	if deleteErr != nil {
		render.ErrorFromAppError(httpResponseWriter, deleteErr)
		return
	}

	render.JSON(httpResponseWriter, http.StatusOK, toUnavailabilityResponse(deletedUnavailability))
}

type CreateUnavailabilityRequest struct {
	StaffID  string `json:"staff_id"`
	StartsAt string `json:"starts_at"`
	EndsAt   string `json:"ends_at"`
	Reason   string `json:"reason"`
}

type UnavailabilityResponse struct {
	ID        string `json:"id"`
	StaffID   string `json:"staff_id"`
	StartsAt  string `json:"starts_at"`
	EndsAt    string `json:"ends_at"`
	Reason    string `json:"reason"`
	CreatedAt string `json:"created_at"`
}

func toUnavailabilityResponse(unavailability *StaffUnavailability) UnavailabilityResponse {
	return UnavailabilityResponse{
		ID:        unavailability.ID.String(),
		StaffID:   unavailability.StaffID.String(),
		StartsAt:  unavailability.StartsAt.Format(time.RFC3339),
		EndsAt:    unavailability.EndsAt.Format(time.RFC3339),
		Reason:    unavailability.Reason,
		CreatedAt: unavailability.CreatedAt.Format(time.RFC3339),
	}
}

func toUnavailabilityResponseList(unavailabilityWindows []*StaffUnavailability) []UnavailabilityResponse {
	responseList := make([]UnavailabilityResponse, 0, len(unavailabilityWindows))
	for _, unavailabilityWindow := range unavailabilityWindows {
		responseList = append(responseList, toUnavailabilityResponse(unavailabilityWindow))
	}
	return responseList
}
