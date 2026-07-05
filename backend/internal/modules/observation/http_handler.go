package observation

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

	mux.Handle("GET /api/v1/patients/{patientFhirId}/observations", clinicalRead(http.HandlerFunc(handler.ListObservationsByPatient)))
	mux.Handle("GET /api/v1/encounters/{encounterFhirId}/observations", clinicalRead(http.HandlerFunc(handler.ListObservationsByEncounter)))
	mux.Handle("POST /api/v1/encounters/{encounterFhirId}/observations", clinicalWrite(http.HandlerFunc(handler.CreateObservation)))
	mux.Handle("PUT /api/v1/observations/{observationFhirId}", clinicalWrite(http.HandlerFunc(handler.UpdateObservation)))
	mux.Handle("DELETE /api/v1/observations/{observationFhirId}", clinicalWrite(http.HandlerFunc(handler.DeleteObservation)))
}

// ListObservationsByPatient godoc
//
//	@Summary		List observations by patient
//	@Description	Returns all observations/vital signs for a patient
//	@Tags			observations
//	@Accept			json
//	@Produce		json
//	@Param			patientFhirId	path	string	true	"Patient FHIR ID"
//	@Success		200				{array}	ObservationResponse
//	@Failure		500				{object}	map[string]string
//	@Router			/patients/{patientFhirId}/observations [get]
func (handler *HTTPHandler) ListObservationsByPatient(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
	patientFhirID := httpRequest.PathValue("patientFhirId")

	observationsList, observationsErr := handler.service.GetObservationsByPatient(httpRequest.Context(), patientFhirID)
	if observationsErr != nil {
		slog.Error("failed to list observations by patient", "error", observationsErr, "patient_fhir_id", patientFhirID, "request_id", middleware.GetRequestID(httpRequest.Context()))
		render.Error(httpResponseWriter, http.StatusInternalServerError, "Erro ao carregar observações do paciente.")
		return
	}

	render.JSON(httpResponseWriter, http.StatusOK, toObservationResponseList(observationsList))
}

// ListObservationsByEncounter godoc
//
//	@Summary		List observations by encounter
//	@Description	Returns all observations/vital signs for an encounter
//	@Tags			observations
//	@Accept			json
//	@Produce		json
//	@Param			encounterFhirId	path	string	true	"Encounter FHIR ID"
//	@Success		200				{array}	ObservationResponse
//	@Failure		500				{object}	map[string]string
//	@Router			/encounters/{encounterFhirId}/observations [get]
func (handler *HTTPHandler) ListObservationsByEncounter(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
	encounterFhirID := httpRequest.PathValue("encounterFhirId")

	observationsList, observationsErr := handler.service.GetObservationsByEncounter(httpRequest.Context(), encounterFhirID)
	if observationsErr != nil {
		slog.Error("failed to list observations by encounter", "error", observationsErr, "encounter_fhir_id", encounterFhirID, "request_id", middleware.GetRequestID(httpRequest.Context()))
		render.Error(httpResponseWriter, http.StatusInternalServerError, "Erro ao carregar observações da consulta.")
		return
	}

	render.JSON(httpResponseWriter, http.StatusOK, toObservationResponseList(observationsList))
}

// CreateObservation godoc
//
//	@Summary		Create an observation
//	@Description	Creates a new observation/vital sign for an encounter
//	@Tags			observations
//	@Accept			json
//	@Produce		json
//	@Param			encounterFhirId	path	string	true	"Encounter FHIR ID"
//	@Param			body			body	CreateObservationRequest	true	"Observation data"
//	@Success		201				{object}	CreateObservationResponse
//	@Failure		400				{object}	map[string]string
//	@Failure		500				{object}	map[string]string
//	@Router			/encounters/{encounterFhirId}/observations [post]
func (handler *HTTPHandler) CreateObservation(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
	encounterFhirID := httpRequest.PathValue("encounterFhirId")

	var payload CreateObservationRequest

	if payloadDecodeErr := json.NewDecoder(httpRequest.Body).Decode(&payload); payloadDecodeErr != nil {
		render.Error(httpResponseWriter, http.StatusBadRequest, "Payload inválido.")
		return
	}

	if payload.PatientFhirID == "" {
		render.Error(httpResponseWriter, http.StatusBadRequest, "O identificador do paciente é obrigatório.")
		return
	}
	if payload.LoincCode == "" {
		render.Error(httpResponseWriter, http.StatusBadRequest, "O código LOINC é obrigatório.")
		return
	}

	newObservation := &Observation{
		EncounterFHIRID: encounterFhirID,
		PatientFHIRID:   payload.PatientFhirID,
		LoincCode:       payload.LoincCode,
		CodeDisplay:     payload.CodeDisplay,
		ValueQuantity:   payload.ValueQuantity,
		ValueUnit:       payload.ValueUnit,
		ObservedAt:      time.Now(),
	}

	createdObservation, createErr := handler.service.CreateObservation(httpRequest.Context(), newObservation)
	if createErr != nil {
		slog.Error("failed to create observation", "error", createErr, "encounter_fhir_id", encounterFhirID, "request_id", middleware.GetRequestID(httpRequest.Context()))
		render.Error(httpResponseWriter, http.StatusInternalServerError, "Erro ao criar observação.")
		return
	}

	render.JSON(httpResponseWriter, http.StatusCreated, ObservationResponse{
		FhirID:          createdObservation.FHIRResourceID,
		EncounterFhirID: createdObservation.EncounterFHIRID,
		PatientFhirID:   createdObservation.PatientFHIRID,
		LoincCode:       createdObservation.LoincCode,
		CodeDisplay:     createdObservation.CodeDisplay,
		ValueQuantity:   createdObservation.ValueQuantity,
		ValueUnit:       createdObservation.ValueUnit,
		CreatedAt:       createdObservation.ObservedAt.Format(time.RFC3339),
	})
}

