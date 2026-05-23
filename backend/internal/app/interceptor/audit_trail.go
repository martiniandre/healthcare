package interceptor

import (
	"context"
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"github.com/healthcare/backend/internal/shared/ctxkeys"
)

func UnaryAuditTrailInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		response, err := handler(ctx, req)

		if isClinicalMethod(info.FullMethod) {
			callerUserID := extractContextValue(ctx, ctxkeys.UserIDKey)
			callerRole := extractContextValue(ctx, ctxkeys.RoleKey)
			correlationID := extractContextValue(ctx, ctxkeys.CorrelationIDKey)

			go func() {
				slog.Info("audit trail",
					"correlation_id", correlationID,
					"caller_user_id", callerUserID,
					"caller_role", callerRole,
					"method", info.FullMethod,
					"access_granted", err == nil,
				)
			}()
		}

		return response, err
	}
}

func isClinicalMethod(fullMethod string) bool {
	clinicalPrefixes := []string{
		"/patients.",
		"/clinical.",
		"/observations.",
		"/encounters.",
	}
	for _, prefix := range clinicalPrefixes {
		if len(fullMethod) >= len(prefix) && fullMethod[:len(prefix)] == prefix {
			return true
		}
	}
	return false
}

func extractContextValue(ctx context.Context, key ctxkeys.ContextKey) string {
	if incomingMetadata, ok := metadata.FromIncomingContext(ctx); ok {
		values := incomingMetadata.Get(string(key))
		if len(values) > 0 {
			return values[0]
		}
	}
	if value, ok := ctx.Value(key).(string); ok {
		return value
	}
	return ""
}

func StreamAuditTrailInterceptor() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		err := handler(srv, stream)

		if isClinicalMethod(info.FullMethod) {
			ctx := stream.Context()
			callerUserID := extractContextValue(ctx, ctxkeys.UserIDKey)
			callerRole := extractContextValue(ctx, ctxkeys.RoleKey)
			correlationID := extractContextValue(ctx, ctxkeys.CorrelationIDKey)

			go func() {
				slog.Info("audit trail",
					"correlation_id", correlationID,
					"caller_user_id", callerUserID,
					"caller_role", callerRole,
					"method", info.FullMethod,
					"access_granted", err == nil,
				)
			}()
		}

		return err
	}
}
