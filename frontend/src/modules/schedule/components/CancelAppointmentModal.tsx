import { useTranslation } from "react-i18next"
import { CalendarX2 } from "lucide-react"
import { Button } from "../../../shared/components/ui/Button"
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "../../../shared/components/ui/Dialog"
import type { Appointment } from "../types"
import { formatDateTime } from "../../../shared/utils/dates"

interface CancelAppointmentModalProps {
  isOpen: boolean
  appointment: Appointment | null
  onClose: () => void
  onConfirm: (appointmentId: string) => void
  isPending: boolean
}

export const CancelAppointmentModal = ({
  isOpen,
  appointment,
  onClose,
  onConfirm,
  isPending,
}: CancelAppointmentModalProps) => {
  const { t } = useTranslation("schedule")

  if (!isOpen || !appointment) {
    return null
  }

  return (
    <Dialog open={isOpen} onOpenChange={(open) => !open && onClose()}>
      <DialogContent className="sm:max-w-[420px]">
        <DialogHeader>
          <DialogTitle className="text-left flex items-center gap-2">
            <CalendarX2 className="w-4 h-4 text-red-500" />
            {t("modals.cancel.title")}
          </DialogTitle>
        </DialogHeader>
        <div className="flex flex-col gap-2 text-left mt-4 text-sm text-gray-700">
          <p>{t("modals.cancel.description")}</p>
          <p className="font-mono text-xs text-gray-500">
            {formatDateTime(appointment.starts_at)} — {appointment.reason || "-"}
          </p>
        </div>
        <div className="flex gap-3 justify-end mt-6">
          <Button variantType="outline" type="button" onClick={onClose}>
            {t("modals.cancel.back")}
          </Button>
          <Button variantType="danger" type="button" disabled={isPending} onClick={() => onConfirm(appointment.id)}>
            {t("modals.cancel.confirm")}
          </Button>
        </div>
      </DialogContent>
    </Dialog>
  )
}
