package dicomstore

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"golang.org/x/oauth2/google"
)

const (
	maxErrorBodyBytes          int64 = 64 << 10
	maxErrorBodyInMessage            = 1024
	defaultHTTPTimeout               = 120 * time.Second
	maxRetryAttempts                 = 3
	dicomContentType                 = "application/dicom"
	dicomJSONContentType             = "application/dicom+json"
	octetStreamContentType           = "application/octet-stream"
	alreadyExistsFailureReason       = "already exists"
	codeMeaningTag                   = "00080104"
	codeValueTag                     = "00080100"
)

var (
	ErrStudyNotFound      = errors.New("dicom study not found in store")
	ErrInstanceNotFound   = errors.New("dicom instance not found in store")
	ErrStoreRejectedBytes = errors.New("dicom store rejected instance")
)

type StoreOutcome string

const (
	StoreOutcomeStored        StoreOutcome = "STORED"
	StoreOutcomeAlreadyExists StoreOutcome = "ALREADY_EXISTS"
)

type InstanceReference struct {
	StudyInstanceUID   string
	SeriesInstanceUID  string
	SOPInstanceUID     string
	SOPClassUID        string
	TransferSyntaxUID  string
	Modality           string
	SeriesDescription  string
	SeriesNumber       int
	InstanceNumber     int
	Rows               int
	Columns            int
	NumberOfFrames     int
	PixelSpacingRow    string
	PixelSpacingColumn string
}

type Client struct {
	httpClient   *http.Client
	baseURL      string
	dicomStoreID string
}

type dicomJSONValue struct {
	ValueRepresentation string          `json:"vr"`
	Values              json.RawMessage `json:"Value"`
}

type failedSOPSequence struct {
	ReferencedSOPInstanceUID string         `json:"00081155"`
	FailureReason            dicomJSONValue `json:"00081197"`
}

type dicomJSONError struct {
	FailedSOPSequence []failedSOPSequence `json:"00081190"`
}

func NewClient(ctx context.Context, projectID, locationID, datasetID, dicomStoreID string) (*Client, error) {
	credentialsScopes := []string{"https://www.googleapis.com/auth/cloud-healthcare"}

	googleHTTPClient, err := google.DefaultClient(ctx, credentialsScopes...)
	if err != nil {
		return nil, fmt.Errorf("failed to create google default client: %w", err)
	}
	googleHTTPClient.Timeout = defaultHTTPTimeout

	baseURL := fmt.Sprintf(
		"https://healthcare.googleapis.com/v1/projects/%s/locations/%s/datasets/%s/dicomStores/%s/dicomWeb",
		projectID, locationID, datasetID, dicomStoreID,
	)

	return &Client{
		httpClient:   googleHTTPClient,
		baseURL:      baseURL,
		dicomStoreID: dicomStoreID,
	}, nil
}

func NewClientWithBaseURL(baseURL string) (*Client, error) {
	parsedBaseURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse DICOMweb base url: %w", err)
	}
	if parsedBaseURL.Scheme == "" || parsedBaseURL.Host == "" {
		return nil, fmt.Errorf("invalid DICOMweb base url %q: expected an absolute url like http://localhost:8080/dicomWeb", baseURL)
	}

	return &Client{
		httpClient: &http.Client{Timeout: defaultHTTPTimeout},
		baseURL:    strings.TrimRight(baseURL, "/"),
	}, nil
}

func (storeClient *Client) StoreInstance(ctx context.Context, instanceReader io.Reader, instanceSize int64) (StoreOutcome, error) {
	endpoint := fmt.Sprintf("%s/studies", storeClient.baseURL)
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, newBoundedReader(instanceReader, instanceSize))
	if err != nil {
		return "", fmt.Errorf("failed to create store request: %w", err)
	}
	request.Header.Set("Content-Type", dicomContentType)
	request.Header.Set("Accept", dicomJSONContentType)
	if instanceSize > 0 {
		request.ContentLength = instanceSize
	}

	response, doError := storeClient.doWithRetry(request)
	if doError != nil {
		return "", doError
	}

	if response.StatusCode == http.StatusConflict && isAlreadyExistsFailure(response.Body) {
		return StoreOutcomeAlreadyExists, nil
	}

	if response.StatusCode == http.StatusConflict {
		return "", fmt.Errorf("%w: %s", ErrStoreRejectedBytes, truncateErrorBody(response.Body))
	}

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return "", fmt.Errorf("%w: unexpected status %d - %s", ErrStoreRejectedBytes, response.StatusCode, truncateErrorBody(response.Body))
	}

	if failedSequences := parseFailedSOPSequences(response.Body); len(failedSequences) > 0 {
		for _, failedSequence := range failedSequences {
			if strings.Contains(strings.ToLower(failureReasonText(failedSequence)), alreadyExistsFailureReason) {
				return StoreOutcomeAlreadyExists, nil
			}
		}
		return "", fmt.Errorf("%w: %s", ErrStoreRejectedBytes, truncateErrorBody(response.Body))
	}

	return StoreOutcomeStored, nil
}

