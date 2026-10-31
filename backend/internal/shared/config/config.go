package config

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"sort"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort              string
	HTTPPort             string
	AppEnv               string
	DBUrl                string
	RedisUrl             string
	SentryDSN            string
	JWTSecret            string
	JWTPreviousSecrets   []string
	OTELExporterEndpoint string
	OTELServiceName      string
	GCPProjectID         string
	GCPLocationID        string
	GCPDatasetID         string
	GCPFHIRStore         string
	GCPDICOMStore        string
	GCPVertexModel       string
	GCSBucketName        string
	StorageRoot          string
	FHIRBaseURL          string

	FrontendURL              string
	TrustedProxies           []string
	EnableTelemetrySimulator bool
}

const minimumJWTSecretLength = 32

type jwtSecretEntry struct {
	variableName string
	secretValue  string
}

func Load() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, relying on environment variables")
	}

	cfg := &Config{
		AppPort:              getEnv("APP_PORT", "50051"),
		HTTPPort:             getEnvAny("PORT", "HTTP_PORT", "8080"),
		AppEnv:               getEnv("APP_ENV", "development"),
		DBUrl:                getEnv("DB_URL", ""),
		RedisUrl:             getEnv("REDIS_URL", "localhost:6379"),
		SentryDSN:            getEnv("SENTRY_DSN", ""),
		JWTSecret:            getEnv("JWT_SECRET", ""),
		JWTPreviousSecrets:   splitCommaSeparated(getEnv("JWT_SECRET_PREVIOUS", "")),
		OTELExporterEndpoint: getEnv("OTEL_EXPORTER_OTLP_ENDPOINT", ""),
		OTELServiceName:      getEnv("OTEL_SERVICE_NAME", "healthcare-api"),
		GCPProjectID:         getEnv("GCP_PROJECT_ID", ""),
		GCPLocationID:        getEnv("GCP_LOCATION_ID", "us-central1"),
		GCPDatasetID:         getEnv("GCP_DATASET_ID", ""),
		GCPFHIRStore:         getEnv("GCP_FHIR_STORE_ID", ""),
		GCPDICOMStore:        getEnv("GCP_DICOM_STORE_ID", "default-dicom"),
		GCPVertexModel:       getEnv("GCP_VERTEX_MODEL", "gemini-2.0-flash-001"),
		GCSBucketName:        getEnv("GCS_BUCKET_NAME", "default-bucket"),
		StorageRoot:          getEnv("STORAGE_ROOT", "dev"),
		FHIRBaseURL:          getEnv("FHIR_BASE_URL", ""),

		FrontendURL:              getEnv("FRONTEND_URL", ""),
		TrustedProxies:           splitCommaSeparated(getEnv("TRUSTED_PROXIES", "")),
		EnableTelemetrySimulator: strings.EqualFold(strings.TrimSpace(getEnv("ENABLE_TELEMETRY_SIMULATOR", "")), "true"),
	}

	if validationErr := cfg.validate(); validationErr != nil {
		return nil, validationErr
	}

	return cfg, nil
}

func (cfg *Config) validate() error {
	requiredFields := map[string]string{
		"DB_URL":     cfg.DBUrl,
		"JWT_SECRET": cfg.JWTSecret,
	}

	if cfg.FHIRBaseURL == "" && cfg.AppEnv != "test" {
		requiredFields["GCP_PROJECT_ID"] = cfg.GCPProjectID
		requiredFields["GCP_DATASET_ID"] = cfg.GCPDatasetID
		requiredFields["GCP_FHIR_STORE_ID"] = cfg.GCPFHIRStore
	}

	var missingFields []string
	for fieldName, fieldValue := range requiredFields {
		if fieldValue == "" {
			missingFields = append(missingFields, fieldName)
		}
	}

	if len(missingFields) > 0 {
		sort.Strings(missingFields)
		return fmt.Errorf("missing required environment variables: %v", missingFields)
	}

	jwtSecretEntries := []jwtSecretEntry{{variableName: "JWT_SECRET", secretValue: cfg.JWTSecret}}
	for _, previousSecret := range cfg.JWTPreviousSecrets {
		jwtSecretEntries = append(jwtSecretEntries, jwtSecretEntry{variableName: "JWT_SECRET_PREVIOUS", secretValue: previousSecret})
	}

	for _, secretEntry := range jwtSecretEntries {
		if len(secretEntry.secretValue) >= minimumJWTSecretLength {
			continue
		}
		if cfg.AppEnv == "production" {
			return fmt.Errorf("%s must be at least %d characters long in production", secretEntry.variableName, minimumJWTSecretLength)
		}
		slog.Warn("JWT secret is shorter than the recommended minimum length", "variable", secretEntry.variableName, "minimum_length", minimumJWTSecretLength)
	}

	return nil
}

func splitCommaSeparated(rawValue string) []string {
	var items []string
	for _, item := range strings.Split(rawValue, ",") {
		if trimmedItem := strings.TrimSpace(item); trimmedItem != "" {
			items = append(items, trimmedItem)
		}
	}
	return items
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getEnvAny(keys ...string) string {
	fallback := keys[len(keys)-1]
	for _, key := range keys[:len(keys)-1] {
		if value, exists := os.LookupEnv(key); exists && value != "" {
			return value
		}
	}
	return fallback
}
