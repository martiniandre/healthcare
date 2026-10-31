import { useTranslation } from "react-i18next"
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/shared/components/ui/alert-dialog"

interface DiscardChangesDialogProperties {
  open: boolean
  onKeepEditing: () => void
  onDiscard: () => void
}

function DiscardChangesDialog({ open, onKeepEditing, onDiscard }: DiscardChangesDialogProperties) {
  const { t } = useTranslation("common")

  return (
    <AlertDialog open={open} onOpenChange={(isOpen) => !isOpen && onKeepEditing()}>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>{t("discardChanges.title")}</AlertDialogTitle>
          <AlertDialogDescription>{t("discardChanges.description")}</AlertDialogDescription>
        </AlertDialogHeader>
        <AlertDialogFooter>
          <AlertDialogCancel onClick={onKeepEditing}>{t("discardChanges.keepEditing")}</AlertDialogCancel>
          <AlertDialogAction variant="destructive" onClick={onDiscard}>
            {t("discardChanges.discard")}
          </AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  )
}

export { DiscardChangesDialog }
