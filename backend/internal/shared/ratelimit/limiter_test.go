package ratelimit

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testBudget = Budget{MaxRequests: 3, Window: time.Minute}

type manualClock struct{ current time.Time }

func (clock *manualClock) now() time.Time { return clock.current }

func TestMemoryLimiterBlocksAfterBudgetAndReportsRetryAfter(testingInstance *testing.T) {
	clock := &manualClock{current: time.Now()}
	limiter := NewMemoryLimiterWithClock(clock.now)

	for requestIndex := 0; requestIndex < testBudget.MaxRequests; requestIndex++ {
		allowed, _, limiterError := limiter.Allow(context.Background(), "client", testBudget)
		require.NoError(testingInstance, limiterError)
		require.True(testingInstance, allowed)
		clock.current = clock.current.Add(10 * time.Second)
	}

	allowed, retryAfter, _ := limiter.Allow(context.Background(), "client", testBudget)
	assert.False(testingInstance, allowed)
	assert.Equal(testingInstance, 30*time.Second, retryAfter)

	clock.current = clock.current.Add(31 * time.Second)
	allowed, _, _ = limiter.Allow(context.Background(), "client", testBudget)
	assert.True(testingInstance, allowed)
}

func TestMemoryLimiterIsolatesKeys(testingInstance *testing.T) {
	limiter := NewMemoryLimiter()
	for requestIndex := 0; requestIndex < testBudget.MaxRequests; requestIndex++ {
		_, _, _ = limiter.Allow(context.Background(), "first", testBudget)
	}

	allowed, _, _ := limiter.Allow(context.Background(), "second", testBudget)
	assert.True(testingInstance, allowed)
}

func TestMemoryLimiterEvictsIdleBuckets(testingInstance *testing.T) {
	clock := &manualClock{current: time.Now()}
	limiter := NewMemoryLimiterWithClock(clock.now)
	_, _, _ = limiter.Allow(context.Background(), "stale", testBudget)

	clock.current = clock.current.Add(2 * time.Minute)
	_, _, _ = limiter.Allow(context.Background(), "fresh", testBudget)

	assert.Equal(testingInstance, 1, limiter.trackedBucketCount())
}

func mustAllow(testingInstance *testing.T, limiter Limiter) bool {
	allowed, _, limiterError := limiter.Allow(context.Background(), "client", testBudget)
	require.NoError(testingInstance, limiterError)
	return allowed
}

func TestRedisLimiterSharesBudgetAcrossInstances(testingInstance *testing.T) {
	fakeRedis := miniredis.RunT(testingInstance)
	firstClient := redis.NewClient(&redis.Options{Addr: fakeRedis.Addr()})
	secondClient := redis.NewClient(&redis.Options{Addr: fakeRedis.Addr()})
	defer firstClient.Close()
	defer secondClient.Close()
	firstInstance := NewRedisLimiter(firstClient, "a")
	secondInstance := NewRedisLimiter(secondClient, "b")

	require.True(testingInstance, mustAllow(testingInstance, firstInstance))
	require.True(testingInstance, mustAllow(testingInstance, secondInstance))
	require.True(testingInstance, mustAllow(testingInstance, firstInstance))

	allowed, retryAfter, limiterError := secondInstance.Allow(context.Background(), "client", testBudget)
	require.NoError(testingInstance, limiterError)
	assert.False(testingInstance, allowed)
	assert.Greater(testingInstance, retryAfter, time.Duration(0))
	assert.LessOrEqual(testingInstance, retryAfter, testBudget.Window)
}

func TestRedisLimiterSetsTTLAndRecoversAfterWindow(testingInstance *testing.T) {
	fakeRedis := miniredis.RunT(testingInstance)
	redisClient := redis.NewClient(&redis.Options{Addr: fakeRedis.Addr()})
	defer redisClient.Close()
	limiter := NewRedisLimiter(redisClient, "a")

	require.True(testingInstance, mustAllow(testingInstance, limiter))
	assert.Greater(testingInstance, fakeRedis.TTL("ratelimit:client"), time.Duration(0))

	fakeRedis.FastForward(testBudget.Window + time.Second)
	assert.False(testingInstance, fakeRedis.Exists("ratelimit:client"))
	assert.True(testingInstance, mustAllow(testingInstance, limiter))
}

func TestRedisLimiterReturnsErrorWhenRedisIsDown(testingInstance *testing.T) {
	fakeRedis := miniredis.RunT(testingInstance)
	redisClient := redis.NewClient(&redis.Options{Addr: fakeRedis.Addr()})
	defer redisClient.Close()
	limiter := NewRedisLimiter(redisClient, "a")
	fakeRedis.Close()

	_, _, limiterError := limiter.Allow(context.Background(), "client", testBudget)
	assert.Error(testingInstance, limiterError)
}

type scriptedLimiter struct {
	calls int
	err   error
}

func (limiter *scriptedLimiter) Allow(context.Context, string, Budget) (bool, time.Duration, error) {
	limiter.calls++
	return true, 0, limiter.err
}

func TestFallbackLimiterUsesLocalLimiterWhenPrimaryFailsAndSkipsDuringCooldown(testingInstance *testing.T) {
	clock := &manualClock{current: time.Now()}
	primary := &scriptedLimiter{err: errors.New("redis down")}
	fallback := NewMemoryLimiterWithClock(clock.now)
	limiter := NewFallbackLimiterWithClock(primary, fallback, clock.now)

	for requestIndex := 0; requestIndex < testBudget.MaxRequests; requestIndex++ {
		allowed, _, limiterError := limiter.Allow(context.Background(), "client", testBudget)
		require.NoError(testingInstance, limiterError)
		require.True(testingInstance, allowed)
	}
	allowed, _, _ := limiter.Allow(context.Background(), "client", testBudget)

	assert.False(testingInstance, allowed)
	assert.Equal(testingInstance, 1, primary.calls)

	clock.current = clock.current.Add(primaryFailureCooldown + time.Second)
	_, _, _ = limiter.Allow(context.Background(), "client", testBudget)
	assert.Equal(testingInstance, 2, primary.calls)
}

func TestFallbackLimiterPrefersPrimaryWhenHealthy(testingInstance *testing.T) {
	primary := &scriptedLimiter{}
	fallback := &scriptedLimiter{}
	limiter := NewFallbackLimiter(primary, fallback)

	_, _, _ = limiter.Allow(context.Background(), "client", testBudget)

	assert.Equal(testingInstance, 1, primary.calls)
	assert.Equal(testingInstance, 0, fallback.calls)
}

func TestNewDefaultLimiterWithoutRedisIsMemoryOnly(testingInstance *testing.T) {
	assert.IsType(testingInstance, &MemoryLimiter{}, NewDefaultLimiter(nil))
}
