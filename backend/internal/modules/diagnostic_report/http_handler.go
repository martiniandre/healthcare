package diagnostic_report

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

	mux.Handle("GET /api/encounters/{encounterFhirId}/reports", clinicalRead(http.HandlerFunc(handler.ListReportsByEncounter)))
	mux.Handle("POST /api/encounters/{encounterFhirId}/reports", clinicalWrite(http.HandlerFunc(handler.CreateReport)))
}

// ListReportsByEncounter godoc
//
//	@Summary		List diagnostic reports by encounter
//	@Description	Returns all diagnostic reports for an encounter
//	@Tags			diagnostic_reports
//	@Accept			json
//	@Produce		json
//	@Param			encounterFhirId	path	string	true	"Encounter FHIR ID"
//	@Success		200				{array}	DiagnosticReportResponse
//	@Failure		500				{object}	map[string]string
//	@Router			/encounters/{encounterFhirId}/reports [get]
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

// CreateReport godoc
//
//	@Summary		Create a diagnostic report
//	@Description	Creates a new diagnostic report for an encounter
//	@Tags			diagnostic_reports
//	@Accept			json
//	@Produce		json
//	@Param			encounterFhirId	path	string	true	"Encounter FHIR ID"
//	@Param			body			body	CreateDiagnosticReportRequest	true	"Report data"
//	@Success		201				{object}	CreateDiagnosticReportResponse
//	@Failure		400				{object}	map[string]string
//	@Failure		500				{object}	map[string]string
//	@Router			/encounters/{encounterFhirId}/reports [post]
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

type DiagnosticReportResponse struct {
	FhirID          string `json:"fhir_id"`
	EncounterFhirID string `json:"encounter_fhir_id"`
	PatientFhirID   string `json:"patient_fhir_id"`
	ReportDisplay   string `json:"report_display"`
	Status          string `json:"status"`
	Conclusion      string `json:"conclusion"`
	CreatedAt       string `json:"created_at"`
}

type CreateDiagnosticReportRequest struct {
	PatientFhirID string `json:"patient_fhir_id"`
	ReportDisplay string `json:"report_display"`
	Conclusion    string `json:"conclusion"`
}

type CreateDiagnosticReportResponse struct {
	FhirID          string `json:"fhir_id"`
	EncounterFhirID string `json:"encounter_fhir_id"`
	PatientFhirID   string `json:"patient_fhir_id"`
	ReportDisplay   string `json:"report_display"`
	Status          string `json:"status"`
	Conclusion      string `json:"conclusion"`
	CreatedAt       string `json:"created_at"`
}
