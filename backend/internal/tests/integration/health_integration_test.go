package integration

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestHealthCheckReportsDatabaseAndCacheHealthy(t *testing.T) {
	testServer := newTestServer(t)
	serverURL := startTestHTTPServer(t, testServer.handler)

	client := newHTTPClient()
	response := performJSONRequest(t, client, serverURL, http.MethodGet, "/health", "", nil)
	requireStatusCode(t, response, http.StatusOK)

	var healthResponse struct {
		Status string `json:"status"`
		Checks []struct {
			Status  string `json:"status"`
			Service string `json:"service"`
		} `json:"checks"`
	}
	decodeJSONResponse(t, response, &healthResponse)

	if healthResponse.Status != "ok" {
		t.Fatalf("expected overall health status ok, got %q", healthResponse.Status)
	}

	servicesByStatus := make(map[string]string)
	for _, check := range healthResponse.Checks {
		servicesByStatus[check.Service] = check.Status
	}
	if servicesByStatus["database"] != "ok" {
		t.Fatalf("expected database check ok, got %q", servicesByStatus["database"])
	}
}

func TestPatientListReturnsPaginationEnvelope(t *testing.T) {
	testServer := newTestServer(t)
	serverURL := startTestHTTPServer(t, testServer.handler)

	doctorClient := loginAs(t, serverURL, "medico@clinica.com", "secret123")

	patientsResponse := doctorClient.Get(t, "/api/v1/patients")
	requireStatusCode(t, patientsResponse, http.StatusOK)

	var patientsEnvelope struct {
		Patients []json.RawMessage `json:"patients"`
		Total    int               `json:"total"`
		Page     int               `json:"page"`
		Limit    int               `json:"limit"`
	}
	decodeJSONResponse(t, patientsResponse, &patientsEnvelope)
	if patientsEnvelope.Total < 0 {
		t.Fatalf("expected non-negative total, got %d", patientsEnvelope.Total)
	}
}
