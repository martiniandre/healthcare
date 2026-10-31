package ratelimit

import (
	"context"
	"fmt"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	redisKeyPrefix        = "ratelimit:"
	redisCallTimeout      = 200 * time.Millisecond
	redisResultFieldCount = 2
)

const slidingWindowLogScript = `
local redisTime = redis.call('TIME')
local nowMicros = tonumber(redisTime[1]) * 1000000 + tonumber(redisTime[2])
local windowMicros = tonumber(ARGV[1])
local maxRequests = tonumber(ARGV[2])
local member = ARGV[3] .. ':' .. nowMicros

redis.call('ZREMRANGEBYSCORE', KEYS[1], '-inf', nowMicros - windowMicros)

if redis.call('ZCARD', KEYS[1]) >= maxRequests then
  local oldest = redis.call('ZRANGE', KEYS[1], 0, 0, 'WITHSCORES')
  local retryMicros = tonumber(oldest[2]) + windowMicros - nowMicros
  if retryMicros < 1 then
    retryMicros = 1
  end
  return {0, retryMicros}
end

redis.call('ZADD', KEYS[1], nowMicros, member)
redis.call('PEXPIRE', KEYS[1], math.ceil(windowMicros / 1000))
return {1, 0}
`

type RedisLimiter struct {
	redisClient    redis.Scripter
	script         *redis.Script
	memberSequence atomic.Uint64
	instanceID     string
}

func NewRedisLimiter(redisClient redis.Scripter, instanceID string) *RedisLimiter {
	return &RedisLimiter{
		redisClient: redisClient,
		script:      redis.NewScript(slidingWindowLogScript),
		instanceID:  instanceID,
	}
}

func (limiter *RedisLimiter) Allow(ctx context.Context, key string, budget Budget) (bool, time.Duration, error) {
	callContext, cancelCall := context.WithTimeout(ctx, redisCallTimeout)
	defer cancelCall()

	uniqueMember := fmt.Sprintf("%s:%d", limiter.instanceID, limiter.memberSequence.Add(1))
	scriptResult, scriptError := limiter.script.Run(
		callContext,
		limiter.redisClient,
		[]string{redisKeyPrefix + key},
		budget.Window.Microseconds(),
		budget.MaxRequests,
		uniqueMember,
	).Slice()
	if scriptError != nil {
		return false, 0, scriptError
	}
	if len(scriptResult) != redisResultFieldCount {
		return false, 0, fmt.Errorf("unexpected rate limit script result length %d", len(scriptResult))
	}

	allowedFlag, flagIsInteger := scriptResult[0].(int64)
	retryMicros, retryIsInteger := scriptResult[1].(int64)
	if !flagIsInteger || !retryIsInteger {
		return false, 0, fmt.Errorf("unexpected rate limit script result types %T %T", scriptResult[0], scriptResult[1])
	}

	return allowedFlag == 1, time.Duration(retryMicros) * time.Microsecond, nil
}

func newInstanceID() string {
	return strconv.FormatInt(time.Now().UnixNano(), 36)
}
