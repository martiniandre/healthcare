package ratelimit

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	primaryFailureCooldown = 5 * time.Second
	primaryWarningInterval = 30 * time.Second
)

type FallbackLimiter struct {
	primary  Limiter
	fallback Limiter
	clock    func() time.Time

	mu                  sync.Mutex
	primarySkippedUntil time.Time
	lastWarningAt       time.Time
}

func NewFallbackLimiter(primary Limiter, fallback Limiter) *FallbackLimiter {
	return NewFallbackLimiterWithClock(primary, fallback, time.Now)
}

func NewFallbackLimiterWithClock(primary Limiter, fallback Limiter, clock func() time.Time) *FallbackLimiter {
	return &FallbackLimiter{primary: primary, fallback: fallback, clock: clock}
}

func NewDefaultLimiter(redisClient *redis.Client) Limiter {
	memoryLimiter := NewMemoryLimiter()
	if redisClient == nil {
		return memoryLimiter
	}
	return NewFallbackLimiter(NewRedisLimiter(redisClient, newInstanceID()), memoryLimiter)
}

func (limiter *FallbackLimiter) Allow(ctx context.Context, key string, budget Budget) (bool, time.Duration, error) {
	now := limiter.clock()
	if limiter.shouldSkipPrimary(now) {
		return limiter.fallback.Allow(ctx, key, budget)
	}

	allowed, retryAfter, primaryError := limiter.primary.Allow(ctx, key, budget)
	if primaryError == nil {
		return allowed, retryAfter, nil
	}

	limiter.recordPrimaryFailure(now, primaryError)
	return limiter.fallback.Allow(ctx, key, budget)
}

func (limiter *FallbackLimiter) shouldSkipPrimary(now time.Time) bool {
	limiter.mu.Lock()
	defer limiter.mu.Unlock()
	return now.Before(limiter.primarySkippedUntil)
}

func (limiter *FallbackLimiter) recordPrimaryFailure(now time.Time, primaryError error) {
	limiter.mu.Lock()
	limiter.primarySkippedUntil = now.Add(primaryFailureCooldown)
	shouldWarn := now.Sub(limiter.lastWarningAt) >= primaryWarningInterval
	if shouldWarn {
		limiter.lastWarningAt = now
	}
	limiter.mu.Unlock()

	if shouldWarn {
		slog.Warn("Shared rate limiter unavailable, using local in-memory fallback", "error", primaryError, "cooldown", primaryFailureCooldown)
	}
}
