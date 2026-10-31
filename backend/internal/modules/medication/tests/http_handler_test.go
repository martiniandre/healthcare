package tests

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/healthcare/backend/internal/modules/medication"
	"github.com/healthcare/backend/internal/modules/medication/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type decodedAlert struct {
	Kind              string            `json:"kind"`
	Severity          string            `json:"severity"`
	Blocking          bool              `json:"blocking"`
	MessageKey        string            `json:"messageKey"`
	Params            map[string]string `json:"params"`
	RelatedResourceID string            `json:"relatedResourceId"`
}

func newHandlerWithAllergy(createdCounter *int) *medication.HTTPHandler {
	fixture := newSafetyDependencies([]medication.PatientAllergy{penicillinAllergy()}, nil)
	medicationService := medication.NewService(countingRepository(createdCounter), fixture.dependencies)
	return medication.NewHTTPHandler(medicationService)
}

func performRequest(handlerFunc http.HandlerFunc, pathKey string, pathValue string, requestBody string, requestContext context.Context) *httptest.ResponseRecorder {
	httpRequest := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(requestBody)).WithContext(requestContext)
	httpRequest.SetPathValue(pathKey, pathValue)
	responseRecorder := httptest.NewRecorder()
	handlerFunc(responseRecorder, httpRequest)
	return responseRecorder
}

func TestCheckMedicationSafetyHandler_ReturnsAlerts(testingInstance *testing.T) {
	createdCounter := 0
	handler := newHandlerWithAllergy(&createdCounter)

	response := performRequest(handler.CheckMedicationSafety, "patientFhirId", "patient-1", `{"medicationCode":"","medicationName":"Ampicilina"}`, contextWithCaller("DOCTOR", testDoctorUserID))

	require.Equal(testingInstance, http.StatusOK, response.Code)
	var body struct {
		Alerts []decodedAlert `json:"alerts"`
	}
	require.NoError(testingInstance, json.Unmarshal(response.Body.Bytes(), &body))
	require.Len(testingInstance, body.Alerts, 1)
	assert.Equal(testingInstance, "allergy", body.Alerts[0].Kind)
	assert.Equal(testingInstance, "contraindicated", body.Alerts[0].Severity)
	assert.True(testingInstance, body.Alerts[0].Blocking)
	assert.Equal(testingInstance, "medicationSafety.allergy.class", body.Alerts[0].MessageKey)
	assert.Equal(testingInstance, "allergy-1", body.Alerts[0].RelatedResourceID)
	assert.Equal(testingInstance, "Penicillins", body.Alerts[0].Params["drugClass"])
	assert.Equal(testingInstance, 0, createdCounter)
}

func TestCheckMedicationSafetyHandler_NoAlertsReturnsEmptyArray(testingInstance *testing.T) {
	handler := medication.NewHTTPHandler(medication.NewService(&mocks.MockMedicationRepository{}, newSafetyDependencies(nil, nil).dependencies))

	response := performRequest(handler.CheckMedicationSafety, "patientFhirId", "patient-1", `{"medicationName":"Dipirona"}`, contextWithCaller("DOCTOR", testDoctorUserID))

	require.Equal(testingInstance, http.StatusOK, response.Code)
	assert.JSONEq(testingInstance, `{"alerts":[]}`, response.Body.String())
}

func TestCheckMedicationSafetyHandler_InvalidInput(testingInstance *testing.T) {
	handler := medication.NewHTTPHandler(medication.NewService(&mocks.MockMedicationRepository{}, newSafetyDependencies(nil, nil).dependencies))

	malformed := performRequest(handler.CheckMedicationSafety, "patientFhirId", "patient-1", `{`, contextWithCaller("DOCTOR", testDoctorUserID))
	missingMedication := performRequest(handler.CheckMedicationSafety, "patientFhirId", "patient-1", `{}`, contextWithCaller("DOCTOR", testDoctorUserID))

	assert.Equal(testingInstance, http.StatusBadRequest, malformed.Code)
	assert.Equal(testingInstance, http.StatusBadRequest, missingMedication.Code)
}

func TestCreateMedicationHandler_AllergyBlocked_Returns409WithAlerts(testingInstance *testing.T) {
	createdCounter := 0
	handler := newHandlerWithAllergy(&createdCounter)
	requestBody := `{"patient_fhir_id":"patient-1","practitioner_fhir_id":"practitioner-1","medication_name":"Amoxicilina 500mg"}`

	response := performRequest(handler.CreateMedication, "encounterFhirId", "encounter-1", requestBody, contextWithCaller("DOCTOR", testDoctorUserID))

	require.Equal(testingInstance, http.StatusConflict, response.Code)
	var body struct {
		Error  string         `json:"error"`
		Code   string         `json:"code"`
		Alerts []decodedAlert `json:"alerts"`
	}
	require.NoError(testingInstance, json.Unmarshal(response.Body.Bytes(), &body))
	assert.Equal(testingInstance, "MEDICATION_ALLERGY_BLOCKED", body.Code)
	assert.NotEmpty(testingInstance, body.Error)
	require.Len(testingInstance, body.Alerts, 1)
	assert.True(testingInstance, body.Alerts[0].Blocking)
	assert.Equal(testingInstance, 0, createdCounter)
}

func TestCreateMedicationHandler_OverrideCreatesAndReturnsWarnings(testingInstance *testing.T) {
	createdCounter := 0
	handler := newHandlerWithAllergy(&createdCounter)
	requestBody := `{"patient_fhir_id":"patient-1","practitioner_fhir_id":"practitioner-1","medication_name":"Amoxicilina 500mg","overrideJustification":"` + validJustification + `"}`

	response := performRequest(handler.CreateMedication, "encounterFhirId", "encounter-1", requestBody, contextWithCaller("DOCTOR", testDoctorUserID))

	require.Equal(testingInstance, http.StatusCreated, response.Code)
	var body struct {
		FhirID   string         `json:"fhir_id"`
		Warnings []decodedAlert `json:"warnings"`
	}
	require.NoError(testingInstance, json.Unmarshal(response.Body.Bytes(), &body))
	assert.Equal(testingInstance, "medication-request-1", body.FhirID)
	require.Len(testingInstance, body.Warnings, 1)
	assert.Equal(testingInstance, "allergy", body.Warnings[0].Kind)
}

func TestCreateMedicationHandler_NoAlertsReturnsEmptyWarnings(testingInstance *testing.T) {
	createdCounter := 0
	handler := medication.NewHTTPHandler(medication.NewService(countingRepository(&createdCounter), newSafetyDependencies(nil, nil).dependencies))

	response := performRequest(handler.CreateMedication, "encounterFhirId", "encounter-1", `{"patient_fhir_id":"patient-1","medication_name":"Dipirona"}`, contextWithCaller("DOCTOR", testDoctorUserID))

	require.Equal(testingInstance, http.StatusCreated, response.Code)
	assert.Contains(testingInstance, response.Body.String(), `"warnings":[]`)
}

func TestCreateMedicationHandler_NurseOverrideIsForbidden(testingInstance *testing.T) {
	createdCounter := 0
	handler := newHandlerWithAllergy(&createdCounter)
	requestBody := `{"patient_fhir_id":"patient-1","medication_name":"Amoxicilina","overrideJustification":"` + validJustification + `"}`

	response := performRequest(handler.CreateMedication, "encounterFhirId", "encounter-1", requestBody, contextWithCaller("NURSE", testNurseUserID))

	assert.Equal(testingInstance, http.StatusForbidden, response.Code)
	assert.Equal(testingInstance, 0, createdCounter)
}
