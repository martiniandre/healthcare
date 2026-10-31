package middleware

import (
	"context"
	"errors"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/healthcare/backend/internal/shared/ratelimit"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsAuthRateLimitedRouteMatchesLoginOnly(testingInstance *testing.T) {
	assert.True(testingInstance, isAuthRateLimitedRoute("/api/v1/auth/login"))
	assert.False(testingInstance, isAuthRateLimitedRoute("/api/v1/auth/logout"))
	assert.False(testingInstance, isAuthRateLimitedRoute("/api/v1/auth/me"))
	assert.False(testingInstance, isAuthRateLimitedRoute("/api/v1/patients"))
}

func mustParseTrustedProxies(testingInstance *testing.T, entries ...string) []*net.IPNet {
	trustedNetworks, parseError := ParseTrustedProxies(entries)
	require.NoError(testingInstance, parseError)
	return trustedNetworks
}

func TestExtractClientIPAddressIgnoresForwardedForWhenNoProxyIsTrusted(testingInstance *testing.T) {
	httpRequest := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", nil)
	httpRequest.RemoteAddr = "10.0.0.5:54321"
	httpRequest.Header.Set("X-Forwarded-For", "203.0.113.7, 70.41.3.1")
	httpRequest.Header.Set("X-Real-IP", "203.0.113.9")

	assert.Equal(testingInstance, "10.0.0.5", extractClientIPAddress(httpRequest, nil))
}

func TestExtractClientIPAddressIgnoresForwardedForFromUntrustedPeer(testingInstance *testing.T) {
	httpRequest := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", nil)
	httpRequest.RemoteAddr = "198.51.100.1:54321"
	httpRequest.Header.Set("X-Forwarded-For", "203.0.113.7")

	assert.Equal(testingInstance, "198.51.100.1", extractClientIPAddress(httpRequest, mustParseTrustedProxies(testingInstance, "10.0.0.0/8")))
}

func TestExtractClientIPAddressUsesRightmostUntrustedHop(testingInstance *testing.T) {
	httpRequest := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", nil)
	httpRequest.RemoteAddr = "10.0.0.5:54321"
	httpRequest.Header.Set("X-Forwarded-For", "1.1.1.1, 203.0.113.7, 10.0.0.9")

	assert.Equal(testingInstance, "203.0.113.7", extractClientIPAddress(httpRequest, mustParseTrustedProxies(testingInstance, "10.0.0.0/8")))
}

func TestExtractClientIPAddressFallsBackToRealIPFromTrustedProxy(testingInstance *testing.T) {
	httpRequest := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", nil)
	httpRequest.RemoteAddr = "10.0.0.5:54321"
	httpRequest.Header.Set("X-Real-IP", "203.0.113.9")

	assert.Equal(testingInstance, "203.0.113.9", extractClientIPAddress(httpRequest, mustParseTrustedProxies(testingInstance, "10.0.0.5")))
}

func TestExtractClientIPAddressFallsBackToRemoteAddr(testingInstance *testing.T) {
	httpRequest := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", nil)
	httpRequest.RemoteAddr = "10.0.0.5:54321"

	assert.Equal(testingInstance, "10.0.0.5", extractClientIPAddress(httpRequest, nil))
}

func TestParseTrustedProxiesRejectsInvalidEntries(testingInstance *testing.T) {
	_, parseError := ParseTrustedProxies([]string{"not-an-ip"})
	assert.Error(testingInstance, parseError)
}

func TestRateLimitMiddlewareBlocksEleventhLoginWithinWindow(testingInstance *testing.T) {

	targetHandler := http.HandlerFunc(func(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
		httpResponseWriter.WriteHeader(http.StatusOK)
	})
	rateLimitedHandler := RateLimit(nil, ratelimit.NewMemoryLimiter())(targetHandler)

	var lastResponse *httptest.ResponseRecorder
	for requestIndex := 0; requestIndex < authRateLimitPerMinute; requestIndex++ {
		httpRequest := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", nil)
		httpRequest.RemoteAddr = "192.168.0.10:5000"
		lastResponse = httptest.NewRecorder()
		rateLimitedHandler.ServeHTTP(lastResponse, httpRequest)
	}
	require.Equal(testingInstance, http.StatusOK, lastResponse.Code)

	blockedRequest := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", nil)
	blockedRequest.RemoteAddr = "192.168.0.10:5000"
	blockedResponse := httptest.NewRecorder()
	rateLimitedHandler.ServeHTTP(blockedResponse, blockedRequest)

	assert.Equal(testingInstance, http.StatusTooManyRequests, blockedResponse.Code)
	assert.NotEmpty(testingInstance, blockedResponse.Header().Get("Retry-After"))
}

