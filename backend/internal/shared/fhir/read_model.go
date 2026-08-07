package fhir

type ResourceMeta struct {
	LastUpdated string `json:"lastUpdated,omitempty"`
	VersionID   string `json:"versionId,omitempty"`
}

type Encounter struct {
	ResourceType string            `json:"resourceType"`
	ID           string            `json:"id,omitempty"`
	Status       string            `json:"status"`
	Subject      Reference         `json:"subject"`
	Period       *Period           `json:"period,omitempty"`
	ReasonCode   []CodeableConcept `json:"reasonCode,omitempty"`
	Participant  []Participant     `json:"participant,omitempty"`
}

type Observation struct {
	ResourceType      string           `json:"resourceType"`
	ID                string           `json:"id,omitempty"`
	Status            string           `json:"status"`
	Code              CodeableConcept  `json:"code"`
	Subject           Reference        `json:"subject"`
	Encounter         Reference        `json:"encounter"`
	ValueQuantity     *ValueQuantity   `json:"valueQuantity,omitempty"`
	EffectiveDateTime string           `json:"effectiveDateTime,omitempty"`
	Issued            string           `json:"issued,omitempty"`
}

type Condition struct {
	ResourceType   string         `json:"resourceType"`
	ID             string         `json:"id,omitempty"`
	ClinicalStatus CodeableConcept `json:"clinicalStatus"`
	Code           CodeableConcept `json:"code"`
	Subject        Reference       `json:"subject"`
	Encounter      Reference       `json:"encounter"`
	OnsetDateTime  string          `json:"onsetDateTime,omitempty"`
	RecordedDate   string          `json:"recordedDate,omitempty"`
}

type AllergyIntolerance struct {
	ResourceType   string            `json:"resourceType"`
	ID             string            `json:"id,omitempty"`
	ClinicalStatus CodeableConcept   `json:"clinicalStatus"`
	Code           CodeableConcept   `json:"code"`
	Patient        Reference         `json:"patient"`
	RecordedDate   string            `json:"recordedDate,omitempty"`
	Reaction       []AllergyReaction `json:"reaction,omitempty"`
}

type MedicationRequest struct {
	ResourceType              string             `json:"resourceType"`
	ID                        string             `json:"id,omitempty"`
	Status                    string             `json:"status"`
	MedicationCodeableConcept CodeableConcept    `json:"medicationCodeableConcept"`
	Subject                   Reference          `json:"subject"`
	Encounter                 Reference          `json:"encounter"`
	Requester                 Requester          `json:"requester,omitempty"`
	AuthoredOn                string             `json:"authoredOn,omitempty"`
	DosageInstruction         []DosageInstruction `json:"dosageInstruction,omitempty"`
}

type Requester struct {
	Agent Reference `json:"agent"`
}

type DiagnosticReport struct {
	ResourceType string         `json:"resourceType"`
	ID           string         `json:"id,omitempty"`
	Status       string         `json:"status"`
	Code         CodeableConcept `json:"code"`
	Subject      Reference       `json:"subject"`
	Encounter    Reference       `json:"encounter"`
	Issued       string          `json:"issued,omitempty"`
	Conclusion   string          `json:"conclusion,omitempty"`
	Meta         *ResourceMeta   `json:"meta,omitempty"`
}

type Patient struct {
	ResourceType string          `json:"resourceType"`
	ID           string          `json:"id,omitempty"`
	Meta         *ResourceMeta   `json:"meta,omitempty"`
	Name         []HumanName     `json:"name,omitempty"`
	BirthDate    string          `json:"birthDate,omitempty"`
	Telecom      []ContactPoint  `json:"telecom,omitempty"`
	Identifier   []Identifier    `json:"identifier,omitempty"`
}

type ImagingStudy struct {
	ResourceType string        `json:"resourceType"`
	ID           string        `json:"id,omitempty"`
	Status       string        `json:"status"`
	Started      string        `json:"started,omitempty"`
	Description  string        `json:"description,omitempty"`
	Modality     []Coding      `json:"modality,omitempty"`
}
