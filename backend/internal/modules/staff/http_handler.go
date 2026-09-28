package staff

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/healthcare/backend/internal/api/middleware"
	"github.com/healthcare/backend/internal/api/render"
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
	mux.Handle("GET /api/v1/staff/employees", middleware.RequirePolicy("GET /api/v1/staff/employees")(http.HandlerFunc(handler.ListEmployees)))
	mux.Handle("POST /api/v1/staff/employees", middleware.RequirePolicy("POST /api/v1/staff/employees")(http.HandlerFunc(handler.CreateEmployee)))
	mux.Handle("PATCH /api/v1/staff/employees/{employeeId}/status", middleware.RequirePolicy("PATCH /api/v1/staff/employees/{employeeId}/status")(http.HandlerFunc(handler.SetEmployeeStatus)))
	mux.Handle("DELETE /api/v1/staff/employees/{employeeId}", middleware.RequirePolicy("DELETE /api/v1/staff/employees/{employeeId}")(http.HandlerFunc(handler.DeleteEmployee)))
	mux.Handle("GET /api/v1/staff/departments", middleware.RequirePolicy("GET /api/v1/staff/departments")(http.HandlerFunc(handler.ListDepartments)))
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
		render.ErrorFromAppError(httpResponseWriter, employeesErr)
		return
	}

	responseList := make([]EmployeeResponse, 0, len(employeesList))
	for _, employee := range employeesList {
		responseList = append(responseList, mapEmployeeToHTTPResponse(employee))
	}
	render.JSON(httpResponseWriter, http.StatusOK, responseList)
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
		CreatedBy    string `json:"created_by"`
		FullName     string `json:"full_name"`
		Email        string `json:"email"`
		Role         string `json:"role"`
		CRMNumber    string `json:"crm_number"`
		DepartmentID string `json:"department_id"`
	}

	if payloadDecodeErr := json.NewDecoder(httpRequest.Body).Decode(&payload); payloadDecodeErr != nil {
		render.Error(httpResponseWriter, http.StatusBadRequest, "Payload inválido.")
		return
	}

	input := CreateEmployeeInput{
		CreatedBy:    payload.CreatedBy,
		FullName:     payload.FullName,
		Email:        payload.Email,
		Role:         payload.Role,
		CRMNumber:    payload.CRMNumber,
		DepartmentID: payload.DepartmentID,
	}

	employee, createErr := handler.service.CreateEmployee(httpRequest.Context(), input)
	if createErr != nil {
		slog.Error("failed to create employee", "error", createErr, "email", payload.Email, "request_id", middleware.GetRequestID(httpRequest.Context()))
		render.ErrorFromAppError(httpResponseWriter, createErr)
		return
	}

	fhirID := ""
	if employee.FHIRResourceID != nil {
		fhirID = *employee.FHIRResourceID
	}
	render.JSON(httpResponseWriter, http.StatusCreated, map[string]string{
		"employee_id":      employee.ID.String(),
		"fhir_resource_id": fhirID,
	})
}

// SetEmployeeStatus godoc
//
//	@Summary		Enable or disable an employee
//	@Description	Toggles the operational access of a healthcare professional without removing the record
//	@Tags			staff
//	@Accept			json
//	@Produce		json
//	@Param			employeeId	path		string						true	"Employee identifier"
//	@Param			body		body		SetEmployeeStatusRequest		true	"Desired status"
//	@Success		200			{object}	EmployeeResponse
//	@Failure		400			{object}	map[string]string
//	@Failure		404			{object}	map[string]string
//	@Router			/staff/employees/{employeeId}/status [patch]
func (handler *HTTPHandler) SetEmployeeStatus(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
	employeeID := httpRequest.PathValue("employeeId")

	var payload struct {
		IsActive bool `json:"is_active"`
	}

	if payloadDecodeErr := json.NewDecoder(httpRequest.Body).Decode(&payload); payloadDecodeErr != nil {
		render.Error(httpResponseWriter, http.StatusBadRequest, "Payload inválido.")
		return
	}

	employee, statusErr := handler.service.SetEmployeeActive(httpRequest.Context(), SetEmployeeActiveInput{
		EmployeeID: employeeID,
		IsActive:   payload.IsActive,
	})
	if statusErr != nil {
		slog.Error("failed to update employee status", "error", statusErr, "employee_id", employeeID, "request_id", middleware.GetRequestID(httpRequest.Context()))
		render.ErrorFromAppError(httpResponseWriter, statusErr)
		return
	}

	render.JSON(httpResponseWriter, http.StatusOK, mapEmployeeToHTTPResponse(employee))
}

