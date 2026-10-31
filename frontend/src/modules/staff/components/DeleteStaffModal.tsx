import { useTranslation } from "react-i18next"
import { UserMinus } from "lucide-react"
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
import type { StaffMember } from "../types"

interface DeleteStaffModalProps {
  member: StaffMember | null
  onClose: () => void
  onConfirm: () => void
  isPending: boolean
}

export const DeleteStaffModal = ({ member, onClose, onConfirm, isPending }: DeleteStaffModalProps) => {
  const { t } = useTranslation("staff")

  return (
    <Dialog open={Boolean(member)} onOpenChange={(isOpen) => !isOpen && onClose()}>
      <DialogContent className="sm:max-w-[420px]">
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2">
            <UserMinus className="size-5 text-destructive" />
            {t("deleteModal.title")}
          </DialogTitle>
          <DialogDescription>
            {t("deleteModal.description")}
            {member && (
              <span className="block mt-2 font-bold text-foreground">{member.fullName}</span>
            )}
          </DialogDescription>
        </DialogHeader>
        <DialogFooter>
          <Button variant="outline" type="button" onClick={onClose}>
            {t("deleteModal.back")}
          </Button>
          <Button variant="destructive" type="button" onClick={onConfirm} disabled={isPending}>
            {isPending && <Spinner />}
            {t("deleteModal.confirm")}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
