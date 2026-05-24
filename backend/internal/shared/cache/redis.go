package cache

import (
	"log/slog"

	"github.com/redis/go-redis/v9"
)

func Connect(redisUrl string) *redis.Client {
	slog.Info("Redis integration is temporarily disabled")
	return nil

	/*
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisUrl,
	})

	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		slog.Warn("Could not connect to Redis, rate limiting will be disabled", "error", err)
		return nil
	}

	slog.Info("Connected to Redis successfully")
	return redisClient
	*/
}
