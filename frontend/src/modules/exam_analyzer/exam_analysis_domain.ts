import type { DetectedFindingInfo, ExamAnalysisStatus, RecommendationInfo } from "./types"

export type QualityTone = "good" | "fair" | "poor"

export const QUALITY_VERDICT_THRESHOLDS = {
  good: 0.85,
  fair: 0.6,
} as const

export interface QualityVerdict {
  tone: QualityTone
  translationKey: string
  percentage: number
}

const normalizeScore = (score: number): number => {
  if (Number.isNaN(score)) return 0
  return Math.min(Math.max(score, 0), 1)
}

export const resolveQualityVerdict = (score: number): QualityVerdict => {
  const normalizedScore = normalizeScore(score)

  if (normalizedScore >= QUALITY_VERDICT_THRESHOLDS.good) {
    return {
      tone: "good",
      translationKey: "card.qualityVerdict.good",
      percentage: Math.round(normalizedScore * 100),
    }
  }

  if (normalizedScore >= QUALITY_VERDICT_THRESHOLDS.fair) {
    return {
      tone: "fair",
      translationKey: "card.qualityVerdict.fair",
      percentage: Math.round(normalizedScore * 100),
    }
  }

  return {
    tone: "poor",
    translationKey: "card.qualityVerdict.poor",
    percentage: Math.round(normalizedScore * 100),
  }
}

export type ConfidenceBand = "strong" | "moderate" | "weak"

export const CONFIDENCE_BAND_THRESHOLDS = {
  strong: 0.8,
  moderate: 0.5,
} as const

export interface ConfidenceReading {
  band: ConfidenceBand
  translationKey: string
  percentage: number
}

export const resolveConfidenceBand = (confidence: number): ConfidenceReading => {
  const normalizedConfidence = normalizeScore(confidence)
  const percentage = Math.round(normalizedConfidence * 100)

  if (normalizedConfidence >= CONFIDENCE_BAND_THRESHOLDS.strong) {
    return { band: "strong", translationKey: "card.confidenceBand.strong", percentage }
  }

  if (normalizedConfidence >= CONFIDENCE_BAND_THRESHOLDS.moderate) {
    return { band: "moderate", translationKey: "card.confidenceBand.moderate", percentage }
  }

  return { band: "weak", translationKey: "card.confidenceBand.weak", percentage }
}

export type FindingTriageTier = "priority" | "verify" | "review" | "informative"

export interface FindingTriageReading {
  tier: FindingTriageTier
  translationKey: string
}

export const resolveFindingTriage = (
  severity: DetectedFindingInfo["severity"],
  confidence: number
): FindingTriageReading => {
  const confidenceBand = resolveConfidenceBand(confidence).band
  const isHighSeverity = severity === "high"
  const isMediumSeverity = severity === "medium"
  const isStrongConfidence = confidenceBand === "strong"
  const isAtLeastModerateConfidence = confidenceBand !== "weak"

  if (isHighSeverity && isStrongConfidence) {
    return { tier: "priority", translationKey: "card.triage.priority" }
  }

  if (isHighSeverity) {
    return { tier: "verify", translationKey: "card.triage.verify" }
  }

  if (isMediumSeverity && isAtLeastModerateConfidence) {
    return { tier: "review", translationKey: "card.triage.review" }
  }

  return { tier: "informative", translationKey: "card.triage.informative" }
}

export type UrgencyTone = "urgent" | "followup" | "routine"

export interface UrgencyReading {
  tone: UrgencyTone
  translationKey: string
}

export const resolveUrgencyTone = (
  urgency: RecommendationInfo["urgency"]
): UrgencyReading => {
  if (urgency === "urgent") {
    return { tone: "urgent", translationKey: "card.urgencyUrgent" }
  }

  if (urgency === "medical_followup") {
    return { tone: "followup", translationKey: "card.urgencyFollowup" }
  }

  return { tone: "routine", translationKey: "card.urgencyNormal" }
}

export type HistoryStatusTone = "completed" | "inProgress" | "attention" | "failed"

export interface HistoryStatusReading {
  tone: HistoryStatusTone
  translationKey: string
}

export const resolveHistoryStatus = (status: ExamAnalysisStatus): HistoryStatusReading => {
  if (status === "completed") {
    return { tone: "completed", translationKey: "history.statusCompleted" }
  }

  if (status === "pending") {
    return { tone: "inProgress", translationKey: "history.statusPending" }
  }

  if (status === "processing") {
    return { tone: "inProgress", translationKey: "history.statusProcessing" }
  }

  if (status === "insufficient_data") {
    return { tone: "attention", translationKey: "history.statusQuality" }
  }

  return { tone: "failed", translationKey: "history.statusFailed" }
}
