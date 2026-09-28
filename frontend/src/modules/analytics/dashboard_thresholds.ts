import type { Severity } from "./metric_tones"

const OCCUPANCY_THRESHOLDS = {
  critical: 85,
  elevated: 65,
} as const

const WAIT_TIME_THRESHOLDS = {
  critical: 45,
  elevated: 25,
} as const

export const resolveOccupancySeverity = (occupancyRate: number): Severity => {
  if (occupancyRate >= OCCUPANCY_THRESHOLDS.critical) return "critical"
  if (occupancyRate >= OCCUPANCY_THRESHOLDS.elevated) return "elevated"
  return "stable"
}

export const resolveWaitTimeSeverity = (minutes: number): Severity => {
  if (minutes >= WAIT_TIME_THRESHOLDS.critical) return "critical"
  if (minutes >= WAIT_TIME_THRESHOLDS.elevated) return "elevated"
  return "stable"
}
