package interceptor

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/healthcare/backend/internal/shared/ctxkeys"
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
		requestID := extractRequestIDFromContext(ctx)
		if requestID == "" {
			requestID = uuid.New().String()
			ctx = context.WithValue(ctx, ctxkeys.RequestIDKey, requestID)
			ctx = context.WithValue(ctx, ctxkeys.CorrelationIDKey, requestID)
		}

		startTime := time.Now()

		slog.Info("grpc request started",
			"request_id", requestID,
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
			"request_id", requestID,
			"method", info.FullMethod,
			"duration_ms", durationMs,
			"grpc_code", grpcCode.String(),
			"error", fmt.Sprintf("%v", err),
		)

		grpc.SetHeader(ctx, metadata.Pairs("x-request-id", requestID))

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
		requestID := extractRequestIDFromContext(stream.Context())
		var newContext context.Context
		if requestID == "" {
			requestID = uuid.New().String()
			newContext = context.WithValue(stream.Context(), ctxkeys.RequestIDKey, requestID)
			newContext = context.WithValue(newContext, ctxkeys.CorrelationIDKey, requestID)
		} else {
			newContext = stream.Context()
		}
		wrappedStream := NewWrappedStream(stream, newContext)

		startTime := time.Now()

		slog.Info("grpc stream started",
			"request_id", requestID,
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
			"request_id", requestID,
			"method", info.FullMethod,
			"duration_ms", durationMs,
			"grpc_code", grpcCode.String(),
			"error", fmt.Sprintf("%v", err),
		)

		return err
	}
}

func extractRequestIDFromContext(ctx context.Context) string {
	if requestID, ok := ctx.Value(ctxkeys.RequestIDKey).(string); ok && requestID != "" {
		return requestID
	}
	if correlationID, ok := ctx.Value(ctxkeys.CorrelationIDKey).(string); ok && correlationID != "" {
		return correlationID
	}
	return ""
}
