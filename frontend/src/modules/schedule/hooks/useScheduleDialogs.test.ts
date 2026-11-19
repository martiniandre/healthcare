import { describe, expect, it } from "vitest"
import { act, renderHook } from "@testing-library/react"
import { appointmentForDialog, useScheduleDialogs } from "./useScheduleDialogs"
import type { Appointment, StaffUnavailability } from "../types"

const appointment = {
  id: "appointment-1",
  patient_fhir_id: "fhir-pat-1",
  staff_id: "staff-1",
  starts_at: "2026-11-20T09:00:00Z",
  ends_at: "2026-11-20T09:30:00Z",
  status: "scheduled",
  reason: "Rotina",
  version: 1,
  created_at: "2026-11-01T10:00:00Z",
} as Appointment

const unavailability = { id: "block-1" } as StaffUnavailability

describe("useScheduleDialogs", () => {
  it("should start with every dialog closed", () => {
    const { result } = renderHook(() => useScheduleDialogs())
    expect(result.current.activeDialog.kind).toBe("none")
    expect(result.current.createDefaults).toBeNull()
  })

  it("should keep only one dialog open when switching from details to edit", () => {
    const { result } = renderHook(() => useScheduleDialogs())
    act(() => result.current.openView(appointment))
    expect(appointmentForDialog(result.current.activeDialog, "view")).toBe(appointment)

    act(() => result.current.openEdit(appointment))
    expect(appointmentForDialog(result.current.activeDialog, "view")).toBeNull()
    expect(appointmentForDialog(result.current.activeDialog, "edit")).toBe(appointment)
  })

  it("should remember the slot chosen for a new appointment after closing", () => {
    const { result } = renderHook(() => useScheduleDialogs())
    act(() => result.current.openCreate({ date: "2026-11-20", startTime: "09:00" }))
    expect(result.current.activeDialog.kind).toBe("create")

    act(() => result.current.closeDialog())
    expect(result.current.activeDialog.kind).toBe("none")
    expect(result.current.createDefaults).toEqual({ date: "2026-11-20", startTime: "09:00" })
  })

  it("should open the unavailability dialogs with their payload", () => {
    const { result } = renderHook(() => useScheduleDialogs())
    act(() => result.current.openDeleteUnavailability(unavailability))
    expect(result.current.activeDialog).toEqual({ kind: "deleteUnavailability", unavailability })

    act(() => result.current.openCreateUnavailability())
    expect(result.current.activeDialog).toEqual({ kind: "createUnavailability" })
  })
})
