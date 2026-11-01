package integration

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/healthcare/backend/internal/shared/database"
	"github.com/jackc/pgx/v5"
)

const (
	isolatedDatabasePrefix       = "healthcare_test_"
	isolatedDatabaseSuffixLength = 16
)

func newIsolatedDatabaseName() (string, error) {
	randomBytes := make([]byte, isolatedDatabaseSuffixLength/2)
	if _, readError := rand.Read(randomBytes); readError != nil {
		return "", fmt.Errorf("failed to generate isolated database name: %w", readError)
	}
	return isolatedDatabasePrefix + hex.EncodeToString(randomBytes), nil
}

func isIsolatedDatabaseName(databaseName string) bool {
	randomSuffix, hasPrefix := strings.CutPrefix(databaseName, isolatedDatabasePrefix)
	if !hasPrefix || len(randomSuffix) != isolatedDatabaseSuffixLength {
		return false
	}
	_, decodeError := hex.DecodeString(randomSuffix)
	return decodeError == nil
}

func databaseURLWithName(baseDatabaseURL string, databaseName string) (string, error) {
	parsedURL, parseError := url.Parse(baseDatabaseURL)
	if parseError != nil {
		return "", fmt.Errorf("invalid database url: %w", parseError)
	}
	parsedURL.Path = "/" + databaseName
	parsedURL.RawPath = ""
	return parsedURL.String(), nil
}

func dropIsolatedDatabase(adminDatabaseURL string, databaseName string) error {
	if !isIsolatedDatabaseName(databaseName) {
		return fmt.Errorf("refusing to drop database %q: it was not created by the integration tests", databaseName)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	adminPool, connectError := database.Connect(ctx, adminDatabaseURL)
	if connectError != nil {
		return fmt.Errorf("failed to connect to admin database: %w", connectError)
	}
	defer adminPool.Close()

	_, dropError := adminPool.Exec(ctx, "DROP DATABASE IF EXISTS "+pgx.Identifier{databaseName}.Sanitize()+" WITH (FORCE)")
	if dropError != nil {
		return fmt.Errorf("failed to drop database %s: %w", databaseName, dropError)
	}
	return nil
}

func (environment *testEnvironment) teardown() {
	environment.teardownOnce.Do(func() {
		if environment.pool != nil {
			environment.pool.Close()
		}
		if environment.redisClient != nil {
			_ = environment.redisClient.Close()
		}
		if environment.createdDatabaseName == "" {
			return
		}
		if dropError := dropIsolatedDatabase(environment.adminDatabaseURL, environment.createdDatabaseName); dropError != nil {
			fmt.Fprintf(os.Stderr, "integration tests: cleanup failed: %v\n", dropError)
		}
	})
}
