package telemetry

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/healthcare/backend/internal/api/middleware"
	"github.com/healthcare/backend/internal/api/render"
	"github.com/healthcare/backend/internal/shared/role"
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
	medicalStaff := middleware.RequireRoles(role.RoleAdmin, role.RoleDoctor, role.RoleNurse, role.RoleReception)
	clinicalStaff := middleware.RequireRoles(role.RoleAdmin, role.RoleDoctor, role.RoleNurse)
	clinicalWrite := middleware.RequireRoles(role.RoleDoctor, role.RoleNurse)

	mux.Handle("GET /api/telemetry/rooms", medicalStaff(http.HandlerFunc(handler.ListRooms)))
	mux.Handle("POST /api/telemetry/rooms/{roomId}/unlock", clinicalStaff(http.HandlerFunc(handler.UnlockRoom)))
	mux.Handle("GET /api/telemetry/rooms/{roomId}/beds", clinicalStaff(http.HandlerFunc(handler.ListBedsByRoom)))
	mux.Handle("POST /api/telemetry/beds/{bedId}/condition", clinicalWrite(http.HandlerFunc(handler.UpdateBedCondition)))
}

func (handler *HTTPHandler) ListRooms(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
	roomsList, roomsErr := handler.service.GetRooms(httpRequest.Context())
	if roomsErr != nil {
		slog.Error("failed to list rooms", "error", roomsErr)
		render.Error(httpResponseWriter, http.StatusInternalServerError, "Erro ao listar salas.")
		return
	}

	render.JSON(httpResponseWriter, http.StatusOK, roomsList)
}

func (handler *HTTPHandler) UnlockRoom(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
	roomIDRaw := httpRequest.PathValue("roomId")

	roomIDParsed, parseErr := uuid.Parse(roomIDRaw)
	if parseErr != nil {
		render.Error(httpResponseWriter, http.StatusBadRequest, "ID de sala inválido.")
		return
	}

	var payload struct {
		Passcode string `json:"passcode"`
	}

	if payloadDecodeErr := json.NewDecoder(httpRequest.Body).Decode(&payload); payloadDecodeErr != nil {
		render.Error(httpResponseWriter, http.StatusBadRequest, "Payload inválido.")
		return
	}

	unlockedRoom, unlockErr := handler.service.UnlockRoom(httpRequest.Context(), roomIDParsed, payload.Passcode)
	if unlockErr != nil {
		slog.Warn("room unlock failed", "room_id", roomIDRaw, "error", unlockErr)
		render.Error(httpResponseWriter, http.StatusUnauthorized, "Senha inválida.")
		return
	}

	render.JSON(httpResponseWriter, http.StatusOK, map[string]interface{}{
		"success":  true,
		"roomName": unlockedRoom.Name,
	})
}

func (handler *HTTPHandler) ListBedsByRoom(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
	roomIDRaw := httpRequest.PathValue("roomId")

	roomIDParsed, parseErr := uuid.Parse(roomIDRaw)
	if parseErr != nil {
		render.Error(httpResponseWriter, http.StatusBadRequest, "ID de sala inválido.")
		return
	}

	bedsList, bedsErr := handler.service.GetBeds(httpRequest.Context(), roomIDParsed)
	if bedsErr != nil {
		slog.Error("failed to list beds", "error", bedsErr, "room_id", roomIDRaw)
		render.Error(httpResponseWriter, http.StatusInternalServerError, "Erro ao listar leitos.")
		return
	}

	render.JSON(httpResponseWriter, http.StatusOK, bedsList)
}

func (handler *HTTPHandler) UpdateBedCondition(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
	bedIDRaw := httpRequest.PathValue("bedId")

	bedIDParsed, parseErr := uuid.Parse(bedIDRaw)
	if parseErr != nil {
		render.Error(httpResponseWriter, http.StatusBadRequest, "ID de leito inválido.")
		return
	}

	var payload struct {
		Bpm         int32   `json:"bpm"`
		Spo2        int32   `json:"spo2"`
		Temperature float64 `json:"temperature"`
		Status      string  `json:"status"`
		Condition   string  `json:"condition"`
	}

	if payloadDecodeErr := json.NewDecoder(httpRequest.Body).Decode(&payload); payloadDecodeErr != nil {
		render.Error(httpResponseWriter, http.StatusBadRequest, "Payload inválido.")
		return
	}

	updateErr := handler.service.UpdateBedCondition(httpRequest.Context(), bedIDParsed, payload.Bpm, payload.Spo2, payload.Temperature, payload.Status, payload.Condition)
	if updateErr != nil {
		slog.Error("failed to update bed condition", "error", updateErr, "bed_id", bedIDRaw)
		render.Error(httpResponseWriter, http.StatusInternalServerError, "Erro ao atualizar leito.")
		return
	}

	render.JSON(httpResponseWriter, http.StatusOK, map[string]bool{"success": true})
}
