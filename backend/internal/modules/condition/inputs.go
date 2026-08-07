package condition

type CreateConditionInput struct {
	PatientFHIRID   string
	EncounterFHIRID string
	ICD10Code       string
	CodeDisplay     string
	ClinicalStatus  string
}

type UpdateConditionInput struct {
	ICD10Code       *string
	CodeDisplay     *string
	ClinicalStatus  *string
	EncounterFHIRID *string
}
