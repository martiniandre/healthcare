package outbox

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/healthcare/backend/internal/shared/eventbus"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type fakeEntry struct {
	record      Record
	availableAt time.Time
	processed   bool
	dead        bool
	lastError   string
}

type fakeStore struct {
	mu           sync.Mutex
	clock        func() time.Time
	entries      []*fakeEntry
	lastExecutor Executor
}

func newFakeStore(clock func() time.Time) *fakeStore {
	return &fakeStore{clock: clock}
}

func (store *fakeStore) Enqueue(ctx context.Context, executor Executor, event eventbus.Event) (string, error) {
	store.mu.Lock()
	defer store.mu.Unlock()
	store.lastExecutor = executor
	payload, marshalError := json.Marshal(event.Data)
	if marshalError != nil {
		return "", marshalError
	}
	eventID := uuid.NewString()
	store.entries = append(store.entries, &fakeEntry{
		record:      Record{ID: eventID, Name: event.Name, Payload: payload, CreatedAt: store.clock()},
		availableAt: store.clock(),
	})
	return eventID, nil
}

func (store *fakeStore) ClaimBatch(ctx context.Context, batchSize int, lease time.Duration) ([]Record, error) {
	store.mu.Lock()
	defer store.mu.Unlock()
	now := store.clock()
	var claimed []Record
	for _, entry := range store.entries {
		if len(claimed) >= batchSize {
			break
		}
		if entry.processed || entry.dead || entry.availableAt.After(now) {
			continue
		}
		entry.record.Attempts++
		entry.availableAt = now.Add(lease)
		claimed = append(claimed, entry.record)
	}
	return claimed, nil
}

func (store *fakeStore) find(eventID string) *fakeEntry {
	for _, entry := range store.entries {
		if entry.record.ID == eventID {
			return entry
		}
	}
	return nil
}

func (store *fakeStore) MarkProcessed(ctx context.Context, eventID string) error {
	store.mu.Lock()
	defer store.mu.Unlock()
	store.find(eventID).processed = true
	return nil
}

func (store *fakeStore) MarkFailed(ctx context.Context, failure Failure) error {
	store.mu.Lock()
	defer store.mu.Unlock()
	entry := store.find(failure.ID)
	entry.lastError = failure.Cause.Error()
	if failure.Dead {
		entry.dead = true
		return nil
	}
	entry.availableAt = failure.RetryAt
	return nil
}

func (store *fakeStore) PurgeProcessed(ctx context.Context, olderThan time.Duration) (int64, error) {
	return 0, nil
}

func (store *fakeStore) Stats(ctx context.Context) (Stats, error) {
	store.mu.Lock()
	defer store.mu.Unlock()
	var stats Stats
	for _, entry := range store.entries {
		switch {
		case entry.dead:
			stats.Dead++
		case !entry.processed:
			stats.Pending++
		}
	}
	return stats, nil
}

type testClock struct {
	mu      sync.Mutex
	current time.Time
}

func (clock *testClock) Now() time.Time {
	clock.mu.Lock()
	defer clock.mu.Unlock()
	return clock.current
}

func (clock *testClock) Advance(duration time.Duration) {
	clock.mu.Lock()
	defer clock.mu.Unlock()
	clock.current = clock.current.Add(duration)
}

func newTestClock() *testClock {
	return &testClock{current: time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)}
}

func fixedJitter() float64 {
	return 0.5
}

func TestBackoffDelayGrowsAndIsCapped(testingInstance *testing.T) {
	previous := time.Duration(0)
	for attempt := 1; attempt <= 5; attempt++ {
		delay := BackoffDelay(attempt, 0.5)
		if delay <= previous {
			testingInstance.Fatalf("expected delay to grow at attempt %d: %v <= %v", attempt, delay, previous)
		}
		previous = delay
	}
	if capped := BackoffDelay(50, 0.5); capped != maxBackoff {
		testingInstance.Fatalf("expected cap %v, got %v", maxBackoff, capped)
	}
}

func TestBackoffDelayJitterBounds(testingInstance *testing.T) {
	base := BackoffDelay(3, 0.5)
	lower := BackoffDelay(3, 0)
	upper := BackoffDelay(3, 1)
	if lower >= base || upper <= base {
		testingInstance.Fatalf("expected jitter around base, got %v %v %v", lower, base, upper)
	}
	if float64(upper) > float64(base)*(1+jitterRatio)+1 || float64(lower) < float64(base)*(1-jitterRatio)-1 {
		testingInstance.Fatalf("jitter out of bounds: %v %v %v", lower, base, upper)
	}
}

