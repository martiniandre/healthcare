package portal

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseEncounterPortalBundle_ReasonDisplayMapping(t *testing.T) {
	testCases := []struct {
		name                string
		bundleJSON          string
		expectedReasonTexts []string
	}{
		{
			name: "reads coding display",
			bundleJSON: `{
				"resourceType": "Bundle",
				"entry": [
					{
						"resource": {
							"resourceType": "Encounter",
							"id": "enc-1",
							"status": "finished",
							"subject": {"reference": "Patient/pat-1"},
							"period": {"start": "2026-08-01T09:00:00Z"},
							"reasonCode": [
								{"coding": [{"system": "http://hl7.org/fhir/sid/icd-10", "code": "R10.9", "display": "Abdominal pain"}]}
							]
						}
					}
				]
			}`,
			expectedReasonTexts: []string{"Abdominal pain"},
		},
		{
			name: "falls back to text when display only",
			bundleJSON: `{
				"resourceType": "Bundle",
				"entry": [
					{
						"resource": {
							"resourceType": "Encounter",
							"id": "enc-2",
							"status": "finished",
							"subject": {"reference": "Patient/pat-1"},
							"period": {"start": "2026-08-02T09:00:00Z"},
							"reasonCode": [
								{"text": "Dor abdominal persistente"}
							]
						}
					}
				]
			}`,
			expectedReasonTexts: []string{"Dor abdominal persistente"},
		},
		{
			name: "keeps empty when no reason code",
			bundleJSON: `{
				"resourceType": "Bundle",
				"entry": [
					{
						"resource": {
							"resourceType": "Encounter",
							"id": "enc-3",
							"status": "finished",
							"subject": {"reference": "Patient/pat-1"},
							"period": {"start": "2026-08-03T09:00:00Z"}
						}
					}
				]
			}`,
			expectedReasonTexts: []string{""},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(subTest *testing.T) {
			parsedEncounters, err := parseEncounterPortalBundle(json.RawMessage(testCase.bundleJSON))

			require.NoError(subTest, err)
			require.Len(subTest, parsedEncounters, len(testCase.expectedReasonTexts))
			for index, expectedReasonText := range testCase.expectedReasonTexts {
				assert.Equal(subTest, expectedReasonText, parsedEncounters[index].ReasonDisplay)
			}
		})
	}
}
