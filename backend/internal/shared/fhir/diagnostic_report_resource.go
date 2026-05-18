package fhir

import "time"

func NewDiagnosticReportResource(patientFHIRID, encounterFHIRID, reportCode, reportDisplay, conclusion string) map[string]interface{} {
	return map[string]interface{}{
		"resourceType": "DiagnosticReport",
		"status":       "final",
		"code": map[string]interface{}{
			"coding": []map[string]interface{}{
				{"system": "http://loinc.org", "code": reportCode, "display": reportDisplay},
			},
			"text": reportDisplay,
		},
		"subject":    map[string]interface{}{"reference": "Patient/" + patientFHIRID},
		"encounter":  map[string]interface{}{"reference": "Encounter/" + encounterFHIRID},
		"issued":     time.Now().Format(time.RFC3339),
		"conclusion": conclusion,
	}
}
