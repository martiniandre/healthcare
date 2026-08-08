package diagnostic_report

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
	mux.Handle("GET /api/v1/encounters/{encounterFhirId}/reports", middleware.RequirePolicy("GET /api/v1/encounters/{encounterFhirId}/reports")(http.HandlerFunc(handler.ListReportsByEncounter)))
	mux.Handle("POST /api/v1/encounters/{encounterFhirId}/reports", middleware.RequirePolicy("POST /api/v1/encounters/{encounterFhirId}/reports")(http.HandlerFunc(handler.CreateReport)))
	mux.Handle("PUT /api/v1/reports/{reportFhirId}", middleware.RequirePolicy("PUT /api/v1/reports/{reportFhirId}")(http.HandlerFunc(handler.UpdateReport)))
	mux.Handle("GET /api/v1/reports/{reportFhirId}/versions", middleware.RequirePolicy("GET /api/v1/reports/{reportFhirId}/versions")(http.HandlerFunc(handler.ListReportVersions)))
	mux.Handle("GET /api/v1/reports/{reportFhirId}/versions/{version}", middleware.RequirePolicy("GET /api/v1/reports/{reportFhirId}/versions/{version}")(http.HandlerFunc(handler.GetReportVersion)))
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
		slog.Error("failed to list reports", "error", reportsErr, "encounter_fhir_id", encounterFhirID, "request_id", middleware.GetRequestID(httpRequest.Context()))
		render.ErrorFromAppError(httpResponseWriter, reportsErr)
		return
	}

	responseList := make([]DiagnosticReportResponse, 0, len(reportsList))
	for _, report := range reportsList {
		responseList = append(responseList, DiagnosticReportResponse{
			FhirID:          report.FHIRResourceID,
			EncounterFhirID: report.EncounterFHIRID,
			PatientFhirID:   report.PatientFHIRID,
			ReportDisplay:   report.ReportDisplay,
			Status:          report.Status,
			Conclusion:      report.Conclusion,
			Version:         report.Version,
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

	var payload CreateDiagnosticReportRequest

	if payloadDecodeErr := json.NewDecoder(httpRequest.Body).Decode(&payload); payloadDecodeErr != nil {
		render.Error(httpResponseWriter, http.StatusBadRequest, "Payload inválido.")
		return
	}

	createdReport, createErr := handler.service.CreateDiagnosticReport(httpRequest.Context(), CreateDiagnosticReportInput{
		EncounterFHIRID: encounterFhirID,
		PatientFHIRID:   payload.PatientFhirID,
		ReportCode:      payload.ReportCode,
		ReportDisplay:   payload.ReportDisplay,
		Conclusion:      payload.Conclusion,
	})
	if createErr != nil {
		slog.Error("failed to create diagnostic report", "error", createErr, "encounter_fhir_id", encounterFhirID, "request_id", middleware.GetRequestID(httpRequest.Context()))
		render.ErrorFromAppError(httpResponseWriter, createErr)
		return
	}

	render.JSON(httpResponseWriter, http.StatusCreated, DiagnosticReportResponse{
		FhirID:          createdReport.FHIRResourceID,
		EncounterFhirID: createdReport.EncounterFHIRID,
		PatientFhirID:   createdReport.PatientFHIRID,
		ReportDisplay:   createdReport.ReportDisplay,
		Status:          createdReport.Status,
		Conclusion:      createdReport.Conclusion,
		Version:         createdReport.Version,
		CreatedAt:       createdReport.IssuedAt.Format(time.RFC3339),
	})
}

// UpdateReport godoc
//
//	@Summary		Update a diagnostic report
//	@Description	Updates an existing diagnostic report, recording a new version
//	@Tags			diagnostic_reports
//	@Accept			json
//	@Produce		json
//	@Param			reportFhirId	path	string	true	"Report FHIR ID"
//	@Param			body			body	UpdateDiagnosticReportRequest	true	"Report fields to update"
//	@Success		200				{object}	DiagnosticReportResponse
//	@Failure		400				{object}	map[string]string
//	@Failure		404				{object}	map[string]string
//	@Failure		500				{object}	map[string]string
//	@Router			/reports/{reportFhirId} [put]
func (handler *HTTPHandler) UpdateReport(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
	reportFhirID := httpRequest.PathValue("reportFhirId")

	var payload UpdateDiagnosticReportRequest
	if payloadDecodeErr := json.NewDecoder(httpRequest.Body).Decode(&payload); payloadDecodeErr != nil {
		render.Error(httpResponseWriter, http.StatusBadRequest, "Payload inválido.")
		return
	}

	updatedReport, updateErr := handler.service.UpdateDiagnosticReport(httpRequest.Context(), reportFhirID, UpdateDiagnosticReportInput{
		ReportCode:    payload.ReportCode,
		ReportDisplay: payload.ReportDisplay,
		Conclusion:    payload.Conclusion,
		Status:        payload.Status,
	})
	if updateErr != nil {
		slog.Error("failed to update diagnostic report", "error", updateErr, "report_fhir_id", reportFhirID, "request_id", middleware.GetRequestID(httpRequest.Context()))
		render.ErrorFromAppError(httpResponseWriter, updateErr)
		return
	}

	render.JSON(httpResponseWriter, http.StatusOK, DiagnosticReportResponse{
		FhirID:          updatedReport.FHIRResourceID,
		EncounterFhirID: updatedReport.EncounterFHIRID,
		PatientFhirID:   updatedReport.PatientFHIRID,
		ReportDisplay:   updatedReport.ReportDisplay,
		Status:          updatedReport.Status,
		Conclusion:      updatedReport.Conclusion,
		Version:         updatedReport.Version,
		CreatedAt:       updatedReport.IssuedAt.Format(time.RFC3339),
	})
}

// ListReportVersions godoc
//
//	@Summary		List diagnostic report versions
//	@Description	Returns the version history of a diagnostic report
//	@Tags			diagnostic_reports
//	@Accept			json
//	@Produce		json
//	@Param			reportFhirId	path	string	true	"Report FHIR ID"
//	@Success		200				{array}	DiagnosticReportVersionResponse
//	@Failure		500				{object}	map[string]string
//	@Router			/reports/{reportFhirId}/versions [get]
func (handler *HTTPHandler) ListReportVersions(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
	reportFhirID := httpRequest.PathValue("reportFhirId")

	versions, listErr := handler.service.GetDiagnosticReportVersions(httpRequest.Context(), reportFhirID)
	if listErr != nil {
		slog.Error("failed to list report versions", "error", listErr, "report_fhir_id", reportFhirID, "request_id", middleware.GetRequestID(httpRequest.Context()))
		render.ErrorFromAppError(httpResponseWriter, listErr)
		return
	}

	responseList := make([]DiagnosticReportVersionResponse, 0, len(versions))
	for _, versionEntry := range versions {
		responseList = append(responseList, toVersionResponse(versionEntry))
	}

	render.JSON(httpResponseWriter, http.StatusOK, responseList)
}

// GetReportVersion godoc
//
//	@Summary		Get a specific diagnostic report version
//	@Description	Returns a specific version snapshot of a diagnostic report
//	@Tags			diagnostic_reports
//	@Accept			json
//	@Produce		json
//	@Param			reportFhirId	path	string	true	"Report FHIR ID"
//	@Param			version			path	string	true	"Version number"
//	@Success		200				{object}	DiagnosticReportVersionResponse
//	@Failure		404				{object}	map[string]string
//	@Failure		500				{object}	map[string]string
//	@Router			/reports/{reportFhirId}/versions/{version} [get]
func (handler *HTTPHandler) GetReportVersion(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
	reportFhirID := httpRequest.PathValue("reportFhirId")
	version := httpRequest.PathValue("version")

	versionEntry, getErr := handler.service.GetDiagnosticReportVersion(httpRequest.Context(), reportFhirID, version)
	if getErr != nil {
		slog.Error("failed to get report version", "error", getErr, "report_fhir_id", reportFhirID, "version", version, "request_id", middleware.GetRequestID(httpRequest.Context()))
		render.ErrorFromAppError(httpResponseWriter, getErr)
		return
	}

	render.JSON(httpResponseWriter, http.StatusOK, toVersionResponse(versionEntry))
}

func toVersionResponse(versionEntry *DiagnosticReportVersion) DiagnosticReportVersionResponse {
	var changedBy *string
	if versionEntry.ChangedBy != nil {
		changedByValue := versionEntry.ChangedBy.String()
		changedBy = &changedByValue
	}
	return DiagnosticReportVersionResponse{
		Version:   versionEntry.Version,
		Snapshot:  versionEntry.Snapshot,
		ChangedBy: changedBy,
		ChangedAt: versionEntry.ChangedAt.Format(time.RFC3339),
	}
}

type DiagnosticReportResponse struct {
	FhirID          string `json:"fhir_id"`
	EncounterFhirID string `json:"encounter_fhir_id"`
	PatientFhirID   string `json:"patient_fhir_id"`
	ReportDisplay   string `json:"report_display"`
	Status          string `json:"status"`
	Conclusion      string `json:"conclusion"`
	Version         string `json:"version"`
	CreatedAt       string `json:"created_at"`
}

type CreateDiagnosticReportRequest struct {
	PatientFhirID string `json:"patient_fhir_id"`
	ReportCode    string `json:"report_code"`
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
	Version         string `json:"version"`
	CreatedAt       string `json:"created_at"`
}

type UpdateDiagnosticReportRequest struct {
	ReportCode    *string `json:"report_code"`
	ReportDisplay *string `json:"report_display"`
	Conclusion    *string `json:"conclusion"`
	Status        *string `json:"status"`
}

type DiagnosticReportVersionResponse struct {
	Version   string          `json:"version"`
	Snapshot  json.RawMessage `json:"snapshot"`
	ChangedBy *string         `json:"changed_by"`
	ChangedAt string          `json:"changed_at"`
}
