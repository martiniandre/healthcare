package fhir

import "time"

type AllergyIntoleranceResource struct {
	ResourceType   string              `json:"resourceType"`
	ID             string              `json:"id,omitempty"`
	ClinicalStatus CodeableConcept     `json:"clinicalStatus"`
	Code           CodeableConcept     `json:"code"`
	Patient        Reference           `json:"patient"`
	RecordedDate   string              `json:"recordedDate"`
	Reaction       []AllergyReaction   `json:"reaction,omitempty"`
}

type AllergyReaction struct {
	Manifestation []CodeableConcept `json:"manifestation"`
}

func NewAllergyIntoleranceResource(patientFHIRID, allergenCode, allergenDisplay, clinicalStatus, reactionDescription string) *AllergyIntoleranceResource {
	return &AllergyIntoleranceResource{
		ResourceType: "AllergyIntolerance",
		ClinicalStatus: CodeableConcept{
			Coding: []Coding{
				{System: "http://terminology.hl7.org/CodeSystem/allergyintolerance-clinical", Code: clinicalStatus},
			},
		},
		Code: CodeableConcept{
			Coding: []Coding{
				{System: "http://www.nlm.nih.gov/research/umls/rxnorm", Code: allergenCode, Display: allergenDisplay},
			},
			Text: allergenDisplay,
		},
		Patient:      Reference{Reference: "Patient/" + patientFHIRID},
		RecordedDate: time.Now().Format(time.RFC3339),
		Reaction: []AllergyReaction{
			{
				Manifestation: []CodeableConcept{
					{Text: reactionDescription},
				},
			},
		},
	}
}
