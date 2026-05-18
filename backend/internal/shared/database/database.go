package database

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(ctx context.Context, dbUrl string) (*pgxpool.Pool, error) {
	if dbUrl == "" {
		slog.Warn("DB_URL is empty, skipping database connection")
		return nil, nil
	}

	poolConfig, err := pgxpool.ParseConfig(dbUrl)
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}

	slog.Info("Connected to PostgreSQL database successfully")
	return pool, nil
}
