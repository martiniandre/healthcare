package mocks

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/healthcare/backend/internal/modules/auth"
)

type MockRepository struct {
	Users         map[string]*auth.User
	RefreshTokens map[string]*auth.RefreshToken
	Err           error
}

func NewMockRepository() *MockRepository {
	return &MockRepository{
		Users:         make(map[string]*auth.User),
		RefreshTokens: make(map[string]*auth.RefreshToken),
	}
}

func (mockRepository *MockRepository) CreateUser(contextParam context.Context, user *auth.User) error {
	if mockRepository.Err != nil {
		return mockRepository.Err
	}
	if _, exists := mockRepository.Users[user.Email]; exists {
		return auth.ErrUserExists
	}
	mockRepository.Users[user.Email] = user
	return nil
}

func (mockRepository *MockRepository) CreateUserWithPatientLink(contextParam context.Context, user *auth.User) error {
	return mockRepository.CreateUser(contextParam, user)
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

func (mockRepository *MockRepository) CreateRefreshToken(contextParam context.Context, refreshToken *auth.RefreshToken) error {
	if mockRepository.Err != nil {
		return mockRepository.Err
	}
	mockRepository.RefreshTokens[refreshToken.TokenHash] = refreshToken
	return nil
}

func (mockRepository *MockRepository) GetRefreshTokenByHash(contextParam context.Context, tokenHash string) (*auth.RefreshToken, error) {
	if mockRepository.Err != nil {
		return nil, mockRepository.Err
	}
	refreshToken, exists := mockRepository.RefreshTokens[tokenHash]
	if !exists {
		return nil, auth.ErrRefreshTokenInvalid
	}
	return refreshToken, nil
}

func (mockRepository *MockRepository) RotateRefreshToken(contextParam context.Context, currentTokenID uuid.UUID, nextToken *auth.RefreshToken) (bool, error) {
	if mockRepository.Err != nil {
		return false, mockRepository.Err
	}
	for _, storedToken := range mockRepository.RefreshTokens {
		if storedToken.ID != currentTokenID {
			continue
		}
		if storedToken.ReplacedBy != nil || storedToken.RevokedAt != nil {
			return false, nil
		}
		rotatedAt := time.Now()
		nextTokenID := nextToken.ID
		storedToken.ReplacedBy = &nextTokenID
		storedToken.RevokedAt = &rotatedAt
		mockRepository.RefreshTokens[nextToken.TokenHash] = nextToken
		return true, nil
	}
	return false, nil
}

func (mockRepository *MockRepository) RevokeRefreshTokenFamily(contextParam context.Context, familyID uuid.UUID) error {
	if mockRepository.Err != nil {
		return mockRepository.Err
	}
	revokedAt := time.Now()
	for _, storedToken := range mockRepository.RefreshTokens {
		if storedToken.FamilyID == familyID && storedToken.RevokedAt == nil {
			storedToken.RevokedAt = &revokedAt
		}
	}
	return nil
}
