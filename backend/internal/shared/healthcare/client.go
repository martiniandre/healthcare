package healthcare

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/oauth2/google"
)

type Client struct {
	projectID   string
	locationID  string
	datasetID   string
	fhirStoreID string
	httpClient  *http.Client
	baseURL     string
}

func NewClient(ctx context.Context, projectID, locationID, datasetID, fhirStoreID string) (*Client, error) {
	credentialsScopes := []string{"https://www.googleapis.com/auth/cloud-healthcare"}

	googleHTTPClient, err := google.DefaultClient(ctx, credentialsScopes...)
	if err != nil {
		return nil, fmt.Errorf("failed to create google default client: %w", err)
	}

	baseURL := fmt.Sprintf(
		"https://healthcare.googleapis.com/v1/projects/%s/locations/%s/datasets/%s/fhirStores/%s/fhir",
		projectID, locationID, datasetID, fhirStoreID,
	)

	return &Client{
		projectID:   projectID,
		locationID:  locationID,
		datasetID:   datasetID,
		fhirStoreID: fhirStoreID,
		httpClient:  googleHTTPClient,
		baseURL:     baseURL,
	}, nil
}

func (healthcareClient *Client) CreateResource(ctx context.Context, resourceType string, resourceBody interface{}) (json.RawMessage, error) {
	bodyBytes, err := json.Marshal(resourceBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal resource: %w", err)
	}

	endpoint := fmt.Sprintf("%s/%s", healthcareClient.baseURL, resourceType)
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	request.Header.Set("Content-Type", "application/fhir+json")

	return healthcareClient.executeRequest(request)
}

func (healthcareClient *Client) GetResource(ctx context.Context, resourceType, resourceID string) (json.RawMessage, error) {
	endpoint := fmt.Sprintf("%s/%s/%s", healthcareClient.baseURL, resourceType, resourceID)
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	request.Header.Set("Accept", "application/fhir+json")

	return healthcareClient.executeRequest(request)
}

func (healthcareClient *Client) SearchResources(ctx context.Context, resourceType, queryParams string) (json.RawMessage, error) {
	endpoint := fmt.Sprintf("%s/%s?%s", healthcareClient.baseURL, resourceType, queryParams)
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	request.Header.Set("Accept", "application/fhir+json")

	return healthcareClient.executeRequest(request)
}

func (healthcareClient *Client) UpdateResource(ctx context.Context, resourceType, resourceID string, resourceBody interface{}) (json.RawMessage, error) {
	bodyBytes, err := json.Marshal(resourceBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal resource: %w", err)
	}

	endpoint := fmt.Sprintf("%s/%s/%s", healthcareClient.baseURL, resourceType, resourceID)
	request, err := http.NewRequestWithContext(ctx, http.MethodPut, endpoint, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	request.Header.Set("Content-Type", "application/fhir+json")

	return healthcareClient.executeRequest(request)
}

func (healthcareClient *Client) executeRequest(request *http.Request) (json.RawMessage, error) {
	response, err := healthcareClient.httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return nil, fmt.Errorf("healthcare api error %d: %s", response.StatusCode, string(responseBody))
	}

	return responseBody, nil
}
