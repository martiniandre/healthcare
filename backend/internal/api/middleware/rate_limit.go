package middleware

import (
	"log/slog"
	"math"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/healthcare/backend/internal/api/render"
	"github.com/healthcare/backend/internal/shared/apperrors"
	"github.com/healthcare/backend/internal/shared/ratelimit"
)

const (
	defaultRateLimitPerMinute = 240
	authRateLimitPerMinute    = 10
	rateLimitWindow           = time.Minute
)

var (
	generalRateLimitBudget = ratelimit.Budget{MaxRequests: defaultRateLimitPerMinute, Window: rateLimitWindow}
	authRateLimitBudget    = ratelimit.Budget{MaxRequests: authRateLimitPerMinute, Window: rateLimitWindow}
)

func isAuthRateLimitedRoute(requestPath string) bool {
	return strings.HasSuffix(requestPath, "/auth/login")
}

func ParseTrustedProxies(trustedProxyEntries []string) ([]*net.IPNet, error) {
	return ratelimit.ParseTrustedProxies(trustedProxyEntries)
}

func extractClientIPAddress(httpRequest *http.Request, trustedNetworks []*net.IPNet) string {
	return ratelimit.ResolveClientIP(
		httpRequest.RemoteAddr,
		httpRequest.Header.Get("X-Forwarded-For"),
		httpRequest.Header.Get("X-Real-IP"),
		trustedNetworks,
	)
}

func RateLimit(trustedNetworks []*net.IPNet, limiter ratelimit.Limiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
			clientIPAddress := extractClientIPAddress(httpRequest, trustedNetworks)

			budget := generalRateLimitBudget
			limiterKey := "http:general:" + clientIPAddress
			if isAuthRateLimitedRoute(httpRequest.URL.Path) {
				budget = authRateLimitBudget
				limiterKey = "http:login:" + clientIPAddress
			}

			allowed, retryAfter, limiterError := limiter.Allow(httpRequest.Context(), limiterKey, budget)
			if limiterError != nil {
				slog.Warn("Rate limiter failed, allowing request", "error", limiterError)
			}
			if limiterError == nil && !allowed {
				retryAfterSeconds := int(math.Ceil(retryAfter.Seconds()))
				if retryAfterSeconds < 1 {
					retryAfterSeconds = 1
				}
				httpResponseWriter.Header().Set("Retry-After", strconv.Itoa(retryAfterSeconds))
				render.ErrorFromAppError(httpResponseWriter, apperrors.ErrRateLimitExceeded)
				return
			}

			next.ServeHTTP(httpResponseWriter, httpRequest)
		})
	}
}