func TestDispatcherDeliversEventWithIDAndContext(testingInstance *testing.T) {
	clock := newTestClock()
	store := newFakeStore(clock.Now)
	bus := NewBus(store)

	var receivedEvent eventbus.Event
	var receivedContextID string
	bus.Subscribe("report.ready", func(ctx context.Context, event eventbus.Event) error {
		receivedEvent = event
		receivedContextID, _ = eventbus.EventIDFromContext(ctx)
		return nil
	})

	if publishError := bus.Publish(context.Background(), eventbus.Event{Name: "report.ready", Data: map[string]any{"resource_id": "abc"}}); publishError != nil {
		testingInstance.Fatal(publishError)
	}

	dispatcher := NewDispatcher(store, bus, Options{Now: clock.Now, Jitter: fixedJitter})
	if delivered := dispatcher.DispatchOnce(context.Background()); delivered != 1 {
		testingInstance.Fatalf("expected 1 delivery, got %d", delivered)
	}

	if receivedEvent.ID == "" || receivedEvent.ID != receivedContextID {
		testingInstance.Fatalf("expected event id in event and context, got %q and %q", receivedEvent.ID, receivedContextID)
	}
	if receivedEvent.Data["resource_id"] != "abc" {
		testingInstance.Fatalf("unexpected payload: %v", receivedEvent.Data)
	}
	if dispatcher.DispatchOnce(context.Background()) != 0 {
		testingInstance.Fatal("processed event must not be claimed again")
	}
}

func TestPublishWithoutHandlersDoesNotEnqueue(testingInstance *testing.T) {
	store := newFakeStore(newTestClock().Now)
	bus := NewBus(store)

	if publishError := bus.Publish(context.Background(), eventbus.Event{Name: "nobody.listens"}); publishError != nil {
		testingInstance.Fatal(publishError)
	}
	if len(store.entries) != 0 {
		testingInstance.Fatalf("expected no entries, got %d", len(store.entries))
	}
}

func TestPublishIsDurableWhenContextIsCancelled(testingInstance *testing.T) {
	store := newFakeStore(newTestClock().Now)
	bus := NewBus(store)
	bus.Subscribe("telemetry.alert", func(ctx context.Context, event eventbus.Event) error { return nil })

	cancelledContext, cancel := context.WithCancel(context.Background())
	cancel()
	if publishError := bus.Publish(cancelledContext, eventbus.Event{Name: "telemetry.alert", Data: map[string]any{}}); publishError != nil {
		testingInstance.Fatal(publishError)
	}
	if len(store.entries) != 1 {
		testingInstance.Fatalf("expected event to be persisted, got %d entries", len(store.entries))
	}
}

func TestPublishReturnsEnqueueError(testingInstance *testing.T) {
	bus := NewBus(failingEnqueueStore{})
	bus.Subscribe("telemetry.alert", func(ctx context.Context, event eventbus.Event) error { return nil })

	if publishError := bus.Publish(context.Background(), eventbus.Event{Name: "telemetry.alert"}); publishError == nil {
		testingInstance.Fatal("expected enqueue error to be returned")
	}
}

type failingEnqueueStore struct {
	*fakeStore
}

func (failingEnqueueStore) Enqueue(ctx context.Context, executor Executor, event eventbus.Event) (string, error) {
	return "", errors.New("database unavailable")
}

type stubExecutor struct{}

func (stubExecutor) Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error) {
	return pgconn.CommandTag{}, nil
}

func (stubExecutor) Query(ctx context.Context, sql string, arguments ...any) (pgx.Rows, error) {
	return nil, nil
}

func (stubExecutor) QueryRow(ctx context.Context, sql string, arguments ...any) pgx.Row {
	return nil
}

func TestPublishWakesDispatcherAfterStandalonePublish(testingInstance *testing.T) {
	store := newFakeStore(newTestClock().Now)
	bus := NewBus(store)
	bus.Subscribe("appointment.created", func(ctx context.Context, event eventbus.Event) error { return nil })

	var wakeups int32
	bus.SetWakeup(func() { atomic.AddInt32(&wakeups, 1) })

	if publishError := bus.Publish(context.Background(), eventbus.Event{Name: "appointment.created"}); publishError != nil {
		testingInstance.Fatal(publishError)
	}
	if atomic.LoadInt32(&wakeups) != 1 {
		testingInstance.Fatalf("expected wakeup after standalone publish, got %d", wakeups)
	}
}

