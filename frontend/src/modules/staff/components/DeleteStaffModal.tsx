import { useTranslation } from "react-i18next"
import { UserMinus } from "lucide-react"
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
} from "../../../shared/components/ui/Dialog"
import { Button } from "../../../shared/components/ui/Button"
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
            <UserMinus className="w-5 h-5 text-red-500" />
            {t("deleteModal.title")}
          </DialogTitle>
          <DialogDescription>
            {t("deleteModal.description")}
            {member && (
              <span className="block mt-2 font-bold text-gray-900">{member.fullName}</span>
            )}
          </DialogDescription>
        </DialogHeader>
        <div className="flex gap-3 justify-end mt-4">
          <Button variantType="outline" type="button" onClick={onClose}>
            {t("deleteModal.back")}
          </Button>
          <Button variantType="danger" type="button" onClick={onConfirm} disabled={isPending}>
            {t("deleteModal.confirm")}
          </Button>
        </div>
      </DialogContent>
    </Dialog>
  )
}
