package tests

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/healthcare/backend/internal/modules/auth"
	"github.com/healthcare/backend/internal/modules/auth/mocks"
	"github.com/healthcare/backend/internal/shared/role"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newHTTPHandlerFixture() (*mocks.MockService, *http.ServeMux) {
	mockService := mocks.NewMockService()
	mockService.User = &auth.User{ID: uuid.New(), Email: "nurse@example.com", Role: role.RoleNurse}
	mockService.Token = "access-token-value"
	mockService.RefreshToken = "refresh-token-value"
	httpHandler := auth.NewHTTPHandler(mockService, true)
	serveMux := http.NewServeMux()
	httpHandler.RegisterRoutes(serveMux)
	return mockService, serveMux
}

func findResponseCookie(responseRecorder *httptest.ResponseRecorder, cookieName string) *http.Cookie {
	for _, responseCookie := range responseRecorder.Result().Cookies() {
		if responseCookie.Name == cookieName {
			return responseCookie
		}
	}
	return nil
}

func refreshRequest(csrfHeader string, csrfCookieValue string, withRefreshCookie bool) *http.Request {
	httpRequest := httptest.NewRequest(http.MethodPost, "/api/v1/auth/refresh", nil)
	if withRefreshCookie {
		httpRequest.AddCookie(&http.Cookie{Name: "refresh_token", Value: "refresh-token-value"})
	}
	if csrfCookieValue != "" {
		httpRequest.AddCookie(&http.Cookie{Name: "csrf_token", Value: csrfCookieValue})
	}
	if csrfHeader != "" {
		httpRequest.Header.Set("X-CSRF-Token", csrfHeader)
	}
	return httpRequest
}

func TestHTTPHandler_Login_SetsAccessRefreshAndCSRFCookies(testingInstance *testing.T) {
	_, serveMux := newHTTPHandlerFixture()
	httpRequest := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", strings.NewReader(`{"email":"nurse@example.com","password":"x"}`))
	responseRecorder := httptest.NewRecorder()

	serveMux.ServeHTTP(responseRecorder, httpRequest)

	require.Equal(testingInstance, http.StatusOK, responseRecorder.Code)
	accessCookie := findResponseCookie(responseRecorder, "token")
	refreshCookie := findResponseCookie(responseRecorder, "refresh_token")
	csrfCookie := findResponseCookie(responseRecorder, "csrf_token")
	require.NotNil(testingInstance, accessCookie)
	require.NotNil(testingInstance, refreshCookie)
	require.NotNil(testingInstance, csrfCookie)

	assert.Equal(testingInstance, 900, accessCookie.MaxAge)
	assert.True(testingInstance, accessCookie.HttpOnly)
	assert.Equal(testingInstance, "/api/v1/auth", refreshCookie.Path)
	assert.True(testingInstance, refreshCookie.HttpOnly)
	assert.True(testingInstance, refreshCookie.Secure)
	assert.Equal(testingInstance, http.SameSiteLaxMode, refreshCookie.SameSite)
	assert.InDelta(testingInstance, 12*3600, refreshCookie.MaxAge, 5)
	assert.False(testingInstance, csrfCookie.HttpOnly)
	assert.Contains(testingInstance, responseRecorder.Body.String(), csrfCookie.Value)
}

func TestHTTPHandler_Refresh_RotatesCookiesAndReturnsSession(testingInstance *testing.T) {
	_, serveMux := newHTTPHandlerFixture()
	responseRecorder := httptest.NewRecorder()

	serveMux.ServeHTTP(responseRecorder, refreshRequest("csrf-1", "csrf-1", true))

	require.Equal(testingInstance, http.StatusOK, responseRecorder.Code)
	assert.NotNil(testingInstance, findResponseCookie(responseRecorder, "token"))
	assert.NotNil(testingInstance, findResponseCookie(responseRecorder, "refresh_token"))
	newCSRFCookie := findResponseCookie(responseRecorder, "csrf_token")
	require.NotNil(testingInstance, newCSRFCookie)
	assert.Contains(testingInstance, responseRecorder.Body.String(), `"csrfToken":"`+newCSRFCookie.Value+`"`)
	assert.Contains(testingInstance, responseRecorder.Body.String(), `"role":"NURSE"`)
}

