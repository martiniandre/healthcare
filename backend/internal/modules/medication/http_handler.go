package medication

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/healthcare/backend/internal/api/middleware"
	"github.com/healthcare/backend/internal/api/render"
	"github.com/healthcare/backend/internal/shared/apperrors"
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
	mux.Handle("GET /api/v1/encounters/{encounterFhirId}/medications", middleware.RequirePolicy("GET /api/v1/encounters/{encounterFhirId}/medications")(http.HandlerFunc(handler.ListMedicationsByEncounter)))
	mux.Handle("POST /api/v1/encounters/{encounterFhirId}/medications", middleware.RequirePolicy("POST /api/v1/encounters/{encounterFhirId}/medications")(http.HandlerFunc(handler.CreateMedication)))
	mux.Handle("POST /api/v1/patients/{patientFhirId}/medications/check", middleware.RequirePolicy("POST /api/v1/patients/{patientFhirId}/medications/check")(http.HandlerFunc(handler.CheckMedicationSafety)))
}

// ListMedicationsByEncounter godoc
//
//	@Summary		List medications by encounter
//	@Description	Returns all medication requests/prescriptions for an encounter
//	@Tags			medications
//	@Accept			json
//	@Produce		json
//	@Param			encounterFhirId	path	string	true	"Encounter FHIR ID"
//	@Success		200				{array}	MedicationResponse
//	@Failure		500				{object}	map[string]string
//	@Router			/encounters/{encounterFhirId}/medications [get]
func (handler *HTTPHandler) ListMedicationsByEncounter(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
	encounterFhirID := httpRequest.PathValue("encounterFhirId")
	if encounterFhirID == "" {
		render.Error(httpResponseWriter, http.StatusBadRequest, "O identificador do atendimento é obrigatório.")
		return
	}

	medicationsList, medicationsErr := handler.service.GetMedicationRequestsByEncounter(httpRequest.Context(), encounterFhirID)
	if medicationsErr != nil {
		slog.Error("failed to list medications", "error", medicationsErr, "encounter_fhir_id", encounterFhirID, "request_id", middleware.GetRequestID(httpRequest.Context()))
		render.ErrorFromAppError(httpResponseWriter, medicationsErr)
		return
	}

	responseList := make([]MedicationResponse, 0, len(medicationsList))
	for _, medication := range medicationsList {
		responseList = append(responseList, MedicationResponse{
			FhirID:             medication.FHIRResourceID,
			EncounterFhirID:    medication.EncounterFHIRID,
			PatientFhirID:      medication.PatientFHIRID,
			PractitionerFhirID: medication.PractitionerFHIRID,
			MedicationCode:     medication.MedicationCode,
			MedicationName:     medication.MedicationName,
			DosageInstructions: medication.DosageInstructions,
			Status:             medication.Status,
			CreatedAt:          medication.IssuedAt.Format(time.RFC3339),
		})
	}

	render.JSON(httpResponseWriter, http.StatusOK, responseList)
}

// CreateMedication godoc
//
//	@Summary		Create a medication
//	@Description	Creates a new medication request/prescription for an encounter. An active allergy matching the drug (by ingredient or class) blocks the request with 409 MEDICATION_ALLERGY_BLOCKED unless a written overrideJustification (min 10 characters) is sent by a DOCTOR or ADMIN. Drug interactions never block and are returned as warnings.
//	@Tags			medications
//	@Accept			json
//	@Produce		json
//	@Param			encounterFhirId	path	string	true	"Encounter FHIR ID"
//	@Param			body			body	CreateMedicationRequest	true	"Medication data"
//	@Success		201				{object}	CreateMedicationResponse
//	@Failure		400				{object}	map[string]string
//	@Failure		403				{object}	map[string]string
//	@Failure		409				{object}	MedicationBlockedResponse
//	@Failure		500				{object}	map[string]string
//	@Router			/encounters/{encounterFhirId}/medications [post]
func (handler *HTTPHandler) CreateMedication(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
	encounterFhirID := httpRequest.PathValue("encounterFhirId")
	if encounterFhirID == "" {
		render.Error(httpResponseWriter, http.StatusBadRequest, "O identificador do atendimento é obrigatório.")
		return
	}

	var payload CreateMedicationRequest

	if payloadDecodeErr := json.NewDecoder(httpRequest.Body).Decode(&payload); payloadDecodeErr != nil {
		render.Error(httpResponseWriter, http.StatusBadRequest, "Payload inválido.")
		return
	}

	input := CreateMedicationInput{
		EncounterFHIRID:       encounterFhirID,
		PatientFHIRID:         payload.PatientFhirID,
		PractitionerFHIRID:    payload.PractitionerFhirID,
		MedicationCode:        payload.MedicationCode,
		MedicationName:        payload.MedicationName,
		DosageInstructions:    payload.DosageInstructions,
		OverrideJustification: payload.OverrideJustification,
	}

	createResult, createErr := handler.service.CreateMedicationRequest(httpRequest.Context(), input)
	if createErr != nil {
		var blockedError *BlockedByAllergyError
		if errors.As(createErr, &blockedError) {
			render.JSON(httpResponseWriter, http.StatusConflict, MedicationBlockedResponse{
				Error:  blockedError.Error(),
				Code:   apperrors.ErrMedicationBlockedByAllergy.Code,
				Alerts: toSafetyAlertResponses(blockedError.Alerts),
			})
			return
		}
		slog.Error("failed to create medication request", "error", createErr, "encounter_fhir_id", encounterFhirID, "request_id", middleware.GetRequestID(httpRequest.Context()))
		render.ErrorFromAppError(httpResponseWriter, createErr)
		return
	}

	createdMedication := createResult.Medication
	render.JSON(httpResponseWriter, http.StatusCreated, CreateMedicationResponse{
		FhirID:             createdMedication.FHIRResourceID,
		EncounterFhirID:    createdMedication.EncounterFHIRID,
		PatientFhirID:      createdMedication.PatientFHIRID,
		PractitionerFhirID: createdMedication.PractitionerFHIRID,
		MedicationCode:     createdMedication.MedicationCode,
		MedicationName:     createdMedication.MedicationName,
		DosageInstructions: createdMedication.DosageInstructions,
		Status:             createdMedication.Status,
		CreatedAt:          createdMedication.IssuedAt.Format(time.RFC3339),
		Warnings:           toSafetyAlertResponses(createResult.Warnings),
	})
}

