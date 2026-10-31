package ratelimit

import (
	"context"
	"sync"
	"time"
)

const idleBucketSweepInterval = time.Minute

type slidingWindowBucket struct {
	timestamps []time.Time
	window     time.Duration
}

func (bucket *slidingWindowBucket) allow(now time.Time, budget Budget) (bool, time.Duration) {
	windowStart := now.Add(-budget.Window)
	firstAliveIndex := 0
	for firstAliveIndex < len(bucket.timestamps) && !bucket.timestamps[firstAliveIndex].After(windowStart) {
		firstAliveIndex++
	}
	bucket.timestamps = bucket.timestamps[firstAliveIndex:]
	bucket.window = budget.Window

	if len(bucket.timestamps) >= budget.MaxRequests {
		retryAfter := bucket.timestamps[0].Add(budget.Window).Sub(now)
		return false, retryAfter
	}

	bucket.timestamps = append(bucket.timestamps, now)
	return true, 0
}

func (bucket *slidingWindowBucket) isIdle(now time.Time) bool {
	if len(bucket.timestamps) == 0 {
		return true
	}
	return !bucket.timestamps[len(bucket.timestamps)-1].After(now.Add(-bucket.window))
}

type MemoryLimiter struct {
	mu          sync.Mutex
	buckets     map[string]*slidingWindowBucket
	lastSweepAt time.Time
	clock       func() time.Time
}

func NewMemoryLimiter() *MemoryLimiter {
	return NewMemoryLimiterWithClock(time.Now)
}

func NewMemoryLimiterWithClock(clock func() time.Time) *MemoryLimiter {
	return &MemoryLimiter{
		buckets:     make(map[string]*slidingWindowBucket),
		lastSweepAt: clock(),
		clock:       clock,
	}
}

func (limiter *MemoryLimiter) Allow(_ context.Context, key string, budget Budget) (bool, time.Duration, error) {
	limiter.mu.Lock()
	defer limiter.mu.Unlock()

	now := limiter.clock()
	limiter.sweepIdleBucketsLocked(now)

	bucket, exists := limiter.buckets[key]
	if !exists {
		bucket = &slidingWindowBucket{window: budget.Window}
		limiter.buckets[key] = bucket
	}

	allowed, retryAfter := bucket.allow(now, budget)
	return allowed, retryAfter, nil
}

func (limiter *MemoryLimiter) sweepIdleBucketsLocked(now time.Time) {
	if now.Sub(limiter.lastSweepAt) < idleBucketSweepInterval {
		return
	}
	limiter.lastSweepAt = now
	for bucketKey, bucket := range limiter.buckets {
		if bucket.isIdle(now) {
			delete(limiter.buckets, bucketKey)
		}
	}
}

func (limiter *MemoryLimiter) trackedBucketCount() int {
	limiter.mu.Lock()
	defer limiter.mu.Unlock()
	return len(limiter.buckets)
}
