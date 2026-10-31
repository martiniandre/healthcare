import { useState, useEffect, useMemo } from "react"
import { useForm, Controller, useWatch } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import { useTranslation } from "react-i18next"
import { isAxiosError } from "axios"
import { Alert, AlertDescription } from "@/shared/components/ui/alert"
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
import { Field, FieldError, FieldLabel } from "@/shared/components/ui/field"
import { Button } from "@/shared/components/ui/button"
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from "@/shared/components/ui/dialog"
import { getNewAppointmentSchema, type NewAppointmentFormData } from "../schedule_schemas"
import { todayDateString } from "../../../shared/utils/validators"
import { getAvailableStartTimeOptions, getEndTimeOptionsForStart } from "../schedule_time_options"
import { useStaffListQuery } from "../../staff/queries"
import { usePatientsQuery, usePatientQuery } from "../../patients/queries"
import { useIdempotencyKey } from "../hooks/useIdempotencyKey"
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
  const isEditMode = appointment !== undefined && appointment !== null
  const { data: staffMembers = [] } = useStaffListQuery()
  const { data: patientsPage } = usePatientsQuery("", "", "", 1, 100)
  const { data: appointmentPatient } = usePatientQuery(
    isEditMode && appointment ? appointment.patient_fhir_id : ""
  )
  const patients = patientsPage?.patients ?? []
  const { getOrCreateKey, resetKey } = useIdempotencyKey()

  const [patientSearch, setPatientSearch] = useState("")
  const [conflictMessage, setConflictMessage] = useState<string | null>(null)

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
    patient.full_name.toLowerCase().includes(patientSearch.toLowerCase())
  )

  const {
    register,
    handleSubmit,
    control,
    reset,
    formState: { errors },
  } = useForm<NewAppointmentFormData>({
    resolver: zodResolver(getNewAppointmentSchema(t)),
    defaultValues: {
      patientFhirId: "",
      staffId: defaultStaffId ?? "",
      date: defaultDate ?? "",
      startTime: "",
      endTime: "",
      reason: "",
    },
  })

  const selectedDate = useWatch({ control, name: "date" }) ?? ""
  const selectedStartTime = useWatch({ control, name: "startTime" }) ?? ""
  const availableStartTimeOptions = getAvailableStartTimeOptions(selectedDate)
  const endTimeOptions = getEndTimeOptionsForStart(selectedStartTime)

  const handleClose = () => {
    setConflictMessage(null)
    setPatientSearch("")
    reset()
    onClose()
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
          endTime: "",
          reason: "",
        })
      }
    }
  }, [isOpen, appointment, isEditMode, defaultStaffId, defaultDate, defaultStartTime, reset])

  if (!isOpen) {
    return null
  }

  const handleFormSubmit = handleSubmit(async (formData) => {
    setConflictMessage(null)
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
        setConflictMessage(t("errors.conflict"))
      }
    }
  })

  return (
    <Dialog open={isOpen} onOpenChange={(open) => !open && handleClose()}>
      <DialogContent className="sm:max-w-[520px]">
        <DialogHeader>
          <DialogTitle className="text-left">
            {isEditMode ? t("modals.edit.title") : t("modals.create.title")}
          </DialogTitle>
        </DialogHeader>
        <form onSubmit={handleFormSubmit} noValidate className="flex flex-col gap-4 text-left mt-4">
          {conflictMessage && (
            <Alert variant="destructive">
              <AlertDescription>{conflictMessage}</AlertDescription>
            </Alert>
          )}
          <Field data-invalid={!!errors.patientFhirId}>
            <FieldLabel htmlFor="appointment-patient">
              {t("modals.create.patient")}
            </FieldLabel>
            <Input
              type="text"
              value={patientSearch}
              onChange={(inputEvent) => setPatientSearch(inputEvent.target.value)}
              placeholder={t("modals.create.patientSearchPlaceholder")}
              className="mb-1"
            />
            <Controller
              control={control}
              name="patientFhirId"
              render={({ field }) => (
                <Select onValueChange={field.onChange} value={field.value}>
                  <SelectTrigger id="appointment-patient" className="w-full" aria-invalid={!!errors.patientFhirId}>
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
            <FieldError>{errors.patientFhirId?.message}</FieldError>
          </Field>

          <Field data-invalid={!!errors.staffId}>
            <FieldLabel htmlFor="appointment-staff">
              {t("modals.create.staff")}
            </FieldLabel>
            <Controller
              control={control}
              name="staffId"
              render={({ field }) => (
                <Select onValueChange={field.onChange} value={field.value}>
                  <SelectTrigger id="appointment-staff" className="w-full" aria-invalid={!!errors.staffId}>
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
            <FieldError>{errors.staffId?.message}</FieldError>
          </Field>

          <div className="grid grid-cols-3 gap-3">
            <Field data-invalid={!!errors.date}>
              <FieldLabel htmlFor="appointment-date">
                {t("modals.create.date")}
              </FieldLabel>
              <Input id="appointment-date" type="date" min={todayDateString()} aria-invalid={!!errors.date} {...register("date")} />
              <FieldError>{errors.date?.message}</FieldError>
            </Field>
            <Field data-invalid={!!errors.startTime}>
              <FieldLabel htmlFor="appointment-start-time">
                {t("modals.create.startTime")}
              </FieldLabel>
              <Controller
                control={control}
                name="startTime"
                render={({ field }) => (
                  <Select onValueChange={field.onChange} value={field.value}>
                    <SelectTrigger id="appointment-start-time" className="w-full" aria-invalid={!!errors.startTime}>
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
              <FieldError>{errors.startTime?.message}</FieldError>
            </Field>
            <Field data-invalid={!!errors.endTime}>
              <FieldLabel htmlFor="appointment-end-time">
                {t("modals.create.endTime")}
              </FieldLabel>
              <Controller
                control={control}
                name="endTime"
                render={({ field }) => (
                  <Select onValueChange={field.onChange} value={field.value}>
                    <SelectTrigger id="appointment-end-time" className="w-full" aria-invalid={!!errors.endTime}>
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
              <FieldError>{errors.endTime?.message}</FieldError>
            </Field>
          </div>

          <Field data-invalid={!!errors.reason}>
            <FieldLabel htmlFor="appointment-reason">
              {t("modals.create.reason")}
            </FieldLabel>
            <Textarea
              id="appointment-reason"
              placeholder={t("modals.create.reasonPlaceholder")}
              aria-invalid={!!errors.reason}
              {...register("reason")}
            />
            <FieldError>{errors.reason?.message}</FieldError>
          </Field>

          <DialogFooter className="mt-4">
            <Button variant="outline" type="button" onClick={handleClose}>
              {t("modals.create.cancel")}
            </Button>
            <Button type="submit" disabled={isPending}>
              {isPending && <Spinner />}
              {isEditMode ? t("modals.edit.confirm") : t("modals.create.confirm")}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}
