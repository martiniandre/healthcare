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

export const ExamAnalysisStatus = {
  PENDING: "pending",
  PROCESSING: "processing",
  COMPLETED: "completed",
  FAILED: "failed",
  INSUFFICIENT_DATA: "insufficient_data"
} as const

export type ExamAnalysisStatus = typeof ExamAnalysisStatus[keyof typeof ExamAnalysisStatus]

export interface ExamAnalysis {
  id: string
  user_id?: string
  patient_fhir_id?: string
  exam_type?: string
  file_name: string
  status: ExamAnalysisStatus
  analysis_response: MedicalAnalysisResponse | { status: typeof ExamAnalysisStatus.INSUFFICIENT_DATA; message: string }
  consent_given: boolean
  anonymized: boolean
  anonymizationReport?: AnonymizationReport | null
  created_at: string
  updated_at: string
}

export interface AnonymizationReport {
  applied: boolean
  steps: string[]
  limitations: string[]
}

export const ExamUploadErrorCode = {
  DEIDENTIFICATION_REQUIRED: "DEIDENTIFICATION_REQUIRED",
  FILE_CONTENT_MISMATCH: "FILE_CONTENT_MISMATCH",
  FILE_CORRUPTED: "FILE_CORRUPTED",
  FILE_EMPTY: "FILE_EMPTY",
} as const

export type ExamUploadErrorCode = typeof ExamUploadErrorCode[keyof typeof ExamUploadErrorCode]

export interface UploadExamPayload {
  file: File
  consent: boolean
  anonymize: boolean
  acknowledgeExternalProcessing?: boolean
  patientFhirId?: string
  onUploadProgress?: (progressPercentage: number) => void
  signal?: AbortSignal
}
