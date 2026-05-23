package interceptor

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/healthcare/backend/internal/shared/apperrors"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

const rateLimitRequestsPerMinute = 60

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

		exceeded, err := checkRateLimit(ctx, redisClient, rateLimitKey)
		if err != nil {
			return handler(ctx, req)
		}
		if exceeded {
			return nil, apperrors.ErrRateLimitExceeded.ToGRPC()
		}

		return handler(ctx, req)
	}
}

func checkRateLimit(ctx context.Context, redisClient *redis.Client, key string) (bool, error) {
	pipe := redisClient.Pipeline()
	incrCmd := pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, time.Minute)
	if _, err := pipe.Exec(ctx); err != nil {
		return false, err
	}
	return incrCmd.Val() > rateLimitRequestsPerMinute, nil
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

		exceeded, err := checkRateLimit(stream.Context(), redisClient, rateLimitKey)
		if err != nil {
			return handler(srv, stream)
		}
		if exceeded {
			return apperrors.ErrRateLimitExceeded.ToGRPC()
		}

		return handler(srv, stream)
	}
}
