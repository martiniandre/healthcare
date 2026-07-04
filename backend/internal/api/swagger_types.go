package api

type ErrorResponse struct {
	Error string `json:"error" example:"error message"`
}

type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Total      int         `json:"total"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
}

// Auth
type LoginRequest struct {
	Email    string `json:"email" example:"user@hospital.com"`
	Password string `json:"password" example:"securepassword"`
}

type LoginResponse struct {
	UserID string `json:"userId" example:"550e8400-e29b-41d4-a716-446655440000"`
	Role   string `json:"role" example:"doctor"`
	Email  string `json:"email" example:"user@hospital.com"`
}

type MeResponse struct {
	UserID    string `json:"userId" example:"550e8400-e29b-41d4-a716-446655440000"`
	Email     string `json:"email" example:"user@hospital.com"`
	FullName  string `json:"fullName" example:"John Doe"`
	Role      string `json:"role" example:"doctor"`
	IsActive  bool   `json:"isActive" example:"true"`
	CreatedAt string `json:"createdAt" example:"2024-01-01T00:00:00Z"`
	UpdatedAt string `json:"updatedAt" example:"2024-01-01T00:00:00Z"`
}

type LogoutResponse struct {
	Message string `json:"message" example:"Logged out successfully"`
}

// Patients
type PatientResponse struct {
	PatientID      string `json:"patient_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	FHIRResourceID string `json:"fhir_resource_id" example:"fhir-12345"`
	FullName       string `json:"full_name" example:"John Doe"`
	BirthDate      string `json:"birth_date" example:"1990-01-15"`
	DocumentID     string `json:"document_id" example:"12345678900"`
	PhoneNumber    string `json:"phone_number" example:"+5511999999999"`
}

type CreatePatientRequest struct {
	FullName    string `json:"full_name" example:"John Doe"`
	BirthDate   string `json:"birth_date" example:"1990-01-15"`
	DocumentID  string `json:"document_id" example:"12345678900"`
	PhoneNumber string `json:"phone_number" example:"+5511999999999"`
}

type CreatePatientResponse struct {
	PatientID      string `json:"patient_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	FHIRResourceID string `json:"fhir_resource_id" example:"fhir-12345"`
}

type PatientDetailResponse struct {
	PatientID      string `json:"patient_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	FHIRResourceID string `json:"fhir_resource_id" example:"fhir-12345"`
	FullName       string `json:"full_name" example:"John Doe"`
	BirthDate      string `json:"birth_date" example:"1990-01-15"`
	DocumentID     string `json:"document_id" example:"12345678900"`
	PhoneNumber    string `json:"phone_number" example:"+5511999999999"`
}

// Encounters
type EncounterResponse struct {
	FhirID         string `json:"fhir_id" example:"encounter-fhir-123"`
	PatientFhirID  string `json:"patient_fhir_id" example:"patient-fhir-123"`
	Status         string `json:"status" example:"finished"`
	ReasonDisplay  string `json:"reason_display" example:"Routine checkup"`
	PractitionerID string `json:"practitioner_id,omitempty" example:"practitioner-fhir-456"`
	CreatedAt      string `json:"created_at" example:"2024-01-01T10:00:00Z"`
}

type CreateEncounterRequest struct {
	ReasonDisplay  string `json:"reason_display" example:"Routine checkup"`
	PractitionerID string `json:"practitioner_id" example:"practitioner-fhir-456"`
}

type CreateEncounterResponse struct {
	FhirID         string `json:"fhir_id" example:"encounter-fhir-123"`
	PatientFhirID  string `json:"patient_fhir_id" example:"patient-fhir-123"`
	Status         string `json:"status" example:"finished"`
	ReasonDisplay  string `json:"reason_display" example:"Routine checkup"`
	PractitionerID string `json:"practitioner_id" example:"practitioner-fhir-456"`
	CreatedAt      string `json:"created_at" example:"2024-01-01T10:00:00Z"`
}

