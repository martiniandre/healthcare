package stats

type Stats struct {
	TotalPatients             int                  `json:"total_patients"`
	FHIRComplianceRate        float64              `json:"fhir_compliance_rate"`
	AvgServiceDurationMinutes float64              `json:"avg_service_duration_minutes"`
	WeeklyConsultations       []WeeklyConsultation `json:"weekly_consultations"`
	ExamModalities            []ExamModality       `json:"exam_modalities"`
	PathologyCases            []PathologyCase      `json:"pathology_cases"`
}

type WeeklyConsultation struct {
	DayName string `json:"dayName"`
	Count   int    `json:"count"`
}

type ExamModality struct {
	Modality   string  `json:"modality"`
	Percentage float64 `json:"percentage"`
	Count      int     `json:"count"`
	Color      string  `json:"color"`
}

type PathologyCase struct {
	Code        string `json:"code"`
	Description string `json:"description"`
	Category    string `json:"category"`
	ActiveCases int    `json:"activeCases"`
	Trend       string `json:"trend"`
}
