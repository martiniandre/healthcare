package tests

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/healthcare/backend/internal/modules/imaging"
	"github.com/healthcare/backend/internal/shared/apperrors"
	"github.com/healthcare/backend/internal/shared/healthcare"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type stubFHIRClient struct {
	healthcare.FHIRClient
	patientResource json.RawMessage
	getError        error
	requestedType   string
	requestedID     string
}

func (stub *stubFHIRClient) GetResource(ctx context.Context, resourceType, resourceID string) (json.RawMessage, error) {
	stub.requestedType = resourceType
	stub.requestedID = resourceID
	if stub.getError != nil {
		return nil, stub.getError
	}
	return stub.patientResource, nil
}

func TestFHIRPatientIdentityReader_ReadsNamesAndDocuments(testingInstance *testing.T) {
	stub := &stubFHIRClient{patientResource: json.RawMessage(`{
		"resourceType": "Patient",
		"id": "patient-123",
		"name": [{"given": ["Maria", "Clara"], "family": "Silva"}],
		"identifier": [{"system": "urn:oid:2.16.840.1.113883.13.237", "value": "529.982.247-25"}, {"system": "other", "value": " "}]
	}`)}
	reader := imaging.NewFHIRPatientIdentityReader(stub)

	identity, readError := reader.ReadPatientIdentity(context.Background(), "patient-123")

	require.NoError(testingInstance, readError)
	assert.Equal(testingInstance, "Patient", stub.requestedType)
	assert.Equal(testingInstance, "patient-123", stub.requestedID)
	assert.Equal(testingInstance, []string{"Maria Clara Silva"}, identity.FullNames)
	assert.Equal(testingInstance, []string{"529.982.247-25"}, identity.DocumentIDs)
}

func TestFHIRPatientIdentityReader_MapsNotFoundToPatientNotFound(testingInstance *testing.T) {
	stub := &stubFHIRClient{getError: &healthcare.NotFoundError{ResourceType: "Patient", ResourceID: "missing"}}
	reader := imaging.NewFHIRPatientIdentityReader(stub)

	_, readError := reader.ReadPatientIdentity(context.Background(), "missing")

	assert.ErrorIs(testingInstance, readError, apperrors.ErrPatientNotFound)
}

func TestFHIRPatientIdentityReader_PropagatesOtherFailures(testingInstance *testing.T) {
	stub := &stubFHIRClient{getError: errors.New("fhir unavailable")}
	reader := imaging.NewFHIRPatientIdentityReader(stub)

	_, readError := reader.ReadPatientIdentity(context.Background(), "patient-123")

	require.Error(testingInstance, readError)
	assert.NotErrorIs(testingInstance, readError, apperrors.ErrPatientNotFound)
}

func TestFHIRPatientIdentityReader_RejectsMalformedPatientPayload(testingInstance *testing.T) {
	stub := &stubFHIRClient{patientResource: json.RawMessage(`not-json`)}
	reader := imaging.NewFHIRPatientIdentityReader(stub)

	_, readError := reader.ReadPatientIdentity(context.Background(), "patient-123")

	require.Error(testingInstance, readError)
}
