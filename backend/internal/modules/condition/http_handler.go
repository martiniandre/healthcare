package condition

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

	mux.Handle("GET /api/patients/{patientFhirId}/conditions", clinicalRead(http.HandlerFunc(handler.ListConditionsByPatient)))
	mux.Handle("POST /api/patients/{patientFhirId}/conditions", clinicalWrite(http.HandlerFunc(handler.CreateCondition)))
}

func (handler *HTTPHandler) ListConditionsByPatient(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
	patientFhirID := httpRequest.PathValue("patientFhirId")

	conditionsList, conditionsErr := handler.service.GetConditionsByPatient(httpRequest.Context(), patientFhirID)
	if conditionsErr != nil {
		slog.Error("failed to list conditions", "error", conditionsErr, "patient_fhir_id", patientFhirID)
		render.Error(httpResponseWriter, http.StatusInternalServerError, "Erro ao carregar diagnósticos do paciente.")
		return
	}

	type conditionResponse struct {
		FhirID         string `json:"fhir_id"`
		PatientFhirID  string `json:"patient_fhir_id"`
		ICD10Code      string `json:"icd10_code"`
		CodeDisplay    string `json:"code_display"`
		ClinicalStatus string `json:"clinical_status"`
		CreatedAt      string `json:"created_at"`
	}

	responseList := make([]conditionResponse, 0, len(conditionsList))
	for _, condition := range conditionsList {
		responseList = append(responseList, conditionResponse{
			FhirID:         condition.FHIRResourceID,
			PatientFhirID:  condition.PatientFHIRID,
			ICD10Code:      condition.ICD10Code,
			CodeDisplay:    condition.CodeDisplay,
			ClinicalStatus: condition.ClinicalStatus,
			CreatedAt:      condition.OnsetAt.Format(time.RFC3339),
		})
	}

	render.JSON(httpResponseWriter, http.StatusOK, responseList)
}

func (handler *HTTPHandler) CreateCondition(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
	patientFhirID := httpRequest.PathValue("patientFhirId")

	var payload struct {
		ICD10Code   string `json:"icd10_code"`
		CodeDisplay string `json:"code_display"`
	}

	if payloadDecodeErr := json.NewDecoder(httpRequest.Body).Decode(&payload); payloadDecodeErr != nil {
		render.Error(httpResponseWriter, http.StatusBadRequest, "Payload inválido.")
		return
	}

	newCondition := &Condition{
		PatientFHIRID:  patientFhirID,
		ICD10Code:      payload.ICD10Code,
		CodeDisplay:    payload.CodeDisplay,
		ClinicalStatus: "active",
		OnsetAt:        time.Now(),
	}

	createdCondition, createErr := handler.service.CreateCondition(httpRequest.Context(), newCondition)
	if createErr != nil {
		slog.Error("failed to create condition", "error", createErr, "patient_fhir_id", patientFhirID)
		render.Error(httpResponseWriter, http.StatusInternalServerError, "Erro ao criar diagnóstico.")
		return
	}

	render.JSON(httpResponseWriter, http.StatusCreated, map[string]interface{}{
		"fhir_id":         createdCondition.FHIRResourceID,
		"patient_fhir_id": createdCondition.PatientFHIRID,
		"icd10_code":      createdCondition.ICD10Code,
		"code_display":    createdCondition.CodeDisplay,
		"clinical_status": createdCondition.ClinicalStatus,
		"created_at":      createdCondition.OnsetAt.Format(time.RFC3339),
	})
}
