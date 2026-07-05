package fhir

type PractitionerResource struct {
	ResourceType string        `json:"resourceType"`
	ID           string        `json:"id,omitempty"`
	Name         []HumanName   `json:"name"`
	Identifier   []Identifier  `json:"identifier,omitempty"`
	Qualification []Qualification `json:"qualification,omitempty"`
}

type Qualification struct {
	Code CodeableConcept `json:"code"`
}

func NewPractitionerResource(fullName, crmNumber string) *PractitionerResource {
	nameParts := splitName(fullName)
	givenName := nameParts[0]
	familyName := ""
	if len(nameParts) > 1 {
		familyName = nameParts[1]
	}

	var identifier []Identifier
	var qualification []Qualification

	if crmNumber != "" {
		identifier = []Identifier{
			{
				System: "https://systems.digital/sus/identifier/crm",
				Value:  crmNumber,
			},
		}
		qualification = []Qualification{
			{
				Code: CodeableConcept{
					Coding: []Coding{
						{
							System: "http://terminology.hl7.org/CodeSystem/v2-0360",
							Code:   "MD",
							Display: "Medical Doctor",
						},
					},
					Text: "CRM " + crmNumber,
				},
			},
		}
	}

	return &PractitionerResource{
		ResourceType: "Practitioner",
		Name: []HumanName{
			{Use: "official", Family: familyName, Given: []string{givenName}},
		},
		Identifier:   identifier,
		Qualification: qualification,
	}
}

func splitName(fullName string) []string {
	result := make([]string, 0, 2)
	current := make([]rune, 0)
	for _, char := range fullName {
		if char == ' ' {
			if len(current) > 0 {
				result = append(result, string(current))
				current = make([]rune, 0)
			}
		} else {
			current = append(current, char)
		}
	}
	if len(current) > 0 {
		result = append(result, string(current))
	}
	return result
}
