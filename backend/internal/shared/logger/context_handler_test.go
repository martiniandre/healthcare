package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"strings"
	"testing"

	"github.com/healthcare/backend/internal/shared/ctxkeys"
)

func newCapturingLogger() (*slog.Logger, *bytes.Buffer) {
	logBuffer := &bytes.Buffer{}
	return slog.New(NewContextHandler(slog.NewJSONHandler(logBuffer, nil))), logBuffer
}

func decodeEntries(testingT *testing.T, logBuffer *bytes.Buffer) []map[string]any {
	testingT.Helper()
	entries := []map[string]any{}
	for _, line := range strings.Split(strings.TrimSpace(logBuffer.String()), "\n") {
		entry := map[string]any{}
		if decodeError := json.Unmarshal([]byte(line), &entry); decodeError != nil {
			testingT.Fatalf("invalid log line %q: %v", line, decodeError)
		}
		entries = append(entries, entry)
	}
	return entries
}

func TestContextHandlerAddsRequestIDFromContext(testingT *testing.T) {
	logger, logBuffer := newCapturingLogger()
	requestContext := context.WithValue(context.Background(), ctxkeys.RequestIDKey, "request-42")

	logger.ErrorContext(requestContext, "save failed", "patient_id", "patient-1")

	entry := decodeEntries(testingT, logBuffer)[0]
	if entry["request_id"] != "request-42" || entry["patient_id"] != "patient-1" {
		testingT.Fatalf("unexpected entry %v", entry)
	}
}

func TestContextHandlerKeepsExplicitRequestIDAndSkipsEmptyContext(testingT *testing.T) {
	logger, logBuffer := newCapturingLogger()
	requestContext := context.WithValue(context.Background(), ctxkeys.RequestIDKey, "from-context")

	logger.InfoContext(requestContext, "explicit", "request_id", "explicit-id")
	logger.Info("no context")

	entries := decodeEntries(testingT, logBuffer)
	if entries[0]["request_id"] != "explicit-id" {
		testingT.Fatalf("expected explicit request id to win, got %v", entries[0])
	}
	if strings.Count(logBuffer.String(), "request_id") != 1 {
		testingT.Fatalf("expected a single request_id attribute, got %q", logBuffer.String())
	}
	if _, hasRequestID := entries[1]["request_id"]; hasRequestID {
		testingT.Fatalf("expected no request_id without context, got %v", entries[1])
	}
}

func TestContextHandlerPreservesAttributesAndGroups(testingT *testing.T) {
	logger, logBuffer := newCapturingLogger()
	requestContext := context.WithValue(context.Background(), ctxkeys.CorrelationIDKey, "correlation-7")

	logger.With("module", "imaging").WarnContext(requestContext, "slow")

	entry := decodeEntries(testingT, logBuffer)[0]
	if entry["module"] != "imaging" || entry["request_id"] != "correlation-7" {
		testingT.Fatalf("unexpected entry %v", entry)
	}
}
