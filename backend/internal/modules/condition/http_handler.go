package condition

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

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
	mux.Handle("GET /api/v1/patients/{patientFhirId}/conditions", middleware.RequirePolicy("GET /api/v1/patients/{patientFhirId}/conditions")(http.HandlerFunc(handler.ListConditionsByPatient)))
	mux.Handle("POST /api/v1/patients/{patientFhirId}/conditions", middleware.RequirePolicy("POST /api/v1/patients/{patientFhirId}/conditions")(http.HandlerFunc(handler.CreateCondition)))
	mux.Handle("PUT /api/v1/patients/{patientFhirId}/conditions/{conditionFhirId}", middleware.RequirePolicy("PUT /api/v1/patients/{patientFhirId}/conditions/{conditionFhirId}")(http.HandlerFunc(handler.UpdateCondition)))
	mux.Handle("DELETE /api/v1/patients/{patientFhirId}/conditions/{conditionFhirId}", middleware.RequirePolicy("DELETE /api/v1/patients/{patientFhirId}/conditions/{conditionFhirId}")(http.HandlerFunc(handler.DeleteCondition)))
}

// ListConditionsByPatient godoc
//
//	@Summary		List conditions by patient
//	@Description	Returns all medical conditions/diagnoses for a patient
//	@Tags			conditions
//	@Accept			json
//	@Produce		json
//	@Param			patientFhirId	path	string	true	"Patient FHIR ID"
//	@Success		200				{array}	ConditionResponse
//	@Failure		500				{object}	map[string]string
//	@Router			/patients/{patientFhirId}/conditions [get]
func (handler *HTTPHandler) ListConditionsByPatient(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
	patientFhirID := httpRequest.PathValue("patientFhirId")

	conditionsList, conditionsErr := handler.service.GetConditionsByPatient(httpRequest.Context(), patientFhirID)
	if conditionsErr != nil {
		slog.Error("failed to list conditions", "error", conditionsErr, "patient_fhir_id", patientFhirID, "request_id", middleware.GetRequestID(httpRequest.Context()))
		render.ErrorFromAppError(httpResponseWriter, conditionsErr)
		return
	}

	responseList := make([]ConditionResponse, 0, len(conditionsList))
	for _, condition := range conditionsList {
		responseList = append(responseList, ConditionResponse{
			FhirID:         condition.FHIRResourceID,
			PatientFhirID:  condition.PatientFHIRID,
			ICD10Code:      condition.ICD10Code,
			CodeDisplay:    condition.CodeDisplay,
			ClinicalStatus: condition.ClinicalStatus,
			CreatedAt:      condition.OnsetAt.Format(time.RFC3339),
		})
	}

	render.JSON(httpResponseWriter, http.StatusOK, responseList)
}

// CreateCondition godoc
//
//	@Summary		Create a condition
//	@Description	Creates a new medical condition/diagnosis for a patient
//	@Tags			conditions
//	@Accept			json
//	@Produce		json
//	@Param			patientFhirId	path	string	true	"Patient FHIR ID"
//	@Param			body			body	CreateConditionRequest	true	"Condition data"
//	@Success		201				{object}	CreateConditionResponse
//	@Failure		400				{object}	map[string]string
//	@Failure		500				{object}	map[string]string
//	@Router			/patients/{patientFhirId}/conditions [post]
func (handler *HTTPHandler) CreateCondition(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
	patientFhirID := httpRequest.PathValue("patientFhirId")

	var payload CreateConditionRequest

	if payloadDecodeErr := json.NewDecoder(httpRequest.Body).Decode(&payload); payloadDecodeErr != nil {
		render.Error(httpResponseWriter, http.StatusBadRequest, "Payload inválido.")
		return
	}

	input := CreateConditionInput{
		PatientFHIRID:   patientFhirID,
		EncounterFHIRID: payload.EncounterID,
		ICD10Code:       payload.ICD10Code,
		CodeDisplay:     payload.CodeDisplay,
	}

	createdCondition, createErr := handler.service.CreateCondition(httpRequest.Context(), input)
	if createErr != nil {
		slog.Error("failed to create condition", "error", createErr, "patient_fhir_id", patientFhirID, "request_id", middleware.GetRequestID(httpRequest.Context()))
		render.ErrorFromAppError(httpResponseWriter, createErr)
		return
	}

	render.JSON(httpResponseWriter, http.StatusCreated, ConditionResponse{
		FhirID:         createdCondition.FHIRResourceID,
		PatientFhirID:  createdCondition.PatientFHIRID,
		ICD10Code:      createdCondition.ICD10Code,
		CodeDisplay:    createdCondition.CodeDisplay,
		ClinicalStatus: createdCondition.ClinicalStatus,
		CreatedAt:      createdCondition.OnsetAt.Format(time.RFC3339),
	})
}

