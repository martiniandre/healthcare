package app

import (
	"time"

	"github.com/healthcare/backend/internal/app/interceptor"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

type Server struct {
	GRPCServer *grpc.Server
}

func NewServer(redisClient *redis.Client) *Server {
	keepaliveParams := grpc.KeepaliveParams(keepalive.ServerParameters{
		MaxConnectionIdle:     5 * time.Minute,
		MaxConnectionAge:      30 * time.Minute,
		MaxConnectionAgeGrace: 5 * time.Second,
		Time:                  2 * time.Minute,
		Timeout:               20 * time.Second,
	})

	keepalivePolicy := grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
		MinTime:             30 * time.Second,
		PermitWithoutStream: true,
	})

	chainedInterceptor := grpc.ChainUnaryInterceptor(
		interceptor.UnaryLoggingInterceptor(),
		interceptor.UnaryTimeoutInterceptor(),
		interceptor.UnaryRateLimitInterceptor(redisClient),
		interceptor.UnaryCSRFInterceptor(),
		interceptor.UnaryAuthInterceptor(),
		interceptor.UnaryAuditTrailInterceptor(),
	)

	chainedStreamInterceptor := grpc.ChainStreamInterceptor(
		interceptor.StreamLoggingInterceptor(),
		interceptor.StreamRateLimitInterceptor(redisClient),
		interceptor.StreamAuthInterceptor(),
		interceptor.StreamAuditTrailInterceptor(),
	)

	grpcServer := grpc.NewServer(
		keepaliveParams,
		keepalivePolicy,
		grpc.MaxConcurrentStreams(256),
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
		chainedInterceptor,
		chainedStreamInterceptor,
	)

	return &Server{GRPCServer: grpcServer}
}
