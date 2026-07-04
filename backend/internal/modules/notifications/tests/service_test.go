package tests

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/healthcare/backend/internal/modules/notifications"
	"github.com/healthcare/backend/internal/modules/notifications/mocks"
	"github.com/stretchr/testify/assert"
)

func TestCreateNotification(testingInstance *testing.T) {
	mockRepository := mocks.NewMockNotificationRepository()
	notificationService := notifications.NewService(mockRepository)
	contextParam := context.Background()
	userID := uuid.New()

	createdNotification, createError := notificationService.CreateNotification(
		contextParam,
		notifications.NotificationTypeSystem,
		"Test Title",
		"Test Body",
		nil,
		"",
		"",
		[]uuid.UUID{userID},
	)

	assert.NoError(testingInstance, createError)
	assert.NotNil(testingInstance, createdNotification)
	assert.Equal(testingInstance, notifications.NotificationTypeSystem, createdNotification.Type)
	assert.Equal(testingInstance, "Test Title", createdNotification.Title)
	assert.Equal(testingInstance, "Test Body", createdNotification.Body)
	assert.Equal(testingInstance, notifications.PriorityLow, createdNotification.Priority)
}

func TestCreateNotificationByRole(testingInstance *testing.T) {
	mockRepository := mocks.NewMockNotificationRepository()
	notificationService := notifications.NewService(mockRepository)
	contextParam := context.Background()

	adminID := uuid.New()
	doctorID := uuid.New()
	mockRepository.UsersByRole["ADMIN"] = []uuid.UUID{adminID}
	mockRepository.UsersByRole["DOCTOR"] = []uuid.UUID{doctorID}

	createdNotification, createError := notificationService.CreateNotificationByRole(
		contextParam,
		notifications.NotificationTypeSystem,
		"System Update",
		"System maintenance scheduled",
		nil,
		"",
		"",
	)

	assert.NoError(testingInstance, createError)
	assert.NotNil(testingInstance, createdNotification)
	assert.Equal(testingInstance, "System Update", createdNotification.Title)
}

func TestCreateNotification_InvalidType(testingInstance *testing.T) {
	mockRepository := mocks.NewMockNotificationRepository()
	notificationService := notifications.NewService(mockRepository)
	contextParam := context.Background()

	_, createError := notificationService.CreateNotificationByRole(
		contextParam,
		notifications.NotificationType("invalid_type"),
		"Test",
		"Test",
		nil,
		"",
		"",
	)

	assert.Error(testingInstance, createError)
	assert.ErrorIs(testingInstance, createError, notifications.ErrInvalidNotificationType)
}

func TestMarkRead(testingInstance *testing.T) {
	mockRepository := mocks.NewMockNotificationRepository()
	notificationService := notifications.NewService(mockRepository)
	contextParam := context.Background()
	userID := uuid.New()
	notificationID := uuid.New()

	markReadError := notificationService.MarkRead(contextParam, notificationID, userID)
	assert.NoError(testingInstance, markReadError)
}

func TestGetUnreadCount(testingInstance *testing.T) {
	mockRepository := mocks.NewMockNotificationRepository()
	notificationService := notifications.NewService(mockRepository)
	contextParam := context.Background()
	userID := uuid.New()
	mockRepository.UnreadCount = 5

	count, countError := notificationService.GetUnreadCount(contextParam, userID)
	assert.NoError(testingInstance, countError)
	assert.Equal(testingInstance, int32(5), count)
}

func TestSubscribeAndBroadcast(testingInstance *testing.T) {
	mockRepository := mocks.NewMockNotificationRepository()
	notificationService := notifications.NewService(mockRepository)
	contextParam := context.Background()

	sub := notificationService.Subscribe(contextParam)
	defer notificationService.Unsubscribe(sub)

	userID := uuid.New()
	createdNotification, createError := notificationService.CreateNotification(
		contextParam,
		notifications.NotificationTypeSystem,
		"Broadcast Test",
		"Test body",
		nil,
		"",
		"",
		[]uuid.UUID{userID},
	)

	assert.NoError(testingInstance, createError)
	assert.NotNil(testingInstance, createdNotification)

	receivedNotification := <-sub.Channel()
	assert.Equal(testingInstance, createdNotification.ID, receivedNotification.ID)
}
