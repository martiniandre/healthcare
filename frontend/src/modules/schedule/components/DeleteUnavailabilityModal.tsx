import { useTranslation } from "react-i18next"
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
} from "../../../shared/components/ui/Dialog"
import { Button } from "../../../shared/components/ui/Button"

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
    <Dialog open={isOpen} onOpenChange={(open) => !open && onClose()}>
      <DialogContent className="sm:max-w-[420px]">
        <DialogHeader>
          <DialogTitle>{t("unavailability.deleteModal.title")}</DialogTitle>
          <DialogDescription>{t("unavailability.deleteModal.description")}</DialogDescription>
        </DialogHeader>
        <div className="flex gap-3 justify-end mt-4">
          <Button variantType="outline" type="button" onClick={onClose}>
            {t("unavailability.deleteModal.back")}
          </Button>
          <Button variantType="destructive" type="button" onClick={onConfirm} disabled={isPending}>
            {t("unavailability.deleteModal.confirm")}
          </Button>
        </div>
      </DialogContent>
    </Dialog>
  )
}