func TestPublishWithEnqueuesOnGivenExecutor(testingInstance *testing.T) {
	store := newFakeStore(newTestClock().Now)
	bus := NewBus(store)
	bus.Subscribe("appointment.created", func(ctx context.Context, event eventbus.Event) error { return nil })

	var wakeups int32
	bus.SetWakeup(func() { atomic.AddInt32(&wakeups, 1) })

	transactionExecutor := stubExecutor{}
	if publishError := bus.PublishWith(context.Background(), transactionExecutor, eventbus.Event{Name: "appointment.created"}); publishError != nil {
		testingInstance.Fatal(publishError)
	}
	if store.lastExecutor != Executor(transactionExecutor) {
		testingInstance.Fatal("expected enqueue on the caller transaction executor")
	}
	if atomic.LoadInt32(&wakeups) != 0 {
		testingInstance.Fatal("must not wake dispatcher before the caller transaction commits")
	}
}

func TestDispatcherRetriesWithBackoffThenSucceeds(testingInstance *testing.T) {
	clock := newTestClock()
	store := newFakeStore(clock.Now)
	bus := NewBus(store)

	var seenEventIDs []string
	var calls int
	bus.Subscribe("exam.complete", func(ctx context.Context, event eventbus.Event) error {
		calls++
		seenEventIDs = append(seenEventIDs, event.ID)
		if calls == 1 {
			return errors.New("temporary failure")
		}
		return nil
	})
	if publishError := bus.Publish(context.Background(), eventbus.Event{Name: "exam.complete", Data: map[string]any{}}); publishError != nil {
		testingInstance.Fatal(publishError)
	}

	dispatcher := NewDispatcher(store, bus, Options{Now: clock.Now, Jitter: fixedJitter})
	dispatcher.DispatchOnce(context.Background())

	if store.entries[0].processed || store.entries[0].dead {
		testingInstance.Fatal("failed delivery must stay pending")
	}
	expectedRetryAt := clock.Now().Add(BackoffDelay(1, 0.5))
	if !store.entries[0].availableAt.Equal(expectedRetryAt) {
		testingInstance.Fatalf("expected retry at %v, got %v", expectedRetryAt, store.entries[0].availableAt)
	}

	if dispatcher.DispatchOnce(context.Background()) != 0 {
		testingInstance.Fatal("event must not be redelivered before backoff elapses")
	}

	clock.Advance(BackoffDelay(1, 0.5) + time.Millisecond)
	dispatcher.DispatchOnce(context.Background())

	if !store.entries[0].processed {
		testingInstance.Fatal("expected event to be processed after retry")
	}
	if len(seenEventIDs) != 2 || seenEventIDs[0] != seenEventIDs[1] {
		testingInstance.Fatalf("redelivery must keep the same event id: %v", seenEventIDs)
	}
}

func TestDispatcherMovesToDeadLetterAfterMaxAttempts(testingInstance *testing.T) {
	clock := newTestClock()
	store := newFakeStore(clock.Now)
	bus := NewBus(store)

	var calls int
	bus.Subscribe("patient.created", func(ctx context.Context, event eventbus.Event) error {
		calls++
		return errors.New("permanent failure")
	})
	bus.Publish(context.Background(), eventbus.Event{Name: "patient.created", Data: map[string]any{}})

	dispatcher := NewDispatcher(store, bus, Options{Now: clock.Now, Jitter: fixedJitter})
	for attempt := 1; attempt <= MaxAttempts+2; attempt++ {
		dispatcher.DispatchOnce(context.Background())
		clock.Advance(10 * time.Minute)
	}

	if calls != MaxAttempts {
		testingInstance.Fatalf("expected %d attempts, got %d", MaxAttempts, calls)
	}
	if !store.entries[0].dead {
		testingInstance.Fatal("expected event in dead letter")
	}
	stats, _ := dispatcher.Stats(context.Background())
	if stats.Dead != 1 || stats.Pending != 0 {
		testingInstance.Fatalf("unexpected stats: %+v", stats)
	}
}

