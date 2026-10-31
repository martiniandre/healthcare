import { useState, useEffect, useMemo, useRef, useCallback, type FormEvent } from "react"
import { useForm, Controller, useWatch, type FieldErrors } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import { useTranslation } from "react-i18next"
import { isAxiosError } from "axios"
import { Alert, AlertDescription, AlertTitle } from "@/shared/components/ui/alert"
import { Input } from "@/shared/components/ui/input"
import {
  Select,
  SelectTrigger,
  SelectValue,
  SelectContent,
  SelectItem,
} from "@/shared/components/ui/select"
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
  getNewAppointmentSchema,
  type NewAppointmentFormData,
  type ScheduleTranslate,
} from "../schedule_schemas"
import { todayDateString } from "../../../shared/utils/validators"
import {
  getAvailableStartTimeOptions,
  getEndTimeOptionsForStart,
  getDefaultAppointmentEndTime,
  getAppointmentEndTimeForNewStart,
  getDurationInMinutes,
} from "../schedule_time_options"
import { useStaffListQuery } from "../../staff/queries"
import { usePatientsQuery, usePatientQuery } from "../../patients/queries"
import { useIdempotencyKey } from "../hooks/useIdempotencyKey"
import { FormErrorSummary, OptionalHint, RequiredIndicator, SelectionSummary } from "./ScheduleFormParts"
import type { Appointment, CreateAppointmentPayload, UpdateAppointmentPayload } from "../types"

interface AppointmentModalProps {
  isOpen: boolean
  onClose: () => void
  onSubmit: (payload: CreateAppointmentPayload | UpdateAppointmentPayload) => Promise<void>
  isPending: boolean
  defaultStaffId?: string
  defaultDate?: string
  defaultStartTime?: string
  appointment?: Appointment | null
}

type SubmitFailure = "conflict" | "generic" | null

const appointmentFieldOrder = ["patientFhirId", "staffId", "date", "startTime", "endTime", "reason"] as const

const reasonMaxLength = 500
const touchTargetClassName = "h-11 md:h-9"
const selectTouchTargetClassName = "w-full data-[size=default]:h-11 md:data-[size=default]:h-9"

const formatLocalDateTime = (dateValue: string, timeValue: string): string => {
  return new Date(`${dateValue}T${timeValue}`).toISOString()
}

const getLocalDateValue = (isoDateTime: string): string => {
  const dateObject = new Date(isoDateTime)
  const year = dateObject.getFullYear()
  const month = String(dateObject.getMonth() + 1).padStart(2, "0")
  const day = String(dateObject.getDate()).padStart(2, "0")
  return `${year}-${month}-${day}`
}

const getLocalTimeValue = (isoDateTime: string): string => {
  const dateObject = new Date(isoDateTime)
  const hour = String(dateObject.getHours()).padStart(2, "0")
  const minute = String(dateObject.getMinutes()).padStart(2, "0")
  return `${hour}:${minute}`
}

const joinDescribedBy = (...identifiers: Array<string | false | undefined>): string | undefined => {
  const joinedIdentifiers = identifiers.filter(Boolean).join(" ")
  return joinedIdentifiers || undefined
}

