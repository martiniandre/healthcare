import { useState, useEffect } from "react"
import { useForm, Controller } from "react-hook-form"
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
import { useStaffListQuery } from "../../staff/queries"
import { usePatientsQuery } from "../../patients/queries"
import { useIdempotencyKey } from "../hooks/useIdempotencyKey"
import type { CreateAppointmentPayload } from "../types"

interface AppointmentModalProps {
  isOpen: boolean
  onClose: () => void
  onSubmit: (payload: CreateAppointmentPayload) => Promise<void>
  isPending: boolean
  defaultStaffId?: string
  defaultDate?: string
  defaultStartTime?: string
}

const formatLocalDateTime = (dateValue: string, timeValue: string): string => {
  return new Date(`${dateValue}T${timeValue}`).toISOString()
}

export const AppointmentModal = ({
  isOpen,
  onClose,
  onSubmit,
  isPending,
  defaultStaffId,
  defaultDate,
  defaultStartTime,
}: AppointmentModalProps) => {
  const { t } = useTranslation("schedule")
  const { data: staffMembers = [] } = useStaffListQuery()
  const { data: patients = [] } = usePatientsQuery("", "", "", 1, 100)
  const { getOrCreateKey, resetKey } = useIdempotencyKey()

  const [patientSearch, setPatientSearch] = useState("")
  const [conflictMessage, setConflictMessage] = useState<string | null>(null)

  const filteredPatients = patients.filter((patient) =>
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

  const handleClose = () => {
    setConflictMessage(null)
    setPatientSearch("")
    reset()
    onClose()
  }

  useEffect(() => {
    if (isOpen) {
      reset({
        patientFhirId: "",
        staffId: defaultStaffId ?? "",
        date: defaultDate ?? "",
        startTime: defaultStartTime ?? "",
        endTime: "",
        reason: "",
      })
    }
  }, [isOpen, defaultStaffId, defaultDate, defaultStartTime, reset])

  if (!isOpen) {
    return null
  }

  const handleFormSubmit = handleSubmit(async (formData) => {
    setConflictMessage(null)
    const createPayload: CreateAppointmentPayload = {
      patient_fhir_id: formData.patientFhirId,
      staff_id: formData.staffId,
      starts_at: formatLocalDateTime(formData.date, formData.startTime),
      ends_at: formatLocalDateTime(formData.date, formData.endTime),
      reason: formData.reason,
      idempotency_key: getOrCreateKey(),
    }

    try {
      await onSubmit(createPayload)
      resetKey()
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
            {t("modals.create.title")}
          </DialogTitle>
        </DialogHeader>
        <form onSubmit={handleFormSubmit} className="flex flex-col gap-4 text-left mt-4">
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
                    {staffMembers.filter((staffMember) => staffMember.fhirResourceId).map((staffMember) => (
                      <SelectItem key={staffMember.id} value={staffMember.fhirResourceId}>
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
              <Input type="date" errorText={errors.date?.message} {...register("date")} />
            </div>
            <div className="flex flex-col gap-1">
              <label className="text-xs font-semibold text-gray-600">
                {t("modals.create.startTime")}
              </label>
              <Input type="time" step={900} errorText={errors.startTime?.message} {...register("startTime")} />
            </div>
            <div className="flex flex-col gap-1">
              <label className="text-xs font-semibold text-gray-600">
                {t("modals.create.endTime")}
              </label>
              <Input type="time" step={900} errorText={errors.endTime?.message} {...register("endTime")} />
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
              {t("modals.create.confirm")}
            </Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  )
}
