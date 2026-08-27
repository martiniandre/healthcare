package middleware

import (
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/healthcare/backend/internal/api/render"
	"github.com/healthcare/backend/internal/shared/apperrors"
)

const (
	defaultRateLimitPerMinute = 240
	authRateLimitPerMinute    = 10
	rateLimitWindow           = time.Minute
)

type slidingWindowCounter struct {
	mu          sync.Mutex
	timestamps  []time.Time
	maxRequests int
}

func newSlidingWindowCounter(maxRequests int) *slidingWindowCounter {
	return &slidingWindowCounter{timestamps: make([]time.Time, 0, maxRequests), maxRequests: maxRequests}
}

func (counter *slidingWindowCounter) allow(now time.Time) bool {
	counter.mu.Lock()
	defer counter.mu.Unlock()

	windowStart := now.Add(-rateLimitWindow)
	firstAliveIndex := 0
	for firstAliveIndex < len(counter.timestamps) && counter.timestamps[firstAliveIndex].Before(windowStart) {
		firstAliveIndex++
	}
	counter.timestamps = counter.timestamps[firstAliveIndex:]

	if len(counter.timestamps) >= counter.maxRequests {
		return false
	}

	counter.timestamps = append(counter.timestamps, now)
	return true
}

type rateLimiterRegistry struct {
	mu             sync.Mutex
	generalBuckets map[string]*slidingWindowCounter
	authBuckets    map[string]*slidingWindowCounter
}

func newRateLimiterRegistry() *rateLimiterRegistry {
	return &rateLimiterRegistry{
		generalBuckets: make(map[string]*slidingWindowCounter),
		authBuckets:    make(map[string]*slidingWindowCounter),
	}
}

func (registry *rateLimiterRegistry) bucketFor(isAuthRoute bool, clientKey string) *slidingWindowCounter {
	registry.mu.Lock()
	defer registry.mu.Unlock()

	if isAuthRoute {
		if existingBucket, exists := registry.authBuckets[clientKey]; exists {
			return existingBucket
		}
		createdBucket := newSlidingWindowCounter(authRateLimitPerMinute)
		registry.authBuckets[clientKey] = createdBucket
		return createdBucket
	}

	if existingBucket, exists := registry.generalBuckets[clientKey]; exists {
		return existingBucket
	}
	createdBucket := newSlidingWindowCounter(defaultRateLimitPerMinute)
	registry.generalBuckets[clientKey] = createdBucket
	return createdBucket
}

var sharedRateLimiterRegistry = newRateLimiterRegistry()

func ResetRateLimits() {
	sharedRateLimiterRegistry = newRateLimiterRegistry()
}

func isAuthRateLimitedRoute(requestPath string) bool {
	return strings.HasSuffix(requestPath, "/auth/login")
}

func extractClientIPAddress(httpRequest *http.Request) string {
	forwardedFor := httpRequest.Header.Get("X-Forwarded-For")
	if forwardedFor != "" {
		firstHop := strings.TrimSpace(strings.Split(forwardedFor, ",")[0])
		if parsedIP := net.ParseIP(firstHop); parsedIP != nil {
			return firstHop
		}
	}

	realIP := httpRequest.Header.Get("X-Real-IP")
	if realIP != "" {
		if parsedIP := net.ParseIP(realIP); parsedIP != nil {
			return realIP
		}
	}

	host, _, splitError := net.SplitHostPort(httpRequest.RemoteAddr)
	if splitError != nil {
		return httpRequest.RemoteAddr
	}
	return host
}

func RateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
		clientKey := extractClientIPAddress(httpRequest)
		isAuthRoute := isAuthRateLimitedRoute(httpRequest.URL.Path)

		bucket := sharedRateLimiterRegistry.bucketFor(isAuthRoute, clientKey)
		if !bucket.allow(time.Now()) {
			httpResponseWriter.Header().Set("Retry-After", strconv.Itoa(int(rateLimitWindow.Seconds())))
			render.ErrorFromAppError(httpResponseWriter, apperrors.ErrRateLimitExceeded)
			return
		}

		next.ServeHTTP(httpResponseWriter, httpRequest)
	})
}
