package fhir

import (
	"fmt"
	"time"
)

type EncounterResource struct {
	ResourceType string            `json:"resourceType"`
	ID           string            `json:"id,omitempty"`
	Status       string            `json:"status"`
	Class        Coding            `json:"class"`
	Subject      Reference         `json:"subject"`
	Participant  []Participant     `json:"participant,omitempty"`
	Period       *Period           `json:"period,omitempty"`
	ReasonCode   []CodeableConcept `json:"reasonCode,omitempty"`
}

type Participant struct {
	Type       []CodeableConcept `json:"type"`
	Individual Reference         `json:"individual"`
}

type Period struct {
	Start string `json:"start"`
	End   string `json:"end,omitempty"`
}

func NewEncounterResource(patientFHIRID, practitionerID, reasonCode, reasonDisplay string) *EncounterResource {
	var participantsArray []Participant
	if practitionerID != "" {
		participantsArray = []Participant{
			{
				Type: []CodeableConcept{
					{
						Coding: []Coding{
							{
								System:  "http://terminology.hl7.org/CodeSystem/v3-ParticipationType",
								Code:    "PPRF",
								Display: "primary performer",
							},
						},
					},
				},
				Individual: Reference{
					Reference: fmt.Sprintf("Practitioner/%s", practitionerID),
				},
			},
		}
	}

	var reasonsCodeableConcepts []CodeableConcept
	if reasonCode != "" {
		reasonsCodeableConcepts = []CodeableConcept{
			{
				Coding: []Coding{
					{
						System:  "http://hl7.org/fhir/sid/icd-10",
						Code:    reasonCode,
						Display: reasonDisplay,
					},
				},
				Text: reasonDisplay,
			},
		}
	}

	return &EncounterResource{
		ResourceType: "Encounter",
		Status:       "finished",
		Class: Coding{
			System:  "http://terminology.hl7.org/CodeSystem/v3-ActCode",
			Code:    "AMB",
			Display: "ambulatory",
		},
		Subject: Reference{
			Reference: fmt.Sprintf("Patient/%s", patientFHIRID),
		},
		Participant: participantsArray,
		ReasonCode:  reasonsCodeableConcepts,
		Period: &Period{
			Start: time.Now().Format(time.RFC3339),
			End:   time.Now().Format(time.RFC3339),
		},
	}
}

