package audit_logs

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/healthcare/backend/internal/api/middleware"
	"github.com/healthcare/backend/internal/api/render"
	"github.com/healthcare/backend/internal/shared/ctxkeys"
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

func (auditLogsHTTPHandler *HTTPHandler) RegisterRoutes(mux *http.ServeMux) {
	adminOnly := middleware.RequireRoles(role.RoleAdmin)
	authenticated := middleware.RequireRoles(role.RoleAdmin, role.RoleDoctor, role.RoleNurse, role.RoleReception, role.RolePatient)

	mux.Handle("GET /api/v1/audit-logs", adminOnly(http.HandlerFunc(auditLogsHTTPHandler.ListAuditLogs)))
	mux.Handle("POST /api/v1/audit-logs", authenticated(http.HandlerFunc(auditLogsHTTPHandler.CreateAuditLog)))
}

// ListAuditLogs godoc
//
//	@Summary		List audit logs
//	@Description	Returns a paginated list of audit logs (admin only)
//	@Tags			audit_logs
//	@Accept			json
//	@Produce		json
//	@Param			limit	query	int	false	"Number of items per page"	default(10)
//	@Param			offset	query	int	false	"Number of items to skip"	default(0)
//	@Success		200		{object}	AuditLogListResponse
//	@Failure		500		{object}	map[string]string
//	@Router			/audit-logs [get]
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
		slog.Error("failed to list audit logs", "error", listError, "request_id", middleware.GetRequestID(httpRequest.Context()))
		render.Error(httpResponseWriter, http.StatusInternalServerError, "Erro ao listar logs de auditoria.")
		return
	}

	render.JSON(httpResponseWriter, http.StatusOK, map[string]interface{}{
		"audit_logs": logs,
		"total":      totalCount,
	})
}

// CreateAuditLog godoc
//
//	@Summary		Create audit log entry
//	@Description	Creates a new audit log entry for tracking access and actions
//	@Tags			audit_logs
//	@Accept			json
//	@Produce		json
//	@Param			body	body	CreateAuditLogRequest	true	"Audit log data"
//	@Success		201		{object}	AuditLog
//	@Failure		400		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Router			/audit-logs [post]
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

	userID, _ := httpRequest.Context().Value(ctxkeys.UserIDKey).(string)
	role, _ := httpRequest.Context().Value(ctxkeys.RoleKey).(string)
	if payload.CallerUserID == "" {
		payload.CallerUserID = userID
	}
	if payload.CallerRole == "" {
		payload.CallerRole = role
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
		slog.Error("failed to create audit log", "error", createError, "request_id", middleware.GetRequestID(httpRequest.Context()))
		render.Error(httpResponseWriter, http.StatusInternalServerError, "Erro ao criar log de auditoria.")
		return
	}

	render.JSON(httpResponseWriter, http.StatusCreated, auditLog)
}

type AuditLogEntryResponse struct {
	ID            string `json:"id"`
	CorrelationID string `json:"correlation_id"`
	CallerUserID  string `json:"caller_user_id"`
	CallerRole    string `json:"caller_role"`
	Method        string `json:"method"`
	AccessGranted bool   `json:"access_granted"`
	CreatedAt     string `json:"created_at"`
}

type AuditLogListResponse struct {
	AuditLogs interface{} `json:"audit_logs"`
	Total     int         `json:"total"`
}

type CreateAuditLogRequest struct {
	CorrelationID string `json:"correlation_id"`
	CallerUserID  string `json:"caller_user_id"`
	CallerRole    string `json:"caller_role"`
	Method        string `json:"method"`
	AccessGranted bool   `json:"access_granted"`
}
