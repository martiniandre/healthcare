package notifications

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/healthcare/backend/internal/api/middleware"
	"github.com/healthcare/backend/internal/api/render"
	"github.com/healthcare/backend/internal/shared/ctxkeys"
	"github.com/healthcare/backend/internal/shared/role"
)

type HTTPHandler struct {
	service Service
}

func NewHTTPHandler(service Service) *HTTPHandler {
	return &HTTPHandler{service: service}
}

func (handler *HTTPHandler) RegisterRoutes(mux *http.ServeMux) {
	authenticatedUser := middleware.RequireRoles(role.RoleAdmin, role.RoleDoctor, role.RoleNurse, role.RoleReception, role.RolePatient)

	mux.Handle("GET /api/v1/notifications", authenticatedUser(http.HandlerFunc(handler.ListNotifications)))
	mux.Handle("POST /api/v1/notifications/{notificationId}/read", authenticatedUser(http.HandlerFunc(handler.MarkRead)))
	mux.Handle("GET /api/v1/notifications/unread-count", authenticatedUser(http.HandlerFunc(handler.GetUnreadCount)))
	mux.Handle("GET /api/v1/notifications/stream", authenticatedUser(http.HandlerFunc(handler.StreamNotifications)))
}

type notificationResponse struct {
	ID           string `json:"id"`
	Type         string `json:"type"`
	Priority     string `json:"priority"`
	Title        string `json:"title"`
	Body         string `json:"body"`
	ResourceType string `json:"resource_type"`
	ResourceID   string `json:"resource_id"`
	IsRead       bool   `json:"is_read"`
	CreatedAt    string `json:"created_at"`
}

