package integration

import (
	"bytes"
	"encoding/binary"
	"hash/crc32"
	"image"
	"image/png"
	"mime/multipart"
	"net/http"
	"testing"
)

const examSecretText = "Paciente Confidencial"

func buildExamPNGWithTextChunk(t *testing.T) []byte {
	t.Helper()
	var encoded bytes.Buffer
	if encodeError := png.Encode(&encoded, image.NewRGBA(image.Rect(0, 0, 8, 8))); encodeError != nil {
		t.Fatalf("failed to encode png: %v", encodeError)
	}
	plain := encoded.Bytes()
	endOffset := len(plain) - 12

	chunkData := []byte("Patient\x00" + examSecretText)
	chunk := binary.BigEndian.AppendUint32(nil, uint32(len(chunkData)))
	chunk = append(chunk, "tEXt"...)
	chunk = append(chunk, chunkData...)
	chunk = binary.BigEndian.AppendUint32(chunk, crc32.ChecksumIEEE(append([]byte("tEXt"), chunkData...)))

	var enriched bytes.Buffer
	enriched.Write(plain[:endOffset])
	enriched.Write(chunk)
	enriched.Write(plain[endOffset:])
	return enriched.Bytes()
}

func postExamUpload(t *testing.T, client *authenticatedClient, fileName string, fileContent []byte, fields map[string]string) *http.Response {
	t.Helper()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	for fieldName, fieldValue := range fields {
		if writeError := writer.WriteField(fieldName, fieldValue); writeError != nil {
			t.Fatalf("failed to write field: %v", writeError)
		}
	}
	filePart, partError := writer.CreateFormFile("file", fileName)
	if partError != nil {
		t.Fatalf("failed to create form file: %v", partError)
	}
	if _, writeError := filePart.Write(fileContent); writeError != nil {
		t.Fatalf("failed to write file: %v", writeError)
	}
	if closeError := writer.Close(); closeError != nil {
		t.Fatalf("failed to close multipart writer: %v", closeError)
	}

	request, requestError := http.NewRequest(http.MethodPost, client.serverURL+"/api/v1/exam-analyses", &body)
	if requestError != nil {
		t.Fatalf("failed to build request: %v", requestError)
	}
	request.Header.Set("Content-Type", writer.FormDataContentType())
	request.Header.Set("X-CSRF-Token", client.csrfToken)

	response, doError := client.httpClient.Do(request)
	if doError != nil {
		t.Fatalf("exam upload request failed: %v", doError)
	}
	return response
}

func TestExamAnalysisRejectsContentThatDoesNotMatchExtension(t *testing.T) {
	testServer := newTestServer(t)
	serverURL := startTestHTTPServer(t, testServer.handler)
	doctorClient := loginAs(t, serverURL, "medico@clinica.com", "secret123")

	response := postExamUpload(t, doctorClient, "scan.png", []byte("%PDF-1.4\n1 0 obj\n<<>>\nendobj\n%%EOF\n"), map[string]string{
		"consent":   "true",
		"anonymize": "true",
	})

	requireStatusCode(t, response, http.StatusUnprocessableEntity)
	if code := requireJSONFieldValue(t, response, "code"); code != "FILE_CONTENT_MISMATCH" {
		t.Fatalf("expected FILE_CONTENT_MISMATCH, got %s", code)
	}
}

func TestExamAnalysisAnonymizationReportRoundTrip(t *testing.T) {
	testServer := newTestServer(t)
	serverURL := startTestHTTPServer(t, testServer.handler)
	doctorClient := loginAs(t, serverURL, "medico@clinica.com", "secret123")

	createResponse := postExamUpload(t, doctorClient, "exame.png", buildExamPNGWithTextChunk(t), map[string]string{
		"consent":   "true",
		"anonymize": "true",
	})
	requireStatusCode(t, createResponse, http.StatusCreated)
	var created map[string]interface{}
	decodeJSONResponse(t, createResponse, &created)

	analysisID, _ := created["id"].(string)
	if analysisID == "" {
		t.Fatalf("created analysis has no id: %v", created)
	}
	assertAnonymizationReport(t, created["anonymizationReport"])

	getResponse := doctorClient.Get(t, "/api/v1/exam-analyses/"+analysisID)
	requireStatusCode(t, getResponse, http.StatusOK)
	var fetched map[string]interface{}
	decodeJSONResponse(t, getResponse, &fetched)
	assertAnonymizationReport(t, fetched["anonymizationReport"])

	listResponse := doctorClient.Get(t, "/api/v1/exam-analyses")
	requireStatusCode(t, listResponse, http.StatusOK)
	var listed []map[string]interface{}
	decodeJSONResponse(t, listResponse, &listed)
	foundInList := false
	for _, listedAnalysis := range listed {
		if listedAnalysis["id"] == analysisID {
			foundInList = true
			assertAnonymizationReport(t, listedAnalysis["anonymizationReport"])
		}
	}
	if !foundInList {
		t.Fatal("created analysis not present in list")
	}

	acknowledgedResponse := postExamUpload(t, doctorClient, "exame.png", buildExamPNGWithTextChunk(t), map[string]string{
		"consent":                       "true",
		"anonymize":                     "false",
		"acknowledgeExternalProcessing": "true",
	})
	requireStatusCode(t, acknowledgedResponse, http.StatusCreated)
	var notAnonymized map[string]interface{}
	decodeJSONResponse(t, acknowledgedResponse, &notAnonymized)
	if notAnonymized["anonymizationReport"] != nil {
		t.Fatalf("expected null anonymizationReport when not anonymized, got %v", notAnonymized["anonymizationReport"])
	}
}

func assertAnonymizationReport(t *testing.T, rawReport interface{}) {
	t.Helper()
	report, isObject := rawReport.(map[string]interface{})
	if !isObject {
		t.Fatalf("expected anonymizationReport object, got %v", rawReport)
	}
	if report["applied"] != true {
		t.Fatalf("expected applied true, got %v", report["applied"])
	}
	limitations, _ := report["limitations"].([]interface{})
	hasBurnedInLimitation := false
	for _, limitation := range limitations {
		if limitation == "burned_in_text_not_removed" {
			hasBurnedInLimitation = true
		}
	}
	if !hasBurnedInLimitation {
		t.Fatalf("expected burned_in_text_not_removed limitation, got %v", limitations)
	}
	steps, _ := report["steps"].([]interface{})
	if len(steps) == 0 {
		t.Fatalf("unexpected steps: %v", steps)
	}
}