// UpdateObservation godoc
//
//	@Summary		Update an observation
//	@Description	Updates an existing observation/vital sign
//	@Tags			observations
//	@Accept			json
//	@Produce		json
//	@Param			observationFhirId	path	string	true	"Observation FHIR ID"
//	@Param			body				body	CreateObservationRequest	true	"Observation data"
//	@Success		200					{array}	ObservationResponse
//	@Failure		400					{object}	map[string]string
//	@Failure		500					{object}	map[string]string
//	@Router			/observations/{observationFhirId} [put]
func (handler *HTTPHandler) UpdateObservation(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
	observationFhirID := httpRequest.PathValue("observationFhirId")

	var payload struct {
		PatientFhirID string  `json:"patient_fhir_id"`
		LoincCode     string  `json:"loinc_code"`
		CodeDisplay   string  `json:"code_display"`
		ValueQuantity float64 `json:"value_quantity"`
		ValueUnit     string  `json:"value_unit"`
	}

	if payloadDecodeErr := json.NewDecoder(httpRequest.Body).Decode(&payload); payloadDecodeErr != nil {
		render.Error(httpResponseWriter, http.StatusBadRequest, "Payload inválido.")
		return
	}

	updatedObservation := &Observation{
		PatientFHIRID:   payload.PatientFhirID,
		LoincCode:       payload.LoincCode,
		CodeDisplay:     payload.CodeDisplay,
		ValueQuantity:   payload.ValueQuantity,
		ValueUnit:       payload.ValueUnit,
	}

	resultObservation, updateErr := handler.service.UpdateObservation(httpRequest.Context(), observationFhirID, updatedObservation)
	if updateErr != nil {
		slog.Error("failed to update observation", "error", updateErr, "observation_fhir_id", observationFhirID, "request_id", middleware.GetRequestID(httpRequest.Context()))
		render.Error(httpResponseWriter, http.StatusInternalServerError, "Erro ao atualizar observação.")
		return
	}

	render.JSON(httpResponseWriter, http.StatusOK, ObservationResponse{
		FhirID:          resultObservation.FHIRResourceID,
		EncounterFhirID: resultObservation.EncounterFHIRID,
		PatientFhirID:   resultObservation.PatientFHIRID,
		LoincCode:       resultObservation.LoincCode,
		CodeDisplay:     resultObservation.CodeDisplay,
		ValueQuantity:   resultObservation.ValueQuantity,
		ValueUnit:       resultObservation.ValueUnit,
		CreatedAt:       resultObservation.ObservedAt.Format(time.RFC3339),
	})
}

// DeleteObservation godoc
//
//	@Summary		Delete an observation
//	@Description	Deletes an existing observation/vital sign
//	@Tags			observations
//	@Accept			json
//	@Produce		json
//	@Param			observationFhirId	path	string	true	"Observation FHIR ID"
//	@Success		204					{object}	nil
//	@Failure		500					{object}	map[string]string
//	@Router			/observations/{observationFhirId} [delete]
func (handler *HTTPHandler) DeleteObservation(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
	observationFhirID := httpRequest.PathValue("observationFhirId")

	if deleteErr := handler.service.DeleteObservation(httpRequest.Context(), observationFhirID); deleteErr != nil {
		slog.Error("failed to delete observation", "error", deleteErr, "observation_fhir_id", observationFhirID, "request_id", middleware.GetRequestID(httpRequest.Context()))
		render.Error(httpResponseWriter, http.StatusInternalServerError, "Erro ao deletar observação.")
		return
	}

	httpResponseWriter.WriteHeader(http.StatusNoContent)
}

func toObservationResponseList(observationsList []*Observation) []ObservationResponse {
	responseList := make([]ObservationResponse, 0, len(observationsList))
	for _, observation := range observationsList {
		responseList = append(responseList, ObservationResponse{
			FhirID:          observation.FHIRResourceID,
			EncounterFhirID: observation.EncounterFHIRID,
			PatientFhirID:   observation.PatientFHIRID,
			LoincCode:       observation.LoincCode,
			CodeDisplay:     observation.CodeDisplay,
			ValueQuantity:   observation.ValueQuantity,
			ValueUnit:       observation.ValueUnit,
			CreatedAt:       observation.ObservedAt.Format(time.RFC3339),
		})
	}
	return responseList
}

type ObservationResponse struct {
	FhirID          string  `json:"fhir_id"`
	EncounterFhirID string  `json:"encounter_fhir_id"`
	PatientFhirID   string  `json:"patient_fhir_id"`
	LoincCode       string  `json:"loinc_code"`
	CodeDisplay     string  `json:"code_display"`
	ValueQuantity   float64 `json:"value_quantity"`
	ValueUnit       string  `json:"value_unit"`
	CreatedAt       string  `json:"created_at"`
}

type CreateObservationRequest struct {
	PatientFhirID string  `json:"patient_fhir_id"`
	LoincCode     string  `json:"loinc_code"`
	CodeDisplay   string  `json:"code_display"`
	ValueQuantity float64 `json:"value_quantity"`
	ValueUnit     string  `json:"value_unit"`
}

type CreateObservationResponse struct {
	FhirID          string  `json:"fhir_id"`
	EncounterFhirID string  `json:"encounter_fhir_id"`
	PatientFhirID   string  `json:"patient_fhir_id"`
	LoincCode       string  `json:"loinc_code"`
	CodeDisplay     string  `json:"code_display"`
	ValueQuantity   float64 `json:"value_quantity"`
	ValueUnit       string  `json:"value_unit"`
	CreatedAt       string  `json:"created_at"`
}
