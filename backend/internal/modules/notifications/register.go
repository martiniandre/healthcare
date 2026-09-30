package notifications

import (
	"context"

	"github.com/google/uuid"
	"github.com/healthcare/backend/internal/shared/eventbus"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Dependency struct {
	DB       *pgxpool.Pool
	EventBus eventbus.Bus
}

func Register(dep Dependency) (Service, *HTTPHandler) {
	repo := NewRepository(dep.DB)
	svc := NewService(repo)
	httpHandler := NewHTTPHandler(svc)

	for _, eventDefinition := range notificationEventDefinitions {
		dep.EventBus.Subscribe(eventDefinition.EventName, subscribeByRoleHandler(svc, eventDefinition.NotificationType))
	}

	return svc, httpHandler
}

var notificationEventDefinitions = []NotificationEventDefinition{
	{EventName: "telemetry.alert", NotificationType: NotificationTypeTelemetryAlert},
	{EventName: "exam.complete", NotificationType: NotificationTypeExamComplete},
	{EventName: "encounter.created", NotificationType: NotificationTypeEncounterCreate},
	{EventName: "patient.created", NotificationType: NotificationTypePatientCreate},
	{EventName: "report.ready", NotificationType: NotificationTypeReportReady},
	{EventName: "system.notification", NotificationType: NotificationTypeSystem},
}

func subscribeByRoleHandler(svc Service, notificationType NotificationType) func(ctx context.Context, event eventbus.Event) error {
	return func(ctx context.Context, event eventbus.Event) error {
		content := notificationContentFromEvent(event)
		actorID := parseActorID(event.Data)
		resourceType, _ := event.Data["resource_type"].(string)
		resourceID, _ := event.Data["resource_id"].(string)
		_, err := svc.CreateNotificationByRole(ctx, notificationType, content, actorID, resourceType, resourceID)
		return err
	}
}

func notificationContentFromEvent(event eventbus.Event) NotificationContent {
	content := NotificationContent{
		Params: map[string]any{},
	}
	content.Title, _ = event.Data["title"].(string)
	content.Body, _ = event.Data["body"].(string)
	content.TitleKey, _ = event.Data["title_key"].(string)
	content.BodyKey, _ = event.Data["body_key"].(string)
	if content.TitleKey == "" && content.BodyKey == "" {
		return content
	}
	if eventParams, paramsAreMap := event.Data["params"].(map[string]any); paramsAreMap {
		for paramName, paramValue := range eventParams {
			content.Params[paramName] = paramValue
		}
	}
	if content.Title == "" {
		content.Title = content.TitleKey
	}
	if content.Body == "" {
		content.Body = content.BodyKey
	}
	return content
}

func parseActorID(data map[string]any) *uuid.UUID {
	actorIDStr, exists := data["actor_id"].(string)
	if !exists || actorIDStr == "" {
		return nil
	}
	parsed, err := uuid.Parse(actorIDStr)
	if err != nil {
		return nil
	}
	return &parsed
}
