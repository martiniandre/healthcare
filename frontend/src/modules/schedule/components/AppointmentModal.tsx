import { useState, useEffect, useMemo } from "react"
import { useForm, Controller, useWatch } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import { useTranslation } from "react-i18next"
import { isAxiosError } from "axios"
import { Input } from "../../../shared/components/ui/Input"
import {
  Select,
  SelectTrigger,
  SelectValue,
  SelectContent,
  SelectItem,
} from "../../../shared/components/ui/Select"
import { Button } from "../../../shared/components/ui/Button"
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "../../../shared/components/ui/Dialog"
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
            <div className="text-xs font-semibold text-red-600 bg-red-50 border border-red-200 rounded-lg px-3 py-2">
              {conflictMessage}
            </div>
          )}
          <div className="flex flex-col gap-1">
            <label className="text-xs font-semibold text-gray-600">
              {t("modals.create.patient")}
            </label>
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
                  <SelectTrigger className="w-full">
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
            {errors.patientFhirId?.message && (
              <span className="text-xs text-red-500 font-medium px-1 mt-1">
                {errors.patientFhirId.message}
              </span>
            )}
          </div>

          <div className="flex flex-col gap-1">
            <label className="text-xs font-semibold text-gray-600">
              {t("modals.create.staff")}
            </label>
            <Controller
              control={control}
              name="staffId"
              render={({ field }) => (
                <Select onValueChange={field.onChange} value={field.value}>
                  <SelectTrigger className="w-full">
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
            {errors.staffId?.message && (
              <span className="text-xs text-red-500 font-medium px-1 mt-1">
                {errors.staffId.message}
              </span>
            )}
          </div>

          <div className="grid grid-cols-3 gap-3">
            <div className="flex flex-col gap-1">
              <label className="text-xs font-semibold text-gray-600">
                {t("modals.create.date")}
              </label>
              <Input type="date" min={todayDateString()} errorText={errors.date?.message} {...register("date")} />
            </div>
            <div className="flex flex-col gap-1">
              <label className="text-xs font-semibold text-gray-600">
                {t("modals.create.startTime")}
              </label>
              <select
                aria-label={t("modals.create.startTime")}
                className="flex h-9 w-full rounded-md border border-input bg-background px-3 py-1 text-sm shadow-sm focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50"
                {...register("startTime")}
              >
                <option value="">{t("modals.create.selectStartTime")}</option>
                {availableStartTimeOptions.map((timeSlot) => (
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
                {t("modals.create.endTime")}
              </label>
              <select
                aria-label={t("modals.create.endTime")}
                className="flex h-9 w-full rounded-md border border-input bg-background px-3 py-1 text-sm shadow-sm focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50"
                {...register("endTime")}
              >
                <option value="">{t("modals.create.selectEndTime")}</option>
                {endTimeOptions.map((timeSlot) => (
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
              {t("modals.create.reason")}
            </label>
            <textarea
              className="flex min-h-[80px] w-full rounded-lg border border-border bg-transparent px-3 py-2 text-sm shadow-sm focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50"
              placeholder={t("modals.create.reasonPlaceholder")}
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
              {t("modals.create.cancel")}
            </Button>
            <Button type="submit" disabled={isPending}>
              {isEditMode ? t("modals.edit.confirm") : t("modals.create.confirm")}
            </Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  )
}
