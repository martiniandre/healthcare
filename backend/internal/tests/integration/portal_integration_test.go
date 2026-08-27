package integration

import (
	"context"
	"fmt"
	"net/http"
	"testing"
)

func TestPortalDashboardReturnsPatientClinicalSummary(t *testing.T) {
	testServer := newTestServer(t)
	serverURL := startTestHTTPServer(t, testServer.handler)

	patientClient := loginAs(t, serverURL, "paciente@mail.com", "secret123")

	seedPortalEncounterForPatient(t, testServer)

	response := patientClient.Get(t, "/api/v1/portal/dashboard")
	requireStatusCode(t, response, http.StatusOK)
	var dashboardBody struct {
		PatientInfo struct {
			FullName string `json:"full_name"`
		} `json:"patient_info"`
		UpcomingEncounters []map[string]interface{} `json:"upcoming_encounters"`
	}
	decodeJSONResponse(t, response, &dashboardBody)
	if dashboardBody.PatientInfo.FullName != "Maria Silva" {
		t.Fatalf("expected patient Maria Silva, got %q", dashboardBody.PatientInfo.FullName)
	}
	if len(dashboardBody.UpcomingEncounters) == 0 {
		t.Fatal("expected upcoming encounter in portal dashboard")
	}
}

func TestPortalPatientEndpointsReturnCollections(t *testing.T) {
	testServer := newTestServer(t)
	serverURL := startTestHTTPServer(t, testServer.handler)

	patientClient := loginAs(t, serverURL, "paciente@mail.com", "secret123")

	for _, endpoint := range []string{
		"/api/v1/portal/encounters",
		"/api/v1/portal/observations",
		"/api/v1/portal/conditions",
		"/api/v1/portal/medications",
		"/api/v1/portal/reports",
		"/api/v1/portal/imaging",
	} {
		requireStatusCode(t, patientClient.Get(t, endpoint), http.StatusOK)
	}
}

func TestDoctorCannotAccessPatientPortal(t *testing.T) {
	testServer := newTestServer(t)
	serverURL := startTestHTTPServer(t, testServer.handler)

	doctorClient := loginAs(t, serverURL, "medico@clinica.com", "secret123")
	requireStatusCode(t, doctorClient.Get(t, "/api/v1/portal/dashboard"), http.StatusForbidden)
}

func seedPortalEncounterForPatient(t *testing.T, testServer *testServer) {
	t.Helper()
	ctx := context.Background()

	var patientFHIRID string
	if queryError := testServer.db.QueryRow(ctx, `
		SELECT l.patient_fhir_id
		FROM patient_user_links l
		JOIN users u ON u.id = l.user_id
		WHERE u.email = 'paciente@mail.com'`).Scan(&patientFHIRID); queryError != nil {
		t.Fatalf("failed to resolve linked patient fhir id: %v", queryError)
	}

	_, err := testServer.fhir.CreateResource(ctx, "Encounter", map[string]interface{}{
		"subject": map[string]string{"reference": fmt.Sprintf("Patient/%s", patientFHIRID)},
		"status":  "planned",
		"period": map[string]string{
			"start": "2026-09-01T09:00:00Z",
		},
		"reasonCode": []map[string]interface{}{
			{"coding": []map[string]interface{}{{"system": "http://snomed.info/sct", "code": "R10", "display": "Dor abdominal"}}},
		},
	})
	if err != nil {
		t.Fatalf("failed to seed portal encounter: %v", err)
	}
}
