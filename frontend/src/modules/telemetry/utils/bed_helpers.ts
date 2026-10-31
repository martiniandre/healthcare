import { BedStatus } from "../../../shared/types"
import type { ClinicalStatus } from "@/shared/components/ui/status-badge"

export const BED_HISTORY_LIMIT = 30

export interface BedVitalSample {
  bpm: number
  spo2: number
}

export type BedHistoryMap = Record<string, BedVitalSample[]>

const SEVERITY_RANK: Record<BedStatus, number> = {
  [BedStatus.Danger]: 0,
  [BedStatus.Warning]: 1,
  [BedStatus.Normal]: 2,
}

export const resolveBedClinicalStatus = (
  status: BedStatus
): Extract<ClinicalStatus, "normal" | "attention" | "critical"> => {
  if (status === BedStatus.Danger) {
    return "critical"
  }
  if (status === BedStatus.Warning) {
    return "attention"
  }
  return "normal"
}

export const sortBedsBySeverity = <TBed extends { bedNumber: string; status: BedStatus }>(beds: TBed[]): TBed[] => {
  return beds
    .map((bed, originalIndex) => ({ bed, originalIndex }))
    .sort((left, right) => {
      const severityDifference = SEVERITY_RANK[left.bed.status] - SEVERITY_RANK[right.bed.status]
      if (severityDifference !== 0) {
        return severityDifference
      }
      const numberDifference = left.bed.bedNumber.localeCompare(right.bed.bedNumber, undefined, { numeric: true })
      if (numberDifference !== 0) {
        return numberDifference
      }
      return left.originalIndex - right.originalIndex
    })
    .map(({ bed }) => bed)
}

export const appendBedSamples = (
  history: BedHistoryMap,
  beds: { id: string; bpm: number; spo2: number }[],
  limit: number = BED_HISTORY_LIMIT
): BedHistoryMap => {
  const nextHistory: BedHistoryMap = { ...history }
  for (const bed of beds) {
    const previousSamples = history[bed.id] ?? []
    nextHistory[bed.id] = [...previousSamples, { bpm: bed.bpm, spo2: bed.spo2 }].slice(-limit)
  }
  return nextHistory
}

const WIDTH = 80
const HEIGHT = 20
const PADDING = 2

export const buildSparklinePoints = (values: number[]): string => {
  if (values.length === 0) {
    return ""
  }
  const minimum = Math.min(...values)
  const maximum = Math.max(...values)
  const range = maximum - minimum
  const horizontalStep = values.length > 1 ? (WIDTH - PADDING * 2) / (values.length - 1) : 0

  return values
    .map((value, index) => {
      const horizontal = PADDING + index * horizontalStep
      const normalized = range === 0 ? 0.5 : (value - minimum) / range
      const vertical = HEIGHT - PADDING - normalized * (HEIGHT - PADDING * 2)
      return `${horizontal.toFixed(1)},${vertical.toFixed(1)}`
    })
    .join(" ")
}
