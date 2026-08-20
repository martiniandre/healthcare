import { useTranslation } from "react-i18next"
import { CalendarX2 } from "lucide-react"
import { usePatientQuery } from "../../patients/queries"
import { formatTime } from "../../../shared/utils/dates"
import type { Appointment } from "../types"

interface AppointmentCardProps {
  appointment: Appointment
  onCancel: (appointment: Appointment) => void
}

const statusBadgeClassNames: Record<string, string> = {
  scheduled: "bg-blue-50 text-blue-700 border-blue-200",
  confirmed: "bg-emerald-50 text-emerald-700 border-emerald-200",
  cancelled: "bg-gray-100 text-gray-500 border-gray-200",
  finished: "bg-violet-50 text-violet-700 border-violet-200",
}

export const AppointmentCard = ({ appointment, onCancel }: AppointmentCardProps) => {
  const { t } = useTranslation("schedule")
  const { data: patient } = usePatientQuery(appointment.patient_fhir_id)

  const canCancel = appointment.status === "scheduled" || appointment.status === "confirmed"

  return (
    <div className={`bg-white border border-border rounded-xl p-4 flex flex-col gap-2 ${appointment.status === "cancelled" ? "opacity-60" : ""}`}>
      <div className="flex items-center justify-between gap-2">
        <span className="font-mono text-xs font-bold text-gray-800">
          {formatTime(appointment.starts_at)} — {formatTime(appointment.ends_at)}
        </span>
        <span className={`text-[10px] font-bold px-2 py-0.5 rounded-full border capitalize ${statusBadgeClassNames[appointment.status] ?? "bg-gray-100 text-gray-500 border-gray-200"}`}>
          {t(`status.${appointment.status}`)}
        </span>
      </div>
      <div className="flex flex-col gap-0.5">
        <span className="text-sm font-bold text-gray-900">
          {patient?.full_name ?? appointment.patient_fhir_id}
        </span>
        <span className="text-xs text-gray-500 truncate">
          {appointment.reason || t("cards.noReason")}
        </span>
      </div>
      {canCancel && (
        <button
          onClick={() => onCancel(appointment)}
          className="mt-1 self-start inline-flex items-center gap-1.5 text-xs font-semibold text-red-500 hover:text-red-700"
        >
          <CalendarX2 className="w-3.5 h-3.5" />
          {t("cards.cancelAppointment")}
        </button>
      )}
    </div>
  )
}
