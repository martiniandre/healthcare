package timeline

import "time"

type TimelineEntry struct {
	ResourceType       string     `json:"resource_type"`
	FHIRResourceID     string     `json:"fhir_resource_id"`
	Title              string     `json:"title"`
	Status             string     `json:"status,omitempty"`
	Summary            string     `json:"summary,omitempty"`
	Code               string     `json:"code,omitempty"`
	RecordedAt         time.Time  `json:"recorded_at"`
	ClinicalDate       *time.Time `json:"clinical_date,omitempty"`
	PeriodEnd          *time.Time `json:"period_end,omitempty"`
	ValueQuantity      *float64   `json:"value_quantity,omitempty"`
	ValueUnit          string     `json:"value_unit,omitempty"`
	ReferenceLow       *float64   `json:"reference_low,omitempty"`
	ReferenceHigh      *float64   `json:"reference_high,omitempty"`
	OnsetDate          string     `json:"onset_date,omitempty"`
	DosageInstructions string     `json:"dosage_instructions,omitempty"`
	Conclusion         string     `json:"conclusion,omitempty"`
	Version            string     `json:"version,omitempty"`
	Modality           string     `json:"modality,omitempty"`
	Reaction           string     `json:"reaction,omitempty"`
	Criticality        string     `json:"criticality,omitempty"`
}

type TimelinePage struct {
	Entries          []TimelineEntry `json:"entries"`
	NextCursor       *time.Time      `json:"next_cursor,omitempty"`
	UnavailableTypes []string        `json:"unavailable_types"`
}
