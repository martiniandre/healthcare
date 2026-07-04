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
