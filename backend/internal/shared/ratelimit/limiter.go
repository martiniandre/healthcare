package ratelimit

import (
	"context"
	"time"
)

type Budget struct {
	MaxRequests int
	Window      time.Duration
}

type Limiter interface {
	Allow(ctx context.Context, key string, budget Budget) (allowed bool, retryAfter time.Duration, err error)
}
