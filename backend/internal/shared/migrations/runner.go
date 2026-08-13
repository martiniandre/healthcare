package migrations

import (
	"errors"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func Run(databaseURL string) error {
	return RunFromSource("file://migrations", databaseURL)
}

func RunFromSource(sourcePath string, databaseURL string) error {
	migrationInstance, err := migrate.New(sourcePath, databaseURL)
	if err != nil {
		return err
	}
	defer migrationInstance.Close()

	err = migrationInstance.Up()
	if errors.Is(err, migrate.ErrNoChange) {
		slog.Info("Database migrations: no changes to apply")
		return nil
	}
	if err != nil {
		return err
	}

	slog.Info("Database migrations applied successfully")
	return nil
}
