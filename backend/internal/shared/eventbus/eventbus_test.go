package eventbus

import (
	"context"
	"testing"
)

func TestPublishSubscribe(testingInstance *testing.T) {
	eventBus := New()
	contextParam := context.Background()

	received := make(chan Event, 1)
	eventBus.Subscribe("test.event", func(ctx context.Context, event Event) error {
		received <- event
		return nil
	})

	publishError := eventBus.Publish(contextParam, Event{
		Name: "test.event",
		Data: map[string]any{"key": "value"},
	})

	if publishError != nil {
		testingInstance.Fatal("publish returned error:", publishError)
	}

	select {
	case event := <-received:
		if event.Name != "test.event" {
			testingInstance.Errorf("expected event name 'test.event', got '%s'", event.Name)
		}
		if event.Data["key"] != "value" {
			testingInstance.Errorf("expected data key 'value', got '%v'", event.Data["key"])
		}
	default:
		testingInstance.Fatal("expected event, got none")
	}
}

func TestPublishNoSubscribers(testingInstance *testing.T) {
	eventBus := New()
	contextParam := context.Background()

	publishError := eventBus.Publish(contextParam, Event{
		Name: "unregistered.event",
		Data: map[string]any{},
	})

	if publishError != nil {
		testingInstance.Fatal("publish returned error:", publishError)
	}
}

func TestMultipleSubscribers(testingInstance *testing.T) {
	eventBus := New()
	contextParam := context.Background()

	count := 0
	eventBus.Subscribe("test.event", func(ctx context.Context, event Event) error {
		count++
		return nil
	})
	eventBus.Subscribe("test.event", func(ctx context.Context, event Event) error {
		count++
		return nil
	})

	eventBus.Publish(contextParam, Event{
		Name: "test.event",
		Data: map[string]any{},
	})

	if count != 2 {
		testingInstance.Errorf("expected 2 handler invocations, got %d", count)
	}
}
