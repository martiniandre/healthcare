package timeline

import (
	"log/slog"
	"net/http"
	"strconv"
	"strings"
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
	mux.Handle("GET /api/v1/patients/{patientFhirId}/timeline", middleware.RequirePolicy("GET /api/v1/patients/{patientFhirId}/timeline")(http.HandlerFunc(handler.GetTimeline)))
}

func (handler *HTTPHandler) GetTimeline(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
	patientFHIRID := httpRequest.PathValue("patientFhirId")

	timelineFilter, parseErr := parseTimelineFilter(httpRequest.URL.Query())
	if parseErr != nil {
		render.ErrorFromAppError(httpResponseWriter, parseErr)
		return
	}

	timelinePage, serviceErr := handler.service.GetTimeline(httpRequest.Context(), patientFHIRID, timelineFilter)
	if serviceErr != nil {
		slog.Error("failed to get clinical timeline",
			"error", serviceErr,
			"patient_id", patientFHIRID,
			"request_id", middleware.GetRequestID(httpRequest.Context()),
		)
		render.ErrorFromAppError(httpResponseWriter, serviceErr)
		return
	}

	render.JSON(httpResponseWriter, http.StatusOK, timelinePage)
}

func parseTimelineFilter(queryParams map[string][]string) (TimelineFilter, error) {
	timelineFilter := TimelineFilter{}

	if limitValues, hasLimit := queryParams["limit"]; hasLimit && len(limitValues) > 0 && limitValues[0] != "" {
		parsedLimit, limitErr := strconv.Atoi(limitValues[0])
		if limitErr != nil || parsedLimit <= 0 || parsedLimit > MaxPageSize {
			return TimelineFilter{}, apperrors.InvalidArgument("invalid limit parameter", map[string]string{"limit": "must be an integer between 1 and 100"})
		}
		timelineFilter.Limit = parsedLimit
	}

	if beforeValues, hasBefore := queryParams["before"]; hasBefore && len(beforeValues) > 0 && beforeValues[0] != "" {
		parsedBefore, beforeErr := time.Parse(time.RFC3339, beforeValues[0])
		if beforeErr != nil {
			return TimelineFilter{}, apperrors.InvalidArgument("invalid before parameter", map[string]string{"before": "must be an RFC3339 timestamp"})
		}
		timelineFilter.Before = &parsedBefore
	}

	if typeValues, hasTypes := queryParams["types"]; hasTypes && len(typeValues) > 0 && typeValues[0] != "" {
		requestedTypes := strings.Split(typeValues[0], ",")
		parsedTypes := make([]string, 0, len(requestedTypes))
		for _, requestedType := range requestedTypes {
			trimmedType := strings.TrimSpace(requestedType)
			if trimmedType != "" {
				parsedTypes = append(parsedTypes, trimmedType)
			}
		}
		timelineFilter.Types = parsedTypes
	}

	if statusValues, hasStatus := queryParams["status"]; hasStatus && len(statusValues) > 0 {
		timelineFilter.Status = statusValues[0]
	}

	return timelineFilter, nil
}
