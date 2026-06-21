package exam_analyzer

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type ExamAnalysis struct {
	ID               uuid.UUID       `db:"id"               json:"id"`
	UserID           *uuid.UUID      `db:"user_id"          json:"user_id,omitempty"`
	PatientFhirID    *string         `db:"patient_fhir_id"  json:"patient_fhir_id,omitempty"`
	ExamType         *string         `db:"exam_type"        json:"exam_type,omitempty"`
	FileName         string          `db:"file_name"        json:"file_name"`
	FilePath         string          `db:"file_path"        json:"file_path"`
	Status           string          `db:"status"           json:"status"`
	AnalysisResponse json.RawMessage `db:"analysis_response" json:"analysis_response"`
	ConsentGiven     bool            `db:"consent_given"    json:"consent_given"`
	Anonymized       bool            `db:"anonymized"       json:"anonymized"`
	CreatedAt        time.Time       `db:"created_at"       json:"created_at"`
	UpdatedAt        time.Time       `db:"updated_at"       json:"updated_at"`
}

type ExamAnalysisAuditLog struct {
	ID          uuid.UUID  `db:"id"           json:"id"`
	AnalysisID  *uuid.UUID `db:"analysis_id"  json:"analysis_id,omitempty"`
	ActionType  string     `db:"action_type"  json:"action_type"`
	PerformedBy string     `db:"performed_by" json:"performed_by"`
	IPAddress   *string    `db:"ip_address"   json:"ip_address,omitempty"`
	Details     *string    `db:"details"      json:"details,omitempty"`
	CreatedAt   time.Time  `db:"created_at"   json:"created_at"`
}

type QualityAssessmentInfo struct {
	Score    float64  `json:"score"`
	Warnings []string `json:"warnings"`
}

type DetectedFindingInfo struct {
	Finding    string  `json:"finding"`
	Confidence float64 `json:"confidence"`
	Severity   string  `json:"severity"`
}

type RecommendationInfo struct {
	Urgency   string   `json:"urgency"`
	NextSteps []string `json:"nextSteps"`
}

type MedicalAnalysisResponse struct {
	ExamType                string                `json:"examType"`
	QualityAssessment       QualityAssessmentInfo `json:"qualityAssessment"`
	DetectedFindings        []DetectedFindingInfo `json:"detectedFindings"`
	PossibleInterpretations []string              `json:"possibleInterpretations"`
	Recommendation          RecommendationInfo    `json:"recommendation"`
	Limitations             []string              `json:"limitations"`
	Disclaimer              string                `json:"disclaimer"`
}
