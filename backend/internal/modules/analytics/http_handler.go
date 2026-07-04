package analytics

import (
	"log/slog"
	"net/http"

	"github.com/healthcare/backend/internal/api/middleware"
	"github.com/healthcare/backend/internal/api/render"
	"github.com/healthcare/backend/internal/shared/role"
)

type HTTPHandler struct {
	analyticsService Service
}

func NewHTTPHandler(analyticsService Service) *HTTPHandler {
	return &HTTPHandler{
		analyticsService: analyticsService,
	}
}

func (analyticsHTTPHandler *HTTPHandler) RegisterRoutes(httpServeMux *http.ServeMux) {
	authorizedRoles := middleware.RequireRoles(role.RoleAdmin, role.RoleDoctor, role.RoleNurse)
	httpServeMux.Handle("GET /api/v1/analytics", authorizedRoles(http.HandlerFunc(analyticsHTTPHandler.GetStats)))
	httpServeMux.Handle("GET /api/v1/analytics/dashboard", authorizedRoles(http.HandlerFunc(analyticsHTTPHandler.GetDashboard)))
	httpServeMux.Handle("GET /api/v1/analytics/dashboard/consultations-per-doctor", authorizedRoles(http.HandlerFunc(analyticsHTTPHandler.GetConsultationsPerDoctor)))
	httpServeMux.Handle("GET /api/v1/analytics/dashboard/occupancy-rate", authorizedRoles(http.HandlerFunc(analyticsHTTPHandler.GetOccupancyRate)))
	httpServeMux.Handle("GET /api/v1/analytics/dashboard/avg-wait-time", authorizedRoles(http.HandlerFunc(analyticsHTTPHandler.GetAvgWaitTime)))
	httpServeMux.Handle("GET /api/v1/analytics/dashboard/top-diagnoses", authorizedRoles(http.HandlerFunc(analyticsHTTPHandler.GetTopDiagnoses)))
}

// GetStats godoc
//
//	@Summary		Get analytics statistics
//	@Description	Returns aggregated healthcare statistics including patient counts, consultation data, and exam modality distribution
//	@Tags			analytics
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	Stats
//	@Failure		500	{object}	map[string]string
//	@Router			/analytics [get]
func (analyticsHTTPHandler *HTTPHandler) GetStats(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
	analyticsData, errorInstance := analyticsHTTPHandler.analyticsService.GetStats(httpRequest.Context())
	if errorInstance != nil {
		slog.Error("failed to get analytics", "error", errorInstance, "request_id", middleware.GetRequestID(httpRequest.Context()))
		render.Error(httpResponseWriter, http.StatusInternalServerError, "Erro ao obter estatísticas.")
		return
	}

	render.JSON(httpResponseWriter, http.StatusOK, analyticsData)
}

// GetDashboard godoc
//
//	@Summary		Get clinical dashboard data
//	@Description	Returns aggregated KPIs for the clinical dashboard including consultations, occupancy, wait times, and diagnoses
//	@Tags			analytics
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	DashboardData
//	@Failure		500	{object}	map[string]string
//	@Router			/analytics/dashboard [get]
func (analyticsHTTPHandler *HTTPHandler) GetDashboard(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
	dashboardData, errorInstance := analyticsHTTPHandler.analyticsService.GetDashboardData(httpRequest.Context())
	if errorInstance != nil {
		slog.Error("failed to get dashboard data", "error", errorInstance, "request_id", middleware.GetRequestID(httpRequest.Context()))
		render.Error(httpResponseWriter, http.StatusInternalServerError, "Erro ao obter dados do dashboard.")
		return
	}

	render.JSON(httpResponseWriter, http.StatusOK, dashboardData)
}

// GetConsultationsPerDoctor godoc
//
//	@Summary		Get consultations per doctor
//	@Description	Returns consultations grouped by doctor for the dashboard
//	@Tags			analytics
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}	DoctorConsultation
//	@Failure		500	{object}	map[string]string
//	@Router			/analytics/dashboard/consultations-per-doctor [get]
func (analyticsHTTPHandler *HTTPHandler) GetConsultationsPerDoctor(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
	consultationsPerDoctor, errorInstance := analyticsHTTPHandler.analyticsService.GetConsultationsPerDoctor(httpRequest.Context())
	if errorInstance != nil {
		slog.Error("failed to get consultations per doctor", "error", errorInstance, "request_id", middleware.GetRequestID(httpRequest.Context()))
		render.Error(httpResponseWriter, http.StatusInternalServerError, "Erro ao obter consultas por médico.")
		return
	}

	render.JSON(httpResponseWriter, http.StatusOK, consultationsPerDoctor)
}

// GetOccupancyRate godoc
//
//	@Summary		Get occupancy rate
//	@Description	Returns bed/room occupancy rate for the dashboard
//	@Tags			analytics
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	OccupancyRate
//	@Failure		500	{object}	map[string]string
//	@Router			/analytics/dashboard/occupancy-rate [get]
func (analyticsHTTPHandler *HTTPHandler) GetOccupancyRate(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
	occupancyRate, errorInstance := analyticsHTTPHandler.analyticsService.GetOccupancyRate(httpRequest.Context())
	if errorInstance != nil {
		slog.Error("failed to get occupancy rate", "error", errorInstance, "request_id", middleware.GetRequestID(httpRequest.Context()))
		render.Error(httpResponseWriter, http.StatusInternalServerError, "Erro ao obter taxa de ocupação.")
		return
	}

	render.JSON(httpResponseWriter, http.StatusOK, occupancyRate)
}

// GetAvgWaitTime godoc
//
//	@Summary		Get average wait time
//	@Description	Returns average wait time by department for the dashboard
//	@Tags			analytics
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	AvgWaitTime
//	@Failure		500	{object}	map[string]string
//	@Router			/analytics/dashboard/avg-wait-time [get]
func (analyticsHTTPHandler *HTTPHandler) GetAvgWaitTime(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
	avgWaitTime, errorInstance := analyticsHTTPHandler.analyticsService.GetAvgWaitTime(httpRequest.Context())
	if errorInstance != nil {
		slog.Error("failed to get average wait time", "error", errorInstance, "request_id", middleware.GetRequestID(httpRequest.Context()))
		render.Error(httpResponseWriter, http.StatusInternalServerError, "Erro ao obter tempo médio de espera.")
		return
	}

	render.JSON(httpResponseWriter, http.StatusOK, avgWaitTime)
}

// GetTopDiagnoses godoc
//
//	@Summary		Get top diagnoses
//	@Description	Returns the most common diagnoses for the dashboard
//	@Tags			analytics
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}	DiagnosisCount
//	@Failure		500	{object}	map[string]string
//	@Router			/analytics/dashboard/top-diagnoses [get]
func (analyticsHTTPHandler *HTTPHandler) GetTopDiagnoses(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
	topDiagnoses, errorInstance := analyticsHTTPHandler.analyticsService.GetTopDiagnoses(httpRequest.Context())
	if errorInstance != nil {
		slog.Error("failed to get top diagnoses", "error", errorInstance, "request_id", middleware.GetRequestID(httpRequest.Context()))
		render.Error(httpResponseWriter, http.StatusInternalServerError, "Erro ao obter principais diagnósticos.")
		return
	}

	render.JSON(httpResponseWriter, http.StatusOK, topDiagnoses)
}
