package notifications

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/healthcare/backend/internal/shared/eventbus"
)

type Dependency struct {
	DB       *pgxpool.Pool
	EventBus eventbus.Bus
}

func Register(dep Dependency) (Service, *HTTPHandler) {
	repo := NewRepository(dep.DB)
	svc := NewService(repo)
	httpHandler := NewHTTPHandler(svc)

	dep.EventBus.Subscribe("telemetry.alert", func(ctx context.Context, event eventbus.Event) error {
		title, _ := event.Data["title"].(string)
		body, _ := event.Data["body"].(string)
		actorID := parseActorID(event.Data)
		resourceType, _ := event.Data["resource_type"].(string)
		resourceID, _ := event.Data["resource_id"].(string)
		_, err := svc.CreateNotificationByRole(ctx, NotificationTypeTelemetryAlert, title, body, actorID, resourceType, resourceID)
		return err
	})

	dep.EventBus.Subscribe("exam.complete", func(ctx context.Context, event eventbus.Event) error {
		title, _ := event.Data["title"].(string)
		body, _ := event.Data["body"].(string)
		actorID := parseActorID(event.Data)
		resourceType, _ := event.Data["resource_type"].(string)
		resourceID, _ := event.Data["resource_id"].(string)
		_, err := svc.CreateNotificationByRole(ctx, NotificationTypeExamComplete, title, body, actorID, resourceType, resourceID)
		return err
	})

	dep.EventBus.Subscribe("encounter.created", func(ctx context.Context, event eventbus.Event) error {
		title, _ := event.Data["title"].(string)
		body, _ := event.Data["body"].(string)
		actorID := parseActorID(event.Data)
		resourceType, _ := event.Data["resource_type"].(string)
		resourceID, _ := event.Data["resource_id"].(string)
		_, err := svc.CreateNotificationByRole(ctx, NotificationTypeEncounterCreate, title, body, actorID, resourceType, resourceID)
		return err
	})

	dep.EventBus.Subscribe("patient.created", func(ctx context.Context, event eventbus.Event) error {
		title, _ := event.Data["title"].(string)
		body, _ := event.Data["body"].(string)
		actorID := parseActorID(event.Data)
		resourceType, _ := event.Data["resource_type"].(string)
		resourceID, _ := event.Data["resource_id"].(string)
		_, err := svc.CreateNotificationByRole(ctx, NotificationTypePatientCreate, title, body, actorID, resourceType, resourceID)
		return err
	})

	dep.EventBus.Subscribe("system.notification", func(ctx context.Context, event eventbus.Event) error {
		title, _ := event.Data["title"].(string)
		body, _ := event.Data["body"].(string)
		actorID := parseActorID(event.Data)
		resourceType, _ := event.Data["resource_type"].(string)
		resourceID, _ := event.Data["resource_id"].(string)
		_, err := svc.CreateNotification(ctx, NotificationTypeSystem, title, body, actorID, resourceType, resourceID, nil)
		return err
	})

	return svc, httpHandler
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
