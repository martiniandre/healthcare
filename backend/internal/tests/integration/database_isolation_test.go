package integration

import (
	"strings"
	"testing"
)

func TestNewIsolatedDatabaseNameIsPrefixedAndUnique(t *testing.T) {
	seenNames := make(map[string]bool)
	for attempt := 0; attempt < 100; attempt++ {
		databaseName, nameError := newIsolatedDatabaseName()
		if nameError != nil {
			t.Fatalf("failed to generate database name: %v", nameError)
		}
		if !strings.HasPrefix(databaseName, isolatedDatabasePrefix) {
			t.Fatalf("expected prefix %q, got %q", isolatedDatabasePrefix, databaseName)
		}
		if !isIsolatedDatabaseName(databaseName) {
			t.Fatalf("generated name %q must be recognised as isolated", databaseName)
		}
		if seenNames[databaseName] {
			t.Fatalf("generated duplicate database name %q", databaseName)
		}
		seenNames[databaseName] = true
	}
}

func TestIsIsolatedDatabaseNameRejectsForeignNames(t *testing.T) {
	foreignNames := []string{
		"postgres",
		"healthcare",
		"healthcare_test",
		"healthcare_test_",
		"healthcare_test_0123",
		"healthcare_test_0123456789abcdef0",
		"healthcare_test_zzzzzzzzzzzzzzzz",
		"other_healthcare_test_0123456789abcdef",
	}
	for _, foreignName := range foreignNames {
		if isIsolatedDatabaseName(foreignName) {
			t.Fatalf("database %q must not be treated as owned by the tests", foreignName)
		}
	}
}

func TestDatabaseURLWithNameReplacesOnlyTheDatabasePath(t *testing.T) {
	rewrittenURL, rewriteError := databaseURLWithName("postgres://user:pass@db.local:5433/postgres?sslmode=disable", "healthcare_test_0123456789abcdef")
	if rewriteError != nil {
		t.Fatalf("unexpected error: %v", rewriteError)
	}
	expectedURL := "postgres://user:pass@db.local:5433/healthcare_test_0123456789abcdef?sslmode=disable"
	if rewrittenURL != expectedURL {
		t.Fatalf("expected %q, got %q", expectedURL, rewrittenURL)
	}
}

func TestDropIsolatedDatabaseRefusesForeignNames(t *testing.T) {
	dropError := dropIsolatedDatabase("postgres://user:pass@127.0.0.1:1/postgres?sslmode=disable", "healthcare")
	if dropError == nil || !strings.Contains(dropError.Error(), "refusing to drop") {
		t.Fatalf("expected refusal to drop a foreign database, got %v", dropError)
	}
}
