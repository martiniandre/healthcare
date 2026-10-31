package eventbus

import (
	"context"
	"log/slog"
	"sync"
)

type Event struct {
	ID   string
	Name string
	Data map[string]any
}

type Handler func(ctx context.Context, event Event) error

type Bus interface {
	Publish(ctx context.Context, event Event) error
	Subscribe(eventName string, handler Handler)
}

type bus struct {
	handlers map[string][]Handler
	mu       sync.RWMutex
}

func New() Bus {
	return &bus{
		handlers: make(map[string][]Handler),
	}
}

func (eventBus *bus) Publish(ctx context.Context, event Event) error {
	eventBus.mu.RLock()
	handlers := eventBus.handlers[event.Name]
	eventBus.mu.RUnlock()

	if len(handlers) == 0 {
		return nil
	}

	var firstError error
	for _, handler := range handlers {
		if err := handler(ctx, event); err != nil {
			slog.Error("event handler failed", "event", event.Name, "error", err)
			if firstError == nil {
				firstError = err
			}
		}
	}

	return firstError
}

func (eventBus *bus) Subscribe(eventName string, handler Handler) {
	eventBus.mu.Lock()
	defer eventBus.mu.Unlock()

	eventBus.handlers[eventName] = append(eventBus.handlers[eventName], handler)
}

type eventIDContextKey struct{}

func ContextWithEventID(ctx context.Context, eventID string) context.Context {
	return context.WithValue(ctx, eventIDContextKey{}, eventID)
}

func EventIDFromContext(ctx context.Context) (string, bool) {
	eventID, exists := ctx.Value(eventIDContextKey{}).(string)
	return eventID, exists && eventID != ""
}
