package fhir

import "time"

type PatientResource struct {
	ResourceType string              `json:"resourceType"`
	ID           string              `json:"id,omitempty"`
	Name         []HumanName         `json:"name"`
	BirthDate    string              `json:"birthDate"`
	Telecom      []ContactPoint      `json:"telecom,omitempty"`
	Identifier   []Identifier        `json:"identifier"`
	BloodType    string              `json:"-"`
	Allergies    []AllergyEntry      `json:"-"`
}

type HumanName struct {
	Use    string   `json:"use"`
	Family string   `json:"family"`
	Given  []string `json:"given"`
}

type ContactPoint struct {
	System string `json:"system"`
	Value  string `json:"value"`
	Use    string `json:"use"`
}

type Identifier struct {
	System string `json:"system"`
	Value  string `json:"value"`
}

type AllergyEntry struct {
	Substance string    `json:"substance"`
	Reaction  string    `json:"reaction"`
	Recorded  time.Time `json:"recorded"`
}

func NewPatientResource(fullName, documentID, phoneNumber, birthDate string) *PatientResource {
	return &PatientResource{
		ResourceType: "Patient",
		Name: []HumanName{
			{Use: "official", Family: fullName},
		},
		BirthDate: birthDate,
		Telecom: []ContactPoint{
			{System: "phone", Value: phoneNumber, Use: "mobile"},
		},
		Identifier: []Identifier{
			{System: "urn:oid:2.16.840.1.113883.13.237", Value: documentID},
		},
	}
}
