import { useState, useEffect, useMemo, useRef, useCallback, type FormEvent } from "react"
import { useForm, Controller, useWatch, type FieldErrors } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import { useTranslation } from "react-i18next"
import { isAxiosError } from "axios"
import {
  Select,
  SelectTrigger,
  SelectValue,
  SelectContent,
  SelectItem,
} from "@/shared/components/ui/select"
import { Alert, AlertDescription, AlertTitle } from "@/shared/components/ui/alert"
import { Input } from "@/shared/components/ui/input"
import { Spinner } from "@/shared/components/ui/spinner"
import { Textarea } from "@/shared/components/ui/textarea"
import { Field, FieldDescription, FieldError, FieldLabel, FieldLegend, FieldSet } from "@/shared/components/ui/field"
import { Button } from "@/shared/components/ui/button"
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
} from "@/shared/components/ui/dialog"
import { DiscardChangesDialog } from "@/shared/components/ui/discard-changes-dialog"
import { useDiscardGuard } from "@/shared/hooks/useDiscardGuard"
import { useLocale } from "@/shared/hooks/useLocale"
import { formatWeekdayDate } from "@/shared/utils/dates"
import {
  getUnavailabilitySchema,
  type UnavailabilityFormData,
  type ScheduleTranslate,
} from "../schedule_schemas"
import { todayDateString } from "../../../shared/utils/validators"
import {
  getTimeSlotOptions,
  getSlotOptionsAfter,
  getDefaultBlockedPeriodEndTime,
  getDurationInMinutes,
} from "../schedule_time_options"
import { FormErrorSummary, OptionalHint, RequiredIndicator, SelectionSummary } from "./ScheduleFormParts"
import type { CreateUnavailabilityPayload } from "../types"

interface CreateUnavailabilityModalProps {
  isOpen: boolean
  onClose: () => void
  onSubmit: (payload: CreateUnavailabilityPayload) => Promise<void>
  isPending: boolean
  staffId: string
  staffName?: string
  defaultDate: string
}

type SubmitFailure = "conflict" | "generic" | null

const unavailabilityFieldOrder = ["date", "startTime", "endTime", "reason"] as const

const reasonMaxLength = 500
const defaultBlockedDurationMinutes = 60
const touchTargetClassName = "h-11 md:h-9"
const selectTouchTargetClassName = "w-full data-[size=default]:h-11 md:data-[size=default]:h-9"
const timeSlotOptions = getTimeSlotOptions()

const formatLocalDateTime = (dateValue: string, timeValue: string): string => {
  return new Date(`${dateValue}T${timeValue}`).toISOString()
}

const joinDescribedBy = (...identifiers: Array<string | false | undefined>): string | undefined => {
  const joinedIdentifiers = identifiers.filter(Boolean).join(" ")
  return joinedIdentifiers || undefined
}

