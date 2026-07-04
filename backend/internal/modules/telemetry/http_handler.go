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

	mux.Handle("GET /api/v1/telemetry/rooms", medicalStaff(http.HandlerFunc(handler.ListRooms)))
	mux.Handle("POST /api/v1/telemetry/rooms/{roomId}/unlock", clinicalStaff(http.HandlerFunc(handler.UnlockRoom)))
	mux.Handle("GET /api/v1/telemetry/rooms/{roomId}/beds", clinicalStaff(http.HandlerFunc(handler.ListBedsByRoom)))
	mux.Handle("POST /api/v1/telemetry/beds/{bedId}/condition", clinicalWrite(http.HandlerFunc(handler.UpdateBedCondition)))
}

// ListRooms godoc
//
//	@Summary		List telemetry rooms
//	@Description	Returns all monitored rooms with telemetry capabilities
//	@Tags			telemetry
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}	TelemetryRoomResponse
//	@Failure		500	{object}	map[string]string
//	@Router			/telemetry/rooms [get]
func (handler *HTTPHandler) ListRooms(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
	roomsList, roomsErr := handler.service.GetRooms(httpRequest.Context())
	if roomsErr != nil {
		slog.Error("failed to list rooms", "error", roomsErr, "request_id", middleware.GetRequestID(httpRequest.Context()))
		render.Error(httpResponseWriter, http.StatusInternalServerError, "Erro ao listar salas.")
		return
	}

	render.JSON(httpResponseWriter, http.StatusOK, roomsList)
}

// UnlockRoom godoc
//
//	@Summary		Unlock a telemetry room
//	@Description	Unlocks a room using a passcode
//	@Tags			telemetry
//	@Accept			json
//	@Produce		json
//	@Param			roomId	path	string	true	"Room UUID"
//	@Param			body	body	UnlockRoomRequest	true	"Passcode"
//	@Success		200		{object}	UnlockRoomResponse
//	@Failure		400		{object}	map[string]string
//	@Failure		401		{object}	map[string]string
//	@Router			/telemetry/rooms/{roomId}/unlock [post]
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

// ListBedsByRoom godoc
//
//	@Summary		List beds in a room
//	@Description	Returns all telemetry beds in a specific room
//	@Tags			telemetry
//	@Accept			json
//	@Produce		json
//	@Param			roomId	path	string	true	"Room UUID"
//	@Success		200		{array}	TelemetryBedResponse
//	@Failure		400		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Router			/telemetry/rooms/{roomId}/beds [get]
func (handler *HTTPHandler) ListBedsByRoom(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
	roomIDRaw := httpRequest.PathValue("roomId")

	roomIDParsed, parseErr := uuid.Parse(roomIDRaw)
	if parseErr != nil {
		render.Error(httpResponseWriter, http.StatusBadRequest, "ID de sala inválido.")
		return
	}

	bedsList, bedsErr := handler.service.GetBeds(httpRequest.Context(), roomIDParsed)
	if bedsErr != nil {
		slog.Error("failed to list beds", "error", bedsErr, "room_id", roomIDRaw, "request_id", middleware.GetRequestID(httpRequest.Context()))
		render.Error(httpResponseWriter, http.StatusInternalServerError, "Erro ao listar leitos.")
		return
	}

	render.JSON(httpResponseWriter, http.StatusOK, bedsList)
}

// UpdateBedCondition godoc
//
//	@Summary		Update bed telemetry condition
//	@Description	Updates vital signs and status for a telemetry bed
//	@Tags			telemetry
//	@Accept			json
//	@Produce		json
//	@Param			bedId	path	string	true	"Bed UUID"
//	@Param			body	body	UpdateBedConditionRequest	true	"Bed condition data"
//	@Success		200		{object}	UpdateBedConditionResponse
//	@Failure		400		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Router			/telemetry/beds/{bedId}/condition [post]
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
		slog.Error("failed to update bed condition", "error", updateErr, "bed_id", bedIDRaw, "request_id", middleware.GetRequestID(httpRequest.Context()))
		render.Error(httpResponseWriter, http.StatusInternalServerError, "Erro ao atualizar leito.")
		return
	}

	render.JSON(httpResponseWriter, http.StatusOK, map[string]bool{"success": true})
}

type TelemetryRoomResponse struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

type UnlockRoomRequest struct {
	Passcode string `json:"passcode"`
}

type UnlockRoomResponse struct {
	Success  bool   `json:"success"`
	RoomName string `json:"roomName"`
}

type TelemetryBedResponse struct {
	ID          string  `json:"id"`
	RoomID      string  `json:"room_id"`
	BedLabel    string  `json:"bed_label"`
	Bpm         int32   `json:"bpm"`
	Spo2        int32   `json:"spo2"`
	Temperature float64 `json:"temperature"`
	Status      string  `json:"status"`
	Condition   string  `json:"condition"`
}

type UpdateBedConditionRequest struct {
	Bpm         int32   `json:"bpm"`
	Spo2        int32   `json:"spo2"`
	Temperature float64 `json:"temperature"`
	Status      string  `json:"status"`
	Condition   string  `json:"condition"`
}

type UpdateBedConditionResponse struct {
	Success bool `json:"success"`
}
