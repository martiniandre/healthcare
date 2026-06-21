package audit_logs

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/healthcare/backend/internal/api/middleware"
	"github.com/healthcare/backend/internal/api/render"
	"github.com/healthcare/backend/internal/modules/auth"
)

type HTTPHandler struct {
	service Service
}

func NewHTTPHandler(service Service) *HTTPHandler {
	return &HTTPHandler{
		service: service,
	}
}

func (auditLogsHTTPHandler *HTTPHandler) RegisterRoutes(mux *http.ServeMux) {
	adminOnly := middleware.RequireRoles(auth.RoleAdmin)
	authenticated := middleware.RequireRoles(auth.RoleAdmin, auth.RoleDoctor, auth.RoleNurse, auth.RoleReception, auth.RolePatient)

	mux.Handle("GET /api/audit-logs", adminOnly(http.HandlerFunc(auditLogsHTTPHandler.ListAuditLogs)))
	mux.Handle("POST /api/audit-logs", authenticated(http.HandlerFunc(auditLogsHTTPHandler.CreateAuditLog)))
}

func (auditLogsHTTPHandler *HTTPHandler) ListAuditLogs(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
	limitStr := httpRequest.URL.Query().Get("limit")
	offsetStr := httpRequest.URL.Query().Get("offset")

	limit := 10
	offset := 0

	if limitStr != "" {
		if parsedLimit, parseLimitError := strconv.Atoi(limitStr); parseLimitError == nil {
			limit = parsedLimit
		}
	}
	if offsetStr != "" {
		if parsedOffset, parseOffsetError := strconv.Atoi(offsetStr); parseOffsetError == nil {
			offset = parsedOffset
		}
	}

	logs, totalCount, listError := auditLogsHTTPHandler.service.ListAuditLogs(httpRequest.Context(), limit, offset)
	if listError != nil {
		slog.Error("failed to list audit logs", "error", listError)
		render.Error(httpResponseWriter, http.StatusInternalServerError, "Erro ao listar logs de auditoria.")
		return
	}

	render.JSON(httpResponseWriter, http.StatusOK, map[string]interface{}{
		"audit_logs": logs,
		"total":      totalCount,
	})
}

func (auditLogsHTTPHandler *HTTPHandler) CreateAuditLog(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
	var payload struct {
		CorrelationID string `json:"correlation_id"`
		CallerUserID  string `json:"caller_user_id"`
		CallerRole    string `json:"caller_role"`
		Method        string `json:"method"`
		AccessGranted bool   `json:"access_granted"`
	}

	if payloadDecodeErr := json.NewDecoder(httpRequest.Body).Decode(&payload); payloadDecodeErr != nil {
		render.Error(httpResponseWriter, http.StatusBadRequest, "Payload inválido.")
		return
	}

	auditLog, createError := auditLogsHTTPHandler.service.CreateAuditLog(
		httpRequest.Context(),
		payload.CorrelationID,
		payload.CallerUserID,
		payload.CallerRole,
		payload.Method,
		payload.AccessGranted,
	)
	if createError != nil {
		slog.Error("failed to create audit log", "error", createError)
		render.Error(httpResponseWriter, http.StatusInternalServerError, "Erro ao criar log de auditoria.")
		return
	}

	render.JSON(httpResponseWriter, http.StatusCreated, auditLog)
}
