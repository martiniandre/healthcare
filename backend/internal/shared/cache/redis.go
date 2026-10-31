package cache

import (
	"context"
	"log/slog"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

const redisPingTimeout = 3 * time.Second

func buildRedisOptions(redisUrl string) (*redis.Options, error) {
	if strings.HasPrefix(redisUrl, "redis://") || strings.HasPrefix(redisUrl, "rediss://") {
		return redis.ParseURL(redisUrl)
	}
	return &redis.Options{Addr: redisUrl}, nil
}

func Connect(redisUrl string) *redis.Client {
	if strings.TrimSpace(redisUrl) == "" {
		slog.Warn("REDIS_URL is empty, shared rate limiting will use the local in-memory fallback")
		return nil
	}

	redisOptions, optionsError := buildRedisOptions(redisUrl)
	if optionsError != nil {
		slog.Warn("Invalid REDIS_URL, shared rate limiting will use the local in-memory fallback", "error", optionsError)
		return nil
	}

	redisClient := redis.NewClient(redisOptions)

	pingContext, cancelPing := context.WithTimeout(context.Background(), redisPingTimeout)
	defer cancelPing()
	if pingError := redisClient.Ping(pingContext).Err(); pingError != nil {
		slog.Warn("Could not connect to Redis, shared rate limiting will use the local in-memory fallback", "error", pingError)
		_ = redisClient.Close()
		return nil
	}

	slog.Info("Connected to Redis successfully")
	return redisClient
}
