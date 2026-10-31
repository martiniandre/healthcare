package outbox

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"math"
	"math/rand"
	"sync"
	"time"

	"github.com/healthcare/backend/internal/shared/eventbus"
)

const (
	MaxAttempts      = 8
	DefaultBatchSize = 50
	DefaultPollEvery = time.Second
	DefaultLease     = time.Minute
	baseBackoff      = 2 * time.Second
	maxBackoff       = 5 * time.Minute
	jitterRatio      = 0.2
	processedTTL     = 24 * time.Hour
	purgeEvery       = 10 * time.Minute
)

type HandlerRegistry interface {
	HandlersFor(eventName string) []eventbus.Handler
}

type Options struct {
	BatchSize    int
	PollInterval time.Duration
	Lease        time.Duration
	MaxAttempts  int
	Now          func() time.Time
	Jitter       func() float64
}

type Dispatcher struct {
	store    Store
	registry HandlerRegistry
	options  Options

	wake    chan struct{}
	mu      sync.Mutex
	cancel  context.CancelFunc
	done    chan struct{}
	started bool
}

func NewDispatcher(store Store, registry HandlerRegistry, options Options) *Dispatcher {
	if options.BatchSize <= 0 {
		options.BatchSize = DefaultBatchSize
	}
	if options.PollInterval <= 0 {
		options.PollInterval = DefaultPollEvery
	}
	if options.Lease <= 0 {
		options.Lease = DefaultLease
	}
	if options.MaxAttempts <= 0 {
		options.MaxAttempts = MaxAttempts
	}
	if options.Now == nil {
		options.Now = time.Now
	}
	if options.Jitter == nil {
		options.Jitter = rand.Float64
	}
	return &Dispatcher{
		store:    store,
		registry: registry,
		options:  options,
		wake:     make(chan struct{}, 1),
		done:     make(chan struct{}),
	}
}

func BackoffDelay(attempts int, jitterSample float64) time.Duration {
	if attempts < 1 {
		attempts = 1
	}
	delay := float64(baseBackoff) * math.Pow(2, float64(attempts-1))
	if delay > float64(maxBackoff) {
		delay = float64(maxBackoff)
	}
	jitterFactor := 1 + (jitterSample*2-1)*jitterRatio
	return time.Duration(delay * jitterFactor)
}

func (dispatcher *Dispatcher) Wake() {
	select {
	case dispatcher.wake <- struct{}{}:
	default:
	}
}

func (dispatcher *Dispatcher) Start(ctx context.Context) {
	dispatcher.mu.Lock()
	if dispatcher.started {
		dispatcher.mu.Unlock()
		return
	}
	dispatcher.started = true
	loopContext, cancel := context.WithCancel(ctx)
	dispatcher.cancel = cancel
	dispatcher.mu.Unlock()
	defer close(dispatcher.done)

	slog.Info("outbox dispatcher started", "poll_interval", dispatcher.options.PollInterval, "batch_size", dispatcher.options.BatchSize)

	ticker := time.NewTicker(dispatcher.options.PollInterval)
	defer ticker.Stop()
	lastPurge := dispatcher.options.Now()

	for {
		dispatcher.drain(loopContext)
		if dispatcher.options.Now().Sub(lastPurge) >= purgeEvery {
			dispatcher.purge(loopContext)
			lastPurge = dispatcher.options.Now()
		}
		select {
		case <-loopContext.Done():
			slog.Info("outbox dispatcher stopped")
			return
		case <-ticker.C:
		case <-dispatcher.wake:
		}
	}
}

func (dispatcher *Dispatcher) Stop() {
	dispatcher.mu.Lock()
	started := dispatcher.started
	cancel := dispatcher.cancel
	dispatcher.mu.Unlock()
	if !started {
		return
	}
	cancel()
	<-dispatcher.done
}

func (dispatcher *Dispatcher) drain(loopContext context.Context) {
	for loopContext.Err() == nil {
		processed := dispatcher.DispatchOnce(loopContext)
		if processed < dispatcher.options.BatchSize {
			return
		}
	}
}

func (dispatcher *Dispatcher) DispatchOnce(loopContext context.Context) int {
	claimed, claimError := dispatcher.store.ClaimBatch(loopContext, dispatcher.options.BatchSize, dispatcher.options.Lease)
	if claimError != nil {
		if loopContext.Err() == nil {
			slog.Error("outbox claim failed", "error", claimError)
		}
		return 0
	}
	for _, record := range claimed {
		dispatcher.deliver(record)
	}
	return len(claimed)
}

func (dispatcher *Dispatcher) deliver(record Record) {
	deliveryContext, cancel := context.WithTimeout(context.Background(), dispatcher.options.Lease)
	defer cancel()

	deliveryError := dispatcher.invokeHandlers(deliveryContext, record)
	if deliveryError == nil {
		if markError := dispatcher.store.MarkProcessed(deliveryContext, record.ID); markError != nil {
			slog.Error("outbox mark processed failed", "event", record.Name, "attempt", record.Attempts, "error", markError)
		}
		return
	}

	failure := Failure{ID: record.ID, Cause: deliveryError}
	if record.Attempts >= dispatcher.options.MaxAttempts {
		failure.Dead = true
		slog.Error("outbox event moved to dead letter", "event", record.Name, "attempt", record.Attempts, "error", deliveryError)
	} else {
		retryDelay := BackoffDelay(record.Attempts, dispatcher.options.Jitter())
		failure.RetryAt = dispatcher.options.Now().Add(retryDelay)
		slog.Warn("outbox delivery failed, will retry", "event", record.Name, "attempt", record.Attempts, "retry_in", retryDelay, "error", deliveryError)
	}
	if markError := dispatcher.store.MarkFailed(deliveryContext, failure); markError != nil {
		slog.Error("outbox mark failed failed", "event", record.Name, "attempt", record.Attempts, "error", markError)
	}
}

func (dispatcher *Dispatcher) invokeHandlers(deliveryContext context.Context, record Record) error {
	var data map[string]any
	if decodeError := json.Unmarshal(record.Payload, &data); decodeError != nil {
		return fmt.Errorf("decode payload: %w", decodeError)
	}
	handlers := dispatcher.registry.HandlersFor(record.Name)
	if len(handlers) == 0 {
		slog.Warn("outbox event has no handlers, discarding", "event", record.Name)
		return nil
	}

	event := eventbus.Event{ID: record.ID, Name: record.Name, Data: data}
	handlerContext := eventbus.ContextWithEventID(deliveryContext, record.ID)

	var handlerErrors []error
	for _, handler := range handlers {
		handlerErrors = append(handlerErrors, runHandler(handlerContext, handler, event))
	}
	return errors.Join(handlerErrors...)
}

func runHandler(handlerContext context.Context, handler eventbus.Handler, event eventbus.Event) (handlerError error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			handlerError = fmt.Errorf("handler panic: %v", recovered)
		}
	}()
	return handler(handlerContext, event)
}

func (dispatcher *Dispatcher) purge(loopContext context.Context) {
	purged, purgeError := dispatcher.store.PurgeProcessed(loopContext, processedTTL)
	if purgeError != nil {
		if loopContext.Err() == nil {
			slog.Warn("outbox purge failed", "error", purgeError)
		}
		return
	}
	if purged > 0 {
		slog.Info("outbox purged processed events", "count", purged)
	}
}

func (dispatcher *Dispatcher) Stats(ctx context.Context) (Stats, error) {
	return dispatcher.store.Stats(ctx)
}
