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
		title, _ := event.Data["title"].(string)
		body, _ := event.Data["body"].(string)
		actorID := parseActorID(event.Data)
		resourceType, _ := event.Data["resource_type"].(string)
		resourceID, _ := event.Data["resource_id"].(string)
		_, err := svc.CreateNotificationByRole(ctx, notificationType, title, body, actorID, resourceType, resourceID)
		return err
	}
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
