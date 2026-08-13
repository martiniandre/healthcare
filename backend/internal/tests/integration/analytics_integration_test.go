package integration

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestAnalyticsEndpointsReturnStructuredData(t *testing.T) {
	testServer := newTestServer(t)
	serverURL := startTestHTTPServer(t, testServer.handler)

	doctorClient := loginAs(t, serverURL, "medico@clinica.com", "secret123")

	statsResponse := doctorClient.Get(t, "/api/v1/analytics")
	requireStatusCode(t, statsResponse, http.StatusOK)
	var statsBody map[string]interface{}
	decodeJSONResponse(t, statsResponse, &statsBody)

	requireStatusCode(t, doctorClient.Get(t, "/api/v1/analytics/dashboard"), http.StatusOK)

	topDiagnosesResponse := doctorClient.Get(t, "/api/v1/analytics/dashboard/top-diagnoses")
	requireStatusCode(t, topDiagnosesResponse, http.StatusOK)
	var topDiagnoses []json.RawMessage
	decodeJSONResponse(t, topDiagnosesResponse, &topDiagnoses)

	requireStatusCode(t, doctorClient.Get(t, "/api/v1/analytics/dashboard/consultations-per-doctor"), http.StatusOK)
	requireStatusCode(t, doctorClient.Get(t, "/api/v1/analytics/dashboard/occupancy-rate"), http.StatusOK)
	requireStatusCode(t, doctorClient.Get(t, "/api/v1/analytics/dashboard/avg-wait-time"), http.StatusOK)
}

func TestPatientCannotAccessAnalytics(t *testing.T) {
	testServer := newTestServer(t)
	serverURL := startTestHTTPServer(t, testServer.handler)

	patientClient := loginAs(t, serverURL, "paciente@mail.com", "secret123")
	requireStatusCode(t, patientClient.Get(t, "/api/v1/analytics"), http.StatusForbidden)
}
