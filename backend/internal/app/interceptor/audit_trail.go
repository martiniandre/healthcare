package interceptor

import (
	"context"
	"log/slog"

	"github.com/healthcare/backend/internal/modules/audit_logs"
	"github.com/healthcare/backend/internal/shared/ctxkeys"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var globalAuditLogsService audit_logs.Service

func SetAuditLogsService(service audit_logs.Service) {
	globalAuditLogsService = service
}

func UnaryAuditTrailInterceptor() grpc.UnaryServerInterceptor {
	return func(
		contextVal context.Context,
		requestVal interface{},
		serverInfo *grpc.UnaryServerInfo,
		unaryHandler grpc.UnaryHandler,
	) (interface{}, error) {
		handlerResponse, executionError := unaryHandler(contextVal, requestVal)

		if isClinicalOrCriticalMethod(serverInfo.FullMethod) {
			callerUserID := extractContextValue(contextVal, ctxkeys.UserIDKey)
			callerRole := extractContextValue(contextVal, ctxkeys.RoleKey)
			correlationID := extractContextValue(contextVal, ctxkeys.CorrelationIDKey)

			go func() {
				if globalAuditLogsService != nil {
					_, logError := globalAuditLogsService.CreateAuditLog(
						context.Background(),
						correlationID,
						callerUserID,
						callerRole,
						serverInfo.FullMethod,
						executionError == nil,
					)
					if logError != nil {
						slog.Error("failed to persist unary audit log", "error", logError)
					}
				} else {
					slog.Info("audit trail",
						"correlation_id", correlationID,
						"caller_user_id", callerUserID,
						"caller_role", callerRole,
						"method", serverInfo.FullMethod,
						"access_granted", executionError == nil,
					)
				}
			}()
		}

		return handlerResponse, executionError
	}
}

func isClinicalOrCriticalMethod(fullMethod string) bool {
	prefixes := []string{
		"/patients.",
		"/clinical.",
		"/observations.",
		"/encounters.",
		"/auth.v1.AuthService/",
		"/staff.v1.StaffService/CreateEmployee",
		"/staff.v1.StaffService/DeactivateEmployee",
		"/telemetry.v1.TelemetryService/UnlockRoom",
	}
	for _, prefix := range prefixes {
		if len(fullMethod) >= len(prefix) && fullMethod[:len(prefix)] == prefix {
			return true
		}
	}
	return false
}

func extractContextValue(contextVal context.Context, key ctxkeys.ContextKey) string {
	if incomingMetadata, metadataOk := metadata.FromIncomingContext(contextVal); metadataOk {
		values := incomingMetadata.Get(string(key))
		if len(values) > 0 {
			return values[0]
		}
	}
	if value, assertOk := contextVal.Value(key).(string); assertOk {
		return value
	}
	return ""
}

func StreamAuditTrailInterceptor() grpc.StreamServerInterceptor {
	return func(
		serviceImpl interface{},
		serverStream grpc.ServerStream,
		streamInfo *grpc.StreamServerInfo,
		streamHandler grpc.StreamHandler,
	) error {
		executionError := streamHandler(serviceImpl, serverStream)

		if isClinicalOrCriticalMethod(streamInfo.FullMethod) {
			contextVal := serverStream.Context()
			callerUserID := extractContextValue(contextVal, ctxkeys.UserIDKey)
			callerRole := extractContextValue(contextVal, ctxkeys.RoleKey)
			correlationID := extractContextValue(contextVal, ctxkeys.CorrelationIDKey)

			go func() {
				if globalAuditLogsService != nil {
					_, logError := globalAuditLogsService.CreateAuditLog(
						context.Background(),
						correlationID,
						callerUserID,
						callerRole,
						streamInfo.FullMethod,
						executionError == nil,
					)
					if logError != nil {
						slog.Error("failed to persist stream audit log", "error", logError)
					}
				} else {
					slog.Info("audit trail",
						"correlation_id", correlationID,
						"caller_user_id", callerUserID,
						"caller_role", callerRole,
						"method", streamInfo.FullMethod,
						"access_granted", executionError == nil,
					)
				}
			}()
		}

		return executionError
	}
}
