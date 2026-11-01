package dicomstore

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestClient(testingInstance *testing.T, handler http.HandlerFunc) (*Client, *httptest.Server) {
	testingInstance.Helper()
	testServer := httptest.NewServer(handler)
	testingInstance.Cleanup(testServer.Close)

	storeClient, clientError := NewClientWithBaseURL(testServer.URL + "/dicomWeb")
	require.NoError(testingInstance, clientError)
	return storeClient, testServer
}

func TestStoreInstanceReportsStoredOutcome(testingInstance *testing.T) {
	storeClient, _ := newTestClient(testingInstance, func(responseWriter http.ResponseWriter, httpRequest *http.Request) {
		assert.Equal(testingInstance, http.MethodPost, httpRequest.Method)
		assert.Equal(testingInstance, "/dicomWeb/studies", httpRequest.URL.Path)
		assert.Equal(testingInstance, dicomContentType, httpRequest.Header.Get("Content-Type"))

		bodyBytes, bodyError := io.ReadAll(httpRequest.Body)
		require.NoError(testingInstance, bodyError)
		assert.Equal(testingInstance, "DICM-PAYLOAD", string(bodyBytes))

		responseWriter.WriteHeader(http.StatusOK)
		responseWriter.Write([]byte(`{}`))
	})

	outcome, storeError := storeClient.StoreInstance(context.Background(), bytes.NewReader([]byte("DICM-PAYLOAD")), 12)

	require.NoError(testingInstance, storeError)
	assert.Equal(testingInstance, StoreOutcomeStored, outcome)
}

func TestStoreInstanceTreatsConflictAlreadyExistsAsSuccess(testingInstance *testing.T) {
	storeClient, _ := newTestClient(testingInstance, func(responseWriter http.ResponseWriter, httpRequest *http.Request) {
		responseWriter.WriteHeader(http.StatusConflict)
		responseWriter.Write([]byte(`{"00081190":{"00081197":{"Value":["Duplicate instance: already exists"]}}}`))
	})

	outcome, storeError := storeClient.StoreInstance(context.Background(), bytes.NewReader([]byte("DICM-PAYLOAD")), 12)

	require.NoError(testingInstance, storeError)
	assert.Equal(testingInstance, StoreOutcomeAlreadyExists, outcome)
}

func TestStoreInstanceTreatsConflictWithFailedSOPSequenceAlreadyExistsAsSuccess(testingInstance *testing.T) {
	storeClient, _ := newTestClient(testingInstance, func(responseWriter http.ResponseWriter, httpRequest *http.Request) {
		responseWriter.WriteHeader(http.StatusConflict)
		responseWriter.Write([]byte(`{"00081190":[{"00081155":{"vr":"UI","Value":["1.2.826.0.1.3680043.8.498.100"]},"00081197":{"vr":"SQ","Value":[{"00080104":{"vr":"SH","Value":["Instance already exists"]}}]}}]}`))
	})

	outcome, storeError := storeClient.StoreInstance(context.Background(), bytes.NewReader([]byte("DICM-PAYLOAD")), 12)

	require.NoError(testingInstance, storeError)
	assert.Equal(testingInstance, StoreOutcomeAlreadyExists, outcome)
}

func TestStoreInstanceFailsOnConflictThatIsNotAlreadyExists(testingInstance *testing.T) {
	storeClient, _ := newTestClient(testingInstance, func(responseWriter http.ResponseWriter, httpRequest *http.Request) {
		responseWriter.WriteHeader(http.StatusConflict)
		responseWriter.Write([]byte(`{"00081190":[{"00081197":{"vr":"SQ","Value":[{"00080104":{"vr":"SH","Value":["SOP class not supported"]}}]}}]}`))
	})

	_, storeError := storeClient.StoreInstance(context.Background(), bytes.NewReader([]byte("DICM-PAYLOAD")), 12)

	assert.ErrorIs(testingInstance, storeError, ErrStoreRejectedBytes)
}

func TestStoreInstanceFailsOnFailedSOPSequenceOnSuccessStatus(testingInstance *testing.T) {
	storeClient, _ := newTestClient(testingInstance, func(responseWriter http.ResponseWriter, httpRequest *http.Request) {
		responseWriter.WriteHeader(http.StatusOK)
		responseWriter.Write([]byte(`{"00081190":[{"00081197":{"vr":"SQ","Value":[{"00080104":{"vr":"SH","Value":["Data set does not match SOP Class"]}}]}}]}`))
	})

	_, storeError := storeClient.StoreInstance(context.Background(), bytes.NewReader([]byte("DICM-PAYLOAD")), 12)

	assert.ErrorIs(testingInstance, storeError, ErrStoreRejectedBytes)
}

