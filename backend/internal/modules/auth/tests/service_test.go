package tests

import (
	"context"
	"testing"

	"github.com/healthcare/backend/internal/modules/auth"
	"github.com/healthcare/backend/internal/modules/auth/mocks"
	"github.com/healthcare/backend/internal/shared/role"
	"github.com/stretchr/testify/assert"
)

func TestService_Register(testingInstance *testing.T) {
	mockRepository := mocks.NewMockRepository()
	authService := auth.NewService(mockRepository)
	contextParam := context.Background()

	user, err := authService.Register(contextParam, "test@example.com", "password123", "Test User", string(role.RoleAdmin))

	assert.NoError(testingInstance, err)
	assert.NotNil(testingInstance, user)
	assert.Equal(testingInstance, "test@example.com", user.Email)
	assert.NotEmpty(testingInstance, user.PasswordHash)

	_, errDuplicate := authService.Register(contextParam, "test@example.com", "password123", "Test User 2", string(role.RoleAdmin))
	assert.ErrorIs(testingInstance, errDuplicate, auth.ErrUserExists)
}

func TestService_Login(testingInstance *testing.T) {
	mockRepository := mocks.NewMockRepository()
	authService := auth.NewService(mockRepository)
	contextParam := context.Background()

	_, errRegister := authService.Register(contextParam, "login@example.com", "securepass", "Login User", string(role.RoleDoctor))
	assert.NoError(testingInstance, errRegister)

	user, token, errLogin := authService.Login(contextParam, "login@example.com", "securepass")

	assert.NoError(testingInstance, errLogin)
	assert.NotNil(testingInstance, user)
	assert.NotEmpty(testingInstance, token)

	_, _, errWrongPass := authService.Login(contextParam, "login@example.com", "wrongpass")
	assert.ErrorIs(testingInstance, errWrongPass, auth.ErrInvalidPassword)

	_, _, errNotFound := authService.Login(contextParam, "notfound@example.com", "securepass")
	assert.ErrorIs(testingInstance, errNotFound, auth.ErrUserNotFound)
}
