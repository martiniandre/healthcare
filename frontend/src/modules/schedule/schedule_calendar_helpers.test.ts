import { describe, expect, it } from "vitest"
import { activeAppointmentsOnly, appointmentsToCalendarEvents, resolveVisibleSlotRange, staffColorForIndex } from "./schedule_calendar_helpers"
import type { Appointment } from "./types"

const buildAppointment = (overrides: Partial<Appointment> = {}): Appointment => ({
  id: "appointment-1",
  patient_fhir_id: "fhir-pat-1",
  staff_id: "staff-1",
  starts_at: "2026-08-20T09:00:00Z",
  ends_at: "2026-08-20T09:30:00Z",
  status: "scheduled",
  reason: "Consulta de rotina",
  version: 1,
  created_at: "2026-08-01T10:00:00Z",
  ...overrides,
})

describe("activeAppointmentsOnly", () => {
  it("removes cancelled appointments from the list", () => {
    const appointments = [
      buildAppointment({ id: "a1", status: "scheduled" }),
      buildAppointment({ id: "a2", status: "cancelled" }),
      buildAppointment({ id: "a3", status: "finished" }),
    ]

    const result = activeAppointmentsOnly(appointments)

    expect(result.map((appointment) => appointment.id)).toEqual(["a1", "a3"])
  })
})

describe("appointmentsToCalendarEvents", () => {
  it("maps appointments into calendar event shapes carrying the staff color", () => {
    const appointment = buildAppointment()
    const staffColor = "#2563eb"

    const events = appointmentsToCalendarEvents([appointment], staffColor)

    expect(events).toHaveLength(1)
    expect(events[0].id).toBe(appointment.id)
    expect(events[0].start).toBe(appointment.starts_at)
    expect(events[0].end).toBe(appointment.ends_at)
    expect(events[0].extendedProps.staffColor).toBe(staffColor)
    expect(events[0].extendedProps.appointment.id).toBe(appointment.id)
  })

  it("skips cancelled appointments when mapping events", () => {
    const events = appointmentsToCalendarEvents(
      [buildAppointment({ status: "cancelled" }), buildAppointment({ id: "a2" })],
      "#2563eb"
    )

    expect(events).toHaveLength(1)
    expect(events[0].id).toBe("a2")
  })

  it("keeps the event title empty when no reason is provided", () => {
    const events = appointmentsToCalendarEvents([buildAppointment({ reason: "" })], "#2563eb")

    expect(events[0].title).toBe("")
  })
})

describe("staffColorForIndex", () => {
  it("returns stable palette colors and wraps around", () => {
    expect(staffColorForIndex(0)).toBe(staffColorForIndex(0))
    expect(staffColorForIndex(0)).not.toBe(staffColorForIndex(1))
    expect(staffColorForIndex(8)).toBe(staffColorForIndex(0))
  })
})

describe("resolveVisibleSlotRange", () => {
  const localIso = (hour: number, minute = 0, dayOffset = 0): string =>
    new Date(2026, 7, 20 + dayOffset, hour, minute).toISOString()

  it("should use the business-day window when there are no events", () => {
    expect(resolveVisibleSlotRange([])).toEqual({ slotMinTime: "06:00:00", slotMaxTime: "22:00:00" })
  })

  it("should keep the default window when events fit inside it", () => {
    expect(resolveVisibleSlotRange([{ start: localIso(9), end: localIso(9, 30) }])).toEqual({
      slotMinTime: "06:00:00",
      slotMaxTime: "22:00:00",
    })
  })

  it("should expand the window so early and late appointments stay visible", () => {
    const range = resolveVisibleSlotRange([
      { start: localIso(4, 30), end: localIso(5) },
      { start: localIso(22, 15), end: localIso(22, 45) },
    ])
    expect(range).toEqual({ slotMinTime: "04:00:00", slotMaxTime: "23:00:00" })
  })

  it("should extend to midnight when an event crosses into the next day", () => {
    const range = resolveVisibleSlotRange([{ start: localIso(23, 30), end: localIso(0, 30, 1) }])
    expect(range.slotMaxTime).toBe("24:00:00")
  })
})