func TestStoreInstanceStopsReadingAtDeclaredSize(testingInstance *testing.T) {
	storeClient, _ := newTestClient(testingInstance, func(responseWriter http.ResponseWriter, httpRequest *http.Request) {
		bodyBytes, bodyError := io.ReadAll(httpRequest.Body)
		require.NoError(testingInstance, bodyError)
		assert.Equal(testingInstance, "DICOM", string(bodyBytes))
		responseWriter.WriteHeader(http.StatusOK)
	})

	oversizedReader := bytes.NewReader([]byte("DICOM-EXTRA-BYTES-IGNORED"))
	outcome, storeError := storeClient.StoreInstance(context.Background(), oversizedReader, 5)

	require.NoError(testingInstance, storeError)
	assert.Equal(testingInstance, StoreOutcomeStored, outcome)
}

func TestRetrieveInstanceReturnsRawPart10Bytes(testingInstance *testing.T) {
	storeClient, _ := newTestClient(testingInstance, func(responseWriter http.ResponseWriter, httpRequest *http.Request) {
		assert.Equal(testingInstance, "/dicomWeb/studies/1.2.3/series/1.2.4/instances/1.2.5", httpRequest.URL.Path)
		responseWriter.Header().Set("Content-Type", dicomContentType)
		responseWriter.WriteHeader(http.StatusOK)
		responseWriter.Write([]byte("DICM-RAW"))
	})

	instanceBytes, contentType, retrieveError := storeClient.RetrieveInstance(context.Background(), "1.2.3", "1.2.4", "1.2.5")

	require.NoError(testingInstance, retrieveError)
	assert.Equal(testingInstance, "DICM-RAW", string(instanceBytes))
	assert.Equal(testingInstance, dicomContentType, contentType)
}

func TestRetrieveInstanceReportsMissingInstance(testingInstance *testing.T) {
	storeClient, _ := newTestClient(testingInstance, func(responseWriter http.ResponseWriter, httpRequest *http.Request) {
		responseWriter.WriteHeader(http.StatusNotFound)
	})

	_, _, retrieveError := storeClient.RetrieveInstance(context.Background(), "1.2.3", "1.2.4", "1.2.5")

	assert.ErrorIs(testingInstance, retrieveError, ErrInstanceNotFound)
}

func TestRetrieveInstanceFrameDefaultsToOctetStream(testingInstance *testing.T) {
	storeClient, _ := newTestClient(testingInstance, func(responseWriter http.ResponseWriter, httpRequest *http.Request) {
		assert.Equal(testingInstance, "/dicomWeb/studies/1.2.3/series/1.2.4/instances/1.2.5/frames/7", httpRequest.URL.Path)
		assert.Equal(testingInstance, octetStreamContentType, httpRequest.Header.Get("Accept"))
		responseWriter.Header().Set("Content-Type", octetStreamContentType)
		responseWriter.WriteHeader(http.StatusOK)
		responseWriter.Write([]byte("RAW-FRAME-BYTES"))
	})

	frameBytes, contentType, frameError := storeClient.RetrieveInstanceFrame(context.Background(), "1.2.3", "1.2.4", "1.2.5", 7, "")

	require.NoError(testingInstance, frameError)
	assert.Equal(testingInstance, "RAW-FRAME-BYTES", string(frameBytes))
	assert.Equal(testingInstance, octetStreamContentType, contentType)
}

func TestRetrieveInstanceFrameForwardsRequestedAcceptHeader(testingInstance *testing.T) {
	storeClient, _ := newTestClient(testingInstance, func(responseWriter http.ResponseWriter, httpRequest *http.Request) {
		assert.Equal(testingInstance, "image/jpeg", httpRequest.Header.Get("Accept"))
		responseWriter.Header().Set("Content-Type", "image/jpeg")
		responseWriter.WriteHeader(http.StatusOK)
		responseWriter.Write([]byte{0xFF, 0xD8, 0xFF, 0xD9})
	})

	frameBytes, contentType, frameError := storeClient.RetrieveInstanceFrame(context.Background(), "1.2.3", "1.2.4", "1.2.5", 1, "image/jpeg")

	require.NoError(testingInstance, frameError)
	assert.Equal(testingInstance, []byte{0xFF, 0xD8, 0xFF, 0xD9}, frameBytes)
	assert.Equal(testingInstance, "image/jpeg", contentType)
}

