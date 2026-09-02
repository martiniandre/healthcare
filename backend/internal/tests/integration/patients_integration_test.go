package integration

import (
	"net/http"
	"testing"
)

func TestPatientCRUDFlow(t *testing.T) {
	testServer := newTestServer(t)
	serverURL := startTestHTTPServer(t, testServer.handler)

	adminClient := loginAs(t, serverURL, "admin@hospital.com", "secret123")

	createResponse := adminClient.Post(t, "/api/v1/patients", map[string]string{
		"full_name":    "Maria Silva",
		"birth_date":   "1990-01-01",
		"document_id":  "529.982.247-25",
		"phone_number": "(11) 98765-4321",
	})
	requireStatusCode(t, createResponse, http.StatusCreated)

	fhirResourceID := readPatientResponseField(t, createResponse, "fhir_resource_id")
	if fhirResourceID == "" {
		t.Fatal("expected fhir_resource_id in create response")
	}

	getResponse := adminClient.Get(t, "/api/v1/patients/"+fhirResourceID)
	requireStatusCode(t, getResponse, http.StatusOK)

	documentID := readPatientResponseField(t, getResponse, "document_id")
	if documentID != "529.982.247-25" {
		t.Fatalf("expected document_id 529.982.247-25, got %q", documentID)
	}

	listResponse := adminClient.Get(t, "/api/v1/patients?page=1&limit=50")
	requireStatusCode(t, listResponse, http.StatusOK)

	var patientsEnvelope struct {
		Patients []map[string]interface{} `json:"patients"`
		Total    int                      `json:"total"`
	}
	decodeJSONResponse(t, listResponse, &patientsEnvelope)
	if len(patientsEnvelope.Patients) < 1 {
		t.Fatal("expected at least one patient in the list")
	}
	if patientsEnvelope.Total < 1 {
		t.Fatalf("expected total >= 1, got %d", patientsEnvelope.Total)
	}
}

func TestCreatePatientRejectsInvalidCPF(t *testing.T) {
	testServer := newTestServer(t)
	serverURL := startTestHTTPServer(t, testServer.handler)

	adminClient := loginAs(t, serverURL, "admin@hospital.com", "secret123")

	response := adminClient.Post(t, "/api/v1/patients", map[string]string{
		"full_name":    "João Souza",
		"birth_date":   "1985-05-15",
		"document_id":  "123.456.789-00",
		"phone_number": "(11) 91234-5678",
	})
	requireStatusCode(t, response, http.StatusBadRequest)
}

func TestCreatePatientRejectsDuplicateDocument(t *testing.T) {
	testServer := newTestServer(t)
	serverURL := startTestHTTPServer(t, testServer.handler)

	adminClient := loginAs(t, serverURL, "admin@hospital.com", "secret123")

	firstResponse := adminClient.Post(t, "/api/v1/patients", map[string]string{
		"full_name":    "Maria Silva",
		"birth_date":   "1990-01-01",
		"document_id":  "529.982.247-25",
		"phone_number": "(11) 98765-4321",
	})
	requireStatusCode(t, firstResponse, http.StatusCreated)

	duplicateResponse := adminClient.Post(t, "/api/v1/patients", map[string]string{
		"full_name":    "Maria Duplicada",
		"birth_date":   "1991-02-02",
		"document_id":  "529.982.247-25",
		"phone_number": "(11) 99876-5432",
	})
	requireStatusCode(t, duplicateResponse, http.StatusConflict)
}

func TestGetNonExistentPatientReturnsNotFound(t *testing.T) {
	testServer := newTestServer(t)
	serverURL := startTestHTTPServer(t, testServer.handler)

	adminClient := loginAs(t, serverURL, "admin@hospital.com", "secret123")

	response := adminClient.Get(t, "/api/v1/patients/nao-existe")
	requireStatusCode(t, response, http.StatusNotFound)
}

func TestPatientRoleCannotCreatePatient(t *testing.T) {
	testServer := newTestServer(t)
	serverURL := startTestHTTPServer(t, testServer.handler)

	patientClient := loginAs(t, serverURL, "paciente@mail.com", "secret123")

	response := patientClient.Post(t, "/api/v1/patients", map[string]string{
		"full_name":    "Paciente Role",
		"birth_date":   "1990-01-01",
		"document_id":  "529.982.247-25",
		"phone_number": "(11) 98765-4321",
	})
	requireStatusCode(t, response, http.StatusForbidden)
}

func readPatientResponseField(t *testing.T, response *http.Response, fieldName string) string {
	t.Helper()
	return requireJSONFieldValue(t, response, fieldName)
}
