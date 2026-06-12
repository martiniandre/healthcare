import { useForm } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import { useTranslation } from "react-i18next"
import { Input } from "../../../../shared/components/ui/Input"
import { Button } from "../../../../shared/components/ui/Button"
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "../../../../shared/components/ui/Dialog"
import { getNewEncounterSchema, type NewEncounterFormData } from "../../patient_schemas"

interface EncounterModalProps {
  isOpen: boolean
  onClose: () => void
  onSubmit: (formData: NewEncounterFormData) => void
  isPending: boolean
}

export const EncounterModal = ({
  isOpen,
  onClose,
  onSubmit,
  isPending
}: EncounterModalProps) => {
  const { t } = useTranslation("patients")

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<NewEncounterFormData>({
    resolver: zodResolver(getNewEncounterSchema(t)),
  })

  if (!isOpen) {
    return null
  }

  return (
    <Dialog open={isOpen} onOpenChange={(open) => !open && onClose()}>
      <DialogContent className="sm:max-w-[460px]">
        <DialogHeader>
          <DialogTitle className="text-left">
            {t("modals.encounter.title")}
          </DialogTitle>
        </DialogHeader>
        <form onSubmit={handleSubmit(onSubmit)} className="flex flex-col gap-4 text-left mt-4">
          <div className="flex flex-col gap-1">
            <label className="text-xs font-semibold text-gray-600">
              {t("modals.encounter.reason")}
            </label>
            <Input
              type="text"
              placeholder={t("modals.encounter.reasonPlaceholder")}
              errorText={errors.reasonDisplay?.message}
              {...register("reasonDisplay")}
            />
          </div>
          <div className="flex gap-3 justify-end mt-4">
            <Button variantType="outline" type="button" onClick={onClose}>
              {t("modal.cancel")}
            </Button>
            <Button type="submit" disabled={isPending}>
              {t("modals.encounter.confirm")}
            </Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  )
}
