package mocks

import (
	"context"
	"time"

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

func (mockRepository *MockRepository) GetUserByID(contextParam context.Context, userID string) (*auth.User, error) {
	if mockRepository.Err != nil {
		return nil, mockRepository.Err
	}
	for _, user := range mockRepository.Users {
		if user.ID.String() == userID {
			return user, nil
		}
	}
	return nil, auth.ErrUserNotFound
}

func (mockRepository *MockRepository) Revoke(contextParam context.Context, tokenDigest string, expiresAt time.Time) error {
	if mockRepository.Err != nil {
		return mockRepository.Err
	}
	return nil
}

func (mockRepository *MockRepository) IsRevoked(contextParam context.Context, tokenDigest string) (bool, error) {
	if mockRepository.Err != nil {
		return false, mockRepository.Err
	}
	return false, nil
}
