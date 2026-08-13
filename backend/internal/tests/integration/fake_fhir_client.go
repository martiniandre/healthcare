package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"sync"

	"github.com/healthcare/backend/internal/shared/healthcare"
)

type inMemoryFHIRClient struct {
	mu              sync.Mutex
	resources       map[string]map[string]json.RawMessage
	resourceCounter int
}

func newInMemoryFHIRClient() *inMemoryFHIRClient {
	return &inMemoryFHIRClient{
		resources: make(map[string]map[string]json.RawMessage),
	}
}

func (fake *inMemoryFHIRClient) reset() {
	fake.mu.Lock()
	defer fake.mu.Unlock()
	fake.resources = make(map[string]map[string]json.RawMessage)
	fake.resourceCounter = 0
}

func (fake *inMemoryFHIRClient) CreateResource(ctx context.Context, resourceType string, resourceBody interface{}) (json.RawMessage, error) {
	fake.mu.Lock()
	defer fake.mu.Unlock()

	rawBody, err := json.Marshal(resourceBody)
	if err != nil {
		return nil, err
	}

	var resource map[string]interface{}
	if err := json.Unmarshal(rawBody, &resource); err != nil {
		return nil, err
	}

	resourceID, hasID := resource["id"].(string)
	if !hasID || resourceID == "" {
		fake.resourceCounter++
		resourceID = fmt.Sprintf("fhir-%d", fake.resourceCounter)
		resource["id"] = resourceID
	}
	resource["resourceType"] = resourceType
	resource["meta"] = map[string]string{"versionId": "1"}

	rawWithID, err := json.Marshal(resource)
	if err != nil {
		return nil, err
	}

	if fake.resources[resourceType] == nil {
		fake.resources[resourceType] = make(map[string]json.RawMessage)
	}
	fake.resources[resourceType][resourceID] = rawWithID

	return rawWithID, nil
}

func (fake *inMemoryFHIRClient) GetResource(ctx context.Context, resourceType, resourceID string) (json.RawMessage, error) {
	fake.mu.Lock()
	defer fake.mu.Unlock()

	if fake.resources[resourceType] == nil {
		return nil, &healthcare.NotFoundError{ResourceType: resourceType, ResourceID: resourceID}
	}
	rawResource, exists := fake.resources[resourceType][resourceID]
	if !exists {
		return nil, &healthcare.NotFoundError{ResourceType: resourceType, ResourceID: resourceID}
	}
	return rawResource, nil
}

func (fake *inMemoryFHIRClient) SearchResources(ctx context.Context, resourceType, queryParams string) (json.RawMessage, error) {
	fake.mu.Lock()
	defer fake.mu.Unlock()

	parsedQuery, _ := url.ParseQuery(queryParams)
	identifierFilter := parsedQuery.Get("identifier")

	entries := make([]map[string]json.RawMessage, 0)
	for _, rawResource := range fake.resources[resourceType] {
		if identifierFilter != "" && !fake.resourceMatchesIdentifier(rawResource, identifierFilter) {
			continue
		}
		entries = append(entries, map[string]json.RawMessage{"resource": rawResource})
	}

	bundle := map[string]interface{}{
		"resourceType": "Bundle",
		"type":         "searchset",
		"total":        len(entries),
		"entry":        entries,
	}
	return json.Marshal(bundle)
}

func (fake *inMemoryFHIRClient) UpdateResource(ctx context.Context, resourceType, resourceID string, resourceBody interface{}) (json.RawMessage, error) {
	fake.mu.Lock()
	defer fake.mu.Unlock()

	if fake.resources[resourceType] == nil {
		return nil, &healthcare.NotFoundError{ResourceType: resourceType, ResourceID: resourceID}
	}
	if _, exists := fake.resources[resourceType][resourceID]; !exists {
		return nil, &healthcare.NotFoundError{ResourceType: resourceType, ResourceID: resourceID}
	}

	rawBody, err := json.Marshal(resourceBody)
	if err != nil {
		return nil, err
	}
	var resource map[string]interface{}
	if err := json.Unmarshal(rawBody, &resource); err != nil {
		return nil, err
	}
	resource["id"] = resourceID
	resource["resourceType"] = resourceType
	resource["meta"] = map[string]string{"versionId": fake.nextVersionID(resourceType, resourceID)}

	rawWithID, err := json.Marshal(resource)
	if err != nil {
		return nil, err
	}
	fake.resources[resourceType][resourceID] = rawWithID
	return rawWithID, nil
}

func (fake *inMemoryFHIRClient) nextVersionID(resourceType, resourceID string) string {
	existingRaw, exists := fake.resources[resourceType][resourceID]
	currentVersion := int64(1)
	if exists {
		var existingResource map[string]interface{}
		if err := json.Unmarshal(existingRaw, &existingResource); err == nil {
			if existingMeta, ok := existingResource["meta"].(map[string]interface{}); ok {
				if rawVersion, ok := existingMeta["versionId"].(string); ok {
					if parsedVersion, parseErr := strconv.ParseInt(rawVersion, 10, 64); parseErr == nil {
						currentVersion = parsedVersion + 1
					}
				}
			}
		}
	}
	return strconv.FormatInt(currentVersion, 10)
}

func (fake *inMemoryFHIRClient) DeleteResource(ctx context.Context, fhirResourcePath string) error {
	fake.mu.Lock()
	defer fake.mu.Unlock()

	resourceParts := strings.SplitN(fhirResourcePath, "/", 2)
	if len(resourceParts) != 2 {
		return nil
	}
	resourceType := resourceParts[0]
	resourceID := resourceParts[1]

	if fake.resources[resourceType] != nil {
		delete(fake.resources[resourceType], resourceID)
	}
	return nil
}

func (fake *inMemoryFHIRClient) resourceMatchesIdentifier(rawResource json.RawMessage, identifierFilter string) bool {
	var resource struct {
		Identifier []struct {
			System string `json:"system"`
			Value  string `json:"value"`
		} `json:"identifier"`
	}
	if err := json.Unmarshal(rawResource, &resource); err != nil {
		return false
	}
	for _, identifier := range resource.Identifier {
		if identifier.System != "" && identifier.Value != "" {
			if identifier.System+"|"+identifier.Value == identifierFilter {
				return true
			}
		}
	}
	return false
}
