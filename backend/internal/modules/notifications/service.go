package notifications

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/healthcare/backend/internal/shared/role"
)

var ErrNotificationNotFound = errors.New("notification not found")
var ErrInvalidNotificationType = errors.New("invalid notification type")

var notificationRoleRoutes = map[NotificationType][]role.Role{
	NotificationTypeTelemetryAlert:  {role.RoleNurse, role.RoleDoctor},
	NotificationTypeExamComplete:    {role.RoleDoctor},
	NotificationTypeEncounterCreate: {role.RoleDoctor, role.RoleReception, role.RolePatient},
	NotificationTypeEncounterUpdate: {role.RoleDoctor, role.RoleReception},
	NotificationTypePatientCreate:   {role.RoleReception, role.RoleAdmin},
	NotificationTypePatientUpdate:   {role.RoleReception, role.RoleDoctor},
	NotificationTypeAuditAlert:      {role.RoleAdmin},
	NotificationTypeSystem:          {role.RoleAdmin, role.RoleDoctor, role.RoleNurse},
}

var notificationPriorityDefaults = map[NotificationType]NotificationPriority{
	NotificationTypeTelemetryAlert:  PriorityCritical,
	NotificationTypeExamComplete:    PriorityHigh,
	NotificationTypeEncounterCreate: PriorityMedium,
	NotificationTypeEncounterUpdate: PriorityMedium,
	NotificationTypePatientCreate:   PriorityLow,
	NotificationTypePatientUpdate:   PriorityLow,
	NotificationTypeAuditAlert:      PriorityHigh,
	NotificationTypeSystem:          PriorityLow,
}

type Subscriber interface {
	ID() string
	Channel() chan *Notification
}

type subscriber struct {
	id      string
	channel chan *Notification
}

func (sub *subscriber) ID() string                  { return sub.id }
func (sub *subscriber) Channel() chan *Notification { return sub.channel }

type Service interface {
	CreateNotification(ctx context.Context, notifType NotificationType, title, body string, actorID *uuid.UUID, resourceType, resourceID string, recipientIDs []uuid.UUID) (*Notification, error)
	CreateNotificationByRole(ctx context.Context, notifType NotificationType, title, body string, actorID *uuid.UUID, resourceType, resourceID string) (*Notification, error)
	ListNotifications(ctx context.Context, userID uuid.UUID, limit, offset int32) ([]*Notification, int32, error)
	MarkRead(ctx context.Context, notificationID, userID uuid.UUID) error
	GetUnreadCount(ctx context.Context, userID uuid.UUID) (int32, error)
	Subscribe(ctx context.Context) Subscriber
	Unsubscribe(sub Subscriber)
}

type service struct {
	repo          Repository
	subscribers   map[string]Subscriber
	subscribersMu sync.RWMutex
}

func NewService(repo Repository) Service {
	return &service{
		repo:        repo,
		subscribers: make(map[string]Subscriber),
	}
}

func (notificationService *service) CreateNotification(ctx context.Context, notifType NotificationType, title, body string, actorID *uuid.UUID, resourceType, resourceID string, recipientIDs []uuid.UUID) (*Notification, error) {
	priority, exists := notificationPriorityDefaults[notifType]
	if !exists {
		priority = PriorityMedium
	}

	notification := &Notification{
		ID:           uuid.New(),
		Type:         notifType,
		Priority:     priority,
		Title:        title,
		Body:         body,
		ActorID:      actorID,
		ResourceType: resourceType,
		ResourceID:   resourceID,
		CreatedAt:    time.Now(),
	}

	err := notificationService.repo.Create(ctx, notification, recipientIDs)
	if err != nil {
		return nil, err
	}

	notificationService.broadcast(notification)
	return notification, nil
}

func (notificationService *service) CreateNotificationByRole(ctx context.Context, notifType NotificationType, title, body string, actorID *uuid.UUID, resourceType, resourceID string) (*Notification, error) {
	roles, exists := notificationRoleRoutes[notifType]
	if !exists {
		return nil, ErrInvalidNotificationType
	}

	recipientIDs, err := notificationService.repo.GetUserIDsByRole(ctx, roles)
	if err != nil {
		return nil, err
	}

	if resourceType != "" && resourceID != "" {
		resourceRecipients, resourceErr := notificationService.repo.GetUserIDsByResource(ctx, resourceType, resourceID)
		if resourceErr == nil && len(resourceRecipients) > 0 {
			existingMap := make(map[uuid.UUID]bool, len(recipientIDs))
			for _, id := range recipientIDs {
				existingMap[id] = true
			}
			for _, id := range resourceRecipients {
				if !existingMap[id] {
					recipientIDs = append(recipientIDs, id)
				}
			}
		}
	}

	return notificationService.CreateNotification(ctx, notifType, title, body, actorID, resourceType, resourceID, recipientIDs)
}

func (notificationService *service) ListNotifications(ctx context.Context, userID uuid.UUID, limit, offset int32) ([]*Notification, int32, error) {
	return notificationService.repo.ListByUserID(ctx, userID, limit, offset)
}

func (notificationService *service) MarkRead(ctx context.Context, notificationID, userID uuid.UUID) error {
	return notificationService.repo.MarkRead(ctx, notificationID, userID)
}

func (notificationService *service) GetUnreadCount(ctx context.Context, userID uuid.UUID) (int32, error) {
	return notificationService.repo.GetUnreadCount(ctx, userID)
}

func (notificationService *service) Subscribe(ctx context.Context) Subscriber {
	sub := &subscriber{
		id:      uuid.New().String(),
		channel: make(chan *Notification, 100),
	}

	notificationService.subscribersMu.Lock()
	notificationService.subscribers[sub.id] = sub
	notificationService.subscribersMu.Unlock()

	go func() {
		<-ctx.Done()
		notificationService.Unsubscribe(sub)
	}()

	return sub
}

func (notificationService *service) Unsubscribe(sub Subscriber) {
	notificationService.subscribersMu.Lock()
	defer notificationService.subscribersMu.Unlock()

	if existingSub, exists := notificationService.subscribers[sub.ID()]; exists {
		close(existingSub.Channel())
		delete(notificationService.subscribers, sub.ID())
	}
}

func (notificationService *service) broadcast(notification *Notification) {
	notificationService.subscribersMu.RLock()
	defer notificationService.subscribersMu.RUnlock()

	for _, sub := range notificationService.subscribers {
		select {
		case sub.Channel() <- notification:
		default:
		}
	}
}
