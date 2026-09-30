package notifications

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/healthcare/backend/internal/shared/role"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Create(ctx context.Context, notification *Notification, recipientIDs []uuid.UUID) error
	ListByUserID(ctx context.Context, userID uuid.UUID, limit, offset int32) ([]*Notification, int32, error)
	MarkRead(ctx context.Context, notificationID, userID uuid.UUID) error
	GetUnreadCount(ctx context.Context, userID uuid.UUID) (int32, error)
	GetUserIDsByRole(ctx context.Context, roles []role.Role) ([]uuid.UUID, error)
	GetUserIDsByResource(ctx context.Context, resourceType, resourceID string) ([]uuid.UUID, error)
	GetRecipientUserIDs(ctx context.Context, notificationID uuid.UUID) ([]uuid.UUID, error)
}

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

func (notificationRepository *repository) Create(ctx context.Context, notification *Notification, recipientIDs []uuid.UUID) error {
	tx, err := notificationRepository.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	insertNotificationQuery := `INSERT INTO notifications (id, type, priority, title, body, title_key, body_key, params, actor_id, resource_type, resource_id, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`
	encodedParams, err := json.Marshal(notification.Params)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, insertNotificationQuery,
		notification.ID, notification.Type, notification.Priority,
		notification.Title, notification.Body, notification.TitleKey, notification.BodyKey, encodedParams,
		notification.ActorID, notification.ResourceType, notification.ResourceID, notification.CreatedAt,
	)
	if err != nil {
		return err
	}

	for _, recipientID := range recipientIDs {
		insertRecipientQuery := `INSERT INTO notification_recipients (notification_id, user_id, is_read, read_at)
			VALUES ($1, $2, $3, $4)`
		_, err = tx.Exec(ctx, insertRecipientQuery,
			notification.ID, recipientID, false, nil,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

type notificationRow struct {
	ID           uuid.UUID
	Type         string
	Priority     string
	Title        string
	Body         string
	TitleKey     string
	BodyKey      string
	Params       []byte
	ActorID      *uuid.UUID
	ResourceType string
	ResourceID   string
	CreatedAt    time.Time
	IsRead       bool
}

func (notificationRepository *repository) ListByUserID(ctx context.Context, userID uuid.UUID, limit, offset int32) ([]*Notification, int32, error) {
	countQuery := `SELECT COUNT(*) FROM notifications n
		INNER JOIN notification_recipients nr ON nr.notification_id = n.id
		WHERE nr.user_id = $1`
	var total int32
	err := notificationRepository.db.QueryRow(ctx, countQuery, userID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	listQuery := `SELECT n.id, n.type, n.priority, n.title, n.body, n.title_key, n.body_key, n.params, n.actor_id, n.resource_type, n.resource_id, n.created_at, nr.is_read
		FROM notifications n
		INNER JOIN notification_recipients nr ON nr.notification_id = n.id
		WHERE nr.user_id = $1
		ORDER BY n.created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := notificationRepository.db.Query(ctx, listQuery, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	notifications := make([]*Notification, 0)
	for rows.Next() {
		var row notificationRow
		err := rows.Scan(
			&row.ID, &row.Type, &row.Priority, &row.Title, &row.Body,
			&row.TitleKey, &row.BodyKey, &row.Params,
			&row.ActorID, &row.ResourceType, &row.ResourceID, &row.CreatedAt, &row.IsRead,
		)
		if err != nil {
			return nil, 0, err
		}
		notificationParams := map[string]any{}
		if len(row.Params) > 0 {
			if unmarshalError := json.Unmarshal(row.Params, &notificationParams); unmarshalError != nil {
				notificationParams = map[string]any{}
			}
		}
		notifications = append(notifications, &Notification{
			ID:           row.ID,
			Type:         NotificationType(row.Type),
			Priority:     NotificationPriority(row.Priority),
			Title:        row.Title,
			Body:         row.Body,
			TitleKey:     row.TitleKey,
			BodyKey:      row.BodyKey,
			Params:       notificationParams,
			ActorID:      row.ActorID,
			ResourceType: row.ResourceType,
			ResourceID:   row.ResourceID,
			IsRead:       row.IsRead,
			CreatedAt:    row.CreatedAt,
		})
	}

	return notifications, total, nil
}

func (notificationRepository *repository) MarkRead(ctx context.Context, notificationID, userID uuid.UUID) error {
	query := `UPDATE notification_recipients SET is_read = true, read_at = NOW()
		WHERE notification_id = $1 AND user_id = $2`
	_, err := notificationRepository.db.Exec(ctx, query, notificationID, userID)
	return err
}

func (notificationRepository *repository) GetUnreadCount(ctx context.Context, userID uuid.UUID) (int32, error) {
	query := `SELECT COUNT(*) FROM notification_recipients WHERE user_id = $1 AND is_read = false`
	var count int32
	err := notificationRepository.db.QueryRow(ctx, query, userID).Scan(&count)
	return count, err
}

func (notificationRepository *repository) GetUserIDsByRole(ctx context.Context, roles []role.Role) ([]uuid.UUID, error) {
	roleStrings := make([]string, 0, len(roles))
	for _, roleValue := range roles {
		roleStrings = append(roleStrings, string(roleValue))
	}

	query := `SELECT id FROM users WHERE role = ANY($1) AND is_active = true`
	rows, err := notificationRepository.db.Query(ctx, query, roleStrings)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	userIDs := make([]uuid.UUID, 0)
	for rows.Next() {
		var userID uuid.UUID
		err := rows.Scan(&userID)
		if err != nil {
			return nil, err
		}
		userIDs = append(userIDs, userID)
	}

	return userIDs, nil
}

func (notificationRepository *repository) GetUserIDsByResource(ctx context.Context, resourceType, resourceID string) ([]uuid.UUID, error) {
	switch resourceType {
	case "room":
		clinicalRoles := []string{string(role.RoleDoctor), string(role.RoleNurse), string(role.RoleAdmin)}
		query := `SELECT id FROM users WHERE role = ANY($1) AND is_active = true`
		rows, err := notificationRepository.db.Query(ctx, query, clinicalRoles)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		userIDs := make([]uuid.UUID, 0)
		for rows.Next() {
			var userID uuid.UUID
			err := rows.Scan(&userID)
			if err != nil {
				return nil, err
			}
			userIDs = append(userIDs, userID)
		}
		return userIDs, nil
	default:
		return nil, nil
	}
}

func (notificationRepository *repository) GetRecipientUserIDs(ctx context.Context, notificationID uuid.UUID) ([]uuid.UUID, error) {
	query := `SELECT user_id FROM notification_recipients WHERE notification_id = $1`
	rows, err := notificationRepository.db.Query(ctx, query, notificationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	userIDs := make([]uuid.UUID, 0)
	for rows.Next() {
		var userID uuid.UUID
		err := rows.Scan(&userID)
		if err != nil {
			return nil, err
		}
		userIDs = append(userIDs, userID)
	}

	return userIDs, nil
}
