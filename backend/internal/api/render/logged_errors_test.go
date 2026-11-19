package render

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/healthcare/backend/internal/shared/apperrors"
	"github.com/healthcare/backend/internal/shared/ctxkeys"
	"github.com/healthcare/backend/internal/shared/logger"
)

func captureLogs(testingT *testing.T) *bytes.Buffer {
	testingT.Helper()
	logBuffer := &bytes.Buffer{}
	previousLogger := slog.Default()
	slog.SetDefault(slog.New(logger.NewContextHandler(slog.NewJSONHandler(logBuffer, &slog.HandlerOptions{Level: slog.LevelDebug}))))
	testingT.Cleanup(func() { slog.SetDefault(previousLogger) })
	return logBuffer
}

func decodeLogEntry(testingT *testing.T, logBuffer *bytes.Buffer) map[string]any {
	testingT.Helper()
	logEntry := map[string]any{}
	if decodeError := json.Unmarshal(logBuffer.Bytes(), &logEntry); decodeError != nil {
		testingT.Fatalf("expected a single JSON log entry, got %q: %v", logBuffer.String(), decodeError)
	}
	return logEntry
}

func TestServiceErrorLogsClientErrorsAsWarningWithRequestID(testingT *testing.T) {
	logBuffer := captureLogs(testingT)
	requestContext := context.WithValue(context.Background(), ctxkeys.RequestIDKey, "request-123")
	responseRecorder := httptest.NewRecorder()

	ServiceError(requestContext, responseRecorder, "failed to load patient", apperrors.ErrPatientNotFound, "patient_id", "patient-1")

	if responseRecorder.Code != http.StatusNotFound {
		testingT.Fatalf("expected status 404, got %d", responseRecorder.Code)
	}
	logEntry := decodeLogEntry(testingT, logBuffer)
	if logEntry["level"] != "WARN" {
		testingT.Fatalf("expected WARN level, got %v", logEntry["level"])
	}
	if logEntry["msg"] != "failed to load patient" || logEntry["request_id"] != "request-123" || logEntry["patient_id"] != "patient-1" {
		testingT.Fatalf("unexpected log entry %v", logEntry)
	}
	if logEntry["status"] != float64(http.StatusNotFound) {
		testingT.Fatalf("expected status attribute 404, got %v", logEntry["status"])
	}
}

func TestServiceErrorLogsUnknownErrorsAsServerErrors(testingT *testing.T) {
	logBuffer := captureLogs(testingT)
	responseRecorder := httptest.NewRecorder()

	ServiceError(context.Background(), responseRecorder, "failed to save", errors.New("database unavailable"))

	if responseRecorder.Code != http.StatusInternalServerError {
		testingT.Fatalf("expected status 500, got %d", responseRecorder.Code)
	}
	logEntry := decodeLogEntry(testingT, logBuffer)
	if logEntry["level"] != "ERROR" || logEntry["error"] != "database unavailable" {
		testingT.Fatalf("unexpected log entry %v", logEntry)
	}
	if _, hasRequestID := logEntry["request_id"]; hasRequestID {
		testingT.Fatalf("expected no request_id without one in context, got %v", logEntry)
	}
}

func TestServiceErrorMessageWritesUserMessageAndLevelFromStatus(testingT *testing.T) {
	logBuffer := captureLogs(testingT)
	requestContext := context.WithValue(context.Background(), ctxkeys.CorrelationIDKey, "correlation-9")
	responseRecorder := httptest.NewRecorder()

	ServiceErrorMessage(requestContext, responseRecorder, http.StatusUnauthorized, "Credenciais inválidas.", "login failed", errors.New("bad password"))

	if responseRecorder.Code != http.StatusUnauthorized {
		testingT.Fatalf("expected status 401, got %d", responseRecorder.Code)
	}
	responseBody := map[string]string{}
	json.Unmarshal(responseRecorder.Body.Bytes(), &responseBody)
	if responseBody["error"] != "Credenciais inválidas." {
		testingT.Fatalf("unexpected response body %v", responseBody)
	}
	logEntry := decodeLogEntry(testingT, logBuffer)
	if logEntry["level"] != "WARN" || logEntry["request_id"] != "correlation-9" {
		testingT.Fatalf("unexpected log entry %v", logEntry)
	}
}
