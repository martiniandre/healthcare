package timeline

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const encounterBundleFixture = `{
  "resourceType": "Bundle",
  "type": "searchset",
  "entry": [
    {
      "resource": {
        "resourceType": "Encounter",
        "id": "enc-9",
        "status": "finished",
        "period": { "start": "2026-08-12T09:00:00Z", "end": "2026-08-12T10:00:00Z" },
        "reasonCode": [{ "text": "Chest pain follow-up" }]
      }
    }
  ]
}`

const observationBundleFixture = `{
  "resourceType": "Bundle",
  "type": "searchset",
  "entry": [
    {
      "resource": {
        "resourceType": "Observation",
        "id": "obs-9",
        "status": "final",
        "code": { "coding": [{ "system": "http://loinc.org", "code": "85354-9", "display": "Blood pressure" }] },
        "effectiveDateTime": "2026-08-02T10:30:00Z",
        "valueQuantity": { "value": 140, "unit": "mmHg" },
        "referenceRange": [{ "low": { "value": 90 }, "high": { "value": 120 } }]
      }
    }
  ]
}`

const conditionBundleFixture = `{
  "resourceType": "Bundle",
  "type": "searchset",
  "entry": [
    {
      "resource": {
        "resourceType": "Condition",
        "id": "cond-9",
        "clinicalStatus": { "coding": [{ "code": "active" }] },
        "code": { "coding": [{ "system": "http://hl7.org/fhir/sid/icd-10", "code": "I10", "display": "Hypertension" }] },
        "onsetDateTime": "2026-06-01T00:00:00Z",
        "recordedDate": "2026-07-15T08:00:00Z"
      }
    }
  ]
}`

const allergyBundleFixture = `{
  "resourceType": "Bundle",
  "type": "searchset",
  "entry": [
    {
      "resource": {
        "resourceType": "AllergyIntolerance",
        "id": "alg-9",
        "clinicalStatus": { "coding": [{ "code": "active" }] },
        "criticality": "high",
        "code": { "coding": [{ "code": "PEN", "display": "Penicillin" }], "text": "Penicillin" },
        "recordedDate": "2026-07-01T12:00:00Z",
        "reaction": [{ "manifestation": [{ "text": "Skin rash" }], "severity": "moderate" }]
      }
    }
  ]
}`

func TestParseEncounterTimelineBundle(t *testing.T) {
	entries, parseErr := parseEncounterTimelineBundle(json.RawMessage(encounterBundleFixture))

	require.NoError(t, parseErr)
	require.Len(t, entries, 1)

	entry := entries[0]
	assert.Equal(t, "Encounter", entry.ResourceType)
	assert.Equal(t, "enc-9", entry.FHIRResourceID)
	assert.Equal(t, "finished", entry.Status)
	assert.Equal(t, "Chest pain follow-up", entry.Title)
	assert.True(t, entry.RecordedAt.Equal(time.Date(2026, 8, 12, 9, 0, 0, 0, time.UTC)))
	require.NotNil(t, entry.PeriodEnd)
	assert.True(t, entry.PeriodEnd.Equal(time.Date(2026, 8, 12, 10, 0, 0, 0, time.UTC)))
}

func TestParseObservationTimelineBundle(t *testing.T) {
	entries, parseErr := parseObservationTimelineBundle(json.RawMessage(observationBundleFixture))

	require.NoError(t, parseErr)
	require.Len(t, entries, 1)

	entry := entries[0]
	assert.Equal(t, "85354-9", entry.Code)
	assert.Equal(t, "Blood pressure", entry.Title)
	require.NotNil(t, entry.ValueQuantity)
	assert.InDelta(t, 140, *entry.ValueQuantity, 0.001)
	assert.Equal(t, "mmHg", entry.ValueUnit)
	require.NotNil(t, entry.ReferenceLow)
	assert.InDelta(t, 90, *entry.ReferenceLow, 0.001)
	require.NotNil(t, entry.ReferenceHigh)
	assert.InDelta(t, 120, *entry.ReferenceHigh, 0.001)
}

func TestParseConditionTimelineBundle(t *testing.T) {
	entries, parseErr := parseConditionTimelineBundle(json.RawMessage(conditionBundleFixture))

	require.NoError(t, parseErr)
	require.Len(t, entries, 1)

	entry := entries[0]
	assert.Equal(t, "Hypertension", entry.Title)
	assert.Equal(t, "I10", entry.Code)
	assert.Equal(t, "active", entry.Status)
	assert.Equal(t, "2026-06-01T00:00:00Z", entry.OnsetDate)
	assert.True(t, entry.RecordedAt.Equal(time.Date(2026, 7, 15, 8, 0, 0, 0, time.UTC)))
}

func TestParseAllergyTimelineBundle(t *testing.T) {
	entries, parseErr := parseAllergyTimelineBundle(json.RawMessage(allergyBundleFixture))

	require.NoError(t, parseErr)
	require.Len(t, entries, 1)

	entry := entries[0]
	assert.Equal(t, "Penicillin", entry.Title)
	assert.Equal(t, "high", entry.Criticality)
	assert.Equal(t, "Skin rash", entry.Reaction)
	assert.True(t, entry.RecordedAt.Equal(time.Date(2026, 7, 1, 12, 0, 0, 0, time.UTC)))
}
