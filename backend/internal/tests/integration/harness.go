package integration

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/healthcare/backend/internal/api"
	"github.com/healthcare/backend/internal/api/middleware"
	"github.com/healthcare/backend/internal/modules/allergy"
	"github.com/healthcare/backend/internal/modules/analytics"
	"github.com/healthcare/backend/internal/modules/audit_logs"
	"github.com/healthcare/backend/internal/modules/auth"
	"github.com/healthcare/backend/internal/modules/condition"
	"github.com/healthcare/backend/internal/modules/diagnostic_report"
	"github.com/healthcare/backend/internal/modules/encounter"
	"github.com/healthcare/backend/internal/modules/exam_analyzer"
	"github.com/healthcare/backend/internal/modules/health"
	"github.com/healthcare/backend/internal/modules/imaging"
	"github.com/healthcare/backend/internal/modules/medication"
	"github.com/healthcare/backend/internal/modules/notifications"
	"github.com/healthcare/backend/internal/modules/observation"
	"github.com/healthcare/backend/internal/modules/patients"
	"github.com/healthcare/backend/internal/modules/portal"
	"github.com/healthcare/backend/internal/modules/schedule"
	"github.com/healthcare/backend/internal/modules/staff"
	"github.com/healthcare/backend/internal/modules/telemetry"
	"github.com/healthcare/backend/internal/shared/cache"
	"github.com/healthcare/backend/internal/shared/database"
	"github.com/healthcare/backend/internal/shared/eventbus"
	"github.com/healthcare/backend/internal/shared/migrations"
	"github.com/healthcare/backend/internal/shared/storage"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"

	"github.com/google/uuid"
)

const (
	defaultAdminDatabaseURL = "postgres://postgres:mauroloremi@localhost:5432/postgres?sslmode=disable"
	testDatabaseName        = "healthcare_test"
	migrationsSourcePath    = "file://../../../migrations"
	testJWTSecret           = "integration-test-secret"
)

var environment = &testEnvironment{}

type testEnvironment struct {
	once        sync.Once
	valid       bool
	pool        *pgxpool.Pool
	redisClient *redis.Client
	databaseURL string
}

func getEnvOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func (environment *testEnvironment) ensure(t *testing.T) {
	environment.once.Do(func() {
		t.Log("preparing integration test environment")

		adminDatabaseURL := getEnvOrDefault("TEST_ADMIN_DB_URL", defaultAdminDatabaseURL)
		parsedAdminURL, parseError := url.Parse(adminDatabaseURL)
		if parseError != nil {
			t.Logf("integration tests skipped: invalid admin database url: %v", parseError)
			return
		}
		parsedAdminURL.Path = "/" + testDatabaseName
		testDatabaseURL := getEnvOrDefault("TEST_DB_URL", parsedAdminURL.String())

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		pool, connectError := database.Connect(ctx, testDatabaseURL)
		if connectError != nil {
			adminPool, adminConnectError := database.Connect(ctx, adminDatabaseURL)
			if adminConnectError != nil {
				t.Logf("integration tests skipped: postgres not reachable: %v", adminConnectError)
				return
			}
			_, createError := adminPool.Exec(ctx, "CREATE DATABASE "+testDatabaseName)
			adminPool.Close()
			if createError != nil && !strings.Contains(createError.Error(), "already exists") {
				t.Logf("integration tests skipped: failed to create test database: %v", createError)
				return
			}
			pool, connectError = database.Connect(ctx, testDatabaseURL)
			if connectError != nil {
				t.Logf("integration tests skipped: cannot connect to test database: %v", connectError)
				return
			}
		}

		if resetError := resetTestSchemaIfStale(pool, testDatabaseURL); resetError != nil {
			t.Logf("integration tests skipped: failed to reset stale test schema: %v", resetError)
			pool.Close()
			return
		}

		if migrationError := migrations.RunFromSource(migrationsSourcePath, testDatabaseURL); migrationError != nil {
			t.Logf("integration tests skipped: migrations failed: %v", migrationError)
			pool.Close()
			return
		}

		if initJWTServerError := auth.InitJWT(testJWTSecret); initJWTServerError != nil {
			t.Logf("integration tests skipped: jwt initialization failed: %v", initJWTServerError)
			pool.Close()
			return
		}

		redisClient := cache.Connect(getEnvOrDefault("REDIS_URL", "localhost:6379"))
		if redisClient != nil {
			if pingError := redisClient.Ping(ctx).Err(); pingError != nil {
				t.Logf("redis ping failed, continuing without cache: %v", pingError)
			}
		}

		environment.pool = pool
		environment.redisClient = redisClient
		environment.databaseURL = testDatabaseURL
		environment.valid = true
	})

	if !environment.valid {
		t.Skip("integration tests skipped: postgres is required")
	}
}