func TestHTTPHandler_Refresh_RequiresCSRF(testingInstance *testing.T) {
	_, serveMux := newHTTPHandlerFixture()

	missingHeaderRecorder := httptest.NewRecorder()
	serveMux.ServeHTTP(missingHeaderRecorder, refreshRequest("", "csrf-1", true))
	assert.Equal(testingInstance, http.StatusForbidden, missingHeaderRecorder.Code)

	mismatchRecorder := httptest.NewRecorder()
	serveMux.ServeHTTP(mismatchRecorder, refreshRequest("other", "csrf-1", true))
	assert.Equal(testingInstance, http.StatusForbidden, mismatchRecorder.Code)
}

func TestHTTPHandler_Refresh_WithoutRefreshCookieIsUnauthorized(testingInstance *testing.T) {
	_, serveMux := newHTTPHandlerFixture()
	responseRecorder := httptest.NewRecorder()

	serveMux.ServeHTTP(responseRecorder, refreshRequest("csrf-1", "csrf-1", false))

	assert.Equal(testingInstance, http.StatusUnauthorized, responseRecorder.Code)
}

func TestHTTPHandler_Refresh_RejectedTokenIsUnauthorizedAndClearsCookies(testingInstance *testing.T) {
	mockService, serveMux := newHTTPHandlerFixture()
	mockService.Err = auth.ErrRefreshTokenReused
	responseRecorder := httptest.NewRecorder()

	serveMux.ServeHTTP(responseRecorder, refreshRequest("csrf-1", "csrf-1", true))

	assert.Equal(testingInstance, http.StatusUnauthorized, responseRecorder.Code)
	refreshCookie := findResponseCookie(responseRecorder, "refresh_token")
	require.NotNil(testingInstance, refreshCookie)
	assert.Equal(testingInstance, -1, refreshCookie.MaxAge)
}

func TestHTTPHandler_Logout_RevokesFamilyAndClearsAllCookies(testingInstance *testing.T) {
	mockService, serveMux := newHTTPHandlerFixture()
	httpRequest := httptest.NewRequest(http.MethodPost, "/api/v1/auth/logout", nil)
	httpRequest.AddCookie(&http.Cookie{Name: "refresh_token", Value: "refresh-token-value"})
	httpRequest.AddCookie(&http.Cookie{Name: "csrf_token", Value: "csrf-1"})
	httpRequest.Header.Set("X-CSRF-Token", "csrf-1")
	responseRecorder := httptest.NewRecorder()

	serveMux.ServeHTTP(responseRecorder, httpRequest)

	require.Equal(testingInstance, http.StatusOK, responseRecorder.Code)
	assert.Equal(testingInstance, []string{"refresh-token-value"}, mockService.LoggedOutTokens)
	for _, cookieName := range []string{"token", "refresh_token", "csrf_token"} {
		clearedCookie := findResponseCookie(responseRecorder, cookieName)
		require.NotNil(testingInstance, clearedCookie, cookieName)
		assert.Equal(testingInstance, -1, clearedCookie.MaxAge, cookieName)
	}
}

func TestHTTPHandler_Logout_WithoutCookiesIsUnauthorized(testingInstance *testing.T) {
	_, serveMux := newHTTPHandlerFixture()
	responseRecorder := httptest.NewRecorder()

	serveMux.ServeHTTP(responseRecorder, httptest.NewRequest(http.MethodPost, "/api/v1/auth/logout", nil))

	assert.Equal(testingInstance, http.StatusUnauthorized, responseRecorder.Code)
}