// Observations
type ObservationResponse struct {
	FhirID          string  `json:"fhir_id" example:"obs-fhir-123"`
	EncounterFhirID string  `json:"encounter_fhir_id" example:"encounter-fhir-123"`
	PatientFhirID   string  `json:"patient_fhir_id" example:"patient-fhir-123"`
	LoincCode       string  `json:"loinc_code" example:"8867-4"`
	CodeDisplay     string  `json:"code_display" example:"Heart rate"`
	ValueQuantity   float64 `json:"value_quantity" example:"72.0"`
	ValueUnit       string  `json:"value_unit" example:"bpm"`
	CreatedAt       string  `json:"created_at" example:"2024-01-01T10:00:00Z"`
}

type CreateObservationRequest struct {
	PatientFhirID string  `json:"patient_fhir_id" example:"patient-fhir-123"`
	LoincCode     string  `json:"loinc_code" example:"8867-4"`
	CodeDisplay   string  `json:"code_display" example:"Heart rate"`
	ValueQuantity float64 `json:"value_quantity" example:"72.0"`
	ValueUnit     string  `json:"value_unit" example:"bpm"`
}

type CreateObservationResponse struct {
	FhirID          string  `json:"fhir_id" example:"obs-fhir-123"`
	EncounterFhirID string  `json:"encounter_fhir_id" example:"encounter-fhir-123"`
	PatientFhirID   string  `json:"patient_fhir_id" example:"patient-fhir-123"`
	LoincCode       string  `json:"loinc_code" example:"8867-4"`
	CodeDisplay     string  `json:"code_display" example:"Heart rate"`
	ValueQuantity   float64 `json:"value_quantity" example:"72.0"`
	ValueUnit       string  `json:"value_unit" example:"bpm"`
	CreatedAt       string  `json:"created_at" example:"2024-01-01T10:00:00Z"`
}

// Conditions
type ConditionResponse struct {
	FhirID         string `json:"fhir_id" example:"cond-fhir-123"`
	PatientFhirID  string `json:"patient_fhir_id" example:"patient-fhir-123"`
	ICD10Code      string `json:"icd10_code" example:"I10"`
	CodeDisplay    string `json:"code_display" example:"Essential hypertension"`
	ClinicalStatus string `json:"clinical_status" example:"active"`
	CreatedAt      string `json:"created_at" example:"2024-01-01T10:00:00Z"`
}

type CreateConditionRequest struct {
	ICD10Code   string `json:"icd10_code" example:"I10"`
	CodeDisplay string `json:"code_display" example:"Essential hypertension"`
}

type CreateConditionResponse struct {
	FhirID         string `json:"fhir_id" example:"cond-fhir-123"`
	PatientFhirID  string `json:"patient_fhir_id" example:"patient-fhir-123"`
	ICD10Code      string `json:"icd10_code" example:"I10"`
	CodeDisplay    string `json:"code_display" example:"Essential hypertension"`
	ClinicalStatus string `json:"clinical_status" example:"active"`
	CreatedAt      string `json:"created_at" example:"2024-01-01T10:00:00Z"`
}

// Allergies
type AllergyResponse struct {
	FhirID          string `json:"fhir_id" example:"allergy-fhir-123"`
	PatientFhirID   string `json:"patient_fhir_id" example:"patient-fhir-123"`
	AllergenCode    string `json:"allergen_code" example:"J30.1"`
	AllergenDisplay string `json:"allergen_display" example:"Peanut"`
	ClinicalStatus  string `json:"clinical_status" example:"active"`
	Reaction        string `json:"reaction" example:"Anaphylaxis"`
	CreatedAt       string `json:"created_at" example:"2024-01-01T10:00:00Z"`
}

type CreateAllergyRequest struct {
	AllergenCode    string `json:"allergen_code" example:"J30.1"`
	AllergenDisplay string `json:"allergen_display" example:"Peanut"`
	Reaction        string `json:"reaction" example:"Anaphylaxis"`
}

