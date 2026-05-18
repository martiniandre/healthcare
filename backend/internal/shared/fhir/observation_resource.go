package fhir

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
