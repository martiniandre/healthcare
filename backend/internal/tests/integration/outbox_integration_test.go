package integration

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/healthcare/backend/internal/shared/eventbus"
	"github.com/healthcare/backend/internal/shared/outbox"
)

func newOutboxTestStore(t *testing.T) *outbox.PostgresStore {
	t.Helper()
	environment.ensure(t)
	if _, truncateError := environment.pool.Exec(context.Background(), "TRUNCATE outbox_events"); truncateError != nil {
		t.Fatalf("failed to reset outbox: %v", truncateError)
	}
	return outbox.NewPostgresStore(environment.pool)
}

func TestOutboxEnqueueRollsBackWithTransaction(t *testing.T) {
	store := newOutboxTestStore(t)
	ctx := context.Background()

	transaction, beginError := environment.pool.Begin(ctx)
	if beginError != nil {
		t.Fatal(beginError)
	}
	if _, enqueueError := store.Enqueue(ctx, transaction, eventbus.Event{Name: "appointment.created", Data: map[string]any{"status": "scheduled"}}); enqueueError != nil {
		t.Fatal(enqueueError)
	}
	if rollbackError := transaction.Rollback(ctx); rollbackError != nil {
		t.Fatal(rollbackError)
	}

	stats, statsError := store.Stats(ctx)
	if statsError != nil {
		t.Fatal(statsError)
	}
	if stats.Pending != 0 {
		t.Fatalf("expected rolled back event to vanish, got %d pending", stats.Pending)
	}

	committed, beginError := environment.pool.Begin(ctx)
	if beginError != nil {
		t.Fatal(beginError)
	}
	if _, enqueueError := store.Enqueue(ctx, committed, eventbus.Event{Name: "appointment.created", Data: map[string]any{}}); enqueueError != nil {
		t.Fatal(enqueueError)
	}
	if commitError := committed.Commit(ctx); commitError != nil {
		t.Fatal(commitError)
	}
	stats, _ = store.Stats(ctx)
	if stats.Pending != 1 {
		t.Fatalf("expected committed event to be pending, got %d", stats.Pending)
	}
}

func TestOutboxDispatchersDoNotDoubleDeliverOnPostgres(t *testing.T) {
	store := newOutboxTestStore(t)
	ctx := context.Background()
	bus := outbox.NewBus(store)

	var deliveries sync.Map
	var duplicates int32
	bus.Subscribe("system.notification", func(handlerContext context.Context, event eventbus.Event) error {
		if _, alreadyDelivered := deliveries.LoadOrStore(event.ID, true); alreadyDelivered {
			atomic.AddInt32(&duplicates, 1)
		}
		return nil
	})

	const totalEvents = 120
	for index := 0; index < totalEvents; index++ {
		if publishError := bus.Publish(ctx, eventbus.Event{Name: "system.notification", Data: map[string]any{"index": index}}); publishError != nil {
			t.Fatal(publishError)
		}
	}

	var waitGroup sync.WaitGroup
	for _, batchSize := range []int{5, 9, 13} {
		dispatcher := outbox.NewDispatcher(store, bus, outbox.Options{BatchSize: batchSize})
		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()
			for dispatcher.DispatchOnce(ctx) > 0 {
			}
		}()
	}
	waitGroup.Wait()

	deliveredCount := 0
	deliveries.Range(func(key, value any) bool {
		deliveredCount++
		return true
	})
	if deliveredCount != totalEvents || duplicates != 0 {
		t.Fatalf("expected %d unique deliveries, got %d with %d duplicates", totalEvents, deliveredCount, duplicates)
	}
	stats, _ := store.Stats(ctx)
	if stats.Pending != 0 {
		t.Fatalf("expected nothing pending, got %d", stats.Pending)
	}
}

func TestOutboxRetryAndDeadLetterOnPostgres(t *testing.T) {
	store := newOutboxTestStore(t)
	ctx := context.Background()
	bus := outbox.NewBus(store)

	var attempts int32
	bus.Subscribe("exam.complete", func(handlerContext context.Context, event eventbus.Event) error {
		atomic.AddInt32(&attempts, 1)
		return errors.New("downstream unavailable")
	})
	if publishError := bus.Publish(ctx, eventbus.Event{Name: "exam.complete", Data: map[string]any{}}); publishError != nil {
		t.Fatal(publishError)
	}

	dispatcher := outbox.NewDispatcher(store, bus, outbox.Options{MaxAttempts: 3, Jitter: func() float64 { return 0.5 }})
	dispatcher.DispatchOnce(ctx)
	if dispatcher.DispatchOnce(ctx) != 0 {
		t.Fatal("event must wait for its backoff before the next attempt")
	}

	for attempt := 2; attempt <= 3; attempt++ {
		if _, expediteError := environment.pool.Exec(ctx, "UPDATE outbox_events SET available_at = NOW() WHERE dead_at IS NULL"); expediteError != nil {
			t.Fatal(expediteError)
		}
		dispatcher.DispatchOnce(ctx)
	}

	if atomic.LoadInt32(&attempts) != 3 {
		t.Fatalf("expected 3 attempts, got %d", attempts)
	}
	stats, _ := store.Stats(ctx)
	if stats.Dead != 1 || stats.Pending != 0 {
		t.Fatalf("expected one dead event, got %+v", stats)
	}

	var lastError string
	scanError := environment.pool.QueryRow(ctx, "SELECT last_error FROM outbox_events").Scan(&lastError)
	if scanError != nil || lastError == "" {
		t.Fatalf("expected last_error to be recorded: %v", scanError)
	}
}

func TestOutboxSurvivesDispatcherRestart(t *testing.T) {
	store := newOutboxTestStore(t)
	ctx := context.Background()

	publishingBus := outbox.NewBus(store)
	publishingBus.Subscribe("report.ready", func(handlerContext context.Context, event eventbus.Event) error { return nil })
	if publishError := publishingBus.Publish(ctx, eventbus.Event{Name: "report.ready", Data: map[string]any{"resource_id": "r-1"}}); publishError != nil {
		t.Fatal(publishError)
	}

	restartedBus := outbox.NewBus(store)
	delivered := make(chan eventbus.Event, 1)
	restartedBus.Subscribe("report.ready", func(handlerContext context.Context, event eventbus.Event) error {
		delivered <- event
		return nil
	})
	dispatcher := outbox.NewDispatcher(store, restartedBus, outbox.Options{PollInterval: 20 * time.Millisecond})
	dispatcherContext, cancel := context.WithCancel(ctx)
	defer cancel()
	go dispatcher.Start(dispatcherContext)
	defer dispatcher.Stop()

	select {
	case event := <-delivered:
		if event.Data["resource_id"] != "r-1" || event.ID == "" {
			t.Fatalf("unexpected event after restart: %+v", event)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("event published before restart was not delivered")
	}
}
