package pb

import "context"

type NotificationServiceServer interface {
	ListNotifications(ctx context.Context, req *ListNotificationsRequest) (*ListNotificationsResponse, error)
	MarkRead(ctx context.Context, req *MarkReadRequest) (*MarkReadResponse, error)
	GetUnreadCount(ctx context.Context, req *GetUnreadCountRequest) (*GetUnreadCountResponse, error)
}

type Notification struct {
	Id           string
	Type         string
	Priority     string
	Title        string
	Body         string
	ActorId      string
	ResourceType string
	ResourceId   string
	IsRead       bool
	CreatedAt    string
}

type ListNotificationsRequest struct {
	Limit  int32
	Offset int32
}

type ListNotificationsResponse struct {
	Notifications []*Notification
	Total         int32
}

type MarkReadRequest struct {
	NotificationId string
}

type MarkReadResponse struct{}

type GetUnreadCountRequest struct{}

type GetUnreadCountResponse struct {
	Count int32
}

func RegisterNotificationServiceServer(_ interface{}, _ NotificationServiceServer) {}
