package timeline

import "time"

type EntryType string

const (
	EntryTypeEncounter   EntryType = "encounter"
	EntryTypeObservation EntryType = "observation"
	EntryTypeCondition   EntryType = "condition"
	EntryTypeAllergy     EntryType = "allergy"
	EntryTypeMedication  EntryType = "medication"
	EntryTypeReport      EntryType = "report"
	EntryTypeImaging     EntryType = "imaging"
)

const (
	DefaultLimit        = 30
	MaxLimit            = 100
	DefaultWindowMonths = 6
)

var AllEntryTypes = []EntryType{
	EntryTypeEncounter,
	EntryTypeObservation,
	EntryTypeCondition,
	EntryTypeAllergy,
	EntryTypeMedication,
	EntryTypeReport,
	EntryTypeImaging,
}

type Entry struct {
	Type         EntryType
	ResourceID   string
	OccurredAt   time.Time
	ClinicalDate *time.Time
	Title        string
	Summary      string
	Status       string
	EncounterID  string
}

func (entry Entry) ID() string {
	return string(entry.Type) + ":" + entry.ResourceID
}

type Window struct {
	From time.Time
	To   time.Time
}

type ListInput struct {
	PatientFHIRID string
	From          string
	To            string
	Types         string
	Cursor        string
	Limit         string
}

type FetchQuery struct {
	PatientFHIRID string
	Window        Window
	MaxEntries    int
}

type Page struct {
	Entries     []Entry
	NextCursor  string
	FailedTypes []EntryType
	Window      Window
}

type ItemResponse struct {
	ID           string  `json:"id"`
	Type         string  `json:"type"`
	OccurredAt   string  `json:"occurredAt"`
	ClinicalDate *string `json:"clinicalDate"`
	Title        string  `json:"title"`
	Summary      string  `json:"summary"`
	Status       string  `json:"status"`
	ResourceID   string  `json:"resourceId"`
	EncounterID  *string `json:"encounterId"`
}

type PartialResponse struct {
	FailedTypes []string `json:"failedTypes"`
}

type WindowResponse struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type TimelineResponse struct {
	Items      []ItemResponse  `json:"items"`
	NextCursor *string         `json:"nextCursor"`
	Partial    PartialResponse `json:"partial"`
	Window     WindowResponse  `json:"window"`
}
