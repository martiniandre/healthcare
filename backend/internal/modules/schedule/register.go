package schedule

import (
	"github.com/healthcare/backend/internal/shared/eventbus"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Dependency struct {
	DB       *pgxpool.Pool
	EventBus eventbus.Bus
}

func Register(dep Dependency) *HTTPHandler {
	scheduleRepository := NewRepository(dep.DB)
	scheduleService := NewService(scheduleRepository, dep.EventBus)
	return NewHTTPHandler(scheduleService)
}
