package portal

import "github.com/healthcare/backend/internal/shared/healthcare"

type Dependency struct {
	FHIRClient healthcare.FHIRClient
}

func Register(dep Dependency) *HTTPHandler {
	repo := NewRepository(dep.FHIRClient)
	svc := NewService(repo)
	return NewHTTPHandler(svc)
}
