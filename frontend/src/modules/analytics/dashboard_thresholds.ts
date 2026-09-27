export type Severity = "critical" | "elevated" | "stable"

export type MetricTone = Severity | "primary" | "secondary" | "neutral"

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

export const toneClassNames: Record<MetricTone, {
  text: string
  surface: string
  bar: string
}> = {
  critical: {
    text: "text-danger",
    surface: "bg-danger/10",
    bar: "bg-danger",
  },
  elevated: {
    text: "text-warning",
    surface: "bg-warning/10",
    bar: "bg-warning",
  },
  stable: {
    text: "text-success",
    surface: "bg-success/10",
    bar: "bg-success",
  },
  primary: {
    text: "text-primary",
    surface: "bg-primary/10",
    bar: "bg-primary",
  },
  secondary: {
    text: "text-secondary",
    surface: "bg-secondary/10",
    bar: "bg-secondary",
  },
  neutral: {
    text: "text-muted-foreground",
    surface: "bg-gray-100",
    bar: "bg-muted-foreground",
  },
}

export const toneFillVariables: Record<MetricTone, string> = {
  critical: "var(--color-danger)",
  elevated: "var(--color-warning)",
  stable: "var(--color-success)",
  primary: "var(--color-primary)",
  secondary: "var(--color-secondary)",
  neutral: "var(--color-muted-foreground)",
}

export const severityFillVariables: Record<Severity, string> = {
  critical: toneFillVariables.critical,
  elevated: toneFillVariables.elevated,
  stable: toneFillVariables.stable,
}
