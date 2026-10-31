import { describe, it, expect } from "vitest"
import { AxiosError, AxiosHeaders, type AxiosResponse } from "axios"
import {
  extractAllergyBlockedAlerts,
  hasBlockingAlert,
  mergeAlerts,
  resolveAlertBadgeStatus,
} from "./medicationSafety"
import type { MedicationSafetyAlert } from "./types"

const buildAlert = (overrides: Partial<MedicationSafetyAlert> = {}): MedicationSafetyAlert => ({
  kind: "interaction",
  severity: "moderate",
  blocking: false,
  messageKey: "interaction.moderate",
  params: {},
  relatedResourceId: "resource-1",
  ...overrides,
})

const buildAxiosError = (status: number, data: unknown): AxiosError => {
  const response = { status, data, statusText: "", headers: {}, config: { headers: new AxiosHeaders() } } as AxiosResponse
  return new AxiosError("failed", "ERR_BAD_REQUEST", undefined, undefined, response)
}

describe("medicationSafety helpers", () => {
  it("should map severities to badge statuses", () => {
    expect(resolveAlertBadgeStatus(buildAlert({ severity: "moderate" }))).toBe("attention")
    expect(resolveAlertBadgeStatus(buildAlert({ severity: "major" }))).toBe("critical")
    expect(resolveAlertBadgeStatus(buildAlert({ severity: "contraindicated" }))).toBe("critical")
    expect(resolveAlertBadgeStatus(buildAlert({ kind: "allergy", severity: "info" }))).toBe("critical")
  })

  it("should detect blocking alerts", () => {
    expect(hasBlockingAlert([buildAlert()])).toBe(false)
    expect(hasBlockingAlert([buildAlert(), buildAlert({ blocking: true, relatedResourceId: "resource-2" })])).toBe(true)
  })

  it("should merge alerts without duplicating the same alert", () => {
    expect(mergeAlerts([buildAlert()], [buildAlert(), buildAlert({ relatedResourceId: "resource-2" })])).toHaveLength(2)
  })

  it("should extract the alerts of an allergy blocked 409", () => {
    const alerts = [buildAlert({ kind: "allergy", blocking: true })]
    const error = buildAxiosError(409, { error: "blocked", code: "MEDICATION_ALLERGY_BLOCKED", alerts })
    expect(extractAllergyBlockedAlerts(error)).toEqual(alerts)
  })

  it("should ignore other errors", () => {
    expect(extractAllergyBlockedAlerts(buildAxiosError(409, { code: "OTHER" }))).toBeNull()
    expect(extractAllergyBlockedAlerts(buildAxiosError(500, { code: "MEDICATION_ALLERGY_BLOCKED" }))).toBeNull()
    expect(extractAllergyBlockedAlerts(new Error("boom"))).toBeNull()
  })
})
