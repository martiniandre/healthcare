package healthcare

import (
	"context"
	"encoding/json"
)

type FHIRClient interface {
	CreateResource(ctx context.Context, resourceType string, resourceBody interface{}) (json.RawMessage, error)
	GetResource(ctx context.Context, resourceType, resourceID string) (json.RawMessage, error)
	SearchResources(ctx context.Context, resourceType, queryParams string) (json.RawMessage, error)
	UpdateResource(ctx context.Context, resourceType, resourceID string, resourceBody interface{}) (json.RawMessage, error)
	DeleteResource(ctx context.Context, fhirResourcePath string) error
}
