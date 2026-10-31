import { useForm, Controller } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import { useTranslation } from "react-i18next"
import {
  Select,
  SelectTrigger,
  SelectValue,
  SelectContent,
  SelectItem,
} from "@/shared/components/ui/select"
import { Input } from "@/shared/components/ui/input"
import { Spinner } from "@/shared/components/ui/spinner"
import { Textarea } from "@/shared/components/ui/textarea"
import { Field, FieldError, FieldLabel } from "@/shared/components/ui/field"
import { Button } from "@/shared/components/ui/button"
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from "@/shared/components/ui/dialog"
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
    control,
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
            <Field data-invalid={!!errors.date}>
              <FieldLabel htmlFor="unavailability-date">
                {t("unavailability.modal.date")}
              </FieldLabel>
              <Input id="unavailability-date" type="date" min={todayDateString()} aria-invalid={!!errors.date} {...register("date")} />
              <FieldError>{errors.date?.message}</FieldError>
            </Field>
            <Field data-invalid={!!errors.startTime}>
              <FieldLabel htmlFor="unavailability-start-time">
                {t("unavailability.modal.startTime")}
              </FieldLabel>
              <Controller
                control={control}
                name="startTime"
                render={({ field }) => (
                  <Select onValueChange={field.onChange} value={field.value}>
                    <SelectTrigger id="unavailability-start-time" className="w-full" aria-invalid={!!errors.startTime}>
                      <SelectValue placeholder={t("unavailability.modal.selectStartTime")} />
                    </SelectTrigger>
                    <SelectContent>
                      {timeSlotOptions.map((timeSlot) => (
                        <SelectItem key={timeSlot.value} value={timeSlot.value}>
                          {timeSlot.label}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                )}
              />
              <FieldError>{errors.startTime?.message}</FieldError>
            </Field>
            <Field data-invalid={!!errors.endTime}>
              <FieldLabel htmlFor="unavailability-end-time">
                {t("unavailability.modal.endTime")}
              </FieldLabel>
              <Controller
                control={control}
                name="endTime"
                render={({ field }) => (
                  <Select onValueChange={field.onChange} value={field.value}>
                    <SelectTrigger id="unavailability-end-time" className="w-full" aria-invalid={!!errors.endTime}>
                      <SelectValue placeholder={t("unavailability.modal.selectEndTime")} />
                    </SelectTrigger>
                    <SelectContent>
                      {timeSlotOptions.map((timeSlot) => (
                        <SelectItem key={timeSlot.value} value={timeSlot.value}>
                          {timeSlot.label}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                )}
              />
              <FieldError>{errors.endTime?.message}</FieldError>
            </Field>
          </div>

          <Field data-invalid={!!errors.reason}>
            <FieldLabel htmlFor="unavailability-reason">
              {t("unavailability.modal.reason")}
            </FieldLabel>
            <Textarea
              id="unavailability-reason"
              placeholder={t("unavailability.modal.reasonPlaceholder")}
              aria-invalid={!!errors.reason}
              {...register("reason")}
            />
            <FieldError>{errors.reason?.message}</FieldError>
          </Field>

          <DialogFooter className="mt-4">
            <Button variant="outline" type="button" onClick={handleClose}>
              {t("unavailability.modal.cancel")}
            </Button>
            <Button type="submit" disabled={isPending}>
              {isPending && <Spinner />}
              {t("unavailability.modal.confirm")}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}
