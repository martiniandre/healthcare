import type { ClinicalStatus } from "../../../shared/components/ui/status-badge"

const statusMaps = {
  allergy: { active: "attention", inactive: "neutral", resolved: "normal" },
  condition: {
    active: "attention",
    recurrence: "attention",
    relapse: "attention",
    inactive: "neutral",
    remission: "normal",
    resolved: "normal",
  },
  medication: {
    active: "info",
    "on-hold": "attention",
    completed: "normal",
    stopped: "neutral",
    cancelled: "neutral",
    draft: "neutral",
    unknown: "neutral",
    "entered-in-error": "critical",
  },
  report: {
    final: "normal",
    registered: "info",
    partial: "info",
    preliminary: "info",
    amended: "attention",
    corrected: "attention",
    appended: "attention",
    cancelled: "neutral",
    unknown: "neutral",
    "entered-in-error": "critical",
  },
  encounter: {
    planned: "info",
    arrived: "info",
    triaged: "info",
    "in-progress": "info",
    onleave: "attention",
    finished: "normal",
    cancelled: "neutral",
    unknown: "neutral",
    "entered-in-error": "critical",
  },
} as const satisfies Record<string, Record<string, ClinicalStatus>>

export type ClinicalStatusKind = keyof typeof statusMaps

export const resolveClinicalStatus = (kind: ClinicalStatusKind, value: string): ClinicalStatus => {
  const statusMap: Record<string, ClinicalStatus> = statusMaps[kind]
  return statusMap[value] ?? "neutral"
}

export const ACTIVE_CONDITION_STATUSES: readonly string[] = ["active", "recurrence", "relapse"]
