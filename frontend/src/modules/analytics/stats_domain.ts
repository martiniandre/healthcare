import type { MetricTone, Severity } from "./metric_tones"
import type { ConsultationsDayData } from "./types"

const FHIR_COMPLIANCE_THRESHOLDS = {
  critical: 90,
  elevated: 98,
} as const

export const resolveFhirComplianceSeverity = (complianceRate: number): Severity => {
  if (complianceRate < FHIR_COMPLIANCE_THRESHOLDS.critical) return "critical"
  if (complianceRate < FHIR_COMPLIANCE_THRESHOLDS.elevated) return "elevated"
  return "stable"
}

export type TrendDirection = "rising" | "falling" | "flat" | "unknown"

const RISING_PREFIXES = ["+"]
const FALLING_PREFIXES = ["-", "\u2212"]

const TREND_MAGNITUDE_THRESHOLDS = {
  critical: 10,
  elevated: 5,
} as const

const parseTrendMagnitude = (trend: string): number | null => {
  const numericCharacters = trend.replace(/[^\d.]/g, "")
  if (numericCharacters.length === 0) return null

  const parsedMagnitude = Number.parseFloat(numericCharacters)
  return Number.isFinite(parsedMagnitude) ? Math.abs(parsedMagnitude) : null
}

export interface TrendReading {
  direction: TrendDirection
  magnitude: number | null
  severity: MetricTone
}

export const resolveTrendDirection = (trend: string): TrendDirection => {
  const normalizedTrend = trend.trim()
  if (normalizedTrend.length === 0) return "unknown"

  const leadingCharacter = normalizedTrend.charAt(0)
  if (RISING_PREFIXES.includes(leadingCharacter)) return "rising"
  if (FALLING_PREFIXES.includes(leadingCharacter)) return "falling"
  return "flat"
}

export const resolveTrendReading = (trend: string): TrendReading => {
  const direction = resolveTrendDirection(trend)
  const magnitude = parseTrendMagnitude(trend)

  if (direction === "falling") return { direction, magnitude, severity: "stable" }
  if (direction !== "rising" || magnitude === null) {
    return { direction, magnitude, severity: "neutral" }
  }
  if (magnitude >= TREND_MAGNITUDE_THRESHOLDS.critical) {
    return { direction, magnitude, severity: "critical" }
  }
  if (magnitude >= TREND_MAGNITUDE_THRESHOLDS.elevated) {
    return { direction, magnitude, severity: "elevated" }
  }
  return { direction, magnitude, severity: "neutral" }
}

export const resolveDayTranslationKey = (dayName: string): string => {
  const segments = dayName.trim().toLowerCase().split(".")
  return `days.${segments[segments.length - 1]}`
}

export interface WeeklyConsultationsSummary {
  total: number
  min: number
  average: number
  peak: number
  peakDayKey: string | null
}

export const summarizeWeeklyConsultations = (
  weeklyConsultations: ConsultationsDayData[]
): WeeklyConsultationsSummary => {
  if (weeklyConsultations.length === 0) {
    return { total: 0, min: 0, average: 0, peak: 0, peakDayKey: null }
  }

  const counts = weeklyConsultations.map((consultation) => consultation.count)
  const total = counts.reduce((sum, count) => sum + count, 0)
  const peak = Math.max(...counts)
  const peakConsultation = weeklyConsultations.find((consultation) => consultation.count === peak)

  return {
    total,
    min: Math.min(...counts),
    average: Math.round(total / weeklyConsultations.length),
    peak,
    peakDayKey: peak > 0 && peakConsultation ? resolveDayTranslationKey(peakConsultation.dayName) : null,
  }
}
