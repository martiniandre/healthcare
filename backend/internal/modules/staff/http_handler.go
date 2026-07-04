package staff

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
	adminOnly := middleware.RequireRoles(role.RoleAdmin)

	mux.Handle("GET /api/v1/staff/employees", medicalStaff(http.HandlerFunc(handler.ListEmployees)))
	mux.Handle("POST /api/v1/staff/employees", adminOnly(http.HandlerFunc(handler.CreateEmployee)))
}

// ListEmployees godoc
//
//	@Summary		List employees
//	@Description	Returns the list of healthcare staff/employees with optional search and role filter
//	@Tags			staff
//	@Accept			json
//	@Produce		json
//	@Param			search	query		string	false	"Search by name or email"
//	@Param			role	query		string	false	"Filter by role (admin, doctor, nurse, reception)"
//	@Success		200		{array}		EmployeeResponse
//	@Failure		500		{object}	map[string]string
//	@Router			/staff/employees [get]
func (handler *HTTPHandler) ListEmployees(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
	search := httpRequest.URL.Query().Get("search")
	role := httpRequest.URL.Query().Get("role")

	employeesList, employeesErr := handler.service.ListEmployees(httpRequest.Context(), search, role)
	if employeesErr != nil {
		slog.Error("failed to list employees", "error", employeesErr, "request_id", middleware.GetRequestID(httpRequest.Context()))
		render.Error(httpResponseWriter, http.StatusInternalServerError, "Erro ao listar corpo clínico.")
		return
	}

	render.JSON(httpResponseWriter, http.StatusOK, employeesList)
}

// CreateEmployee godoc
//
//	@Summary		Create a new employee
//	@Description	Registers a new healthcare professional as an employee
//	@Tags			staff
//	@Accept			json
//	@Produce		json
//	@Param			body	body		CreateEmployeeRequest	true	"Employee data"
//	@Success		201		{object}	CreateEmployeeResponse
//	@Failure		400		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Router			/staff/employees [post]
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
		slog.Error("failed to create employee", "error", createErr, "email", payload.Email, "request_id", middleware.GetRequestID(httpRequest.Context()))
		render.Error(httpResponseWriter, http.StatusInternalServerError, "Erro ao registrar profissional.")
		return
	}

	render.JSON(httpResponseWriter, http.StatusCreated, map[string]string{
		"employee_id": employee.ID.String(),
	})
}

type EmployeeResponse struct {
	ID        string `json:"id"`
	FullName  string `json:"full_name"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	CRMNumber string `json:"crm_number"`
}

type CreateEmployeeRequest struct {
	UserID    string `json:"user_id"`
	FullName  string `json:"full_name"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	CRMNumber string `json:"crm_number"`
}

type CreateEmployeeResponse struct {
	EmployeeID string `json:"employee_id"`
}
