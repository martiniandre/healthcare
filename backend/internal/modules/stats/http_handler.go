package stats

import (
	"log/slog"
	"net/http"

	"github.com/healthcare/backend/internal/api/middleware"
	"github.com/healthcare/backend/internal/api/render"
	"github.com/healthcare/backend/internal/modules/auth"
)

type HTTPHandler struct {
	statsService Service
}

func NewHTTPHandler(statsService Service) *HTTPHandler {
	return &HTTPHandler{
		statsService: statsService,
	}
}

func (statsHTTPHandler *HTTPHandler) RegisterRoutes(httpServeMux *http.ServeMux) {
	authorizedRoles := middleware.RequireRoles(auth.RoleAdmin, auth.RoleDoctor, auth.RoleNurse)
	httpServeMux.Handle("GET /api/stats", authorizedRoles(http.HandlerFunc(statsHTTPHandler.GetStats)))
}

func (statsHTTPHandler *HTTPHandler) GetStats(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
	statsData, errorInstance := statsHTTPHandler.statsService.GetStats(httpRequest.Context())
	if errorInstance != nil {
		slog.Error("failed to get stats", "error", errorInstance, "request_id", middleware.GetRequestID(httpRequest.Context()))
		render.Error(httpResponseWriter, http.StatusInternalServerError, "Erro ao obter estatísticas.")
		return
	}

	render.JSON(httpResponseWriter, http.StatusOK, statsData)
}
