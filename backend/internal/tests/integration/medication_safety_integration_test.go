package integration

import (
	"context"
	"net/http"
	"testing"
)

type safetyAlertBody struct {
	Kind              string            `json:"kind"`
	Severity          string            `json:"severity"`
	Blocking          bool              `json:"blocking"`
	MessageKey        string            `json:"messageKey"`
	Params            map[string]string `json:"params"`
	RelatedResourceID string            `json:"relatedResourceId"`
}

func TestMedicationSafetyBlocksAllergyUntilDoctorOverrides(t *testing.T) {
	testServer := newTestServer(t)
	serverURL := startTestHTTPServer(t, testServer.handler)

	receptionClient := loginAs(t, serverURL, "recepcao@hospital.com", "secret123")
	doctorClient := loginAs(t, serverURL, "medico@clinica.com", "secret123")
	nurseClient := loginAs(t, serverURL, "enfermeiro@hospital.com", "secret123")

	patientResponse := receptionClient.Post(t, "/api/v1/patients", map[string]string{
		"full_name":    "Marta Souza",
		"birth_date":   "1980-05-10",
		"document_id":  "529.982.247-25",
		"phone_number": "(11) 98123-4567",
	})
	requireStatusCode(t, patientResponse, http.StatusCreated)
	patientFhirID := requireJSONFieldValue(t, patientResponse, "fhir_resource_id")

	encounterResponse := doctorClient.Post(t, "/api/v1/patients/"+patientFhirID+"/encounters", map[string]string{
		"practitioner_id": "practitioner-1",
		"reason_code":     "J03",
		"reason_display":  "Amigdalite",
	})
	requireStatusCode(t, encounterResponse, http.StatusCreated)
	encounterFhirID := requireJSONFieldValue(t, encounterResponse, "fhir_id")

	allergyResponse := doctorClient.Post(t, "/api/v1/patients/"+patientFhirID+"/allergies", map[string]string{
		"allergen_code":    "91936005",
		"allergen_display": "Alergia a amoxicilina",
		"reaction":         "Urticária",
	})
	requireStatusCode(t, allergyResponse, http.StatusCreated)
	allergyFhirID := requireJSONFieldValue(t, allergyResponse, "fhir_id")

	checkResponse := doctorClient.Post(t, "/api/v1/patients/"+patientFhirID+"/medications/check", map[string]string{"medicationName": "Ampicilina 500mg"})
	requireStatusCode(t, checkResponse, http.StatusOK)
	var checkBody struct {
		Alerts []safetyAlertBody `json:"alerts"`
	}
	decodeJSONResponse(t, checkResponse, &checkBody)
	if len(checkBody.Alerts) != 1 {
		t.Fatalf("expected exactly 1 alert, got %d", len(checkBody.Alerts))
	}
	if !checkBody.Alerts[0].Blocking || checkBody.Alerts[0].Kind != "allergy" || checkBody.Alerts[0].Severity != "contraindicated" {
		t.Fatalf("unexpected alert: %+v", checkBody.Alerts[0])
	}
	if checkBody.Alerts[0].RelatedResourceID != allergyFhirID {
		t.Fatalf("expected alert to reference allergy %s, got %s", allergyFhirID, checkBody.Alerts[0].RelatedResourceID)
	}

	requireStatusCode(t, nurseClient.Post(t, "/api/v1/patients/"+patientFhirID+"/medications/check", map[string]string{"medicationName": "Ampicilina"}), http.StatusForbidden)

	prescriptionPath := "/api/v1/encounters/" + encounterFhirID + "/medications"
	prescription := map[string]string{
		"patient_fhir_id":      patientFhirID,
		"practitioner_fhir_id": "practitioner-1",
		"medication_name":      "Ampicilina 500mg",
		"dosage_instructions":  "1 capsula 6/6h",
	}

	blockedResponse := doctorClient.Post(t, prescriptionPath, prescription)
	requireStatusCode(t, blockedResponse, http.StatusConflict)
	var blockedBody struct {
		Error  string            `json:"error"`
		Code   string            `json:"code"`
		Alerts []safetyAlertBody `json:"alerts"`
	}
	decodeJSONResponse(t, blockedResponse, &blockedBody)
	if blockedBody.Code != "MEDICATION_ALLERGY_BLOCKED" || len(blockedBody.Alerts) != 1 {
		t.Fatalf("unexpected blocked body: %+v", blockedBody)
	}

	prescription["overrideJustification"] = "curta"
	requireStatusCode(t, doctorClient.Post(t, prescriptionPath, prescription), http.StatusConflict)

	prescription["overrideJustification"] = "Paciente tolerou a classe em tratamento anterior, beneficio supera o risco"
	requireStatusCode(t, nurseClient.Post(t, prescriptionPath, prescription), http.StatusForbidden)

	overrideResponse := doctorClient.Post(t, prescriptionPath, prescription)
	requireStatusCode(t, overrideResponse, http.StatusCreated)
	var overrideBody struct {
		FhirID   string            `json:"fhir_id"`
		Warnings []safetyAlertBody `json:"warnings"`
	}
	decodeJSONResponse(t, overrideResponse, &overrideBody)
	if overrideBody.FhirID == "" || len(overrideBody.Warnings) != 1 {
		t.Fatalf("unexpected override body: %+v", overrideBody)
	}

	var overrideCount int
	overrideQueryError := testServer.db.QueryRow(context.Background(),
		`SELECT count(*) FROM medication_safety_overrides WHERE patient_fhir_id = $1 AND medication_request_fhir_id = $2`,
		patientFhirID, overrideBody.FhirID).Scan(&overrideCount)
	if overrideQueryError != nil {
		t.Fatalf("failed to query overrides: %v", overrideQueryError)
	}
	if overrideCount != 1 {
		t.Fatalf("expected exactly 1 linked override, got %d", overrideCount)
	}

	var totalOverrideRows int
	if countError := testServer.db.QueryRow(context.Background(), `SELECT count(*) FROM medication_safety_overrides`).Scan(&totalOverrideRows); countError != nil {
		t.Fatalf("failed to count overrides: %v", countError)
	}
	if totalOverrideRows != 1 {
		t.Fatalf("expected rejected attempts to leave no override rows, got %d", totalOverrideRows)
	}

	var auditCount int
	auditQueryError := testServer.db.QueryRow(context.Background(),
		`SELECT count(*) FROM audit_logs WHERE action = 'medication_safety_override' AND resource_id = $1`,
		overrideBody.FhirID).Scan(&auditCount)
	if auditQueryError != nil {
		t.Fatalf("failed to query audit logs: %v", auditQueryError)
	}
	if auditCount != 1 {
		t.Fatalf("expected exactly 1 audit log for the override, got %d", auditCount)
	}
}
