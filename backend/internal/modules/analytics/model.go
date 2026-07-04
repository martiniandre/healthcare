package analytics

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

type DashboardData struct {
	ConsultationsToday      int                  `json:"consultations_today"`
	ConsultationsTrend      string               `json:"consultations_trend"`
	OccupancyRate           float64              `json:"occupancy_rate"`
	OccupancyTotalBeds      int                  `json:"occupancy_total_beds"`
	OccupancyOccupiedBeds   int                  `json:"occupancy_occupied_beds"`
	AvgWaitTimeMinutes      float64              `json:"avg_wait_time_minutes"`
	ActivePatients          int                  `json:"active_patients"`
	ExamsToday              int                  `json:"exams_today"`
	NewDiagnosesToday       int                  `json:"new_diagnoses_today"`
	ConsultationsPerDoctor  []DoctorConsultation `json:"consultations_per_doctor"`
	WaitTimeByDepartment    []DepartmentWaitTime `json:"wait_time_by_department"`
	TopDiagnoses            []DiagnosisCount     `json:"top_diagnoses"`
}

type DoctorConsultation struct {
	DoctorName string `json:"doctor_name"`
	Specialty  string `json:"specialty"`
	Count      int    `json:"count"`
}

type OccupancyRate struct {
	Rate         float64 `json:"rate"`
	TotalBeds    int     `json:"total_beds"`
	OccupiedBeds int     `json:"occupied_beds"`
}

type AvgWaitTime struct {
	AverageMinutes float64              `json:"average_minutes"`
	ByDepartment   []DepartmentWaitTime `json:"by_department"`
}

type DepartmentWaitTime struct {
	Department string  `json:"department"`
	Minutes    float64 `json:"minutes"`
}

type DiagnosisCount struct {
	ICD10Code   string `json:"icd10_code"`
	Description string `json:"description"`
	Count       int    `json:"count"`
}
