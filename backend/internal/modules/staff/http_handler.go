package staff

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
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

func (handler *HTTPHandler) RegisterRoutes(mux *http.ServeMux) {
	medicalStaff := middleware.RequireRoles(auth.RoleAdmin, auth.RoleDoctor, auth.RoleNurse, auth.RoleReception)
	adminOnly := middleware.RequireRoles(auth.RoleAdmin)

	mux.Handle("GET /api/staff/employees", medicalStaff(http.HandlerFunc(handler.ListEmployees)))
	mux.Handle("POST /api/staff/employees", adminOnly(http.HandlerFunc(handler.CreateEmployee)))
}

func (handler *HTTPHandler) ListEmployees(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
	employeesList, employeesErr := handler.service.ListEmployees(httpRequest.Context())
	if employeesErr != nil {
		slog.Error("failed to list employees", "error", employeesErr)
		render.Error(httpResponseWriter, http.StatusInternalServerError, "Erro ao listar corpo clínico.")
		return
	}

	render.JSON(httpResponseWriter, http.StatusOK, employeesList)
}

func (handler *HTTPHandler) CreateEmployee(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
	var payload struct {
		UserID    string `json:"user_id"`
		FullName  string `json:"full_name"`
		Email     string `json:"email"`
		Role      string `json:"role"`
		CRMNumber string `json:"crm_number"`
	}

	if payloadDecodeErr := json.NewDecoder(httpRequest.Body).Decode(&payload); payloadDecodeErr != nil {
		render.Error(httpResponseWriter, http.StatusBadRequest, "Payload inválido.")
		return
	}

	userIDParsed, parseErr := uuid.Parse(payload.UserID)
	if parseErr != nil {
		render.Error(httpResponseWriter, http.StatusBadRequest, "User ID inválido.")
		return
	}

	employee, createErr := handler.service.CreateEmployee(httpRequest.Context(), userIDParsed, payload.FullName, payload.Email, payload.Role, payload.CRMNumber)
	if createErr != nil {
		slog.Error("failed to create employee", "error", createErr, "email", payload.Email)
		render.Error(httpResponseWriter, http.StatusInternalServerError, "Erro ao registrar profissional.")
		return
	}

	render.JSON(httpResponseWriter, http.StatusCreated, map[string]string{
		"employee_id": employee.ID.String(),
	})
}
