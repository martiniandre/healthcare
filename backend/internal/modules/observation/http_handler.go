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
	clinicalWrite := middleware.RequireRoles(role.RoleDoctor, role.RoleNurse)
	clinicalRead := middleware.RequireRoles(role.RoleAdmin, role.RoleDoctor, role.RoleNurse)

	mux.Handle("GET /api/patients/{patientFhirId}/observations", clinicalRead(http.HandlerFunc(handler.ListObservationsByPatient)))
	mux.Handle("GET /api/encounters/{encounterFhirId}/observations", clinicalRead(http.HandlerFunc(handler.ListObservationsByEncounter)))
	mux.Handle("POST /api/encounters/{encounterFhirId}/observations", clinicalWrite(http.HandlerFunc(handler.CreateObservation)))
}

func (handler *HTTPHandler) ListObservationsByPatient(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
	patientFhirID := httpRequest.PathValue("patientFhirId")

	observationsList, observationsErr := handler.service.GetObservationsByPatient(httpRequest.Context(), patientFhirID)
	if observationsErr != nil {
		slog.Error("failed to list observations by patient", "error", observationsErr, "patient_fhir_id", patientFhirID)
		render.Error(httpResponseWriter, http.StatusInternalServerError, "Erro ao carregar observações do paciente.")
		return
	}

	render.JSON(httpResponseWriter, http.StatusOK, toObservationResponseList(observationsList))
}

func (handler *HTTPHandler) ListObservationsByEncounter(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
	encounterFhirID := httpRequest.PathValue("encounterFhirId")

	observationsList, observationsErr := handler.service.GetObservationsByEncounter(httpRequest.Context(), encounterFhirID)
	if observationsErr != nil {
		slog.Error("failed to list observations by encounter", "error", observationsErr, "encounter_fhir_id", encounterFhirID)
		render.Error(httpResponseWriter, http.StatusInternalServerError, "Erro ao carregar observações da consulta.")
		return
	}

	render.JSON(httpResponseWriter, http.StatusOK, toObservationResponseList(observationsList))
}

func (handler *HTTPHandler) CreateObservation(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
	encounterFhirID := httpRequest.PathValue("encounterFhirId")

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
		slog.Error("failed to create observation", "error", createErr, "encounter_fhir_id", encounterFhirID)
		render.Error(httpResponseWriter, http.StatusInternalServerError, "Erro ao criar observação.")
		return
	}

	render.JSON(httpResponseWriter, http.StatusCreated, map[string]interface{}{
		"fhir_id":           createdObservation.FHIRResourceID,
		"encounter_fhir_id": createdObservation.EncounterFHIRID,
		"patient_fhir_id":   createdObservation.PatientFHIRID,
		"loinc_code":        createdObservation.LoincCode,
		"code_display":      createdObservation.CodeDisplay,
		"value_quantity":    createdObservation.ValueQuantity,
		"value_unit":        createdObservation.ValueUnit,
		"created_at":        createdObservation.ObservedAt.Format(time.RFC3339),
	})
}

func toObservationResponseList(observationsList []*Observation) []map[string]interface{} {
	responseList := make([]map[string]interface{}, 0, len(observationsList))
	for _, observation := range observationsList {
		responseList = append(responseList, map[string]interface{}{
			"fhir_id":           observation.FHIRResourceID,
			"encounter_fhir_id": observation.EncounterFHIRID,
			"patient_fhir_id":   observation.PatientFHIRID,
			"loinc_code":        observation.LoincCode,
			"code_display":      observation.CodeDisplay,
			"value_quantity":    observation.ValueQuantity,
			"value_unit":        observation.ValueUnit,
			"created_at":        observation.ObservedAt.Format(time.RFC3339),
		})
	}
	return responseList
}
