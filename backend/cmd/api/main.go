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
	"strings"
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

	patientsRepository := patients.NewRepository(fhirClient)
	patientsService := patients.NewService(patientsRepository)

	validateHTTPAuth := func(httpResponseWriter http.ResponseWriter, httpRequest *http.Request, allowedRoles []auth.Role) (context.Context, bool) {
		cookie, cookieError := httpRequest.Cookie("token")
		if cookieError != nil {
			httpResponseWriter.Header().Set("Content-Type", "application/json")
			httpResponseWriter.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(httpResponseWriter).Encode(map[string]string{"error": "Não autenticado."})
			return nil, false
		}

		claims, jwtValidationErr := auth.ValidateJWT(cookie.Value)
		if jwtValidationErr != nil {
			httpResponseWriter.Header().Set("Content-Type", "application/json")
			httpResponseWriter.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(httpResponseWriter).Encode(map[string]string{"error": "Sessão expirada."})
			return nil, false
		}

		roleStr, roleClaimExists := claims["role"].(string)
		if !roleClaimExists {
			httpResponseWriter.Header().Set("Content-Type", "application/json")
			httpResponseWriter.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(httpResponseWriter).Encode(map[string]string{"error": "Função de usuário inválida."})
			return nil, false
		}

		callerRole := auth.Role(roleStr)
		roleAllowed := false
		for _, allowedRole := range allowedRoles {
			if callerRole == allowedRole {
				roleAllowed = true
				break
			}
		}

		if !roleAllowed {
			httpResponseWriter.Header().Set("Content-Type", "application/json")
			httpResponseWriter.WriteHeader(http.StatusForbidden)
			json.NewEncoder(httpResponseWriter).Encode(map[string]string{"error": "Acesso negado."})
			return nil, false
		}

		if httpRequest.Method == http.MethodPost || httpRequest.Method == http.MethodPut || httpRequest.Method == http.MethodDelete {
			csrfHeader := httpRequest.Header.Get("X-CSRF-Token")
			csrfCookie, csrfCookieErr := httpRequest.Cookie("csrf_token")
			if csrfCookieErr != nil || csrfHeader == "" || csrfHeader != csrfCookie.Value {
				httpResponseWriter.Header().Set("Content-Type", "application/json")
				httpResponseWriter.WriteHeader(http.StatusForbidden)
				json.NewEncoder(httpResponseWriter).Encode(map[string]string{"error": "Token CSRF inválido ou ausente."})
				return nil, false
			}
		}

		userIDStr, _ := claims["user_id"].(string)
		contextWithValues := context.WithValue(httpRequest.Context(), "user_id", userIDStr)
		contextWithValues = context.WithValue(contextWithValues, "role", roleStr)
		return contextWithValues, true
	}

	httpServeMux.HandleFunc("/api/patients", func(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
		if corsOptionHandler(httpResponseWriter, httpRequest) {
			return
		}

		if httpRequest.Method == http.MethodGet {
			contextWithValues, authIsOk := validateHTTPAuth(httpResponseWriter, httpRequest, []auth.Role{auth.RoleAdmin, auth.RoleDoctor, auth.RoleNurse, auth.RoleReception})
			if !authIsOk {
				return
			}

			patientsList, patientListErr := patientsService.ListPatients(contextWithValues)
			if patientListErr != nil {
				httpResponseWriter.Header().Set("Content-Type", "application/json")
				httpResponseWriter.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(httpResponseWriter).Encode(map[string]string{"error": "Erro ao listar pacientes."})
				return
			}

			type patientResponse struct {
				PatientID      string `json:"patient_id"`
				FHIRResourceID string `json:"fhir_resource_id"`
				FullName       string `json:"full_name"`
				BirthDate      string `json:"birth_date"`
				DocumentID     string `json:"document_id"`
				PhoneNumber    string `json:"phone_number"`
			}

			responseList := make([]patientResponse, 0, len(patientsList))
			for _, patient := range patientsList {
				responseList = append(responseList, patientResponse{
					PatientID:      patient.ID.String(),
					FHIRResourceID: patient.FHIRResourceID,
					FullName:       patient.FullName,
					BirthDate:      patient.BirthDate.Format("2006-01-02"),
					DocumentID:     patient.DocumentID,
					PhoneNumber:    patient.PhoneNumber,
				})
			}

			httpResponseWriter.Header().Set("Content-Type", "application/json")
			httpResponseWriter.WriteHeader(http.StatusOK)
			json.NewEncoder(httpResponseWriter).Encode(responseList)
			return
		}

		if httpRequest.Method == http.MethodPost {
			contextWithValues, authIsOk := validateHTTPAuth(httpResponseWriter, httpRequest, []auth.Role{auth.RoleAdmin, auth.RoleReception})
			if !authIsOk {
				return
			}

			var payload struct {
				FullName    string `json:"full_name"`
				BirthDate   string `json:"birth_date"`
				DocumentID  string `json:"document_id"`
				PhoneNumber string `json:"phone_number"`
			}

			if payloadDecodeErr := json.NewDecoder(httpRequest.Body).Decode(&payload); payloadDecodeErr != nil {
				httpResponseWriter.Header().Set("Content-Type", "application/json")
				httpResponseWriter.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(httpResponseWriter).Encode(map[string]string{"error": "Payload inválido."})
				return
			}

			patient, createPatientErr := patientsService.CreatePatient(contextWithValues, payload.FullName, payload.BirthDate, payload.DocumentID, payload.PhoneNumber)
			if createPatientErr != nil {
				httpResponseWriter.Header().Set("Content-Type", "application/json")
				if errors.Is(createPatientErr, patients.ErrPatientAlreadyExists) {
					httpResponseWriter.WriteHeader(http.StatusConflict)
					json.NewEncoder(httpResponseWriter).Encode(map[string]string{"error": "Paciente com este documento já cadastrado."})
					return
				}
				httpResponseWriter.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(httpResponseWriter).Encode(map[string]string{"error": "Erro ao criar paciente."})
				return
			}

			httpResponseWriter.Header().Set("Content-Type", "application/json")
			httpResponseWriter.WriteHeader(http.StatusCreated)
			json.NewEncoder(httpResponseWriter).Encode(map[string]string{
				"patient_id":       patient.ID.String(),
				"fhir_resource_id": patient.FHIRResourceID,
			})
			return
		}

		http.Error(httpResponseWriter, "Method Not Allowed", http.StatusMethodNotAllowed)
	})

	httpServeMux.HandleFunc("/api/patients/", func(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
		if corsOptionHandler(httpResponseWriter, httpRequest) {
			return
		}

		if httpRequest.Method != http.MethodGet {
			http.Error(httpResponseWriter, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		contextWithValues, authIsOk := validateHTTPAuth(httpResponseWriter, httpRequest, []auth.Role{auth.RoleAdmin, auth.RoleDoctor, auth.RoleNurse, auth.RoleReception})
		if !authIsOk {
			return
		}

		fhirResourceID := strings.TrimPrefix(httpRequest.URL.Path, "/api/patients/")
		if fhirResourceID == "" {
			httpResponseWriter.Header().Set("Content-Type", "application/json")
			httpResponseWriter.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(httpResponseWriter).Encode(map[string]string{"error": "ID do paciente ausente."})
			return
		}

		patient, getPatientErr := patientsService.GetPatient(contextWithValues, fhirResourceID)
		if getPatientErr != nil {
			httpResponseWriter.Header().Set("Content-Type", "application/json")
			httpResponseWriter.WriteHeader(http.StatusNotFound)
			json.NewEncoder(httpResponseWriter).Encode(map[string]string{"error": "Paciente não encontrado."})
			return
		}

		httpResponseWriter.Header().Set("Content-Type", "application/json")
		httpResponseWriter.WriteHeader(http.StatusOK)
		json.NewEncoder(httpResponseWriter).Encode(map[string]interface{}{
			"patient_id":       patient.ID.String(),
			"fhir_resource_id": patient.FHIRResourceID,
			"full_name":        patient.FullName,
			"birth_date":       patient.BirthDate.Format("2006-01-02"),
			"document_id":      patient.DocumentID,
			"phone_number":     patient.PhoneNumber,
		})
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
