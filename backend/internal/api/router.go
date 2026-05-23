package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"
	"log/slog"

	"github.com/google/uuid"
	"github.com/healthcare/backend/internal/api/middleware"
	"github.com/healthcare/backend/internal/modules/auth"
	"github.com/healthcare/backend/internal/modules/clinical"
	"github.com/healthcare/backend/internal/modules/imaging"
	"github.com/healthcare/backend/internal/modules/patients"
	"github.com/healthcare/backend/internal/modules/staff"
	"github.com/healthcare/backend/internal/modules/telemetry"
)

func NewRouter(
	authService auth.Service,
	patientsService patients.Service,
	clinicalService clinical.Service,
	imagingHTTPHandler *imaging.HTTPHandler,
	staffService staff.Service,
	telemetryService telemetry.Service,
	secureCookies bool,
) http.Handler {
	httpServeMux := http.NewServeMux()

	httpServeMux.HandleFunc("/api/auth/login", func(writer http.ResponseWriter, request *http.Request) {

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
			json.NewEncoder(writer).Encode(map[string]string{"error": "Credenciais invÃ¡lidas."})
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


	httpServeMux.HandleFunc("/api/patients", func(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {

		if httpRequest.Method == http.MethodGet {
			contextWithValues, authIsOk := middleware.ValidateHTTPAuth(httpResponseWriter, httpRequest, []auth.Role{auth.RoleAdmin, auth.RoleDoctor, auth.RoleNurse, auth.RoleReception})
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
			contextWithValues, authIsOk := middleware.ValidateHTTPAuth(httpResponseWriter, httpRequest, []auth.Role{auth.RoleAdmin, auth.RoleReception})
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
				json.NewEncoder(httpResponseWriter).Encode(map[string]string{"error": "Payload invÃ¡lido."})
				return
			}

			patient, createPatientErr := patientsService.CreatePatient(contextWithValues, payload.FullName, payload.BirthDate, payload.DocumentID, payload.PhoneNumber)
			if createPatientErr != nil {
				httpResponseWriter.Header().Set("Content-Type", "application/json")
				if errors.Is(createPatientErr, patients.ErrPatientAlreadyExists) {
					httpResponseWriter.WriteHeader(http.StatusConflict)
					json.NewEncoder(httpResponseWriter).Encode(map[string]string{"error": "Paciente com este documento jÃ¡ cadastrado."})
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

			contextWithValues, authIsOk := middleware.ValidateHTTPAuth(httpResponseWriter, httpRequest, []auth.Role{auth.RoleAdmin, auth.RoleDoctor, auth.RoleNurse, auth.RoleReception})
			if !authIsOk {
				return
			}

			patient, getPatientErr := patientsService.GetPatient(contextWithValues, parts[0])
			if getPatientErr != nil {
				httpResponseWriter.Header().Set("Content-Type", "application/json")
				httpResponseWriter.WriteHeader(http.StatusNotFound)
				json.NewEncoder(httpResponseWriter).Encode(map[string]string{"error": "Paciente nÃ£o encontrado."})
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
					contextWithValues, authIsOk := middleware.ValidateHTTPAuth(httpResponseWriter, httpRequest, []auth.Role{auth.RoleAdmin, auth.RoleDoctor, auth.RoleNurse, auth.RoleReception})
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
					contextWithValues, authIsOk := middleware.ValidateHTTPAuth(httpResponseWriter, httpRequest, []auth.Role{auth.RoleAdmin, auth.RoleDoctor, auth.RoleNurse, auth.RoleReception})
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
						json.NewEncoder(httpResponseWriter).Encode(map[string]string{"error": "Payload invÃ¡lido."})
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

				contextWithValues, authIsOk := middleware.ValidateHTTPAuth(httpResponseWriter, httpRequest, []auth.Role{auth.RoleAdmin, auth.RoleDoctor, auth.RoleNurse})
				if !authIsOk {
					return
				}

				observationsList, observationsErr := clinicalService.GetObservationsByPatient(contextWithValues, patientFHIRID)
				if observationsErr != nil {
					httpResponseWriter.Header().Set("Content-Type", "application/json")
					httpResponseWriter.WriteHeader(http.StatusInternalServerError)
					json.NewEncoder(httpResponseWriter).Encode(map[string]string{"error": "Erro ao carregar observaÃ§Ãµes do paciente."})
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
					contextWithValues, authIsOk := middleware.ValidateHTTPAuth(httpResponseWriter, httpRequest, []auth.Role{auth.RoleAdmin, auth.RoleDoctor, auth.RoleNurse})
					if !authIsOk {
						return
					}

					conditionsList, conditionsErr := clinicalService.GetConditionsByPatient(contextWithValues, patientFHIRID)
					if conditionsErr != nil {
						httpResponseWriter.Header().Set("Content-Type", "application/json")
						httpResponseWriter.WriteHeader(http.StatusInternalServerError)
						json.NewEncoder(httpResponseWriter).Encode(map[string]string{"error": "Erro ao carregar diagnÃ³sticos do paciente."})
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
					contextWithValues, authIsOk := middleware.ValidateHTTPAuth(httpResponseWriter, httpRequest, []auth.Role{auth.RoleDoctor, auth.RoleNurse})
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
						json.NewEncoder(httpResponseWriter).Encode(map[string]string{"error": "Payload invÃ¡lido."})
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
						json.NewEncoder(httpResponseWriter).Encode(map[string]string{"error": "Erro ao criar diagnÃ³stico."})
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

			if subResource == "allergies" {
				if httpRequest.Method == http.MethodGet {
					contextWithValues, authIsOk := middleware.ValidateHTTPAuth(httpResponseWriter, httpRequest, []auth.Role{auth.RoleAdmin, auth.RoleDoctor, auth.RoleNurse})
					if !authIsOk {
						return
					}

					allergiesList, allergiesErr := clinicalService.GetAllergyIntolerancesByPatient(contextWithValues, patientFHIRID)
					if allergiesErr != nil {
						httpResponseWriter.Header().Set("Content-Type", "application/json")
						httpResponseWriter.WriteHeader(http.StatusInternalServerError)
						json.NewEncoder(httpResponseWriter).Encode(map[string]string{"error": "Erro ao carregar alergias do paciente."})
						return
					}

					type allergyResponse struct {
						FHIRID          string `json:"fhir_id"`
						PatientFHIRID   string `json:"patient_fhir_id"`
						AllergenCode    string `json:"allergen_code"`
						AllergenDisplay string `json:"allergen_display"`
						ClinicalStatus  string `json:"clinical_status"`
						Reaction        string `json:"reaction"`
						CreatedAt       string `json:"created_at"`
					}

					responseList := make([]allergyResponse, 0, len(allergiesList))
					for _, allergy := range allergiesList {
						responseList = append(responseList, allergyResponse{
							FHIRID:          allergy.FHIRResourceID,
							PatientFHIRID:   allergy.PatientFHIRID,
							AllergenCode:    allergy.AllergenCode,
							AllergenDisplay: allergy.AllergenDisplay,
							ClinicalStatus:  allergy.ClinicalStatus,
							Reaction:        allergy.Reaction,
							CreatedAt:       allergy.RecordedAt.Format(time.RFC3339),
						})
					}

					httpResponseWriter.Header().Set("Content-Type", "application/json")
					httpResponseWriter.WriteHeader(http.StatusOK)
					json.NewEncoder(httpResponseWriter).Encode(responseList)
					return
				}

				if httpRequest.Method == http.MethodPost {
					contextWithValues, authIsOk := middleware.ValidateHTTPAuth(httpResponseWriter, httpRequest, []auth.Role{auth.RoleDoctor, auth.RoleNurse})
					if !authIsOk {
						return
					}

					var payload struct {
						AllergenCode    string `json:"allergen_code"`
						AllergenDisplay string `json:"allergen_display"`
						Reaction        string `json:"reaction"`
					}

					if payloadDecodeErr := json.NewDecoder(httpRequest.Body).Decode(&payload); payloadDecodeErr != nil {
						httpResponseWriter.Header().Set("Content-Type", "application/json")
						httpResponseWriter.WriteHeader(http.StatusBadRequest)
						json.NewEncoder(httpResponseWriter).Encode(map[string]string{"error": "Payload inválido."})
						return
					}

					newAllergy := &clinical.AllergyIntolerance{
						PatientFHIRID:   patientFHIRID,
						AllergenCode:    payload.AllergenCode,
						AllergenDisplay: payload.AllergenDisplay,
						ClinicalStatus:  "active",
						Reaction:        payload.Reaction,
						RecordedAt:      time.Now(),
					}

					createdAllergy, createErr := clinicalService.CreateAllergyIntolerance(contextWithValues, newAllergy)
					if createErr != nil {
						httpResponseWriter.Header().Set("Content-Type", "application/json")
						httpResponseWriter.WriteHeader(http.StatusInternalServerError)
						json.NewEncoder(httpResponseWriter).Encode(map[string]string{"error": "Erro ao criar alergia."})
						return
					}

					httpResponseWriter.Header().Set("Content-Type", "application/json")
					httpResponseWriter.WriteHeader(http.StatusCreated)
					json.NewEncoder(httpResponseWriter).Encode(map[string]interface{}{
						"fhir_id":          createdAllergy.FHIRResourceID,
						"patient_fhir_id":  createdAllergy.PatientFHIRID,
						"allergen_code":    createdAllergy.AllergenCode,
						"allergen_display": createdAllergy.AllergenDisplay,
						"clinical_status":  createdAllergy.ClinicalStatus,
						"reaction":         createdAllergy.Reaction,
						"created_at":       createdAllergy.RecordedAt.Format(time.RFC3339),
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
					contextWithValues, authIsOk := middleware.ValidateHTTPAuth(httpResponseWriter, httpRequest, []auth.Role{auth.RoleAdmin, auth.RoleDoctor, auth.RoleNurse})
					if !authIsOk {
						return
					}

					observationsList, observationsErr := clinicalService.GetObservationsByEncounter(contextWithValues, encounterFHIRID)
					if observationsErr != nil {
						httpResponseWriter.Header().Set("Content-Type", "application/json")
						httpResponseWriter.WriteHeader(http.StatusInternalServerError)
						json.NewEncoder(httpResponseWriter).Encode(map[string]string{"error": "Erro ao carregar observaÃ§Ãµes da consulta."})
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
					contextWithValues, authIsOk := middleware.ValidateHTTPAuth(httpResponseWriter, httpRequest, []auth.Role{auth.RoleDoctor, auth.RoleNurse})
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
						json.NewEncoder(httpResponseWriter).Encode(map[string]string{"error": "Payload invÃ¡lido."})
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
						json.NewEncoder(httpResponseWriter).Encode(map[string]string{"error": "Erro ao criar observaÃ§Ã£o."})
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
					contextWithValues, authIsOk := middleware.ValidateHTTPAuth(httpResponseWriter, httpRequest, []auth.Role{auth.RoleAdmin, auth.RoleDoctor, auth.RoleNurse})
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
					contextWithValues, authIsOk := middleware.ValidateHTTPAuth(httpResponseWriter, httpRequest, []auth.Role{auth.RoleDoctor, auth.RoleNurse})
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
						json.NewEncoder(httpResponseWriter).Encode(map[string]string{"error": "Payload invÃ¡lido."})
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
			if subResource == "medications" {
				if httpRequest.Method == http.MethodGet {
					contextWithValues, authIsOk := middleware.ValidateHTTPAuth(httpResponseWriter, httpRequest, []auth.Role{auth.RoleAdmin, auth.RoleDoctor, auth.RoleNurse})
					if !authIsOk {
						return
					}

					medicationsList, medicationsErr := clinicalService.GetMedicationRequestsByEncounter(contextWithValues, encounterFHIRID)
					if medicationsErr != nil {
						httpResponseWriter.Header().Set("Content-Type", "application/json")
						httpResponseWriter.WriteHeader(http.StatusInternalServerError)
						json.NewEncoder(httpResponseWriter).Encode(map[string]string{"error": "Erro ao carregar prescrições da consulta."})
						return
					}

					type medicationResponse struct {
						FHIRID             string `json:"fhir_id"`
						EncounterFHIRID    string `json:"encounter_fhir_id"`
						PatientFHIRID      string `json:"patient_fhir_id"`
						PractitionerFHIRID string `json:"practitioner_fhir_id"`
						MedicationCode     string `json:"medication_code"`
						MedicationName     string `json:"medication_name"`
						DosageInstructions string `json:"dosage_instructions"`
						Status             string `json:"status"`
						CreatedAt          string `json:"created_at"`
					}

					responseList := make([]medicationResponse, 0, len(medicationsList))
					for _, med := range medicationsList {
						responseList = append(responseList, medicationResponse{
							FHIRID:             med.FHIRResourceID,
							EncounterFHIRID:    med.EncounterFHIRID,
							PatientFHIRID:      med.PatientFHIRID,
							PractitionerFHIRID: med.PractitionerFHIRID,
							MedicationCode:     med.MedicationCode,
							MedicationName:     med.MedicationName,
							DosageInstructions: med.DosageInstructions,
							Status:             med.Status,
							CreatedAt:          med.IssuedAt.Format(time.RFC3339),
						})
					}

					httpResponseWriter.Header().Set("Content-Type", "application/json")
					httpResponseWriter.WriteHeader(http.StatusOK)
					json.NewEncoder(httpResponseWriter).Encode(responseList)
					return
				}

				if httpRequest.Method == http.MethodPost {
					contextWithValues, authIsOk := middleware.ValidateHTTPAuth(httpResponseWriter, httpRequest, []auth.Role{auth.RoleDoctor})
					if !authIsOk {
						return
					}

					var payload struct {
						PatientFHIRID      string `json:"patient_fhir_id"`
						PractitionerFHIRID string `json:"practitioner_fhir_id"`
						MedicationCode     string `json:"medication_code"`
						MedicationName     string `json:"medication_name"`
						DosageInstructions string `json:"dosage_instructions"`
					}

					if payloadDecodeErr := json.NewDecoder(httpRequest.Body).Decode(&payload); payloadDecodeErr != nil {
						httpResponseWriter.Header().Set("Content-Type", "application/json")
						httpResponseWriter.WriteHeader(http.StatusBadRequest)
						json.NewEncoder(httpResponseWriter).Encode(map[string]string{"error": "Payload inválido."})
						return
					}

					newMedication := &clinical.MedicationRequest{
						EncounterFHIRID:    encounterFHIRID,
						PatientFHIRID:      payload.PatientFHIRID,
						PractitionerFHIRID: payload.PractitionerFHIRID,
						MedicationCode:     payload.MedicationCode,
						MedicationName:     payload.MedicationName,
						DosageInstructions: payload.DosageInstructions,
						Status:             "active",
						IssuedAt:           time.Now(),
					}

					createdMedication, createErr := clinicalService.CreateMedicationRequest(contextWithValues, newMedication)
					if createErr != nil {
						httpResponseWriter.Header().Set("Content-Type", "application/json")
						httpResponseWriter.WriteHeader(http.StatusInternalServerError)
						json.NewEncoder(httpResponseWriter).Encode(map[string]string{"error": "Erro ao criar prescrição."})
						return
					}

					httpResponseWriter.Header().Set("Content-Type", "application/json")
					httpResponseWriter.WriteHeader(http.StatusCreated)
					json.NewEncoder(httpResponseWriter).Encode(map[string]interface{}{
						"fhir_id":              createdMedication.FHIRResourceID,
						"encounter_fhir_id":    createdMedication.EncounterFHIRID,
						"patient_fhir_id":      createdMedication.PatientFHIRID,
						"practitioner_fhir_id": createdMedication.PractitionerFHIRID,
						"medication_code":      createdMedication.MedicationCode,
						"medication_name":      createdMedication.MedicationName,
						"dosage_instructions":  createdMedication.DosageInstructions,
						"status":               createdMedication.Status,
						"created_at":           createdMedication.IssuedAt.Format(time.RFC3339),
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

		imagingHTTPHandler.HandleStudy(httpResponseWriter, httpRequest)
	})

	httpServeMux.HandleFunc("/api/staff/employees", func(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
		if httpRequest.Method == http.MethodGet {
			contextWithValues, authIsOk := middleware.ValidateHTTPAuth(httpResponseWriter, httpRequest, []auth.Role{auth.RoleAdmin, auth.RoleDoctor, auth.RoleNurse, auth.RoleReception})
			if !authIsOk {
				return
			}

			employeesList, employeesErr := staffService.ListEmployees(contextWithValues)
			if employeesErr != nil {
				httpResponseWriter.Header().Set("Content-Type", "application/json")
				httpResponseWriter.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(httpResponseWriter).Encode(map[string]string{"error": "Erro ao listar corpo clínico."})
				return
			}

			httpResponseWriter.Header().Set("Content-Type", "application/json")
			httpResponseWriter.WriteHeader(http.StatusOK)
			json.NewEncoder(httpResponseWriter).Encode(employeesList)
			return
		}

		if httpRequest.Method == http.MethodPost {
			contextWithValues, authIsOk := middleware.ValidateHTTPAuth(httpResponseWriter, httpRequest, []auth.Role{auth.RoleAdmin})
			if !authIsOk {
				return
			}

			var payload struct {
				UserID    string `json:"user_id"`
				FullName  string `json:"full_name"`
				Email     string `json:"email"`
				Role      string `json:"role"`
				CRMNumber string `json:"crm_number"`
			}

			if payloadDecodeErr := json.NewDecoder(httpRequest.Body).Decode(&payload); payloadDecodeErr != nil {
				httpResponseWriter.Header().Set("Content-Type", "application/json")
				httpResponseWriter.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(httpResponseWriter).Encode(map[string]string{"error": "Payload inválido."})
				return
			}

			userIDParsed, err := uuid.Parse(payload.UserID)
			if err != nil {
				httpResponseWriter.Header().Set("Content-Type", "application/json")
				httpResponseWriter.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(httpResponseWriter).Encode(map[string]string{"error": "User ID inválido."})
				return
			}

			employee, createErr := staffService.CreateEmployee(contextWithValues, userIDParsed, payload.FullName, payload.Email, payload.Role, payload.CRMNumber)
			if createErr != nil {
				httpResponseWriter.Header().Set("Content-Type", "application/json")
				httpResponseWriter.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(httpResponseWriter).Encode(map[string]string{"error": "Erro ao registrar profissional."})
				return
			}

			httpResponseWriter.Header().Set("Content-Type", "application/json")
			httpResponseWriter.WriteHeader(http.StatusCreated)
			json.NewEncoder(httpResponseWriter).Encode(map[string]string{
				"employee_id": employee.ID.String(),
			})
			return
		}

		http.Error(httpResponseWriter, "Method Not Allowed", http.StatusMethodNotAllowed)
	})

	httpServeMux.HandleFunc("/api/telemetry/rooms", func(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
		if httpRequest.Method != http.MethodGet {
			http.Error(httpResponseWriter, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		contextWithValues, authIsOk := middleware.ValidateHTTPAuth(httpResponseWriter, httpRequest, []auth.Role{auth.RoleAdmin, auth.RoleDoctor, auth.RoleNurse, auth.RoleReception})
		if !authIsOk {
			return
		}

		roomsList, roomsErr := telemetryService.GetRooms(contextWithValues)
		if roomsErr != nil {
			httpResponseWriter.Header().Set("Content-Type", "application/json")
			httpResponseWriter.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(httpResponseWriter).Encode(map[string]string{"error": "Erro ao listar salas."})
			return
		}

		httpResponseWriter.Header().Set("Content-Type", "application/json")
		httpResponseWriter.WriteHeader(http.StatusOK)
		json.NewEncoder(httpResponseWriter).Encode(roomsList)
	})

	httpServeMux.HandleFunc("/api/telemetry/rooms/", func(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
		remainingPath := strings.TrimPrefix(httpRequest.URL.Path, "/api/telemetry/rooms/")
		if remainingPath == "" {
			http.Error(httpResponseWriter, "Not Found", http.StatusNotFound)
			return
		}

		parts := strings.Split(remainingPath, "/")
		if len(parts) == 2 {
			roomIDString := parts[0]
			subResource := parts[1]

			roomIDParsed, err := uuid.Parse(roomIDString)
			if err != nil {
				http.Error(httpResponseWriter, "Bad Request", http.StatusBadRequest)
				return
			}

			if subResource == "unlock" {
				if httpRequest.Method != http.MethodPost {
					http.Error(httpResponseWriter, "Method Not Allowed", http.StatusMethodNotAllowed)
					return
				}

				contextWithValues, authIsOk := middleware.ValidateHTTPAuth(httpResponseWriter, httpRequest, []auth.Role{auth.RoleAdmin, auth.RoleDoctor, auth.RoleNurse})
				if !authIsOk {
					return
				}

				var payload struct {
					Passcode string `json:"passcode"`
				}

				if payloadDecodeErr := json.NewDecoder(httpRequest.Body).Decode(&payload); payloadDecodeErr != nil {
					http.Error(httpResponseWriter, "Bad Request", http.StatusBadRequest)
					return
				}

				unlockedRoom, unlockErr := telemetryService.UnlockRoom(contextWithValues, roomIDParsed, payload.Passcode)
				if unlockErr != nil {
					httpResponseWriter.Header().Set("Content-Type", "application/json")
					httpResponseWriter.WriteHeader(http.StatusUnauthorized)
					json.NewEncoder(httpResponseWriter).Encode(map[string]string{"error": "Senha inválida."})
					return
				}

				httpResponseWriter.Header().Set("Content-Type", "application/json")
				httpResponseWriter.WriteHeader(http.StatusOK)
				json.NewEncoder(httpResponseWriter).Encode(map[string]interface{}{
					"success":  true,
					"roomName": unlockedRoom.Name,
				})
				return
			}

			if subResource == "beds" {
				if httpRequest.Method != http.MethodGet {
					http.Error(httpResponseWriter, "Method Not Allowed", http.StatusMethodNotAllowed)
					return
				}

				contextWithValues, authIsOk := middleware.ValidateHTTPAuth(httpResponseWriter, httpRequest, []auth.Role{auth.RoleAdmin, auth.RoleDoctor, auth.RoleNurse})
				if !authIsOk {
					return
				}

				bedsList, bedsErr := telemetryService.GetBeds(contextWithValues, roomIDParsed)
				if bedsErr != nil {
					httpResponseWriter.Header().Set("Content-Type", "application/json")
					httpResponseWriter.WriteHeader(http.StatusInternalServerError)
					json.NewEncoder(httpResponseWriter).Encode(map[string]string{"error": "Erro ao listar leitos."})
					return
				}

				httpResponseWriter.Header().Set("Content-Type", "application/json")
				httpResponseWriter.WriteHeader(http.StatusOK)
				json.NewEncoder(httpResponseWriter).Encode(bedsList)
				return
			}
		}

		http.Error(httpResponseWriter, "Not Found", http.StatusNotFound)
	})

	httpServeMux.HandleFunc("/api/telemetry/beds/", func(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
		remainingPath := strings.TrimPrefix(httpRequest.URL.Path, "/api/telemetry/beds/")
		if remainingPath == "" {
			http.Error(httpResponseWriter, "Not Found", http.StatusNotFound)
			return
		}

		parts := strings.Split(remainingPath, "/")
		if len(parts) == 2 {
			bedIDString := parts[0]
			subResource := parts[1]

			bedIDParsed, err := uuid.Parse(bedIDString)
			if err != nil {
				http.Error(httpResponseWriter, "Bad Request", http.StatusBadRequest)
				return
			}

			if subResource == "condition" {
				if httpRequest.Method != http.MethodPost {
					http.Error(httpResponseWriter, "Method Not Allowed", http.StatusMethodNotAllowed)
					return
				}

				contextWithValues, authIsOk := middleware.ValidateHTTPAuth(httpResponseWriter, httpRequest, []auth.Role{auth.RoleDoctor, auth.RoleNurse})
				if !authIsOk {
					return
				}

				var payload struct {
					Bpm         int32   `json:"bpm"`
					Spo2        int32   `json:"spo2"`
					Temperature float64 `json:"temperature"`
					Status      string  `json:"status"`
					Condition   string  `json:"condition"`
				}

				if payloadDecodeErr := json.NewDecoder(httpRequest.Body).Decode(&payload); payloadDecodeErr != nil {
					http.Error(httpResponseWriter, "Bad Request", http.StatusBadRequest)
					return
				}

				updateErr := telemetryService.UpdateBedCondition(contextWithValues, bedIDParsed, payload.Bpm, payload.Spo2, payload.Temperature, payload.Status, payload.Condition)
				if updateErr != nil {
					httpResponseWriter.Header().Set("Content-Type", "application/json")
					httpResponseWriter.WriteHeader(http.StatusInternalServerError)
					json.NewEncoder(httpResponseWriter).Encode(map[string]string{"error": "Erro ao atualizar leito."})
					return
				}

				httpResponseWriter.Header().Set("Content-Type", "application/json")
				httpResponseWriter.WriteHeader(http.StatusOK)
				json.NewEncoder(httpResponseWriter).Encode(map[string]bool{"success": true})
				return
			}
		}

		http.Error(httpResponseWriter, "Not Found", http.StatusNotFound)
	})

	return middleware.CORS(secureCookies)(httpServeMux)
}