func (storeClient *Client) RetrieveInstance(ctx context.Context, studyUID, seriesUID, instanceUID string) ([]byte, string, error) {
	endpoint := fmt.Sprintf("%s/studies/%s/series/%s/instances/%s",
		storeClient.baseURL,
		url.PathEscape(studyUID),
		url.PathEscape(seriesUID),
		url.PathEscape(instanceUID),
	)

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create retrieve request: %w", err)
	}
	request.Header.Set("Accept", dicomContentType)

	response, doError := storeClient.doWithRetry(request)
	if doError != nil {
		return nil, "", doError
	}
	if response.StatusCode == http.StatusNotFound {
		return nil, "", ErrInstanceNotFound
	}
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return nil, "", fmt.Errorf("dicom store: unexpected status %d - %s", response.StatusCode, truncateErrorBody(response.Body))
	}
	return response.Body, response.ContentType, nil
}

func (storeClient *Client) RetrieveInstanceFrame(ctx context.Context, studyUID, seriesUID, instanceUID string, frameNumber int, acceptHeader string) ([]byte, string, error) {
	endpoint := fmt.Sprintf("%s/studies/%s/series/%s/instances/%s/frames/%d",
		storeClient.baseURL,
		url.PathEscape(studyUID),
		url.PathEscape(seriesUID),
		url.PathEscape(instanceUID),
		frameNumber,
	)

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create frame request: %w", err)
	}

	if strings.TrimSpace(acceptHeader) == "" {
		acceptHeader = octetStreamContentType
	}
	request.Header.Set("Accept", acceptHeader)

	response, doError := storeClient.doWithRetry(request)
	if doError != nil {
		return nil, "", doError
	}
	if response.StatusCode == http.StatusNotFound {
		return nil, "", ErrInstanceNotFound
	}
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return nil, "", fmt.Errorf("dicom store: unexpected status %d - %s", response.StatusCode, truncateErrorBody(response.Body))
	}
	return response.Body, response.ContentType, nil
}

func (storeClient *Client) RetrieveInstanceMetadata(ctx context.Context, studyUID, seriesUID, instanceUID string) ([]InstanceReference, error) {
	endpoint := fmt.Sprintf("%s/studies/%s/series/%s/instances/%s/metadata",
		storeClient.baseURL,
		url.PathEscape(studyUID),
		url.PathEscape(seriesUID),
		url.PathEscape(instanceUID),
	)

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create metadata request: %w", err)
	}
	request.Header.Set("Accept", dicomJSONContentType)

	response, doError := storeClient.doWithRetry(request)
	if doError != nil {
		return nil, doError
	}
	if response.StatusCode == http.StatusNotFound {
		return nil, ErrInstanceNotFound
	}
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return nil, fmt.Errorf("dicom store: unexpected status %d - %s", response.StatusCode, truncateErrorBody(response.Body))
	}

	var metadataArray []map[string]dicomJSONValue
	if unmarshalError := json.Unmarshal(response.Body, &metadataArray); unmarshalError != nil {
		return nil, fmt.Errorf("failed to parse instance metadata: %w", unmarshalError)
	}
	references := make([]InstanceReference, 0, len(metadataArray))
	for _, metadataEntry := range metadataArray {
		references = append(references, instanceReferenceFromDICOMJSON(metadataEntry))
	}
	return references, nil
}