func (handler *HTTPHandler) ListNotifications(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
	authenticatedContext, authPassed := middleware.ValidateHTTPAuth(httpResponseWriter, httpRequest, []role.Role{role.RoleAdmin, role.RoleDoctor, role.RoleNurse, role.RoleReception, role.RolePatient})
	if !authPassed {
		return
	}

	userIDString, _ := authenticatedContext.Value(ctxkeys.UserIDKey).(string)
	userID, parseError := uuid.Parse(userIDString)
	if parseError != nil {
		render.Error(httpResponseWriter, http.StatusBadRequest, "invalid user ID")
		return
	}

	limit := int32(50)
	offset := int32(0)

	if limitParam := httpRequest.URL.Query().Get("limit"); limitParam != "" {
		if parsedLimit, err := strconv.Atoi(limitParam); err == nil && parsedLimit > 0 && parsedLimit <= 100 {
			limit = int32(parsedLimit)
		}
	}
	if offsetParam := httpRequest.URL.Query().Get("offset"); offsetParam != "" {
		if parsedOffset, err := strconv.Atoi(offsetParam); err == nil && parsedOffset >= 0 {
			offset = int32(parsedOffset)
		}
	}

	notifications, total, listError := handler.service.ListNotifications(authenticatedContext, userID, limit, offset)
	if listError != nil {
		slog.Error("failed to list notifications", "error", listError)
		render.Error(httpResponseWriter, http.StatusInternalServerError, "failed to list notifications")
		return
	}

	responseItems := make([]notificationResponse, 0, len(notifications))
	for _, notification := range notifications {
		responseItems = append(responseItems, notificationResponse{
			ID:           notification.ID.String(),
			Type:         string(notification.Type),
			Priority:     string(notification.Priority),
			Title:        notification.Title,
			Body:         notification.Body,
			ResourceType: notification.ResourceType,
			ResourceID:   notification.ResourceID,
			IsRead:       notification.IsRead,
			CreatedAt:    notification.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	render.JSON(httpResponseWriter, http.StatusOK, map[string]any{
		"notifications": responseItems,
		"total":         total,
	})
}

func (handler *HTTPHandler) MarkRead(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
	authenticatedContext, authPassed := middleware.ValidateHTTPAuth(httpResponseWriter, httpRequest, []role.Role{role.RoleAdmin, role.RoleDoctor, role.RoleNurse, role.RoleReception, role.RolePatient})
	if !authPassed {
		return
	}

	userIDString, _ := authenticatedContext.Value(ctxkeys.UserIDKey).(string)
	userID, parseUserError := uuid.Parse(userIDString)
	if parseUserError != nil {
		render.Error(httpResponseWriter, http.StatusBadRequest, "invalid user ID")
		return
	}

	notificationID, parseNotifError := uuid.Parse(httpRequest.PathValue("notificationId"))
	if parseNotifError != nil {
		render.Error(httpResponseWriter, http.StatusBadRequest, "invalid notification ID")
		return
	}

	err := handler.service.MarkRead(authenticatedContext, notificationID, userID)
	if err != nil {
		slog.Error("failed to mark notification as read", "error", err)
		render.Error(httpResponseWriter, http.StatusInternalServerError, "failed to mark as read")
		return
	}

	render.JSON(httpResponseWriter, http.StatusOK, map[string]bool{"success": true})
}

func (handler *HTTPHandler) GetUnreadCount(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
	authenticatedContext, authPassed := middleware.ValidateHTTPAuth(httpResponseWriter, httpRequest, []role.Role{role.RoleAdmin, role.RoleDoctor, role.RoleNurse, role.RoleReception, role.RolePatient})
	if !authPassed {
		return
	}

	userIDString, _ := authenticatedContext.Value(ctxkeys.UserIDKey).(string)
	userID, parseError := uuid.Parse(userIDString)
	if parseError != nil {
		render.Error(httpResponseWriter, http.StatusBadRequest, "invalid user ID")
		return
	}

	count, err := handler.service.GetUnreadCount(authenticatedContext, userID)
	if err != nil {
		slog.Error("failed to get unread count", "error", err)
		render.Error(httpResponseWriter, http.StatusInternalServerError, "failed to get unread count")
		return
	}

	render.JSON(httpResponseWriter, http.StatusOK, map[string]int32{"count": count})
}

func (handler *HTTPHandler) StreamNotifications(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
	authenticatedContext, authPassed := middleware.ValidateHTTPAuth(httpResponseWriter, httpRequest, []role.Role{role.RoleAdmin, role.RoleDoctor, role.RoleNurse, role.RoleReception, role.RolePatient})
	if !authPassed {
		return
	}

	flusher, flushSupported := httpResponseWriter.(http.Flusher)
	if !flushSupported {
		render.Error(httpResponseWriter, http.StatusInternalServerError, "streaming not supported")
		return
	}

	httpResponseWriter.Header().Set("Content-Type", "text/event-stream")
	httpResponseWriter.Header().Set("Cache-Control", "no-cache")
	httpResponseWriter.Header().Set("Connection", "keep-alive")

	notificationChannel := handler.service.Subscribe(authenticatedContext)
	defer handler.service.Unsubscribe(notificationChannel)

	httpResponseWriter.Write([]byte(": connected\n\n"))
	flusher.Flush()

	for {
		select {
		case <-httpRequest.Context().Done():
			return
		case notification, channelOpen := <-notificationChannel.Channel():
			if !channelOpen {
				return
			}

			eventData, jsonError := json.Marshal(notificationResponse{
				ID:           notification.ID.String(),
				Type:         string(notification.Type),
				Priority:     string(notification.Priority),
				Title:        notification.Title,
				Body:         notification.Body,
				ResourceType: notification.ResourceType,
				ResourceID:   notification.ResourceID,
				IsRead:       notification.IsRead,
				CreatedAt:    notification.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			})
			if jsonError != nil {
				slog.Error("failed to marshal SSE event", "error", jsonError)
				continue
			}

			_, writeError := fmt.Fprintf(httpResponseWriter, "event: notification\ndata: %s\n\n", eventData)
			if writeError != nil {
				return
			}
			flusher.Flush()
		}
	}
}