type testServer struct {
	handler    http.Handler
	db         *pgxpool.Pool
	fhir       *inMemoryFHIRClient
	grpcServer *grpc.Server
}

func newTestServer(t *testing.T) *testServer {
	t.Helper()
	environment.ensure(t)

	grpcServer := grpc.NewServer()
	eventBus := eventbus.New()
	fhirClient := newInMemoryFHIRClient()

	authService := auth.Register(grpcServer, auth.Dependency{DB: environment.pool, EventBus: eventBus})
	staffService := staff.Register(grpcServer, staff.Dependency{DB: environment.pool, FHIRClient: fhirClient})
	patientsService := patients.Register(grpcServer, patients.Dependency{FHIRClient: fhirClient, EventBus: eventBus})
	encounterService := encounter.Register(grpcServer, encounter.Dependency{FHIRClient: fhirClient, EventBus: eventBus})
	observationService := observation.Register(grpcServer, observation.Dependency{FHIRClient: fhirClient})
	conditionService := condition.Register(grpcServer, condition.Dependency{FHIRClient: fhirClient})
	allergyService := allergy.Register(grpcServer, allergy.Dependency{FHIRClient: fhirClient})
	medicationService := medication.Register(grpcServer, medication.Dependency{FHIRClient: fhirClient})
	auditLogsService := audit_logs.Register(grpcServer, audit_logs.Dependency{DB: environment.pool})
	middleware.SetHTTPAuditRecorder(audit_logs.NewHTTPAuditRecorder(auditLogsService))
	diagnosticReportService := diagnostic_report.Register(grpcServer, diagnostic_report.Dependency{
		FHIRClient:   fhirClient,
		EventBus:     eventBus,
		DB:           environment.pool,
		AuditService: auditLogsService,
	})
	storageClient := storage.NewStorageClient("integration-test-bucket", "integration-test")
	imagingService := imaging.Register(grpcServer, imaging.Dependency{
		DB:           environment.pool,
		Storage:      storageClient,
		Redis:        environment.redisClient,
		AuditService: auditLogsService,
	})
	telemetryService := telemetry.Register(grpcServer, telemetry.Dependency{DB: environment.pool, EventBus: eventBus})
	health.Register(grpcServer, health.Dependency{DB: environment.pool, Redis: environment.redisClient})
	analyticsHTTPHandler := analytics.Register(analytics.Dependency{DB: environment.pool, FHIRClient: fhirClient})
	portalHTTPHandler := portal.Register(portal.Dependency{FHIRClient: fhirClient, DB: environment.pool})
	_, notificationsHTTPHandler := notifications.Register(notifications.Dependency{DB: environment.pool, EventBus: eventBus})
	scheduleHTTPHandler := schedule.Register(schedule.Dependency{DB: environment.pool, EventBus: eventBus, AuditService: auditLogsService})
	examAnalyzerRepository, examAnalyzerService, examAnalyzerWorker := exam_analyzer.Register(grpcServer, exam_analyzer.Dependency{DB: environment.pool, Storage: storageClient, EventBus: eventBus})

	secureCookies := false
	middleware.ResetRateLimits()

	router := api.NewRouter(
		secureCookies,
		true,
		auth.NewHTTPHandler(authService, secureCookies),
		patients.NewHTTPHandler(patientsService),
		encounter.NewHTTPHandler(encounterService),
		observation.NewHTTPHandler(observationService),
		condition.NewHTTPHandler(conditionService),
		allergy.NewHTTPHandler(allergyService),
		medication.NewHTTPHandler(medicationService),
		diagnostic_report.NewHTTPHandler(diagnosticReportService),
		imaging.NewHTTPHandler(imagingService),
		staff.NewHTTPHandler(staffService),
		telemetry.NewHTTPHandler(telemetryService),
		exam_analyzer.NewHTTPHandler(examAnalyzerRepository, examAnalyzerService, examAnalyzerWorker),
		analyticsHTTPHandler,
		portalHTTPHandler,
		audit_logs.NewHTTPHandler(auditLogsService),
		notificationsHTTPHandler,
		scheduleHTTPHandler,
		health.NewHTTPHandler(environment.pool, environment.redisClient),
	)

	seedTestUsers(t, environment.pool)
	resetTestData(t, environment.pool)
	fhirClient.reset()
	seedPortalPatientIdentity(t, environment.pool, fhirClient)

	return &testServer{
		handler:    router,
		db:         environment.pool,
		fhir:       fhirClient,
		grpcServer: grpcServer,
	}
}

