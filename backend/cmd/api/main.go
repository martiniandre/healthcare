package main

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/healthcare/backend/internal/api"
	"github.com/healthcare/backend/internal/app"
	"github.com/healthcare/backend/internal/modules/audit_logs"
	"github.com/healthcare/backend/internal/modules/auth"
	"github.com/healthcare/backend/internal/modules/clinical"
	"github.com/healthcare/backend/internal/modules/exam_analyzer"
	"github.com/healthcare/backend/internal/modules/health"
	"github.com/healthcare/backend/internal/modules/imaging"
	"github.com/healthcare/backend/internal/modules/patients"
	"github.com/healthcare/backend/internal/modules/staff"
	"github.com/healthcare/backend/internal/modules/stats"
	"github.com/healthcare/backend/internal/modules/telemetry"
	"github.com/healthcare/backend/internal/shared/cache"
	"github.com/healthcare/backend/internal/shared/config"
	"github.com/healthcare/backend/internal/shared/database"
	"github.com/healthcare/backend/internal/shared/healthcare"
	"github.com/healthcare/backend/internal/shared/logger"
	"github.com/healthcare/backend/internal/shared/migrations"
	"github.com/healthcare/backend/internal/shared/storage"
)

func main() {
	appConfig, loadError := config.Load()
	if loadError != nil {
		slog.Error("Invalid configuration", "error", loadError)
		os.Exit(1)
	}

	logger.Init(appConfig.AppEnv, appConfig.SentryDSN)
	slog.Info("Starting Healthcare API", "env", appConfig.AppEnv, "port", appConfig.AppPort)

	if initJWTError := auth.InitJWT(appConfig.JWTSecret); initJWTError != nil {
		slog.Error("Failed to initialize JWT", "error", initJWTError)
		os.Exit(1)
	}

	if migrationError := migrations.Run(appConfig.DBUrl); migrationError != nil {
		slog.Error("Failed to run database migrations", "error", migrationError)
		os.Exit(1)
	}

	mainContext := context.Background()

	databasePool, connectionError := database.Connect(mainContext, appConfig.DBUrl)
	if connectionError != nil {
		slog.Error("Failed to connect to database", "error", connectionError)
		os.Exit(1)
	}
	defer databasePool.Close()

	fhirClient, fhirClientError := healthcare.NewClient(mainContext, appConfig.GCPProjectID, appConfig.GCPLocationID, appConfig.GCPDatasetID, appConfig.GCPFHIRStore)
	if fhirClientError != nil {
		slog.Error("Failed to initialize Healthcare API client", "error", fhirClientError)
		os.Exit(1)
	}

	redisClient := cache.Connect(appConfig.RedisUrl)

	applicationServer := app.NewServer(redisClient)

	authService := auth.Register(applicationServer.GRPCServer, auth.Dependency{DB: databasePool})
	staffService := staff.Register(applicationServer.GRPCServer, staff.Dependency{DB: databasePool})
	patientsService := patients.Register(applicationServer.GRPCServer, patients.Dependency{FHIRClient: fhirClient})
	clinicalService := clinical.Register(applicationServer.GRPCServer, clinical.Dependency{FHIRClient: fhirClient})
	storageClient, storageClientErr := storage.NewGCSClient(mainContext)
	if storageClientErr != nil {
		slog.Warn("Failed to initialize GCS client, falling back to dummy", "error", storageClientErr)
		storageClient = storage.NewStorageClient()
	}
	imagingService := imaging.Register(applicationServer.GRPCServer, imaging.Dependency{DB: databasePool, Storage: storageClient, Redis: redisClient, BucketName: appConfig.GCSBucketName})
	telemetryService := telemetry.Register(applicationServer.GRPCServer, telemetry.Dependency{DB: databasePool})
	telemetrySimulator := telemetry.StartSimulator(mainContext, databasePool)
	health.Register(applicationServer.GRPCServer, health.Dependency{DB: databasePool, Redis: redisClient})
	statsHTTPHandler := stats.Register(stats.Dependency{DB: databasePool, FHIRClient: fhirClient})
	auditLogsService := audit_logs.Register(applicationServer.GRPCServer, audit_logs.Dependency{DB: databasePool})

	exam_analyzer.Register(exam_analyzer.Dependency{DB: databasePool, ProjectID: appConfig.GCPProjectID, LocationID: appConfig.GCPLocationID, VertexModel: appConfig.GCPVertexModel})
	go exam_analyzer.WorkerInstance.Start(mainContext)

	imagingWorker := imaging.NewWorker(imaging.NewRepository(databasePool), redisClient, fhirClient)
	go imagingWorker.Start(mainContext)

	secureCookies := appConfig.AppEnv != "development" && appConfig.AppEnv != "test"

	authHTTPHandler := auth.NewHTTPHandler(authService, secureCookies)
	patientsHTTPHandler := patients.NewHTTPHandler(patientsService)
	clinicalHTTPHandler := clinical.NewHTTPHandler(clinicalService)
	imagingHTTPHandler := imaging.NewHTTPHandler(imagingService)
	staffHTTPHandler := staff.NewHTTPHandler(staffService)
	telemetryHTTPHandler := telemetry.NewHTTPHandler(telemetryService)
	examAnalyzerHTTPHandler := exam_analyzer.NewHTTPHandler(exam_analyzer.Repo, exam_analyzer.Svc, exam_analyzer.WorkerInstance)
	auditLogsHTTPHandler := audit_logs.NewHTTPHandler(auditLogsService)

	router := api.NewRouter(
		secureCookies,
		authHTTPHandler,
		patientsHTTPHandler,
		clinicalHTTPHandler,
		imagingHTTPHandler,
		staffHTTPHandler,
		telemetryHTTPHandler,
		examAnalyzerHTTPHandler,
		statsHTTPHandler,
		auditLogsHTTPHandler,
	)

	tcpListener, listenerError := net.Listen("tcp", ":"+appConfig.AppPort)
	if listenerError != nil {
		slog.Error("Failed to listen", "error", listenerError)
		os.Exit(1)
	}

	go func() {
		slog.Info("gRPC server listening", "address", tcpListener.Addr().String())
		if serveError := applicationServer.GRPCServer.Serve(tcpListener); serveError != nil {
			slog.Error("Failed to serve gRPC", "error", serveError)
		}
	}()

	httpServer := &http.Server{
		Addr:              ":" + appConfig.HTTPPort,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       120 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}

	go func() {
		slog.Info("HTTP server listening", "address", ":"+appConfig.HTTPPort)
		if serveError := httpServer.ListenAndServe(); serveError != nil && !errors.Is(serveError, http.ErrServerClosed) {
			slog.Error("Failed to serve HTTP", "error", serveError)
		}
	}()

	quitSignalChannel := make(chan os.Signal, 1)
	signal.Notify(quitSignalChannel, os.Interrupt, syscall.SIGTERM)
	<-quitSignalChannel

	slog.Info("Shutting down servers gracefully...")
	exam_analyzer.WorkerInstance.Stop()
	imagingWorker.Stop()
	telemetrySimulator.Stop()
	applicationServer.GRPCServer.GracefulStop()

	ctxShutdownTimeout, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdown()

	if shutdownError := httpServer.Shutdown(ctxShutdownTimeout); shutdownError != nil {
		slog.Error("Failed to shutdown HTTP server", "error", shutdownError)
	}

	slog.Info("Server stopped")
}