func (handler *HTTPHandler) UpdateCondition(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
	conditionFhirID := httpRequest.PathValue("conditionFhirId")

	var payload struct {
		ICD10Code      *string `json:"icd10_code"`
		CodeDisplay    *string `json:"code_display"`
		ClinicalStatus *string `json:"clinical_status"`
		EncounterID    *string `json:"encounter_id"`
	}

	if payloadDecodeErr := json.NewDecoder(httpRequest.Body).Decode(&payload); payloadDecodeErr != nil {
		render.Error(httpResponseWriter, http.StatusBadRequest, "Payload inválido.")
		return
	}

	input := UpdateConditionInput{
		ICD10Code:       payload.ICD10Code,
		CodeDisplay:     payload.CodeDisplay,
		ClinicalStatus:  payload.ClinicalStatus,
		EncounterFHIRID: payload.EncounterID,
	}

	conditionResult, updateErr := handler.service.UpdateCondition(httpRequest.Context(), conditionFhirID, input)
	if updateErr != nil {
		slog.Error("failed to update condition", "error", updateErr, "condition_fhir_id", conditionFhirID, "request_id", middleware.GetRequestID(httpRequest.Context()))
		render.ErrorFromAppError(httpResponseWriter, updateErr)
		return
	}

	render.JSON(httpResponseWriter, http.StatusOK, ConditionResponse{
		FhirID:         conditionResult.FHIRResourceID,
		PatientFhirID:  conditionResult.PatientFHIRID,
		ICD10Code:      conditionResult.ICD10Code,
		CodeDisplay:    conditionResult.CodeDisplay,
		ClinicalStatus: conditionResult.ClinicalStatus,
		CreatedAt:      conditionResult.OnsetAt.Format(time.RFC3339),
	})
}

func (handler *HTTPHandler) DeleteCondition(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
	conditionFhirID := httpRequest.PathValue("conditionFhirId")

	if deleteErr := handler.service.DeleteCondition(httpRequest.Context(), conditionFhirID); deleteErr != nil {
		slog.Error("failed to delete condition", "error", deleteErr, "condition_fhir_id", conditionFhirID, "request_id", middleware.GetRequestID(httpRequest.Context()))
		render.ErrorFromAppError(httpResponseWriter, deleteErr)
		return
	}

	httpResponseWriter.WriteHeader(http.StatusNoContent)
}

type ConditionResponse struct {
	FhirID         string `json:"fhir_id"`
	PatientFhirID  string `json:"patient_fhir_id"`
	ICD10Code      string `json:"icd10_code"`
	CodeDisplay    string `json:"code_display"`
	ClinicalStatus string `json:"clinical_status"`
	CreatedAt      string `json:"created_at"`
}

type CreateConditionRequest struct {
	ICD10Code   string `json:"icd10_code"`
	CodeDisplay string `json:"code_display"`
	EncounterID string `json:"encounter_id"`
}

type CreateConditionResponse struct {
	FhirID         string `json:"fhir_id"`
	PatientFhirID  string `json:"patient_fhir_id"`
	ICD10Code      string `json:"icd10_code"`
	CodeDisplay    string `json:"code_display"`
	ClinicalStatus string `json:"clinical_status"`
	CreatedAt      string `json:"created_at"`
}
