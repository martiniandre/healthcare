package auth

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/healthcare/backend/internal/modules/auth/pb"
	"github.com/healthcare/backend/internal/shared/apperrors"
	"github.com/healthcare/backend/internal/shared/role"
	"github.com/healthcare/backend/internal/shared/validator"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

type GRPCHandler struct {
	service Service
}

func NewGRPCHandler(service Service) *GRPCHandler {
	return &GRPCHandler{service: service}
}

func (handler *GRPCHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	violations := make(map[string]string)
	if strings.TrimSpace(req.Email) == "" || !validator.IsValidEmail(req.Email) {
		violations["email"] = "invalid email format"
	}
	if strings.TrimSpace(req.Password) == "" {
		violations["password"] = "password is required"
	}
	if len(violations) > 0 {
		return nil, apperrors.ErrBadRequest.WithFields(violations)
	}

	user, sessionTokens, loginError := handler.service.Login(ctx, req.Email, req.Password, extractGRPCClientInfo(ctx))
	if loginError != nil {
		if errors.Is(loginError, ErrUserNotFound) || errors.Is(loginError, ErrInvalidPassword) || errors.Is(loginError, ErrUserInactive) {
			return nil, apperrors.ErrInvalidCredentials.ToGRPC()
		}
		return nil, apperrors.ErrInternalServer.ToGRPC()
	}

	csrfToken := uuid.New().String()

	remainingSessionSeconds := int(time.Until(sessionTokens.SessionExpiresAt).Seconds())
	if remainingSessionSeconds < 1 {
		remainingSessionSeconds = 1
	}
	accessTokenSeconds := int(accessTokenLifetime.Seconds())
	if remainingSessionSeconds < accessTokenSeconds {
		accessTokenSeconds = remainingSessionSeconds
	}

	grpc.SetHeader(ctx, metadata.Pairs(
		"set-cookie", fmt.Sprintf("token=%s; HttpOnly; Secure; Path=%s; Max-Age=%d", sessionTokens.AccessToken, accessTokenCookiePath, accessTokenSeconds),
		"set-cookie", fmt.Sprintf("refresh_token=%s; HttpOnly; Secure; Path=%s; SameSite=Lax; Max-Age=%d", sessionTokens.RefreshToken, refreshTokenCookiePath, remainingSessionSeconds),
		"set-cookie", fmt.Sprintf("csrf_token=%s; Secure; Path=%s; SameSite=Lax; Max-Age=%d", csrfToken, csrfTokenCookiePath, remainingSessionSeconds),
	))

	return &pb.LoginResponse{
		Token:  sessionTokens.AccessToken,
		UserId: user.ID.String(),
		Role:   string(user.Role),
	}, nil
}

func extractGRPCClientInfo(ctx context.Context) ClientInfo {
	clientInfo := ClientInfo{}
	if incomingMetadata, metadataAvailable := metadata.FromIncomingContext(ctx); metadataAvailable {
		if userAgents := incomingMetadata.Get("user-agent"); len(userAgents) > 0 {
			clientInfo.UserAgent = userAgents[0]
		}
	}
	if callerPeer, peerAvailable := peer.FromContext(ctx); peerAvailable && callerPeer.Addr != nil {
		remoteHost, _, splitError := net.SplitHostPort(callerPeer.Addr.String())
		if splitError != nil {
			remoteHost = callerPeer.Addr.String()
		}
		clientInfo.IP = remoteHost
	}
	return clientInfo
}

func (handler *GRPCHandler) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	violations := make(map[string]string)
	if strings.TrimSpace(req.Email) == "" || !validator.IsValidEmail(req.Email) {
		violations["email"] = "invalid email format"
	}
	if len(req.Password) < 8 {
		violations["password"] = "password must be at least 8 characters long"
	}
	if strings.TrimSpace(req.FullName) == "" {
		violations["full_name"] = "full name is required"
	}
	if len(violations) > 0 {
		return nil, apperrors.ErrBadRequest.WithFields(violations)
	}

	user, registerError := handler.service.Register(ctx, req.Email, req.Password, req.FullName, string(role.RolePatient))
	if registerError != nil {
		if errors.Is(registerError, ErrUserExists) {
			return nil, apperrors.ErrUserAlreadyExists.ToGRPC()
		}
		if errors.Is(registerError, role.ErrInvalidRole) {
			return nil, apperrors.ErrBadRequest.ToGRPC()
		}
		return nil, apperrors.ErrInternalServer.ToGRPC()
	}

	return &pb.RegisterResponse{
		UserId: user.ID.String(),
	}, nil
}

func (handler *GRPCHandler) Logout(ctx context.Context, req *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	incomingMetadata, metadataAvailable := metadata.FromIncomingContext(ctx)
	if metadataAvailable {
		for _, cookieHeader := range incomingMetadata.Get("cookie") {
			for _, cookiePart := range strings.Split(cookieHeader, ";") {
				cookiePart = strings.TrimSpace(cookiePart)
				if strings.HasPrefix(cookiePart, "token=") {
					RevokeToken(strings.TrimPrefix(cookiePart, "token="))
				}
				if strings.HasPrefix(cookiePart, "refresh_token=") {
					handler.service.Logout(ctx, strings.TrimPrefix(cookiePart, "refresh_token="))
				}
			}
		}
	}

	grpc.SetHeader(ctx, metadata.Pairs(
		"set-cookie", "token=; HttpOnly; Secure; Path=/; Max-Age=0",
		"set-cookie", "refresh_token=; HttpOnly; Secure; Path=/api/v1/auth; SameSite=Lax; Max-Age=0",
		"set-cookie", "csrf_token=; Secure; Path=/; SameSite=Lax; Max-Age=0",
	))
	return &pb.LogoutResponse{}, nil
}
