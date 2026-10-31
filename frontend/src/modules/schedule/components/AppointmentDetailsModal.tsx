import { useTranslation } from "react-i18next"
import { CalendarClock } from "lucide-react"
import { StatusBadge } from "@/shared/components/ui/status-badge"
import { resolveAppointmentClinicalStatus } from "../appointment_status"
import { Button } from "@/shared/components/ui/button"
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from "@/shared/components/ui/dialog"
import { useStaffListQuery } from "../../staff/queries"
import { usePatientQuery } from "../../patients/queries"
import { formatDateTime } from "../../../shared/utils/dates"
import type { Appointment } from "../types"

interface AppointmentDetailsModalProps {
  isOpen: boolean
  appointment: Appointment | null
  onClose: () => void
  onEdit: (appointment: Appointment) => void
  onCancel: (appointment: Appointment) => void
}

const getStatusLabelKey = (status: Appointment["status"]): string => {
  return `status.${status}`
}

export const AppointmentDetailsModal = ({
  isOpen,
  appointment,
  onClose,
  onEdit,
  onCancel,
}: AppointmentDetailsModalProps) => {
  const { t } = useTranslation("schedule")
  const { data: staffMembers = [] } = useStaffListQuery()
  const { data: patient } = usePatientQuery(appointment?.patient_fhir_id ?? "")

  if (!isOpen || !appointment) {
    return null
  }

  const assignedStaffMember = staffMembers.find((staffMember) => staffMember.id === appointment.staff_id)

  return (
    <Dialog open={isOpen} onOpenChange={(open) => !open && onClose()}>
      <DialogContent className="sm:max-w-[440px]">
        <DialogHeader>
          <DialogTitle className="text-left flex items-center gap-2">
            <CalendarClock className="size-4 text-primary" />
            {t("modals.details.title")}
          </DialogTitle>
        </DialogHeader>
        <div className="flex flex-col gap-4 text-left mt-4">
          <div className="flex flex-col gap-1 text-sm">
            <span className="text-xs font-semibold text-muted-foreground">{t("modals.details.patient")}</span>
            <span className="font-semibold text-foreground">
              {patient?.full_name ?? t("cards.unknownPatient")}
            </span>
          </div>
          <div className="flex flex-col gap-1 text-sm">
            <span className="text-xs font-semibold text-muted-foreground">{t("modals.details.staff")}</span>
            <span className="text-foreground">
              {assignedStaffMember
                ? `${assignedStaffMember.fullName} — ${assignedStaffMember.role}`
                : "-"}
            </span>
          </div>
          <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <div className="flex flex-col gap-1 text-sm">
              <span className="text-xs font-semibold text-muted-foreground">{t("modals.details.date")}</span>
              <span className="text-foreground">{formatDateTime(appointment.starts_at)}</span>
            </div>
            <div className="flex flex-col gap-1 text-sm">
              <span className="text-xs font-semibold text-muted-foreground">{t("modals.details.end")}</span>
              <span className="text-foreground">{formatDateTime(appointment.ends_at)}</span>
            </div>
          </div>
          <div className="flex flex-col gap-1 text-sm">
            <span className="text-xs font-semibold text-muted-foreground">{t("modals.details.reason")}</span>
            <span className="text-foreground">{appointment.reason || t("cards.noReason")}</span>
          </div>
          <div className="flex flex-col gap-1 text-sm">
            <span className="text-xs font-semibold text-muted-foreground">{t("modals.details.status")}</span>
            <StatusBadge status={resolveAppointmentClinicalStatus(appointment.status)}>
              {t(getStatusLabelKey(appointment.status))}
            </StatusBadge>
          </div>
        </div>
        <DialogFooter className="mt-6">
          <Button variant="outline" type="button" onClick={() => onEdit(appointment)}>
            {t("modals.details.edit")}
          </Button>
          <Button variant="outline" type="button" onClick={() => onCancel(appointment)}>
            {t("modals.details.cancel")}
          </Button>
          <Button type="button" onClick={onClose}>
            {t("modals.details.close")}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}