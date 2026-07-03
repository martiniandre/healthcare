package interceptor

import (
	"context"
	"strings"

	"github.com/healthcare/backend/internal/modules/auth"
	"github.com/healthcare/backend/internal/shared/apperrors"
	"github.com/healthcare/backend/internal/shared/ctxkeys"
	"github.com/healthcare/backend/internal/shared/role"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func UnaryAuthInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		if publicMethods[info.FullMethod] {
			return handler(ctx, req)
		}

		incomingMetadata, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, apperrors.ErrMissingToken.ToGRPC()
		}

		token := extractJWTFromCookie(incomingMetadata)
		if token == "" {
			return nil, apperrors.ErrMissingToken.ToGRPC()
		}

		claims, err := auth.ValidateJWT(token)
		if err != nil {
			return nil, apperrors.ErrTokenExpired.ToGRPC()
		}

		roleStr, ok := claims["role"].(string)
		if !ok {
			return nil, apperrors.ErrTokenExpired.ToGRPC()
		}

		callerRole := role.Role(roleStr)
		if permErr := checkPermission(info.FullMethod, callerRole); permErr != nil {
			return nil, permErr
		}

		userIDStr, _ := claims["user_id"].(string)
		ctx = context.WithValue(ctx, ctxkeys.UserIDKey, userIDStr)
		ctx = context.WithValue(ctx, ctxkeys.RoleKey, roleStr)

		return handler(ctx, req)
	}
}

func checkPermission(fullMethod string, callerRole role.Role) error {
	allowedRoles, methodIsDefined := methodPermissions[fullMethod]
	if !methodIsDefined {
		return apperrors.ErrMethodNotInPermissionMatrix.ToGRPC()
	}

	for _, allowedRole := range allowedRoles {
		if callerRole == allowedRole {
			return nil
		}
	}

	return apperrors.ErrPermissionDenied.ToGRPC()
}

func extractJWTFromCookie(incomingMetadata metadata.MD) string {
	for _, cookieHeader := range incomingMetadata.Get("cookie") {
		for _, cookiePart := range strings.Split(cookieHeader, ";") {
			cookiePart = strings.TrimSpace(cookiePart)
			if strings.HasPrefix(cookiePart, "token=") {
				return strings.TrimPrefix(cookiePart, "token=")
			}
		}
	}
	return ""
}

type wrappedStream struct {
	grpc.ServerStream
	wrappedContext context.Context
}

func (stream *wrappedStream) Context() context.Context {
	return stream.wrappedContext
}

func NewWrappedStream(stream grpc.ServerStream, newContext context.Context) grpc.ServerStream {
	return &wrappedStream{
		ServerStream:   stream,
		wrappedContext: newContext,
	}
}

func StreamAuthInterceptor() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		if publicMethods[info.FullMethod] {
			return handler(srv, stream)
		}

		incomingMetadata, ok := metadata.FromIncomingContext(stream.Context())
		if !ok {
			return apperrors.ErrMissingToken.ToGRPC()
		}

		token := extractJWTFromCookie(incomingMetadata)
		if token == "" {
			return apperrors.ErrMissingToken.ToGRPC()
		}

		claims, err := auth.ValidateJWT(token)
		if err != nil {
			return apperrors.ErrTokenExpired.ToGRPC()
		}

		roleStr, ok := claims["role"].(string)
		if !ok {
			return apperrors.ErrTokenExpired.ToGRPC()
		}

		callerRole := role.Role(roleStr)
		if permErr := checkPermission(info.FullMethod, callerRole); permErr != nil {
			return permErr
		}

		userIDStr, _ := claims["user_id"].(string)
		newContext := context.WithValue(stream.Context(), ctxkeys.UserIDKey, userIDStr)
		newContext = context.WithValue(newContext, ctxkeys.RoleKey, roleStr)

		return handler(srv, NewWrappedStream(stream, newContext))
	}
}
