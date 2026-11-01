package tests

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"errors"
	"hash/crc32"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/healthcare/backend/internal/modules/exam_analyzer"
	"github.com/healthcare/backend/internal/modules/exam_analyzer/mocks"
	"github.com/healthcare/backend/internal/shared/apperrors"
	"github.com/healthcare/backend/internal/shared/filecontent"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	secretPatientName = "Maria Secreta Silva"
	originalFileName  = "maria_secreta_raio_x.png"
)

type recordingStorage struct {
	mutex   sync.Mutex
	objects map[string][]byte
}

func newRecordingStorage() *recordingStorage {
	return &recordingStorage{objects: map[string][]byte{}}
}

func (storageClient *recordingStorage) Upload(ctx context.Context, key string, reader io.Reader, contentType string) error {
	content, readErr := io.ReadAll(reader)
	if readErr != nil {
		return readErr
	}
	storageClient.mutex.Lock()
	defer storageClient.mutex.Unlock()
	storageClient.objects[key] = content
	return nil
}

func (storageClient *recordingStorage) GetPresignedURL(ctx context.Context, key string, expiration time.Duration) (string, error) {
	return "", nil
}

func (storageClient *recordingStorage) Delete(ctx context.Context, key string) error {
	storageClient.mutex.Lock()
	defer storageClient.mutex.Unlock()
	delete(storageClient.objects, key)
	return nil
}

func (storageClient *recordingStorage) Read(ctx context.Context, key string) (io.ReadCloser, error) {
	storageClient.mutex.Lock()
	defer storageClient.mutex.Unlock()
	content, exists := storageClient.objects[key]
	if !exists {
		return nil, errors.New("object not found")
	}
	return io.NopCloser(bytes.NewReader(content)), nil
}

func (storageClient *recordingStorage) objectCount() int {
	storageClient.mutex.Lock()
	defer storageClient.mutex.Unlock()
	return len(storageClient.objects)
}

func (storageClient *recordingStorage) firstObject() []byte {
	storageClient.mutex.Lock()
	defer storageClient.mutex.Unlock()
	for _, content := range storageClient.objects {
		return content
	}
	return nil
}

type analysisCall struct {
	fileName  string
	fileBytes []byte
}

type spyService struct {
	exam_analyzer.Service
	mutex sync.Mutex
	calls []analysisCall
}

func (spy *spyService) AnalyzeExamFile(ctx context.Context, fileName string, fileBytes []byte) (*exam_analyzer.MedicalAnalysisResponse, string, error) {
	spy.mutex.Lock()
	defer spy.mutex.Unlock()
	spy.calls = append(spy.calls, analysisCall{fileName: fileName, fileBytes: fileBytes})
	return &exam_analyzer.MedicalAnalysisResponse{ExamType: "Teste"}, "completed", nil
}

func (spy *spyService) recordedCalls() []analysisCall {
	spy.mutex.Lock()
	defer spy.mutex.Unlock()
	return append([]analysisCall{}, spy.calls...)
}

func buildPNGChunk(chunkType string, data []byte) []byte {
	chunk := binary.BigEndian.AppendUint32(nil, uint32(len(data)))
	chunk = append(chunk, chunkType...)
	chunk = append(chunk, data...)
	return binary.BigEndian.AppendUint32(chunk, crc32.ChecksumIEEE(append([]byte(chunkType), data...)))
}

func buildPNGWithPatientText(testingInstance *testing.T) []byte {
	testingInstance.Helper()
	canvas := image.NewRGBA(image.Rect(0, 0, 8, 8))
	for row := 0; row < 8; row++ {
		for column := 0; column < 8; column++ {
			canvas.Set(column, row, color.RGBA{R: uint8(row * 30), G: uint8(column * 30), B: 90, A: 255})
		}
	}
	var buffer bytes.Buffer
	require.NoError(testingInstance, png.Encode(&buffer, canvas))
	plain := buffer.Bytes()
	endOffset := len(plain) - 12

	var enriched bytes.Buffer
	enriched.Write(plain[:endOffset])
	enriched.Write(buildPNGChunk("tEXt", []byte("Patient\x00"+secretPatientName)))
	enriched.Write(plain[endOffset:])
	return enriched.Bytes()
}

