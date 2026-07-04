package medication

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
	clinicalRead := middleware.RequireRoles(role.RoleAdmin, role.RoleDoctor, role.RoleNurse)

	mux.Handle("GET /api/v1/encounters/{encounterFhirId}/medications", clinicalRead(http.HandlerFunc(handler.ListMedicationsByEncounter)))
	mux.Handle("POST /api/v1/encounters/{encounterFhirId}/medications", middleware.RequireRoles(role.RoleDoctor)(http.HandlerFunc(handler.CreateMedication)))
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

	medicationsList, medicationsErr := handler.service.GetMedicationRequestsByEncounter(httpRequest.Context(), encounterFhirID)
	if medicationsErr != nil {
		slog.Error("failed to list medications", "error", medicationsErr, "encounter_fhir_id", encounterFhirID, "request_id", middleware.GetRequestID(httpRequest.Context()))
		render.Error(httpResponseWriter, http.StatusInternalServerError, "Erro ao carregar prescrições da consulta.")
		return
	}

	type medicationResponse struct {
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

	responseList := make([]medicationResponse, 0, len(medicationsList))
	for _, medication := range medicationsList {
		responseList = append(responseList, medicationResponse{
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
//	@Description	Creates a new medication request/prescription for an encounter
//	@Tags			medications
//	@Accept			json
//	@Produce		json
//	@Param			encounterFhirId	path	string	true	"Encounter FHIR ID"
//	@Param			body			body	CreateMedicationRequest	true	"Medication data"
//	@Success		201				{object}	CreateMedicationResponse
//	@Failure		400				{object}	map[string]string
//	@Failure		500				{object}	map[string]string
//	@Router			/encounters/{encounterFhirId}/medications [post]
func (handler *HTTPHandler) CreateMedication(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
	encounterFhirID := httpRequest.PathValue("encounterFhirId")

	var payload struct {
		PatientFhirID      string `json:"patient_fhir_id"`
		PractitionerFhirID string `json:"practitioner_fhir_id"`
		MedicationCode     string `json:"medication_code"`
		MedicationName     string `json:"medication_name"`
		DosageInstructions string `json:"dosage_instructions"`
	}

	if payloadDecodeErr := json.NewDecoder(httpRequest.Body).Decode(&payload); payloadDecodeErr != nil {
		render.Error(httpResponseWriter, http.StatusBadRequest, "Payload inválido.")
		return
	}

	newMedication := &Medication{
		EncounterFHIRID:    encounterFhirID,
		PatientFHIRID:      payload.PatientFhirID,
		PractitionerFHIRID: payload.PractitionerFhirID,
		MedicationCode:     payload.MedicationCode,
		MedicationName:     payload.MedicationName,
		DosageInstructions: payload.DosageInstructions,
		Status:             "active",
		IssuedAt:           time.Now(),
	}

	createdMedication, createErr := handler.service.CreateMedicationRequest(httpRequest.Context(), newMedication)
	if createErr != nil {
		slog.Error("failed to create medication request", "error", createErr, "encounter_fhir_id", encounterFhirID, "request_id", middleware.GetRequestID(httpRequest.Context()))
		render.Error(httpResponseWriter, http.StatusInternalServerError, "Erro ao criar prescrição.")
		return
	}

	render.JSON(httpResponseWriter, http.StatusCreated, map[string]interface{}{
		"fhir_id":              createdMedication.FHIRResourceID,
		"encounter_fhir_id":    createdMedication.EncounterFHIRID,
		"patient_fhir_id":      createdMedication.PatientFHIRID,
		"practitioner_fhir_id": createdMedication.PractitionerFHIRID,
		"medication_code":      createdMedication.MedicationCode,
		"medication_name":      createdMedication.MedicationName,
		"dosage_instructions":  createdMedication.DosageInstructions,
		"status":               createdMedication.Status,
		"created_at":           createdMedication.IssuedAt.Format(time.RFC3339),
	})
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
	PatientFhirID      string `json:"patient_fhir_id"`
	PractitionerFhirID string `json:"practitioner_fhir_id"`
	MedicationCode     string `json:"medication_code"`
	MedicationName     string `json:"medication_name"`
	DosageInstructions string `json:"dosage_instructions"`
}

type CreateMedicationResponse struct {
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
