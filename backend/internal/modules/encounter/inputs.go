package encounter

type CreateEncounterInput struct {
	PatientFHIRID  string
	PractitionerID string
	ReasonCode     string
	ReasonDisplay  string
}

type UpdateEncounterInput struct {
	ReasonCode     *string
	ReasonDisplay  *string
	PractitionerID *string
	Status         *string
}
