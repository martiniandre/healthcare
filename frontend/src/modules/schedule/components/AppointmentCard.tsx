import { useTranslation } from "react-i18next"
import { CalendarX2 } from "lucide-react"
import { StatusBadge } from "@/shared/components/ui/status-badge"
import { Button } from "@/shared/components/ui/button"
import { Card, CardContent } from "@/shared/components/ui/card"
import { usePatientQuery } from "../../patients/queries"
import { formatTime } from "../../../shared/utils/dates"
import { resolveAppointmentClinicalStatus } from "../appointment_status"
import type { Appointment } from "../types"

interface AppointmentCardProps {
  appointment: Appointment
  onCancel: (appointment: Appointment) => void
}

export const AppointmentCard = ({ appointment, onCancel }: AppointmentCardProps) => {
  const { t } = useTranslation("schedule")
  const { data: patient } = usePatientQuery(appointment.patient_fhir_id)

  const canCancel = appointment.status === "scheduled" || appointment.status === "confirmed"

  return (
    <Card className={`py-4 gap-2 ${appointment.status === "cancelled" ? "opacity-60" : ""}`}>
      <CardContent className="flex flex-col gap-2 px-4">
        <div className="flex items-center justify-between gap-2">
          <span className="font-mono text-xs font-bold tabular-nums text-foreground">
            {formatTime(appointment.starts_at)} — {formatTime(appointment.ends_at)}
          </span>
          <StatusBadge status={resolveAppointmentClinicalStatus(appointment.status)}>
            {t(`status.${appointment.status}`)}
          </StatusBadge>
        </div>
        <div className="flex flex-col gap-0.5">
          <span className="text-sm font-bold text-foreground">
            {patient?.full_name ?? t("cards.unknownPatient")}
          </span>
          <span className="text-xs text-muted-foreground truncate">
            {appointment.reason || t("cards.noReason")}
          </span>
        </div>
        {canCancel && (
          <Button
            type="button"
            variant="ghost"
            size="sm"
            onClick={() => onCancel(appointment)}
            className="self-start text-xs font-semibold text-destructive hover:text-destructive"
          >
            <CalendarX2 />
            {t("cards.cancelAppointment")}
          </Button>
        )}
      </CardContent>
    </Card>
  )
}
