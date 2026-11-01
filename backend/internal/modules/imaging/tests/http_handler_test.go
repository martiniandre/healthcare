package tests

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/healthcare/backend/internal/modules/imaging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func buildMultipartUploadRequest(testingInstance *testing.T, fileBytes []byte, formFields map[string]string) *http.Request {
	testingInstance.Helper()

	var bodyBuffer bytes.Buffer
	multipartWriter := multipart.NewWriter(&bodyBuffer)
	for fieldName, fieldValue := range formFields {
		require.NoError(testingInstance, multipartWriter.WriteField(fieldName, fieldValue))
	}
	filePart, partError := multipartWriter.CreateFormFile("file", "scan.dcm")
	require.NoError(testingInstance, partError)
	_, writeError := filePart.Write(fileBytes)
	require.NoError(testingInstance, writeError)
	require.NoError(testingInstance, multipartWriter.Close())

	request := httptest.NewRequest(http.MethodPost, "/api/v1/patients/"+fixturePatientFhirID+"/studies", &bodyBuffer)
	request.Header.Set("Content-Type", multipartWriter.FormDataContentType())
	request.SetPathValue("patientFhirId", fixturePatientFhirID)
	return request
}

func TestHTTPHandler_UploadPatientStudy_ReturnsHeaderDerivedFields(testingInstance *testing.T) {
	harness := newUploadHarness()
	handler := imaging.NewHTTPHandler(harness.service)
	recorder := httptest.NewRecorder()

	handler.UploadPatientStudy(recorder, buildMultipartUploadRequest(testingInstance, newDICOMFixture().bytes(), map[string]string{"modality": "CT"}))

	require.Equal(testingInstance, http.StatusCreated, recorder.Code)
	var responseBody map[string]interface{}
	require.NoError(testingInstance, json.Unmarshal(recorder.Body.Bytes(), &responseBody))
	assert.Equal(testingInstance, fixtureStudyInstanceUID, responseBody["studyInstanceUid"])
	assert.Equal(testingInstance, fixtureStudyInstanceUID, responseBody["study_instance_uid"])
	assert.Equal(testingInstance, "MR", responseBody["modality"])
	assert.Equal(testingInstance, "MR - BRAIN WITHOUT CONTRAST", responseBody["title"])
	assert.Equal(testingInstance, "2026-01-15T10:15:00Z", responseBody["studyDate"])
	assert.Equal(testingInstance, "ACC-0042", responseBody["accessionNumber"])
	assert.Equal(testingInstance, fixtureSeriesUID, responseBody["seriesInstanceUid"])
	assert.Equal(testingInstance, fixtureSOPInstanceUID, responseBody["sopInstanceUid"])
	assert.Equal(testingInstance, imaging.ImagingStudyStatusReceiving, responseBody["status"])
	assert.NotEmpty(testingInstance, responseBody["id"])
	assert.Equal(testingInstance, fixturePatientFhirID, responseBody["patient_fhir_id"])

	warnings, isList := responseBody["warnings"].([]interface{})
	require.True(testingInstance, isList)
	require.Len(testingInstance, warnings, 1)
	firstWarning := warnings[0].(map[string]interface{})
	assert.Equal(testingInstance, imaging.WarningCodeModalityOverridden, firstWarning["code"])
	assert.Equal(testingInstance, "imaging.upload.warning.modality_overridden", firstWarning["message"])
}

func TestHTTPHandler_UploadPatientStudy_WarningsIsAnEmptyListWhenClean(testingInstance *testing.T) {
	harness := newUploadHarness()
	handler := imaging.NewHTTPHandler(harness.service)
	recorder := httptest.NewRecorder()

	handler.UploadPatientStudy(recorder, buildMultipartUploadRequest(testingInstance, newDICOMFixture().bytes(), nil))

	require.Equal(testingInstance, http.StatusCreated, recorder.Code)
	assert.Contains(testingInstance, recorder.Body.String(), `"warnings":[]`)
}