func TestDispatcherRecoversFromHandlerPanic(testingInstance *testing.T) {
	clock := newTestClock()
	store := newFakeStore(clock.Now)
	bus := NewBus(store)
	bus.Subscribe("encounter.created", func(ctx context.Context, event eventbus.Event) error {
		panic("boom")
	})
	bus.Publish(context.Background(), eventbus.Event{Name: "encounter.created", Data: map[string]any{}})

	dispatcher := NewDispatcher(store, bus, Options{Now: clock.Now, Jitter: fixedJitter})
	dispatcher.DispatchOnce(context.Background())

	if store.entries[0].processed || store.entries[0].lastError == "" {
		testingInstance.Fatal("panic must be recorded as failure and keep the event pending")
	}
}

func TestDispatcherLeaseExpiryRedeliversCrashedClaim(testingInstance *testing.T) {
	clock := newTestClock()
	store := newFakeStore(clock.Now)
	bus := NewBus(store)

	var calls int
	bus.Subscribe("report.ready", func(ctx context.Context, event eventbus.Event) error {
		calls++
		return nil
	})
	bus.Publish(context.Background(), eventbus.Event{Name: "report.ready", Data: map[string]any{}})

	claimed, _ := store.ClaimBatch(context.Background(), 10, DefaultLease)
	if len(claimed) != 1 {
		testingInstance.Fatalf("expected simulated claim, got %d", len(claimed))
	}

	dispatcher := NewDispatcher(store, bus, Options{Now: clock.Now, Jitter: fixedJitter})
	if dispatcher.DispatchOnce(context.Background()) != 0 {
		testingInstance.Fatal("leased event must not be claimed by another dispatcher")
	}
	clock.Advance(DefaultLease + time.Second)
	if dispatcher.DispatchOnce(context.Background()) != 1 || calls != 1 {
		testingInstance.Fatal("expected redelivery after lease expiry")
	}
}

func TestTwoDispatchersDoNotDoubleDeliver(testingInstance *testing.T) {
	clock := newTestClock()
	store := newFakeStore(clock.Now)
	bus := NewBus(store)

	var deliveries sync.Map
	var duplicates int32
	bus.Subscribe("system.notification", func(ctx context.Context, event eventbus.Event) error {
		if _, alreadyDelivered := deliveries.LoadOrStore(event.ID, true); alreadyDelivered {
			atomic.AddInt32(&duplicates, 1)
		}
		return nil
	})

	const totalEvents = 200
	for index := 0; index < totalEvents; index++ {
		bus.Publish(context.Background(), eventbus.Event{Name: "system.notification", Data: map[string]any{}})
	}

	first := NewDispatcher(store, bus, Options{Now: clock.Now, BatchSize: 7})
	second := NewDispatcher(store, bus, Options{Now: clock.Now, BatchSize: 11})

	var waitGroup sync.WaitGroup
	for _, dispatcher := range []*Dispatcher{first, second} {
		waitGroup.Add(1)
		go func(current *Dispatcher) {
			defer waitGroup.Done()
			for current.DispatchOnce(context.Background()) > 0 {
			}
		}(dispatcher)
	}
	waitGroup.Wait()

	deliveredCount := 0
	deliveries.Range(func(key, value any) bool {
		deliveredCount++
		return true
	})
	if deliveredCount != totalEvents || duplicates != 0 {
		testingInstance.Fatalf("expected %d unique deliveries without duplicates, got %d unique and %d duplicates", totalEvents, deliveredCount, duplicates)
	}
}

func TestDispatcherStartAndStop(testingInstance *testing.T) {
	store := newFakeStore(time.Now)
	bus := NewBus(store)

	delivered := make(chan string, 1)
	bus.Subscribe("report.ready", func(ctx context.Context, event eventbus.Event) error {
		delivered <- event.ID
		return nil
	})

	dispatcher := NewDispatcher(store, bus, Options{PollInterval: 10 * time.Millisecond})
	bus.SetWakeup(dispatcher.Wake)
	go dispatcher.Start(context.Background())

	bus.Publish(context.Background(), eventbus.Event{Name: "report.ready", Data: map[string]any{}})

	select {
	case <-delivered:
	case <-time.After(2 * time.Second):
		testingInstance.Fatal("timed out waiting for delivery")
	}

	stopped := make(chan struct{})
	go func() {
		dispatcher.Stop()
		close(stopped)
	}()
	select {
	case <-stopped:
	case <-time.After(2 * time.Second):
		testingInstance.Fatal("timed out waiting for stop")
	}
}

func TestStopWithoutStartReturns(testingInstance *testing.T) {
	dispatcher := NewDispatcher(newFakeStore(time.Now), NewBus(nil), Options{})
	dispatcher.Stop()
}
