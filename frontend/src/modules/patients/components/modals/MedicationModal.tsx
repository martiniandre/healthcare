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
import { getNewMedicationSchema, type NewMedicationFormData } from "../../patient_schemas"

interface MedicationModalProps {
  isOpen: boolean
  onClose: () => void
  onSubmit: (formData: NewMedicationFormData) => void
  isPending: boolean
}

export const MedicationModal = ({
  isOpen,
  onClose,
  onSubmit,
  isPending
}: MedicationModalProps) => {
  const { t } = useTranslation("patients")

  const {
    register,
    handleSubmit,
    formState: { errors },
    reset,
  } = useForm<NewMedicationFormData>({
    resolver: zodResolver(getNewMedicationSchema(t)),
  })

  if (!isOpen) {
    return null
  }

  const handleFormSubmit = (formData: NewMedicationFormData) => {
    onSubmit(formData)
    reset()
  }

  return (
      <Dialog open={isOpen} onOpenChange={(open) => !open && onClose()}>
        <DialogContent className="sm:max-w-[460px]">
          <DialogHeader>
            <DialogTitle className="text-left">
              {t("modals.medication.title")}
            </DialogTitle>
          </DialogHeader>
          <form onSubmit={handleSubmit(handleFormSubmit)} className="flex flex-col gap-4 text-left mt-4">
            <div className="flex flex-col gap-1">
              <label className="text-xs font-semibold text-gray-600">
                {t("modals.medication.name")}
              </label>
              <Input
                type="text"
                placeholder={t("modals.medication.namePlaceholder")}
                errorText={errors.medicationDisplay?.message}
                {...register("medicationDisplay")}
              />
            </div>

            <div className="flex flex-col gap-1">
              <label className="text-xs font-semibold text-gray-600">
                {t("modals.medication.dosage")}
              </label>
              <textarea
                className="w-full h-24 px-3 py-2 text-sm border border-border rounded-lg outline-none focus:ring-2 focus:ring-primary/20 focus:border-primary transition-all resize-none"
                placeholder={t("modals.medication.dosagePlaceholder")}
                {...register("dosageInstruction")}
              />
              {errors.dosageInstruction?.message && (
                <span className="text-xs text-rose-500 font-medium">
                  {errors.dosageInstruction.message}
                </span>
              )}
            </div>

            <div className="flex gap-3 justify-end mt-4">
              <Button variantType="outline" type="button" onClick={onClose}>
                {t("modal.cancel")}
              </Button>
              <Button type="submit" disabled={isPending}>
                {t("modals.medication.confirm")}
              </Button>
            </div>
          </form>
        </DialogContent>
      </Dialog>
  )
}
