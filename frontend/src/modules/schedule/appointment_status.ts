import type { ClinicalStatus } from "@/shared/components/ui/status-badge"
import type { Appointment } from "./types"

const APPOINTMENT_STATUS_PRESENTATION: Record<Appointment["status"], ClinicalStatus> = {
  scheduled: "info",
  confirmed: "normal",
  cancelled: "critical",
  finished: "neutral",
}

export const resolveAppointmentClinicalStatus = (status: Appointment["status"]): ClinicalStatus => {
  return APPOINTMENT_STATUS_PRESENTATION[status] ?? "neutral"
}
