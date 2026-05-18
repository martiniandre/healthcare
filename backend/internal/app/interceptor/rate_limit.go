package interceptor

import (
	"context"
	"fmt"
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

		currentCount, err := redisClient.Incr(ctx, rateLimitKey).Result()
		if err != nil {
			return handler(ctx, req)
		}

		if currentCount == 1 {
			redisClient.Expire(ctx, rateLimitKey, time.Minute)
		}

		if currentCount > rateLimitRequestsPerMinute {
			return nil, apperrors.ErrRateLimitExceeded.ToGRPC()
		}

		return handler(ctx, req)
	}
}

func extractClientIP(ctx context.Context) string {
	if incomingMetadata, ok := metadata.FromIncomingContext(ctx); ok {
		forwardedFor := incomingMetadata.Get("x-forwarded-for")
		if len(forwardedFor) > 0 {
			return forwardedFor[0]
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

		currentCount, err := redisClient.Incr(stream.Context(), rateLimitKey).Result()
		if err != nil {
			return handler(srv, stream)
		}

		if currentCount == 1 {
			redisClient.Expire(stream.Context(), rateLimitKey, time.Minute)
		}

		if currentCount > rateLimitRequestsPerMinute {
			return apperrors.ErrRateLimitExceeded.ToGRPC()
		}

		return handler(srv, stream)
	}
}
