package tests

import (
	"context"
	"testing"

	"github.com/healthcare/backend/internal/modules/auth"
	"github.com/healthcare/backend/internal/modules/auth/mocks"
	"github.com/healthcare/backend/internal/shared/eventbus"
	"github.com/healthcare/backend/internal/shared/role"
	"github.com/stretchr/testify/assert"
)

type mockAuthEventBus struct {
	PublishedEvents []eventbus.Event
}

func (mockBus *mockAuthEventBus) Publish(ctx context.Context, event eventbus.Event) error {
	mockBus.PublishedEvents = append(mockBus.PublishedEvents, event)
	return nil
}

func (mockBus *mockAuthEventBus) Subscribe(eventName string, handler eventbus.Handler) {}

func TestService_Register(testingInstance *testing.T) {
	mockRepository := mocks.NewMockRepository()
	authService := auth.NewService(mockRepository, nil)
	contextParam := context.Background()

	user, err := authService.Register(contextParam, "test@example.com", "password123", "Test User", string(role.RoleAdmin))

	assert.NoError(testingInstance, err)
	assert.NotNil(testingInstance, user)
	assert.Equal(testingInstance, "test@example.com", user.Email)
	assert.NotEmpty(testingInstance, user.PasswordHash)

	_, errDuplicate := authService.Register(contextParam, "test@example.com", "password123", "Test User 2", string(role.RoleAdmin))
	assert.ErrorIs(testingInstance, errDuplicate, auth.ErrUserExists)
}

func TestService_Login_PublishesEvent(testingInstance *testing.T) {
	mockRepository := mocks.NewMockRepository()
	eventBus := &mockAuthEventBus{}
	authService := auth.NewService(mockRepository, eventBus)
	contextParam := context.Background()

	_, errRegister := authService.Register(contextParam, "event@example.com", "securepass", "Event User", string(role.RoleDoctor))
	assert.NoError(testingInstance, errRegister)

	user, token, errLogin := authService.Login(contextParam, "event@example.com", "securepass")

	assert.NoError(testingInstance, errLogin)
	assert.NotNil(testingInstance, user)
	assert.NotEmpty(testingInstance, token)

	assert.Len(testingInstance, eventBus.PublishedEvents, 1)
	assert.Equal(testingInstance, "system.notification", eventBus.PublishedEvents[0].Name)
	assert.Equal(testingInstance, "Login Realizado", eventBus.PublishedEvents[0].Data["title"])
	assert.Equal(testingInstance, "Login realizado com sucesso", eventBus.PublishedEvents[0].Data["body"])
	assert.Equal(testingInstance, "user", eventBus.PublishedEvents[0].Data["resource_type"])
}

func TestService_Login(testingInstance *testing.T) {
	mockRepository := mocks.NewMockRepository()
	authService := auth.NewService(mockRepository, nil)
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
