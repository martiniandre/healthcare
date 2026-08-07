package allergy

type CreateAllergyInput struct {
	PatientFHIRID   string
	AllergenCode    string
	AllergenDisplay string
	ClinicalStatus  string
	Reaction        string
}

type UpdateAllergyInput struct {
	AllergenCode    *string
	AllergenDisplay *string
	ClinicalStatus  *string
	Reaction        *string
}
