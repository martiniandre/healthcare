import { isAxiosError } from "axios"
import type { ClinicalStatus } from "../../shared/components/ui/status-badge"
import type { MedicationAlertSeverity, MedicationSafetyAlert } from "./types"

export const MEDICATION_ALLERGY_BLOCKED_CODE = "MEDICATION_ALLERGY_BLOCKED"

const SEVERITY_STATUS: Record<MedicationAlertSeverity, ClinicalStatus> = {
  info: "info",
  moderate: "attention",
  major: "critical",
  contraindicated: "critical",
}

export const resolveAlertBadgeStatus = (alert: MedicationSafetyAlert): ClinicalStatus => {
  if (alert.kind === "allergy") return "critical"
  return SEVERITY_STATUS[alert.severity] ?? "attention"
}

export const hasBlockingAlert = (alerts: MedicationSafetyAlert[]): boolean =>
  alerts.some((alert) => alert.blocking)

export const sortAlertsByKind = (alerts: MedicationSafetyAlert[]): MedicationSafetyAlert[] =>
  [...alerts].sort((firstAlert, secondAlert) => {
    if (firstAlert.kind === secondAlert.kind) return 0
    return firstAlert.kind === "allergy" ? -1 : 1
  })

const buildAlertIdentity = (alert: MedicationSafetyAlert): string =>
  `${alert.kind}|${alert.messageKey}|${alert.relatedResourceId}`

export const mergeAlerts = (
  firstAlerts: MedicationSafetyAlert[],
  secondAlerts: MedicationSafetyAlert[],
): MedicationSafetyAlert[] => {
  const alertsByIdentity = new Map<string, MedicationSafetyAlert>()
  for (const alert of [...firstAlerts, ...secondAlerts]) {
    alertsByIdentity.set(buildAlertIdentity(alert), alert)
  }
  return [...alertsByIdentity.values()]
}

export const extractAllergyBlockedAlerts = (error: unknown): MedicationSafetyAlert[] | null => {
  if (!isAxiosError(error) || error.response?.status !== 409) return null
  const responseData = error.response.data as { code?: string; alerts?: MedicationSafetyAlert[] } | undefined
  if (responseData?.code !== MEDICATION_ALLERGY_BLOCKED_CODE) return null
  return Array.isArray(responseData.alerts) ? responseData.alerts : []
}

export type AlertTranslateFunction = (key: string, options?: Record<string, unknown>) => string

export const buildAlertMessage = (translate: AlertTranslateFunction, alert: MedicationSafetyAlert): string =>
  translate(`medicationSafety.messages.${alert.messageKey}`, {
    ...alert.params,
    defaultValue: translate("medicationSafety.messages.generic", alert.params),
  })
