package notifications

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/healthcare/backend/internal/modules/notifications/pb"
	"github.com/healthcare/backend/internal/shared/apperrors"
	"github.com/healthcare/backend/internal/shared/ctxkeys"
)

func mapNotificationError(err error) error {
	if errors.Is(err, ErrNotificationNotFound) {
		return apperrors.ErrNotificationNotFound.ToGRPC()
	}
	return apperrors.ToGRPCStatus(err)
}

type GRPCHandler struct {
	service Service
}

func NewGRPCHandler(service Service) *GRPCHandler {
	return &GRPCHandler{service: service}
}

func (handler *GRPCHandler) ListNotifications(ctx context.Context, req *pb.ListNotificationsRequest) (*pb.ListNotificationsResponse, error) {
	userIDValue := ctx.Value(ctxkeys.UserIDKey)
	userID, ok := userIDValue.(string)
	if !ok || userID == "" {
		return nil, apperrors.ErrMissingToken.ToGRPC()
	}

	parsedUserID, parseError := uuid.Parse(userID)
	if parseError != nil {
		return nil, apperrors.ErrBadRequest.ToGRPC()
	}

	limit := req.Limit
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	offset := req.Offset
	if offset < 0 {
		offset = 0
	}

	notifications, total, listError := handler.service.ListNotifications(ctx, parsedUserID, limit, offset)
	if listError != nil {
		return nil, mapNotificationError(listError)
	}

	notificationResponses := make([]*pb.Notification, 0, len(notifications))
	for _, notification := range notifications {
		notificationResponses = append(notificationResponses, notificationToPB(notification))
	}

	return &pb.ListNotificationsResponse{
		Notifications: notificationResponses,
		Total:         total,
	}, nil
}

func (handler *GRPCHandler) MarkRead(ctx context.Context, req *pb.MarkReadRequest) (*pb.MarkReadResponse, error) {
	userIDValue := ctx.Value(ctxkeys.UserIDKey)
	userID, ok := userIDValue.(string)
	if !ok || userID == "" {
		return nil, apperrors.ErrMissingToken.ToGRPC()
	}

	parsedUserID, parseUserError := uuid.Parse(userID)
	if parseUserError != nil {
		return nil, apperrors.ErrBadRequest.ToGRPC()
	}

	notificationID, parseNotifError := uuid.Parse(req.NotificationId)
	if parseNotifError != nil {
		return nil, apperrors.ErrBadRequest.ToGRPC()
	}

	err := handler.service.MarkRead(ctx, notificationID, parsedUserID)
	if err != nil {
		return nil, mapNotificationError(err)
	}

	return &pb.MarkReadResponse{}, nil
}

func (handler *GRPCHandler) GetUnreadCount(ctx context.Context, req *pb.GetUnreadCountRequest) (*pb.GetUnreadCountResponse, error) {
	userIDValue := ctx.Value(ctxkeys.UserIDKey)
	userID, ok := userIDValue.(string)
	if !ok || userID == "" {
		return nil, apperrors.ErrMissingToken.ToGRPC()
	}

	parsedUserID, parseError := uuid.Parse(userID)
	if parseError != nil {
		return nil, apperrors.ErrBadRequest.ToGRPC()
	}

	count, err := handler.service.GetUnreadCount(ctx, parsedUserID)
	if err != nil {
		return nil, mapNotificationError(err)
	}

	return &pb.GetUnreadCountResponse{Count: count}, nil
}

func notificationToPB(notification *Notification) *pb.Notification {
	var actorID string
	if notification.ActorID != nil {
		actorID = notification.ActorID.String()
	}

	return &pb.Notification{
		Id:           notification.ID.String(),
		Type:         string(notification.Type),
		Priority:     string(notification.Priority),
		Title:        notification.Title,
		Body:         notification.Body,
		ActorId:      actorID,
		ResourceType: notification.ResourceType,
		ResourceId:   notification.ResourceID,
		CreatedAt:    notification.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
