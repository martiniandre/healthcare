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
	clinicalWrite := middleware.RequireRoles(role.RoleAdmin, role.RoleDoctor, role.RoleNurse)
	clinicalRead := middleware.RequireRoles(role.RoleAdmin, role.RoleDoctor, role.RoleNurse)

	mux.Handle("GET /api/v1/patients/{patientFhirId}/encounters", clinicalRead(http.HandlerFunc(handler.ListEncountersByPatient)))
	mux.Handle("POST /api/v1/patients/{patientFhirId}/encounters", clinicalWrite(http.HandlerFunc(handler.CreateEncounter)))
	mux.Handle("GET /api/v1/encounters/{encounterFhirId}", clinicalRead(http.HandlerFunc(handler.GetEncounter)))
	mux.Handle("PUT /api/v1/encounters/{encounterFhirId}", clinicalWrite(http.HandlerFunc(handler.UpdateEncounter)))
	mux.Handle("DELETE /api/v1/encounters/{encounterFhirId}", clinicalWrite(http.HandlerFunc(handler.DeleteEncounter)))
}

// ListEncountersByPatient godoc
//
//	@Summary		List encounters by patient
//	@Description	Returns all encounters/consultations for a patient
//	@Tags			encounters
//	@Accept			json
//	@Produce		json
//	@Param			patientFhirId	path	string	true	"Patient FHIR ID"
//	@Success		200				{array}	EncounterResponse
//	@Failure		500				{object}	map[string]string
//	@Router			/patients/{patientFhirId}/encounters [get]
func (handler *HTTPHandler) ListEncountersByPatient(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
	patientFhirID := httpRequest.PathValue("patientFhirId")

	encountersList, encountersErr := handler.service.GetEncountersByPatient(httpRequest.Context(), patientFhirID)
	if encountersErr != nil {
		slog.Error("failed to list encounters", "error", encountersErr, "patient_fhir_id", patientFhirID, "request_id", middleware.GetRequestID(httpRequest.Context()))
		render.Error(httpResponseWriter, http.StatusInternalServerError, "Erro ao carregar consultas do paciente.")
		return
	}

	responseList := make([]EncounterResponse, 0, len(encountersList))
	for _, encounter := range encountersList {
		responseList = append(responseList, EncounterResponse{
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

// CreateEncounter godoc
//
//	@Summary		Create an encounter
//	@Description	Creates a new encounter/consultation for a patient
//	@Tags			encounters
//	@Accept			json
//	@Produce		json
//	@Param			patientFhirId	path	string	true	"Patient FHIR ID"
//	@Param			body			body	CreateEncounterRequest	true	"Encounter data"
//	@Success		201				{object}	CreateEncounterResponse
//	@Failure		400				{object}	map[string]string
//	@Failure		500				{object}	map[string]string
//	@Router			/patients/{patientFhirId}/encounters [post]
func (handler *HTTPHandler) CreateEncounter(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
	patientFhirID := httpRequest.PathValue("patientFhirId")

	var payload CreateEncounterRequest

	if payloadDecodeErr := json.NewDecoder(httpRequest.Body).Decode(&payload); payloadDecodeErr != nil {
		render.Error(httpResponseWriter, http.StatusBadRequest, "Payload inválido.")
		return
	}

	if payload.ReasonDisplay == "" {
		render.Error(httpResponseWriter, http.StatusBadRequest, "O motivo da consulta é obrigatório.")
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
		slog.Error("failed to create encounter", "error", createErr, "patient_fhir_id", patientFhirID, "request_id", middleware.GetRequestID(httpRequest.Context()))
		render.Error(httpResponseWriter, http.StatusInternalServerError, "Erro ao criar consulta.")
		return
	}

	render.JSON(httpResponseWriter, http.StatusCreated, EncounterResponse{
		FhirID:         createdEncounter.FHIRResourceID,
		PatientFhirID:  createdEncounter.PatientFHIRID,
		Status:         createdEncounter.Status,
		ReasonDisplay:  createdEncounter.ReasonDisplay,
		PractitionerID: createdEncounter.PractitionerID,
		CreatedAt:      createdEncounter.StartedAt.Format(time.RFC3339),
	})
}

func (handler *HTTPHandler) GetEncounter(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
	encounterFhirID := httpRequest.PathValue("encounterFhirId")

	encounter, encounterErr := handler.service.GetEncounter(httpRequest.Context(), encounterFhirID)
	if encounterErr != nil {
		slog.Error("failed to get encounter", "error", encounterErr, "encounter_fhir_id", encounterFhirID, "request_id", middleware.GetRequestID(httpRequest.Context()))
		render.Error(httpResponseWriter, http.StatusInternalServerError, "Erro ao carregar consulta.")
		return
	}

	render.JSON(httpResponseWriter, http.StatusOK, EncounterResponse{
		FhirID:         encounter.FHIRResourceID,
		PatientFhirID:  encounter.PatientFHIRID,
		Status:         encounter.Status,
		ReasonDisplay:  encounter.ReasonDisplay,
		PractitionerID: encounter.PractitionerID,
		CreatedAt:      encounter.StartedAt.Format(time.RFC3339),
	})
}

func (handler *HTTPHandler) UpdateEncounter(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
	encounterFhirID := httpRequest.PathValue("encounterFhirId")

	var payload struct {
		ReasonDisplay  string `json:"reason_display"`
		Status         string `json:"status"`
		PractitionerID string `json:"practitioner_id"`
	}

	if payloadDecodeErr := json.NewDecoder(httpRequest.Body).Decode(&payload); payloadDecodeErr != nil {
		render.Error(httpResponseWriter, http.StatusBadRequest, "Payload inválido.")
		return
	}

	updatedEncounter := &Encounter{
		PatientFHIRID:  "", // will be fetched from existing
		PractitionerID: payload.PractitionerID,
		ReasonDisplay:  payload.ReasonDisplay,
		Status:         payload.Status,
	}

	currentEncounter, getErr := handler.service.GetEncounter(httpRequest.Context(), encounterFhirID)
	if getErr != nil {
		slog.Error("failed to get encounter for update", "error", getErr, "encounter_fhir_id", encounterFhirID, "request_id", middleware.GetRequestID(httpRequest.Context()))
		render.Error(httpResponseWriter, http.StatusInternalServerError, "Erro ao atualizar consulta.")
		return
	}
	updatedEncounter.PatientFHIRID = currentEncounter.PatientFHIRID
	updatedEncounter.ReasonCode = currentEncounter.ReasonCode

	result, updateErr := handler.service.UpdateEncounter(httpRequest.Context(), encounterFhirID, updatedEncounter)
	if updateErr != nil {
		slog.Error("failed to update encounter", "error", updateErr, "encounter_fhir_id", encounterFhirID, "request_id", middleware.GetRequestID(httpRequest.Context()))
		render.Error(httpResponseWriter, http.StatusInternalServerError, "Erro ao atualizar consulta.")
		return
	}

	render.JSON(httpResponseWriter, http.StatusOK, EncounterResponse{
		FhirID:         result.FHIRResourceID,
		PatientFhirID:  result.PatientFHIRID,
		Status:         result.Status,
		ReasonDisplay:  result.ReasonDisplay,
		PractitionerID: result.PractitionerID,
		CreatedAt:      result.StartedAt.Format(time.RFC3339),
	})
}

func (handler *HTTPHandler) DeleteEncounter(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
	encounterFhirID := httpRequest.PathValue("encounterFhirId")

	deleteErr := handler.service.DeleteEncounter(httpRequest.Context(), encounterFhirID)
	if deleteErr != nil {
		slog.Error("failed to delete encounter", "error", deleteErr, "encounter_fhir_id", encounterFhirID, "request_id", middleware.GetRequestID(httpRequest.Context()))
		render.Error(httpResponseWriter, http.StatusInternalServerError, "Erro ao excluir consulta.")
		return
	}

	httpResponseWriter.WriteHeader(http.StatusNoContent)
}

type EncounterResponse struct {
	FhirID         string `json:"fhir_id"`
	PatientFhirID  string `json:"patient_fhir_id"`
	Status         string `json:"status"`
	ReasonDisplay  string `json:"reason_display"`
	PractitionerID string `json:"practitioner_id,omitempty"`
	CreatedAt      string `json:"created_at"`
}

type CreateEncounterRequest struct {
	ReasonDisplay  string `json:"reason_display"`
	PractitionerID string `json:"practitioner_id"`
}

type CreateEncounterResponse struct {
	FhirID         string `json:"fhir_id"`
	PatientFhirID  string `json:"patient_fhir_id"`
	Status         string `json:"status"`
	ReasonDisplay  string `json:"reason_display"`
	PractitionerID string `json:"practitioner_id"`
	CreatedAt      string `json:"created_at"`
}
