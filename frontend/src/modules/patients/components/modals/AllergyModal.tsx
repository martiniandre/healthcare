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
import { getNewAllergySchema, type NewAllergyFormData } from "../../patient_schemas"

interface AllergyModalProps {
  isOpen: boolean
  onClose: () => void
  onSubmit: (formData: NewAllergyFormData) => void
  isPending: boolean
}

export const AllergyModal = ({
  isOpen,
  onClose,
  onSubmit,
  isPending
}: AllergyModalProps) => {
  const { t } = useTranslation("patients")

  const {
    register,
    handleSubmit,
    formState: { errors },
    reset,
  } = useForm<NewAllergyFormData>({
    resolver: zodResolver(getNewAllergySchema(t)),
  })

  if (!isOpen) {
    return null
  }

  const handleFormSubmit = (formData: NewAllergyFormData) => {
    onSubmit(formData)
    reset()
  }

  return (
    <Dialog open={isOpen} onOpenChange={(open) => !open && onClose()}>
      <DialogContent className="sm:max-w-[460px]">
        <DialogHeader>
          <DialogTitle className="text-left">
            {t("modals.allergy.title")}
          </DialogTitle>
        </DialogHeader>
        <form onSubmit={handleSubmit(handleFormSubmit)} className="flex flex-col gap-4 text-left mt-4">
          <div className="flex flex-col gap-1">
            <label className="text-xs font-semibold text-gray-600">
              {t("modals.allergy.code")}
            </label>
            <Input
              type="text"
              placeholder={t("modals.allergy.codePlaceholder")}
              errorText={errors.allergenCode?.message}
              {...register("allergenCode")}
            />
          </div>

          <div className="flex flex-col gap-1">
            <label className="text-xs font-semibold text-gray-600">
              {t("modals.allergy.display")}
            </label>
            <Input
              type="text"
              placeholder={t("modals.allergy.displayPlaceholder")}
              errorText={errors.allergenDisplay?.message}
              {...register("allergenDisplay")}
            />
          </div>

          <div className="flex flex-col gap-1">
            <label className="text-xs font-semibold text-gray-600">
              {t("modals.allergy.reaction")}
            </label>
            <Input
              type="text"
              placeholder={t("modals.allergy.reactionPlaceholder")}
              errorText={errors.reaction?.message}
              {...register("reaction")}
            />
          </div>

          <div className="flex gap-3 justify-end mt-4">
            <Button variantType="outline" type="button" onClick={onClose}>
              {t("modal.cancel")}
            </Button>
            <Button type="submit" disabled={isPending}>
              {t("modals.allergy.confirm")}
            </Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  )
}