func (storeClient *Client) DeleteStudy(ctx context.Context, studyUID string) error {
	endpoint := fmt.Sprintf("%s/studies/%s", storeClient.baseURL, url.PathEscape(studyUID))

	request, err := http.NewRequestWithContext(ctx, http.MethodDelete, endpoint, nil)
	if err != nil {
		return fmt.Errorf("failed to create delete study request: %w", err)
	}

	response, doError := storeClient.doWithRetry(request)
	if doError != nil {
		return doError
	}
	if response.StatusCode == http.StatusNotFound {
		return ErrStudyNotFound
	}
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return fmt.Errorf("dicom store: unexpected status %d - %s", response.StatusCode, truncateErrorBody(response.Body))
	}
	return nil
}

func instanceReferenceFromDICOMJSON(metadataEntry map[string]dicomJSONValue) InstanceReference {
	spacingValues := firstTwoStrings(dicomJSONValueAt(metadataEntry, "00280030"))

	reference := InstanceReference{
		StudyInstanceUID:  firstString(dicomJSONValueAt(metadataEntry, "0020000D")),
		SeriesInstanceUID: firstString(dicomJSONValueAt(metadataEntry, "0020000E")),
		SOPInstanceUID:    firstString(dicomJSONValueAt(metadataEntry, "00080018")),
		SOPClassUID:       firstString(dicomJSONValueAt(metadataEntry, "00080016")),
		TransferSyntaxUID: firstString(dicomJSONValueAt(metadataEntry, "00020010")),
		Modality:          firstString(dicomJSONValueAt(metadataEntry, "00080060")),
		SeriesDescription: firstString(dicomJSONValueAt(metadataEntry, "0008103E")),
		SeriesNumber:      firstInteger(dicomJSONValueAt(metadataEntry, "00200011")),
		InstanceNumber:    firstInteger(dicomJSONValueAt(metadataEntry, "00200013")),
		Rows:              firstInteger(dicomJSONValueAt(metadataEntry, "00280010")),
		Columns:           firstInteger(dicomJSONValueAt(metadataEntry, "00280011")),
		NumberOfFrames:    firstInteger(dicomJSONValueAt(metadataEntry, "00280008")),
	}
	if len(spacingValues) == 2 {
		reference.PixelSpacingRow = spacingValues[0]
		reference.PixelSpacingColumn = spacingValues[1]
	}
	return reference
}

func dicomJSONValueAt(metadataEntry map[string]dicomJSONValue, tagKey string) json.RawMessage {
	entry, present := metadataEntry[tagKey]
	if !present {
		return nil
	}
	return entry.Values
}

func decodeDICOMJSONValues(rawValues json.RawMessage) []string {
	if len(rawValues) == 0 {
		return nil
	}

	var stringValues []string
	if unmarshalError := json.Unmarshal(rawValues, &stringValues); unmarshalError == nil {
		return stringValues
	}

	var numericValues []float64
	if unmarshalError := json.Unmarshal(rawValues, &numericValues); unmarshalError == nil {
		convertedValues := make([]string, 0, len(numericValues))
		for _, numericValue := range numericValues {
			convertedValues = append(convertedValues, strconv.FormatFloat(numericValue, 'f', -1, 64))
		}
		return convertedValues
	}

	return nil
}

func firstString(rawValues json.RawMessage) string {
	decodedValues := decodeDICOMJSONValues(rawValues)
	if len(decodedValues) == 0 {
		return ""
	}
	return strings.TrimSpace(decodedValues[0])
}

func firstTwoStrings(rawValues json.RawMessage) []string {
	decodedValues := decodeDICOMJSONValues(rawValues)
	if len(decodedValues) < 2 {
		return decodedValues
	}
	return decodedValues[:2]
}

func firstInteger(rawValues json.RawMessage) int {
	decodedValue := firstString(rawValues)
	if decodedValue == "" {
		return 0
	}
	parsedValue, parseError := strconv.ParseFloat(decodedValue, 64)
	if parseError != nil {
		return 0
	}
	return int(parsedValue)
}

func parseFailedSOPSequences(responseBody []byte) []failedSOPSequence {
	var decodedError dicomJSONError
	if unmarshalError := json.Unmarshal(responseBody, &decodedError); unmarshalError != nil {
		return nil
	}
	return decodedError.FailedSOPSequence
}

