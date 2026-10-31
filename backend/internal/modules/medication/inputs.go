package medication

type CreateMedicationInput struct {
	EncounterFHIRID       string
	PatientFHIRID         string
	PractitionerFHIRID    string
	MedicationCode        string
	MedicationName        string
	DosageInstructions    string
	OverrideJustification string
}

type CreateMedicationResult struct {
	*Medication
	Warnings []SafetyAlert
}
