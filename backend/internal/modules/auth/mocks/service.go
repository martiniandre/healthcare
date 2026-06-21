package mocks

import (
	"context"

	"github.com/healthcare/backend/internal/modules/auth"
)

type MockService struct {
	User  *auth.User
	Token string
	Err   error
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

func (mockService *MockService) Login(contextParam context.Context, email, password string) (*auth.User, string, error) {
	if mockService.Err != nil {
		return nil, "", mockService.Err
	}
	return mockService.User, mockService.Token, nil
}

func (mockService *MockService) Me(contextParam context.Context, userID string) (*auth.User, error) {
	if mockService.Err != nil {
		return nil, mockService.Err
	}
	return mockService.User, nil
}
