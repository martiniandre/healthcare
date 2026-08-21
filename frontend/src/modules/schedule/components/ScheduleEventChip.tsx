import { usePatientQuery } from "../../patients/queries"
import type { Appointment } from "../types"

interface ScheduleEventChipProps {
  appointment: Appointment
  staffColor: string
}

const formatChipTime = (dateTimeValue: string): string => {
  return new Date(dateTimeValue).toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" })
}

export const ScheduleEventChip = ({ appointment, staffColor }: ScheduleEventChipProps) => {
  const { data: patient } = usePatientQuery(appointment.patient_fhir_id)
  const patientLabel = patient?.full_name ?? appointment.patient_fhir_id

  return (
    <div
      className="flex items-center gap-1 overflow-hidden rounded px-1 py-0.5 text-[11px] font-semibold leading-tight text-white"
      style={{ backgroundColor: staffColor }}
    >
      <span className="shrink-0 tabular-nums opacity-90">{formatChipTime(appointment.starts_at)}</span>
      <span className="truncate">{patientLabel}</span>
    </div>
  )
}
