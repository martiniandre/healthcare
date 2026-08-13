package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"testing"
)

type authenticatedClient struct {
	httpClient *http.Client
	csrfToken  string
	serverURL  string
	email      string
	role       string
}

func newHTTPClient() *http.Client {
	cookieJar, _ := cookiejar.New(nil)
	return &http.Client{Jar: cookieJar}
}

func startTestHTTPServer(t *testing.T, serverHandler http.Handler) string {
	t.Helper()
	httpServer := httptest.NewServer(serverHandler)
	t.Cleanup(httpServer.Close)
	return httpServer.URL
}

func loginAs(t *testing.T, serverURL string, email string, password string) *authenticatedClient {
	t.Helper()

	httpClient := newHTTPClient()

	loginBody := map[string]string{"email": email, "password": password}
	response := performJSONRequest(t, httpClient, serverURL, http.MethodPost, "/api/v1/auth/login", "", loginBody)
	if response.StatusCode != http.StatusOK {
		t.Fatalf("login as %s failed with status %d", email, response.StatusCode)
	}

	var loginResponse struct {
		Role string `json:"role"`
	}
	if decodeError := json.NewDecoder(response.Body).Decode(&loginResponse); decodeError != nil {
		t.Fatalf("failed to decode login response: %v", decodeError)
	}
	response.Body.Close()

	csrfToken := readCookieValue(t, httpClient, serverURL, "csrf_token")

	return &authenticatedClient{
		httpClient: httpClient,
		csrfToken:  csrfToken,
		serverURL:  serverURL,
		email:      email,
		role:       loginResponse.Role,
	}
}

func (authenticated *authenticatedClient) Get(t *testing.T, path string) *http.Response {
	t.Helper()
	return performJSONRequest(t, authenticated.httpClient, authenticated.serverURL, http.MethodGet, path, "", nil)
}

func (authenticated *authenticatedClient) Post(t *testing.T, path string, payload interface{}) *http.Response {
	t.Helper()
	return performJSONRequest(t, authenticated.httpClient, authenticated.serverURL, http.MethodPost, path, authenticated.csrfToken, payload)
}

func (authenticated *authenticatedClient) Put(t *testing.T, path string, payload interface{}) *http.Response {
	t.Helper()
	return performJSONRequest(t, authenticated.httpClient, authenticated.serverURL, http.MethodPut, path, authenticated.csrfToken, payload)
}

func (authenticated *authenticatedClient) Delete(t *testing.T, path string) *http.Response {
	t.Helper()
	return performJSONRequest(t, authenticated.httpClient, authenticated.serverURL, http.MethodDelete, path, authenticated.csrfToken, nil)
}

func performJSONRequest(t *testing.T, httpClient *http.Client, serverURL string, method string, path string, csrfToken string, payload interface{}) *http.Response {
	t.Helper()

	var requestBody bytes.Buffer
	if payload != nil {
		if marshalError := json.NewEncoder(&requestBody).Encode(payload); marshalError != nil {
			t.Fatalf("failed to encode request payload: %v", marshalError)
		}
	}

	request, requestError := http.NewRequest(method, serverURL+path, &requestBody)
	if requestError != nil {
		t.Fatalf("failed to build request: %v", requestError)
	}
	if payload != nil {
		request.Header.Set("Content-Type", "application/json")
	}
	if csrfToken != "" {
		request.Header.Set("X-CSRF-Token", csrfToken)
	}

	response, doError := httpClient.Do(request)
	if doError != nil {
		t.Fatalf("request %s %s failed: %v", method, path, doError)
	}
	return response
}

func readCookieValue(t *testing.T, httpClient *http.Client, serverURL string, cookieName string) string {
	t.Helper()

	parsedURL, parseError := url.Parse(serverURL)
	if parseError != nil {
		t.Fatalf("failed to parse server url: %v", parseError)
	}
	for _, cookie := range httpClient.Jar.Cookies(parsedURL) {
		if cookie.Name == cookieName {
			return cookie.Value
		}
	}
	return ""
}

func decodeJSONResponse(t *testing.T, response *http.Response, target interface{}) {
	t.Helper()
	defer response.Body.Close()
	if decodeError := json.NewDecoder(response.Body).Decode(target); decodeError != nil {
		t.Fatalf("failed to decode response body: %v", decodeError)
	}
}

func responseStatusCode(response *http.Response) int {
	return response.StatusCode
}

func requireStatusCode(t *testing.T, response *http.Response, expectedStatus int) {
	t.Helper()
	if response.StatusCode != expectedStatus {
		t.Fatalf("expected status %d, got %d", expectedStatus, response.StatusCode)
	}
}

func requireJSONFieldValue(t *testing.T, response *http.Response, fieldName string) string {
	t.Helper()
	var decodedBody map[string]interface{}
	decodeJSONResponse(t, response, &decodedBody)
	fieldValue, exists := decodedBody[fieldName]
	if !exists {
		t.Fatalf("response body missing field %q: %v", fieldName, decodedBody)
	}
	stringValue := fmt.Sprintf("%v", fieldValue)
	if stringValue == "" {
		t.Fatalf("response field %q is empty", fieldName)
	}
	return stringValue
}