// CheckMedicationSafety godoc
//
//	@Summary		Check medication safety
//	@Description	Cross-checks a requested medication against the patient's active allergies (blocking) and active medications (warning only)
//	@Tags			medications
//	@Accept			json
//	@Produce		json
//	@Param			patientFhirId	path	string	true	"Patient FHIR ID"
//	@Param			body			body	CheckMedicationSafetyRequest	true	"Medication to check"
//	@Success		200				{object}	CheckMedicationSafetyResponse
//	@Failure		400				{object}	map[string]string
//	@Failure		503				{object}	map[string]string
//	@Router			/patients/{patientFhirId}/medications/check [post]
func (handler *HTTPHandler) CheckMedicationSafety(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
	patientFhirID := httpRequest.PathValue("patientFhirId")
	if patientFhirID == "" {
		render.Error(httpResponseWriter, http.StatusBadRequest, "O identificador do paciente é obrigatório.")
		return
	}

	var payload CheckMedicationSafetyRequest
	if payloadDecodeErr := json.NewDecoder(httpRequest.Body).Decode(&payload); payloadDecodeErr != nil {
		render.Error(httpResponseWriter, http.StatusBadRequest, "Payload inválido.")
		return
	}

	alerts, checkErr := handler.service.CheckSafety(httpRequest.Context(), CheckSafetyInput{
		PatientFHIRID:  patientFhirID,
		MedicationCode: payload.MedicationCode,
		MedicationName: payload.MedicationName,
	})
	if checkErr != nil {
		slog.Error("failed to check medication safety", "error", checkErr, "patient_fhir_id", patientFhirID, "request_id", middleware.GetRequestID(httpRequest.Context()))
		render.ErrorFromAppError(httpResponseWriter, checkErr)
		return
	}

	render.JSON(httpResponseWriter, http.StatusOK, CheckMedicationSafetyResponse{Alerts: toSafetyAlertResponses(alerts)})
}

func toSafetyAlertResponses(alerts []SafetyAlert) []SafetyAlertResponse {
	responses := make([]SafetyAlertResponse, 0, len(alerts))
	for _, alert := range alerts {
		params := alert.Params
		if params == nil {
			params = map[string]string{}
		}
		responses = append(responses, SafetyAlertResponse{
			Kind:              alert.Kind,
			Severity:          alert.Severity,
			Blocking:          alert.Blocking,
			MessageKey:        alert.MessageKey,
			Params:            params,
			RelatedResourceID: alert.RelatedResourceID,
		})
	}
	return responses
}

type MedicationResponse struct {
	FhirID             string `json:"fhir_id"`
	EncounterFhirID    string `json:"encounter_fhir_id"`
	PatientFhirID      string `json:"patient_fhir_id"`
	PractitionerFhirID string `json:"practitioner_fhir_id"`
	MedicationCode     string `json:"medication_code"`
	MedicationName     string `json:"medication_name"`
	DosageInstructions string `json:"dosage_instructions"`
	Status             string `json:"status"`
	CreatedAt          string `json:"created_at"`
}

type CreateMedicationRequest struct {
	PatientFhirID         string `json:"patient_fhir_id"`
	PractitionerFhirID    string `json:"practitioner_fhir_id"`
	MedicationCode        string `json:"medication_code"`
	MedicationName        string `json:"medication_name"`
	DosageInstructions    string `json:"dosage_instructions"`
	OverrideJustification string `json:"overrideJustification,omitempty"`
}

type CreateMedicationResponse struct {
	FhirID             string                `json:"fhir_id"`
	EncounterFhirID    string                `json:"encounter_fhir_id"`
	PatientFhirID      string                `json:"patient_fhir_id"`
	PractitionerFhirID string                `json:"practitioner_fhir_id"`
	MedicationCode     string                `json:"medication_code"`
	MedicationName     string                `json:"medication_name"`
	DosageInstructions string                `json:"dosage_instructions"`
	Status             string                `json:"status"`
	CreatedAt          string                `json:"created_at"`
	Warnings           []SafetyAlertResponse `json:"warnings"`
}

type SafetyAlertResponse struct {
	Kind              string            `json:"kind" example:"allergy" enums:"allergy,interaction"`
	Severity          string            `json:"severity" example:"contraindicated" enums:"info,moderate,major,contraindicated"`
	Blocking          bool              `json:"blocking"`
	MessageKey        string            `json:"messageKey" example:"medicationSafety.allergy.class"`
	Params            map[string]string `json:"params"`
	RelatedResourceID string            `json:"relatedResourceId" example:"allergy-fhir-id"`
}

type CheckMedicationSafetyRequest struct {
	MedicationCode string `json:"medicationCode"`
	MedicationName string `json:"medicationName"`
}

type CheckMedicationSafetyResponse struct {
	Alerts []SafetyAlertResponse `json:"alerts"`
}

type MedicationBlockedResponse struct {
	Error  string                `json:"error"`
	Code   string                `json:"code" example:"MEDICATION_ALLERGY_BLOCKED"`
	Alerts []SafetyAlertResponse `json:"alerts"`
}
