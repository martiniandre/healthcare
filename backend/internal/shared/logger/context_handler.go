package logger

import (
	"context"
	"log/slog"

	"github.com/healthcare/backend/internal/shared/ctxkeys"
)

const requestIDAttributeKey = "request_id"

type ContextHandler struct {
	inner slog.Handler
}

func NewContextHandler(inner slog.Handler) *ContextHandler {
	return &ContextHandler{inner: inner}
}

func (handler *ContextHandler) Enabled(requestContext context.Context, level slog.Level) bool {
	return handler.inner.Enabled(requestContext, level)
}

func (handler *ContextHandler) Handle(requestContext context.Context, record slog.Record) error {
	if requestContext != nil && !recordHasAttribute(record, requestIDAttributeKey) {
		if requestID := ctxkeys.RequestID(requestContext); requestID != "" {
			record.AddAttrs(slog.String(requestIDAttributeKey, requestID))
		}
	}
	return handler.inner.Handle(requestContext, record)
}

func (handler *ContextHandler) WithAttrs(attributes []slog.Attr) slog.Handler {
	return &ContextHandler{inner: handler.inner.WithAttrs(attributes)}
}

func (handler *ContextHandler) WithGroup(groupName string) slog.Handler {
	return &ContextHandler{inner: handler.inner.WithGroup(groupName)}
}

func recordHasAttribute(record slog.Record, attributeKey string) bool {
	found := false
	record.Attrs(func(attribute slog.Attr) bool {
		if attribute.Key == attributeKey {
			found = true
			return false
		}
		return true
	})
	return found
}
