package notifications

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/healthcare/backend/internal/app/policy"
)

var ErrNotificationNotFound = errors.New("notification not found")
var ErrInvalidNotificationType = errors.New("invalid notification type")

var notificationPriorityDefaults = map[NotificationType]NotificationPriority{
	NotificationTypeTelemetryAlert:  PriorityCritical,
	NotificationTypeExamComplete:    PriorityHigh,
	NotificationTypeEncounterCreate: PriorityMedium,
	NotificationTypeEncounterUpdate: PriorityMedium,
	NotificationTypePatientCreate:   PriorityLow,
	NotificationTypePatientUpdate:   PriorityLow,
	NotificationTypeAuditAlert:      PriorityHigh,
	NotificationTypeReportReady:     PriorityHigh,
	NotificationTypeSystem:          PriorityLow,
}

type Subscriber interface {
	ID() string
	UserID() uuid.UUID
	Channel() chan *Notification
}

type subscriber struct {
	id      string
	userID  uuid.UUID
	channel chan *Notification
}

func (sub *subscriber) ID() string                  { return sub.id }
func (sub *subscriber) UserID() uuid.UUID           { return sub.userID }
func (sub *subscriber) Channel() chan *Notification { return sub.channel }

type Service interface {
	CreateNotification(ctx context.Context, notifType NotificationType, title, body string, actorID *uuid.UUID, resourceType, resourceID string, recipientIDs []uuid.UUID) (*Notification, error)
	CreateNotificationByRole(ctx context.Context, notifType NotificationType, title, body string, actorID *uuid.UUID, resourceType, resourceID string) (*Notification, error)
	ListNotifications(ctx context.Context, userID uuid.UUID, limit, offset int32) ([]*Notification, int32, error)
	ListNotificationEventDefinitions(ctx context.Context) ([]NotificationEventDefinition, error)
	MarkRead(ctx context.Context, notificationID, userID uuid.UUID) error
	GetUnreadCount(ctx context.Context, userID uuid.UUID) (int32, error)
	Subscribe(ctx context.Context, userID uuid.UUID) Subscriber
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
	roles, exists := policy.RolesForNotificationType(string(notifType))
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

func (notificationService *service) ListNotificationEventDefinitions(ctx context.Context) ([]NotificationEventDefinition, error) {
	eventDefinitions := make([]NotificationEventDefinition, 0, len(notificationEventDefinitions))
	for _, eventDefinition := range notificationEventDefinitions {
		recipientRoles := make([]string, 0)
		if allowedRoles, rolesDefined := policy.RolesForNotificationType(string(eventDefinition.NotificationType)); rolesDefined {
			for _, allowedRole := range allowedRoles {
				recipientRoles = append(recipientRoles, string(allowedRole))
			}
		}
		eventDefinitions = append(eventDefinitions, NotificationEventDefinition{
			EventName:        eventDefinition.EventName,
			NotificationType: eventDefinition.NotificationType,
			Priority:         notificationPriorityDefaults[eventDefinition.NotificationType],
			RecipientRoles:   recipientRoles,
		})
	}
	return eventDefinitions, nil
}

func (notificationService *service) MarkRead(ctx context.Context, notificationID, userID uuid.UUID) error {
	return notificationService.repo.MarkRead(ctx, notificationID, userID)
}

func (notificationService *service) GetUnreadCount(ctx context.Context, userID uuid.UUID) (int32, error) {
	return notificationService.repo.GetUnreadCount(ctx, userID)
}

func (notificationService *service) Subscribe(ctx context.Context, userID uuid.UUID) Subscriber {
	sub := &subscriber{
		id:      uuid.New().String(),
		userID:  userID,
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
	recipientUserIDs, recipientError := notificationService.repo.GetRecipientUserIDs(context.Background(), notification.ID)
	if recipientError != nil {
		return
	}

	recipientSet := make(map[uuid.UUID]struct{}, len(recipientUserIDs))
	for _, recipientUserID := range recipientUserIDs {
		recipientSet[recipientUserID] = struct{}{}
	}

	notificationService.subscribersMu.RLock()
	defer notificationService.subscribersMu.RUnlock()

	for _, sub := range notificationService.subscribers {
		if _, isRecipient := recipientSet[sub.UserID()]; !isRecipient {
			continue
		}
		select {
		case sub.Channel() <- notification:
		default:
		}
	}
}
