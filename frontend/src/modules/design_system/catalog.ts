import type { ContrastKind } from "./contrast"

export interface ColorToken {
  name: string
  against: string
  kind: ContrastKind
}

export const COLOR_TOKENS: readonly ColorToken[] = [
  { name: "background", against: "foreground", kind: "text" },
  { name: "card", against: "card-foreground", kind: "text" },
  { name: "primary", against: "primary-foreground", kind: "text" },
  { name: "secondary", against: "secondary-foreground", kind: "text" },
  { name: "muted", against: "muted-foreground", kind: "text" },
  { name: "accent", against: "accent-foreground", kind: "text" },
  { name: "destructive", against: "destructive-foreground", kind: "text" },
  { name: "success", against: "success-foreground", kind: "text" },
  { name: "warning", against: "warning-foreground", kind: "text" },
  { name: "danger", against: "danger-foreground", kind: "text" },
  { name: "border", against: "background", kind: "ui" },
  { name: "ring", against: "background", kind: "ui" },
]

export const TYPOGRAPHY_SCALE = [
  { token: "text-xs", className: "text-xs" },
  { token: "text-sm", className: "text-sm" },
  { token: "text-base", className: "text-base" },
  { token: "text-lg", className: "text-lg" },
  { token: "text-xl", className: "text-xl" },
  { token: "text-2xl", className: "text-2xl" },
  { token: "text-3xl", className: "text-3xl" },
] as const

export const TABULAR_SAMPLE = "0123456789 · 1.234,56 · 98,6 °C"

export const BUTTON_VARIANTS = ["default", "secondary", "outline", "ghost", "destructive", "link"] as const
export const BUTTON_SIZES = ["sm", "default", "lg", "icon"] as const
export const BADGE_VARIANTS = ["default", "secondary", "outline", "success", "warning", "info", "destructive"] as const
export const CLINICAL_STATUSES = ["normal", "attention", "critical", "info", "neutral"] as const
export const ALERT_VARIANTS = ["default", "info", "success", "warning", "destructive"] as const

export const SECTION_IDS = ["colors", "typography", "components", "guidelines"] as const
export const COMPONENT_GROUP_IDS = ["actions", "feedback", "forms", "navigation", "overlays", "data"] as const

export const TABLE_ROWS = ["row1", "row2", "row3"] as const
export const PAGINATION_PAGES = [1, 2, 3] as const
export const GUIDELINE_IDS = ["status", "touch", "i18n", "tokens"] as const
