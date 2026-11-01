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

func buildUploadRequestWithFile(testingInstance *testing.T, fileName string, fileBytes []byte, formFields map[string]string) *http.Request {
	testingInstance.Helper()

	var bodyBuffer bytes.Buffer
	multipartWriter := multipart.NewWriter(&bodyBuffer)
	for fieldName, fieldValue := range formFields {
		require.NoError(testingInstance, multipartWriter.WriteField(fieldName, fieldValue))
	}
	filePart, partError := multipartWriter.CreateFormFile("file", fileName)
	require.NoError(testingInstance, partError)
	_, writeError := filePart.Write(fileBytes)
	require.NoError(testingInstance, writeError)
	require.NoError(testingInstance, multipartWriter.Close())

	request := httptest.NewRequest(http.MethodPost, "/api/v1/patients/"+fixturePatientFhirID+"/studies", &bodyBuffer)
	request.Header.Set("Content-Type", multipartWriter.FormDataContentType())
	request.SetPathValue("patientFhirId", fixturePatientFhirID)
	return request
}

func decodeErrorBody(testingInstance *testing.T, recorder *httptest.ResponseRecorder) map[string]string {
	testingInstance.Helper()
	var responseBody map[string]string
	require.NoError(testingInstance, json.Unmarshal(recorder.Body.Bytes(), &responseBody))
	return responseBody
}

func TestHTTPHandler_UploadPatientStudy_SingleFileReportsStudyFields(testingInstance *testing.T) {
	harness := newUploadHarness()
	handler := imaging.NewHTTPHandler(harness.service)
	recorder := httptest.NewRecorder()

	handler.UploadPatientStudy(recorder, buildUploadRequestWithFile(testingInstance, "scan.dcm", newDICOMFixture().bytes(), nil))

	require.Equal(testingInstance, http.StatusCreated, recorder.Code)
	var responseBody map[string]interface{}
	require.NoError(testingInstance, json.Unmarshal(recorder.Body.Bytes(), &responseBody))
	assert.NotEmpty(testingInstance, responseBody["studyId"])
	assert.Equal(testingInstance, responseBody["id"], responseBody["studyId"])
	assert.Equal(testingInstance, true, responseBody["studyCreated"])
	assert.Equal(testingInstance, float64(1), responseBody["instancesAccepted"])
	assert.Equal(testingInstance, []interface{}{}, responseBody["instancesRejected"])
}

func TestHTTPHandler_UploadPatientStudy_AppendReportsStudyNotCreated(testingInstance *testing.T) {
	harness := newUploadHarness()
	handler := imaging.NewHTTPHandler(harness.service)
	handler.UploadPatientStudy(httptest.NewRecorder(), buildUploadRequestWithFile(testingInstance, "scan.dcm", instanceFixtureWithSOP("1").bytes(), nil))
	recorder := httptest.NewRecorder()

	handler.UploadPatientStudy(recorder, buildUploadRequestWithFile(testingInstance, "second.dcm", instanceFixtureWithSOP("2").bytes(), nil))

	require.Equal(testingInstance, http.StatusCreated, recorder.Code)
	var responseBody map[string]interface{}
	require.NoError(testingInstance, json.Unmarshal(recorder.Body.Bytes(), &responseBody))
	assert.Equal(testingInstance, false, responseBody["studyCreated"])
	assert.Equal(testingInstance, float64(1), responseBody["instancesAccepted"])
}

func TestHTTPHandler_UploadPatientStudy_ZIPReportsAcceptedAndRejectedInstances(testingInstance *testing.T) {
	harness := newUploadHarness()
	handler := imaging.NewHTTPHandler(harness.service)
	archiveBytes := buildZIPArchive(testingInstance, []zipEntry{
		{name: "first.dcm", content: instanceFixtureWithSOP("1").bytes()},
		{name: "second.dcm", content: instanceFixtureWithSOP("2").bytes()},
		{name: "notes.txt", content: []byte("not dicom")},
	})
	recorder := httptest.NewRecorder()

	handler.UploadPatientStudy(recorder, buildUploadRequestWithFile(testingInstance, "study.zip", archiveBytes, map[string]string{"title": "Knee series"}))

	require.Equal(testingInstance, http.StatusCreated, recorder.Code)
	var responseBody struct {
		ID                string `json:"id"`
		StudyID           string `json:"studyId"`
		Title             string `json:"title"`
		StudyCreated      bool   `json:"studyCreated"`
		InstancesAccepted int    `json:"instancesAccepted"`
		InstancesRejected []struct {
			FileName string `json:"fileName"`
			Code     string `json:"code"`
			Reason   string `json:"reason"`
		} `json:"instancesRejected"`
	}
	require.NoError(testingInstance, json.Unmarshal(recorder.Body.Bytes(), &responseBody))
	assert.Equal(testingInstance, responseBody.ID, responseBody.StudyID)
	assert.Equal(testingInstance, "Knee series", responseBody.Title)
	assert.True(testingInstance, responseBody.StudyCreated)
	assert.Equal(testingInstance, 2, responseBody.InstancesAccepted)
	require.Len(testingInstance, responseBody.InstancesRejected, 1)
	assert.Equal(testingInstance, "notes.txt", responseBody.InstancesRejected[0].FileName)
	assert.Equal(testingInstance, "INVALID_DICOM", responseBody.InstancesRejected[0].Code)
	assert.NotEmpty(testingInstance, responseBody.InstancesRejected[0].Reason)
}

