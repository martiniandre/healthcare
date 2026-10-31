import { useTranslation } from "react-i18next"
import { CalendarX2 } from "lucide-react"
import { Spinner } from "@/shared/components/ui/spinner"
import { Button } from "@/shared/components/ui/button"
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from "@/shared/components/ui/dialog"
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
            <CalendarX2 className="size-4 text-destructive" />
            {t("modals.cancel.title")}
          </DialogTitle>
        </DialogHeader>
        <div className="flex flex-col gap-2 text-left mt-4 text-sm text-foreground">
          <p>{t("modals.cancel.description")}</p>
          <p className="font-mono text-xs text-muted-foreground">
            {formatDateTime(appointment.starts_at)} — {appointment.reason || "-"}
          </p>
        </div>
        <DialogFooter className="mt-6">
          <Button variant="outline" type="button" onClick={onClose}>
            {t("modals.cancel.back")}
          </Button>
          <Button variant="destructive" type="button" disabled={isPending} onClick={() => onConfirm(appointment.id)}>
            {isPending && <Spinner />}
            {t("modals.cancel.confirm")}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
