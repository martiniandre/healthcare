package schedule

import (
	"net/http"

	"github.com/healthcare/backend/internal/modules/audit_logs"
	"github.com/healthcare/backend/internal/shared/eventbus"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Dependency struct {
	DB           *pgxpool.Pool
	EventBus     eventbus.Bus
	AuditService audit_logs.Service
}

type RegisterHandler struct {
	appointmentHandler    *HTTPHandler
	unavailabilityHandler *UnavailabilityHTTPHandler
}

func (registerHandler *RegisterHandler) RegisterRoutes(mux *http.ServeMux) {
	registerHandler.appointmentHandler.RegisterRoutes(mux)
	registerHandler.unavailabilityHandler.RegisterRoutes(mux)
}

func Register(dep Dependency) *RegisterHandler {
	scheduleRepository := NewRepository(dep.DB)
	scheduleService := NewService(scheduleRepository, dep.EventBus, dep.AuditService)

	unavailabilityRepository := NewUnavailabilityRepository(dep.DB)
	unavailabilityService := NewUnavailabilityService(unavailabilityRepository)

	return &RegisterHandler{
		appointmentHandler:    NewHTTPHandler(scheduleService),
		unavailabilityHandler: NewUnavailabilityHTTPHandler(unavailabilityService),
	}
}