func TestRateLimitMiddlewareKeepsOtherClientsUnaffected(testingInstance *testing.T) {

	targetHandler := http.HandlerFunc(func(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
		httpResponseWriter.WriteHeader(http.StatusOK)
	})
	rateLimitedHandler := RateLimit(nil, ratelimit.NewMemoryLimiter())(targetHandler)

	for requestIndex := 0; requestIndex < authRateLimitPerMinute; requestIndex++ {
		httpRequest := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", nil)
		httpRequest.RemoteAddr = "192.168.0.20:5000"
		responseRecorder := httptest.NewRecorder()
		rateLimitedHandler.ServeHTTP(responseRecorder, httpRequest)
		require.Equal(testingInstance, http.StatusOK, responseRecorder.Code)
	}

	otherClientRequest := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", nil)
	otherClientRequest.RemoteAddr = "192.168.0.30:5000"
	otherClientResponse := httptest.NewRecorder()
	rateLimitedHandler.ServeHTTP(otherClientResponse, otherClientRequest)

	assert.Equal(testingInstance, http.StatusOK, otherClientResponse.Code)
}

type failingLimiter struct{}

func (failingLimiter) Allow(context.Context, string, ratelimit.Budget) (bool, time.Duration, error) {
	return false, 0, errors.New("limiter down")
}

func okHandler() http.Handler {
	return http.HandlerFunc(func(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
		httpResponseWriter.WriteHeader(http.StatusOK)
	})
}

func TestRateLimitMiddlewareFailsOpenWhenLimiterErrors(testingInstance *testing.T) {
	rateLimitedHandler := RateLimit(nil, failingLimiter{})(okHandler())

	httpRequest := httptest.NewRequest(http.MethodGet, "/api/v1/patients", nil)
	responseRecorder := httptest.NewRecorder()
	rateLimitedHandler.ServeHTTP(responseRecorder, httpRequest)

	assert.Equal(testingInstance, http.StatusOK, responseRecorder.Code)
}

func TestRateLimitMiddlewareSetsRetryAfterFromRemainingWindow(testingInstance *testing.T) {
	currentTime := time.Now()
	limiter := ratelimit.NewMemoryLimiterWithClock(func() time.Time { return currentTime })
	rateLimitedHandler := RateLimit(nil, limiter)(okHandler())

	for requestIndex := 0; requestIndex < authRateLimitPerMinute; requestIndex++ {
		httpRequest := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", nil)
		httpRequest.RemoteAddr = "192.168.0.40:5000"
		rateLimitedHandler.ServeHTTP(httptest.NewRecorder(), httpRequest)
	}

	currentTime = currentTime.Add(45 * time.Second)
	blockedRequest := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", nil)
	blockedRequest.RemoteAddr = "192.168.0.40:5000"
	blockedResponse := httptest.NewRecorder()
	rateLimitedHandler.ServeHTTP(blockedResponse, blockedRequest)

	assert.Equal(testingInstance, http.StatusTooManyRequests, blockedResponse.Code)
	assert.Equal(testingInstance, "15", blockedResponse.Header().Get("Retry-After"))
}

func TestRateLimitMiddlewareKeepsLoginAndGeneralBudgetsSeparate(testingInstance *testing.T) {
	rateLimitedHandler := RateLimit(nil, ratelimit.NewMemoryLimiter())(okHandler())

	for requestIndex := 0; requestIndex <= authRateLimitPerMinute; requestIndex++ {
		httpRequest := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", nil)
		httpRequest.RemoteAddr = "192.168.0.50:5000"
		rateLimitedHandler.ServeHTTP(httptest.NewRecorder(), httpRequest)
	}

	generalRequest := httptest.NewRequest(http.MethodGet, "/api/v1/patients", nil)
	generalRequest.RemoteAddr = "192.168.0.50:5000"
	generalResponse := httptest.NewRecorder()
	rateLimitedHandler.ServeHTTP(generalResponse, generalRequest)

	assert.Equal(testingInstance, http.StatusOK, generalResponse.Code)
}
