package clinical

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
	clinicalWrite := middleware.RequireRoles(role.RoleDoctor, role.RoleNurse)
	clinicalRead := middleware.RequireRoles(role.RoleAdmin, role.RoleDoctor, role.RoleNurse)

	mux.Handle("GET /api/patients/{patientFhirId}/encounters", medicalStaff(http.HandlerFunc(handler.ListEncountersByPatient)))
	mux.Handle("POST /api/patients/{patientFhirId}/encounters", medicalStaff(http.HandlerFunc(handler.CreateEncounter)))

	mux.Handle("GET /api/patients/{patientFhirId}/observations", clinicalRead(http.HandlerFunc(handler.ListObservationsByPatient)))

	mux.Handle("GET /api/patients/{patientFhirId}/conditions", clinicalRead(http.HandlerFunc(handler.ListConditionsByPatient)))
	mux.Handle("POST /api/patients/{patientFhirId}/conditions", clinicalWrite(http.HandlerFunc(handler.CreateCondition)))

	mux.Handle("GET /api/patients/{patientFhirId}/allergies", clinicalRead(http.HandlerFunc(handler.ListAllergiesByPatient)))
	mux.Handle("POST /api/patients/{patientFhirId}/allergies", clinicalWrite(http.HandlerFunc(handler.CreateAllergy)))

	mux.Handle("GET /api/encounters/{encounterFhirId}/observations", clinicalRead(http.HandlerFunc(handler.ListObservationsByEncounter)))
	mux.Handle("POST /api/encounters/{encounterFhirId}/observations", clinicalWrite(http.HandlerFunc(handler.CreateObservation)))

	mux.Handle("GET /api/encounters/{encounterFhirId}/reports", clinicalRead(http.HandlerFunc(handler.ListReportsByEncounter)))
	mux.Handle("POST /api/encounters/{encounterFhirId}/reports", clinicalWrite(http.HandlerFunc(handler.CreateReport)))

	mux.Handle("GET /api/encounters/{encounterFhirId}/medications", clinicalRead(http.HandlerFunc(handler.ListMedicationsByEncounter)))
	mux.Handle("POST /api/encounters/{encounterFhirId}/medications", middleware.RequireRoles(role.RoleDoctor)(http.HandlerFunc(handler.CreateMedication)))
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

