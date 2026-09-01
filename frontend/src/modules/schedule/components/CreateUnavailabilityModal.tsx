import { useForm } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import { useTranslation } from "react-i18next"
import { Input } from "../../../shared/components/ui/Input"
import { Button } from "../../../shared/components/ui/Button"
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "../../../shared/components/ui/Dialog"
import { getUnavailabilitySchema, type UnavailabilityFormData } from "../schedule_schemas"
import { todayDateString } from "../../../shared/utils/validators"
import { getTimeSlotOptions } from "../schedule_time_options"
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

const timeSlotOptions = getTimeSlotOptions()

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
        <form onSubmit={handleFormSubmit} noValidate className="flex flex-col gap-4 text-left mt-4">
          <div className="grid grid-cols-3 gap-3">
            <div className="flex flex-col gap-1">
              <label className="text-xs font-semibold text-gray-600">
                {t("unavailability.modal.date")}
              </label>
              <Input type="date" min={todayDateString()} errorText={errors.date?.message} {...register("date")} />
            </div>
            <div className="flex flex-col gap-1">
              <label className="text-xs font-semibold text-gray-600">
                {t("unavailability.modal.startTime")}
              </label>
              <select
                aria-label={t("unavailability.modal.startTime")}
                className="flex h-9 w-full rounded-md border border-input bg-background px-3 py-1 text-sm shadow-sm focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50"
                {...register("startTime")}
              >
                <option value="">{t("unavailability.modal.selectStartTime")}</option>
                {timeSlotOptions.map((timeSlot) => (
                  <option key={timeSlot.value} value={timeSlot.value}>
                    {timeSlot.label}
                  </option>
                ))}
              </select>
              {errors.startTime?.message && (
                <span className="text-xs text-red-500 font-medium px-1 mt-1">
                  {errors.startTime.message}
                </span>
              )}
            </div>
            <div className="flex flex-col gap-1">
              <label className="text-xs font-semibold text-gray-600">
                {t("unavailability.modal.endTime")}
              </label>
              <select
                aria-label={t("unavailability.modal.endTime")}
                className="flex h-9 w-full rounded-md border border-input bg-background px-3 py-1 text-sm shadow-sm focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50"
                {...register("endTime")}
              >
                <option value="">{t("unavailability.modal.selectEndTime")}</option>
                {timeSlotOptions.map((timeSlot) => (
                  <option key={timeSlot.value} value={timeSlot.value}>
                    {timeSlot.label}
                  </option>
                ))}
              </select>
              {errors.endTime?.message && (
                <span className="text-xs text-red-500 font-medium px-1 mt-1">
                  {errors.endTime.message}
                </span>
              )}
            </div>
          </div>

          <div className="flex flex-col gap-1">
            <label className="text-xs font-semibold text-gray-600">
              {t("unavailability.modal.reason")}
            </label>
            <textarea
              className="flex min-h-[80px] w-full rounded-lg border border-border bg-transparent px-3 py-2 text-sm shadow-sm focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50"
              placeholder={t("unavailability.modal.reasonPlaceholder")}
              {...register("reason")}
            />
            {errors.reason?.message && (
              <span className="text-xs text-red-500 font-medium px-1 mt-1">
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