func TestHTTPHandler_UploadPatientStudy_ZIPErrorShapes(testingInstance *testing.T) {
	testCases := []struct {
		name         string
		fileName     string
		fileBytes    func(*testing.T) []byte
		settings     func(*imaging.ServiceSettings)
		expectedCode int
		expectedText string
	}{
		{
			name:         "non zip payload",
			fileName:     "study.zip",
			fileBytes:    func(*testing.T) []byte { return []byte("this is not a zip") },
			expectedCode: http.StatusUnprocessableEntity,
			expectedText: "ZIP_INVALID",
		},
		{
			name:     "unsafe entry",
			fileName: "study.zip",
			fileBytes: func(testingInstance *testing.T) []byte {
				return buildZIPArchive(testingInstance, []zipEntry{{name: "../evil.dcm", content: instanceFixtureWithSOP("1").bytes()}})
			},
			expectedCode: http.StatusUnprocessableEntity,
			expectedText: "ZIP_UNSAFE_ENTRY",
		},
		{
			name:     "entry limit",
			fileName: "study.zip",
			fileBytes: func(testingInstance *testing.T) []byte {
				return buildZIPArchive(testingInstance, []zipEntry{
					{name: "a.dcm", content: instanceFixtureWithSOP("1").bytes()},
					{name: "b.dcm", content: instanceFixtureWithSOP("2").bytes()},
				})
			},
			settings:     func(settings *imaging.ServiceSettings) { settings.ArchiveLimits.MaxEntries = 1 },
			expectedCode: http.StatusRequestEntityTooLarge,
			expectedText: "ZIP_LIMITS_EXCEEDED",
		},
		{
			name:     "no valid instances",
			fileName: "study.zip",
			fileBytes: func(testingInstance *testing.T) []byte {
				return buildZIPArchive(testingInstance, []zipEntry{{name: "notes.txt", content: []byte("text")}})
			},
			expectedCode: http.StatusUnprocessableEntity,
			expectedText: "DICOM_NO_VALID_INSTANCES",
		},
	}

	for _, testCase := range testCases {
		testingInstance.Run(testCase.name, func(testingInstance *testing.T) {
			settings := imaging.DefaultServiceSettings()
			if testCase.settings != nil {
				testCase.settings(&settings)
			}
			harness := newUploadHarnessWithSettings(settings)
			handler := imaging.NewHTTPHandler(harness.service)
			recorder := httptest.NewRecorder()

			handler.UploadPatientStudy(recorder, buildUploadRequestWithFile(testingInstance, testCase.fileName, testCase.fileBytes(testingInstance), nil))

			require.Equal(testingInstance, testCase.expectedCode, recorder.Code)
			responseBody := decodeErrorBody(testingInstance, recorder)
			assert.Equal(testingInstance, testCase.expectedText, responseBody["code"])
			assert.NotEmpty(testingInstance, responseBody["error"])
		})
	}
}

func TestHTTPHandler_UploadPatientStudy_PatientMismatchInArchiveLeaksNoPatientData(testingInstance *testing.T) {
	harness := newUploadHarness()
	handler := imaging.NewHTTPHandler(harness.service)
	mismatched := instanceFixtureWithSOP("2")
	mismatched.patientID = "111.444.777-35"
	mismatched.patientName = "CONFIDENCIAL^PACIENTE"
	archiveBytes := buildZIPArchive(testingInstance, []zipEntry{
		{name: "ok.dcm", content: instanceFixtureWithSOP("1").bytes()},
		{name: "mismatch.dcm", content: mismatched.bytes()},
	})
	recorder := httptest.NewRecorder()

	handler.UploadPatientStudy(recorder, buildUploadRequestWithFile(testingInstance, "study.zip", archiveBytes, nil))

	require.Equal(testingInstance, http.StatusUnprocessableEntity, recorder.Code)
	assert.Equal(testingInstance, "DICOM_PATIENT_MISMATCH", decodeErrorBody(testingInstance, recorder)["code"])
	assert.NotContains(testingInstance, recorder.Body.String(), "CONFIDENCIAL")
	assert.NotContains(testingInstance, recorder.Body.String(), "111.444")
	assert.Empty(testingInstance, harness.storage.uploadedObjects)
}

