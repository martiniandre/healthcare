import { useTranslation } from "react-i18next"
import { CalendarX2 } from "lucide-react"
import { Spinner } from "@/shared/components/ui/spinner"
import { Button } from "@/shared/components/ui/button"
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
} from "@/shared/components/ui/dialog"
import type { Appointment } from "../types"
import { formatDateTime, formatTime } from "../../../shared/utils/dates"

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
    <Dialog open={isOpen} onOpenChange={(open) => !open && !isPending && onClose()}>
      <DialogContent className="sm:max-w-[420px]">
        <DialogHeader>
          <DialogTitle className="text-left flex items-center gap-2">
            <CalendarX2 className="size-4 text-destructive" />
            {t("modals.cancel.title")}
          </DialogTitle>
        </DialogHeader>
        <DialogDescription className="text-left text-foreground">
          {t("modals.cancel.description")}
        </DialogDescription>
        <div className="flex flex-col gap-1 rounded-md border bg-muted/50 px-3 py-2 text-left text-sm">
          <span className="font-medium text-foreground">
            {formatDateTime(appointment.starts_at)} – {formatTime(appointment.ends_at)}
          </span>
          <span className="text-xs text-muted-foreground">{appointment.reason || t("cards.noReason")}</span>
        </div>
        <p className="text-left text-xs text-muted-foreground">{t("modals.cancel.hint")}</p>
        <DialogFooter>
          <Button variant="outline" type="button" autoFocus className="h-11 md:h-9" disabled={isPending} onClick={onClose}>
            {t("modals.cancel.back")}
          </Button>
          <Button variant="destructive" type="button" className="h-11 md:h-9" disabled={isPending} aria-busy={isPending} onClick={() => onConfirm(appointment.id)}>
            {isPending && <Spinner />}
            {t("modals.cancel.confirm")}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