func TestRetrieveInstanceMetadataParsesDICOMJSON(testingInstance *testing.T) {
	storeClient, _ := newTestClient(testingInstance, func(responseWriter http.ResponseWriter, httpRequest *http.Request) {
		assert.Equal(testingInstance, "/dicomWeb/studies/1.2.3/series/1.2.4/instances/1.2.5/metadata", httpRequest.URL.Path)
		assert.Equal(testingInstance, dicomJSONContentType, httpRequest.Header.Get("Accept"))
		responseWriter.WriteHeader(http.StatusOK)
		responseWriter.Write([]byte(`[{
			"00080016": {"vr":"UI","Value":["1.2.840.10008.5.1.4.1.1.2"]},
			"00080018": {"vr":"UI","Value":["1.2.826.0.1.3680043.8.498.100"]},
			"00080060": {"vr":"CS","Value":["CT"]},
			"0008103E": {"vr":"LO","Value":["AXIAL 1.0MM"]},
			"0020000D": {"vr":"UI","Value":["1.2.3"]},
			"0020000E": {"vr":"UI","Value":["1.2.4"]},
			"00200011": {"vr":"IS","Value":["2"]},
			"00200013": {"vr":"IS","Value":["7"]},
			"00280008": {"vr":"IS","Value":["1"]},
			"00280010": {"vr":"US","Value":[512]},
			"00280011": {"vr":"US","Value":[512]},
			"00280030": {"vr":"DS","Value":["0.703125","0.703125"]}
		}]`))
	})

	references, metadataError := storeClient.RetrieveInstanceMetadata(context.Background(), "1.2.3", "1.2.4", "1.2.5")

	require.NoError(testingInstance, metadataError)
	require.Len(testingInstance, references, 1)
	assert.Equal(testingInstance, "1.2.3", references[0].StudyInstanceUID)
	assert.Equal(testingInstance, "1.2.4", references[0].SeriesInstanceUID)
	assert.Equal(testingInstance, "1.2.826.0.1.3680043.8.498.100", references[0].SOPInstanceUID)
	assert.Equal(testingInstance, "1.2.840.10008.5.1.4.1.1.2", references[0].SOPClassUID)
	assert.Equal(testingInstance, "CT", references[0].Modality)
	assert.Equal(testingInstance, "AXIAL 1.0MM", references[0].SeriesDescription)
	assert.Equal(testingInstance, 2, references[0].SeriesNumber)
	assert.Equal(testingInstance, 7, references[0].InstanceNumber)
	assert.Equal(testingInstance, 512, references[0].Rows)
	assert.Equal(testingInstance, 512, references[0].Columns)
	assert.Equal(testingInstance, 1, references[0].NumberOfFrames)
	assert.Equal(testingInstance, "0.703125", references[0].PixelSpacingRow)
	assert.Equal(testingInstance, "0.703125", references[0].PixelSpacingColumn)
}

func TestDeleteStudyReturnsNoContent(testingInstance *testing.T) {
	storeClient, _ := newTestClient(testingInstance, func(responseWriter http.ResponseWriter, httpRequest *http.Request) {
		assert.Equal(testingInstance, http.MethodDelete, httpRequest.Method)
		assert.Equal(testingInstance, "/dicomWeb/studies/1.2.3", httpRequest.URL.Path)
		responseWriter.WriteHeader(http.StatusNoContent)
	})

	assert.NoError(testingInstance, storeClient.DeleteStudy(context.Background(), "1.2.3"))
}

func TestDeleteStudyReportsMissingStudy(testingInstance *testing.T) {
	storeClient, _ := newTestClient(testingInstance, func(responseWriter http.ResponseWriter, httpRequest *http.Request) {
		responseWriter.WriteHeader(http.StatusNotFound)
	})

	assert.ErrorIs(testingInstance, storeClient.DeleteStudy(context.Background(), "1.2.3"), ErrStudyNotFound)
}

func TestRetryExhaustionSurfacesServerError(testingInstance *testing.T) {
	attemptCount := 0
	storeClient, _ := newTestClient(testingInstance, func(responseWriter http.ResponseWriter, httpRequest *http.Request) {
		attemptCount++
		responseWriter.WriteHeader(http.StatusBadRequest)
		responseWriter.Write([]byte("malformed request"))
	})

	_, storeError := storeClient.StoreInstance(context.Background(), bytes.NewReader([]byte("DICM")), 4)

	assert.ErrorIs(testingInstance, storeError, ErrStoreRejectedBytes)
	assert.Equal(testingInstance, 1, attemptCount)
}

func TestNewClientWithBaseURLRejectsRelativeURL(testingInstance *testing.T) {
	_, clientError := NewClientWithBaseURL("/dicomWeb")

	assert.Error(testingInstance, clientError)
}
