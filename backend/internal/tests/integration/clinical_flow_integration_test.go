package integration

import (
	"context"
	"net/http"
	"testing"
)

func TestClinicalFlowFromEncounterToMedication(t *testing.T) {
	testServer := newTestServer(t)
	serverURL := startTestHTTPServer(t, testServer.handler)

	receptionClient := loginAs(t, serverURL, "recepcao@hospital.com", "secret123")
	doctorClient := loginAs(t, serverURL, "medico@clinica.com", "secret123")

	patientResponse := receptionClient.Post(t, "/api/v1/patients", map[string]string{
		"full_name":    "Carlos Pereira",
		"birth_date":   "1975-03-20",
		"document_id":  "111.444.777-35",
		"phone_number": "(11) 98123-4567",
	})
	requireStatusCode(t, patientResponse, http.StatusCreated)
	patientFhirID := requireJSONFieldValue(t, patientResponse, "fhir_resource_id")

	encounterResponse := doctorClient.Post(t, "/api/v1/patients/"+patientFhirID+"/encounters", map[string]string{
		"practitioner_id": "practitioner-1",
		"reason_code":     "R51",
		"reason_display":  "Dor de cabeça",
	})
	requireStatusCode(t, encounterResponse, http.StatusCreated)
	encounterFhirID := requireJSONFieldValue(t, encounterResponse, "fhir_id")
	if encounterFhirID == "" {
		t.Fatal("expected fhir_id in encounter response")
	}

	observationResponse := doctorClient.Post(t, "/api/v1/encounters/"+encounterFhirID+"/observations", map[string]interface{}{
		"patient_fhir_id": patientFhirID,
		"loinc_code":      "8867-4",
		"code_display":    "Frequência cardíaca",
		"value_quantity":  72,
		"value_unit":      "bpm",
	})
	requireStatusCode(t, observationResponse, http.StatusCreated)
	observationFhirID := requireJSONFieldValue(t, observationResponse, "fhir_id")

	conditionResponse := doctorClient.Post(t, "/api/v1/patients/"+patientFhirID+"/conditions", map[string]string{
		"icd10_code":   "I10",
		"code_display": "Hipertensão essencial",
		"encounter_id": encounterFhirID,
	})
	requireStatusCode(t, conditionResponse, http.StatusCreated)
	conditionFhirID := requireJSONFieldValue(t, conditionResponse, "fhir_id")

	allergyResponse := doctorClient.Post(t, "/api/v1/patients/"+patientFhirID+"/allergies", map[string]string{
		"allergen_code":    "91934008",
		"allergen_display": "Alergia à penicilina",
		"reaction":         "Urticária",
	})
	requireStatusCode(t, allergyResponse, http.StatusCreated)
	allergyFhirID := requireJSONFieldValue(t, allergyResponse, "fhir_id")

	medicationResponse := doctorClient.Post(t, "/api/v1/encounters/"+encounterFhirID+"/medications", map[string]string{
		"patient_fhir_id":      patientFhirID,
		"practitioner_fhir_id": "practitioner-1",
		"medication_name":      "Amoxicilina 500mg",
		"dosage_instructions":  "1 comprimido 8/8h por 7 dias",
	})
	requireStatusCode(t, medicationResponse, http.StatusCreated)

	requireStatusCode(t, doctorClient.Get(t, "/api/v1/patients/"+patientFhirID+"/observations"), http.StatusOK)
	requireStatusCode(t, doctorClient.Get(t, "/api/v1/encounters/"+encounterFhirID+"/observations"), http.StatusOK)
	requireStatusCode(t, doctorClient.Get(t, "/api/v1/patients/"+patientFhirID+"/conditions"), http.StatusOK)
	requireStatusCode(t, doctorClient.Get(t, "/api/v1/patients/"+patientFhirID+"/allergies"), http.StatusOK)
	requireStatusCode(t, doctorClient.Get(t, "/api/v1/encounters/"+encounterFhirID+"/medications"), http.StatusOK)
	requireStatusCode(t, doctorClient.Get(t, "/api/v1/patients/"+patientFhirID+"/encounters"), http.StatusOK)
	requireStatusCode(t, doctorClient.Get(t, "/api/v1/encounters/"+encounterFhirID), http.StatusOK)

	requireStatusCode(t, doctorClient.Delete(t, "/api/v1/observations/"+observationFhirID), http.StatusNoContent)
	requireStatusCode(t, doctorClient.Delete(t, "/api/v1/patients/"+patientFhirID+"/conditions/"+conditionFhirID), http.StatusNoContent)
	requireStatusCode(t, doctorClient.Delete(t, "/api/v1/patients/"+patientFhirID+"/allergies/"+allergyFhirID), http.StatusNoContent)
}

func TestCreateEncounterRequiresPractitionerAndReason(t *testing.T) {
	testServer := newTestServer(t)
	serverURL := startTestHTTPServer(t, testServer.handler)

	doctorClient := loginAs(t, serverURL, "medico@clinica.com", "secret123")

	response := doctorClient.Post(t, "/api/v1/patients/patient-1/encounters", map[string]string{
		"reason_display": "Consulta de rotina",
	})
	requireStatusCode(t, response, http.StatusBadRequest)
}

func TestGetNonExistentEncounterReturnsNotFound(t *testing.T) {
	testServer := newTestServer(t)
	serverURL := startTestHTTPServer(t, testServer.handler)

	doctorClient := loginAs(t, serverURL, "medico@clinica.com", "secret123")

	response := doctorClient.Get(t, "/api/v1/encounters/nao-existe")
	requireStatusCode(t, response, http.StatusNotFound)
}

func TestNotificationsAreCreatedOnLogin(t *testing.T) {
	testServer := newTestServer(t)
	serverURL := startTestHTTPServer(t, testServer.handler)

	adminClient := loginAs(t, serverURL, "admin@hospital.com", "secret123")

	notificationsResponse := adminClient.Get(t, "/api/v1/notifications")
	requireStatusCode(t, notificationsResponse, http.StatusOK)

	var notificationsBody struct {
		Notifications []map[string]interface{} `json:"notifications"`
		Total         int                      `json:"total"`
	}
	decodeJSONResponse(t, notificationsResponse, &notificationsBody)

	var notificationCount int
	countError := testServer.db.QueryRow(context.Background(),
		"SELECT count(*) FROM notifications").Scan(&notificationCount)
	if countError != nil {
		t.Fatalf("failed to query notifications: %v", countError)
	}
	if notificationCount == 0 {
		t.Fatal("expected notifications to be created after login event")
	}
	if len(notificationsBody.Notifications) == 0 {
		t.Fatal("expected notifications list to be non-empty")
	}
	if notificationsBody.Total != notificationCount {
		t.Fatalf("expected total %d, got %d", notificationCount, notificationsBody.Total)
	}
}
