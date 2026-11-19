package render

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/healthcare/backend/internal/shared/apperrors"
)

func ServiceError(requestContext context.Context, httpResponseWriter http.ResponseWriter, logMessage string, err error, attributes ...any) {
	logFailure(requestContext, statusFromError(err), logMessage, err, attributes)
	ErrorFromAppError(httpResponseWriter, err)
}

func ServiceErrorMessage(requestContext context.Context, httpResponseWriter http.ResponseWriter, statusCode int, userMessage string, logMessage string, err error, attributes ...any) {
	logFailure(requestContext, statusCode, logMessage, err, attributes)
	Error(httpResponseWriter, statusCode, userMessage)
}

func statusFromError(err error) int {
	var appError apperrors.AppError
	if errors.As(err, &appError) && appError.HTTPCode != 0 {
		return appError.HTTPCode
	}
	return http.StatusInternalServerError
}

func logFailure(requestContext context.Context, statusCode int, logMessage string, err error, attributes []any) {
	logAttributes := make([]any, 0, len(attributes)+4)
	logAttributes = append(logAttributes, "error", err, "status", statusCode)
	logAttributes = append(logAttributes, attributes...)
	level := slog.LevelError
	if statusCode < http.StatusInternalServerError {
		level = slog.LevelWarn
	}
	slog.Log(requestContext, level, logMessage, logAttributes...)
}
