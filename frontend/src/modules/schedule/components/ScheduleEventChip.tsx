import { useTranslation } from "react-i18next"
import { usePatientQuery } from "../../patients/queries"
import type { Appointment } from "../types"

interface ScheduleEventChipProps {
  appointment: Appointment
  staffColor: string
}

const formatChipTime = (dateTimeValue: string): string => {
  return new Date(dateTimeValue).toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" })
}

const tintedBackground = (hexColor: string, alpha: number): string => {
  const normalizedHex = hexColor.replace("#", "")
  const red = parseInt(normalizedHex.slice(0, 2), 16)
  const green = parseInt(normalizedHex.slice(2, 4), 16)
  const blue = parseInt(normalizedHex.slice(4, 6), 16)
  return `rgba(${red}, ${green}, ${blue}, ${alpha})`
}

export const ScheduleEventChip = ({ appointment, staffColor }: ScheduleEventChipProps) => {
  const { t } = useTranslation("schedule")
  const { data: patient } = usePatientQuery(appointment.patient_fhir_id)
  const patientLabel = patient?.full_name ?? t("cards.unknownPatient")

  return (
    <div
      className="schedule-event-ticket flex items-center gap-1.5 rounded-md px-1.5 py-0.5 text-[11px] leading-tight"
      style={{
        backgroundColor: tintedBackground(staffColor, 0.12),
        borderColor: tintedBackground(staffColor, 0.28),
      }}
    >
      <span className="shrink-0 w-1 self-stretch rounded-full" style={{ backgroundColor: staffColor }} />
      <span className="shrink-0 tabular-nums font-semibold opacity-80" style={{ color: staffColor }}>
        {formatChipTime(appointment.starts_at)}
      </span>
      <span className="truncate font-semibold text-gray-800">{patientLabel}</span>
    </div>
  )
}
