package timeline

import "github.com/healthcare/backend/internal/shared/healthcare"

type Dependency struct {
	FHIRClient healthcare.FHIRClient
}

func Register(dep Dependency) *HTTPHandler {
	timelineRepository := NewRepository(dep.FHIRClient)
	timelineService := NewService(timelineRepository)
	return NewHTTPHandler(timelineService)
}
