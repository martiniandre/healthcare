import { useForm } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import { useTranslation } from "react-i18next"
import { Input } from "../../../shared/components/ui/Input"
import { Label } from "../../../shared/components/ui/Label"
import { Textarea } from "../../../shared/components/ui/Textarea"
import { Button } from "../../../shared/components/ui/Button"
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "../../../shared/components/ui/Dialog"
import { getUnavailabilitySchema, type UnavailabilityFormData } from "../schedule_schemas"
import type { CreateUnavailabilityPayload } from "../types"

interface CreateUnavailabilityModalProps {
  isOpen: boolean
  onClose: () => void
  onSubmit: (payload: CreateUnavailabilityPayload) => Promise<void>
  isPending: boolean
  staffId: string
  defaultDate: string
}

const formatLocalDateTime = (dateValue: string, timeValue: string): string => {
  return new Date(`${dateValue}T${timeValue}`).toISOString()
}

export const CreateUnavailabilityModal = ({
  isOpen,
  onClose,
  onSubmit,
  isPending,
  staffId,
  defaultDate,
}: CreateUnavailabilityModalProps) => {
  const { t } = useTranslation("schedule")

  const {
    register,
    handleSubmit,
    reset,
    formState: { errors },
  } = useForm<UnavailabilityFormData>({
    resolver: zodResolver(getUnavailabilitySchema(t)),
    defaultValues: {
      date: defaultDate,
      startTime: "",
      endTime: "",
      reason: "",
    },
  })

  const handleClose = () => {
    reset()
    onClose()
  }

  if (!isOpen) {
    return null
  }

  const handleFormSubmit = handleSubmit(async (formData) => {
    const createPayload: CreateUnavailabilityPayload = {
      staff_id: staffId,
      starts_at: formatLocalDateTime(formData.date, formData.startTime),
      ends_at: formatLocalDateTime(formData.date, formData.endTime),
      reason: formData.reason ?? "",
    }

    await onSubmit(createPayload)
    handleClose()
  })

  return (
    <Dialog open={isOpen} onOpenChange={(open) => !open && handleClose()}>
      <DialogContent className="sm:max-w-[480px]">
        <DialogHeader>
          <DialogTitle className="text-left">
            {t("unavailability.modal.title")}
          </DialogTitle>
        </DialogHeader>
        <form onSubmit={handleFormSubmit} className="flex flex-col gap-4 text-left mt-4">
          <div className="grid grid-cols-3 gap-3">
            <div className="flex flex-col gap-1">
              <Label>
                {t("unavailability.modal.date")}
              </Label>
              <Input type="date" errorText={errors.date?.message} {...register("date")} />
            </div>
            <div className="flex flex-col gap-1">
              <Label>
                {t("unavailability.modal.startTime")}
              </Label>
              <Input type="time" step={900} errorText={errors.startTime?.message} {...register("startTime")} />
            </div>
            <div className="flex flex-col gap-1">
              <Label>
                {t("unavailability.modal.endTime")}
              </Label>
              <Input type="time" step={900} errorText={errors.endTime?.message} {...register("endTime")} />
            </div>
          </div>

          <div className="flex flex-col gap-1">
            <Label>
              {t("unavailability.modal.reason")}
            </Label>
            <Textarea
              placeholder={t("unavailability.modal.reasonPlaceholder")}
              {...register("reason")}
            />
            {errors.reason?.message && (
              <span className="text-xs text-danger font-medium px-1 mt-1">
                {errors.reason.message}
              </span>
            )}
          </div>

          <div className="flex gap-3 justify-end mt-4">
            <Button variantType="outline" type="button" onClick={handleClose}>
              {t("unavailability.modal.cancel")}
            </Button>
            <Button type="submit" disabled={isPending}>
              {t("unavailability.modal.confirm")}
            </Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  )
}
