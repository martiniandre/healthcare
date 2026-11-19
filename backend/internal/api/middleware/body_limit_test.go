package middleware

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func readBodyHandler(capturedError *error) http.Handler {
	return http.HandlerFunc(func(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
		_, readError := io.ReadAll(httpRequest.Body)
		*capturedError = readError
		httpResponseWriter.WriteHeader(http.StatusOK)
	})
}

func TestBodyLimitRejectsDeclaredOversizedJSONBody(testingInstance *testing.T) {
	var readError error
	handler := BodyLimit(8)(readBodyHandler(&readError))
	request := httptest.NewRequest(http.MethodPost, "/api/v1/patients", strings.NewReader(`{"name":"too long"}`))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	assert.Equal(testingInstance, http.StatusRequestEntityTooLarge, recorder.Code)
}

func TestBodyLimitStopsReadingStreamedOversizedJSONBody(testingInstance *testing.T) {
	var readError error
	handler := BodyLimit(8)(readBodyHandler(&readError))
	request := httptest.NewRequest(http.MethodPost, "/api/v1/patients", io.NopCloser(strings.NewReader(`{"name":"too long"}`)))
	request.ContentLength = -1
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	assert.Error(testingInstance, readError)
}

func TestBodyLimitAllowsSmallJSONBody(testingInstance *testing.T) {
	var readError error
	handler := BodyLimit(1024)(readBodyHandler(&readError))
	request := httptest.NewRequest(http.MethodPost, "/api/v1/patients", strings.NewReader(`{"name":"ok"}`))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	assert.Equal(testingInstance, http.StatusOK, recorder.Code)
	assert.NoError(testingInstance, readError)
}

func TestBodyLimitSkipsMultipartUploads(testingInstance *testing.T) {
	var readError error
	handler := BodyLimit(8)(readBodyHandler(&readError))
	request := httptest.NewRequest(http.MethodPost, "/api/v1/imaging/upload", strings.NewReader(strings.Repeat("x", 64)))
	request.Header.Set("Content-Type", "multipart/form-data; boundary=abc")
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	assert.Equal(testingInstance, http.StatusOK, recorder.Code)
	assert.NoError(testingInstance, readError)
}
