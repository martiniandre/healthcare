package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/healthcare/backend/internal/modules/auth"
	"github.com/healthcare/backend/internal/shared/role"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type stubAuditRecorder struct {
	recordedEntries []AuditEntry
}

func (stub *stubAuditRecorder) RecordHTTPAudit(contextVal context.Context, auditEntry AuditEntry) error {
	stub.recordedEntries = append(stub.recordedEntries, auditEntry)
	return nil
}

func installStubAuditRecorder(testingInstance *testing.T) *stubAuditRecorder {
	recorderStub := &stubAuditRecorder{}
	SetHTTPAuditRecorder(recorderStub)
	testingInstance.Cleanup(func() {
		SetHTTPAuditRecorder(nil)
	})
	return recorderStub
}

func buildAuditedPipeline(handlerToPipe http.Handler) http.Handler {
	return RequestID(AuditTrail(handlerToPipe))
}

func TestAuditTrail_RecordsGrantedClinicalRequest(testingInstance *testing.T) {
	adminToken := initializeJWTForMiddlewareTest(testingInstance)
	recorderStub := installStubAuditRecorder(testingInstance)

	httpRequest := httptest.NewRequest(http.MethodPost, "/api/v1/patients", nil)
	httpRequest.AddCookie(&http.Cookie{Name: "token", Value: adminToken})
	httpResponseRecorder := httptest.NewRecorder()

	auditedPipeline := buildAuditedPipeline(http.HandlerFunc(func(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
		httpResponseWriter.WriteHeader(http.StatusCreated)
	}))

	auditedPipeline.ServeHTTP(httpResponseRecorder, httpRequest)

	require.Len(testingInstance, recorderStub.recordedEntries, 1)
	recordedEntry := recorderStub.recordedEntries[0]
	assert.True(testingInstance, recordedEntry.AccessGranted)
	assert.Equal(testingInstance, "patient", recordedEntry.ResourceType)
	assert.Equal(testingInstance, "create", recordedEntry.Action)
	assert.Equal(testingInstance, "user-1", recordedEntry.CallerUserID)
	assert.Equal(testingInstance, string(role.RoleAdmin), recordedEntry.CallerRole)
	assert.NotEmpty(testingInstance, recordedEntry.CorrelationID)
	assert.Equal(testingInstance, "POST /api/v1/patients", recordedEntry.Route)
}

func TestAuditTrail_RecordsDeniedRequest(testingInstance *testing.T) {
	adminToken := initializeJWTForMiddlewareTest(testingInstance)
	recorderStub := installStubAuditRecorder(testingInstance)

	httpRequest := httptest.NewRequest(http.MethodGet, "/api/v1/patients", nil)
	httpRequest.AddCookie(&http.Cookie{Name: "token", Value: adminToken})
	httpResponseRecorder := httptest.NewRecorder()

	auditedPipeline := buildAuditedPipeline(http.HandlerFunc(func(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
		httpResponseWriter.WriteHeader(http.StatusForbidden)
	}))

	auditedPipeline.ServeHTTP(httpResponseRecorder, httpRequest)

	require.Len(testingInstance, recorderStub.recordedEntries, 1)
	assert.False(testingInstance, recorderStub.recordedEntries[0].AccessGranted)
	assert.Equal(testingInstance, http.StatusForbidden, httpResponseRecorder.Code)
}

func TestAuditTrail_ResolvesResourceIDFromPath(testingInstance *testing.T) {
	adminToken := initializeJWTForMiddlewareTest(testingInstance)
	recorderStub := installStubAuditRecorder(testingInstance)

	httpRequest := httptest.NewRequest(http.MethodGet, "/api/v1/patients/fhir-patient-123", nil)
	httpRequest.AddCookie(&http.Cookie{Name: "token", Value: adminToken})
	httpResponseRecorder := httptest.NewRecorder()

	auditedPipeline := buildAuditedPipeline(http.HandlerFunc(func(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
		httpResponseWriter.WriteHeader(http.StatusOK)
	}))

	auditedPipeline.ServeHTTP(httpResponseRecorder, httpRequest)

	require.Len(testingInstance, recorderStub.recordedEntries, 1)
	recordedEntry := recorderStub.recordedEntries[0]
	assert.Equal(testingInstance, "patient", recordedEntry.ResourceType)
	assert.Equal(testingInstance, "read", recordedEntry.Action)
	assert.Equal(testingInstance, "fhir-patient-123", recordedEntry.ResourceID)
}

func TestAuditTrail_SkipsHealthEndpoint(testingInstance *testing.T) {
	recorderStub := installStubAuditRecorder(testingInstance)

	httpRequest := httptest.NewRequest(http.MethodGet, "/health", nil)
	httpResponseRecorder := httptest.NewRecorder()

	auditedPipeline := buildAuditedPipeline(http.HandlerFunc(func(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
		httpResponseWriter.WriteHeader(http.StatusOK)
	}))

	auditedPipeline.ServeHTTP(httpResponseRecorder, httpRequest)

	assert.Empty(testingInstance, recorderStub.recordedEntries)
}

