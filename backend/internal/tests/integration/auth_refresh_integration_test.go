package integration

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"testing"
)

func readRefreshCookieValue(t *testing.T, authenticated *authenticatedClient) string {
	t.Helper()
	parsedURL, parseError := url.Parse(authenticated.serverURL + "/api/v1/auth/refresh")
	if parseError != nil {
		t.Fatalf("failed to parse url: %v", parseError)
	}
	for _, cookie := range authenticated.httpClient.Jar.Cookies(parsedURL) {
		if cookie.Name == "refresh_token" {
			return cookie.Value
		}
	}
	t.Fatal("expected refresh_token cookie after login")
	return ""
}

func performRefreshWithCookie(t *testing.T, serverURL string, refreshToken string, csrfToken string) *http.Response {
	t.Helper()
	request, requestError := http.NewRequest(http.MethodPost, serverURL+"/api/v1/auth/refresh", nil)
	if requestError != nil {
		t.Fatalf("failed to build refresh request: %v", requestError)
	}
	request.AddCookie(&http.Cookie{Name: "refresh_token", Value: refreshToken})
	request.AddCookie(&http.Cookie{Name: "csrf_token", Value: csrfToken})
	request.Header.Set("X-CSRF-Token", csrfToken)
	response, doError := http.DefaultClient.Do(request)
	if doError != nil {
		t.Fatalf("refresh request failed: %v", doError)
	}
	return response
}

func TestRefreshRotatesSessionAndKeepsAccessWorking(t *testing.T) {
	testServer := newTestServer(t)
	serverURL := startTestHTTPServer(t, testServer.handler)

	adminClient := loginAs(t, serverURL, "admin@hospital.com", "secret123")
	originalRefreshToken := readRefreshCookieValue(t, adminClient)

	refreshResponse := adminClient.Post(t, "/api/v1/auth/refresh", nil)
	requireStatusCode(t, refreshResponse, http.StatusOK)
	var refreshBody struct {
		Role      string `json:"role"`
		CSRFToken string `json:"csrfToken"`
	}
	if decodeError := json.NewDecoder(refreshResponse.Body).Decode(&refreshBody); decodeError != nil {
		t.Fatalf("failed to decode refresh response: %v", decodeError)
	}
	refreshResponse.Body.Close()

	if refreshBody.Role != "ADMIN" || refreshBody.CSRFToken == "" {
		t.Fatalf("unexpected refresh body: %+v", refreshBody)
	}
	if rotatedRefreshToken := readRefreshCookieValue(t, adminClient); rotatedRefreshToken == originalRefreshToken {
		t.Fatal("expected refresh token to be rotated")
	}

	requireStatusCode(t, adminClient.Get(t, "/api/v1/auth/me"), http.StatusOK)
}

func TestRefreshTokenReuseRevokesWholeSession(t *testing.T) {
	testServer := newTestServer(t)
	serverURL := startTestHTTPServer(t, testServer.handler)

	adminClient := loginAs(t, serverURL, "admin@hospital.com", "secret123")
	originalRefreshToken := readRefreshCookieValue(t, adminClient)

	requireStatusCode(t, adminClient.Post(t, "/api/v1/auth/refresh", nil), http.StatusOK)

	reuseResponse := performRefreshWithCookie(t, serverURL, originalRefreshToken, "attacker-csrf")
	requireStatusCode(t, reuseResponse, http.StatusUnauthorized)

	legitimateAfterReuse := performRefreshWithCookie(t, serverURL, readRefreshCookieValue(t, adminClient), "attacker-csrf")
	requireStatusCode(t, legitimateAfterReuse, http.StatusUnauthorized)
}

func TestRefreshRequiresCSRFHeader(t *testing.T) {
	testServer := newTestServer(t)
	serverURL := startTestHTTPServer(t, testServer.handler)

	adminClient := loginAs(t, serverURL, "admin@hospital.com", "secret123")

	response := performJSONRequest(t, adminClient.httpClient, serverURL, http.MethodPost, "/api/v1/auth/refresh", "", nil)
	requireStatusCode(t, response, http.StatusForbidden)
}

func TestRefreshRejectsDeactivatedUser(t *testing.T) {
	testServer := newTestServer(t)
	serverURL := startTestHTTPServer(t, testServer.handler)

	doctorClient := loginAs(t, serverURL, "medico@clinica.com", "secret123")

	ctx := context.Background()
	if _, deactivateError := testServer.db.Exec(ctx, "UPDATE users SET is_active = FALSE WHERE email = $1", "medico@clinica.com"); deactivateError != nil {
		t.Fatalf("failed to deactivate user: %v", deactivateError)
	}
	t.Cleanup(func() {
		testServer.db.Exec(ctx, "UPDATE users SET is_active = TRUE WHERE email = $1", "medico@clinica.com")
	})

	requireStatusCode(t, doctorClient.Post(t, "/api/v1/auth/refresh", nil), http.StatusUnauthorized)
}

func TestLogoutRevokesOnlyCurrentSessionFamily(t *testing.T) {
	testServer := newTestServer(t)
	serverURL := startTestHTTPServer(t, testServer.handler)

	firstDevice := loginAs(t, serverURL, "admin@hospital.com", "secret123")
	secondDevice := loginAs(t, serverURL, "admin@hospital.com", "secret123")
	firstDeviceRefreshToken := readRefreshCookieValue(t, firstDevice)

	requireStatusCode(t, firstDevice.Post(t, "/api/v1/auth/logout", nil), http.StatusOK)

	requireStatusCode(t, performRefreshWithCookie(t, serverURL, firstDeviceRefreshToken, "x"), http.StatusUnauthorized)
	requireStatusCode(t, secondDevice.Post(t, "/api/v1/auth/refresh", nil), http.StatusOK)
}
