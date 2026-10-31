package tests

import (
	"context"
	"encoding/json"
	"errors"
	"net/url"
	"testing"

	"github.com/healthcare/backend/internal/modules/medication"
	"github.com/healthcare/backend/internal/shared/healthcare"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type searchRecordingFHIRClient struct {
	healthcare.FHIRClient
	responseBody  json.RawMessage
	searchErr     error
	recordedType  string
	recordedQuery string
}

func (client *searchRecordingFHIRClient) SearchResources(ctx context.Context, resourceType string, queryParams string) (json.RawMessage, error) {
	client.recordedType = resourceType
	client.recordedQuery = queryParams
	return client.responseBody, client.searchErr
}

func TestGetActiveMedicationRequestsByPatient_QueriesPatientAndKeepsOnlyActive(testingInstance *testing.T) {
	bundle := `{"resourceType":"Bundle","entry":[
		{"resource":{"resourceType":"MedicationRequest","id":"request-1","status":"active","subject":{"reference":"Patient/patient-1"},"medicationCodeableConcept":{"text":"Warfarin"}}},
		{"resource":{"resourceType":"MedicationRequest","id":"request-2","status":"stopped","subject":{"reference":"Patient/patient-1"},"medicationCodeableConcept":{"text":"Ibuprofen"}}}
	]}`
	fhirClient := &searchRecordingFHIRClient{responseBody: json.RawMessage(bundle)}
	repository := medication.NewRepository(fhirClient)

	medications, err := repository.GetActiveMedicationRequestsByPatient(context.Background(), "patient-1")

	require.NoError(testingInstance, err)
	require.Len(testingInstance, medications, 1)
	assert.Equal(testingInstance, "request-1", medications[0].FHIRResourceID)
	assert.Equal(testingInstance, "Warfarin", medications[0].MedicationName)
	assert.Equal(testingInstance, "MedicationRequest", fhirClient.recordedType)
	parsedQuery, parseErr := url.ParseQuery(fhirClient.recordedQuery)
	require.NoError(testingInstance, parseErr)
	assert.Equal(testingInstance, "Patient/patient-1", parsedQuery.Get("subject"))
	assert.Equal(testingInstance, "active", parsedQuery.Get("status"))
}

func TestGetActiveMedicationRequestsByPatient_SearchFailureIsReturned(testingInstance *testing.T) {
	searchFailure := errors.New("search failed")
	repository := medication.NewRepository(&searchRecordingFHIRClient{searchErr: searchFailure})

	medications, err := repository.GetActiveMedicationRequestsByPatient(context.Background(), "patient-1")

	assert.Nil(testingInstance, medications)
	assert.ErrorIs(testingInstance, err, searchFailure)
}
