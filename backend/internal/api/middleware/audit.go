package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	"github.com/healthcare/backend/internal/modules/auth"
)

type AuditEntry struct {
	CorrelationID string
	CallerUserID  string
	CallerRole    string
	Route         string
	AccessGranted bool
	ResourceType  string
	ResourceID    string
	Action        string
}

type HTTPAuditRecorder interface {
	RecordHTTPAudit(contextVal context.Context, auditEntry AuditEntry) error
}

var globalHTTPAuditRecorder HTTPAuditRecorder

func SetHTTPAuditRecorder(auditRecorder HTTPAuditRecorder) {
	globalHTTPAuditRecorder = auditRecorder
}

func AuditTrail(next http.Handler) http.Handler {
	return http.HandlerFunc(func(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
		requestRecorder := &statusRecorder{ResponseWriter: httpResponseWriter, statusCode: http.StatusOK}
		next.ServeHTTP(requestRecorder, httpRequest)

		if globalHTTPAuditRecorder == nil {
			return
		}
		if shouldSkipHTTPAudit(httpRequest.Method, httpRequest.URL.Path) {
			return
		}

		callerUserID, callerRole := resolveCallerIdentity(httpRequest)
		resourceType, resourceID, action := matchHTTPRouteMetadata(httpRequest.Method, httpRequest.URL.Path)

		auditError := globalHTTPAuditRecorder.RecordHTTPAudit(context.Background(), AuditEntry{
			CorrelationID: GetRequestID(httpRequest.Context()),
			CallerUserID:  callerUserID,
			CallerRole:    callerRole,
			Route:         httpRequest.Method + " " + httpRequest.URL.Path,
			AccessGranted: requestRecorder.statusCode < http.StatusBadRequest,
			ResourceType:  resourceType,
			ResourceID:    resourceID,
			Action:        action,
		})
		if auditError != nil {
			slog.Error("failed to persist http audit log", "error", auditError, "route", httpRequest.URL.Path)
		}
	})
}

func shouldSkipHTTPAudit(httpMethod string, requestPath string) bool {
	if httpMethod == http.MethodOptions {
		return true
	}
	if httpMethod == http.MethodPost && requestPath == "/api/v1/auth/login" {
		return true
	}
	if requestPath == "/health" || requestPath == "/api/v1/health" {
		return true
	}
	if strings.HasPrefix(requestPath, "/swagger/") {
		return true
	}
	return false
}

func resolveCallerIdentity(httpRequest *http.Request) (string, string) {
	cookie, cookieError := httpRequest.Cookie("token")
	if cookieError != nil {
		return "", ""
	}
	claims, jwtValidationError := auth.ValidateJWT(cookie.Value)
	if jwtValidationError != nil {
		return "", ""
	}
	callerUserID, _ := claims["user_id"].(string)
	callerRole, _ := claims["role"].(string)
	return callerUserID, callerRole
}

type httpRouteAuditMetadata struct {
	resourceType    string
	action          string
	resourceIDParam string
}

