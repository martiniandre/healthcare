export interface QualityAssessmentInfo {
  score: number
  warnings: string[]
}

export interface DetectedFindingInfo {
  finding: string
  confidence: number
  severity: "low" | "medium" | "high"
}

export interface RecommendationInfo {
  urgency: "normal" | "medical_followup" | "urgent"
  nextSteps: string[]
}

export interface MedicalAnalysisResponse {
  examType: string
  qualityAssessment: QualityAssessmentInfo
  detectedFindings: DetectedFindingInfo[]
  possibleInterpretations: string[]
  recommendation: RecommendationInfo
  limitations: string[]
  disclaimer: string
}

export interface ExamAnalysis {
  id: string
  user_id?: string
  patient_fhir_id?: string
  exam_type?: string
  file_name: string
  file_path: string
  status: "pending" | "processing" | "completed" | "failed" | "insufficient_data"
  analysis_response: MedicalAnalysisResponse | { status: "insufficient_data"; message: string }
  consent_given: boolean
  anonymized: boolean
  created_at: string
  updated_at: string
}
