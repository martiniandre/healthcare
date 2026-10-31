package middleware

import (
	"fmt"
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
