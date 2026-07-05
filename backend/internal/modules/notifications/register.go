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

	dep.EventBus.Subscribe("telemetry.alert", subscribeByRoleHandler(svc, NotificationTypeTelemetryAlert))
	dep.EventBus.Subscribe("exam.complete", subscribeByRoleHandler(svc, NotificationTypeExamComplete))
	dep.EventBus.Subscribe("encounter.created", subscribeByRoleHandler(svc, NotificationTypeEncounterCreate))
	dep.EventBus.Subscribe("patient.created", subscribeByRoleHandler(svc, NotificationTypePatientCreate))
	dep.EventBus.Subscribe("system.notification", subscribeByRoleHandler(svc, NotificationTypeSystem))

	return svc, httpHandler
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