func buildPlainJPEG(testingInstance *testing.T) []byte {
	testingInstance.Helper()
	var buffer bytes.Buffer
	require.NoError(testingInstance, jpeg.Encode(&buffer, image.NewRGBA(image.Rect(0, 0, 8, 8)), nil))
	return buffer.Bytes()
}

func buildMinimalPDF() []byte {
	return []byte("%PDF-1.4\n1 0 obj\n<< /Type /Catalog >>\nendobj\ntrailer\n<< /Root 1 0 R >>\n%%EOF\n")
}

func buildMultipartUpload(testingInstance *testing.T, fileName string, fileContent []byte, fields map[string]string) (*bytes.Buffer, string) {
	testingInstance.Helper()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	for fieldName, fieldValue := range fields {
		require.NoError(testingInstance, writer.WriteField(fieldName, fieldValue))
	}
	filePart, partErr := writer.CreateFormFile("file", fileName)
	require.NoError(testingInstance, partErr)
	_, writeErr := filePart.Write(fileContent)
	require.NoError(testingInstance, writeErr)
	require.NoError(testingInstance, writer.Close())
	return body, writer.FormDataContentType()
}

type uploadFixture struct {
	repository *mocks.MockExamAnalysisRepository
	storage    *recordingStorage
	handler    *exam_analyzer.HTTPHandler
	service    exam_analyzer.Service
	worker     *exam_analyzer.Worker
}

func newUploadFixture() uploadFixture {
	repository := mocks.NewMockExamAnalysisRepository()
	storageClient := newRecordingStorage()
	service := exam_analyzer.NewService(repository, storageClient, "", "", "")
	worker := exam_analyzer.NewWorker(repository, service, storageClient, nil)
	return uploadFixture{
		repository: repository,
		storage:    storageClient,
		handler:    exam_analyzer.NewHTTPHandler(repository, service, worker),
		service:    service,
		worker:     worker,
	}
}

func (fixture uploadFixture) upload(testingInstance *testing.T, fileName string, fileContent []byte, fields map[string]string) *httptest.ResponseRecorder {
	testingInstance.Helper()
	body, contentType := buildMultipartUpload(testingInstance, fileName, fileContent, fields)
	httpRequest := httptest.NewRequest(http.MethodPost, "/api/v1/exam-analyses", body)
	httpRequest.Header.Set("Content-Type", contentType)
	recorder := httptest.NewRecorder()
	fixture.handler.CreateAnalysis(recorder, httpRequest)
	return recorder
}

func decodeBody(testingInstance *testing.T, recorder *httptest.ResponseRecorder) map[string]any {
	testingInstance.Helper()
	var decoded map[string]any
	require.NoError(testingInstance, json.Unmarshal(recorder.Body.Bytes(), &decoded))
	return decoded
}

func TestHTTPCreateAnalysis_RejectsContentExtensionMismatch(testingInstance *testing.T) {
	testCases := []struct {
		name         string
		fileName     string
		content      []byte
		expectedCode string
	}{
		{"png that is pdf", "scan.png", buildMinimalPDF(), "FILE_CONTENT_MISMATCH"},
		{"pdf that is jpeg", "laudo.pdf", buildPlainJPEG(testingInstance), "FILE_CONTENT_MISMATCH"},
		{"empty file", "scan.png", []byte{}, "FILE_EMPTY"},
		{"text renamed to png", "scan.png", []byte("plain text, not an image"), "FILE_CORRUPTED"},
		{"truncated jpeg", "scan.jpg", buildPlainJPEG(testingInstance)[:20], "FILE_CORRUPTED"},
	}

	for _, testCase := range testCases {
		testingInstance.Run(testCase.name, func(testingInstance *testing.T) {
			fixture := newUploadFixture()

			recorder := fixture.upload(testingInstance, testCase.fileName, testCase.content, map[string]string{
				"consent":   "true",
				"anonymize": "true",
			})

			assert.Equal(testingInstance, http.StatusUnprocessableEntity, recorder.Code)
			assert.Equal(testingInstance, testCase.expectedCode, decodeBody(testingInstance, recorder)["code"])
			assert.Equal(testingInstance, 0, fixture.storage.objectCount())
			assert.Empty(testingInstance, fixture.repository.Analyses)
		})
	}
}