func resetTestSchemaIfStale(pool *pgxpool.Pool, testDatabaseURL string) error {
	parsedTestURL, parseError := url.Parse(testDatabaseURL)
	if parseError != nil {
		return parseError
	}
	databaseName := strings.TrimPrefix(parsedTestURL.Path, "/")
	if !strings.HasSuffix(databaseName, "test") {
		return fmt.Errorf("refusing to reset schema of non-test database %q", databaseName)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	_, execError := pool.Exec(ctx, "DROP SCHEMA public CASCADE; CREATE SCHEMA public;")
	return execError
}

func seedTestUsers(t *testing.T, pool *pgxpool.Pool) {
	t.Helper()
	ctx := context.Background()

	seedUsers := []struct {
		email    string
		password string
		fullName string
		role     string
	}{
		{email: "admin@hospital.com", password: "secret123", fullName: "Administrador Teste", role: "ADMIN"},
		{email: "medico@clinica.com", password: "secret123", fullName: "Médico Teste", role: "DOCTOR"},
		{email: "enfermeiro@hospital.com", password: "secret123", fullName: "Enfermeiro Teste", role: "NURSE"},
		{email: "recepcao@hospital.com", password: "secret123", fullName: "Recepcionista Teste", role: "RECEPTION"},
		{email: "paciente@mail.com", password: "secret123", fullName: "Paciente Teste", role: "PATIENT"},
	}

	for _, seedUser := range seedUsers {
		passwordHash, hashError := bcrypt.GenerateFromPassword([]byte(seedUser.password), bcrypt.MinCost)
		if hashError != nil {
			t.Fatalf("failed to hash seed user password: %v", hashError)
		}
		_, insertError := pool.Exec(ctx, `
			INSERT INTO users (id, email, password_hash, full_name, role, is_active, created_at, updated_at)
			VALUES (uuid_generate_v4(), $1, $2, $3, $4, true, NOW(), NOW())
			ON CONFLICT (email) DO NOTHING`,
			seedUser.email, string(passwordHash), seedUser.fullName, seedUser.role)
		if insertError != nil {
			t.Fatalf("failed to seed user %s: %v", seedUser.email, insertError)
		}
	}
}

func seedPortalPatientIdentity(t *testing.T, pool *pgxpool.Pool, fhirClient *inMemoryFHIRClient) {
	t.Helper()
	ctx := context.Background()

	var patientUserID uuid.UUID
	if queryError := pool.QueryRow(ctx, `SELECT id FROM users WHERE email = 'paciente@mail.com'`).Scan(&patientUserID); queryError != nil {
		t.Fatalf("failed to fetch patient user: %v", queryError)
	}

	patientResource := map[string]interface{}{
		"id":         patientUserID.String(),
		"name":       []map[string]interface{}{{"given": []string{"Maria"}, "family": "Silva"}},
		"birthDate":  "1990-01-01",
		"gender":     "female",
		"identifier": []map[string]string{{"system": "http://hospital.com.br/patients/cpf", "value": "52998224725"}},
	}
	if _, seedError := fhirClient.CreateResource(ctx, "Patient", patientResource); seedError != nil {
		t.Fatalf("failed to seed portal patient: %v", seedError)
	}

	_, linkError := pool.Exec(ctx, `
		INSERT INTO patient_user_links (user_id, patient_fhir_id)
		VALUES ($1, $2)
		ON CONFLICT (user_id) DO UPDATE SET patient_fhir_id = EXCLUDED.patient_fhir_id`,
		patientUserID, patientUserID.String())
	if linkError != nil {
		t.Fatalf("failed to seed patient user link: %v", linkError)
	}
}

func resetTestData(t *testing.T, pool *pgxpool.Pool) {
	t.Helper()
	ctx := context.Background()

	appTables := []string{
		"employees",
		"appointments",
		"idempotency_keys",
		"notifications",
		"notification_recipients",
		"audit_logs",
		"diagnostic_report_versions",
		"imaging_studies",
		"telemetry_rooms",
		"telemetry_beds",
		"exam_analyses",
		"exam_analyses_audit_logs",
		"patient_user_links",
		"revoked_tokens",
	}

	_, truncateError := pool.Exec(ctx, "TRUNCATE "+strings.Join(appTables, ", ")+" RESTART IDENTITY CASCADE")
	if truncateError != nil {
		t.Fatalf("failed to reset test data: %v", truncateError)
	}
}
