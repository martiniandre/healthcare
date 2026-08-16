package healthcare

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewClientWithBaseURLBuildsClient(t *testing.T) {
	client, clientError := NewClientWithBaseURL("http://localhost:8081/fhir/")
	if clientError != nil {
		t.Fatalf("unexpected error: %v", clientError)
	}
	if client.baseURL != "http://localhost:8081/fhir" {
		t.Fatalf("expected trimmed base url, got %q", client.baseURL)
	}
	if client.httpClient == nil {
		t.Fatal("expected http client to be set")
	}
}

func TestNewClientWithBaseURLRejectsRelativeURL(t *testing.T) {
	_, clientError := NewClientWithBaseURL("/fhir")
	if clientError == nil {
		t.Fatal("expected error for relative url")
	}
}

func TestNewClientWithBaseURLRejectsMalformedURL(t *testing.T) {
	_, clientError := NewClientWithBaseURL("http://")
	if clientError == nil {
		t.Fatal("expected error for malformed url")
	}
}

func TestSearchResourcesAccumulatesEntriesAcrossAllPages(t *testing.T) {
	var paginationServer *httptest.Server
	paginationServer = httptest.NewServer(http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		responseWriter.Header().Set("Content-Type", "application/fhir+json")
		switch request.URL.Path {
		case "/fhir/Observation":
			responseWriter.Write([]byte(fmt.Sprintf(`{
				"resourceType": "Bundle",
				"type": "searchset",
				"total": 3,
				"entry": [{"resource": {"id": "obs-1"}}],
				"link": [{"relation": "next", "url": "%s/page-2"}]
			}`, paginationServer.URL)))
		case "/page-2":
			responseWriter.Write([]byte(fmt.Sprintf(`{
				"resourceType": "Bundle",
				"type": "searchset",
				"total": 3,
				"entry": [{"resource": {"id": "obs-2"}}],
				"link": [{"relation": "next", "url": "%s/page-3"}]
			}`, paginationServer.URL)))
		case "/page-3":
			responseWriter.Write([]byte(`{
				"resourceType": "Bundle",
				"type": "searchset",
				"total": 3,
				"entry": [{"resource": {"id": "obs-3"}}]
			}`))
		default:
			http.NotFound(responseWriter, request)
		}
	}))
	defer paginationServer.Close()

	client, clientError := NewClientWithBaseURL(paginationServer.URL + "/fhir")
	if clientError != nil {
		t.Fatalf("unexpected error: %v", clientError)
	}

	searchResult, searchError := client.SearchResources(context.Background(), "Observation", "encounter=Encounter%2Fenc-9")
	if searchError != nil {
		t.Fatalf("unexpected search error: %v", searchError)
	}

	var mergedBundle struct {
		Entry []struct {
			Resource map[string]interface{} `json:"resource"`
		} `json:"entry"`
	}
	if parseError := json.Unmarshal(searchResult, &mergedBundle); parseError != nil {
		t.Fatalf("failed to parse merged bundle: %v", parseError)
	}
	if len(mergedBundle.Entry) != 3 {
		t.Fatalf("expected 3 entries across all pages, got %d", len(mergedBundle.Entry))
	}

	observedIDs := make(map[string]bool)
	for _, entry := range mergedBundle.Entry {
		resourceID, _ := entry.Resource["id"].(string)
		observedIDs[resourceID] = true
	}
	for _, expectedID := range []string{"obs-1", "obs-2", "obs-3"} {
		if !observedIDs[expectedID] {
			t.Fatalf("expected entry %q to be present in merged bundle", expectedID)
		}
	}
}