func TestHTTPCreateAnalysis_AcknowledgementRuleUnchanged(testingInstance *testing.T) {
	fixture := newUploadFixture()
	pngContent := buildPNGWithPatientText(testingInstance)

	rejected := fixture.upload(testingInstance, "scan.png", pngContent, map[string]string{"consent": "true", "anonymize": "false"})
	assert.Equal(testingInstance, http.StatusUnprocessableEntity, rejected.Code)
	assert.Equal(testingInstance, apperrors.ErrDeidentificationRequired.Message, decodeBody(testingInstance, rejected)["error"])
	assert.Empty(testingInstance, fixture.repository.Analyses)

	accepted := fixture.upload(testingInstance, "scan.png", pngContent, map[string]string{
		"consent":                       "true",
		"anonymize":                     "false",
		"acknowledgeExternalProcessing": "true",
	})
	assert.Equal(testingInstance, http.StatusCreated, accepted.Code)
	body := decodeBody(testingInstance, accepted)
	assert.Equal(testingInstance, false, body["anonymized"])
	assert.Nil(testingInstance, body["anonymizationReport"])
}

func TestHTTPCreateAnalysis_AnonymizedPersistsAndReturnsReport(testingInstance *testing.T) {
	fixture := newUploadFixture()
	original := buildPNGWithPatientText(testingInstance)

	recorder := fixture.upload(testingInstance, originalFileName, original, map[string]string{"consent": "true", "anonymize": "true"})

	require.Equal(testingInstance, http.StatusCreated, recorder.Code)
	body := decodeBody(testingInstance, recorder)
	assert.Equal(testingInstance, true, body["anonymized"])

	report, isObject := body["anonymizationReport"].(map[string]any)
	require.True(testingInstance, isObject, "anonymizationReport must be an object")
	assert.Equal(testingInstance, true, report["applied"])
	assert.ElementsMatch(testingInstance, []any{filecontent.StepMetadataRemoved, filecontent.StepFileNameNeutralized}, report["steps"])
	assert.Equal(testingInstance, []any{filecontent.LimitationBurnedInText}, report["limitations"])

	require.Len(testingInstance, fixture.repository.Analyses, 1)
	for _, persisted := range fixture.repository.Analyses {
		require.NotNil(testingInstance, persisted.AnonymizationReport)
		assert.True(testingInstance, persisted.AnonymizationReport.Applied)
		assert.Equal(testingInstance, "maria_secreta_raio_x.png", persisted.FileName)
	}

	assert.Equal(testingInstance, original, fixture.storage.firstObject(), "storage keeps the original")
}

func TestService_CreateAnalysis_ReturnsAppErrors(testingInstance *testing.T) {
	fixture := newUploadFixture()
	analysisID := uuid.New()

	mismatchErr := fixture.service.CreateAnalysis(context.Background(), &exam_analyzer.ExamAnalysis{ID: analysisID, FileName: "a.png"}, bytes.NewReader(buildMinimalPDF()))
	assert.Equal(testingInstance, apperrors.ErrFileContentMismatch, mismatchErr)

	emptyErr := fixture.service.CreateAnalysis(context.Background(), &exam_analyzer.ExamAnalysis{ID: analysisID, FileName: "a.png"}, strings.NewReader(""))
	assert.Equal(testingInstance, apperrors.ErrFileEmpty, emptyErr)

	oversized := bytes.NewReader(append(buildPNGWithPatientText(testingInstance), make([]byte, filecontent.MaxFileBytes)...))
	oversizedErr := fixture.service.CreateAnalysis(context.Background(), &exam_analyzer.ExamAnalysis{ID: analysisID, FileName: "a.png"}, oversized)
	assert.Equal(testingInstance, apperrors.ErrPayloadTooLarge, oversizedErr)

	assert.Equal(testingInstance, 0, fixture.storage.objectCount())
}

