package fhir

import "time"

func NewConditionResource(patientFHIRID, encounterFHIRID, icdCode, display, clinicalStatus string, onsetAt time.Time) map[string]interface{} {
	return map[string]interface{}{
		"resourceType": "Condition",
		"clinicalStatus": map[string]interface{}{
			"coding": []map[string]interface{}{
				{"system": "http://terminology.hl7.org/CodeSystem/condition-clinical", "code": clinicalStatus},
			},
		},
		"code": map[string]interface{}{
			"coding": []map[string]interface{}{
				{"system": "http://hl7.org/fhir/sid/icd-10", "code": icdCode, "display": display},
			},
			"text": display,
		},
		"subject":         map[string]interface{}{"reference": "Patient/" + patientFHIRID},
		"encounter":       map[string]interface{}{"reference": "Encounter/" + encounterFHIRID},
		"onsetDateTime":   onsetAt.Format(time.RFC3339),
		"recordedDate":    time.Now().Format(time.RFC3339),
	}
}