type CreateAllergyResponse struct {
	FhirID          string `json:"fhir_id" example:"allergy-fhir-123"`
	PatientFhirID   string `json:"patient_fhir_id" example:"patient-fhir-123"`
	AllergenCode    string `json:"allergen_code" example:"J30.1"`
	AllergenDisplay string `json:"allergen_display" example:"Peanut"`
	ClinicalStatus  string `json:"clinical_status" example:"active"`
	Reaction        string `json:"reaction" example:"Anaphylaxis"`
	CreatedAt       string `json:"created_at" example:"2024-01-01T10:00:00Z"`
}

// Medications
type MedicationResponse struct {
	FhirID             string `json:"fhir_id" example:"med-fhir-123"`
	EncounterFhirID    string `json:"encounter_fhir_id" example:"encounter-fhir-123"`
	PatientFhirID      string `json:"patient_fhir_id" example:"patient-fhir-123"`
	PractitionerFhirID string `json:"practitioner_fhir_id" example:"practitioner-fhir-456"`
	MedicationCode     string `json:"medication_code" example:"12345"`
	MedicationName     string `json:"medication_name" example:"Amoxicillin"`
	DosageInstructions string `json:"dosage_instructions" example:"500mg three times daily"`
	Status             string `json:"status" example:"active"`
	CreatedAt          string `json:"created_at" example:"2024-01-01T10:00:00Z"`
}

type CreateMedicationRequest struct {
	PatientFhirID      string `json:"patient_fhir_id" example:"patient-fhir-123"`
	PractitionerFhirID string `json:"practitioner_fhir_id" example:"practitioner-fhir-456"`
	MedicationCode     string `json:"medication_code" example:"12345"`
	MedicationName     string `json:"medication_name" example:"Amoxicillin"`
	DosageInstructions string `json:"dosage_instructions" example:"500mg three times daily"`
}

type CreateMedicationResponse struct {
	FhirID             string `json:"fhir_id" example:"med-fhir-123"`
	EncounterFhirID    string `json:"encounter_fhir_id" example:"encounter-fhir-123"`
	PatientFhirID      string `json:"patient_fhir_id" example:"patient-fhir-123"`
	PractitionerFhirID string `json:"practitioner_fhir_id" example:"practitioner-fhir-456"`
	MedicationCode     string `json:"medication_code" example:"12345"`
	MedicationName     string `json:"medication_name" example:"Amoxicillin"`
	DosageInstructions string `json:"dosage_instructions" example:"500mg three times daily"`
	Status             string `json:"status" example:"active"`
	CreatedAt          string `json:"created_at" example:"2024-01-01T10:00:00Z"`
}

// Diagnostic Reports
type DiagnosticReportResponse struct {
	FhirID          string `json:"fhir_id" example:"report-fhir-123"`
	EncounterFhirID string `json:"encounter_fhir_id" example:"encounter-fhir-123"`
	PatientFhirID   string `json:"patient_fhir_id" example:"patient-fhir-123"`
	ReportDisplay   string `json:"report_display" example:"Chest X-ray Report"`
	Status          string `json:"status" example:"final"`
	Conclusion      string `json:"conclusion" example:"Normal findings"`
	CreatedAt       string `json:"created_at" example:"2024-01-01T10:00:00Z"`
}

type CreateDiagnosticReportRequest struct {
	PatientFhirID string `json:"patient_fhir_id" example:"patient-fhir-123"`
	ReportDisplay string `json:"report_display" example:"Chest X-ray Report"`
	Conclusion    string `json:"conclusion" example:"Normal findings"`
}

type CreateDiagnosticReportResponse struct {
	FhirID          string `json:"fhir_id" example:"report-fhir-123"`
	EncounterFhirID string `json:"encounter_fhir_id" example:"encounter-fhir-123"`
	PatientFhirID   string `json:"patient_fhir_id" example:"patient-fhir-123"`
	ReportDisplay   string `json:"report_display" example:"Chest X-ray Report"`
	Status          string `json:"status" example:"final"`
	Conclusion      string `json:"conclusion" example:"Normal findings"`
	CreatedAt       string `json:"created_at" example:"2024-01-01T10:00:00Z"`
}