func TestAuditTrail_SkipsSwaggerPath(testingInstance *testing.T) {
	recorderStub := installStubAuditRecorder(testingInstance)

	httpRequest := httptest.NewRequest(http.MethodGet, "/swagger/index.html", nil)
	httpResponseRecorder := httptest.NewRecorder()

	auditedPipeline := buildAuditedPipeline(http.HandlerFunc(func(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
		httpResponseWriter.WriteHeader(http.StatusOK)
	}))

	auditedPipeline.ServeHTTP(httpResponseRecorder, httpRequest)

	assert.Empty(testingInstance, recorderStub.recordedEntries)
}

func TestAuditTrail_NoRecorderSkipsAuditWrite(testingInstance *testing.T) {
	adminToken := initializeJWTForMiddlewareTest(testingInstance)

	httpRequest := httptest.NewRequest(http.MethodGet, "/api/v1/patients", nil)
	httpRequest.AddCookie(&http.Cookie{Name: "token", Value: adminToken})
	httpResponseRecorder := httptest.NewRecorder()

	auditedPipeline := buildAuditedPipeline(http.HandlerFunc(func(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
		httpResponseWriter.WriteHeader(http.StatusOK)
	}))

	auditedPipeline.ServeHTTP(httpResponseRecorder, httpRequest)

	assert.Equal(testingInstance, http.StatusOK, httpResponseRecorder.Code)
}

func TestMatchHTTPRouteMetadata(testingInstance *testing.T) {
	testCases := []struct {
		routeName      string
		httpMethod     string
		requestPath    string
		expectedType   string
		expectedID     string
		expectedAction string
	}{
		{routeName: "list patients", httpMethod: http.MethodGet, requestPath: "/api/v1/patients", expectedType: "patient", expectedAction: "list"},
		{routeName: "create observation", httpMethod: http.MethodPost, requestPath: "/api/v1/encounters/enc-1/observations", expectedType: "observation", expectedAction: "create", expectedID: "enc-1"},
		{routeName: "update report", httpMethod: http.MethodPut, requestPath: "/api/v1/reports/report-77", expectedType: "diagnostic_report", expectedAction: "update", expectedID: "report-77"},
		{routeName: "my appointments", httpMethod: http.MethodGet, requestPath: "/api/v1/appointments/my", expectedType: "appointment", expectedAction: "read"},
		{routeName: "appointment by id", httpMethod: http.MethodGet, requestPath: "/api/v1/appointments/appt-9", expectedType: "appointment", expectedAction: "read", expectedID: "appt-9"},
		{routeName: "portal report read", httpMethod: http.MethodGet, requestPath: "/api/v1/portal/reports", expectedType: "diagnostic_report", expectedAction: "read"},
		{routeName: "login", httpMethod: http.MethodPost, requestPath: "/api/v1/auth/login", expectedType: "auth", expectedAction: "login"},
		{routeName: "cancel appointment", httpMethod: http.MethodPost, requestPath: "/api/v1/appointments/appt-9/cancel", expectedType: "appointment", expectedAction: "cancel", expectedID: "appt-9"},
		{routeName: "stream notifications", httpMethod: http.MethodGet, requestPath: "/api/v1/notifications/stream", expectedType: "notification", expectedAction: "read"},
		{routeName: "unregistered route", httpMethod: http.MethodGet, requestPath: "/api/v1/unregistered", expectedType: "", expectedAction: ""},
	}

	for _, testCase := range testCases {
		resourceType, resourceID, action := matchHTTPRouteMetadata(testCase.httpMethod, testCase.requestPath)
		assert.Equal(testingInstance, testCase.expectedType, resourceType, "%s resource type", testCase.routeName)
		assert.Equal(testingInstance, testCase.expectedID, resourceID, "%s resource id", testCase.routeName)
		assert.Equal(testingInstance, testCase.expectedAction, action, "%s action", testCase.routeName)
	}
}

func TestResolveCallerIdentity_ValidToken(testingInstance *testing.T) {
	require.NoError(testingInstance, auth.InitJWT("middleware-test-secret"))
	adminToken, tokenError := auth.GenerateJWT("user-42", string(role.RoleNurse), "nurse@healthcare.com")
	require.NoError(testingInstance, tokenError)

	httpRequest := httptest.NewRequest(http.MethodGet, "/api/v1/patients", nil)
	httpRequest.AddCookie(&http.Cookie{Name: "token", Value: adminToken})

	callerUserID, callerRole := resolveCallerIdentity(httpRequest)
	assert.Equal(testingInstance, "user-42", callerUserID)
	assert.Equal(testingInstance, string(role.RoleNurse), callerRole)
}

func TestResolveCallerIdentity_MissingCookie(testingInstance *testing.T) {
	httpRequest := httptest.NewRequest(http.MethodGet, "/api/v1/patients", nil)

	callerUserID, callerRole := resolveCallerIdentity(httpRequest)
	assert.Empty(testingInstance, callerUserID)
	assert.Empty(testingInstance, callerRole)
}
