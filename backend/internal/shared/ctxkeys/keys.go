package ctxkeys

type ContextKey string

const (
	UserIDKey        ContextKey = "user_id"
	RoleKey          ContextKey = "role"
	CorrelationIDKey ContextKey = "correlation_id"
)
