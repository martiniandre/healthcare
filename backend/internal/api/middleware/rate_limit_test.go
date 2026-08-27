package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSlidingWindowCounterAllowsUpToLimitThenBlocks(testingInstance *testing.T) {
	counter := newSlidingWindowCounter(3)
	baseTime := time.Now()

	assert.True(testingInstance, counter.allow(baseTime))
	assert.True(testingInstance, counter.allow(baseTime.Add(time.Second)))
	assert.True(testingInstance, counter.allow(baseTime.Add(2*time.Second)))
	assert.False(testingInstance, counter.allow(baseTime.Add(3*time.Second)))
}

func TestSlidingWindowCounterReleasesExpiredTimestamps(testingInstance *testing.T) {
	counter := newSlidingWindowCounter(2)
	baseTime := time.Now()

	assert.True(testingInstance, counter.allow(baseTime))
	assert.True(testingInstance, counter.allow(baseTime.Add(time.Second)))
	assert.False(testingInstance, counter.allow(baseTime.Add(2*time.Second)))
	assert.True(testingInstance, counter.allow(baseTime.Add(rateLimitWindow+time.Second)))
}

func TestIsAuthRateLimitedRouteMatchesLoginOnly(testingInstance *testing.T) {
	assert.True(testingInstance, isAuthRateLimitedRoute("/api/v1/auth/login"))
	assert.False(testingInstance, isAuthRateLimitedRoute("/api/v1/auth/logout"))
	assert.False(testingInstance, isAuthRateLimitedRoute("/api/v1/auth/me"))
	assert.False(testingInstance, isAuthRateLimitedRoute("/api/v1/patients"))
}

func TestExtractClientIPAddressPrefersForwardedFor(testingInstance *testing.T) {
	httpRequest := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", nil)
	httpRequest.RemoteAddr = "10.0.0.5:54321"
	httpRequest.Header.Set("X-Forwarded-For", "203.0.113.7, 70.41.3.1")

	assert.Equal(testingInstance, "203.0.113.7", extractClientIPAddress(httpRequest))
}

func TestExtractClientIPAddressFallsBackToRemoteAddr(testingInstance *testing.T) {
	httpRequest := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", nil)
	httpRequest.RemoteAddr = "10.0.0.5:54321"

	assert.Equal(testingInstance, "10.0.0.5", extractClientIPAddress(httpRequest))
}

func TestRateLimitMiddlewareBlocksEleventhLoginWithinWindow(testingInstance *testing.T) {
	ResetRateLimits()
	defer ResetRateLimits()

	targetHandler := http.HandlerFunc(func(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
		httpResponseWriter.WriteHeader(http.StatusOK)
	})
	rateLimitedHandler := RateLimit(targetHandler)

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
	ResetRateLimits()
	defer ResetRateLimits()

	targetHandler := http.HandlerFunc(func(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
		httpResponseWriter.WriteHeader(http.StatusOK)
	})
	rateLimitedHandler := RateLimit(targetHandler)

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
