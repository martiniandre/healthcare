package fhir

import "time"

func NewMedicationRequestResource(patientFHIRID, encounterFHIRID, practitionerFHIRID, medicationCode, medicationName, dosageInstructions string) map[string]interface{} {
	return map[string]interface{}{
		"resourceType": "MedicationRequest",
		"status":       "active",
		"intent":       "order",
		"medicationCodeableConcept": map[string]interface{}{
			"coding": []map[string]interface{}{
				{"system": "http://www.nlm.nih.gov/research/umls/rxnorm", "code": medicationCode, "display": medicationName},
			},
			"text": medicationName,
		},
		"subject":    map[string]interface{}{"reference": "Patient/" + patientFHIRID},
		"encounter":  map[string]interface{}{"reference": "Encounter/" + encounterFHIRID},
		"requester":  map[string]interface{}{"reference": "Practitioner/" + practitionerFHIRID},
		"authoredOn": time.Now().Format(time.RFC3339),
		"dosageInstruction": []map[string]interface{}{
			{"text": dosageInstructions},
		},
	}
}
