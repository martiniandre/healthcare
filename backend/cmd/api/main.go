package main

import (
	"context"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/healthcare/backend/internal/app"
	"github.com/healthcare/backend/internal/modules/auth"
	"github.com/healthcare/backend/internal/modules/clinical"
	"github.com/healthcare/backend/internal/modules/health"
	"github.com/healthcare/backend/internal/modules/imaging"
	"github.com/healthcare/backend/internal/modules/patients"
	"github.com/healthcare/backend/internal/modules/staff"
	"github.com/healthcare/backend/internal/modules/telemetry"
	"github.com/healthcare/backend/internal/shared/cache"
	"github.com/healthcare/backend/internal/shared/config"
	"github.com/healthcare/backend/internal/shared/database"
	"github.com/healthcare/backend/internal/shared/healthcare"
	"github.com/healthcare/backend/internal/shared/logger"
	"github.com/healthcare/backend/internal/shared/migrations"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("Invalid configuration", "error", err)
		os.Exit(1)
	}

	logger.Init(cfg.AppEnv, cfg.SentryDSN)
	slog.Info("Starting Healthcare API", "env", cfg.AppEnv, "port", cfg.AppPort)

	if err := auth.InitJWT(cfg.JWTSecret); err != nil {
		slog.Error("Failed to initialize JWT", "error", err)
		os.Exit(1)
	}

	if err := migrations.Run(cfg.DBUrl); err != nil {
		slog.Error("Failed to run database migrations", "error", err)
		os.Exit(1)
	}

	ctx := context.Background()

	dbPool, err := database.Connect(ctx, cfg.DBUrl)
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer dbPool.Close()

	fhirClient, err := healthcare.NewClient(ctx, cfg.GCPProjectID, cfg.GCPLocationID, cfg.GCPDatasetID, cfg.GCPFHIRStore)
	if err != nil {
		slog.Error("Failed to initialize Healthcare API client", "error", err)
		os.Exit(1)
	}

	redisClient := cache.Connect(cfg.RedisUrl)

	server := app.NewServer(redisClient)

	auth.Register(server.GRPCServer, dbPool)
	staff.Register(server.GRPCServer, dbPool)
	patients.Register(server.GRPCServer, fhirClient)
	clinical.Register(server.GRPCServer, fhirClient)
	imaging.Register(server.GRPCServer, dbPool, redisClient, cfg.GCSBucketName)
	telemetry.Register(server.GRPCServer, dbPool)
	health.Register(server.GRPCServer, dbPool, redisClient)

	imagingWorker := imaging.NewWorker(imaging.NewRepository(dbPool), redisClient, fhirClient)
	go imagingWorker.Start(ctx)

	listener, err := net.Listen("tcp", ":"+cfg.AppPort)
	if err != nil {
		slog.Error("Failed to listen", "error", err)
		os.Exit(1)
	}

	go func() {
		slog.Info("gRPC server listening", "address", listener.Addr().String())
		if err := server.GRPCServer.Serve(listener); err != nil {
			slog.Error("Failed to serve", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	slog.Info("Shutting down gRPC server gracefully...")
	imagingWorker.Stop()
	server.GRPCServer.GracefulStop()
	time.Sleep(1 * time.Second)
	slog.Info("Server stopped")
}
