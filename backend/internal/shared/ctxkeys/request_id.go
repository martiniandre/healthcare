package ctxkeys

import "context"

func RequestID(requestContext context.Context) string {
	if requestID, found := requestContext.Value(RequestIDKey).(string); found {
		return requestID
	}
	if correlationID, found := requestContext.Value(CorrelationIDKey).(string); found {
		return correlationID
	}
	return ""
}
