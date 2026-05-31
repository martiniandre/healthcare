package stats

import (
	"github.com/healthcare/backend/internal/shared/healthcare"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Register(databasePool *pgxpool.Pool, fhirClient healthcare.FHIRClient) (Service, *HTTPHandler) {
	statsRepository := NewRepository(databasePool, fhirClient)
	statsService := NewService(statsRepository)
	statsHTTPHandler := NewHTTPHandler(statsService)
	return statsService, statsHTTPHandler
}
