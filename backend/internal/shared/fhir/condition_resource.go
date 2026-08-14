package fhir

import "time"

type ConditionResource struct {
	ResourceType   string          `json:"resourceType"`
	ID             string          `json:"id,omitempty"`
	ClinicalStatus CodeableConcept `json:"clinicalStatus"`
	Code           CodeableConcept `json:"code"`
	Subject        Reference       `json:"subject"`
	Encounter      *Reference      `json:"encounter,omitempty"`
	OnsetDateTime  string          `json:"onsetDateTime"`
	RecordedDate   string          `json:"recordedDate"`
}

func NewConditionResource(patientFHIRID, encounterFHIRID, icdCode, display, clinicalStatus string, onsetAt time.Time) *ConditionResource {
	conditionResource := &ConditionResource{
		ResourceType: "Condition",
		ClinicalStatus: CodeableConcept{
			Coding: []Coding{
				{System: "http://terminology.hl7.org/CodeSystem/condition-clinical", Code: clinicalStatus},
			},
		},
		Code: CodeableConcept{
			Coding: []Coding{
				{System: "http://hl7.org/fhir/sid/icd-10", Code: icdCode, Display: display},
			},
			Text: display,
		},
		Subject:       Reference{Reference: "Patient/" + patientFHIRID},
		OnsetDateTime: onsetAt.Format(time.RFC3339),
		RecordedDate:  time.Now().Format(time.RFC3339),
	}
	if encounterFHIRID != "" {
		conditionResource.Encounter = &Reference{Reference: "Encounter/" + encounterFHIRID}
	}
	return conditionResource
}
