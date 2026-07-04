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
	httpServeMux.Handle("GET /api/analytics", authorizedRoles(http.HandlerFunc(analyticsHTTPHandler.GetStats)))
}

func (analyticsHTTPHandler *HTTPHandler) GetStats(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
	analyticsData, errorInstance := analyticsHTTPHandler.analyticsService.GetStats(httpRequest.Context())
	if errorInstance != nil {
		slog.Error("failed to get analytics", "error", errorInstance, "request_id", middleware.GetRequestID(httpRequest.Context()))
		render.Error(httpResponseWriter, http.StatusInternalServerError, "Erro ao obter estatísticas.")
		return
	}

	render.JSON(httpResponseWriter, http.StatusOK, analyticsData)
}
