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
	"github.com/healthcare/backend/internal/modules/clinical"
	"github.com/healthcare/backend/internal/modules/health"
	"github.com/healthcare/backend/internal/modules/imaging"
	"github.com/healthcare/backend/internal/modules/patients"
	"github.com/healthcare/backend/internal/modules/staff"
	"github.com/healthcare/backend/internal/modules/telemetry"
	"github.com/healthcare/backend/internal/shared/cache"
	"github.com/healthcare/backend/internal/shared/config"
	"github.com/healthcare/backend/internal/shared/ctxkeys"
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

	authService := auth.Register(applicationServer.GRPCServer, databasePool)
	staff.Register(applicationServer.GRPCServer, databasePool)
	patientsService := patients.Register(applicationServer.GRPCServer, fhirClient)
	clinicalService := clinical.Register(applicationServer.GRPCServer, fhirClient)
	storageClient, err := storage.NewGCSClient(mainContext)
	if err != nil {
		slog.Warn("Failed to initialize GCS client, falling back to dummy", "error", err)
		storageClient = storage.NewStorageClient()
	}
	imagingService := imaging.Register(applicationServer.GRPCServer, databasePool, storageClient, redisClient, appConfig.GCSBucketName)
	telemetry.Register(applicationServer.GRPCServer, databasePool)
	health.Register(applicationServer.GRPCServer, databasePool, redisClient)

	imagingWorker := imaging.NewWorker(imaging.NewRepository(databasePool), redisClient, fhirClient)
	go imagingWorker.Start(mainContext)


	httpServeMux := http.NewServeMux()

	corsOptionHandler := func(writer http.ResponseWriter, request *http.Request) bool {
		writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		writer.Header().Set("Access-Control-Allow-Credentials", "true")
		writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-CSRF-Token, Authorization")
		writer.Header().Set("Vary", "Origin")
		writer.Header().Set("X-Content-Type-Options", "nosniff")
		writer.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		writer.Header().Set("Content-Security-Policy", "default-src 'none'; frame-ancestors 'none'")
		if request.Method == http.MethodOptions {
			writer.WriteHeader(http.StatusOK)
			return true
		}
		return false
	}

	secureCookies := appConfig.AppEnv != "development" && appConfig.AppEnv != "test"

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
			Secure:   secureCookies,
			SameSite: http.SameSiteLaxMode,
			MaxAge:   86400,
		})

		http.SetCookie(writer, &http.Cookie{
			Name:     "csrf_token",
			Value:    crossSiteRequestForgeryToken,
			Path:     "/",
			HttpOnly: false,
			Secure:   secureCookies,
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
			Secure:   secureCookies,
			SameSite: http.SameSiteLaxMode,
			MaxAge:   -1,
		})

		http.SetCookie(writer, &http.Cookie{
			Name:     "csrf_token",
			Value:    "",
			Path:     "/",
			HttpOnly: false,
			Secure:   secureCookies,
			SameSite: http.SameSiteLaxMode,
			MaxAge:   -1,
		})

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		json.NewEncoder(writer).Encode(map[string]string{"message": "Logged out successfully"})
	})

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
		contextWithValues := context.WithValue(httpRequest.Context(), ctxkeys.UserIDKey, userIDStr)
		contextWithValues = context.WithValue(contextWithValues, ctxkeys.RoleKey, roleStr)
		return contextWithValues, true
	}

	imagingHTTPHandler := imaging.NewHTTPHandler(imagingService, validateHTTPAuth)

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

		remainingPath := strings.TrimPrefix(httpRequest.URL.Path, "/api/patients/")
		if remainingPath == "" {
			httpResponseWriter.Header().Set("Content-Type", "application/json")
			httpResponseWriter.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(httpResponseWriter).Encode(map[string]string{"error": "ID do paciente ausente."})
			return
		}

		parts := strings.Split(remainingPath, "/")

		if len(parts) == 1 {
			if httpRequest.Method != http.MethodGet {
				http.Error(httpResponseWriter, "Method Not Allowed", http.StatusMethodNotAllowed)
				return
			}

			contextWithValues, authIsOk := validateHTTPAuth(httpResponseWriter, httpRequest, []auth.Role{auth.RoleAdmin, auth.RoleDoctor, auth.RoleNurse, auth.RoleReception})
			if !authIsOk {
				return
			}

			patient, getPatientErr := patientsService.GetPatient(contextWithValues, parts[0])
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
			return
		}

		if len(parts) == 2 {
			patientFHIRID := parts[0]
			subResource := parts[1]

			if subResource == "encounters" {
				if httpRequest.Method == http.MethodGet {
					contextWithValues, authIsOk := validateHTTPAuth(httpResponseWriter, httpRequest, []auth.Role{auth.RoleAdmin, auth.RoleDoctor, auth.RoleNurse, auth.RoleReception})
					if !authIsOk {
						return
					}

					encountersList, encountersErr := clinicalService.GetEncountersByPatient(contextWithValues, patientFHIRID)
					if encountersErr != nil {
						httpResponseWriter.Header().Set("Content-Type", "application/json")
						httpResponseWriter.WriteHeader(http.StatusInternalServerError)
						json.NewEncoder(httpResponseWriter).Encode(map[string]string{"error": "Erro ao carregar consultas do paciente."})
						return
					}

					type encounterResponse struct {
						FHIRID         string `json:"fhir_id"`
						PatientFHIRID  string `json:"patient_fhir_id"`
						Status         string `json:"status"`
						ReasonDisplay  string `json:"reason_display"`
						PractitionerID string `json:"practitioner_id,omitempty"`
						CreatedAt      string `json:"created_at"`
					}

					responseList := make([]encounterResponse, 0, len(encountersList))
					for _, encounter := range encountersList {
						responseList = append(responseList, encounterResponse{
							FHIRID:         encounter.FHIRResourceID,
							PatientFHIRID:  encounter.PatientFHIRID,
							Status:         encounter.Status,
							ReasonDisplay:  encounter.ReasonDisplay,
							PractitionerID: encounter.PractitionerID,
							CreatedAt:      encounter.StartedAt.Format(time.RFC3339),
						})
					}

					httpResponseWriter.Header().Set("Content-Type", "application/json")
					httpResponseWriter.WriteHeader(http.StatusOK)
					json.NewEncoder(httpResponseWriter).Encode(responseList)
					return
				}

				if httpRequest.Method == http.MethodPost {
					contextWithValues, authIsOk := validateHTTPAuth(httpResponseWriter, httpRequest, []auth.Role{auth.RoleAdmin, auth.RoleDoctor, auth.RoleNurse, auth.RoleReception})
					if !authIsOk {
						return
					}

					var payload struct {
						ReasonDisplay  string `json:"reason_display"`
						PractitionerID string `json:"practitioner_id"`
					}

					if payloadDecodeErr := json.NewDecoder(httpRequest.Body).Decode(&payload); payloadDecodeErr != nil {
						httpResponseWriter.Header().Set("Content-Type", "application/json")
						httpResponseWriter.WriteHeader(http.StatusBadRequest)
						json.NewEncoder(httpResponseWriter).Encode(map[string]string{"error": "Payload inválido."})
						return
					}

					newEncounter := &clinical.Encounter{
						PatientFHIRID:  patientFHIRID,
						PractitionerID: payload.PractitionerID,
						ReasonDisplay:  payload.ReasonDisplay,
						Status:         "finished",
						StartedAt:      time.Now(),
					}

					createdEncounter, createErr := clinicalService.CreateEncounter(contextWithValues, newEncounter)
					if createErr != nil {
						httpResponseWriter.Header().Set("Content-Type", "application/json")
						httpResponseWriter.WriteHeader(http.StatusInternalServerError)
						json.NewEncoder(httpResponseWriter).Encode(map[string]string{"error": "Erro ao criar consulta."})
						return
					}

					httpResponseWriter.Header().Set("Content-Type", "application/json")
					httpResponseWriter.WriteHeader(http.StatusCreated)
					json.NewEncoder(httpResponseWriter).Encode(map[string]interface{}{
						"fhir_id":         createdEncounter.FHIRResourceID,
						"patient_fhir_id": createdEncounter.PatientFHIRID,
						"status":          createdEncounter.Status,
						"reason_display":  createdEncounter.ReasonDisplay,
						"practitioner_id": createdEncounter.PractitionerID,
						"created_at":      createdEncounter.StartedAt.Format(time.RFC3339),
					})
					return
				}

				http.Error(httpResponseWriter, "Method Not Allowed", http.StatusMethodNotAllowed)
				return
			}

			if subResource == "observations" {
				if httpRequest.Method != http.MethodGet {
					http.Error(httpResponseWriter, "Method Not Allowed", http.StatusMethodNotAllowed)
					return
				}

				contextWithValues, authIsOk := validateHTTPAuth(httpResponseWriter, httpRequest, []auth.Role{auth.RoleAdmin, auth.RoleDoctor, auth.RoleNurse})
				if !authIsOk {
					return
				}

				observationsList, observationsErr := clinicalService.GetObservationsByPatient(contextWithValues, patientFHIRID)
				if observationsErr != nil {
					httpResponseWriter.Header().Set("Content-Type", "application/json")
					httpResponseWriter.WriteHeader(http.StatusInternalServerError)
					json.NewEncoder(httpResponseWriter).Encode(map[string]string{"error": "Erro ao carregar observações do paciente."})
					return
				}

				type observationResponse struct {
					FHIRID          string  `json:"fhir_id"`
					EncounterFHIRID string  `json:"encounter_fhir_id"`
					PatientFHIRID   string  `json:"patient_fhir_id"`
					LoincCode       string  `json:"loinc_code"`
					CodeDisplay     string  `json:"code_display"`
					ValueQuantity   float64 `json:"value_quantity"`
					ValueUnit       string  `json:"value_unit"`
					CreatedAt       string  `json:"created_at"`
				}

				responseList := make([]observationResponse, 0, len(observationsList))
				for _, observation := range observationsList {
					responseList = append(responseList, observationResponse{
						FHIRID:          observation.FHIRResourceID,
						EncounterFHIRID: observation.EncounterFHIRID,
						PatientFHIRID:   observation.PatientFHIRID,
						LoincCode:       observation.LoincCode,
						CodeDisplay:     observation.CodeDisplay,
						ValueQuantity:   observation.ValueQuantity,
						ValueUnit:       observation.ValueUnit,
						CreatedAt:       observation.ObservedAt.Format(time.RFC3339),
					})
				}

				httpResponseWriter.Header().Set("Content-Type", "application/json")
				httpResponseWriter.WriteHeader(http.StatusOK)
				json.NewEncoder(httpResponseWriter).Encode(responseList)
				return
			}

			if subResource == "conditions" {
				if httpRequest.Method == http.MethodGet {
					contextWithValues, authIsOk := validateHTTPAuth(httpResponseWriter, httpRequest, []auth.Role{auth.RoleAdmin, auth.RoleDoctor, auth.RoleNurse})
					if !authIsOk {
						return
					}

					conditionsList, conditionsErr := clinicalService.GetConditionsByPatient(contextWithValues, patientFHIRID)
					if conditionsErr != nil {
						httpResponseWriter.Header().Set("Content-Type", "application/json")
						httpResponseWriter.WriteHeader(http.StatusInternalServerError)
						json.NewEncoder(httpResponseWriter).Encode(map[string]string{"error": "Erro ao carregar diagnósticos do paciente."})
						return
					}

					type conditionResponse struct {
						FHIRID         string `json:"fhir_id"`
						PatientFHIRID  string `json:"patient_fhir_id"`
						ICD10Code      string `json:"icd10_code"`
						CodeDisplay    string `json:"code_display"`
						ClinicalStatus string `json:"clinical_status"`
						CreatedAt      string `json:"created_at"`
					}

					responseList := make([]conditionResponse, 0, len(conditionsList))
					for _, condition := range conditionsList {
						responseList = append(responseList, conditionResponse{
							FHIRID:         condition.FHIRResourceID,
							PatientFHIRID:  condition.PatientFHIRID,
							ICD10Code:      condition.ICD10Code,
							CodeDisplay:    condition.CodeDisplay,
							ClinicalStatus: condition.ClinicalStatus,
							CreatedAt:      condition.OnsetAt.Format(time.RFC3339),
						})
					}

					httpResponseWriter.Header().Set("Content-Type", "application/json")
					httpResponseWriter.WriteHeader(http.StatusOK)
					json.NewEncoder(httpResponseWriter).Encode(responseList)
					return
				}

				if httpRequest.Method == http.MethodPost {
					contextWithValues, authIsOk := validateHTTPAuth(httpResponseWriter, httpRequest, []auth.Role{auth.RoleDoctor, auth.RoleNurse})
					if !authIsOk {
						return
					}

					var payload struct {
						ICD10Code   string `json:"icd10_code"`
						CodeDisplay string `json:"code_display"`
					}

					if payloadDecodeErr := json.NewDecoder(httpRequest.Body).Decode(&payload); payloadDecodeErr != nil {
						httpResponseWriter.Header().Set("Content-Type", "application/json")
						httpResponseWriter.WriteHeader(http.StatusBadRequest)
						json.NewEncoder(httpResponseWriter).Encode(map[string]string{"error": "Payload inválido."})
						return
					}

					newCondition := &clinical.Condition{
						PatientFHIRID:  patientFHIRID,
						ICD10Code:      payload.ICD10Code,
						CodeDisplay:    payload.CodeDisplay,
						ClinicalStatus: "active",
						OnsetAt:        time.Now(),
					}

					createdCondition, createErr := clinicalService.CreateCondition(contextWithValues, newCondition)
					if createErr != nil {
						httpResponseWriter.Header().Set("Content-Type", "application/json")
						httpResponseWriter.WriteHeader(http.StatusInternalServerError)
						json.NewEncoder(httpResponseWriter).Encode(map[string]string{"error": "Erro ao criar diagnóstico."})
						return
					}

					httpResponseWriter.Header().Set("Content-Type", "application/json")
					httpResponseWriter.WriteHeader(http.StatusCreated)
					json.NewEncoder(httpResponseWriter).Encode(map[string]interface{}{
						"fhir_id":         createdCondition.FHIRResourceID,
						"patient_fhir_id": createdCondition.PatientFHIRID,
						"icd10_code":      createdCondition.ICD10Code,
						"code_display":    createdCondition.CodeDisplay,
						"clinical_status": createdCondition.ClinicalStatus,
						"created_at":      createdCondition.OnsetAt.Format(time.RFC3339),
					})
					return
				}

				http.Error(httpResponseWriter, "Method Not Allowed", http.StatusMethodNotAllowed)
				return
			}

			if subResource == "studies" {
				imagingHTTPHandler.HandlePatientStudies(httpResponseWriter, httpRequest, patientFHIRID)
				return
			}
		}

		http.Error(httpResponseWriter, "Not Found", http.StatusNotFound)
	})

	httpServeMux.HandleFunc("/api/encounters/", func(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
		if corsOptionHandler(httpResponseWriter, httpRequest) {
			return
		}

		remainingPath := strings.TrimPrefix(httpRequest.URL.Path, "/api/encounters/")
		if remainingPath == "" {
			httpResponseWriter.Header().Set("Content-Type", "application/json")
			httpResponseWriter.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(httpResponseWriter).Encode(map[string]string{"error": "ID da consulta ausente."})
			return
		}

		parts := strings.Split(remainingPath, "/")
		if len(parts) == 2 {
			encounterFHIRID := parts[0]
			subResource := parts[1]

			if subResource == "observations" {
				if httpRequest.Method == http.MethodGet {
					contextWithValues, authIsOk := validateHTTPAuth(httpResponseWriter, httpRequest, []auth.Role{auth.RoleAdmin, auth.RoleDoctor, auth.RoleNurse})
					if !authIsOk {
						return
					}

					observationsList, observationsErr := clinicalService.GetObservationsByEncounter(contextWithValues, encounterFHIRID)
					if observationsErr != nil {
						httpResponseWriter.Header().Set("Content-Type", "application/json")
						httpResponseWriter.WriteHeader(http.StatusInternalServerError)
						json.NewEncoder(httpResponseWriter).Encode(map[string]string{"error": "Erro ao carregar observações da consulta."})
						return
					}

					type observationResponse struct {
						FHIRID          string  `json:"fhir_id"`
						EncounterFHIRID string  `json:"encounter_fhir_id"`
						PatientFHIRID   string  `json:"patient_fhir_id"`
						LoincCode       string  `json:"loinc_code"`
						CodeDisplay     string  `json:"code_display"`
						ValueQuantity   float64 `json:"value_quantity"`
						ValueUnit       string  `json:"value_unit"`
						CreatedAt       string  `json:"created_at"`
					}

					responseList := make([]observationResponse, 0, len(observationsList))
					for _, observation := range observationsList {
						responseList = append(responseList, observationResponse{
							FHIRID:          observation.FHIRResourceID,
							EncounterFHIRID: observation.EncounterFHIRID,
							PatientFHIRID:   observation.PatientFHIRID,
							LoincCode:       observation.LoincCode,
							CodeDisplay:     observation.CodeDisplay,
							ValueQuantity:   observation.ValueQuantity,
							ValueUnit:       observation.ValueUnit,
							CreatedAt:       observation.ObservedAt.Format(time.RFC3339),
						})
					}

					httpResponseWriter.Header().Set("Content-Type", "application/json")
					httpResponseWriter.WriteHeader(http.StatusOK)
					json.NewEncoder(httpResponseWriter).Encode(responseList)
					return
				}

				if httpRequest.Method == http.MethodPost {
					contextWithValues, authIsOk := validateHTTPAuth(httpResponseWriter, httpRequest, []auth.Role{auth.RoleDoctor, auth.RoleNurse})
					if !authIsOk {
						return
					}

					var payload struct {
						PatientFHIRID string  `json:"patient_fhir_id"`
						LoincCode     string  `json:"loinc_code"`
						CodeDisplay   string  `json:"code_display"`
						ValueQuantity float64 `json:"value_quantity"`
						ValueUnit     string  `json:"value_unit"`
					}

					if payloadDecodeErr := json.NewDecoder(httpRequest.Body).Decode(&payload); payloadDecodeErr != nil {
						httpResponseWriter.Header().Set("Content-Type", "application/json")
						httpResponseWriter.WriteHeader(http.StatusBadRequest)
						json.NewEncoder(httpResponseWriter).Encode(map[string]string{"error": "Payload inválido."})
						return
					}

					newObservation := &clinical.Observation{
						EncounterFHIRID: encounterFHIRID,
						PatientFHIRID:   payload.PatientFHIRID,
						LoincCode:       payload.LoincCode,
						CodeDisplay:     payload.CodeDisplay,
						ValueQuantity:   payload.ValueQuantity,
						ValueUnit:       payload.ValueUnit,
						ObservedAt:      time.Now(),
					}

					createdObservation, createErr := clinicalService.CreateObservation(contextWithValues, newObservation)
					if createErr != nil {
						httpResponseWriter.Header().Set("Content-Type", "application/json")
						httpResponseWriter.WriteHeader(http.StatusInternalServerError)
						json.NewEncoder(httpResponseWriter).Encode(map[string]string{"error": "Erro ao criar observação."})
						return
					}

					httpResponseWriter.Header().Set("Content-Type", "application/json")
					httpResponseWriter.WriteHeader(http.StatusCreated)
					json.NewEncoder(httpResponseWriter).Encode(map[string]interface{}{
						"fhir_id":           createdObservation.FHIRResourceID,
						"encounter_fhir_id": createdObservation.EncounterFHIRID,
						"patient_fhir_id":   createdObservation.PatientFHIRID,
						"loinc_code":        createdObservation.LoincCode,
						"code_display":      createdObservation.CodeDisplay,
						"value_quantity":    createdObservation.ValueQuantity,
						"value_unit":        createdObservation.ValueUnit,
						"created_at":        createdObservation.ObservedAt.Format(time.RFC3339),
					})
					return
				}

				http.Error(httpResponseWriter, "Method Not Allowed", http.StatusMethodNotAllowed)
				return
			}

			if subResource == "reports" {
				if httpRequest.Method == http.MethodGet {
					contextWithValues, authIsOk := validateHTTPAuth(httpResponseWriter, httpRequest, []auth.Role{auth.RoleAdmin, auth.RoleDoctor, auth.RoleNurse})
					if !authIsOk {
						return
					}

					reportsList, reportsErr := clinicalService.GetDiagnosticReportsByEncounter(contextWithValues, encounterFHIRID)
					if reportsErr != nil {
						httpResponseWriter.Header().Set("Content-Type", "application/json")
						httpResponseWriter.WriteHeader(http.StatusInternalServerError)
						json.NewEncoder(httpResponseWriter).Encode(map[string]string{"error": "Erro ao carregar laudos da consulta."})
						return
					}

					type reportResponse struct {
						FHIRID          string `json:"fhir_id"`
						EncounterFHIRID string `json:"encounter_fhir_id"`
						PatientFHIRID   string `json:"patient_fhir_id"`
						ReportDisplay   string `json:"report_display"`
						Status          string `json:"status"`
						Conclusion      string `json:"conclusion"`
						CreatedAt       string `json:"created_at"`
					}

					responseList := make([]reportResponse, 0, len(reportsList))
					for _, report := range reportsList {
						responseList = append(responseList, reportResponse{
							FHIRID:          report.FHIRResourceID,
							EncounterFHIRID: report.EncounterFHIRID,
							PatientFHIRID:   report.PatientFHIRID,
							ReportDisplay:   report.ReportDisplay,
							Status:          report.Status,
							Conclusion:      report.Conclusion,
							CreatedAt:       report.IssuedAt.Format(time.RFC3339),
						})
					}

					httpResponseWriter.Header().Set("Content-Type", "application/json")
					httpResponseWriter.WriteHeader(http.StatusOK)
					json.NewEncoder(httpResponseWriter).Encode(responseList)
					return
				}

				if httpRequest.Method == http.MethodPost {
					contextWithValues, authIsOk := validateHTTPAuth(httpResponseWriter, httpRequest, []auth.Role{auth.RoleDoctor, auth.RoleNurse})
					if !authIsOk {
						return
					}

					var payload struct {
						PatientFHIRID string `json:"patient_fhir_id"`
						ReportDisplay string `json:"report_display"`
						Conclusion    string `json:"conclusion"`
					}

					if payloadDecodeErr := json.NewDecoder(httpRequest.Body).Decode(&payload); payloadDecodeErr != nil {
						httpResponseWriter.Header().Set("Content-Type", "application/json")
						httpResponseWriter.WriteHeader(http.StatusBadRequest)
						json.NewEncoder(httpResponseWriter).Encode(map[string]string{"error": "Payload inválido."})
						return
					}

					newReport := &clinical.DiagnosticReport{
						EncounterFHIRID: encounterFHIRID,
						PatientFHIRID:   payload.PatientFHIRID,
						ReportCode:      "24323-8",
						ReportDisplay:   payload.ReportDisplay,
						Status:          "final",
						Conclusion:      payload.Conclusion,
						IssuedAt:        time.Now(),
					}

					createdReport, createErr := clinicalService.CreateDiagnosticReport(contextWithValues, newReport)
					if createErr != nil {
						httpResponseWriter.Header().Set("Content-Type", "application/json")
						httpResponseWriter.WriteHeader(http.StatusInternalServerError)
						json.NewEncoder(httpResponseWriter).Encode(map[string]string{"error": "Erro ao criar laudo."})
						return
					}

					httpResponseWriter.Header().Set("Content-Type", "application/json")
					httpResponseWriter.WriteHeader(http.StatusCreated)
					json.NewEncoder(httpResponseWriter).Encode(map[string]interface{}{
						"fhir_id":           createdReport.FHIRResourceID,
						"encounter_fhir_id": createdReport.EncounterFHIRID,
						"patient_fhir_id":   createdReport.PatientFHIRID,
						"report_display":    createdReport.ReportDisplay,
						"status":            createdReport.Status,
						"conclusion":        createdReport.Conclusion,
						"created_at":        createdReport.IssuedAt.Format(time.RFC3339),
					})
					return
				}

				http.Error(httpResponseWriter, "Method Not Allowed", http.StatusMethodNotAllowed)
				return
			}
		}

		http.Error(httpResponseWriter, "Not Found", http.StatusNotFound)
	})

	httpServeMux.HandleFunc("/api/studies/", func(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
		if corsOptionHandler(httpResponseWriter, httpRequest) {
			return
		}

		imagingHTTPHandler.HandleStudy(httpResponseWriter, httpRequest)
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
		Addr:              ":" + appConfig.HTTPPort,
		Handler:           httpServeMux,
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
	imagingWorker.Stop()
	applicationServer.GRPCServer.GracefulStop()

	ctxShutdownTimeout, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdown()

	if shutdownError := httpServer.Shutdown(ctxShutdownTimeout); shutdownError != nil {
		slog.Error("Failed to shutdown HTTP server", "error", shutdownError)
	}

	slog.Info("Server stopped")
}
