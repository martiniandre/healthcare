package timeline

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/healthcare/backend/internal/api/middleware"
	"github.com/healthcare/backend/internal/api/render"
)

const timelineRoute = "GET /api/v1/patients/{patientFhirId}/timeline"

type HTTPHandler struct {
	service Service
}

func NewHTTPHandler(service Service) *HTTPHandler {
	return &HTTPHandler{service: service}
}

func (handler *HTTPHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.Handle(timelineRoute, middleware.RequirePolicy(timelineRoute)(http.HandlerFunc(handler.GetTimeline)))
}

// GetTimeline godoc
//
//	@Summary		Get patient clinical timeline
//	@Description	Returns a chronological feed merging the patient's Encounter, Observation, Condition, AllergyIntolerance, MedicationRequest, DiagnosticReport and ImagingStudy records, ordered by recorded time (newest first) with cursor pagination. When some resource types cannot be loaded, the rest is returned and partial.failedTypes lists the unavailable types. Defaults to the last 6 months.
//	@Tags			timeline
//	@Accept			json
//	@Produce		json
//	@Param			patientFhirId	path		string	true	"Patient FHIR ID"
//	@Param			from			query		string	false	"Window start (ISO date or RFC3339). Defaults to 6 months before now"
//	@Param			to				query		string	false	"Window end (ISO date or RFC3339). Defaults to now"
//	@Param			types			query		string	false	"Comma separated list: encounter,observation,condition,allergy,medication,report,imaging"
//	@Param			cursor			query		string	false	"Opaque cursor returned as nextCursor"
//	@Param			limit			query		int		false	"Page size (1-100)"	default(30)
//	@Success		200				{object}	TimelineResponse
//	@Failure		400				{object}	map[string]string
//	@Failure		404				{object}	map[string]string
//	@Failure		503				{object}	map[string]string
//	@Router			/patients/{patientFhirId}/timeline [get]
func (handler *HTTPHandler) GetTimeline(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
	queryValues := httpRequest.URL.Query()
	input := ListInput{
		PatientFHIRID: httpRequest.PathValue("patientFhirId"),
		From:          queryValues.Get("from"),
		To:            queryValues.Get("to"),
		Types:         queryValues.Get("types"),
		Cursor:        queryValues.Get("cursor"),
		Limit:         queryValues.Get("limit"),
	}

	page, serviceError := handler.service.GetTimeline(httpRequest.Context(), input)
	if serviceError != nil {
		slog.Error("failed to get timeline", "error", serviceError, "patient_fhir_id", input.PatientFHIRID, "request_id", middleware.GetRequestID(httpRequest.Context()))
		render.ErrorFromAppError(httpResponseWriter, serviceError)
		return
	}

	render.JSON(httpResponseWriter, http.StatusOK, toTimelineResponse(page))
}

func toTimelineResponse(page *Page) TimelineResponse {
	items := make([]ItemResponse, 0, len(page.Entries))
	for _, entry := range page.Entries {
		item := ItemResponse{
			ID:         entry.ID(),
			Type:       string(entry.Type),
			OccurredAt: entry.OccurredAt.UTC().Format(time.RFC3339),
			Title:      entry.Title,
			Summary:    entry.Summary,
			Status:     entry.Status,
			ResourceID: entry.ResourceID,
		}
		if entry.ClinicalDate != nil {
			clinicalDate := entry.ClinicalDate.UTC().Format(time.RFC3339)
			item.ClinicalDate = &clinicalDate
		}
		if entry.EncounterID != "" {
			encounterID := entry.EncounterID
			item.EncounterID = &encounterID
		}
		items = append(items, item)
	}

	failedTypes := make([]string, 0, len(page.FailedTypes))
	for _, failedType := range page.FailedTypes {
		failedTypes = append(failedTypes, string(failedType))
	}

	response := TimelineResponse{
		Items:   items,
		Partial: PartialResponse{FailedTypes: failedTypes},
		Window: WindowResponse{
			From: page.Window.From.UTC().Format(time.RFC3339),
			To:   page.Window.To.UTC().Format(time.RFC3339),
		},
	}
	if page.NextCursor != "" {
		nextCursor := page.NextCursor
		response.NextCursor = &nextCursor
	}
	return response
}
