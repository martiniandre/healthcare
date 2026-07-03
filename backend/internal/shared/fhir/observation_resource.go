package fhir

import (
	"fmt"
	"time"
)

type ObservationResource struct {
	ResourceType string              `json:"resourceType"`
	ID           string              `json:"id,omitempty"`
	Status       string              `json:"status"`
	Category     []CodeableConcept   `json:"category"`
	Code         CodeableConcept     `json:"code"`
	Subject      Reference           `json:"subject"`
	EffectiveDateTime string         `json:"effectiveDateTime"`
	ValueQuantity *ValueQuantity     `json:"valueQuantity,omitempty"`
	ValueString  string              `json:"valueString,omitempty"`
}

type CodeableConcept struct {
	Coding []Coding `json:"coding"`
	Text   string   `json:"text,omitempty"`
}

type Coding struct {
	System  string `json:"system"`
	Code    string `json:"code"`
	Display string `json:"display,omitempty"`
}

type Reference struct {
	Reference string `json:"reference"`
}

type ValueQuantity struct {
	Value  float64 `json:"value"`
	Unit   string  `json:"unit"`
	System string  `json:"system"`
	Code   string  `json:"code"`
}

func NewObservationResource(patientFHIRID, encounterFHIRID, loincCode, codeDisplay string, valueQuantity float64, valueUnit string) *ObservationResource {
	return &ObservationResource{
		ResourceType: "Observation",
		Status:       "final",
		Category: []CodeableConcept{
			{
				Coding: []Coding{
					{
						System:  "http://terminology.hl7.org/CodeSystem/observation-category",
						Code:    "vital-signs",
						Display: "Vital Signs",
					},
				},
			},
		},
		Code: CodeableConcept{
			Coding: []Coding{
				{
					System:  "http://loinc.org",
					Code:    loincCode,
					Display: codeDisplay,
				},
			},
			Text: codeDisplay,
		},
		Subject: Reference{
			Reference: fmt.Sprintf("Patient/%s", patientFHIRID),
		},
		EffectiveDateTime: time.Now().Format(time.RFC3339),
		ValueQuantity: &ValueQuantity{
			Value:  valueQuantity,
			Unit:   valueUnit,
			System: "http://unitsofmeasure.org",
			Code:   valueUnit,
		},
	}
}

