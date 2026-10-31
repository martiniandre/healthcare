package outbox

import (
	"context"
	"sync"
	"time"

	"github.com/healthcare/backend/internal/shared/eventbus"
)

const enqueueTimeout = 5 * time.Second

type Bus struct {
	store     Store
	handlers  map[string][]eventbus.Handler
	mu        sync.RWMutex
	onEnqueue func()
}

func NewBus(store Store) *Bus {
	return &Bus{
		store:    store,
		handlers: make(map[string][]eventbus.Handler),
	}
}

func (outboxBus *Bus) Subscribe(eventName string, handler eventbus.Handler) {
	outboxBus.mu.Lock()
	defer outboxBus.mu.Unlock()
	outboxBus.handlers[eventName] = append(outboxBus.handlers[eventName], handler)
}

func (outboxBus *Bus) HandlersFor(eventName string) []eventbus.Handler {
	outboxBus.mu.RLock()
	defer outboxBus.mu.RUnlock()
	registered := outboxBus.handlers[eventName]
	copied := make([]eventbus.Handler, len(registered))
	copy(copied, registered)
	return copied
}

func (outboxBus *Bus) SetWakeup(wakeup func()) {
	outboxBus.mu.Lock()
	defer outboxBus.mu.Unlock()
	outboxBus.onEnqueue = wakeup
}

func (outboxBus *Bus) Publish(ctx context.Context, event eventbus.Event) error {
	return outboxBus.PublishWith(ctx, nil, event)
}

func (outboxBus *Bus) PublishWith(ctx context.Context, executor Executor, event eventbus.Event) error {
	outboxBus.mu.RLock()
	hasHandlers := len(outboxBus.handlers[event.Name]) > 0
	wakeup := outboxBus.onEnqueue
	outboxBus.mu.RUnlock()

	if !hasHandlers {
		return nil
	}

	enqueueContext := ctx
	if executor == nil {
		var cancel context.CancelFunc
		enqueueContext, cancel = context.WithTimeout(context.WithoutCancel(ctx), enqueueTimeout)
		defer cancel()
	}
	if _, enqueueError := outboxBus.store.Enqueue(enqueueContext, executor, event); enqueueError != nil {
		return enqueueError
	}
	if wakeup != nil && executor == nil {
		wakeup()
	}
	return nil
}
