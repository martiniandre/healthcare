package auth

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/healthcare/backend/internal/modules/auth/pb"
	"github.com/healthcare/backend/internal/shared/apperrors"
	"github.com/healthcare/backend/internal/shared/validator"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
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

	user, token, loginError := handler.service.Login(ctx, req.Email, req.Password)
	if loginError != nil {
		if errors.Is(loginError, ErrUserNotFound) || errors.Is(loginError, ErrInvalidPassword) || errors.Is(loginError, ErrUserInactive) {
			return nil, apperrors.ErrInvalidCredentials.ToGRPC()
		}
		return nil, apperrors.ErrInternalServer.ToGRPC()
	}

	csrfToken := uuid.New().String()

	grpc.SetHeader(ctx, metadata.Pairs(
		"set-cookie", "token="+token+"; HttpOnly; Secure; Path=/",
		"set-cookie", "csrf_token="+csrfToken+"; Secure; Path=/; SameSite=Lax",
	))

	return &pb.LoginResponse{
		Token:  token,
		UserId: user.ID.String(),
		Role:   string(user.Role),
	}, nil
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

	user, registerError := handler.service.Register(ctx, req.Email, req.Password, req.FullName, req.Role)
	if registerError != nil {
		if errors.Is(registerError, ErrUserExists) {
			return nil, apperrors.ErrUserAlreadyExists.ToGRPC()
		}
		if errors.Is(registerError, ErrInvalidRole) {
			return nil, apperrors.ErrBadRequest.ToGRPC()
		}
		return nil, apperrors.ErrInternalServer.ToGRPC()
	}

	return &pb.RegisterResponse{
		UserId: user.ID.String(),
	}, nil
}

func (handler *GRPCHandler) Logout(ctx context.Context, req *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	grpc.SetHeader(ctx, metadata.Pairs(
		"set-cookie", "token=; HttpOnly; Secure; Path=/; Max-Age=0",
		"set-cookie", "csrf_token=; Secure; Path=/; SameSite=Lax; Max-Age=0",
	))
	return &pb.LogoutResponse{}, nil
}
