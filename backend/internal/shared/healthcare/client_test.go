package healthcare

import "testing"

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
