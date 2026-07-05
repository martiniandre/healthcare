package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/healthcare/backend/internal/modules/notifications"
	"github.com/healthcare/backend/internal/shared/role"
)

type MockNotificationRepository struct {
	Notifications []*notifications.Notification
	Recipients    map[uuid.UUID][]uuid.UUID
	UsersByRole   map[string][]uuid.UUID
	UnreadCount   int32
	CreateError   error
	MarkReadError error
	GetUsersError error
}

func NewMockNotificationRepository() *MockNotificationRepository {
	return &MockNotificationRepository{
		Notifications: make([]*notifications.Notification, 0),
		Recipients:    make(map[uuid.UUID][]uuid.UUID),
		UsersByRole:   make(map[string][]uuid.UUID),
	}
}

func (mockRepo *MockNotificationRepository) Create(ctx context.Context, notification *notifications.Notification, recipientIDs []uuid.UUID) error {
	if mockRepo.CreateError != nil {
		return mockRepo.CreateError
	}
	mockRepo.Notifications = append(mockRepo.Notifications, notification)
	mockRepo.Recipients[notification.ID] = recipientIDs
	return nil
}

func (mockRepo *MockNotificationRepository) ListByUserID(ctx context.Context, userID uuid.UUID, limit, offset int32) ([]*notifications.Notification, int32, error) {
	filtered := make([]*notifications.Notification, 0)
	for _, n := range mockRepo.Notifications {
		recipients := mockRepo.Recipients[n.ID]
		for _, recipientID := range recipients {
			if recipientID == userID {
				filtered = append(filtered, n)
				break
			}
		}
	}
	return filtered, int32(len(filtered)), nil
}

func (mockRepo *MockNotificationRepository) MarkRead(ctx context.Context, notificationID, userID uuid.UUID) error {
	if mockRepo.MarkReadError != nil {
		return mockRepo.MarkReadError
	}
	return nil
}

func (mockRepo *MockNotificationRepository) GetUnreadCount(ctx context.Context, userID uuid.UUID) (int32, error) {
	return mockRepo.UnreadCount, nil
}

func (mockRepo *MockNotificationRepository) GetUserIDsByRole(ctx context.Context, roles []role.Role) ([]uuid.UUID, error) {
	if mockRepo.GetUsersError != nil {
		return nil, mockRepo.GetUsersError
	}
	result := make([]uuid.UUID, 0)
	for _, roleValue := range roles {
		if users, exists := mockRepo.UsersByRole[string(roleValue)]; exists {
			result = append(result, users...)
		}
	}
	return result, nil
}

func (mockRepo *MockNotificationRepository) GetUserIDsByResource(ctx context.Context, resourceType, resourceID string) ([]uuid.UUID, error) {
	return nil, nil
}
