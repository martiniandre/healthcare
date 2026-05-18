package main

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/healthcare/backend/internal/app"
	"github.com/healthcare/backend/internal/modules/auth"
	authpb "github.com/healthcare/backend/internal/modules/auth/pb"
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

	authRepository := auth.NewRepository(databasePool)
	authService := auth.NewService(authRepository)
	authGRPCHandler := auth.NewGRPCHandler(authService)

	authpb.RegisterAuthServiceServer(applicationServer.GRPCServer, authGRPCHandler)
	staff.Register(applicationServer.GRPCServer, databasePool)
	patients.Register(applicationServer.GRPCServer, fhirClient)
	clinical.Register(applicationServer.GRPCServer, fhirClient)
	imaging.Register(applicationServer.GRPCServer, databasePool, redisClient, appConfig.GCSBucketName)
	telemetry.Register(applicationServer.GRPCServer, databasePool)
	health.Register(applicationServer.GRPCServer, databasePool, redisClient)

	imagingWorker := imaging.NewWorker(imaging.NewRepository(databasePool), redisClient, fhirClient)
	go imagingWorker.Start(mainContext)

	go func() {
		seedingContext, seedingCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer seedingCancel()

		doctorEmailAddress := "medico@clinica.com"
		_, _, doctorLoginError := authService.Login(seedingContext, doctorEmailAddress, "senha123")
		if doctorLoginError != nil && errors.Is(doctorLoginError, auth.ErrUserNotFound) {
			_, registerError := authService.Register(seedingContext, doctorEmailAddress, "senha123", "Dr. Guilherme Araujo", "RoleDoctor")
			if registerError != nil {
				slog.Error("Failed to seed doctor user", "error", registerError)
			} else {
				slog.Info("Successfully seeded doctor user", "email", doctorEmailAddress)
			}
		}

		adminEmailAddress := "admin@clinica.com"
		_, _, adminLoginError := authService.Login(seedingContext, adminEmailAddress, "admin123")
		if adminLoginError != nil && errors.Is(adminLoginError, auth.ErrUserNotFound) {
			_, registerError := authService.Register(seedingContext, adminEmailAddress, "admin123", "Administrador Central", "RoleAdmin")
			if registerError != nil {
				slog.Error("Failed to seed admin user", "error", registerError)
			} else {
				slog.Info("Successfully seeded admin user", "email", adminEmailAddress)
			}
		}
	}()

	httpServeMux := http.NewServeMux()

	corsOptionHandler := func(writer http.ResponseWriter, request *http.Request) bool {
		writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		writer.Header().Set("Access-Control-Allow-Credentials", "true")
		writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-CSRF-Token, Authorization")
		if request.Method == http.MethodOptions {
			writer.WriteHeader(http.StatusOK)
			return true
		}
		return false
	}

	httpServeMux.HandleFunc("/api/auth/login", func(writer http.ResponseWriter, request *http.Request) {
		if corsOptionHandler(writer, request) {
			return
		}

		if request.Method != http.MethodPost {
			http.Error(writer, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		var loginRequestPayload struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		if decodeError := json.NewDecoder(request.Body).Decode(&loginRequestPayload); decodeError != nil {
			http.Error(writer, "Bad Request", http.StatusBadRequest)
			return
		}

		authenticatedUser, jsonWebToken, authError := authService.Login(request.Context(), loginRequestPayload.Email, loginRequestPayload.Password)
		if authError != nil {
			slog.Warn("Login failed", "email", loginRequestPayload.Email, "error", authError)
			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(writer).Encode(map[string]string{"error": "Credenciais inválidas."})
			return
		}

		crossSiteRequestForgeryToken := uuid.New().String()

		http.SetCookie(writer, &http.Cookie{
			Name:     "token",
			Value:    jsonWebToken,
			Path:     "/",
			HttpOnly: true,
			Secure:   false,
			SameSite: http.SameSiteLaxMode,
			MaxAge:   86400,
		})

		http.SetCookie(writer, &http.Cookie{
			Name:     "csrf_token",
			Value:    crossSiteRequestForgeryToken,
			Path:     "/",
			HttpOnly: false,
			Secure:   false,
			SameSite: http.SameSiteLaxMode,
			MaxAge:   86400,
		})

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		json.NewEncoder(writer).Encode(map[string]interface{}{
			"token":  jsonWebToken,
			"userId": authenticatedUser.ID.String(),
			"role":   string(authenticatedUser.Role),
			"email":  authenticatedUser.Email,
		})
	})

	httpServeMux.HandleFunc("/api/auth/logout", func(writer http.ResponseWriter, request *http.Request) {
		if corsOptionHandler(writer, request) {
			return
		}

		if request.Method != http.MethodPost {
			http.Error(writer, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		http.SetCookie(writer, &http.Cookie{
			Name:     "token",
			Value:    "",
			Path:     "/",
			HttpOnly: true,
			Secure:   false,
			SameSite: http.SameSiteLaxMode,
			MaxAge:   -1,
		})

		http.SetCookie(writer, &http.Cookie{
			Name:     "csrf_token",
			Value:    "",
			Path:     "/",
			HttpOnly: false,
			Secure:   false,
			SameSite: http.SameSiteLaxMode,
			MaxAge:   -1,
		})

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		json.NewEncoder(writer).Encode(map[string]string{"message": "Logged out successfully"})
	})

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
		Addr:    ":8080",
		Handler: httpServeMux,
	}

	go func() {
		slog.Info("HTTP server listening", "address", ":8080")
		if serveError := httpServer.ListenAndServe(); serveError != nil && !errors.Is(serveError, http.ErrServerClosed) {
			slog.Error("Failed to serve HTTP", "error", serveError)
		}
	}()

	quitSignalChannel := make(chan os.Signal, 1)
	signal.Notify(quitSignalChannel, os.Interrupt, syscall.SIGTERM)
	<-quitSignalChannel

	slog.Info("Shutting down servers gracefully...")
	imagingWorker.Stop()
	applicationServer.GRPCServer.GracefulStop()

	ctxShutdownTimeout, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdown()

	if shutdownError := httpServer.Shutdown(ctxShutdownTimeout); shutdownError != nil {
		slog.Error("Failed to shutdown HTTP server", "error", shutdownError)
	}

	slog.Info("Server stopped")
}
