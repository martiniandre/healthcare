import { useTranslation } from "react-i18next"
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
} from "@/shared/components/ui/dialog"
import { Spinner } from "@/shared/components/ui/spinner"
import { Button } from "@/shared/components/ui/button"

interface DeleteUnavailabilityModalProps {
  isOpen: boolean
  onClose: () => void
  onConfirm: () => void
  isPending: boolean
}

export const DeleteUnavailabilityModal = ({
  isOpen,
  onClose,
  onConfirm,
  isPending,
}: DeleteUnavailabilityModalProps) => {
  const { t } = useTranslation("schedule")

  if (!isOpen) {
    return null
  }

  return (
    <Dialog open={isOpen} onOpenChange={(open) => !open && !isPending && onClose()}>
      <DialogContent className="sm:max-w-[420px]">
        <DialogHeader>
          <DialogTitle>{t("unavailability.deleteModal.title")}</DialogTitle>
          <DialogDescription>{t("unavailability.deleteModal.description")}</DialogDescription>
        </DialogHeader>
        <DialogFooter>
          <Button variant="outline" type="button" autoFocus className="h-11 md:h-9" disabled={isPending} onClick={onClose}>
            {t("unavailability.deleteModal.back")}
          </Button>
          <Button variant="destructive" type="button" className="h-11 md:h-9" onClick={onConfirm} disabled={isPending} aria-busy={isPending}>
            {isPending && <Spinner />}
            {t("unavailability.deleteModal.confirm")}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
