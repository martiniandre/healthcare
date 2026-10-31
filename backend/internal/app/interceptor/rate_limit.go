package interceptor

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/healthcare/backend/internal/shared/apperrors"
	"github.com/healthcare/backend/internal/shared/ratelimit"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

var generalRateLimitPerMinute = getEnvRateLimit("RATE_LIMIT_PER_MINUTE", 6000)
var authRateLimitPerMinute = getEnvRateLimit("RATE_LIMIT_AUTH_PER_MINUTE", 1200)

func getEnvRateLimit(envKey string, defaultValue int) int {
	rawValue := os.Getenv(envKey)
	if rawValue == "" {
		return defaultValue
	}
	parsedValue, parseError := strconv.Atoi(rawValue)
	if parseError != nil || parsedValue <= 0 {
		return defaultValue
	}
	return parsedValue
}

func resolveBudget(fullMethod string) ratelimit.Budget {
	maxRequests := generalRateLimitPerMinute
	if strings.HasPrefix(fullMethod, "/auth.v1.AuthService/") {
		maxRequests = authRateLimitPerMinute
	}
	return ratelimit.Budget{MaxRequests: maxRequests, Window: time.Minute}
}

func isRateLimitExceeded(ctx context.Context, limiter ratelimit.Limiter, key string, budget ratelimit.Budget) bool {
	allowed, _, limiterError := limiter.Allow(ctx, key, budget)
	if limiterError != nil {
		slog.Warn("gRPC rate limiter failed, allowing request", "error", limiterError)
		return false
	}
	return !allowed
}

func UnaryRateLimitInterceptor(limiter ratelimit.Limiter, trustedNetworks []*net.IPNet) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		rateLimitKey := fmt.Sprintf("grpc:%s:%s", info.FullMethod, extractClientIP(ctx, trustedNetworks))
		if isRateLimitExceeded(ctx, limiter, rateLimitKey, resolveBudget(info.FullMethod)) {
			return nil, apperrors.ErrRateLimitExceeded.ToGRPC()
		}

		return handler(ctx, req)
	}
}

func extractClientIP(ctx context.Context, trustedNetworks []*net.IPNet) string {
	peerAddress := ""
	if peerInfo, ok := peer.FromContext(ctx); ok && peerInfo.Addr != nil {
		peerAddress = peerInfo.Addr.String()
	}
	if peerAddress == "" {
		return "unknown"
	}

	forwardedFor := ""
	realIP := ""
	if incomingMetadata, ok := metadata.FromIncomingContext(ctx); ok {
		forwardedFor = strings.Join(incomingMetadata.Get("x-forwarded-for"), ",")
		if realIPValues := incomingMetadata.Get("x-real-ip"); len(realIPValues) > 0 {
			realIP = realIPValues[0]
		}
	}

	return ratelimit.ResolveClientIP(peerAddress, forwardedFor, realIP, trustedNetworks)
}

func StreamRateLimitInterceptor(limiter ratelimit.Limiter, trustedNetworks []*net.IPNet) grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		rateLimitKey := fmt.Sprintf("grpc-stream:%s:%s", info.FullMethod, extractClientIP(stream.Context(), trustedNetworks))
		if isRateLimitExceeded(stream.Context(), limiter, rateLimitKey, resolveBudget(info.FullMethod)) {
			return apperrors.ErrRateLimitExceeded.ToGRPC()
		}

		return handler(srv, stream)
	}
}
