package mocks

import (
	"context"

	"github.com/healthcare/backend/internal/modules/auth"
)

type MockRepository struct {
	Users map[string]*auth.User
	Err   error
}

func NewMockRepository() *MockRepository {
	return &MockRepository{
		Users: make(map[string]*auth.User),
	}
}

func (mockRepository *MockRepository) CreateUser(contextParam context.Context, user *auth.User) error {
	if mockRepository.Err != nil {
		return mockRepository.Err
	}
	mockRepository.Users[user.Email] = user
	return nil
}

func (mockRepository *MockRepository) GetUserByEmail(contextParam context.Context, email string) (*auth.User, error) {
	if mockRepository.Err != nil {
		return nil, mockRepository.Err
	}
	user, exists := mockRepository.Users[email]
	if !exists {
		return nil, auth.ErrUserNotFound
	}
	return user, nil
}
