package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/healthcare/backend/internal/modules/auth"
	"github.com/healthcare/backend/internal/shared/ctxkeys"
	"github.com/healthcare/backend/internal/shared/role"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func initializeJWTForMiddlewareTest(testingInstance *testing.T) string {
	require.NoError(testingInstance, auth.InitJWT("middleware-test-secret"))
	token, tokenErr := auth.GenerateJWT("user-1", string(role.RoleAdmin), "admin@healthcare.com")
	require.NoError(testingInstance, tokenErr)
	return token
}

func TestRequirePolicy_AllowsAuthorizedRole(testingInstance *testing.T) {
	adminToken := initializeJWTForMiddlewareTest(testingInstance)

	httpRequest := httptest.NewRequest(http.MethodGet, "/api/v1/patients", nil)
	httpRequest.AddCookie(&http.Cookie{Name: "token", Value: adminToken})
	httpResponseRecorder := httptest.NewRecorder()

	capturedRole := ""
	guardedHandler := RequirePolicy("GET /api/v1/patients")(http.HandlerFunc(func(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
		capturedRole, _ = httpRequest.Context().Value(ctxkeys.RoleKey).(string)
		httpResponseWriter.WriteHeader(http.StatusOK)
	}))

	guardedHandler.ServeHTTP(httpResponseRecorder, httpRequest)

	assert.Equal(testingInstance, http.StatusOK, httpResponseRecorder.Code)
	assert.Equal(testingInstance, string(role.RoleAdmin), capturedRole)
}

func TestRequirePolicy_DeniesUnauthorizedRole(testingInstance *testing.T) {
	require.NoError(testingInstance, auth.InitJWT("middleware-test-secret"))
	patientToken, tokenErr := auth.GenerateJWT("patient-1", string(role.RolePatient), "patient@healthcare.com")
	require.NoError(testingInstance, tokenErr)

	httpRequest := httptest.NewRequest(http.MethodGet, "/api/v1/audit-logs", nil)
	httpRequest.AddCookie(&http.Cookie{Name: "token", Value: patientToken})
	httpResponseRecorder := httptest.NewRecorder()

	guardedHandler := RequirePolicy("GET /api/v1/audit-logs")(http.HandlerFunc(func(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
		httpResponseWriter.WriteHeader(http.StatusOK)
	}))

	guardedHandler.ServeHTTP(httpResponseRecorder, httpRequest)

	assert.Equal(testingInstance, http.StatusForbidden, httpResponseRecorder.Code)
}

func TestRequirePolicy_DeniesUndefinedRoute(testingInstance *testing.T) {
	adminToken := initializeJWTForMiddlewareTest(testingInstance)

	httpRequest := httptest.NewRequest(http.MethodGet, "/api/v1/unregistered-route", nil)
	httpRequest.AddCookie(&http.Cookie{Name: "token", Value: adminToken})
	httpResponseRecorder := httptest.NewRecorder()

	guardedHandler := RequirePolicy("GET /api/v1/unregistered-route")(http.HandlerFunc(func(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
		httpResponseWriter.WriteHeader(http.StatusOK)
	}))

	guardedHandler.ServeHTTP(httpResponseRecorder, httpRequest)

	assert.Equal(testingInstance, http.StatusForbidden, httpResponseRecorder.Code)
}

func TestRequirePolicy_DeniesMissingToken(testingInstance *testing.T) {
	httpRequest := httptest.NewRequest(http.MethodGet, "/api/v1/patients", nil)
	httpResponseRecorder := httptest.NewRecorder()

	guardedHandler := RequirePolicy("GET /api/v1/patients")(http.HandlerFunc(func(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
		httpResponseWriter.WriteHeader(http.StatusOK)
	}))

	guardedHandler.ServeHTTP(httpResponseRecorder, httpRequest)

	assert.Equal(testingInstance, http.StatusUnauthorized, httpResponseRecorder.Code)
}

func TestValidatePolicyAuth_PropagatesIdentityContext(testingInstance *testing.T) {
	adminToken := initializeJWTForMiddlewareTest(testingInstance)

	httpRequest := httptest.NewRequest(http.MethodGet, "/api/v1/patients", nil)
	httpRequest.AddCookie(&http.Cookie{Name: "token", Value: adminToken})
	httpResponseRecorder := httptest.NewRecorder()

	authenticatedContext, authorizationPassed := ValidatePolicyAuth(httpResponseRecorder, httpRequest, "GET /api/v1/patients")

	assert.True(testingInstance, authorizationPassed)
	userID, _ := authenticatedContext.Value(ctxkeys.UserIDKey).(string)
	assert.Equal(testingInstance, "user-1", userID)
	callerRole, _ := authenticatedContext.Value(ctxkeys.RoleKey).(string)
	assert.Equal(testingInstance, string(role.RoleAdmin), callerRole)
}
