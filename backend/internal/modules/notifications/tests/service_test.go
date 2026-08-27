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

func TestCreateNotificationByRole_RoutesReportReadyToPatients(testingInstance *testing.T) {
	mockRepository := mocks.NewMockNotificationRepository()
	notificationService := notifications.NewService(mockRepository)
	contextParam := context.Background()

	patientID := uuid.New()
	mockRepository.UsersByRole["PATIENT"] = []uuid.UUID{patientID}

	createdNotification, createError := notificationService.CreateNotificationByRole(
		contextParam,
		notifications.NotificationTypeReportReady,
		"Laudo Disponível",
		"O laudo está pronto para consulta.",
		nil,
		"diagnostic_report",
		"report-123",
	)

	assert.NoError(testingInstance, createError)
	assert.NotNil(testingInstance, createdNotification)
	assert.Equal(testingInstance, notifications.PriorityHigh, createdNotification.Priority)
	recipientIDs := mockRepository.Recipients[createdNotification.ID]
	assert.Len(testingInstance, recipientIDs, 1)
	assert.Equal(testingInstance, patientID, recipientIDs[0])
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

	userID := uuid.New()
	sub := notificationService.Subscribe(contextParam, userID)
	defer notificationService.Unsubscribe(sub)

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

func TestBroadcastDoesNotDeliverToNonRecipients(testingInstance *testing.T) {
	mockRepository := mocks.NewMockNotificationRepository()
	notificationService := notifications.NewService(mockRepository)
	contextParam := context.Background()

	recipientUserID := uuid.New()
	outsiderUserID := uuid.New()

	recipientSubscription := notificationService.Subscribe(contextParam, recipientUserID)
	defer notificationService.Unsubscribe(recipientSubscription)
	outsiderSubscription := notificationService.Subscribe(contextParam, outsiderUserID)
	defer notificationService.Unsubscribe(outsiderSubscription)

	createdNotification, createError := notificationService.CreateNotification(
		contextParam,
		notifications.NotificationTypeTelemetryAlert,
		"Clinical Alert",
		"Paciente X apresenta condição crítica",
		nil,
		"bed",
		uuid.New().String(),
		[]uuid.UUID{recipientUserID},
	)

	assert.NoError(testingInstance, createError)

	receivedNotification := <-recipientSubscription.Channel()
	assert.Equal(testingInstance, createdNotification.ID, receivedNotification.ID)

	select {
	case leakedNotification := <-outsiderSubscription.Channel():
		assert.Fail(testingInstance, "non-recipient subscriber received a notification", leakedNotification.Title)
	default:
	}
}

func TestBroadcastFailsClosedWhenRecipientLookupFails(testingInstance *testing.T) {
	mockRepository := mocks.NewMockNotificationRepository()
	mockRepository.GetRecipientsError = assert.AnError
	notificationService := notifications.NewService(mockRepository)
	contextParam := context.Background()

	userID := uuid.New()
	sub := notificationService.Subscribe(contextParam, userID)
	defer notificationService.Unsubscribe(sub)

	_, createError := notificationService.CreateNotification(
		contextParam,
		notifications.NotificationTypeSystem,
		"Failing Lookup",
		"Test body",
		nil,
		"",
		"",
		[]uuid.UUID{userID},
	)

	assert.NoError(testingInstance, createError)

	select {
	case leakedNotification := <-sub.Channel():
		assert.Fail(testingInstance, "notification delivered despite recipient lookup failure", leakedNotification.Title)
	default:
	}
}
