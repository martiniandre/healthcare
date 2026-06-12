package patients

import (
	"encoding/json"
	"errors"
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

func (handler *HTTPHandler) RegisterRoutes(mux *http.ServeMux) {
	medicalStaff := middleware.RequireRoles(auth.RoleAdmin, auth.RoleDoctor, auth.RoleNurse, auth.RoleReception)
	adminOrReception := middleware.RequireRoles(auth.RoleAdmin, auth.RoleReception)

	mux.Handle("GET /api/patients", medicalStaff(http.HandlerFunc(handler.ListPatients)))
	mux.Handle("POST /api/patients", adminOrReception(http.HandlerFunc(handler.CreatePatient)))
	mux.Handle("GET /api/patients/{patientFhirId}", medicalStaff(http.HandlerFunc(handler.GetPatient)))
}

func (handler *HTTPHandler) ListPatients(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
	search := httpRequest.URL.Query().Get("search")
	sortField := httpRequest.URL.Query().Get("sortField")
	sortDirection := httpRequest.URL.Query().Get("sortDirection")
	page, _ := strconv.Atoi(httpRequest.URL.Query().Get("page"))
	if page <= 0 {
		page = 1
	}
	limit, _ := strconv.Atoi(httpRequest.URL.Query().Get("limit"))
	if limit <= 0 {
		limit = 50
	}

	patientsList, listError := handler.service.ListPatients(httpRequest.Context(), search, sortField, sortDirection, page, limit)
	if listError != nil {
		slog.Error("failed to list patients", "error", listError, "request_id", middleware.GetRequestID(httpRequest.Context()))
		render.Error(httpResponseWriter, http.StatusInternalServerError, "Erro ao listar pacientes.")
		return
	}

	type patientResponse struct {
		PatientID      string `json:"patient_id"`
		FHIRResourceID string `json:"fhir_resource_id"`
		FullName       string `json:"full_name"`
		BirthDate      string `json:"birth_date"`
		DocumentID     string `json:"document_id"`
		PhoneNumber    string `json:"phone_number"`
	}

	responseList := make([]patientResponse, 0, len(patientsList))
	for _, patient := range patientsList {
		responseList = append(responseList, patientResponse{
			PatientID:      patient.ID.String(),
			FHIRResourceID: patient.FHIRResourceID,
			FullName:       patient.FullName,
			BirthDate:      patient.BirthDate.Format("2006-01-02"),
			DocumentID:     patient.DocumentID,
			PhoneNumber:    patient.PhoneNumber,
		})
	}

	render.JSON(httpResponseWriter, http.StatusOK, responseList)
}

func (handler *HTTPHandler) CreatePatient(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
	var payload struct {
		FullName    string `json:"full_name"`
		BirthDate   string `json:"birth_date"`
		DocumentID  string `json:"document_id"`
		PhoneNumber string `json:"phone_number"`
	}

	if payloadDecodeErr := json.NewDecoder(httpRequest.Body).Decode(&payload); payloadDecodeErr != nil {
		render.Error(httpResponseWriter, http.StatusBadRequest, "Payload inválido.")
		return
	}

	patient, createPatientErr := handler.service.CreatePatient(httpRequest.Context(), payload.FullName, payload.BirthDate, payload.DocumentID, payload.PhoneNumber)
	if createPatientErr != nil {
		slog.Error("failed to create patient", "error", createPatientErr, "document_id", payload.DocumentID)
		if errors.Is(createPatientErr, ErrPatientAlreadyExists) {
			render.Error(httpResponseWriter, http.StatusConflict, "Paciente com este documento já cadastrado.")
			return
		}
		render.Error(httpResponseWriter, http.StatusInternalServerError, "Erro ao criar paciente.")
		return
	}

	render.JSON(httpResponseWriter, http.StatusCreated, map[string]string{
		"patient_id":       patient.ID.String(),
		"fhir_resource_id": patient.FHIRResourceID,
	})
}

func (handler *HTTPHandler) GetPatient(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
	patientFhirID := httpRequest.PathValue("patientFhirId")
	if patientFhirID == "" {
		render.Error(httpResponseWriter, http.StatusBadRequest, "ID do paciente ausente.")
		return
	}

	patient, getPatientErr := handler.service.GetPatient(httpRequest.Context(), patientFhirID)
	if getPatientErr != nil {
		slog.Error("patient not found", "error", getPatientErr, "patient_fhir_id", patientFhirID)
		render.Error(httpResponseWriter, http.StatusNotFound, "Paciente não encontrado.")
		return
	}

	render.JSON(httpResponseWriter, http.StatusOK, map[string]interface{}{
		"patient_id":       patient.ID.String(),
		"fhir_resource_id": patient.FHIRResourceID,
		"full_name":        patient.FullName,
		"birth_date":       patient.BirthDate.Format("2006-01-02"),
		"document_id":      patient.DocumentID,
		"phone_number":     patient.PhoneNumber,
	})
}
