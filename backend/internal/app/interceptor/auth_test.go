package interceptor

import (
	"context"
	"testing"

	"github.com/healthcare/backend/internal/modules/auth"
	"github.com/healthcare/backend/internal/shared/role"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func generateTokenForInterceptorTest(testingInstance *testing.T, userRole role.Role) string {
	require.NoError(testingInstance, auth.InitJWT("interceptor-test-secret"))
	token, tokenErr := auth.GenerateJWT("user-1", string(userRole), "user@healthcare.com")
	require.NoError(testingInstance, tokenErr)
	return token
}

func grpcContextWithToken(token string) context.Context {
	return metadata.NewIncomingContext(context.Background(), metadata.Pairs("cookie", "token="+token))
}

func TestUnaryAuthInterceptor_PublicMethodBypassesAuthentication(testingInstance *testing.T) {
	authInterceptor := UnaryAuthInterceptor()

	_, interceptorErr := authInterceptor(context.Background(), nil, &grpc.UnaryServerInfo{FullMethod: "/auth.v1.AuthService/Login"}, func(ctx context.Context, req interface{}) (interface{}, error) {
		return "ok", nil
	})

	assert.NoError(testingInstance, interceptorErr)
}

func TestUnaryAuthInterceptor_AllowsAuthorizedRole(testingInstance *testing.T) {
	authInterceptor := UnaryAuthInterceptor()
	authorizedContext := grpcContextWithToken(generateTokenForInterceptorTest(testingInstance, role.RoleAdmin))

	_, interceptorErr := authInterceptor(authorizedContext, nil, &grpc.UnaryServerInfo{FullMethod: "/audit_logs.v1.AuditLogsService/ListAuditLogs"}, func(ctx context.Context, req interface{}) (interface{}, error) {
		return "ok", nil
	})

	assert.NoError(testingInstance, interceptorErr)
}

func TestUnaryAuthInterceptor_DeniesUnauthorizedRole(testingInstance *testing.T) {
	authInterceptor := UnaryAuthInterceptor()
	unauthorizedContext := grpcContextWithToken(generateTokenForInterceptorTest(testingInstance, role.RoleNurse))

	_, interceptorErr := authInterceptor(unauthorizedContext, nil, &grpc.UnaryServerInfo{FullMethod: "/audit_logs.v1.AuditLogsService/ListAuditLogs"}, func(ctx context.Context, req interface{}) (interface{}, error) {
		return "ok", nil
	})

	assert.Equal(testingInstance, codes.PermissionDenied, status.Code(interceptorErr))
}

func TestUnaryAuthInterceptor_DeniesMissingToken(testingInstance *testing.T) {
	authInterceptor := UnaryAuthInterceptor()

	_, interceptorErr := authInterceptor(context.Background(), nil, &grpc.UnaryServerInfo{FullMethod: "/audit_logs.v1.AuditLogsService/ListAuditLogs"}, func(ctx context.Context, req interface{}) (interface{}, error) {
		return "ok", nil
	})

	assert.Equal(testingInstance, codes.Unauthenticated, status.Code(interceptorErr))
}

func TestUnaryAuthInterceptor_DeniesUnregisteredMethod(testingInstance *testing.T) {
	authInterceptor := UnaryAuthInterceptor()
	authorizedContext := grpcContextWithToken(generateTokenForInterceptorTest(testingInstance, role.RoleAdmin))

	_, interceptorErr := authInterceptor(authorizedContext, nil, &grpc.UnaryServerInfo{FullMethod: "/unregistered.v1.Service/UnregisteredMethod"}, func(ctx context.Context, req interface{}) (interface{}, error) {
		return "ok", nil
	})

	assert.Equal(testingInstance, codes.PermissionDenied, status.Code(interceptorErr))
}

func TestStreamAuthInterceptor_PublicMethodBypassesAuthentication(testingInstance *testing.T) {
	authInterceptor := StreamAuthInterceptor()
	mockStream := &grpcServerStreamStub{}

	interceptorErr := authInterceptor(nil, mockStream, &grpc.StreamServerInfo{FullMethod: "/grpc.health.v1.Health/Watch"}, func(srv interface{}, stream grpc.ServerStream) error {
		return nil
	})

	assert.NoError(testingInstance, interceptorErr)
}

type grpcServerStreamStub struct {
	grpc.ServerStream
}

func (stream *grpcServerStreamStub) Context() context.Context {
	return context.Background()
}
