package fhir

type ImagingStudyResource struct {
	ResourceType string         `json:"resourceType"`
	ID           string         `json:"id,omitempty"`
	Status       string         `json:"status"`
	Subject      Reference      `json:"subject"`
	Started      string         `json:"started,omitempty"`
	Description  string         `json:"description,omitempty"`
	Series       []Series       `json:"series,omitempty"`
}

type Series struct {
	Uid      string     `json:"uid"`
	Number   int        `json:"number,omitempty"`
	Modality Coding     `json:"modality"`
	BodySite *Coding    `json:"bodySite,omitempty"`
	Instance []Instance `json:"instance,omitempty"`
}

type Instance struct {
	Uid      string `json:"uid"`
	SopClass Coding `json:"sopClass"`
	Number   int    `json:"number,omitempty"`
}

func NewImagingStudyResource(patientFHIRID, status, started, description string, series []Series) *ImagingStudyResource {
	return &ImagingStudyResource{
		ResourceType: "ImagingStudy",
		Status:       status,
		Subject:      Reference{Reference: "Patient/" + patientFHIRID},
		Started:      started,
		Description:  description,
		Series:       series,
	}
}