export const CreateUnavailabilityModal = ({
  isOpen,
  onClose,
  onSubmit,
  isPending,
  staffId,
  staffName,
  defaultDate,
}: CreateUnavailabilityModalProps) => {
  const { t } = useTranslation("schedule")
  const locale = useLocale()
  const [submitFailure, setSubmitFailure] = useState<SubmitFailure>(null)
  const [conflictSignature, setConflictSignature] = useState<string | null>(null)
  const [isLocallySubmitting, setIsLocallySubmitting] = useState(false)
  const isSubmittingReference = useRef(false)
  const conflictAlertReference = useRef<HTMLDivElement | null>(null)

  const translateSchedule = useCallback<ScheduleTranslate>(
    (key, options) => t(key, options) as string,
    [t]
  )
  const unavailabilitySchema = useMemo(() => getUnavailabilitySchema(translateSchedule), [translateSchedule])

  const {
    register,
    handleSubmit,
    control,
    reset,
    setValue,
    getValues,
    setFocus,
    trigger,
    formState: { errors, isDirty, submitCount },
  } = useForm<UnavailabilityFormData>({
    resolver: zodResolver(unavailabilitySchema),
    mode: "onTouched",
    shouldFocusError: false,
    defaultValues: {
      date: defaultDate,
      startTime: "",
      endTime: "",
      reason: "",
    },
  })

  const selectedDate = useWatch({ control, name: "date" }) ?? ""
  const selectedStartTime = useWatch({ control, name: "startTime" }) ?? ""
  const selectedEndTime = useWatch({ control, name: "endTime" }) ?? ""
  const reasonValue = useWatch({ control, name: "reason" }) ?? ""
  const endTimeOptions = getSlotOptionsAfter(selectedStartTime)

  const handleClose = useCallback(() => {
    setSubmitFailure(null)
    setConflictSignature(null)
    reset()
    onClose()
  }, [onClose, reset])

  const { isConfirmationOpen, requestClose, keepEditing, discardChanges } = useDiscardGuard({
    isDirty,
    onClose: handleClose,
  })

  const handleDialogOpenChange = (open: boolean) => {
    if (!open && !isPending && !isLocallySubmitting) {
      requestClose()
    }
  }

  useEffect(() => {
    if (isOpen) {
      reset({ date: defaultDate, startTime: "", endTime: "", reason: "" })
    }
  }, [isOpen, defaultDate, reset])

  const currentSignature = `${selectedDate}|${selectedStartTime}|${selectedEndTime}`
  const isConflictVisible = submitFailure === "conflict" && conflictSignature === currentSignature

  useEffect(() => {
    if (isConflictVisible) {
      conflictAlertReference.current?.scrollIntoView?.({ block: "nearest" })
    }
  }, [isConflictVisible])

  const summary = useMemo(() => {
    if (!selectedDate || !selectedStartTime) {
      return null
    }
    const dateLabel = formatWeekdayDate(selectedDate, locale)
    const durationMinutes = getDurationInMinutes(selectedStartTime, selectedEndTime)
    const whenLabel = selectedEndTime && durationMinutes > 0
      ? t("form.summary.when", {
          date: dateLabel,
          start: selectedStartTime,
          end: selectedEndTime,
          duration: t("form.summary.minutes", { count: durationMinutes }),
        })
      : t("form.summary.whenOpenEnded", { date: dateLabel, start: selectedStartTime })
    return staffName ? t("unavailability.summary.forStaff", { when: whenLabel, staff: staffName }) : whenLabel
  }, [selectedDate, selectedStartTime, selectedEndTime, staffName, locale, t])

  const errorSummaryItems = submitCount > 0
    ? [
        { name: "date", label: t("unavailability.modal.date"), message: errors.date?.message },
        { name: "startTime", label: t("unavailability.modal.startTime"), message: errors.startTime?.message },
        { name: "endTime", label: t("unavailability.modal.endTime"), message: errors.endTime?.message },
        { name: "reason", label: t("unavailability.modal.reason"), message: errors.reason?.message },
      ].filter((item): item is { name: string; label: string; message: string } => Boolean(item.message))
    : []

  if (!isOpen) {
    return null
  }

  const handleStartTimeChange = (nextStartTime: string, notifyChange: (value: string) => void) => {
    const previousStartTime = getValues("startTime")
    const previousEndTime = getValues("endTime")
    notifyChange(nextStartTime)
    const previousDuration = getDurationInMinutes(previousStartTime, previousEndTime)
    const nextEndTime = getDefaultBlockedPeriodEndTime(
      nextStartTime,
      previousDuration > 0 ? previousDuration : defaultBlockedDurationMinutes
    )
    const shouldRevalidateEnd = Boolean(errors.endTime) || submitCount > 0
    if (nextEndTime !== previousEndTime) {
      setValue("endTime", nextEndTime, { shouldDirty: true, shouldValidate: shouldRevalidateEnd })
    } else if (shouldRevalidateEnd) {
      void trigger("endTime")
    }
  }

  const focusFirstInvalidField = (invalidFields: FieldErrors<UnavailabilityFormData>) => {
    const firstInvalidName = unavailabilityFieldOrder.find((fieldName) => invalidFields[fieldName])
    if (firstInvalidName) {
      setFocus(firstInvalidName)
    }
  }

  const submitValidForm = async (formData: UnavailabilityFormData) => {
    if (isSubmittingReference.current || isPending) {
      return
    }
    isSubmittingReference.current = true
    setIsLocallySubmitting(true)
    setSubmitFailure(null)
    const createPayload: CreateUnavailabilityPayload = {
      staff_id: staffId,
      starts_at: formatLocalDateTime(formData.date, formData.startTime),
      ends_at: formatLocalDateTime(formData.date, formData.endTime),
      reason: formData.reason ?? "",
    }

    try {
      await onSubmit(createPayload)
      handleClose()
    } catch (submitError) {
      if (isAxiosError(submitError) && submitError.response?.status === 409) {
        setConflictSignature(`${formData.date}|${formData.startTime}|${formData.endTime}`)
        setSubmitFailure("conflict")
      } else {
        setSubmitFailure("generic")
      }
    } finally {
      isSubmittingReference.current = false
      setIsLocallySubmitting(false)
    }
  }

  const handleFormSubmit = (event: FormEvent<HTMLFormElement>) =>
    handleSubmit(submitValidForm, focusFirstInvalidField)(event)

  const isBusy = isPending || isLocallySubmitting

  return (
    <>
      <Dialog open={isOpen} onOpenChange={handleDialogOpenChange}>
        <DialogContent className="flex max-h-[calc(100dvh-2rem)] flex-col gap-0 p-0 sm:max-w-[560px]">
          <DialogHeader className="border-b px-6 py-4 pr-12 text-left">
            <DialogTitle className="text-left">{t("unavailability.modal.title")}</DialogTitle>
            <DialogDescription>
              {staffName
                ? t("unavailability.modal.descriptionForStaff", { staff: staffName })
                : t("unavailability.modal.description")}
            </DialogDescription>
          </DialogHeader>
          <form onSubmit={handleFormSubmit} noValidate className="flex min-h-0 flex-1 flex-col text-left">
            <div className="flex flex-1 flex-col gap-5 overflow-y-auto px-6 py-4">
              <p className="text-xs text-muted-foreground">
                <span aria-hidden="true" className="text-destructive">*</span> {t("form.requiredNote")}
              </p>

              <FormErrorSummary
                items={errorSummaryItems}
                onSelectItem={(fieldName) => setFocus(fieldName as keyof UnavailabilityFormData)}
              />

              {submitFailure === "generic" && (
                <Alert variant="destructive">
                  <AlertTitle>{t("unavailability.errors.submitFailedTitle")}</AlertTitle>
                  <AlertDescription>{t("unavailability.errors.submitFailed")}</AlertDescription>
                </Alert>
              )}

              <FieldSet className="gap-3">
                <FieldLegend variant="label" className="mb-0 text-xs font-semibold uppercase tracking-wide text-muted-foreground">
                  {t("form.sections.when")}
                </FieldLegend>
                <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
                  <Field data-invalid={!!errors.date} className="md:col-span-2">
                    <FieldLabel htmlFor="unavailability-date">
                      {t("unavailability.modal.date")}
                      <RequiredIndicator />
                    </FieldLabel>
                    <Input
                      id="unavailability-date"
                      type="date"
                      autoFocus
                      min={todayDateString()}
                      className={touchTargetClassName}
                      aria-required="true"
                      aria-invalid={!!errors.date}
                      aria-describedby={joinDescribedBy(!!errors.date && "unavailability-date-error")}
                      {...register("date")}
                    />
                    <FieldError id="unavailability-date-error">{errors.date?.message}</FieldError>
                  </Field>
                  <Field data-invalid={!!errors.startTime}>
                    <FieldLabel htmlFor="unavailability-start-time">
                      {t("unavailability.modal.startTime")}
                      <RequiredIndicator />
                    </FieldLabel>
                    <Controller
                      control={control}
                      name="startTime"
                      render={({ field }) => (
                        <Select
                          onValueChange={(nextStartTime) => handleStartTimeChange(nextStartTime, field.onChange)}
                          value={field.value}
                        >
                          <SelectTrigger
                            id="unavailability-start-time"
                            ref={field.ref}
                            onBlur={field.onBlur}
                            className={selectTouchTargetClassName}
                            aria-required="true"
                            aria-invalid={!!errors.startTime}
                            aria-describedby={joinDescribedBy(!!errors.startTime && "unavailability-start-error")}
                          >
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
                    <FieldError id="unavailability-start-error">{errors.startTime?.message}</FieldError>
                  </Field>
                  <Field data-invalid={!!errors.endTime}>
                    <FieldLabel htmlFor="unavailability-end-time">
                      {t("unavailability.modal.endTime")}
                      <RequiredIndicator />
                    </FieldLabel>
                    <Controller
                      control={control}
                      name="endTime"
                      render={({ field }) => (
                        <Select onValueChange={field.onChange} value={field.value}>
                          <SelectTrigger
                            id="unavailability-end-time"
                            ref={field.ref}
                            onBlur={field.onBlur}
                            className={selectTouchTargetClassName}
                            aria-required="true"
                            aria-invalid={!!errors.endTime}
                            aria-describedby={joinDescribedBy(
                              "unavailability-end-hint",
                              !!errors.endTime && "unavailability-end-error"
                            )}
                          >
                            <SelectValue placeholder={t("unavailability.modal.selectEndTime")} />
                          </SelectTrigger>
                          <SelectContent>
                            {endTimeOptions.map((timeSlot) => (
                              <SelectItem key={timeSlot.value} value={timeSlot.value}>
                                {timeSlot.label}
                              </SelectItem>
                            ))}
                          </SelectContent>
                        </Select>
                      )}
                    />
                    <FieldDescription id="unavailability-end-hint" className="text-xs">
                      {t("unavailability.modal.endTimeHint")}
                    </FieldDescription>
                    <FieldError id="unavailability-end-error">{errors.endTime?.message}</FieldError>
                  </Field>
                </div>

                {isConflictVisible && (
                  <Alert variant="destructive" ref={conflictAlertReference} data-testid="conflict-alert">
                    <AlertTitle>{t("unavailability.errors.conflictTitle")}</AlertTitle>
                    <AlertDescription>{t("unavailability.errors.conflict")}</AlertDescription>
                  </Alert>
                )}

                {summary ? (
                  <SelectionSummary id="unavailability-summary" text={summary} />
                ) : (
                  <p
                    id="unavailability-summary"
                    role="status"
                    aria-live="polite"
                    className="text-xs text-muted-foreground"
                  >
                    {t("form.summary.empty")}
                  </p>
                )}
              </FieldSet>

              <FieldSet className="gap-3">
                <FieldLegend variant="label" className="mb-0 text-xs font-semibold uppercase tracking-wide text-muted-foreground">
                  {t("form.sections.why")}
                </FieldLegend>
                <Field data-invalid={!!errors.reason}>
                  <FieldLabel htmlFor="unavailability-reason">
                    {t("unavailability.modal.reason")}
                    <OptionalHint />
                  </FieldLabel>
                  <Textarea
                    id="unavailability-reason"
                    rows={3}
                    maxLength={reasonMaxLength}
                    placeholder={t("unavailability.modal.reasonPlaceholder")}
                    aria-invalid={!!errors.reason}
                    aria-describedby={joinDescribedBy(
                      "unavailability-reason-counter",
                      !!errors.reason && "unavailability-reason-error"
                    )}
                    {...register("reason")}
                  />
                  <div className="flex items-start justify-between gap-3">
                    <FieldError id="unavailability-reason-error">{errors.reason?.message}</FieldError>
                    <span
                      id="unavailability-reason-counter"
                      className="ml-auto text-xs text-muted-foreground"
                    >
                      {t("form.characterCount", { count: reasonValue.length, max: reasonMaxLength })}
                    </span>
                  </div>
                </Field>
              </FieldSet>
            </div>

            <DialogFooter className="border-t bg-background px-6 py-4">
              <Button
                variant="outline"
                type="button"
                className={touchTargetClassName}
                onClick={() => handleDialogOpenChange(false)}
                disabled={isBusy}
              >
                {t("unavailability.modal.cancel")}
              </Button>
              <Button type="submit" className={touchTargetClassName} disabled={isBusy} aria-busy={isBusy}>
                {isBusy && <Spinner />}
                {t("unavailability.modal.confirm")}
              </Button>
            </DialogFooter>
          </form>
        </DialogContent>
      </Dialog>
      <DiscardChangesDialog open={isConfirmationOpen} onKeepEditing={keepEditing} onDiscard={discardChanges} />
    </>
  )
}
