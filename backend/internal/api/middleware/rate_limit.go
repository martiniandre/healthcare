package middleware

import (
	"fmt"
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
	idleBucketSweepInterval   = time.Minute
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

func (counter *slidingWindowCounter) isIdle(now time.Time) bool {
	counter.mu.Lock()
	defer counter.mu.Unlock()

	if len(counter.timestamps) == 0 {
		return true
	}
	return counter.timestamps[len(counter.timestamps)-1].Before(now.Add(-rateLimitWindow))
}

type rateLimiterRegistry struct {
	mu             sync.Mutex
	generalBuckets map[string]*slidingWindowCounter
	authBuckets    map[string]*slidingWindowCounter
	lastSweepAt    time.Time
}

func newRateLimiterRegistry() *rateLimiterRegistry {
	return &rateLimiterRegistry{
		generalBuckets: make(map[string]*slidingWindowCounter),
		authBuckets:    make(map[string]*slidingWindowCounter),
		lastSweepAt:    time.Now(),
	}
}

func evictIdleBuckets(buckets map[string]*slidingWindowCounter, now time.Time) {
	for clientKey, bucket := range buckets {
		if bucket.isIdle(now) {
			delete(buckets, clientKey)
		}
	}
}

func (registry *rateLimiterRegistry) sweepIdleBucketsLocked(now time.Time) {
	if now.Sub(registry.lastSweepAt) < idleBucketSweepInterval {
		return
	}
	registry.lastSweepAt = now
	evictIdleBuckets(registry.generalBuckets, now)
	evictIdleBuckets(registry.authBuckets, now)
}

func (registry *rateLimiterRegistry) bucketFor(isAuthRoute bool, clientKey string) *slidingWindowCounter {
	registry.mu.Lock()
	defer registry.mu.Unlock()

	registry.sweepIdleBucketsLocked(time.Now())

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

func ParseTrustedProxies(trustedProxyEntries []string) ([]*net.IPNet, error) {
	trustedNetworks := make([]*net.IPNet, 0, len(trustedProxyEntries))
	for _, rawEntry := range trustedProxyEntries {
		trimmedEntry := strings.TrimSpace(rawEntry)
		if trimmedEntry == "" {
			continue
		}
		if !strings.Contains(trimmedEntry, "/") {
			parsedIP := net.ParseIP(trimmedEntry)
			if parsedIP == nil {
				return nil, fmt.Errorf("invalid trusted proxy %q", trimmedEntry)
			}
			if parsedIP.To4() != nil {
				trimmedEntry += "/32"
			} else {
				trimmedEntry += "/128"
			}
		}
		_, parsedNetwork, parseError := net.ParseCIDR(trimmedEntry)
		if parseError != nil {
			return nil, fmt.Errorf("invalid trusted proxy %q: %w", trimmedEntry, parseError)
		}
		trustedNetworks = append(trustedNetworks, parsedNetwork)
	}
	return trustedNetworks, nil
}

func isTrustedProxy(trustedNetworks []*net.IPNet, candidateIP net.IP) bool {
	for _, trustedNetwork := range trustedNetworks {
		if trustedNetwork.Contains(candidateIP) {
			return true
		}
	}
	return false
}

func extractClientIPAddress(httpRequest *http.Request, trustedNetworks []*net.IPNet) string {
	remoteHost, _, splitError := net.SplitHostPort(httpRequest.RemoteAddr)
	if splitError != nil {
		remoteHost = httpRequest.RemoteAddr
	}

	remoteIP := net.ParseIP(remoteHost)
	if remoteIP == nil || !isTrustedProxy(trustedNetworks, remoteIP) {
		return remoteHost
	}

	forwardedFor := httpRequest.Header.Get("X-Forwarded-For")
	if forwardedFor != "" {
		forwardedHops := strings.Split(forwardedFor, ",")
		for hopIndex := len(forwardedHops) - 1; hopIndex >= 0; hopIndex-- {
			hop := strings.TrimSpace(forwardedHops[hopIndex])
			parsedHop := net.ParseIP(hop)
			if parsedHop == nil {
				break
			}
			if !isTrustedProxy(trustedNetworks, parsedHop) {
				return hop
			}
		}
	}

	realIP := strings.TrimSpace(httpRequest.Header.Get("X-Real-IP"))
	if realIP != "" && net.ParseIP(realIP) != nil {
		return realIP
	}

	return remoteHost
}

func RateLimit(trustedNetworks []*net.IPNet) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
			clientKey := extractClientIPAddress(httpRequest, trustedNetworks)
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
}
