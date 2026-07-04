package encounter

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

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

	mux.Handle("GET /api/patients/{patientFhirId}/encounters", medicalStaff(http.HandlerFunc(handler.ListEncountersByPatient)))
	mux.Handle("POST /api/patients/{patientFhirId}/encounters", medicalStaff(http.HandlerFunc(handler.CreateEncounter)))
}

func (handler *HTTPHandler) ListEncountersByPatient(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
	patientFhirID := httpRequest.PathValue("patientFhirId")

	encountersList, encountersErr := handler.service.GetEncountersByPatient(httpRequest.Context(), patientFhirID)
	if encountersErr != nil {
		slog.Error("failed to list encounters", "error", encountersErr, "patient_fhir_id", patientFhirID)
		render.Error(httpResponseWriter, http.StatusInternalServerError, "Erro ao carregar consultas do paciente.")
		return
	}

	type encounterResponse struct {
		FhirID         string `json:"fhir_id"`
		PatientFhirID  string `json:"patient_fhir_id"`
		Status         string `json:"status"`
		ReasonDisplay  string `json:"reason_display"`
		PractitionerID string `json:"practitioner_id,omitempty"`
		CreatedAt      string `json:"created_at"`
	}

	responseList := make([]encounterResponse, 0, len(encountersList))
	for _, encounter := range encountersList {
		responseList = append(responseList, encounterResponse{
			FhirID:         encounter.FHIRResourceID,
			PatientFhirID:  encounter.PatientFHIRID,
			Status:         encounter.Status,
			ReasonDisplay:  encounter.ReasonDisplay,
			PractitionerID: encounter.PractitionerID,
			CreatedAt:      encounter.StartedAt.Format(time.RFC3339),
		})
	}

	render.JSON(httpResponseWriter, http.StatusOK, responseList)
}

func (handler *HTTPHandler) CreateEncounter(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
	patientFhirID := httpRequest.PathValue("patientFhirId")

	var payload struct {
		ReasonDisplay  string `json:"reason_display"`
		PractitionerID string `json:"practitioner_id"`
	}

	if payloadDecodeErr := json.NewDecoder(httpRequest.Body).Decode(&payload); payloadDecodeErr != nil {
		render.Error(httpResponseWriter, http.StatusBadRequest, "Payload inválido.")
		return
	}

	newEncounter := &Encounter{
		PatientFHIRID:  patientFhirID,
		PractitionerID: payload.PractitionerID,
		ReasonDisplay:  payload.ReasonDisplay,
		Status:         "finished",
		StartedAt:      time.Now(),
	}

	createdEncounter, createErr := handler.service.CreateEncounter(httpRequest.Context(), newEncounter)
	if createErr != nil {
		slog.Error("failed to create encounter", "error", createErr, "patient_fhir_id", patientFhirID)
		render.Error(httpResponseWriter, http.StatusInternalServerError, "Erro ao criar consulta.")
		return
	}

	render.JSON(httpResponseWriter, http.StatusCreated, map[string]interface{}{
		"fhir_id":         createdEncounter.FHIRResourceID,
		"patient_fhir_id": createdEncounter.PatientFHIRID,
		"status":          createdEncounter.Status,
		"reason_display":  createdEncounter.ReasonDisplay,
		"practitioner_id": createdEncounter.PractitionerID,
		"created_at":      createdEncounter.StartedAt.Format(time.RFC3339),
	})
}
