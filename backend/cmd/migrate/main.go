// Command migrate applies pending database migrations and exits.
// It reads DB_URL (or DATABASE_URL) from the environment and MIGRATIONS_PATH
// (default "file://migrations").
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
