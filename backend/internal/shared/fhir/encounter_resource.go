package fhir

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
