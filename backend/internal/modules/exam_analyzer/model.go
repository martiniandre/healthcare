package exam_analyzer

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type ExamAnalysis struct {
	ID               uuid.UUID       `db:"id"`
	UserID           *uuid.UUID      `db:"user_id"`
	PatientFhirID    *string         `db:"patient_fhir_id"`
	ExamType         *string         `db:"exam_type"`
	FileName         string          `db:"file_name"`
	FilePath         string          `db:"file_path"`
	Status           string          `db:"status"`
	AnalysisResponse json.RawMessage `db:"analysis_response"`
	ConsentGiven     bool            `db:"consent_given"`
	Anonymized       bool            `db:"anonymized"`
	CreatedAt        time.Time       `db:"created_at"`
	UpdatedAt        time.Time       `db:"updated_at"`
}

type ExamAnalysisAuditLog struct {
	ID          uuid.UUID  `db:"id"`
	AnalysisID  *uuid.UUID `db:"analysis_id"`
	ActionType  string     `db:"action_type"`
	PerformedBy string     `db:"performed_by"`
	IPAddress   *string    `db:"ip_address"`
	Details     *string    `db:"details"`
	CreatedAt   time.Time  `db:"created_at"`
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
