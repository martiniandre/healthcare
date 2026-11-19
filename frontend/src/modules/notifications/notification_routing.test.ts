import { describe, expect, it } from "vitest"
import { resolveNotificationTarget, sortNotificationsByUrgency } from "./notification_routing"
import type { NotificationItem } from "./types"

const buildNotification = (overrides: Partial<NotificationItem> = {}): NotificationItem => ({
  id: "notification-1",
  type: "system",
  priority: "medium",
  title: "Aviso",
  body: "Corpo",
  resource_type: "",
  resource_id: "",
  is_read: false,
  created_at: "2026-11-16T10:00:00Z",
  ...overrides,
})

describe("sortNotificationsByUrgency", () => {
  it("should list unread before read and critical before low", () => {
    const sorted = sortNotificationsByUrgency([
      buildNotification({ id: "low-unread", priority: "low", created_at: "2026-11-16T12:00:00Z" }),
      buildNotification({ id: "critical-read", priority: "critical", is_read: true }),
      buildNotification({ id: "critical-unread", priority: "critical", created_at: "2026-11-16T08:00:00Z" }),
    ])

    expect(sorted.map((notification) => notification.id)).toEqual(["critical-unread", "low-unread", "critical-read"])
  })

  it("should order notifications with the same urgency from newest to oldest", () => {
    const sorted = sortNotificationsByUrgency([
      buildNotification({ id: "older", created_at: "2026-11-16T08:00:00Z" }),
      buildNotification({ id: "newer", created_at: "2026-11-16T09:00:00Z" }),
    ])

    expect(sorted.map((notification) => notification.id)).toEqual(["newer", "older"])
  })

  it("should not mutate the original list", () => {
    const original = [buildNotification({ id: "first", priority: "low" }), buildNotification({ id: "second", priority: "critical" })]
    sortNotificationsByUrgency(original)
    expect(original[0].id).toBe("first")
  })
})

describe("resolveNotificationTarget", () => {
  it("should send staff to the encounter of the patient", () => {
    const notification = buildNotification({ type: "encounter_created", resource_id: "enc-9", params: { patient_id: "fhir-pat-1" } })
    expect(resolveNotificationTarget(notification, true)).toBe("/patients/fhir-pat-1/encounters/enc-9")
  })

  it("should send staff to the patient record for a ready report", () => {
    const notification = buildNotification({ type: "report_ready", resource_id: "rep-1", params: { patient_id: "fhir-pat-1" } })
    expect(resolveNotificationTarget(notification, true)).toBe("/patients/fhir-pat-1")
  })

  it("should send staff to telemetry for bed alerts", () => {
    expect(resolveNotificationTarget(buildNotification({ type: "telemetry_alert" }), true)).toBe("/telemetry")
  })

  it("should send patients to the reports section of the portal", () => {
    expect(resolveNotificationTarget(buildNotification({ type: "report_ready" }), false)).toBe("/portal?tab=reports")
  })

  it("should return null when there is nothing to open", () => {
    expect(resolveNotificationTarget(buildNotification({ type: "system" }), true)).toBeNull()
    expect(resolveNotificationTarget(buildNotification({ type: "report_ready" }), true)).toBeNull()
  })
})