func TestHTTPHandler_UploadPatientStudy_MismatchReturns422WithStableCode(testingInstance *testing.T) {
	harness := newUploadHarness()
	harness.identity.Identity = imaging.PatientIdentity{FullNames: []string{"Pedro Santos"}, DocumentIDs: []string{"11144477735"}}
	handler := imaging.NewHTTPHandler(harness.service)
	recorder := httptest.NewRecorder()

	handler.UploadPatientStudy(recorder, buildMultipartUploadRequest(testingInstance, newDICOMFixture().bytes(), nil))

	require.Equal(testingInstance, http.StatusUnprocessableEntity, recorder.Code)
	var responseBody map[string]string
	require.NoError(testingInstance, json.Unmarshal(recorder.Body.Bytes(), &responseBody))
	assert.Equal(testingInstance, "DICOM_PATIENT_MISMATCH", responseBody["code"])
	assert.NotEmpty(testingInstance, responseBody["error"])
	assert.NotContains(testingInstance, recorder.Body.String(), "SILVA")
	assert.NotContains(testingInstance, recorder.Body.String(), "529.982")
}

func TestHTTPHandler_UploadPatientStudy_DuplicateReturns409WithStableCode(testingInstance *testing.T) {
	harness := newUploadHarness()
	handler := imaging.NewHTTPHandler(harness.service)
	fileBytes := newDICOMFixture().bytes()

	handler.UploadPatientStudy(httptest.NewRecorder(), buildMultipartUploadRequest(testingInstance, fileBytes, nil))
	recorder := httptest.NewRecorder()
	handler.UploadPatientStudy(recorder, buildMultipartUploadRequest(testingInstance, fileBytes, nil))

	require.Equal(testingInstance, http.StatusConflict, recorder.Code)
	var responseBody map[string]string
	require.NoError(testingInstance, json.Unmarshal(recorder.Body.Bytes(), &responseBody))
	assert.Equal(testingInstance, "DICOM_DUPLICATE_INSTANCE", responseBody["code"])
}

func TestHTTPHandler_UploadPatientStudy_InvalidDICOMReturns400WithReason(testingInstance *testing.T) {
	harness := newUploadHarness()
	handler := imaging.NewHTTPHandler(harness.service)
	fixture := newDICOMFixture()
	fixture.studyInstanceUID = ""
	recorder := httptest.NewRecorder()

	handler.UploadPatientStudy(recorder, buildMultipartUploadRequest(testingInstance, fixture.bytes(), nil))

	require.Equal(testingInstance, http.StatusBadRequest, recorder.Code)
	var responseBody map[string]string
	require.NoError(testingInstance, json.Unmarshal(recorder.Body.Bytes(), &responseBody))
	assert.Equal(testingInstance, "INVALID_DICOM", responseBody["code"])
	assert.Contains(testingInstance, responseBody["error"], "StudyInstanceUID is missing")
}

func TestHTTPHandler_UploadPatientStudy_MissingFileReturns400(testingInstance *testing.T) {
	harness := newUploadHarness()
	handler := imaging.NewHTTPHandler(harness.service)

	var bodyBuffer bytes.Buffer
	multipartWriter := multipart.NewWriter(&bodyBuffer)
	require.NoError(testingInstance, multipartWriter.WriteField("title", "Scan"))
	require.NoError(testingInstance, multipartWriter.Close())
	request := httptest.NewRequest(http.MethodPost, "/api/v1/patients/"+fixturePatientFhirID+"/studies", &bodyBuffer)
	request.Header.Set("Content-Type", multipartWriter.FormDataContentType())
	request.SetPathValue("patientFhirId", fixturePatientFhirID)
	recorder := httptest.NewRecorder()

	handler.UploadPatientStudy(recorder, request)

	assert.Equal(testingInstance, http.StatusBadRequest, recorder.Code)
}
