package fhir

import "time"

type DiagnosticReportResource struct {
	ResourceType string         `json:"resourceType"`
	ID           string         `json:"id,omitempty"`
	Status       string         `json:"status"`
	Code         CodeableConcept `json:"code"`
	Subject      Reference       `json:"subject"`
	Encounter    Reference       `json:"encounter"`
	Issued       string         `json:"issued"`
	Conclusion   string         `json:"conclusion,omitempty"`
}

func NewDiagnosticReportResource(patientFHIRID, encounterFHIRID, reportCode, reportDisplay, conclusion string) *DiagnosticReportResource {
	return &DiagnosticReportResource{
		ResourceType: "DiagnosticReport",
		Status:       "final",
		Code: CodeableConcept{
			Coding: []Coding{
				{System: "http://loinc.org", Code: reportCode, Display: reportDisplay},
			},
			Text: reportDisplay,
		},
		Subject:    Reference{Reference: "Patient/" + patientFHIRID},
		Encounter:  Reference{Reference: "Encounter/" + encounterFHIRID},
		Issued:     time.Now().Format(time.RFC3339),
		Conclusion: conclusion,
	}
}
