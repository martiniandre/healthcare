package config

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort             string
	HTTPPort            string
	AppEnv              string
	DBUrl               string
	RedisUrl            string
	SentryDSN           string
	JWTSecret           string
	OTELExporterEndpoint string
	OTELServiceName     string
	GCPProjectID        string
	GCPLocationID       string
	GCPDatasetID        string
	GCPFHIRStore        string
	GCPDICOMStore       string
	GCPVertexModel      string
	GCSBucketName       string
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
		OTELExporterEndpoint: getEnv("OTEL_EXPORTER_OTLP_ENDPOINT", ""),
		OTELServiceName:      getEnv("OTEL_SERVICE_NAME", "healthcare-api"),
		GCPProjectID:         getEnv("GCP_PROJECT_ID", ""),
		GCPLocationID:        getEnv("GCP_LOCATION_ID", "us-central1"),
		GCPDatasetID:         getEnv("GCP_DATASET_ID", ""),
		GCPFHIRStore:         getEnv("GCP_FHIR_STORE_ID", ""),
		GCPDICOMStore:        getEnv("GCP_DICOM_STORE_ID", "default-dicom"),
		GCPVertexModel:       getEnv("GCP_VERTEX_MODEL", "gemini-2.0-flash-001"),
		GCSBucketName:        getEnv("GCS_BUCKET_NAME", "default-bucket"),
	}

	if validationErr := cfg.validate(); validationErr != nil {
		return nil, validationErr
	}

	return cfg, nil
}

func (cfg *Config) validate() error {
	requiredFields := map[string]string{
		"DB_URL":            cfg.DBUrl,
		"JWT_SECRET":        cfg.JWTSecret,
		"GCP_PROJECT_ID":    cfg.GCPProjectID,
		"GCP_DATASET_ID":    cfg.GCPDatasetID,
		"GCP_FHIR_STORE_ID": cfg.GCPFHIRStore,
	}

	var missingFields []string
	for fieldName, fieldValue := range requiredFields {
		if fieldValue == "" {
			missingFields = append(missingFields, fieldName)
		}
	}

	if len(missingFields) > 0 {
		return errors.New(fmt.Sprintf("missing required environment variables: %v", missingFields))
	}

	return nil
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
