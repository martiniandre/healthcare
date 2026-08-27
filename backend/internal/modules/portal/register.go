package portal

import (
	"github.com/healthcare/backend/internal/shared/healthcare"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Dependency struct {
	FHIRClient healthcare.FHIRClient
	DB         *pgxpool.Pool
}

func Register(dep Dependency) *HTTPHandler {
	repo := NewRepository(dep.FHIRClient, dep.DB)
	svc := NewService(repo)
	return NewHTTPHandler(svc)
}
