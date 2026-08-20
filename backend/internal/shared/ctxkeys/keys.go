package ctxkeys

type ContextKey string

const (
	UserIDKey        ContextKey = "user_id"
	RoleKey          ContextKey = "role"
	EmailKey         ContextKey = "email"
	CorrelationIDKey ContextKey = "correlation_id"
	RequestIDKey     ContextKey = "request_id"
)
