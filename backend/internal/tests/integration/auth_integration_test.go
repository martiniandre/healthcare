package integration

import (
	"net/http"
	"testing"
)

func TestLoginWithValidCredentialsReturnsSession(t *testing.T) {
	testServer := newTestServer(t)
	serverURL := startTestHTTPServer(t, testServer.handler)

	adminClient := loginAs(t, serverURL, "admin@hospital.com", "secret123")

	if adminClient.role != "ADMIN" {
		t.Fatalf("expected role ADMIN, got %q", adminClient.role)
	}
	if adminClient.csrfToken == "" {
		t.Fatal("expected csrf token to be set after login")
	}

	requireStatusCode(t, adminClient.Get(t, "/api/v1/auth/me"), http.StatusOK)
}

func TestLoginWithInvalidPasswordReturnsUnauthorized(t *testing.T) {
	testServer := newTestServer(t)
	serverURL := startTestHTTPServer(t, testServer.handler)

	client := newHTTPClient()
	response := performJSONRequest(t, client, serverURL, http.MethodPost, "/api/v1/auth/login", "", map[string]string{
		"email":    "admin@hospital.com",
		"password": "wrong-password",
	})
	requireStatusCode(t, response, http.StatusUnauthorized)
}

func TestUnauthenticatedRequestIsRejected(t *testing.T) {
	testServer := newTestServer(t)
	serverURL := startTestHTTPServer(t, testServer.handler)

	client := newHTTPClient()
	response := performJSONRequest(t, client, serverURL, http.MethodGet, "/api/v1/patients", "", nil)
	requireStatusCode(t, response, http.StatusUnauthorized)
}

func TestWriteRequestsRequireCSRFHeader(t *testing.T) {
	testServer := newTestServer(t)
	serverURL := startTestHTTPServer(t, testServer.handler)

	adminClient := loginAs(t, serverURL, "admin@hospital.com", "secret123")

	response := performJSONRequest(t, adminClient.httpClient, serverURL, http.MethodPost, "/api/v1/patients", "", map[string]string{
		"full_name":    "Maria Silva",
		"birth_date":   "1990-01-01",
		"document_id":  "529.982.247-25",
		"phone_number": "(11) 98765-4321",
	})
	requireStatusCode(t, response, http.StatusForbidden)
}

func TestLogoutRevokesSession(t *testing.T) {
	testServer := newTestServer(t)
	serverURL := startTestHTTPServer(t, testServer.handler)

	adminClient := loginAs(t, serverURL, "admin@hospital.com", "secret123")

	requireStatusCode(t, adminClient.Post(t, "/api/v1/auth/logout", nil), http.StatusOK)
	requireStatusCode(t, adminClient.Get(t, "/api/v1/auth/me"), http.StatusUnauthorized)
}