var httpRouteAuditMetadataList = []struct {
	routePattern string
	httpRouteAuditMetadata
}{
	{routePattern: "GET /api/v1/auth/me", httpRouteAuditMetadata: httpRouteAuditMetadata{resourceType: "auth", action: "read"}},
	{routePattern: "POST /api/v1/auth/logout", httpRouteAuditMetadata: httpRouteAuditMetadata{resourceType: "auth", action: "logout"}},

	{routePattern: "GET /api/v1/audit-logs", httpRouteAuditMetadata: httpRouteAuditMetadata{resourceType: "audit_log", action: "list"}},
	{routePattern: "POST /api/v1/audit-logs", httpRouteAuditMetadata: httpRouteAuditMetadata{resourceType: "audit_log", action: "create"}},

	{routePattern: "GET /api/v1/patients", httpRouteAuditMetadata: httpRouteAuditMetadata{resourceType: "patient", action: "list"}},
	{routePattern: "POST /api/v1/patients", httpRouteAuditMetadata: httpRouteAuditMetadata{resourceType: "patient", action: "create"}},
	{routePattern: "GET /api/v1/patients/{patientFhirId}", httpRouteAuditMetadata: httpRouteAuditMetadata{resourceType: "patient", action: "read", resourceIDParam: "patientFhirId"}},

	{routePattern: "GET /api/v1/patients/{patientFhirId}/encounters", httpRouteAuditMetadata: httpRouteAuditMetadata{resourceType: "encounter", action: "list", resourceIDParam: "patientFhirId"}},
	{routePattern: "POST /api/v1/patients/{patientFhirId}/encounters", httpRouteAuditMetadata: httpRouteAuditMetadata{resourceType: "encounter", action: "create", resourceIDParam: "patientFhirId"}},
	{routePattern: "GET /api/v1/encounters/{encounterFhirId}", httpRouteAuditMetadata: httpRouteAuditMetadata{resourceType: "encounter", action: "read", resourceIDParam: "encounterFhirId"}},
	{routePattern: "PUT /api/v1/encounters/{encounterFhirId}", httpRouteAuditMetadata: httpRouteAuditMetadata{resourceType: "encounter", action: "update", resourceIDParam: "encounterFhirId"}},
	{routePattern: "DELETE /api/v1/encounters/{encounterFhirId}", httpRouteAuditMetadata: httpRouteAuditMetadata{resourceType: "encounter", action: "delete", resourceIDParam: "encounterFhirId"}},

	{routePattern: "GET /api/v1/patients/{patientFhirId}/observations", httpRouteAuditMetadata: httpRouteAuditMetadata{resourceType: "observation", action: "list", resourceIDParam: "patientFhirId"}},
	{routePattern: "GET /api/v1/encounters/{encounterFhirId}/observations", httpRouteAuditMetadata: httpRouteAuditMetadata{resourceType: "observation", action: "list", resourceIDParam: "encounterFhirId"}},
	{routePattern: "POST /api/v1/encounters/{encounterFhirId}/observations", httpRouteAuditMetadata: httpRouteAuditMetadata{resourceType: "observation", action: "create", resourceIDParam: "encounterFhirId"}},
	{routePattern: "POST /api/v1/encounters/{encounterFhirId}/observations/batch", httpRouteAuditMetadata: httpRouteAuditMetadata{resourceType: "observation", action: "create", resourceIDParam: "encounterFhirId"}},
	{routePattern: "PUT /api/v1/observations/{observationFhirId}", httpRouteAuditMetadata: httpRouteAuditMetadata{resourceType: "observation", action: "update", resourceIDParam: "observationFhirId"}},
	{routePattern: "DELETE /api/v1/observations/{observationFhirId}", httpRouteAuditMetadata: httpRouteAuditMetadata{resourceType: "observation", action: "delete", resourceIDParam: "observationFhirId"}},

	{routePattern: "GET /api/v1/patients/{patientFhirId}/conditions", httpRouteAuditMetadata: httpRouteAuditMetadata{resourceType: "condition", action: "list", resourceIDParam: "patientFhirId"}},
	{routePattern: "POST /api/v1/patients/{patientFhirId}/conditions", httpRouteAuditMetadata: httpRouteAuditMetadata{resourceType: "condition", action: "create", resourceIDParam: "patientFhirId"}},
	{routePattern: "PUT /api/v1/patients/{patientFhirId}/conditions/{conditionFhirId}", httpRouteAuditMetadata: httpRouteAuditMetadata{resourceType: "condition", action: "update", resourceIDParam: "conditionFhirId"}},
	{routePattern: "DELETE /api/v1/patients/{patientFhirId}/conditions/{conditionFhirId}", httpRouteAuditMetadata: httpRouteAuditMetadata{resourceType: "condition", action: "delete", resourceIDParam: "conditionFhirId"}},

	{routePattern: "GET /api/v1/patients/{patientFhirId}/allergies", httpRouteAuditMetadata: httpRouteAuditMetadata{resourceType: "allergy", action: "list", resourceIDParam: "patientFhirId"}},
	{routePattern: "POST /api/v1/patients/{patientFhirId}/allergies", httpRouteAuditMetadata: httpRouteAuditMetadata{resourceType: "allergy", action: "create", resourceIDParam: "patientFhirId"}},
	{routePattern: "PUT /api/v1/patients/{patientFhirId}/allergies/{allergyFhirId}", httpRouteAuditMetadata: httpRouteAuditMetadata{resourceType: "allergy", action: "update", resourceIDParam: "allergyFhirId"}},
	{routePattern: "DELETE /api/v1/patients/{patientFhirId}/allergies/{allergyFhirId}", httpRouteAuditMetadata: httpRouteAuditMetadata{resourceType: "allergy", action: "delete", resourceIDParam: "allergyFhirId"}},

	{routePattern: "GET /api/v1/encounters/{encounterFhirId}/medications", httpRouteAuditMetadata: httpRouteAuditMetadata{resourceType: "medication", action: "list", resourceIDParam: "encounterFhirId"}},
	{routePattern: "POST /api/v1/encounters/{encounterFhirId}/medications", httpRouteAuditMetadata: httpRouteAuditMetadata{resourceType: "medication", action: "create", resourceIDParam: "encounterFhirId"}},

	{routePattern: "GET /api/v1/encounters/{encounterFhirId}/reports", httpRouteAuditMetadata: httpRouteAuditMetadata{resourceType: "diagnostic_report", action: "list", resourceIDParam: "encounterFhirId"}},
	{routePattern: "POST /api/v1/encounters/{encounterFhirId}/reports", httpRouteAuditMetadata: httpRouteAuditMetadata{resourceType: "diagnostic_report", action: "create", resourceIDParam: "encounterFhirId"}},
	{routePattern: "PUT /api/v1/reports/{reportFhirId}", httpRouteAuditMetadata: httpRouteAuditMetadata{resourceType: "diagnostic_report", action: "update", resourceIDParam: "reportFhirId"}},
	{routePattern: "GET /api/v1/reports/{reportFhirId}/versions", httpRouteAuditMetadata: httpRouteAuditMetadata{resourceType: "diagnostic_report", action: "read", resourceIDParam: "reportFhirId"}},
	{routePattern: "GET /api/v1/reports/{reportFhirId}/versions/{version}", httpRouteAuditMetadata: httpRouteAuditMetadata{resourceType: "diagnostic_report", action: "read", resourceIDParam: "reportFhirId"}},

	{routePattern: "GET /api/v1/patients/{patientFhirId}/studies", httpRouteAuditMetadata: httpRouteAuditMetadata{resourceType: "imaging_study", action: "list", resourceIDParam: "patientFhirId"}},
	{routePattern: "POST /api/v1/patients/{patientFhirId}/studies", httpRouteAuditMetadata: httpRouteAuditMetadata{resourceType: "imaging_study", action: "create", resourceIDParam: "patientFhirId"}},
	{routePattern: "GET /api/v1/studies/{studyId}", httpRouteAuditMetadata: httpRouteAuditMetadata{resourceType: "imaging_study", action: "read", resourceIDParam: "studyId"}},

	{routePattern: "GET /api/v1/exam-analyses", httpRouteAuditMetadata: httpRouteAuditMetadata{resourceType: "exam_analysis", action: "list"}},
	{routePattern: "POST /api/v1/exam-analyses", httpRouteAuditMetadata: httpRouteAuditMetadata{resourceType: "exam_analysis", action: "create"}},
	{routePattern: "GET /api/v1/exam-analyses/{analysisId}", httpRouteAuditMetadata: httpRouteAuditMetadata{resourceType: "exam_analysis", action: "read", resourceIDParam: "analysisId"}},
	{routePattern: "DELETE /api/v1/exam-analyses/{analysisId}", httpRouteAuditMetadata: httpRouteAuditMetadata{resourceType: "exam_analysis", action: "delete", resourceIDParam: "analysisId"}},

	{routePattern: "POST /api/v1/appointments", httpRouteAuditMetadata: httpRouteAuditMetadata{resourceType: "appointment", action: "create"}},
	{routePattern: "GET /api/v1/appointments", httpRouteAuditMetadata: httpRouteAuditMetadata{resourceType: "appointment", action: "list"}},
	{routePattern: "GET /api/v1/appointments/my", httpRouteAuditMetadata: httpRouteAuditMetadata{resourceType: "appointment", action: "read"}},
	{routePattern: "GET /api/v1/appointments/{appointmentId}", httpRouteAuditMetadata: httpRouteAuditMetadata{resourceType: "appointment", action: "read", resourceIDParam: "appointmentId"}},
	{routePattern: "POST /api/v1/appointments/{appointmentId}/cancel", httpRouteAuditMetadata: httpRouteAuditMetadata{resourceType: "appointment", action: "cancel", resourceIDParam: "appointmentId"}},
	{routePattern: "PUT /api/v1/appointments/{appointmentId}", httpRouteAuditMetadata: httpRouteAuditMetadata{resourceType: "appointment", action: "update", resourceIDParam: "appointmentId"}},

	{routePattern: "POST /api/v1/schedule/unavailability", httpRouteAuditMetadata: httpRouteAuditMetadata{resourceType: "unavailability", action: "create"}},
	{routePattern: "GET /api/v1/schedule/unavailability", httpRouteAuditMetadata: httpRouteAuditMetadata{resourceType: "unavailability", action: "list"}},
	{routePattern: "DELETE /api/v1/schedule/unavailability/{unavailabilityId}", httpRouteAuditMetadata: httpRouteAuditMetadata{resourceType: "unavailability", action: "delete", resourceIDParam: "unavailabilityId"}},

	{routePattern: "GET /api/v1/staff/employees", httpRouteAuditMetadata: httpRouteAuditMetadata{resourceType: "employee", action: "list"}},
	{routePattern: "POST /api/v1/staff/employees", httpRouteAuditMetadata: httpRouteAuditMetadata{resourceType: "employee", action: "create"}},

	{routePattern: "GET /api/v1/telemetry/rooms", httpRouteAuditMetadata: httpRouteAuditMetadata{resourceType: "telemetry", action: "list"}},
	{routePattern: "POST /api/v1/telemetry/rooms/{roomId}/unlock", httpRouteAuditMetadata: httpRouteAuditMetadata{resourceType: "telemetry", action: "unlock", resourceIDParam: "roomId"}},
	{routePattern: "GET /api/v1/telemetry/rooms/{roomId}/beds", httpRouteAuditMetadata: httpRouteAuditMetadata{resourceType: "telemetry", action: "list", resourceIDParam: "roomId"}},
	{routePattern: "POST /api/v1/telemetry/beds/{bedId}/condition", httpRouteAuditMetadata: httpRouteAuditMetadata{resourceType: "telemetry", action: "update", resourceIDParam: "bedId"}},

	{routePattern: "GET /api/v1/notifications", httpRouteAuditMetadata: httpRouteAuditMetadata{resourceType: "notification", action: "list"}},
	{routePattern: "POST /api/v1/notifications/{notificationId}/read", httpRouteAuditMetadata: httpRouteAuditMetadata{resourceType: "notification", action: "mark_read", resourceIDParam: "notificationId"}},
	{routePattern: "GET /api/v1/notifications/unread-count", httpRouteAuditMetadata: httpRouteAuditMetadata{resourceType: "notification", action: "read"}},
	{routePattern: "GET /api/v1/notifications/stream", httpRouteAuditMetadata: httpRouteAuditMetadata{resourceType: "notification", action: "read"}},

	{routePattern: "GET /api/v1/portal/dashboard", httpRouteAuditMetadata: httpRouteAuditMetadata{resourceType: "portal", action: "read"}},
	{routePattern: "GET /api/v1/portal/encounters", httpRouteAuditMetadata: httpRouteAuditMetadata{resourceType: "encounter", action: "read"}},
	{routePattern: "GET /api/v1/portal/observations", httpRouteAuditMetadata: httpRouteAuditMetadata{resourceType: "observation", action: "read"}},
	{routePattern: "GET /api/v1/portal/conditions", httpRouteAuditMetadata: httpRouteAuditMetadata{resourceType: "condition", action: "read"}},
	{routePattern: "GET /api/v1/portal/medications", httpRouteAuditMetadata: httpRouteAuditMetadata{resourceType: "medication", action: "read"}},
	{routePattern: "GET /api/v1/portal/reports", httpRouteAuditMetadata: httpRouteAuditMetadata{resourceType: "diagnostic_report", action: "read"}},
	{routePattern: "GET /api/v1/portal/imaging", httpRouteAuditMetadata: httpRouteAuditMetadata{resourceType: "imaging_study", action: "read"}},

	{routePattern: "GET /api/v1/analytics", httpRouteAuditMetadata: httpRouteAuditMetadata{resourceType: "analytics", action: "read"}},
	{routePattern: "GET /api/v1/analytics/dashboard", httpRouteAuditMetadata: httpRouteAuditMetadata{resourceType: "analytics", action: "read"}},
	{routePattern: "GET /api/v1/analytics/dashboard/consultations-per-doctor", httpRouteAuditMetadata: httpRouteAuditMetadata{resourceType: "analytics", action: "read"}},
	{routePattern: "GET /api/v1/analytics/dashboard/occupancy-rate", httpRouteAuditMetadata: httpRouteAuditMetadata{resourceType: "analytics", action: "read"}},
	{routePattern: "GET /api/v1/analytics/dashboard/avg-wait-time", httpRouteAuditMetadata: httpRouteAuditMetadata{resourceType: "analytics", action: "read"}},
	{routePattern: "GET /api/v1/analytics/dashboard/top-diagnoses", httpRouteAuditMetadata: httpRouteAuditMetadata{resourceType: "analytics", action: "read"}},
}

