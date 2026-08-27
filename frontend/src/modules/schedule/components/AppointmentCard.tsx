import { useTranslation } from "react-i18next"
import { CalendarX2 } from "lucide-react"
import { usePatientQuery } from "../../patients/queries"
import { formatTime } from "../../../shared/utils/dates"
import { Badge } from "../../../shared/components/ui/Badge"
import { Button } from "../../../shared/components/ui/Button"
import type { Appointment } from "../types"

interface AppointmentCardProps {
  appointment: Appointment
  onCancel: (appointment: Appointment) => void
}

const statusBadgeVariants: Record<string, "info" | "success" | "muted" | "secondary"> = {
  scheduled: "info",
  confirmed: "success",
  cancelled: "muted",
  finished: "secondary",
}

export const AppointmentCard = ({ appointment, onCancel }: AppointmentCardProps) => {
  const { t } = useTranslation("schedule")
  const { data: patient } = usePatientQuery(appointment.patient_fhir_id)

  const canCancel = appointment.status === "scheduled" || appointment.status === "confirmed"

  return (
    <div className={`bg-surface border border-border rounded-xl p-4 flex flex-col gap-2 ${appointment.status === "cancelled" ? "opacity-60" : ""}`}>
      <div className="flex items-center justify-between gap-2">
        <span className="font-mono text-xs font-bold text-foreground">
          {formatTime(appointment.starts_at)} — {formatTime(appointment.ends_at)}
        </span>
        <Badge variant={statusBadgeVariants[appointment.status] ?? "muted"} className="capitalize">
          {t(`status.${appointment.status}`)}
        </Badge>
      </div>
      <div className="flex flex-col gap-0.5">
        <span className="text-sm font-bold text-foreground">
          {patient?.full_name ?? appointment.patient_fhir_id}
        </span>
        <span className="text-xs text-muted-foreground truncate">
          {appointment.reason || t("cards.noReason")}
        </span>
      </div>
      {canCancel && (
        <Button
          type="button"
          variantType="danger"
          size="sm"
          onClick={() => onCancel(appointment)}
          className="mt-1 self-start"
        >
          <CalendarX2 className="w-3.5 h-3.5" />
          {t("cards.cancelAppointment")}
        </Button>
      )}
    </div>
  )
}
