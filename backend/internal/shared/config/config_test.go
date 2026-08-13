package config

import "testing"

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
