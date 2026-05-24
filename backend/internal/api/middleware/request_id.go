package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/healthcare/backend/internal/shared/ctxkeys"
)

const RequestIDHeader = "X-Request-ID"

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
		requestID := httpRequest.Header.Get(RequestIDHeader)
		if requestID == "" {
			requestID = uuid.New().String()
		}

		httpResponseWriter.Header().Set(RequestIDHeader, requestID)

		enrichedContext := context.WithValue(httpRequest.Context(), ctxkeys.CorrelationIDKey, requestID)
		next.ServeHTTP(httpResponseWriter, httpRequest.WithContext(enrichedContext))
	})
}

func GetRequestID(requestContext context.Context) string {
	correlationID, typeAssertionOk := requestContext.Value(ctxkeys.CorrelationIDKey).(string)
	if !typeAssertionOk {
		return ""
	}
	return correlationID
}
