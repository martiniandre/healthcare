package integration

import (
	"bytes"
	"context"
	"encoding/binary"
	"mime/multipart"
	"net/http"
	"testing"
)

type minimalDICOM struct {
	studyInstanceUID string
	patientID        string
	patientName      string
}

func (instance minimalDICOM) bytes() []byte {
	var fileBuffer bytes.Buffer
	fileBuffer.Write(make([]byte, 128))
	fileBuffer.WriteString("DICM")

	writeExplicitElement(&fileBuffer, 0x0002, 0x0002, "UI", "1.2.840.10008.5.1.4.1.1.4")
	writeExplicitElement(&fileBuffer, 0x0002, 0x0003, "UI", "1.2.826.0.1.3680043.8.498.9001")
	writeExplicitElement(&fileBuffer, 0x0002, 0x0010, "UI", "1.2.840.10008.1.2.1")
	writeExplicitElement(&fileBuffer, 0x0008, 0x0020, "DA", "20260115")
	writeExplicitElement(&fileBuffer, 0x0008, 0x0030, "TM", "101500")
	writeExplicitElement(&fileBuffer, 0x0008, 0x0050, "SH", "ACC-9")
	writeExplicitElement(&fileBuffer, 0x0008, 0x0060, "CS", "MR")
	writeExplicitElement(&fileBuffer, 0x0008, 0x1030, "LO", "KNEE")
	writeExplicitElement(&fileBuffer, 0x0010, 0x0010, "PN", instance.patientName)
	writeExplicitElement(&fileBuffer, 0x0010, 0x0020, "LO", instance.patientID)
	writeExplicitElement(&fileBuffer, 0x0020, 0x000D, "UI", instance.studyInstanceUID)
	writeExplicitElement(&fileBuffer, 0x0020, 0x000E, "UI", "1.2.826.0.1.3680043.8.498.9002")
	return fileBuffer.Bytes()
}

func writeExplicitElement(targetBuffer *bytes.Buffer, groupNumber uint16, elementNumber uint16, representation string, value string) {
	tagBytes := make([]byte, 4)
	binary.LittleEndian.PutUint16(tagBytes[0:2], groupNumber)
	binary.LittleEndian.PutUint16(tagBytes[2:4], elementNumber)
	targetBuffer.Write(tagBytes)
	targetBuffer.WriteString(representation)

	valueBytes := []byte(value)
	if len(valueBytes)%2 == 1 {
		valueBytes = append(valueBytes, 0)
	}
	lengthBytes := make([]byte, 2)
	binary.LittleEndian.PutUint16(lengthBytes, uint16(len(valueBytes)))
	targetBuffer.Write(lengthBytes)
	targetBuffer.Write(valueBytes)
}

func uploadStudy(t *testing.T, client *authenticatedClient, patientFhirID string, fileBytes []byte) *http.Response {
	t.Helper()

	var bodyBuffer bytes.Buffer
	multipartWriter := multipart.NewWriter(&bodyBuffer)
	filePart, partError := multipartWriter.CreateFormFile("file", "knee.dcm")
	if partError != nil {
		t.Fatalf("failed to create multipart file part: %v", partError)
	}
	if _, writeError := filePart.Write(fileBytes); writeError != nil {
		t.Fatalf("failed to write multipart file part: %v", writeError)
	}
	if closeError := multipartWriter.Close(); closeError != nil {
		t.Fatalf("failed to close multipart writer: %v", closeError)
	}

	request, requestError := http.NewRequest(http.MethodPost, client.serverURL+"/api/v1/patients/"+patientFhirID+"/studies", &bodyBuffer)
	if requestError != nil {
		t.Fatalf("failed to build upload request: %v", requestError)
	}
	request.Header.Set("Content-Type", multipartWriter.FormDataContentType())
	request.Header.Set("X-CSRF-Token", client.csrfToken)

	response, doError := client.httpClient.Do(request)
	if doError != nil {
		t.Fatalf("upload request failed: %v", doError)
	}
	return response
}

func resolveSeededPatientFhirID(t *testing.T, testServer *testServer) string {
	t.Helper()
	var patientFhirID string
	if queryError := testServer.db.QueryRow(context.Background(), `SELECT id::text FROM users WHERE email = 'paciente@mail.com'`).Scan(&patientFhirID); queryError != nil {
		t.Fatalf("failed to resolve seeded patient id: %v", queryError)
	}
	return patientFhirID
}

func countRows(t *testing.T, testServer *testServer, query string) int {
	t.Helper()
	var rowCount int
	if queryError := testServer.db.QueryRow(context.Background(), query).Scan(&rowCount); queryError != nil {
		t.Fatalf("failed to count rows: %v", queryError)
	}
	return rowCount
}

