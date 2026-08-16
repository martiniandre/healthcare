package encounter

import (
	"testing"

	"github.com/healthcare/backend/internal/shared/fhir"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMapFHIREncounterToDomain_ReasonMapping(t *testing.T) {
	testCases := []struct {
		name                string
		fhirEncounter       *fhir.Encounter
		expectedReasonCode  string
		expectedReasonText  string
	}{
		{
			name: "code display and text maps code and display",
			fhirEncounter: &fhir.Encounter{
				ReasonCode: []fhir.CodeableConcept{
					{
						Coding: []fhir.Coding{
							{System: "http://hl7.org/fhir/sid/icd-10", Code: "R10.9", Display: "Abdominal pain"},
						},
						Text: "Unspecified abdominal pain",
					},
				},
			},
			expectedReasonCode: "R10.9",
			expectedReasonText: "Abdominal pain",
		},
		{
			name: "display only persists as text",
			fhirEncounter: &fhir.Encounter{
				ReasonCode: []fhir.CodeableConcept{
					{
						Coding: []fhir.Coding{},
						Text:   "Dor abdominal persistente",
					},
				},
			},
			expectedReasonCode: "",
			expectedReasonText: "Dor abdominal persistente",
		},
		{
			name: "code and display without text",
			fhirEncounter: &fhir.Encounter{
				ReasonCode: []fhir.CodeableConcept{
					{
						Coding: []fhir.Coding{
							{System: "http://hl7.org/fhir/sid/icd-10", Code: "A00", Display: "Cholera"},
						},
					},
				},
			},
			expectedReasonCode: "A00",
			expectedReasonText: "Cholera",
		},
		{
			name:          "no reason code leaves fields empty",
			fhirEncounter: &fhir.Encounter{},
			expectedReasonCode: "",
			expectedReasonText: "",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(subTest *testing.T) {
			mappedEncounter := mapFHIREncounterToDomain(testCase.fhirEncounter)

			require.NotNil(subTest, mappedEncounter)
			assert.Equal(subTest, testCase.expectedReasonCode, mappedEncounter.ReasonCode)
			assert.Equal(subTest, testCase.expectedReasonText, mappedEncounter.ReasonDisplay)
		})
	}
}
