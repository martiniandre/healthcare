package interceptor

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func UnaryLoggingInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		correlationID := uuid.New().String()
		ctx = context.WithValue(ctx, "correlation_id", correlationID)

		startTime := time.Now()

		slog.Info("grpc request started",
			"correlation_id", correlationID,
			"method", info.FullMethod,
		)

		response, err := handler(ctx, req)

		durationMs := time.Since(startTime).Milliseconds()
		grpcCode := codes.OK
		if err != nil {
			grpcCode = status.Code(err)
		}

		logLevel := slog.LevelInfo
		if err != nil {
			logLevel = slog.LevelError
		}

		slog.Log(ctx, logLevel, "grpc request completed",
			"correlation_id", correlationID,
			"method", info.FullMethod,
			"duration_ms", durationMs,
			"grpc_code", grpcCode.String(),
			"error", fmt.Sprintf("%v", err),
		)

		grpc.SetHeader(ctx, metadata.Pairs("x-correlation-id", correlationID))

		return response, err
	}
}

func StreamLoggingInterceptor() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		correlationID := uuid.New().String()
		newContext := context.WithValue(stream.Context(), "correlation_id", correlationID)
		wrappedStream := NewWrappedStream(stream, newContext)

		startTime := time.Now()

		slog.Info("grpc stream started",
			"correlation_id", correlationID,
			"method", info.FullMethod,
		)

		err := handler(srv, wrappedStream)

		durationMs := time.Since(startTime).Milliseconds()
		grpcCode := codes.OK
		if err != nil {
			grpcCode = status.Code(err)
		}

		logLevel := slog.LevelInfo
		if err != nil {
			logLevel = slog.LevelError
		}

		slog.Log(newContext, logLevel, "grpc stream completed",
			"correlation_id", correlationID,
			"method", info.FullMethod,
			"duration_ms", durationMs,
			"grpc_code", grpcCode.String(),
			"error", fmt.Sprintf("%v", err),
		)

		return err
	}
}
