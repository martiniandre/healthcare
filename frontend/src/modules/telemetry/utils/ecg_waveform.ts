import { BedStatus, CardiacCondition } from "@/shared/types"

export const ECG_COLORS = {
  background: "#090d16",
  grid: "rgba(16, 185, 129, 0.15)",
  trace: "#10b981",
  arrestTrace: "#ef4444",
} as const

export const ECG_GRID_SIZE = 15
export const ECG_PLOT_POINTS = 300

const DEFAULT_CYCLE_PERIOD = 60

export const cyclePeriodFor = (condition: CardiacCondition | undefined): number => {
  if (condition === CardiacCondition.Bradycardia) return 110
  if (condition === CardiacCondition.Tachycardia) return 35
  return DEFAULT_CYCLE_PERIOD
}

export const computeEcgSample = (
  traceOffset: number,
  baseline: number,
  condition: CardiacCondition | undefined
): number => {
  if (condition === CardiacCondition.CardiacArrest) {
    return baseline + Math.sin(traceOffset * 0.15) * 0.8
  }

  const phase = traceOffset % cyclePeriodFor(condition)

  if (phase < 6) return baseline
  if (phase < 12) return baseline - Math.sin(((phase - 6) * Math.PI) / 6) * 12
  if (phase < 16) return baseline
  if (phase < 18) return baseline + 8
  if (phase < 21) return baseline + 8 - ((phase - 18) / 3) * 65
  if (phase < 24) return baseline - 57 + ((phase - 21) / 3) * 78
  if (phase < 27) return baseline + 21 - ((phase - 24) / 3) * 21
  if (phase < 33) return baseline
  if (phase < 43) return baseline - Math.sin(((phase - 33) * Math.PI) / 10) * 18
  return baseline
}

export interface SimulatedVitals {
  bpm: number
  spo2: number
  status: BedStatus
}

export const simulatedVitalsFor = (condition: CardiacCondition): SimulatedVitals => {
  if (condition === CardiacCondition.Bradycardia) return { bpm: 45, spo2: 93, status: BedStatus.Warning }
  if (condition === CardiacCondition.Tachycardia) return { bpm: 135, spo2: 95, status: BedStatus.Warning }
  if (condition === CardiacCondition.CardiacArrest) return { bpm: 0, spo2: 0, status: BedStatus.Danger }
  return { bpm: 75, spo2: 98, status: BedStatus.Normal }
}
