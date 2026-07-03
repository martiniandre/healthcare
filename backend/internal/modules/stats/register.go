package stats

import (
	"github.com/healthcare/backend/internal/shared/healthcare"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Dependency struct {
	DB         *pgxpool.Pool
	FHIRClient healthcare.FHIRClient
}

func Register(dep Dependency) *HTTPHandler {
	statsRepository := NewRepository(dep.DB, dep.FHIRClient)
	statsService := NewService(statsRepository)
	return NewHTTPHandler(statsService)
}
