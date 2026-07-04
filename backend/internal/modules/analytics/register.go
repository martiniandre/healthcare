package analytics

import (
	"github.com/healthcare/backend/internal/shared/healthcare"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Dependency struct {
	DB         *pgxpool.Pool
	FHIRClient healthcare.FHIRClient
}

func Register(dep Dependency) *HTTPHandler {
	analyticsRepository := NewRepository(dep.DB, dep.FHIRClient)
	analyticsService := NewService(analyticsRepository)
	return NewHTTPHandler(analyticsService)
}
