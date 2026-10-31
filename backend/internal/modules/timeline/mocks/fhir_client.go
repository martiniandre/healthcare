package mocks

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
)

var ErrUnsupportedOperation = errors.New("operation not supported by mock fhir client")

type SearchCall struct {
	ResourceType string
	QueryParams  string
}

type MockFHIRClient struct {
	GetResourceFn         func(ctx context.Context, resourceType string, resourceID string) (json.RawMessage, error)
	SearchResourcesPageFn func(ctx context.Context, resourceType string, queryParams string) (json.RawMessage, error)

	mutex       sync.Mutex
	searchCalls []SearchCall
}

func (mockClient *MockFHIRClient) SearchCalls() []SearchCall {
	mockClient.mutex.Lock()
	defer mockClient.mutex.Unlock()
	return append([]SearchCall(nil), mockClient.searchCalls...)
}

func (mockClient *MockFHIRClient) GetResource(ctx context.Context, resourceType string, resourceID string) (json.RawMessage, error) {
	if mockClient.GetResourceFn != nil {
		return mockClient.GetResourceFn(ctx, resourceType, resourceID)
	}
	return json.RawMessage(`{}`), nil
}

func (mockClient *MockFHIRClient) SearchResourcesPage(ctx context.Context, resourceType string, queryParams string) (json.RawMessage, error) {
	mockClient.mutex.Lock()
	mockClient.searchCalls = append(mockClient.searchCalls, SearchCall{ResourceType: resourceType, QueryParams: queryParams})
	mockClient.mutex.Unlock()
	if mockClient.SearchResourcesPageFn != nil {
		return mockClient.SearchResourcesPageFn(ctx, resourceType, queryParams)
	}
	return json.RawMessage(`{"resourceType":"Bundle","entry":[]}`), nil
}

func (mockClient *MockFHIRClient) CreateResource(ctx context.Context, resourceType string, resourceBody interface{}) (json.RawMessage, error) {
	return nil, ErrUnsupportedOperation
}

func (mockClient *MockFHIRClient) SearchResources(ctx context.Context, resourceType string, queryParams string) (json.RawMessage, error) {
	return nil, ErrUnsupportedOperation
}

func (mockClient *MockFHIRClient) UpdateResource(ctx context.Context, resourceType string, resourceID string, resourceBody interface{}) (json.RawMessage, error) {
	return nil, ErrUnsupportedOperation
}

func (mockClient *MockFHIRClient) DeleteResource(ctx context.Context, fhirResourcePath string) error {
	return ErrUnsupportedOperation
}
