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
import { getNewConditionSchema, type NewConditionFormData } from "../../patient_schemas"

interface ConditionModalProps {
  isOpen: boolean
  onClose: () => void
  onSubmit: (formData: NewConditionFormData) => void
  isPending: boolean
}

export const ConditionModal = ({
  isOpen,
  onClose,
  onSubmit,
  isPending
}: ConditionModalProps) => {
  const { t } = useTranslation("patients")

  const {
    register,
    handleSubmit,
    formState: { errors },
    reset,
  } = useForm<NewConditionFormData>({
    resolver: zodResolver(getNewConditionSchema(t)),
  })

  if (!isOpen) {
    return null
  }

  const handleFormSubmit = (formData: NewConditionFormData) => {
    onSubmit({
      ...formData,
      icd10Code: formData.icd10Code.toUpperCase().trim(),
    })
    reset()
  }

  return (
    <Dialog open={isOpen} onOpenChange={(open) => !open && onClose()}>
      <DialogContent className="sm:max-w-[460px]">
        <DialogHeader>
          <DialogTitle className="text-left">
            {t("modals.condition.title")}
          </DialogTitle>
        </DialogHeader>
        <form onSubmit={handleSubmit(handleFormSubmit)} className="flex flex-col gap-4 text-left mt-4">
          <div className="flex flex-col gap-1">
            <label className="text-xs font-semibold text-gray-600">
              {t("modals.condition.code")}
            </label>
            <Input
              type="text"
              placeholder={t("modals.condition.codePlaceholder")}
              errorText={errors.icd10Code?.message}
              {...register("icd10Code")}
            />
          </div>

          <div className="flex flex-col gap-1">
            <label className="text-xs font-semibold text-gray-600">
              {t("modals.condition.display")}
            </label>
            <Input
              type="text"
              placeholder={t("modals.condition.displayPlaceholder")}
              errorText={errors.codeDisplay?.message}
              {...register("codeDisplay")}
            />
          </div>

          <div className="flex gap-3 justify-end mt-4">
            <Button variantType="outline" type="button" onClick={onClose}>
              {t("modal.cancel")}
            </Button>
            <Button type="submit" disabled={isPending}>
              {t("modals.condition.confirm")}
            </Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  )
}