func TestHTTPHandler_UploadPatientStudy_MissingFileWhenOnlyFieldsSent(testingInstance *testing.T) {
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

	require.Equal(testingInstance, http.StatusBadRequest, recorder.Code)
	assert.Equal(testingInstance, "Arquivo DICOM ausente.", decodeErrorBody(testingInstance, recorder)["error"])
}

func TestHTTPHandler_GetStudyDownloadURL_SingleInstanceReturnsURL(testingInstance *testing.T) {
	harness := newUploadHarness()
	result, uploadError := harness.upload(newDICOMFixture(), "", "")
	require.NoError(testingInstance, uploadError)
	handler := imaging.NewHTTPHandler(harness.service)
	request := httptest.NewRequest(http.MethodGet, "/api/v1/studies/"+result.Study.ID.String()+"/download-url", nil)
	request.SetPathValue("studyId", result.Study.ID.String())
	recorder := httptest.NewRecorder()

	handler.GetStudyDownloadURL(recorder, request)

	require.Equal(testingInstance, http.StatusOK, recorder.Code)
	var responseBody map[string]string
	require.NoError(testingInstance, json.Unmarshal(recorder.Body.Bytes(), &responseBody))
	assert.NotEmpty(testingInstance, responseBody["downloadUrl"])
	assert.NotEmpty(testingInstance, responseBody["expiresAt"])
}

func TestHTTPHandler_GetStudyDownloadURL_MultipleInstancesReturns409(testingInstance *testing.T) {
	harness := newUploadHarness()
	first, uploadError := harness.upload(instanceFixtureWithSOP("1"), "", "")
	require.NoError(testingInstance, uploadError)
	_, secondError := harness.upload(instanceFixtureWithSOP("2"), "", "")
	require.NoError(testingInstance, secondError)
	handler := imaging.NewHTTPHandler(harness.service)
	request := httptest.NewRequest(http.MethodGet, "/api/v1/studies/"+first.Study.ID.String()+"/download-url", nil)
	request.SetPathValue("studyId", first.Study.ID.String())
	recorder := httptest.NewRecorder()

	handler.GetStudyDownloadURL(recorder, request)

	require.Equal(testingInstance, http.StatusConflict, recorder.Code)
	responseBody := decodeErrorBody(testingInstance, recorder)
	assert.Equal(testingInstance, "MULTIPLE_INSTANCES", responseBody["code"])
	assert.Contains(testingInstance, responseBody["error"], "instances")
}

func TestHTTPHandler_ListStudyInstances_ReturnsEveryInstanceWithURL(testingInstance *testing.T) {
	harness := newUploadHarness()
	result, uploadError := uploadArchiveThrough(harness, buildZIPArchive(testingInstance, []zipEntry{
		{name: "a.dcm", content: instanceFixtureWithSOP("1").bytes()},
		{name: "b.dcm", content: instanceFixtureWithSOP("2").bytes()},
	}), "study.zip")
	require.NoError(testingInstance, uploadError)
	handler := imaging.NewHTTPHandler(harness.service)
	request := httptest.NewRequest(http.MethodGet, "/api/v1/studies/"+result.Study.ID.String()+"/instances", nil)
	request.SetPathValue("studyId", result.Study.ID.String())
	recorder := httptest.NewRecorder()

	handler.ListStudyInstances(recorder, request)

	require.Equal(testingInstance, http.StatusOK, recorder.Code)
	var responseBody []map[string]interface{}
	require.NoError(testingInstance, json.Unmarshal(recorder.Body.Bytes(), &responseBody))
	require.Len(testingInstance, responseBody, 2)
	for _, instanceEntry := range responseBody {
		for _, requiredField := range []string{"id", "sopInstanceUid", "seriesInstanceUid", "instanceNumber", "fileName", "downloadUrl", "expiresAt"} {
			assert.Contains(testingInstance, instanceEntry, requiredField)
		}
		assert.NotEmpty(testingInstance, instanceEntry["downloadUrl"])
	}
}

func TestHTTPHandler_ListStudyInstances_UnknownStudyReturns404(testingInstance *testing.T) {
	harness := newUploadHarness()
	handler := imaging.NewHTTPHandler(harness.service)
	unknownStudyID := "00000000-0000-0000-0000-000000000000"
	request := httptest.NewRequest(http.MethodGet, "/api/v1/studies/"+unknownStudyID+"/instances", nil)
	request.SetPathValue("studyId", unknownStudyID)
	recorder := httptest.NewRecorder()

	handler.ListStudyInstances(recorder, request)

	assert.Equal(testingInstance, http.StatusNotFound, recorder.Code)
}