func runWorkerUntilIdle(testingInstance *testing.T, fixture uploadFixture, analysisID uuid.UUID) {
	testingInstance.Helper()
	workerContext, cancel := context.WithCancel(context.Background())
	defer cancel()
	go fixture.worker.Start(workerContext)

	require.Eventually(testingInstance, func() bool {
		record, fetchErr := fixture.repository.GetAnalysis(context.Background(), analysisID)
		return fetchErr == nil && record.Status == "completed"
	}, 5*time.Second, 20*time.Millisecond)
}

func TestWorker_SendsNeutralNameAndStrippedBytesWhenAnonymized(testingInstance *testing.T) {
	repository := mocks.NewMockExamAnalysisRepository()
	storageClient := newRecordingStorage()
	spy := &spyService{Service: exam_analyzer.NewService(repository, storageClient, "", "", "")}
	worker := exam_analyzer.NewWorker(repository, spy, storageClient, nil)
	fixture := uploadFixture{repository: repository, storage: storageClient, service: spy, worker: worker, handler: exam_analyzer.NewHTTPHandler(repository, spy, worker)}
	original := buildPNGWithPatientText(testingInstance)

	recorder := fixture.upload(testingInstance, originalFileName, original, map[string]string{
		"consent":       "true",
		"anonymize":     "true",
		"patientFhirId": "patient-12345",
	})
	require.Equal(testingInstance, http.StatusCreated, recorder.Code)

	var analysisID uuid.UUID
	for identifier := range repository.Analyses {
		analysisID = identifier
	}
	runWorkerUntilIdle(testingInstance, fixture, analysisID)

	calls := spy.recordedCalls()
	require.Len(testingInstance, calls, 1)
	assert.Equal(testingInstance, "exam.png", calls[0].fileName)
	assert.NotContains(testingInstance, string(calls[0].fileBytes), secretPatientName)
	assert.NotContains(testingInstance, calls[0].fileName, "maria")
	assert.NotContains(testingInstance, calls[0].fileName, "patient-12345")
	assert.NoError(testingInstance, filecontent.Validate(filecontent.KindPNG, "exam.png", calls[0].fileBytes))
	assert.Contains(testingInstance, string(original), secretPatientName)

	record, fetchErr := repository.GetAnalysis(context.Background(), analysisID)
	require.NoError(testingInstance, fetchErr)
	require.NotNil(testingInstance, record.AnonymizationReport)
	assert.Contains(testingInstance, record.AnonymizationReport.Limitations, filecontent.LimitationBurnedInText)
	assert.Equal(testingInstance, originalFileName, record.FileName)
}

func TestWorker_UsesOriginalNameAndBytesWhenNotAnonymized(testingInstance *testing.T) {
	repository := mocks.NewMockExamAnalysisRepository()
	storageClient := newRecordingStorage()
	spy := &spyService{Service: exam_analyzer.NewService(repository, storageClient, "", "", "")}
	worker := exam_analyzer.NewWorker(repository, spy, storageClient, nil)
	fixture := uploadFixture{repository: repository, storage: storageClient, service: spy, worker: worker, handler: exam_analyzer.NewHTTPHandler(repository, spy, worker)}
	original := buildPNGWithPatientText(testingInstance)

	recorder := fixture.upload(testingInstance, originalFileName, original, map[string]string{
		"consent":                       "true",
		"anonymize":                     "false",
		"acknowledgeExternalProcessing": "true",
	})
	require.Equal(testingInstance, http.StatusCreated, recorder.Code)

	var analysisID uuid.UUID
	for identifier := range repository.Analyses {
		analysisID = identifier
	}
	runWorkerUntilIdle(testingInstance, fixture, analysisID)

	calls := spy.recordedCalls()
	require.Len(testingInstance, calls, 1)
	assert.Equal(testingInstance, originalFileName, calls[0].fileName)
	assert.Equal(testingInstance, original, calls[0].fileBytes)

	record, fetchErr := repository.GetAnalysis(context.Background(), analysisID)
	require.NoError(testingInstance, fetchErr)
	assert.Nil(testingInstance, record.AnonymizationReport)
}
