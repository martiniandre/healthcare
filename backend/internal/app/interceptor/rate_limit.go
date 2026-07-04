package interceptor

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/healthcare/backend/internal/shared/apperrors"
	"github.com/redis/go-redis/v9"
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

func UnaryRateLimitInterceptor(redisClient *redis.Client) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		if redisClient == nil {
			return handler(ctx, req)
		}

		clientIP := extractClientIP(ctx)
		rateLimitKey := fmt.Sprintf("rate_limit:%s:%s", info.FullMethod, clientIP)
		effectiveLimit := generalRateLimitPerMinute
		if strings.HasPrefix(info.FullMethod, "/auth.v1.AuthService/") {
			effectiveLimit = authRateLimitPerMinute
		}

		exceeded, err := checkRateLimit(ctx, redisClient, rateLimitKey, effectiveLimit)
		if err != nil {
			return handler(ctx, req)
		}
		if exceeded {
			return nil, apperrors.ErrRateLimitExceeded.ToGRPC()
		}

		return handler(ctx, req)
	}
}

func checkRateLimit(ctx context.Context, redisClient *redis.Client, key string, limit int) (bool, error) {
	pipe := redisClient.Pipeline()
	incrCmd := pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, time.Minute)
	if _, err := pipe.Exec(ctx); err != nil {
		return false, err
	}
	return incrCmd.Val() > int64(limit), nil
}

func extractClientIP(ctx context.Context) string {
	if incomingMetadata, ok := metadata.FromIncomingContext(ctx); ok {
		forwardedFor := incomingMetadata.Get("x-forwarded-for")
		if len(forwardedFor) > 0 {
			clientIP := strings.TrimSpace(strings.Split(forwardedFor[0], ",")[0])
			if clientIP != "" {
				return clientIP
			}
		}
	}

	if peerInfo, ok := peer.FromContext(ctx); ok {
		return peerInfo.Addr.String()
	}

	return "unknown"
}

func StreamRateLimitInterceptor(redisClient *redis.Client) grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		if redisClient == nil {
			return handler(srv, stream)
		}

		clientIP := extractClientIP(stream.Context())
		rateLimitKey := fmt.Sprintf("rate_limit:stream:%s:%s", info.FullMethod, clientIP)
		effectiveLimit := generalRateLimitPerMinute
		if strings.HasPrefix(info.FullMethod, "/auth.v1.AuthService/") {
			effectiveLimit = authRateLimitPerMinute
		}

		exceeded, err := checkRateLimit(stream.Context(), redisClient, rateLimitKey, effectiveLimit)
		if err != nil {
			return handler(srv, stream)
		}
		if exceeded {
			return apperrors.ErrRateLimitExceeded.ToGRPC()
		}

		return handler(srv, stream)
	}
}
