package allergy

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/healthcare/backend/internal/api/middleware"
	"github.com/healthcare/backend/internal/api/render"
	"github.com/healthcare/backend/internal/shared/role"
)

type HTTPHandler struct {
	service Service
}

func NewHTTPHandler(service Service) *HTTPHandler {
	return &HTTPHandler{
		service: service,
	}
}

func (handler *HTTPHandler) RegisterRoutes(mux *http.ServeMux) {
	clinicalWrite := middleware.RequireRoles(role.RoleDoctor, role.RoleNurse)
	clinicalRead := middleware.RequireRoles(role.RoleAdmin, role.RoleDoctor, role.RoleNurse)

	mux.Handle("GET /api/v1/patients/{patientFhirId}/allergies", clinicalRead(http.HandlerFunc(handler.ListAllergiesByPatient)))
	mux.Handle("POST /api/v1/patients/{patientFhirId}/allergies", clinicalWrite(http.HandlerFunc(handler.CreateAllergy)))
}

// ListAllergiesByPatient godoc
//
//	@Summary		List allergies by patient
//	@Description	Returns all allergy intolerances for a patient
//	@Tags			allergies
//	@Accept			json
//	@Produce		json
//	@Param			patientFhirId	path	string	true	"Patient FHIR ID"
//	@Success		200				{array}	AllergyResponse
//	@Failure		500				{object}	map[string]string
//	@Router			/patients/{patientFhirId}/allergies [get]
func (handler *HTTPHandler) ListAllergiesByPatient(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
	patientFhirID := httpRequest.PathValue("patientFhirId")

	allergiesList, allergiesErr := handler.service.GetAllergyIntolerancesByPatient(httpRequest.Context(), patientFhirID)
	if allergiesErr != nil {
		slog.Error("failed to list allergies", "error", allergiesErr, "patient_fhir_id", patientFhirID, "request_id", middleware.GetRequestID(httpRequest.Context()))
		render.Error(httpResponseWriter, http.StatusInternalServerError, "Erro ao carregar alergias do paciente.")
		return
	}

	type allergyResponse struct {
		FhirID          string `json:"fhir_id"`
		PatientFhirID   string `json:"patient_fhir_id"`
		AllergenCode    string `json:"allergen_code"`
		AllergenDisplay string `json:"allergen_display"`
		ClinicalStatus  string `json:"clinical_status"`
		Reaction        string `json:"reaction"`
		CreatedAt       string `json:"created_at"`
	}

	responseList := make([]allergyResponse, 0, len(allergiesList))
	for _, allergy := range allergiesList {
		responseList = append(responseList, allergyResponse{
			FhirID:          allergy.FHIRResourceID,
			PatientFhirID:   allergy.PatientFHIRID,
			AllergenCode:    allergy.AllergenCode,
			AllergenDisplay: allergy.AllergenDisplay,
			ClinicalStatus:  allergy.ClinicalStatus,
			Reaction:        allergy.Reaction,
			CreatedAt:       allergy.RecordedAt.Format(time.RFC3339),
		})
	}

	render.JSON(httpResponseWriter, http.StatusOK, responseList)
}

// CreateAllergy godoc
//
//	@Summary		Create an allergy
//	@Description	Creates a new allergy intolerance record for a patient
//	@Tags			allergies
//	@Accept			json
//	@Produce		json
//	@Param			patientFhirId	path	string	true	"Patient FHIR ID"
//	@Param			body			body	CreateAllergyRequest	true	"Allergy data"
//	@Success		201				{object}	CreateAllergyResponse
//	@Failure		400				{object}	map[string]string
//	@Failure		500				{object}	map[string]string
//	@Router			/patients/{patientFhirId}/allergies [post]
func (handler *HTTPHandler) CreateAllergy(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
	patientFhirID := httpRequest.PathValue("patientFhirId")

	var payload struct {
		AllergenCode    string `json:"allergen_code"`
		AllergenDisplay string `json:"allergen_display"`
		Reaction        string `json:"reaction"`
	}

	if payloadDecodeErr := json.NewDecoder(httpRequest.Body).Decode(&payload); payloadDecodeErr != nil {
		render.Error(httpResponseWriter, http.StatusBadRequest, "Payload inválido.")
		return
	}

	newAllergy := &Allergy{
		PatientFHIRID:   patientFhirID,
		AllergenCode:    payload.AllergenCode,
		AllergenDisplay: payload.AllergenDisplay,
		ClinicalStatus:  "active",
		Reaction:        payload.Reaction,
		RecordedAt:      time.Now(),
	}

	createdAllergy, createErr := handler.service.CreateAllergyIntolerance(httpRequest.Context(), newAllergy)
	if createErr != nil {
		slog.Error("failed to create allergy", "error", createErr, "patient_fhir_id", patientFhirID, "request_id", middleware.GetRequestID(httpRequest.Context()))
		render.Error(httpResponseWriter, http.StatusInternalServerError, "Erro ao criar alergia.")
		return
	}

	render.JSON(httpResponseWriter, http.StatusCreated, map[string]interface{}{
		"fhir_id":          createdAllergy.FHIRResourceID,
		"patient_fhir_id":  createdAllergy.PatientFHIRID,
		"allergen_code":    createdAllergy.AllergenCode,
		"allergen_display": createdAllergy.AllergenDisplay,
		"clinical_status":  createdAllergy.ClinicalStatus,
		"reaction":         createdAllergy.Reaction,
		"created_at":       createdAllergy.RecordedAt.Format(time.RFC3339),
	})
}

type AllergyResponse struct {
	FhirID          string `json:"fhir_id"`
	PatientFhirID   string `json:"patient_fhir_id"`
	AllergenCode    string `json:"allergen_code"`
	AllergenDisplay string `json:"allergen_display"`
	ClinicalStatus  string `json:"clinical_status"`
	Reaction        string `json:"reaction"`
	CreatedAt       string `json:"created_at"`
}

type CreateAllergyRequest struct {
	AllergenCode    string `json:"allergen_code"`
	AllergenDisplay string `json:"allergen_display"`
	Reaction        string `json:"reaction"`
}

type CreateAllergyResponse struct {
	FhirID          string `json:"fhir_id"`
	PatientFhirID   string `json:"patient_fhir_id"`
	AllergenCode    string `json:"allergen_code"`
	AllergenDisplay string `json:"allergen_display"`
	ClinicalStatus  string `json:"clinical_status"`
	Reaction        string `json:"reaction"`
	CreatedAt       string `json:"created_at"`
}