func failureReasonText(failedSequence failedSOPSequence) string {
	var codeSequenceItems []map[string]dicomJSONValue
	if unmarshalError := json.Unmarshal(failedSequence.FailureReason.Values, &codeSequenceItems); unmarshalError == nil {
		for _, codeSequenceItem := range codeSequenceItems {
			if codeMeaning := firstString(dicomJSONValueAt(codeSequenceItem, codeMeaningTag)); codeMeaning != "" {
				return codeMeaning
			}
			if codeValue := firstString(dicomJSONValueAt(codeSequenceItem, codeValueTag)); codeValue != "" {
				return codeValue
			}
		}
	}
	return firstString(failedSequence.FailureReason.Values)
}

func isAlreadyExistsFailure(responseBody []byte) bool {
	failedSequences := parseFailedSOPSequences(responseBody)
	if len(failedSequences) == 0 {
		return strings.Contains(strings.ToLower(string(responseBody)), alreadyExistsFailureReason)
	}
	for _, failedSequence := range failedSequences {
		if strings.Contains(strings.ToLower(failureReasonText(failedSequence)), alreadyExistsFailureReason) {
			return true
		}
	}
	return false
}

func truncateErrorBody(body []byte) string {
	if len(body) > maxErrorBodyInMessage {
		return string(body[:maxErrorBodyInMessage]) + "..."
	}
	return string(body)
}

type storeResponse struct {
	Body        []byte
	StatusCode  int
	ContentType string
}

func (storeClient *Client) doWithRetry(request *http.Request) (storeResponse, error) {
	retryDelays := []time.Duration{100 * time.Millisecond, 300 * time.Millisecond, 900 * time.Millisecond}

	var lastError error
	for attemptIndex := 0; attemptIndex <= maxRetryAttempts; attemptIndex++ {
		if request.GetBody != nil {
			reusableBody, bodyError := request.GetBody()
			if bodyError == nil {
				request.Body = reusableBody
			}
		}

		response, requestError := storeClient.httpClient.Do(request)
		if requestError != nil {
			lastError = fmt.Errorf("dicom store request failed: %w", requestError)
			if attemptIndex < maxRetryAttempts {
				time.Sleep(retryDelays[attemptIndex])
				continue
			}
			return storeResponse{}, lastError
		}

		var responseReader io.Reader = response.Body
		if response.StatusCode < 200 || response.StatusCode >= 300 {
			responseReader = io.LimitReader(response.Body, maxErrorBodyBytes)
		}

		responseBody, readError := io.ReadAll(responseReader)
		response.Body.Close()

		if readError != nil {
			return storeResponse{}, fmt.Errorf("failed to read dicom store response: %w", readError)
		}

		isRetryableStatus := response.StatusCode == http.StatusTooManyRequests ||
			response.StatusCode == http.StatusBadGateway ||
			response.StatusCode == http.StatusServiceUnavailable ||
			response.StatusCode == http.StatusGatewayTimeout

		if isRetryableStatus && attemptIndex < maxRetryAttempts {
			lastError = fmt.Errorf("dicom store: retryable status %d - %s", response.StatusCode, truncateErrorBody(responseBody))
			time.Sleep(retryDelays[attemptIndex])
			continue
		}

		slog.Debug("dicom store response", "status", response.StatusCode, "url", request.URL.Path)
		return storeResponse{
			Body:        responseBody,
			StatusCode:  response.StatusCode,
			ContentType: response.Header.Get("Content-Type"),
		}, nil
	}

	return storeResponse{}, lastError
}

type boundedReader struct {
	sourceReader  io.Reader
	remainingSize int64
	bounded       bool
}

func newBoundedReader(sourceReader io.Reader, totalSize int64) *boundedReader {
	return &boundedReader{
		sourceReader:  sourceReader,
		remainingSize: totalSize,
		bounded:       totalSize > 0,
	}
}

func (readerInstance *boundedReader) Read(buffer []byte) (int, error) {
	if readerInstance.bounded {
		if readerInstance.remainingSize == 0 {
			return 0, io.EOF
		}
		if int64(len(buffer)) > readerInstance.remainingSize {
			buffer = buffer[:readerInstance.remainingSize]
		}
	}
	bytesRead, readError := readerInstance.sourceReader.Read(buffer)
	readerInstance.remainingSize -= int64(bytesRead)
	return bytesRead, readError
}
