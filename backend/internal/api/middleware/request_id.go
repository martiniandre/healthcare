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

		enrichedContext := context.WithValue(httpRequest.Context(), ctxkeys.RequestIDKey, requestID)
		enrichedContext = context.WithValue(enrichedContext, ctxkeys.CorrelationIDKey, requestID)
		next.ServeHTTP(httpResponseWriter, httpRequest.WithContext(enrichedContext))
	})
}

func GetRequestID(requestContext context.Context) string {
	requestID, typeAssertionOk := requestContext.Value(ctxkeys.RequestIDKey).(string)
	if typeAssertionOk {
		return requestID
	}
	correlationID, correlationOk := requestContext.Value(ctxkeys.CorrelationIDKey).(string)
	if correlationOk {
		return correlationID
	}
	return ""
}
