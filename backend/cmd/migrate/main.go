package main

import (
	"log/slog"
	"os"

	"github.com/healthcare/backend/internal/shared/migrations"
)

func main() {
	databaseURL := os.Getenv("DB_URL")
	if databaseURL == "" {
		databaseURL = os.Getenv("DATABASE_URL")
	}
	if databaseURL == "" {
		slog.Error("DB_URL is required")
		os.Exit(1)
	}
	sourcePath := os.Getenv("MIGRATIONS_PATH")
	if sourcePath == "" {
		sourcePath = "file://migrations"
	}
	if err := migrations.RunFromSource(sourcePath, databaseURL); err != nil {
		slog.Error("Failed to run database migrations", "error", err)
		os.Exit(1)
	}
}
