package fhir

import "time"

type MedicationRequestResource struct {
	ResourceType               string                `json:"resourceType"`
	ID                         string                `json:"id,omitempty"`
	Status                     string                `json:"status"`
	Intent                     string                `json:"intent"`
	MedicationCodeableConcept  CodeableConcept       `json:"medicationCodeableConcept"`
	Subject                    Reference             `json:"subject"`
	Encounter                  Reference             `json:"encounter"`
	Requester                  Reference             `json:"requester"`
	AuthoredOn                 string                `json:"authoredOn"`
	DosageInstruction          []DosageInstruction   `json:"dosageInstruction"`
}

type DosageInstruction struct {
	Text string `json:"text"`
}

func NewMedicationRequestResource(patientFHIRID, encounterFHIRID, practitionerFHIRID, medicationCode, medicationName, dosageInstructions string) *MedicationRequestResource {
	return &MedicationRequestResource{
		ResourceType: "MedicationRequest",
		Status:       "active",
		Intent:       "order",
		MedicationCodeableConcept: CodeableConcept{
			Coding: []Coding{
				{System: "http://www.nlm.nih.gov/research/umls/rxnorm", Code: medicationCode, Display: medicationName},
			},
			Text: medicationName,
		},
		Subject:   Reference{Reference: "Patient/" + patientFHIRID},
		Encounter: Reference{Reference: "Encounter/" + encounterFHIRID},
		Requester: Reference{Reference: "Practitioner/" + practitionerFHIRID},
		AuthoredOn: time.Now().Format(time.RFC3339),
		DosageInstruction: []DosageInstruction{
			{Text: dosageInstructions},
		},
	}
}
