package fhir

import "time"

type MedicationRequestResource struct {
	ResourceType              string              `json:"resourceType"`
	ID                        string              `json:"id,omitempty"`
	Status                    string              `json:"status"`
	Intent                    string              `json:"intent"`
	MedicationCodeableConcept CodeableConcept     `json:"medicationCodeableConcept"`
	Subject                   Reference           `json:"subject"`
	Encounter                 *Reference          `json:"encounter,omitempty"`
	Requester                 Reference          `json:"requester,omitempty"`
	AuthoredOn                string              `json:"authoredOn"`
	DosageInstruction         []DosageInstruction `json:"dosageInstruction"`
}

type DosageInstruction struct {
	Text string `json:"text"`
}

func NewMedicationRequestResource(patientFHIRID, encounterFHIRID, practitionerFHIRID, medicationCode, medicationName, dosageInstructions string) *MedicationRequestResource {
	resource := &MedicationRequestResource{
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
		AuthoredOn: time.Now().Format(time.RFC3339),
		DosageInstruction: []DosageInstruction{
			{Text: dosageInstructions},
		},
	}
	if encounterFHIRID != "" {
		resource.Encounter = &Reference{Reference: "Encounter/" + encounterFHIRID}
	}
	if practitionerFHIRID != "" {
		resource.Requester = Reference{Reference: "Practitioner/" + practitionerFHIRID}
	}
	return resource
}
