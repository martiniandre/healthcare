package mocks

import (
	"context"
	"time"

	"github.com/healthcare/backend/internal/modules/auth"
)

type MockService struct {
	User            *auth.User
	Token           string
	RefreshToken    string
	LoggedOutTokens []string
	Err             error
}

func NewMockService() *MockService {
	return &MockService{}
}

func (mockService *MockService) Register(contextParam context.Context, email, password, fullName, role string) (*auth.User, error) {
	if mockService.Err != nil {
		return nil, mockService.Err
	}
	return mockService.User, nil
}

func (mockService *MockService) Login(contextParam context.Context, email, password string, clientInfo auth.ClientInfo) (*auth.User, *auth.SessionTokens, error) {
	if mockService.Err != nil {
		return nil, nil, mockService.Err
	}
	return mockService.User, mockService.sessionTokens(), nil
}

func (mockService *MockService) Refresh(contextParam context.Context, presentedRefreshToken string, clientInfo auth.ClientInfo) (*auth.User, *auth.SessionTokens, error) {
	if mockService.Err != nil {
		return nil, nil, mockService.Err
	}
	return mockService.User, mockService.sessionTokens(), nil
}

func (mockService *MockService) Logout(contextParam context.Context, presentedRefreshToken string) error {
	mockService.LoggedOutTokens = append(mockService.LoggedOutTokens, presentedRefreshToken)
	return mockService.Err
}

func (mockService *MockService) sessionTokens() *auth.SessionTokens {
	return &auth.SessionTokens{
		AccessToken:      mockService.Token,
		RefreshToken:     mockService.RefreshToken,
		SessionExpiresAt: time.Now().Add(12 * time.Hour),
	}
}

func (mockService *MockService) Me(contextParam context.Context, userID string) (*auth.User, error) {
	if mockService.Err != nil {
		return nil, mockService.Err
	}
	return mockService.User, nil
}
