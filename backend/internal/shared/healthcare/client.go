package healthcare

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/oauth2/google"
)

const maxHealthcareErrorBodyBytes int64 = 64 << 10
const defaultHTTPTimeout = 30 * time.Second
const maxRetryAttempts = 3
const maxPaginationPages = 10

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

	googleHTTPClient.Timeout = defaultHTTPTimeout

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

	endpoint := fmt.Sprintf("%s/%s", healthcareClient.baseURL, url.PathEscape(resourceType))
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	request.Header.Set("Content-Type", "application/fhir+json")

	return healthcareClient.doWithRetry(request)
}

func (healthcareClient *Client) GetResource(ctx context.Context, resourceType, resourceID string) (json.RawMessage, error) {
	endpoint := fmt.Sprintf("%s/%s/%s", healthcareClient.baseURL, url.PathEscape(resourceType), url.PathEscape(resourceID))
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	request.Header.Set("Accept", "application/fhir+json")

	return healthcareClient.doWithRetry(request)
}

func (healthcareClient *Client) SearchResources(ctx context.Context, resourceType, queryParams string) (json.RawMessage, error) {
	endpoint := fmt.Sprintf("%s/%s?%s", healthcareClient.baseURL, url.PathEscape(resourceType), queryParams)
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	request.Header.Set("Accept", "application/fhir+json")

	bundleResponse, bundleError := healthcareClient.doWithRetry(request)
	if bundleError != nil {
		return nil, bundleError
	}

	var bundle map[string]interface{}
	if parseError := json.Unmarshal(bundleResponse, &bundle); parseError != nil {
		return nil, fmt.Errorf("failed to parse search bundle: %w", parseError)
	}

	for pageCount := 0; pageCount < maxPaginationPages; pageCount++ {
		links, hasLinks := bundle["link"].([]interface{})
		if !hasLinks {
			break
		}

		var nextPageURL string
		for _, rawLink := range links {
			linkEntry, ok := rawLink.(map[string]interface{})
			if !ok {
				continue
			}
			if relation, ok := linkEntry["relation"].(string); ok && relation == "next" {
				nextPageURL, _ = linkEntry["url"].(string)
				break
			}
		}

		if nextPageURL == "" {
			break
		}

		nextRequest, requestError := http.NewRequestWithContext(ctx, http.MethodGet, nextPageURL, nil)
		if requestError != nil {
			return nil, fmt.Errorf("failed to create pagination request: %w", requestError)
		}
		nextRequest.Header.Set("Accept", "application/fhir+json")

		nextResponse, nextError := healthcareClient.doWithRetry(nextRequest)
		if nextError != nil {
			return nil, nextError
		}

		var nextBundle map[string]interface{}
		if parseError := json.Unmarshal(nextResponse, &nextBundle); parseError != nil {
			return nil, fmt.Errorf("failed to parse next page bundle: %w", parseError)
		}

		if nextEntries, hasEntries := nextBundle["entry"].([]interface{}); hasEntries {
			currentEntries, _ := bundle["entry"].([]interface{})
			bundle["entry"] = append(currentEntries, nextEntries...)
		}

		bundle = nextBundle
	}

	mergedResponse, marshalError := json.Marshal(bundle)
	if marshalError != nil {
		return nil, fmt.Errorf("failed to marshal merged bundle: %w", marshalError)
	}

	return mergedResponse, nil
}

func (healthcareClient *Client) UpdateResource(ctx context.Context, resourceType, resourceID string, resourceBody interface{}) (json.RawMessage, error) {
	bodyBytes, err := json.Marshal(resourceBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal resource: %w", err)
	}

	endpoint := fmt.Sprintf("%s/%s/%s", healthcareClient.baseURL, url.PathEscape(resourceType), url.PathEscape(resourceID))
	request, err := http.NewRequestWithContext(ctx, http.MethodPut, endpoint, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	request.Header.Set("Content-Type", "application/fhir+json")

	return healthcareClient.doWithRetry(request)
}

func (healthcareClient *Client) DeleteResource(ctx context.Context, fhirResourcePath string) error {
	endpoint := fmt.Sprintf("%s/%s", healthcareClient.baseURL, fhirResourcePath)
	request, err := http.NewRequestWithContext(ctx, http.MethodDelete, endpoint, nil)
	if err != nil {
		return fmt.Errorf("failed to create delete request: %w", err)
	}
	request.Header.Set("Accept", "application/fhir+json")

	response, err := healthcareClient.doWithRetry(request)
	if err != nil {
		return err
	}

	_ = response
	return nil
}

func (healthcareClient *Client) doWithRetry(request *http.Request) (json.RawMessage, error) {
	retryDelays := []time.Duration{100 * time.Millisecond, 300 * time.Millisecond, 900 * time.Millisecond}

	var lastError error
	for attemptIndex := 0; attemptIndex <= maxRetryAttempts; attemptIndex++ {
		response, requestError := healthcareClient.httpClient.Do(request)
		if requestError != nil {
			lastError = fmt.Errorf("request failed: %w", requestError)
			if attemptIndex < maxRetryAttempts {
				time.Sleep(retryDelays[attemptIndex])
				continue
			}
			return nil, lastError
		}

		var responseReader io.Reader = response.Body
		if response.StatusCode < 200 || response.StatusCode >= 300 {
			responseReader = io.LimitReader(response.Body, maxHealthcareErrorBodyBytes)
		}

		responseBody, readError := io.ReadAll(responseReader)
		response.Body.Close()

		if readError != nil {
			lastError = fmt.Errorf("failed to read response body: %w", readError)
			return nil, lastError
		}

		if response.StatusCode == http.StatusTooManyRequests || response.StatusCode == http.StatusBadGateway ||
			response.StatusCode == http.StatusServiceUnavailable || response.StatusCode == http.StatusGatewayTimeout {
			if attemptIndex < maxRetryAttempts {
				lastError = fmt.Errorf("healthcare api: unexpected status %d", response.StatusCode)
				time.Sleep(retryDelays[attemptIndex])
				continue
			}
			return nil, fmt.Errorf("healthcare api: unexpected status %d", response.StatusCode)
		}

		if response.StatusCode < 200 || response.StatusCode >= 300 {
			return nil, fmt.Errorf("healthcare api: unexpected status %d", response.StatusCode)
		}

		return responseBody, nil
	}

	return nil, lastError
}