func (handler *HTTPHandler) ListConditionsByPatient(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
	patientFhirID := httpRequest.PathValue("patientFhirId")

	conditionsList, conditionsErr := handler.service.GetConditionsByPatient(httpRequest.Context(), patientFhirID)
	if conditionsErr != nil {
		slog.Error("failed to list conditions", "error", conditionsErr, "patient_fhir_id", patientFhirID)
		render.Error(httpResponseWriter, http.StatusInternalServerError, "Erro ao carregar diagnósticos do paciente.")
		return
	}

	type conditionResponse struct {
		FhirID         string `json:"fhir_id"`
		PatientFhirID  string `json:"patient_fhir_id"`
		ICD10Code      string `json:"icd10_code"`
		CodeDisplay    string `json:"code_display"`
		ClinicalStatus string `json:"clinical_status"`
		CreatedAt      string `json:"created_at"`
	}

	responseList := make([]conditionResponse, 0, len(conditionsList))
	for _, condition := range conditionsList {
		responseList = append(responseList, conditionResponse{
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

func (handler *HTTPHandler) CreateCondition(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
	patientFhirID := httpRequest.PathValue("patientFhirId")

	var payload struct {
		ICD10Code   string `json:"icd10_code"`
		CodeDisplay string `json:"code_display"`
	}

	if payloadDecodeErr := json.NewDecoder(httpRequest.Body).Decode(&payload); payloadDecodeErr != nil {
		render.Error(httpResponseWriter, http.StatusBadRequest, "Payload inválido.")
		return
	}

	newCondition := &Condition{
		PatientFHIRID:  patientFhirID,
		ICD10Code:      payload.ICD10Code,
		CodeDisplay:    payload.CodeDisplay,
		ClinicalStatus: "active",
		OnsetAt:        time.Now(),
	}

	createdCondition, createErr := handler.service.CreateCondition(httpRequest.Context(), newCondition)
	if createErr != nil {
		slog.Error("failed to create condition", "error", createErr, "patient_fhir_id", patientFhirID)
		render.Error(httpResponseWriter, http.StatusInternalServerError, "Erro ao criar diagnóstico.")
		return
	}

	render.JSON(httpResponseWriter, http.StatusCreated, map[string]interface{}{
		"fhir_id":         createdCondition.FHIRResourceID,
		"patient_fhir_id": createdCondition.PatientFHIRID,
		"icd10_code":      createdCondition.ICD10Code,
		"code_display":    createdCondition.CodeDisplay,
		"clinical_status": createdCondition.ClinicalStatus,
		"created_at":      createdCondition.OnsetAt.Format(time.RFC3339),
	})
}

func (handler *HTTPHandler) ListAllergiesByPatient(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
	patientFhirID := httpRequest.PathValue("patientFhirId")

	allergiesList, allergiesErr := handler.service.GetAllergyIntolerancesByPatient(httpRequest.Context(), patientFhirID)
	if allergiesErr != nil {
		slog.Error("failed to list allergies", "error", allergiesErr, "patient_fhir_id", patientFhirID)
		render.Error(httpResponseWriter, http.StatusInternalServerError, "Erro ao carregar alergias do paciente.")
		return
	}

	type allergyResponse struct {
		FhirID          string `json:"fhir_id"`
		PatientFhirID   string `json:"patient_fhir_id"`
		AllergenCode    string `json:"allergen_code"`
		AllergenDisplay string `json:"allergen_display"`
		ClinicalStatus  string `json:"clinical_status"`
		Reaction        string `json:"reaction"`
		CreatedAt       string `json:"created_at"`
	}

	responseList := make([]allergyResponse, 0, len(allergiesList))
	for _, allergy := range allergiesList {
		responseList = append(responseList, allergyResponse{
			FhirID:          allergy.FHIRResourceID,
			PatientFhirID:   allergy.PatientFHIRID,
			AllergenCode:    allergy.AllergenCode,
			AllergenDisplay: allergy.AllergenDisplay,
			ClinicalStatus:  allergy.ClinicalStatus,
			Reaction:        allergy.Reaction,
			CreatedAt:       allergy.RecordedAt.Format(time.RFC3339),
		})
	}

	render.JSON(httpResponseWriter, http.StatusOK, responseList)
}

func (handler *HTTPHandler) CreateAllergy(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
	patientFhirID := httpRequest.PathValue("patientFhirId")

	var payload struct {
		AllergenCode    string `json:"allergen_code"`
		AllergenDisplay string `json:"allergen_display"`
		Reaction        string `json:"reaction"`
	}

	if payloadDecodeErr := json.NewDecoder(httpRequest.Body).Decode(&payload); payloadDecodeErr != nil {
		render.Error(httpResponseWriter, http.StatusBadRequest, "Payload inválido.")
		return
	}

	newAllergy := &AllergyIntolerance{
		PatientFHIRID:   patientFhirID,
		AllergenCode:    payload.AllergenCode,
		AllergenDisplay: payload.AllergenDisplay,
		ClinicalStatus:  "active",
		Reaction:        payload.Reaction,
		RecordedAt:      time.Now(),
	}

	createdAllergy, createErr := handler.service.CreateAllergyIntolerance(httpRequest.Context(), newAllergy)
	if createErr != nil {
		slog.Error("failed to create allergy", "error", createErr, "patient_fhir_id", patientFhirID)
		render.Error(httpResponseWriter, http.StatusInternalServerError, "Erro ao criar alergia.")
		return
	}

	render.JSON(httpResponseWriter, http.StatusCreated, map[string]interface{}{
		"fhir_id":          createdAllergy.FHIRResourceID,
		"patient_fhir_id":  createdAllergy.PatientFHIRID,
		"allergen_code":    createdAllergy.AllergenCode,
		"allergen_display": createdAllergy.AllergenDisplay,
		"clinical_status":  createdAllergy.ClinicalStatus,
		"reaction":         createdAllergy.Reaction,
		"created_at":       createdAllergy.RecordedAt.Format(time.RFC3339),
	})
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

func (handler *HTTPHandler) ListReportsByEncounter(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
	encounterFhirID := httpRequest.PathValue("encounterFhirId")

	reportsList, reportsErr := handler.service.GetDiagnosticReportsByEncounter(httpRequest.Context(), encounterFhirID)
	if reportsErr != nil {
		slog.Error("failed to list reports", "error", reportsErr, "encounter_fhir_id", encounterFhirID)
		render.Error(httpResponseWriter, http.StatusInternalServerError, "Erro ao carregar laudos da consulta.")
		return
	}

	type reportResponse struct {
		FhirID          string `json:"fhir_id"`
		EncounterFhirID string `json:"encounter_fhir_id"`
		PatientFhirID   string `json:"patient_fhir_id"`
		ReportDisplay   string `json:"report_display"`
		Status          string `json:"status"`
		Conclusion      string `json:"conclusion"`
		CreatedAt       string `json:"created_at"`
	}

	responseList := make([]reportResponse, 0, len(reportsList))
	for _, report := range reportsList {
		responseList = append(responseList, reportResponse{
			FhirID:          report.FHIRResourceID,
			EncounterFhirID: report.EncounterFHIRID,
			PatientFhirID:   report.PatientFHIRID,
			ReportDisplay:   report.ReportDisplay,
			Status:          report.Status,
			Conclusion:      report.Conclusion,
			CreatedAt:       report.IssuedAt.Format(time.RFC3339),
		})
	}

	render.JSON(httpResponseWriter, http.StatusOK, responseList)
}

func (handler *HTTPHandler) CreateReport(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
	encounterFhirID := httpRequest.PathValue("encounterFhirId")

	var payload struct {
		PatientFhirID string `json:"patient_fhir_id"`
		ReportDisplay string `json:"report_display"`
		Conclusion    string `json:"conclusion"`
	}

	if payloadDecodeErr := json.NewDecoder(httpRequest.Body).Decode(&payload); payloadDecodeErr != nil {
		render.Error(httpResponseWriter, http.StatusBadRequest, "Payload inválido.")
		return
	}

	newReport := &DiagnosticReport{
		EncounterFHIRID: encounterFhirID,
		PatientFHIRID:   payload.PatientFhirID,
		ReportCode:      "24323-8",
		ReportDisplay:   payload.ReportDisplay,
		Status:          "final",
		Conclusion:      payload.Conclusion,
		IssuedAt:        time.Now(),
	}

	createdReport, createErr := handler.service.CreateDiagnosticReport(httpRequest.Context(), newReport)
	if createErr != nil {
		slog.Error("failed to create diagnostic report", "error", createErr, "encounter_fhir_id", encounterFhirID)
		render.Error(httpResponseWriter, http.StatusInternalServerError, "Erro ao criar laudo.")
		return
	}

	render.JSON(httpResponseWriter, http.StatusCreated, map[string]interface{}{
		"fhir_id":           createdReport.FHIRResourceID,
		"encounter_fhir_id": createdReport.EncounterFHIRID,
		"patient_fhir_id":   createdReport.PatientFHIRID,
		"report_display":    createdReport.ReportDisplay,
		"status":            createdReport.Status,
		"conclusion":        createdReport.Conclusion,
		"created_at":        createdReport.IssuedAt.Format(time.RFC3339),
	})
}

func (handler *HTTPHandler) ListMedicationsByEncounter(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
	encounterFhirID := httpRequest.PathValue("encounterFhirId")

	medicationsList, medicationsErr := handler.service.GetMedicationRequestsByEncounter(httpRequest.Context(), encounterFhirID)
	if medicationsErr != nil {
		slog.Error("failed to list medications", "error", medicationsErr, "encounter_fhir_id", encounterFhirID)
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

	newMedication := &MedicationRequest{
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
		slog.Error("failed to create medication request", "error", createErr, "encounter_fhir_id", encounterFhirID)
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
