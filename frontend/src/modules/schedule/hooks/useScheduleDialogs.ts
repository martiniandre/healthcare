import { useCallback, useState } from "react"
import type { Appointment, StaffUnavailability } from "../types"

export interface CreateAppointmentDefaults {
  date: string
  startTime: string
}

export type ScheduleDialog =
  | { kind: "none" }
  | { kind: "create"; defaults: CreateAppointmentDefaults | null }
  | { kind: "view"; appointment: Appointment }
  | { kind: "edit"; appointment: Appointment }
  | { kind: "cancel"; appointment: Appointment }
  | { kind: "createUnavailability" }
  | { kind: "deleteUnavailability"; unavailability: StaffUnavailability }

const CLOSED_DIALOG: ScheduleDialog = { kind: "none" }

export const appointmentForDialog = (dialog: ScheduleDialog, kind: "view" | "edit" | "cancel"): Appointment | null =>
  dialog.kind === kind ? dialog.appointment : null

export const useScheduleDialogs = () => {
  const [activeDialog, setActiveDialog] = useState<ScheduleDialog>(CLOSED_DIALOG)
  const [lastCreateDefaults, setLastCreateDefaults] = useState<CreateAppointmentDefaults | null>(null)

  const closeDialog = useCallback(() => setActiveDialog(CLOSED_DIALOG), [])

  const openCreate = useCallback((defaults: CreateAppointmentDefaults | null = null) => {
    setLastCreateDefaults(defaults)
    setActiveDialog({ kind: "create", defaults })
  }, [])

  const openView = useCallback((appointment: Appointment) => setActiveDialog({ kind: "view", appointment }), [])
  const openEdit = useCallback((appointment: Appointment) => setActiveDialog({ kind: "edit", appointment }), [])
  const openCancel = useCallback((appointment: Appointment) => setActiveDialog({ kind: "cancel", appointment }), [])
  const openCreateUnavailability = useCallback(() => setActiveDialog({ kind: "createUnavailability" }), [])
  const openDeleteUnavailability = useCallback(
    (unavailability: StaffUnavailability) => setActiveDialog({ kind: "deleteUnavailability", unavailability }),
    []
  )

  return {
    activeDialog,
    createDefaults: lastCreateDefaults,
    closeDialog,
    openCreate,
    openView,
    openEdit,
    openCancel,
    openCreateUnavailability,
    openDeleteUnavailability,
  }
}
