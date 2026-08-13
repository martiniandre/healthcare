package integration

import (
	"net/http"
	"testing"
)

func TestDiagnosticReportLifecycleWithVersions(t *testing.T) {
	testServer := newTestServer(t)
	serverURL := startTestHTTPServer(t, testServer.handler)

	receptionClient := loginAs(t, serverURL, "recepcao@hospital.com", "secret123")
	doctorClient := loginAs(t, serverURL, "medico@clinica.com", "secret123")

	patientResponse := receptionClient.Post(t, "/api/v1/patients", map[string]string{
		"full_name":    "Fernanda Lima",
		"birth_date":   "1988-07-11",
		"document_id":  "111.444.777-35",
		"phone_number": "(21) 99876-5432",
	})
	requireStatusCode(t, patientResponse, http.StatusCreated)
	patientFhirID := requireJSONFieldValue(t, patientResponse, "fhir_resource_id")

	encounterResponse := doctorClient.Post(t, "/api/v1/patients/"+patientFhirID+"/encounters", map[string]string{
		"practitioner_id": "practitioner-1",
		"reason_code":     "R10",
		"reason_display":  "Dor abdominal",
	})
	requireStatusCode(t, encounterResponse, http.StatusCreated)
	encounterFhirID := requireJSONFieldValue(t, encounterResponse, "fhir_id")

	createResponse := doctorClient.Post(t, "/api/v1/encounters/"+encounterFhirID+"/reports", map[string]string{
		"patient_fhir_id": patientFhirID,
		"report_code":     "8867-4",
		"report_display":  "Frequência cardíaca",
		"conclusion":      "Dentro da normalidade",
	})
	requireStatusCode(t, createResponse, http.StatusCreated)
	reportFhirID := requireJSONFieldValue(t, createResponse, "fhir_id")

	listResponse := doctorClient.Get(t, "/api/v1/encounters/"+encounterFhirID+"/reports")
	requireStatusCode(t, listResponse, http.StatusOK)
	var reportsList []map[string]interface{}
	decodeJSONResponse(t, listResponse, &reportsList)
	if len(reportsList) == 0 {
		t.Fatal("expected at least one report for the encounter")
	}

	updateResponse := doctorClient.Put(t, "/api/v1/reports/"+reportFhirID, map[string]string{
		"conclusion": "Reavaliar em 30 dias",
		"status":     "final",
	})
	requireStatusCode(t, updateResponse, http.StatusOK)

	versionsResponse := doctorClient.Get(t, "/api/v1/reports/"+reportFhirID+"/versions")
	requireStatusCode(t, versionsResponse, http.StatusOK)
	var versionsList []map[string]interface{}
	decodeJSONResponse(t, versionsResponse, &versionsList)
	if len(versionsList) < 2 {
		t.Fatalf("expected at least 2 report versions, got %d", len(versionsList))
	}

	requireStatusCode(t, doctorClient.Get(t, "/api/v1/reports/"+reportFhirID+"/versions/1"), http.StatusOK)
}

func TestCreateReportRejectsInvalidLOINC(t *testing.T) {
	testServer := newTestServer(t)
	serverURL := startTestHTTPServer(t, testServer.handler)

	doctorClient := loginAs(t, serverURL, "medico@clinica.com", "secret123")

	response := doctorClient.Post(t, "/api/v1/encounters/encounter-1/reports", map[string]string{
		"patient_fhir_id": "patient-1",
		"report_code":     "NOT-A-LOINC",
		"report_display":  "Exame inválido",
	})
	requireStatusCode(t, response, http.StatusBadRequest)
}
