package patients

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/healthcare/backend/internal/api/middleware"
	"github.com/healthcare/backend/internal/api/render"
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
	mux.Handle("GET /api/v1/patients", middleware.RequirePolicy("GET /api/v1/patients")(http.HandlerFunc(handler.ListPatients)))
	mux.Handle("POST /api/v1/patients", middleware.RequirePolicy("POST /api/v1/patients")(http.HandlerFunc(handler.CreatePatient)))
	mux.Handle("GET /api/v1/patients/{patientFhirId}", middleware.RequirePolicy("GET /api/v1/patients/{patientFhirId}")(http.HandlerFunc(handler.GetPatient)))
}

// ListPatients godoc
//
//	@Summary		List all patients
//	@Description	Returns a paginated list of patients with optional search/filter
//	@Tags			patients
//	@Accept			json
//	@Produce		json
//	@Param			search			query		string	false	"Search term"
//	@Param			sortField		query		string	false	"Sort field"
//	@Param			sortDirection	query		string	false	"Sort direction (asc/desc)"
//	@Param			page			query		int		false	"Page number"			default(1)
//	@Param			limit			query		int		false	"Items per page"		default(50)
//	@Success		200				{object}	PatientListEnvelope
//	@Failure		500				{object}	map[string]string
//	@Router			/patients [get]
func (handler *HTTPHandler) ListPatients(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
	search := httpRequest.URL.Query().Get("search")
	sortField := httpRequest.URL.Query().Get("sortField")
	sortDirection := httpRequest.URL.Query().Get("sortDirection")
	page, _ := strconv.Atoi(httpRequest.URL.Query().Get("page"))
	if page <= 0 {
		page = 1
	}
	limit, _ := strconv.Atoi(httpRequest.URL.Query().Get("limit"))
	if limit <= 0 {
		limit = 50
	}

	patientsList, totalCount, listError := handler.service.ListPatients(httpRequest.Context(), search, sortField, sortDirection, page, limit)
	if listError != nil {
		slog.Error("failed to list patients", "error", listError, "request_id", middleware.GetRequestID(httpRequest.Context()))
		render.ErrorFromAppError(httpResponseWriter, listError)
		return
	}

	responseList := make([]PatientListResponse, 0, len(patientsList))
	for _, patient := range patientsList {
		responseList = append(responseList, PatientListResponse{
			PatientID:      patient.ID.String(),
			FHIRResourceID: patient.FHIRResourceID,
			FullName:       patient.FullName,
			BirthDate:      patient.BirthDate.Format("2006-01-02"),
			DocumentID:     patient.DocumentID,
			PhoneNumber:    patient.PhoneNumber,
		})
	}

	render.JSON(httpResponseWriter, http.StatusOK, PatientListEnvelope{
		Patients: responseList,
		Total:    totalCount,
		Page:     page,
		Limit:    limit,
	})
}

// CreatePatient godoc
//
//	@Summary		Create a new patient
//	@Description	Creates a new patient record and FHIR resource
//	@Tags			patients
//	@Accept			json
//	@Produce		json
//	@Param			body	body		CreatePatientRequest	true	"Patient data"
//	@Success		201		{object}	CreatePatientResponse
//	@Failure		400		{object}	map[string]string
//	@Failure		409		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Router			/patients [post]
func (handler *HTTPHandler) CreatePatient(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
	var payload CreatePatientRequest

	if payloadDecodeErr := json.NewDecoder(httpRequest.Body).Decode(&payload); payloadDecodeErr != nil {
		render.Error(httpResponseWriter, http.StatusBadRequest, "Payload inválido.")
		return
	}

	input := CreatePatientInput{
		FullName:    payload.FullName,
		BirthDate:   payload.BirthDate,
		DocumentID:  payload.DocumentID,
		PhoneNumber: payload.PhoneNumber,
	}

	patient, createPatientErr := handler.service.CreatePatient(httpRequest.Context(), input)
	if createPatientErr != nil {
		slog.Error("failed to create patient", "error", createPatientErr, "document_id", payload.DocumentID, "request_id", middleware.GetRequestID(httpRequest.Context()))
		render.ErrorFromAppError(httpResponseWriter, createPatientErr)
		return
	}

	render.JSON(httpResponseWriter, http.StatusCreated, CreatePatientResponse{
		PatientID:      patient.ID.String(),
		FHIRResourceID: patient.FHIRResourceID,
	})
}

// GetPatient godoc
//
//	@Summary		Get patient by FHIR ID
//	@Description	Returns a single patient by their FHIR resource ID
//	@Tags			patients
//	@Accept			json
//	@Produce		json
//	@Param			patientFhirId	path		string	true	"Patient FHIR ID"
//	@Success		200				{object}	map[string]interface{}
//	@Failure		400				{object}	map[string]string
//	@Failure		404				{object}	map[string]string
//	@Router			/patients/{patientFhirId} [get]
func (handler *HTTPHandler) GetPatient(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
	patientFhirID := httpRequest.PathValue("patientFhirId")
	if patientFhirID == "" {
		render.Error(httpResponseWriter, http.StatusBadRequest, "ID do paciente não informado.")
		return
	}

	patient, getPatientErr := handler.service.GetPatient(httpRequest.Context(), patientFhirID)
	if getPatientErr != nil {
		slog.Error("failed to get patient", "error", getPatientErr, "patient_fhir_id", patientFhirID, "request_id", middleware.GetRequestID(httpRequest.Context()))
		render.ErrorFromAppError(httpResponseWriter, getPatientErr)
		return
	}

	render.JSON(httpResponseWriter, http.StatusOK, PatientListResponse{
		PatientID:      patient.ID.String(),
		FHIRResourceID: patient.FHIRResourceID,
		FullName:       patient.FullName,
		BirthDate:      patient.BirthDate.Format("2006-01-02"),
		DocumentID:     patient.DocumentID,
		PhoneNumber:    patient.PhoneNumber,
	})
}

type PatientListEnvelope struct {
	Patients []PatientListResponse `json:"patients"`
	Total    int                   `json:"total"`
	Page     int                   `json:"page"`
	Limit    int                   `json:"limit"`
}

type PatientListResponse struct {
	PatientID      string `json:"patient_id"`
	FHIRResourceID string `json:"fhir_resource_id"`
	FullName       string `json:"full_name"`
	BirthDate      string `json:"birth_date"`
	DocumentID     string `json:"document_id"`
	PhoneNumber    string `json:"phone_number"`
}

type CreatePatientRequest struct {
	FullName    string `json:"full_name"`
	BirthDate   string `json:"birth_date"`
	DocumentID  string `json:"document_id"`
	PhoneNumber string `json:"phone_number"`
}

type CreatePatientResponse struct {
	PatientID      string `json:"patient_id"`
	FHIRResourceID string `json:"fhir_resource_id"`
}
