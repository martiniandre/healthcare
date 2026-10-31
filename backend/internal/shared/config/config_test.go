package config

import (
	"strings"
	"testing"
)

const validProductionJWTSecret = "0123456789abcdef0123456789abcdef"

func TestValidateRequiresDatabaseAndJWTSecret(t *testing.T) {
	cfg := &Config{
		DBUrl:        "postgres://localhost:5432/healthcare",
		JWTSecret:    "secret",
		GCPProjectID: "project",
		GCPDatasetID: "dataset",
		GCPFHIRStore: "store",
		FHIRBaseURL:  "",
	}
	if validationErr := cfg.validate(); validationErr != nil {
		t.Fatalf("expected valid config, got error: %v", validationErr)
	}
}

func TestValidateRejectsMissingDatabase(t *testing.T) {
	cfg := &Config{
		JWTSecret: "secret",
	}
	if validationErr := cfg.validate(); validationErr == nil {
		t.Fatal("expected error when DB_URL is missing")
	}
}

func TestValidateRejectsMissingJWTSecret(t *testing.T) {
	cfg := &Config{
		DBUrl: "postgres://localhost:5432/healthcare",
	}
	if validationErr := cfg.validate(); validationErr == nil {
		t.Fatal("expected error when JWT_SECRET is missing")
	}
}

func TestValidateRequiresGCPFieldsByDefault(t *testing.T) {
	cfg := &Config{
		DBUrl:     "postgres://localhost:5432/healthcare",
		JWTSecret: "secret",
	}
	validationErr := cfg.validate()
	if validationErr == nil {
		t.Fatal("expected error when GCP fields are missing in default mode")
	}
}

func TestValidateAllowsMissingGCPFieldsWhenFHIRBaseURLIsSet(t *testing.T) {
	cfg := &Config{
		DBUrl:       "postgres://localhost:5432/healthcare",
		JWTSecret:   "secret",
		FHIRBaseURL: "http://localhost:8081/fhir",
	}
	if validationErr := cfg.validate(); validationErr != nil {
		t.Fatalf("expected valid config in offline mode, got error: %v", validationErr)
	}
}

func TestValidateAllowsMissingGCPFieldsInTestEnvironment(t *testing.T) {
	cfg := &Config{
		DBUrl:     "postgres://localhost:5432/healthcare",
		JWTSecret: "secret",
		AppEnv:    "test",
	}
	if validationErr := cfg.validate(); validationErr != nil {
		t.Fatalf("expected valid config in test environment, got error: %v", validationErr)
	}
}

func TestValidateReportsMissingFieldsInSortedOrder(t *testing.T) {
	cfg := &Config{AppEnv: "test"}
	validationErr := cfg.validate()
	if validationErr == nil {
		t.Fatal("expected error when required fields are missing")
	}
	if !strings.Contains(validationErr.Error(), "[DB_URL JWT_SECRET]") {
		t.Fatalf("expected sorted missing fields, got: %v", validationErr)
	}
}

func TestValidateRejectsShortJWTSecretInProduction(t *testing.T) {
	cfg := &Config{
		DBUrl:       "postgres://localhost:5432/healthcare",
		JWTSecret:   "short-secret",
		AppEnv:      "production",
		FHIRBaseURL: "http://localhost:8081/fhir",
	}
	if validationErr := cfg.validate(); validationErr == nil {
		t.Fatal("expected error when JWT_SECRET is too short in production")
	}
}

func TestValidateAcceptsLongJWTSecretInProduction(t *testing.T) {
	cfg := &Config{
		DBUrl:       "postgres://localhost:5432/healthcare",
		JWTSecret:   validProductionJWTSecret,
		AppEnv:      "production",
		FHIRBaseURL: "http://localhost:8081/fhir",
	}
	if validationErr := cfg.validate(); validationErr != nil {
		t.Fatalf("expected valid production config, got error: %v", validationErr)
	}
}

func TestSplitCommaSeparatedTrimsAndDropsEmptyEntries(t *testing.T) {
	items := splitCommaSeparated(" 10.0.0.1 , ,192.168.0.0/16,")
	if len(items) != 2 || items[0] != "10.0.0.1" || items[1] != "192.168.0.0/16" {
		t.Fatalf("unexpected items: %v", items)
	}
}
