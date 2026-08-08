package portal

import (
	"log/slog"
	"net/http"

	"github.com/healthcare/backend/internal/api/middleware"
	"github.com/healthcare/backend/internal/api/render"
	"github.com/healthcare/backend/internal/shared/ctxkeys"
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
	mux.Handle("GET /api/v1/portal/dashboard", middleware.RequirePolicy("GET /api/v1/portal/dashboard")(http.HandlerFunc(handler.GetDashboard)))
	mux.Handle("GET /api/v1/portal/encounters", middleware.RequirePolicy("GET /api/v1/portal/encounters")(http.HandlerFunc(handler.GetEncounters)))
	mux.Handle("GET /api/v1/portal/observations", middleware.RequirePolicy("GET /api/v1/portal/observations")(http.HandlerFunc(handler.GetObservations)))
	mux.Handle("GET /api/v1/portal/conditions", middleware.RequirePolicy("GET /api/v1/portal/conditions")(http.HandlerFunc(handler.GetConditions)))
	mux.Handle("GET /api/v1/portal/medications", middleware.RequirePolicy("GET /api/v1/portal/medications")(http.HandlerFunc(handler.GetMedications)))
	mux.Handle("GET /api/v1/portal/reports", middleware.RequirePolicy("GET /api/v1/portal/reports")(http.HandlerFunc(handler.GetReports)))
	mux.Handle("GET /api/v1/portal/imaging", middleware.RequirePolicy("GET /api/v1/portal/imaging")(http.HandlerFunc(handler.GetImaging)))
}

func (handler *HTTPHandler) extractPatientFHIRID(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) string {
	userID, ok := httpRequest.Context().Value(ctxkeys.UserIDKey).(string)
	if !ok || userID == "" {
		render.Error(httpResponseWriter, http.StatusUnauthorized, "Usuário não autenticado.")
		return ""
	}
	return userID
}

func (handler *HTTPHandler) GetDashboard(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
	patientFHIRID := handler.extractPatientFHIRID(httpResponseWriter, httpRequest)
	if patientFHIRID == "" {
		return
	}

	dashboard, serviceError := handler.service.GetDashboard(httpRequest.Context(), patientFHIRID)
	if serviceError != nil {
		slog.Error("failed to get portal dashboard", "error", serviceError, "patient_id", patientFHIRID, "request_id", middleware.GetRequestID(httpRequest.Context()))
		render.Error(httpResponseWriter, http.StatusInternalServerError, "Erro ao carregar dashboard.")
		return
	}

	render.JSON(httpResponseWriter, http.StatusOK, dashboard)
}

func (handler *HTTPHandler) GetEncounters(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
	patientFHIRID := handler.extractPatientFHIRID(httpResponseWriter, httpRequest)
	if patientFHIRID == "" {
		return
	}

	encounters, serviceError := handler.service.GetEncounters(httpRequest.Context(), patientFHIRID)
	if serviceError != nil {
		slog.Error("failed to get portal encounters", "error", serviceError, "patient_id", patientFHIRID, "request_id", middleware.GetRequestID(httpRequest.Context()))
		render.Error(httpResponseWriter, http.StatusInternalServerError, "Erro ao carregar consultas.")
		return
	}

	render.JSON(httpResponseWriter, http.StatusOK, encounters)
}

func (handler *HTTPHandler) GetObservations(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
	patientFHIRID := handler.extractPatientFHIRID(httpResponseWriter, httpRequest)
	if patientFHIRID == "" {
		return
	}

	observations, serviceError := handler.service.GetObservations(httpRequest.Context(), patientFHIRID)
	if serviceError != nil {
		slog.Error("failed to get portal observations", "error", serviceError, "patient_id", patientFHIRID, "request_id", middleware.GetRequestID(httpRequest.Context()))
		render.Error(httpResponseWriter, http.StatusInternalServerError, "Erro ao carregar sinais vitais.")
		return
	}

	render.JSON(httpResponseWriter, http.StatusOK, observations)
}

func (handler *HTTPHandler) GetConditions(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
	patientFHIRID := handler.extractPatientFHIRID(httpResponseWriter, httpRequest)
	if patientFHIRID == "" {
		return
	}

	conditions, serviceError := handler.service.GetConditions(httpRequest.Context(), patientFHIRID)
	if serviceError != nil {
		slog.Error("failed to get portal conditions", "error", serviceError, "patient_id", patientFHIRID, "request_id", middleware.GetRequestID(httpRequest.Context()))
		render.Error(httpResponseWriter, http.StatusInternalServerError, "Erro ao carregar condições.")
		return
	}

	render.JSON(httpResponseWriter, http.StatusOK, conditions)
}

func (handler *HTTPHandler) GetMedications(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
	patientFHIRID := handler.extractPatientFHIRID(httpResponseWriter, httpRequest)
	if patientFHIRID == "" {
		return
	}

	medications, serviceError := handler.service.GetMedications(httpRequest.Context(), patientFHIRID)
	if serviceError != nil {
		slog.Error("failed to get portal medications", "error", serviceError, "patient_id", patientFHIRID, "request_id", middleware.GetRequestID(httpRequest.Context()))
		render.Error(httpResponseWriter, http.StatusInternalServerError, "Erro ao carregar medicamentos.")
		return
	}

	render.JSON(httpResponseWriter, http.StatusOK, medications)
}

func (handler *HTTPHandler) GetReports(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
	patientFHIRID := handler.extractPatientFHIRID(httpResponseWriter, httpRequest)
	if patientFHIRID == "" {
		return
	}

	reports, serviceError := handler.service.GetReports(httpRequest.Context(), patientFHIRID)
	if serviceError != nil {
		slog.Error("failed to get portal reports", "error", serviceError, "patient_id", patientFHIRID, "request_id", middleware.GetRequestID(httpRequest.Context()))
		render.Error(httpResponseWriter, http.StatusInternalServerError, "Erro ao carregar exames.")
		return
	}

	render.JSON(httpResponseWriter, http.StatusOK, reports)
}

func (handler *HTTPHandler) GetImaging(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
	patientFHIRID := handler.extractPatientFHIRID(httpResponseWriter, httpRequest)
	if patientFHIRID == "" {
		return
	}

	imaging, serviceError := handler.service.GetImaging(httpRequest.Context(), patientFHIRID)
	if serviceError != nil {
		slog.Error("failed to get portal imaging", "error", serviceError, "patient_id", patientFHIRID, "request_id", middleware.GetRequestID(httpRequest.Context()))
		render.Error(httpResponseWriter, http.StatusInternalServerError, "Erro ao carregar imagens.")
		return
	}

	render.JSON(httpResponseWriter, http.StatusOK, imaging)
}
