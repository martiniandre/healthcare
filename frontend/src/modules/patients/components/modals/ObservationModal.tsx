import { useForm, Controller } from "react-hook-form"
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
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "../../../../shared/components/ui/Select"
import { getNewObservationSchema, type NewObservationFormData } from "../../patient_schemas"

interface ObservationModalProps {
  isOpen: boolean
  onClose: () => void
  onSubmit: (formData: NewObservationFormData) => void
  isPending: boolean
}

export const ObservationModal = ({
  isOpen,
  onClose,
  onSubmit,
  isPending
}: ObservationModalProps) => {
  const { t } = useTranslation("patients")

  const {
    register,
    handleSubmit,
    control,
    formState: { errors },
  } = useForm<NewObservationFormData>({
    resolver: zodResolver(getNewObservationSchema(t)),
    defaultValues: {
      loincCode: "8867-4",
    },
  })

  if (!isOpen) {
    return null
  }

  return (
      <Dialog open={isOpen} onOpenChange={(open) => !open && onClose()}>
        <DialogContent className="sm:max-w-[460px]">
          <DialogHeader>
            <DialogTitle className="text-left">
              {t("modals.observation.title")}
            </DialogTitle>
          </DialogHeader>
          <form onSubmit={handleSubmit(onSubmit)} className="flex flex-col gap-4 text-left mt-4">
            <div className="flex flex-col gap-1">
              <label className="text-xs font-semibold text-gray-600">
                {t("modals.observation.selectMetric")}
              </label>
              <Controller
                control={control}
                name="loincCode"
                render={({ field }) => (
                  <Select onValueChange={field.onChange} defaultValue={field.value}>
                    <SelectTrigger className="w-full">
                      <SelectValue placeholder={t("modals.observation.selectMetric")} />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="8867-4">{t("modals.observation.heartRate")}</SelectItem>
                      <SelectItem value="8310-5">{t("modals.observation.temperature")}</SelectItem>
                      <SelectItem value="85354-9">{t("modals.observation.bloodPressure")}</SelectItem>
                    </SelectContent>
                  </Select>
                )}
              />
            </div>

            <div className="flex flex-col gap-1">
              <label className="text-xs font-semibold text-gray-600">
                {t("modals.observation.value")}
              </label>
              <Input
                type="number"
                step="any"
                placeholder={t("modals.observation.valuePlaceholder")}
                errorText={errors.valueQuantity?.message}
                {...register("valueQuantity", { valueAsNumber: true })}
              />
            </div>

            <div className="flex gap-3 justify-end mt-4">
              <Button variantType="outline" type="button" onClick={onClose}>
                {t("modal.cancel")}
              </Button>
              <Button type="submit" disabled={isPending}>
                {t("modals.observation.confirm")}
              </Button>
            </div>
          </form>
        </DialogContent>
      </Dialog>
  )
}