func TestImagingUploadFillsMetadataFromDICOMHeader(t *testing.T) {
	testServer := newTestServer(t)
	serverURL := startTestHTTPServer(t, testServer.handler)
	doctorClient := loginAs(t, serverURL, "medico@clinica.com", "secret123")
	patientFhirID := resolveSeededPatientFhirID(t, testServer)

	dicomFile := minimalDICOM{studyInstanceUID: "1.2.826.0.1.3680043.8.498.9000", patientID: "529.982.247-25", patientName: "SILVA^MARIA"}
	response := uploadStudy(t, doctorClient, patientFhirID, dicomFile.bytes())
	requireStatusCode(t, response, http.StatusCreated)

	var uploadBody struct {
		ID                string `json:"id"`
		Title             string `json:"title"`
		Modality          string `json:"modality"`
		Status            string `json:"status"`
		StudyInstanceUid  string `json:"studyInstanceUid"`
		StudyDate         string `json:"studyDate"`
		AccessionNumber   string `json:"accessionNumber"`
		SeriesInstanceUid string `json:"seriesInstanceUid"`
		SopInstanceUid    string `json:"sopInstanceUid"`
		Warnings          []struct {
			Code string `json:"code"`
		} `json:"warnings"`
	}
	decodeJSONResponse(t, response, &uploadBody)

	if uploadBody.StudyInstanceUid != dicomFile.studyInstanceUID || uploadBody.Modality != "MR" || uploadBody.Title != "MR - KNEE" {
		t.Fatalf("unexpected upload body: %+v", uploadBody)
	}
	if uploadBody.Status != "RECEIVING" || uploadBody.StudyDate != "2026-01-15T10:15:00Z" || uploadBody.AccessionNumber != "ACC-9" {
		t.Fatalf("unexpected upload body: %+v", uploadBody)
	}
	if uploadBody.SeriesInstanceUid != "1.2.826.0.1.3680043.8.498.9002" || uploadBody.SopInstanceUid != "1.2.826.0.1.3680043.8.498.9001" {
		t.Fatalf("unexpected series or sop uid: %+v", uploadBody)
	}
	if len(uploadBody.Warnings) != 0 {
		t.Fatalf("expected no warnings, got %+v", uploadBody.Warnings)
	}

	if countRows(t, testServer, `SELECT COUNT(*) FROM imaging_studies WHERE study_instance_uid = '1.2.826.0.1.3680043.8.498.9000' AND accession_number = 'ACC-9' AND study_datetime IS NOT NULL AND total_instances = 1`) != 1 {
		t.Fatal("expected the study row with header metadata to be persisted")
	}
	if countRows(t, testServer, `SELECT COUNT(*) FROM imaging_series`) != 1 || countRows(t, testServer, `SELECT COUNT(*) FROM imaging_instances WHERE ingest_status = 'STORED'`) != 1 {
		t.Fatal("expected one series and one stored instance to be indexed")
	}

	requireStatusCode(t, doctorClient.Get(t, "/api/v1/studies/"+uploadBody.ID), http.StatusOK)
}

func TestImagingUploadRejectsFileFromAnotherPatient(t *testing.T) {
	testServer := newTestServer(t)
	serverURL := startTestHTTPServer(t, testServer.handler)
	doctorClient := loginAs(t, serverURL, "medico@clinica.com", "secret123")
	patientFhirID := resolveSeededPatientFhirID(t, testServer)

	dicomFile := minimalDICOM{studyInstanceUID: "1.2.826.0.1.3680043.8.498.9100", patientID: "111.444.777-35", patientName: "SANTOS^PEDRO"}
	response := uploadStudy(t, doctorClient, patientFhirID, dicomFile.bytes())
	requireStatusCode(t, response, http.StatusUnprocessableEntity)

	var errorBody map[string]string
	decodeJSONResponse(t, response, &errorBody)
	if errorBody["code"] != "DICOM_PATIENT_MISMATCH" {
		t.Fatalf("expected DICOM_PATIENT_MISMATCH, got %v", errorBody)
	}
	if countRows(t, testServer, `SELECT COUNT(*) FROM imaging_studies`) != 0 {
		t.Fatal("a rejected upload must not create a study")
	}
}

func TestImagingUploadRejectsDuplicateStudy(t *testing.T) {
	testServer := newTestServer(t)
	serverURL := startTestHTTPServer(t, testServer.handler)
	doctorClient := loginAs(t, serverURL, "medico@clinica.com", "secret123")
	patientFhirID := resolveSeededPatientFhirID(t, testServer)

	dicomFile := minimalDICOM{studyInstanceUID: "1.2.826.0.1.3680043.8.498.9200", patientID: "529.982.247-25", patientName: "SILVA^MARIA"}
	requireStatusCode(t, uploadStudy(t, doctorClient, patientFhirID, dicomFile.bytes()), http.StatusCreated)

	duplicateResponse := uploadStudy(t, doctorClient, patientFhirID, dicomFile.bytes())
	requireStatusCode(t, duplicateResponse, http.StatusConflict)
	var errorBody map[string]string
	decodeJSONResponse(t, duplicateResponse, &errorBody)
	if errorBody["code"] != "DICOM_DUPLICATE_INSTANCE" {
		t.Fatalf("expected DICOM_DUPLICATE_INSTANCE, got %v", errorBody)
	}
	if countRows(t, testServer, `SELECT COUNT(*) FROM imaging_studies`) != 1 {
		t.Fatal("a duplicate upload must not create a second study")
	}
}

func TestImagingUploadWarnsWhenDICOMIsAnonymized(t *testing.T) {
	testServer := newTestServer(t)
	serverURL := startTestHTTPServer(t, testServer.handler)
	doctorClient := loginAs(t, serverURL, "medico@clinica.com", "secret123")
	patientFhirID := resolveSeededPatientFhirID(t, testServer)

	dicomFile := minimalDICOM{studyInstanceUID: "1.2.826.0.1.3680043.8.498.9300", patientID: "ANON-17", patientName: "ANONYMOUS"}
	response := uploadStudy(t, doctorClient, patientFhirID, dicomFile.bytes())
	requireStatusCode(t, response, http.StatusCreated)

	var uploadBody struct {
		Warnings []struct {
			Code string `json:"code"`
		} `json:"warnings"`
	}
	decodeJSONResponse(t, response, &uploadBody)
	if len(uploadBody.Warnings) != 1 || uploadBody.Warnings[0].Code != "anonymized_study" {
		t.Fatalf("expected the anonymized_study warning, got %+v", uploadBody.Warnings)
	}
}
