package fhir

import "time"

func NewAllergyIntoleranceResource(patientFHIRID, allergenCode, allergenDisplay, clinicalStatus, reactionDescription string) map[string]interface{} {
	return map[string]interface{}{
		"resourceType": "AllergyIntolerance",
		"clinicalStatus": map[string]interface{}{
			"coding": []map[string]interface{}{
				{"system": "http://terminology.hl7.org/CodeSystem/allergyintolerance-clinical", "code": clinicalStatus},
			},
		},
		"code": map[string]interface{}{
			"coding": []map[string]interface{}{
				{"system": "http://www.nlm.nih.gov/research/umls/rxnorm", "code": allergenCode, "display": allergenDisplay},
			},
			"text": allergenDisplay,
		},
		"patient":      map[string]interface{}{"reference": "Patient/" + patientFHIRID},
		"recordedDate": time.Now().Format(time.RFC3339),
		"reaction": []map[string]interface{}{
			{
				"manifestation": []map[string]interface{}{
					{"text": reactionDescription},
				},
			},
		},
	}
}