export const AppointmentModal = ({
  isOpen,
  onClose,
  onSubmit,
  isPending,
  defaultStaffId,
  defaultDate,
  defaultStartTime,
  appointment,
}: AppointmentModalProps) => {
  const { t } = useTranslation("schedule")
  const locale = useLocale()
  const isEditMode = appointment !== undefined && appointment !== null
  const { data: staffMembers = [] } = useStaffListQuery()
  const { data: patientsPage } = usePatientsQuery("", "", "", 1, 100)
  const { data: appointmentPatient } = usePatientQuery(
    isEditMode && appointment ? appointment.patient_fhir_id : ""
  )
  const patients = useMemo(() => patientsPage?.patients ?? [], [patientsPage])
  const { getOrCreateKey, resetKey } = useIdempotencyKey()

  const [patientSearch, setPatientSearch] = useState("")
  const [submitFailure, setSubmitFailure] = useState<SubmitFailure>(null)
  const [conflictSignature, setConflictSignature] = useState<string | null>(null)
  const [isLocallySubmitting, setIsLocallySubmitting] = useState(false)
  const isSubmittingReference = useRef(false)
  const conflictAlertReference = useRef<HTMLDivElement | null>(null)

  const translateSchedule = useCallback<ScheduleTranslate>(
    (key, options) => t(key, options) as string,
    [t]
  )
  const appointmentSchema = useMemo(() => getNewAppointmentSchema(translateSchedule), [translateSchedule])

  const patientCandidates = useMemo(() => {
    if (!isEditMode || !appointmentPatient) {
      return patients
    }
    const alreadyListed = patients.some(
      (patient) => patient.fhir_resource_id === appointmentPatient.fhir_resource_id
    )
    return alreadyListed ? patients : [...patients, appointmentPatient]
  }, [isEditMode, appointmentPatient, patients])

  const filteredPatients = patientCandidates.filter((patient) =>
    patient.full_name.toLowerCase().includes(patientSearch.trim().toLowerCase())
  )

  const {
    register,
    handleSubmit,
    control,
    reset,
    setValue,
    getValues,
    trigger,
    setFocus,
    formState: { errors, isDirty, submitCount },
  } = useForm<NewAppointmentFormData>({
    resolver: zodResolver(appointmentSchema),
    mode: "onTouched",
    shouldFocusError: false,
    defaultValues: {
      patientFhirId: "",
      staffId: defaultStaffId ?? "",
      date: defaultDate ?? "",
      startTime: "",
      endTime: "",
      reason: "",
    },
  })

  const selectedPatientId = useWatch({ control, name: "patientFhirId" }) ?? ""
  const selectedStaffId = useWatch({ control, name: "staffId" }) ?? ""
  const selectedDate = useWatch({ control, name: "date" }) ?? ""
  const selectedStartTime = useWatch({ control, name: "startTime" }) ?? ""
  const selectedEndTime = useWatch({ control, name: "endTime" }) ?? ""
  const reasonValue = useWatch({ control, name: "reason" }) ?? ""
  const availableStartTimeOptions = getAvailableStartTimeOptions(selectedDate)
  const endTimeOptions = getEndTimeOptionsForStart(selectedStartTime)

  const handleClose = useCallback(() => {
    setSubmitFailure(null)
    setConflictSignature(null)
    setPatientSearch("")
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
      if (isEditMode && appointment) {
        reset({
          patientFhirId: appointment.patient_fhir_id,
          staffId: appointment.staff_id,
          date: getLocalDateValue(appointment.starts_at),
          startTime: getLocalTimeValue(appointment.starts_at),
          endTime: getLocalTimeValue(appointment.ends_at),
          reason: appointment.reason || "",
        })
      } else {
        reset({
          patientFhirId: "",
          staffId: defaultStaffId ?? "",
          date: defaultDate ?? "",
          startTime: defaultStartTime ?? "",
          endTime: defaultStartTime ? getDefaultAppointmentEndTime(defaultStartTime) : "",
          reason: "",
        })
      }
    }
  }, [isOpen, appointment, isEditMode, defaultStaffId, defaultDate, defaultStartTime, reset])

  const currentSignature = `${selectedStaffId}|${selectedDate}|${selectedStartTime}|${selectedEndTime}`
  const isConflictVisible = submitFailure === "conflict" && conflictSignature === currentSignature

  useEffect(() => {
    if (isConflictVisible) {
      conflictAlertReference.current?.scrollIntoView?.({ block: "nearest" })
    }
  }, [isConflictVisible])

  const selectedPatientName =
    patientCandidates.find((patient) => patient.fhir_resource_id === selectedPatientId)?.full_name ?? ""

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
    const staffName = staffMembers.find((staffMember) => staffMember.id === selectedStaffId)?.fullName ?? ""
    return {
      text: staffName ? t("form.summary.withStaff", { when: whenLabel, staff: staffName }) : whenLabel,
      detail: selectedPatientName ? t("form.summary.forPatient", { patient: selectedPatientName }) : "",
    }
  }, [selectedDate, selectedStartTime, selectedEndTime, selectedStaffId, selectedPatientName, staffMembers, locale, t])

  const errorSummaryItems = submitCount > 0
    ? [
        { name: "patientFhirId", label: t("modals.create.patient"), message: errors.patientFhirId?.message },
        { name: "staffId", label: t("modals.create.staff"), message: errors.staffId?.message },
        { name: "date", label: t("modals.create.date"), message: errors.date?.message },
        { name: "startTime", label: t("modals.create.startTime"), message: errors.startTime?.message },
        { name: "endTime", label: t("modals.create.endTime"), message: errors.endTime?.message },
        { name: "reason", label: t("modals.create.reason"), message: errors.reason?.message },
      ].filter((item): item is { name: string; label: string; message: string } => Boolean(item.message))
    : []

  if (!isOpen) {
    return null
  }

  const handleStartTimeChange = (nextStartTime: string, notifyChange: (value: string) => void) => {
    const previousStartTime = getValues("startTime")
    const previousEndTime = getValues("endTime")
    notifyChange(nextStartTime)
    const nextEndTime = getAppointmentEndTimeForNewStart(previousStartTime, previousEndTime, nextStartTime)
    const shouldRevalidateEnd = Boolean(errors.endTime) || submitCount > 0
    if (nextEndTime !== previousEndTime) {
      setValue("endTime", nextEndTime, { shouldDirty: true, shouldValidate: shouldRevalidateEnd })
    } else if (shouldRevalidateEnd) {
      void trigger("endTime")
    }
  }

  const focusFirstInvalidField = (invalidFields: FieldErrors<NewAppointmentFormData>) => {
    const firstInvalidName = appointmentFieldOrder.find((fieldName) => invalidFields[fieldName])
    if (firstInvalidName) {
      setFocus(firstInvalidName)
    }
  }

  const submitValidForm = async (formData: NewAppointmentFormData) => {
    if (isSubmittingReference.current || isPending) {
      return
    }
    isSubmittingReference.current = true
    setIsLocallySubmitting(true)
    setSubmitFailure(null)
    const startsAt = formatLocalDateTime(formData.date, formData.startTime)
    const endsAt = formatLocalDateTime(formData.date, formData.endTime)

    try {
      if (isEditMode && appointment) {
        const updatePayload: UpdateAppointmentPayload = {
          patient_fhir_id: formData.patientFhirId,
          staff_id: formData.staffId,
          starts_at: startsAt,
          ends_at: endsAt,
          reason: formData.reason,
        }
        await onSubmit(updatePayload)
      } else {
        const createPayload: CreateAppointmentPayload = {
          patient_fhir_id: formData.patientFhirId,
          staff_id: formData.staffId,
          starts_at: startsAt,
          ends_at: endsAt,
          reason: formData.reason,
          idempotency_key: getOrCreateKey(),
        }
        await onSubmit(createPayload)
        resetKey()
      }
      handleClose()
    } catch (submitError) {
      if (isAxiosError(submitError) && submitError.response?.status === 409) {
        setConflictSignature(`${formData.staffId}|${formData.date}|${formData.startTime}|${formData.endTime}`)
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
  const hasNoPatientMatch = patientSearch.trim() !== "" && filteredPatients.length === 0

  return (
    <>
      <Dialog open={isOpen} onOpenChange={handleDialogOpenChange}>
        <DialogContent className="flex max-h-[calc(100dvh-2rem)] flex-col gap-0 p-0 sm:max-w-[640px]">
          <DialogHeader className="border-b px-6 py-4 pr-12 text-left">
            <DialogTitle className="text-left">
              {isEditMode ? t("modals.edit.title") : t("modals.create.title")}
            </DialogTitle>
            <DialogDescription>
              {isEditMode ? t("modals.edit.description") : t("modals.create.description")}
            </DialogDescription>
          </DialogHeader>
          <form onSubmit={handleFormSubmit} noValidate className="flex min-h-0 flex-1 flex-col text-left">
            <div className="flex flex-1 flex-col gap-5 overflow-y-auto px-6 py-4">
              <p className="text-xs text-muted-foreground">
                <span aria-hidden="true" className="text-destructive">*</span> {t("form.requiredNote")}
              </p>

              <FormErrorSummary
                items={errorSummaryItems}
                onSelectItem={(fieldName) => setFocus(fieldName as keyof NewAppointmentFormData)}
              />

              {submitFailure === "generic" && (
                <Alert variant="destructive">
                  <AlertTitle>{t("errors.submitFailedTitle")}</AlertTitle>
                  <AlertDescription>{t("errors.submitFailed")}</AlertDescription>
                </Alert>
              )}

              <FieldSet className="gap-3">
                <FieldLegend variant="label" className="mb-0 text-xs font-semibold uppercase tracking-wide text-muted-foreground">
                  {t("form.sections.who")}
                </FieldLegend>
                <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
                  <Field data-invalid={!!errors.patientFhirId}>
                    <FieldLabel htmlFor="appointment-patient-search">
                      {t("modals.create.patient")}
                      <RequiredIndicator />
                    </FieldLabel>
                    <Input
                      id="appointment-patient-search"
                      type="search"
                      autoFocus
                      autoComplete="off"
                      value={patientSearch}
                      onChange={(inputEvent) => setPatientSearch(inputEvent.target.value)}
                      placeholder={t("modals.create.patientSearchPlaceholder")}
                      aria-describedby="appointment-patient-hint"
                      className={touchTargetClassName}
                    />
                    <Controller
                      control={control}
                      name="patientFhirId"
                      render={({ field }) => (
                        <Select onValueChange={field.onChange} value={field.value}>
                          <SelectTrigger
                            id="appointment-patient"
                            ref={field.ref}
                            onBlur={field.onBlur}
                            className={selectTouchTargetClassName}
                            aria-label={t("modals.create.selectPatient")}
                            aria-required="true"
                            aria-invalid={!!errors.patientFhirId}
                            aria-describedby={joinDescribedBy(
                              "appointment-patient-hint",
                              !!errors.patientFhirId && "appointment-patient-error"
                            )}
                          >
                            <SelectValue placeholder={t("modals.create.selectPatient")} />
                          </SelectTrigger>
                          <SelectContent>
                            {filteredPatients.map((patient) => (
                              <SelectItem key={patient.fhir_resource_id} value={patient.fhir_resource_id}>
                                {patient.full_name}
                              </SelectItem>
                            ))}
                          </SelectContent>
                        </Select>
                      )}
                    />
                    <FieldDescription id="appointment-patient-hint" aria-live="polite" className="text-xs">
                      {hasNoPatientMatch ? t("modals.create.patientNoMatch") : t("modals.create.patientHint")}
                    </FieldDescription>
                    <FieldError id="appointment-patient-error">{errors.patientFhirId?.message}</FieldError>
                  </Field>

                  <Field data-invalid={!!errors.staffId}>
                    <FieldLabel htmlFor="appointment-staff">
                      {t("modals.create.staff")}
                      <RequiredIndicator />
                    </FieldLabel>
                    <Controller
                      control={control}
                      name="staffId"
                      render={({ field }) => (
                        <Select onValueChange={field.onChange} value={field.value}>
                          <SelectTrigger
                            id="appointment-staff"
                            ref={field.ref}
                            onBlur={field.onBlur}
                            className={selectTouchTargetClassName}
                            aria-required="true"
                            aria-invalid={!!errors.staffId}
                            aria-describedby={joinDescribedBy(!!errors.staffId && "appointment-staff-error")}
                          >
                            <SelectValue placeholder={t("modals.create.selectStaff")} />
                          </SelectTrigger>
                          <SelectContent>
                            {staffMembers.filter((staffMember) => staffMember.id).map((staffMember) => (
                              <SelectItem key={staffMember.id} value={staffMember.id}>
                                {staffMember.fullName} — {staffMember.role}
                              </SelectItem>
                            ))}
                          </SelectContent>
                        </Select>
                      )}
                    />
                    <FieldError id="appointment-staff-error">{errors.staffId?.message}</FieldError>
                  </Field>
                </div>
              </FieldSet>

              <FieldSet className="gap-3">
                <FieldLegend variant="label" className="mb-0 text-xs font-semibold uppercase tracking-wide text-muted-foreground">
                  {t("form.sections.when")}
                </FieldLegend>
                <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
                  <Field data-invalid={!!errors.date} className="md:col-span-2">
                    <FieldLabel htmlFor="appointment-date">
                      {t("modals.create.date")}
                      <RequiredIndicator />
                    </FieldLabel>
                    <Input
                      id="appointment-date"
                      type="date"
                      min={todayDateString()}
                      className={touchTargetClassName}
                      aria-required="true"
                      aria-invalid={!!errors.date}
                      aria-describedby={joinDescribedBy(!!errors.date && "appointment-date-error")}
                      {...register("date", {
                        onChange: () => {
                          if (getValues("startTime")) {
                            void trigger("startTime")
                          }
                        },
                      })}
                    />
                    <FieldError id="appointment-date-error">{errors.date?.message}</FieldError>
                  </Field>
                  <Field data-invalid={!!errors.startTime}>
                    <FieldLabel htmlFor="appointment-start-time">
                      {t("modals.create.startTime")}
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
                            id="appointment-start-time"
                            ref={field.ref}
                            onBlur={field.onBlur}
                            className={selectTouchTargetClassName}
                            aria-required="true"
                            aria-invalid={!!errors.startTime}
                            aria-describedby={joinDescribedBy(
                              "appointment-start-hint",
                              !!errors.startTime && "appointment-start-error"
                            )}
                          >
                            <SelectValue placeholder={t("modals.create.selectStartTime")} />
                          </SelectTrigger>
                          <SelectContent>
                            {availableStartTimeOptions.map((timeSlot) => (
                              <SelectItem key={timeSlot.value} value={timeSlot.value}>
                                {timeSlot.label}
                              </SelectItem>
                            ))}
                          </SelectContent>
                        </Select>
                      )}
                    />
                    <FieldDescription id="appointment-start-hint" className="text-xs">
                      {t("modals.create.startTimeHint")}
                    </FieldDescription>
                    <FieldError id="appointment-start-error">{errors.startTime?.message}</FieldError>
                  </Field>
                  <Field data-invalid={!!errors.endTime}>
                    <FieldLabel htmlFor="appointment-end-time">
                      {t("modals.create.endTime")}
                      <RequiredIndicator />
                    </FieldLabel>
                    <Controller
                      control={control}
                      name="endTime"
                      render={({ field }) => (
                        <Select onValueChange={field.onChange} value={field.value}>
                          <SelectTrigger
                            id="appointment-end-time"
                            ref={field.ref}
                            onBlur={field.onBlur}
                            className={selectTouchTargetClassName}
                            aria-required="true"
                            aria-invalid={!!errors.endTime}
                            aria-describedby={joinDescribedBy(
                              "appointment-end-hint",
                              !!errors.endTime && "appointment-end-error"
                            )}
                          >
                            <SelectValue placeholder={t("modals.create.selectEndTime")} />
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
                    <FieldDescription id="appointment-end-hint" className="text-xs">
                      {selectedStartTime ? t("modals.create.endTimeHint") : t("modals.create.endTimeHintNoStart")}
                    </FieldDescription>
                    <FieldError id="appointment-end-error">{errors.endTime?.message}</FieldError>
                  </Field>
                </div>

                {isConflictVisible && (
                  <Alert variant="destructive" ref={conflictAlertReference} data-testid="conflict-alert">
                    <AlertTitle>{t("errors.conflictTitle")}</AlertTitle>
                    <AlertDescription>{t("errors.conflict")}</AlertDescription>
                  </Alert>
                )}

                {summary ? (
                  <SelectionSummary id="appointment-summary" text={summary.text} detail={summary.detail} />
                ) : (
                  <p
                    id="appointment-summary"
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
                  <FieldLabel htmlFor="appointment-reason">
                    {t("modals.create.reason")}
                    <OptionalHint />
                  </FieldLabel>
                  <Textarea
                    id="appointment-reason"
                    rows={3}
                    maxLength={reasonMaxLength}
                    placeholder={t("modals.create.reasonPlaceholder")}
                    aria-invalid={!!errors.reason}
                    aria-describedby={joinDescribedBy(
                      "appointment-reason-counter",
                      !!errors.reason && "appointment-reason-error"
                    )}
                    {...register("reason")}
                  />
                  <div className="flex items-start justify-between gap-3">
                    <FieldError id="appointment-reason-error">{errors.reason?.message}</FieldError>
                    <span
                      id="appointment-reason-counter"
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
                {t("modals.create.cancel")}
              </Button>
              <Button type="submit" className={touchTargetClassName} disabled={isBusy} aria-busy={isBusy}>
                {isBusy && <Spinner />}
                {isEditMode ? t("modals.edit.confirm") : t("modals.create.confirm")}
              </Button>
            </DialogFooter>
          </form>
        </DialogContent>
      </Dialog>
      <DiscardChangesDialog open={isConfirmationOpen} onKeepEditing={keepEditing} onDiscard={discardChanges} />
    </>
  )
}
