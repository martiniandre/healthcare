package portal

import "time"

type PortalDashboard struct {
	PatientInfo        PatientInfo          `json:"patient_info"`
	UpcomingEncounters []PortalEncounter    `json:"upcoming_encounters"`
	RecentObservations []PortalObservation  `json:"recent_observations"`
	ActiveConditions   []PortalCondition    `json:"active_conditions"`
	ActiveMedications  []PortalMedication   `json:"active_medications"`
	RecentReports      []PortalReport       `json:"recent_reports"`
	RecentImaging      []PortalImaging      `json:"recent_imaging"`
}

type PatientInfo struct {
	FHIRResourceID string `json:"fhir_resource_id"`
	FullName       string `json:"full_name"`
	BirthDate      string `json:"birth_date"`
	DocumentID     string `json:"document_id"`
}

type PortalEncounter struct {
	FHIRResourceID string     `json:"fhir_resource_id"`
	Status         string     `json:"status"`
	ReasonDisplay  string     `json:"reason_display"`
	StartedAt      time.Time  `json:"started_at"`
	EndedAt        *time.Time `json:"ended_at,omitempty"`
}

type PortalObservation struct {
	FHIRResourceID string    `json:"fhir_resource_id"`
	CodeDisplay    string    `json:"code_display"`
	LoincCode      string    `json:"loinc_code"`
	ValueQuantity  float64   `json:"value_quantity"`
	ValueUnit      string    `json:"value_unit"`
	ObservedAt     time.Time `json:"observed_at"`
}

type PortalCondition struct {
	FHIRResourceID string `json:"fhir_resource_id"`
	CodeDisplay    string `json:"code_display"`
	ICD10Code      string `json:"icd10_code"`
	ClinicalStatus string `json:"clinical_status"`
	OnsetAt        string `json:"onset_at"`
}

type PortalMedication struct {
	FHIRResourceID     string `json:"fhir_resource_id"`
	MedicationName     string `json:"medication_name"`
	DosageInstructions string `json:"dosage_instructions"`
	Status             string `json:"status"`
	IssuedAt           string `json:"issued_at"`
}

type PortalReport struct {
	FHIRResourceID string `json:"fhir_resource_id"`
	ReportDisplay  string `json:"report_display"`
	Status         string `json:"status"`
	Conclusion     string `json:"conclusion"`
	IssuedAt       string `json:"issued_at"`
}

type PortalImaging struct {
	FHIRResourceID string `json:"fhir_resource_id"`
	Title          string `json:"title"`
	Modality       string `json:"modality"`
	Status         string `json:"status"`
	CreatedAt      string `json:"created_at"`
}
