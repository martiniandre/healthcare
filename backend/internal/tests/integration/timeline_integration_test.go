package integration

import (
	"context"
	"net/http"
	"testing"
	"time"
)

type timelineResponseBody struct {
	Items []struct {
		ID         string `json:"id"`
		Type       string `json:"type"`
		OccurredAt string `json:"occurredAt"`
	} `json:"items"`
	NextCursor *string `json:"nextCursor"`
	Partial    struct {
		FailedTypes []string `json:"failedTypes"`
	} `json:"partial"`
}

func seedTimelineResources(t *testing.T, testServer *testServer, patientFHIRID string) {
	t.Helper()
	ctx := context.Background()
	now := time.Now().UTC()
	patientReference := map[string]string{"reference": "Patient/" + patientFHIRID}

	seeds := []struct {
		resourceType string
		resource     map[string]interface{}
	}{
		{"Patient", map[string]interface{}{"id": patientFHIRID, "name": []map[string]interface{}{{"given": []string{"Joao"}, "family": "Souza"}}}},
		{"Encounter", map[string]interface{}{"id": "timeline-enc", "status": "finished", "subject": patientReference, "period": map[string]string{"start": now.Add(-72 * time.Hour).Format(time.RFC3339)}}},
		{"Observation", map[string]interface{}{"id": "timeline-obs", "status": "final", "subject": patientReference, "effectiveDateTime": now.Add(-24 * time.Hour).Format(time.RFC3339)}},
		{"Condition", map[string]interface{}{"id": "timeline-cond", "subject": patientReference, "recordedDate": now.Add(-48 * time.Hour).Format(time.RFC3339), "onsetDateTime": now.Add(-9000 * time.Hour).Format(time.RFC3339)}},
		{"DiagnosticReport", map[string]interface{}{"id": "timeline-old-report", "status": "final", "subject": patientReference, "issued": now.AddDate(-2, 0, 0).Format(time.RFC3339)}},
	}
	for _, seed := range seeds {
		if _, seedError := testServer.fhir.CreateResource(ctx, seed.resourceType, seed.resource); seedError != nil {
			t.Fatalf("failed to seed %s: %v", seed.resourceType, seedError)
		}
	}
}

func TestTimelineReturnsMergedEntriesNewestFirstWithinDefaultWindow(t *testing.T) {
	testServer := newTestServer(t)
	serverURL := startTestHTTPServer(t, testServer.handler)
	seedTimelineResources(t, testServer, "timeline-patient")

	doctorClient := loginAs(t, serverURL, "medico@clinica.com", "secret123")

	response := doctorClient.Get(t, "/api/v1/patients/timeline-patient/timeline")
	requireStatusCode(t, response, http.StatusOK)
	var body timelineResponseBody
	decodeJSONResponse(t, response, &body)

	expectedOrder := []string{"observation:timeline-obs", "condition:timeline-cond", "encounter:timeline-enc"}
	if len(body.Items) != len(expectedOrder) {
		t.Fatalf("expected %d items inside default window, got %d", len(expectedOrder), len(body.Items))
	}
	for itemIndex, expectedID := range expectedOrder {
		if body.Items[itemIndex].ID != expectedID {
			t.Fatalf("item %d: expected %s, got %s", itemIndex, expectedID, body.Items[itemIndex].ID)
		}
	}
	if len(body.Partial.FailedTypes) != 0 {
		t.Fatalf("expected no failed types, got %v", body.Partial.FailedTypes)
	}
	if body.NextCursor != nil {
		t.Fatal("expected no next cursor")
	}
}

func TestTimelinePaginatesWithCursorAndFiltersByType(t *testing.T) {
	testServer := newTestServer(t)
	serverURL := startTestHTTPServer(t, testServer.handler)
	seedTimelineResources(t, testServer, "timeline-patient")

	doctorClient := loginAs(t, serverURL, "medico@clinica.com", "secret123")

	firstResponse := doctorClient.Get(t, "/api/v1/patients/timeline-patient/timeline?limit=2")
	requireStatusCode(t, firstResponse, http.StatusOK)
	var firstPage timelineResponseBody
	decodeJSONResponse(t, firstResponse, &firstPage)
	if len(firstPage.Items) != 2 || firstPage.NextCursor == nil {
		t.Fatalf("expected 2 items and a cursor on first page, got %d items", len(firstPage.Items))
	}

	secondResponse := doctorClient.Get(t, "/api/v1/patients/timeline-patient/timeline?limit=2&cursor="+*firstPage.NextCursor)
	requireStatusCode(t, secondResponse, http.StatusOK)
	var secondPage timelineResponseBody
	decodeJSONResponse(t, secondResponse, &secondPage)
	if len(secondPage.Items) != 1 || secondPage.Items[0].ID != "encounter:timeline-enc" || secondPage.NextCursor != nil {
		t.Fatalf("unexpected second page: %+v", secondPage)
	}

	filteredResponse := doctorClient.Get(t, "/api/v1/patients/timeline-patient/timeline?types=condition")
	requireStatusCode(t, filteredResponse, http.StatusOK)
	var filteredPage timelineResponseBody
	decodeJSONResponse(t, filteredResponse, &filteredPage)
	if len(filteredPage.Items) != 1 || filteredPage.Items[0].Type != "condition" {
		t.Fatalf("expected only the condition entry, got %+v", filteredPage.Items)
	}
}

func TestTimelineRejectsInvalidParametersAndUnknownPatient(t *testing.T) {
	testServer := newTestServer(t)
	serverURL := startTestHTTPServer(t, testServer.handler)
	seedTimelineResources(t, testServer, "timeline-patient")

	doctorClient := loginAs(t, serverURL, "medico@clinica.com", "secret123")

	requireStatusCode(t, doctorClient.Get(t, "/api/v1/patients/timeline-patient/timeline?types=appointment"), http.StatusBadRequest)
	requireStatusCode(t, doctorClient.Get(t, "/api/v1/patients/timeline-patient/timeline?limit=1000"), http.StatusBadRequest)
	requireStatusCode(t, doctorClient.Get(t, "/api/v1/patients/timeline-patient/timeline?from=garbage"), http.StatusBadRequest)
	requireStatusCode(t, doctorClient.Get(t, "/api/v1/patients/does-not-exist/timeline"), http.StatusNotFound)
}

func TestTimelineIsForbiddenForPatientAndReceptionRoles(t *testing.T) {
	testServer := newTestServer(t)
	serverURL := startTestHTTPServer(t, testServer.handler)
	seedTimelineResources(t, testServer, "timeline-patient")

	patientClient := loginAs(t, serverURL, "paciente@mail.com", "secret123")
	requireStatusCode(t, patientClient.Get(t, "/api/v1/patients/timeline-patient/timeline"), http.StatusForbidden)

	receptionClient := loginAs(t, serverURL, "recepcao@hospital.com", "secret123")
	requireStatusCode(t, receptionClient.Get(t, "/api/v1/patients/timeline-patient/timeline"), http.StatusForbidden)
}
