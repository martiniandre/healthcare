package notifications

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/healthcare/backend/internal/shared/eventbus"
	"github.com/healthcare/backend/internal/shared/role"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockSubscriberRepository struct {
	createdNotification *Notification
	createdRecipients   []uuid.UUID
	usersByRole         map[string][]uuid.UUID
}

func (mockRepo *mockSubscriberRepository) Create(ctx context.Context, notification *Notification, recipientIDs []uuid.UUID) error {
	mockRepo.createdNotification = notification
	mockRepo.createdRecipients = recipientIDs
	return nil
}

func (mockRepo *mockSubscriberRepository) ListByUserID(ctx context.Context, userID uuid.UUID, limit, offset int32) ([]*Notification, int32, error) {
	return nil, 0, nil
}

func (mockRepo *mockSubscriberRepository) MarkRead(ctx context.Context, notificationID, userID uuid.UUID) error {
	return nil
}

func (mockRepo *mockSubscriberRepository) GetUnreadCount(ctx context.Context, userID uuid.UUID) (int32, error) {
	return 0, nil
}

func (mockRepo *mockSubscriberRepository) GetUserIDsByRole(ctx context.Context, roles []role.Role) ([]uuid.UUID, error) {
	userIDs := make([]uuid.UUID, 0)
	for _, roleValue := range roles {
		userIDs = append(userIDs, mockRepo.usersByRole[string(roleValue)]...)
	}
	return userIDs, nil
}

func (mockRepo *mockSubscriberRepository) GetUserIDsByResource(ctx context.Context, resourceType, resourceID string) ([]uuid.UUID, error) {
	return nil, nil
}

func (mockRepo *mockSubscriberRepository) GetRecipientUserIDs(ctx context.Context, notificationID uuid.UUID) ([]uuid.UUID, error) {
	return mockRepo.createdRecipients, nil
}

func TestSubscribeByRoleHandler_PersistsReportReadyNotificationForPatients(testingInstance *testing.T) {
	patientID := uuid.New()
	mockRepository := &mockSubscriberRepository{
		usersByRole: map[string][]uuid.UUID{"PATIENT": {patientID}},
	}
	notificationService := NewService(mockRepository)
	reportReadyHandler := subscribeByRoleHandler(notificationService, NotificationTypeReportReady)

	handlerError := reportReadyHandler(context.Background(), eventbus.Event{
		Name: "report.ready",
		Data: map[string]any{
			"title":         "Laudo Disponível",
			"body":          "O laudo está pronto para consulta.",
			"resource_type": "diagnostic_report",
			"resource_id":   "report-456",
		},
	})

	require.NoError(testingInstance, handlerError)
	require.NotNil(testingInstance, mockRepository.createdNotification)
	assert.Equal(testingInstance, NotificationTypeReportReady, mockRepository.createdNotification.Type)
	assert.Equal(testingInstance, "Laudo Disponível", mockRepository.createdNotification.Title)
	assert.Equal(testingInstance, "diagnostic_report", mockRepository.createdNotification.ResourceType)
	assert.Equal(testingInstance, "report-456", mockRepository.createdNotification.ResourceID)
	assert.Equal(testingInstance, []uuid.UUID{patientID}, mockRepository.createdRecipients)
}