// Staff
type EmployeeResponse struct {
	ID        string `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	FullName  string `json:"full_name" example:"Dr. John Doe"`
	Email     string `json:"email" example:"john.doe@hospital.com"`
	Role      string `json:"role" example:"doctor"`
	CRMNumber string `json:"crm_number" example:"CRM-SP-123456"`
}

type CreateEmployeeRequest struct {
	UserID    string `json:"user_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	FullName  string `json:"full_name" example:"Dr. John Doe"`
	Email     string `json:"email" example:"john.doe@hospital.com"`
	Role      string `json:"role" example:"doctor"`
	CRMNumber string `json:"crm_number" example:"CRM-SP-123456"`
}

type CreateEmployeeResponse struct {
	EmployeeID string `json:"employee_id" example:"550e8400-e29b-41d4-a716-446655440000"`
}

// Telemetry
type TelemetryRoomResponse struct {
	ID             string `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name           string `json:"name" example:"Room 101"`
	Status         string `json:"status" example:"locked"`
}

type UnlockRoomRequest struct {
	Passcode string `json:"passcode" example:"1234"`
}

type UnlockRoomResponse struct {
	Success  bool   `json:"success" example:"true"`
	RoomName string `json:"roomName" example:"Room 101"`
}

type TelemetryBedResponse struct {
	ID          string  `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	RoomID      string  `json:"room_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	BedLabel    string  `json:"bed_label" example:"A"`
	Bpm         int32   `json:"bpm" example:"72"`
	Spo2        int32   `json:"spo2" example:"98"`
	Temperature float64 `json:"temperature" example:"36.5"`
	Status      string  `json:"status" example:"occupied"`
	Condition   string  `json:"condition" example:"stable"`
}

type UpdateBedConditionRequest struct {
	Bpm         int32   `json:"bpm" example:"72"`
	Spo2        int32   `json:"spo2" example:"98"`
	Temperature float64 `json:"temperature" example:"36.5"`
	Status      string  `json:"status" example:"occupied"`
	Condition   string  `json:"condition" example:"stable"`
}

type UpdateBedConditionResponse struct {
	Success bool `json:"success" example:"true"`
}

// Audit Logs
type AuditLogEntry struct {
	ID             string `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	CorrelationID  string `json:"correlation_id" example:"corr-123"`
	CallerUserID   string `json:"caller_user_id" example:"user-123"`
	CallerRole     string `json:"caller_role" example:"admin"`
	Method         string `json:"method" example:"GET /api/patients"`
	AccessGranted  bool   `json:"access_granted" example:"true"`
	CreatedAt      string `json:"created_at" example:"2024-01-01T10:00:00Z"`
}

type AuditLogListResponse struct {
	AuditLogs []AuditLogEntry `json:"audit_logs"`
	Total     int             `json:"total" example:"100"`
}

type CreateAuditLogRequest struct {
	CorrelationID string `json:"correlation_id" example:"corr-123"`
	CallerUserID  string `json:"caller_user_id" example:"user-123"`
	CallerRole    string `json:"caller_role" example:"admin"`
	Method        string `json:"method" example:"GET /api/patients"`
	AccessGranted bool   `json:"access_granted" example:"true"`
}

// Exam Analyzer
type ExamAnalysisResponse struct {
	ID               string `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	UserID           string `json:"user_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	PatientFhirID    string `json:"patient_fhir_id" example:"patient-fhir-123"`
	ExamType         string `json:"exam_type" example:"X-Ray"`
	FileName         string `json:"file_name" example:"xray.dcm"`
	Status           string `json:"status" example:"pending"`
	AnalysisResponse string `json:"analysis_response" example:"{\"status\":\"pending\"}"`
	ConsentGiven     bool   `json:"consent_given" example:"true"`
	Anonymized       bool   `json:"anonymized" example:"false"`
	CreatedAt        string `json:"created_at" example:"2024-01-01T10:00:00Z"`
	UpdatedAt        string `json:"updated_at" example:"2024-01-01T10:00:00Z"`
}

type DeleteAnalysisResponse struct {
	Success string `json:"success" example:"Análise e arquivo excluídos com sucesso."`
}
