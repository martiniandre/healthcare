package migrations

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source"
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

func ShouldRunOnStart(appEnv string, flagValue string) bool {
	if parsed, err := strconv.ParseBool(strings.TrimSpace(flagValue)); err == nil {
		return parsed
	}
	return appEnv != "production"
}

// LatestSourceVersion returns the highest migration version available in the source.
func LatestSourceVersion(sourcePath string) (uint, error) {
	sourceDriver, err := source.Open(sourcePath)
	if err != nil {
		return 0, err
	}
	defer sourceDriver.Close()

	version, err := sourceDriver.First()
	if err != nil {
		return 0, fmt.Errorf("no migrations found in %s: %w", sourcePath, err)
	}
	for {
		next, nextErr := sourceDriver.Next(version)
		if nextErr != nil {
			if errors.Is(nextErr, os.ErrNotExist) {
				return version, nil
			}
			return 0, nextErr
		}
		version = next
	}
}

// EvaluateSchemaVersion compares the applied schema version with the latest available one.
func EvaluateSchemaVersion(current uint, dirty bool, latest uint) error {
	if dirty {
		return fmt.Errorf("database schema is dirty at version %d: fix it and re-run the migrate job", current)
	}
	if current < latest {
		return fmt.Errorf("database schema version %d is behind the latest migration %d: run the migrate job before starting the API", current, latest)
	}
	return nil
}

// Verify checks, without applying anything, that the database schema is up to date.
func Verify(databaseURL string) error {
	return VerifyFromSource("file://migrations", databaseURL)
}

func VerifyFromSource(sourcePath string, databaseURL string) error {
	latest, err := LatestSourceVersion(sourcePath)
	if err != nil {
		return err
	}
	migrationInstance, err := migrate.New(sourcePath, databaseURL)
	if err != nil {
		return err
	}
	defer migrationInstance.Close()

	current, dirty, err := migrationInstance.Version()
	if errors.Is(err, migrate.ErrNilVersion) {
		return fmt.Errorf("database has no migrations applied (latest is %d): run the migrate job before starting the API", latest)
	}
	if err != nil {
		return err
	}
	return EvaluateSchemaVersion(current, dirty, latest)
}
