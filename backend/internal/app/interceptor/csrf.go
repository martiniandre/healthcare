package interceptor

import (
	"context"
	"strings"

	"github.com/healthcare/backend/internal/shared/apperrors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var csrfExemptMethods = map[string]bool{
	"/auth.v1.AuthService/Login":    true,
	"/auth.v1.AuthService/Register": true,
	"/auth.v1.AuthService/Logout":   true,
	"/grpc.health.v1.Health/Check":  true,
	"/grpc.health.v1.Health/Watch":  true,
}

func UnaryCSRFInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		if csrfExemptMethods[info.FullMethod] {
			return handler(ctx, req)
		}

		if isReadOnlyMethod(info.FullMethod) {
			return handler(ctx, req)
		}

		incomingMetadata, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, apperrors.ErrCSRFTokenMismatch.ToGRPC()
		}

		csrfTokenFromHeader := extractFirstMetadataValue(incomingMetadata, "x-csrf-token")
		csrfTokenFromCookie := extractCSRFCookieToken(incomingMetadata)

		if csrfTokenFromHeader == "" || csrfTokenFromCookie == "" {
			return nil, apperrors.ErrCSRFTokenMismatch.ToGRPC()
		}

		if csrfTokenFromHeader != csrfTokenFromCookie {
			return nil, apperrors.ErrCSRFTokenMismatch.ToGRPC()
		}

		return handler(ctx, req)
	}
}

func isReadOnlyMethod(fullMethod string) bool {
	methodName := fullMethod
	if lastSlash := strings.LastIndex(fullMethod, "/"); lastSlash >= 0 {
		methodName = fullMethod[lastSlash+1:]
	}

	readOnlyPrefixes := []string{"Get", "List", "Search", "Health", "Check"}
	for _, prefix := range readOnlyPrefixes {
		if strings.HasPrefix(methodName, prefix) {
			return true
		}
	}
	return false
}

func extractFirstMetadataValue(incomingMetadata metadata.MD, key string) string {
	values := incomingMetadata.Get(key)
	if len(values) == 0 {
		return ""
	}
	return strings.TrimSpace(values[0])
}

func extractCSRFCookieToken(incomingMetadata metadata.MD) string {
	cookieValues := incomingMetadata.Get("cookie")
	for _, cookieHeader := range cookieValues {
		for _, cookiePart := range strings.Split(cookieHeader, ";") {
			cookiePart = strings.TrimSpace(cookiePart)
			if strings.HasPrefix(cookiePart, "csrf_token=") {
				return strings.TrimPrefix(cookiePart, "csrf_token=")
			}
		}
	}
	return ""
}