// DeleteEmployee godoc
//
//	@Summary		Delete an employee
//	@Description	Deactivates a healthcare professional, preserving the record for auditing purposes
//	@Tags			staff
//	@Accept			json
//	@Produce		json
//	@Param			employeeId	path		string					true	"Employee identifier"
//	@Success		200			{object}	EmployeeResponse
//	@Failure		400			{object}	map[string]string
//	@Failure		404			{object}	map[string]string
//	@Router			/staff/employees/{employeeId} [delete]
func (handler *HTTPHandler) DeleteEmployee(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
	employeeID := httpRequest.PathValue("employeeId")

	employee, deleteErr := handler.service.SetEmployeeActive(httpRequest.Context(), SetEmployeeActiveInput{
		EmployeeID: employeeID,
		IsActive:   false,
	})
	if deleteErr != nil {
		slog.Error("failed to delete employee", "error", deleteErr, "employee_id", employeeID, "request_id", middleware.GetRequestID(httpRequest.Context()))
		render.ErrorFromAppError(httpResponseWriter, deleteErr)
		return
	}

	render.JSON(httpResponseWriter, http.StatusOK, mapEmployeeToHTTPResponse(employee))
}

// ListDepartments godoc
//
//	@Summary		List departments
//	@Description	Returns the active department catalog used to assign professionals to hospital units
//	@Tags			staff
//	@Accept			json
//	@Produce		json
//	@Success		200		{array}		DepartmentResponse
//	@Failure		500		{object}	map[string]string
//	@Router			/staff/departments [get]
func (handler *HTTPHandler) ListDepartments(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
	departmentsList, departmentsErr := handler.service.ListDepartments(httpRequest.Context())
	if departmentsErr != nil {
		slog.Error("failed to list departments", "error", departmentsErr, "request_id", middleware.GetRequestID(httpRequest.Context()))
		render.ErrorFromAppError(httpResponseWriter, departmentsErr)
		return
	}

	responseList := make([]DepartmentResponse, 0, len(departmentsList))
	for _, department := range departmentsList {
		responseList = append(responseList, DepartmentResponse{
			ID:   department.ID.String(),
			Name: department.Name,
		})
	}
	render.JSON(httpResponseWriter, http.StatusOK, responseList)
}

func mapEmployeeToHTTPResponse(employee *Employee) EmployeeResponse {
	crmValue := ""
	if employee.CRMNumber != nil {
		crmValue = *employee.CRMNumber
	}
	return EmployeeResponse{
		ID:             employee.ID.String(),
		FullName:       employee.FullName,
		Email:          employee.Email,
		Role:           string(employee.Role),
		CRMNumber:      crmValue,
		DepartmentID:   employee.DepartmentID.String(),
		DepartmentName: employee.DepartmentName,
		FHIRResourceID: employee.FHIRResourceID,
		IsActive:       employee.IsActive,
	}
}

type EmployeeResponse struct {
	ID             string  `json:"id"`
	FullName       string  `json:"full_name"`
	Email          string  `json:"email"`
	Role           string  `json:"role"`
	CRMNumber      string  `json:"crm_number"`
	DepartmentID   string  `json:"department_id"`
	DepartmentName string  `json:"department_name"`
	FHIRResourceID *string `json:"fhir_resource_id"`
	IsActive       bool    `json:"is_active"`
}

type SetEmployeeStatusRequest struct {
	IsActive bool `json:"is_active"`
}

type DepartmentResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
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