func matchHTTPRouteMetadata(httpMethod string, requestPath string) (string, string, string) {
	bestMatchScore := -1
	bestResourceType := ""
	bestResourceID := ""
	bestAction := ""
	for _, routeEntry := range httpRouteAuditMetadataList {
		entryMethod, entryPath, patternValid := splitRoutePattern(routeEntry.routePattern)
		if !patternValid || entryMethod != httpMethod {
			continue
		}
		literalCount, resourceID, matches := matchRouteSegments(entryPath, requestPath, routeEntry.resourceIDParam)
		if !matches {
			continue
		}
		if literalCount > bestMatchScore {
			bestMatchScore = literalCount
			bestResourceType = routeEntry.resourceType
			bestResourceID = resourceID
			bestAction = routeEntry.action
		}
	}
	if bestMatchScore == -1 {
		return "", "", ""
	}
	return bestResourceType, bestResourceID, bestAction
}

func splitRoutePattern(routePattern string) (string, string, bool) {
	spaceIndex := strings.Index(routePattern, " ")
	if spaceIndex == -1 {
		return "", "", false
	}
	return routePattern[:spaceIndex], routePattern[spaceIndex+1:], true
}

func matchRouteSegments(routePatternPath string, requestPath string, resourceIDParam string) (int, string, bool) {
	patternSegments := strings.Split(strings.Trim(routePatternPath, "/"), "/")
	requestSegments := strings.Split(strings.Trim(requestPath, "/"), "/")
	if len(patternSegments) != len(requestSegments) {
		return 0, "", false
	}
	literalCount := 0
	resourceID := ""
	for segmentIndex, patternSegment := range patternSegments {
		if strings.HasPrefix(patternSegment, "{") && strings.HasSuffix(patternSegment, "}") {
			paramName := patternSegment[1 : len(patternSegment)-1]
			if paramName == resourceIDParam {
				resourceID = requestSegments[segmentIndex]
			}
			continue
		}
		if patternSegment != requestSegments[segmentIndex] {
			return 0, "", false
		}
		literalCount++
	}
	return literalCount, resourceID, true
}
