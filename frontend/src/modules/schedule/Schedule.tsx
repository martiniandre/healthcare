import { useState } from "react"
import { useTranslation } from "react-i18next"
import { isAxiosError } from "axios"
import { CalendarPlus, CalendarClock, CalendarOff, Loader2 } from "lucide-react"
import { useStaffListQuery } from "../staff/queries"
import {
  Select,
  SelectTrigger,
  SelectValue,
  SelectContent,
  SelectItem,
} from "../../shared/components/ui/Select"
import { Input } from "../../shared/components/ui/Input"
import { Button } from "../../shared/components/ui/Button"
import { AppointmentModal } from "./components/AppointmentModal"
import { CancelAppointmentModal } from "./components/CancelAppointmentModal"
import { AppointmentCard } from "./components/AppointmentCard"
import { UnavailabilityCard } from "./components/UnavailabilityCard"
import { CreateUnavailabilityModal } from "./components/CreateUnavailabilityModal"
import { DeleteUnavailabilityModal } from "./components/DeleteUnavailabilityModal"
import {
  useStaffDayAppointmentsQuery,
  useCreateAppointmentMutation,
  useCancelAppointmentMutation,
  useStaffUnavailabilityQuery,
  useCreateUnavailabilityMutation,
  useDeleteUnavailabilityMutation,
} from "./queries"
import { toast } from "../../shared/store/toast_store"
import type { Appointment, CreateAppointmentPayload, StaffUnavailability, CreateUnavailabilityPayload } from "./types"

const todayDate = (): string => {
  const localDate = new Date()
  const year = localDate.getFullYear()
  const month = String(localDate.getMonth() + 1).padStart(2, "0")
  const day = String(localDate.getDate()).padStart(2, "0")
  return `${year}-${month}-${day}`
}

export const Schedule = () => {
  const { t } = useTranslation("schedule")
  const { data: staffMembers = [] } = useStaffListQuery()

  const [selectedStaffId, setSelectedStaffId] = useState("")
  const [selectedDate, setSelectedDate] = useState(todayDate())
  const [isCreateModalOpen, setIsCreateModalOpen] = useState(false)
  const [appointmentToCancel, setAppointmentToCancel] = useState<Appointment | null>(null)
  const [isCreateUnavailabilityOpen, setIsCreateUnavailabilityOpen] = useState(false)
  const [unavailabilityToDelete, setUnavailabilityToDelete] = useState<StaffUnavailability | null>(null)

  const staffMemberOptions = staffMembers.filter((staffMember) => staffMember.fhirResourceId)

  const { data: appointments = [], isLoading: isDayLoading } = useStaffDayAppointmentsQuery(
    selectedStaffId,
    selectedDate
  )

  const { data: unavailabilityWindows = [] } = useStaffUnavailabilityQuery(selectedStaffId)

  const createAppointmentMutation = useCreateAppointmentMutation()
  const cancelAppointmentMutation = useCancelAppointmentMutation()
  const createUnavailabilityMutation = useCreateUnavailabilityMutation()
  const deleteUnavailabilityMutation = useDeleteUnavailabilityMutation()

  const sortedAppointments = [...appointments].sort(
    (firstAppointment, secondAppointment) =>
      new Date(firstAppointment.starts_at).getTime() - new Date(secondAppointment.starts_at).getTime()
  )

  const handleCreateAppointment = async (payload: CreateAppointmentPayload) => {
    try {
      await createAppointmentMutation.mutateAsync(payload)
      toast.success(t("toasts.createSuccess"))
    } catch (createError) {
      if (isAxiosError(createError) && createError.response?.status === 409) {
        throw createError
      }
      toast.error(t("toasts.createError"))
      throw createError
    }
  }

  const handleConfirmCancel = (appointmentId: string) => {
    cancelAppointmentMutation.mutate(appointmentId, {
      onSuccess: () => {
        toast.success(t("toasts.cancelSuccess"))
        setAppointmentToCancel(null)
      },
      onError: () => {
        toast.error(t("toasts.cancelError"))
      },
    })
  }

  const handleCreateUnavailability = async (payload: CreateUnavailabilityPayload) => {
    try {
      await createUnavailabilityMutation.mutateAsync(payload)
      toast.success(t("unavailability.toasts.createSuccess"))
    } catch {
      toast.error(t("unavailability.toasts.createError"))
    }
  }

  const handleDeleteUnavailability = (unavailabilityId: string) => {
    deleteUnavailabilityMutation.mutate(unavailabilityId, {
      onSuccess: () => {
        toast.success(t("unavailability.toasts.deleteSuccess"))
        setUnavailabilityToDelete(null)
      },
      onError: () => {
        toast.error(t("unavailability.toasts.deleteError"))
      },
    })
  }

  const dayUnavailabilityWindows = unavailabilityWindows.filter((window) => {
    const windowDate = window.starts_at.slice(0, 10)
    return windowDate === selectedDate
  })

  return (
    <div className="flex-1 p-4 sm:p-6 md:p-8 flex flex-col gap-4 md:gap-6 max-w-7xl mx-auto w-full">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <h1 className="text-xl font-bold text-gray-900">{t("title")}</h1>
          <p className="text-sm text-gray-500">{t("subtitle")}</p>
        </div>
        <div className="flex gap-2">
          <Button
            variantType="outline"
            onClick={() => setIsCreateUnavailabilityOpen(true)}
            disabled={!selectedStaffId}
          >
            <CalendarOff className="w-4 h-4" />
            {t("unavailability.createButton")}
          </Button>
          <Button onClick={() => setIsCreateModalOpen(true)} disabled={!selectedStaffId}>
            <CalendarPlus className="w-4 h-4" />
            {t("newAppointment")}
          </Button>
        </div>
      </div>

      <div className="bg-white border border-border rounded-xl p-4 flex flex-col sm:flex-row gap-4">
        <div className="flex flex-col gap-1 flex-1">
          <label className="text-xs font-semibold text-gray-600">{t("selectStaff")}</label>
          <Select onValueChange={setSelectedStaffId} value={selectedStaffId}>
            <SelectTrigger className="w-full">
              <SelectValue placeholder={t("selectStaffPlaceholder")} />
            </SelectTrigger>
            <SelectContent>
              {staffMemberOptions.map((staffMember) => (
                <SelectItem key={staffMember.id} value={staffMember.fhirResourceId}>
                  {staffMember.fullName} — {staffMember.role}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>
        <div className="flex flex-col gap-1">
          <label className="text-xs font-semibold text-gray-600">{t("selectDate")}</label>
          <Input type="date" value={selectedDate} onChange={(inputEvent) => setSelectedDate(inputEvent.target.value)} />
        </div>
      </div>

      <div className="flex flex-col gap-4">
        <div className="flex items-center gap-2">
          <CalendarClock className="w-4 h-4 text-primary" />
          <h2 className="text-sm font-bold text-gray-700">
            {t("dayAgenda", { date: new Date(`${selectedDate}T12:00:00`).toLocaleDateString() })}
          </h2>
        </div>

        {isDayLoading ? (
          <div className="flex items-center justify-center py-16 bg-white border border-border rounded-xl">
            <Loader2 className="w-6 h-6 text-primary animate-spin" />
          </div>
        ) : !selectedStaffId ? (
          <div className="py-16 text-center bg-white border border-border rounded-xl">
            <p className="text-sm text-gray-500">{t("emptyState.selectStaff")}</p>
          </div>
        ) : sortedAppointments.length === 0 && dayUnavailabilityWindows.length === 0 ? (
          <div className="py-16 text-center bg-white border border-border rounded-xl">
            <p className="text-sm text-gray-500">{t("emptyState.noAppointments")}</p>
          </div>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-4">
            {dayUnavailabilityWindows.map((unavailability) => (
              <UnavailabilityCard
                key={unavailability.id}
                unavailability={unavailability}
                onDelete={setUnavailabilityToDelete}
              />
            ))}
            {sortedAppointments.map((appointment) => (
              <AppointmentCard
                key={appointment.id}
                appointment={appointment}
                onCancel={setAppointmentToCancel}
              />
            ))}
          </div>
        )}
      </div>

      <AppointmentModal
        isOpen={isCreateModalOpen}
        onClose={() => setIsCreateModalOpen(false)}
        onSubmit={handleCreateAppointment}
        isPending={createAppointmentMutation.isPending}
        defaultStaffId={selectedStaffId}
        defaultDate={selectedDate}
      />

      <CancelAppointmentModal
        isOpen={appointmentToCancel !== null}
        appointment={appointmentToCancel}
        onClose={() => setAppointmentToCancel(null)}
        onConfirm={handleConfirmCancel}
        isPending={cancelAppointmentMutation.isPending}
      />

      <CreateUnavailabilityModal
        isOpen={isCreateUnavailabilityOpen}
        onClose={() => setIsCreateUnavailabilityOpen(false)}
        onSubmit={handleCreateUnavailability}
        isPending={createUnavailabilityMutation.isPending}
        staffId={selectedStaffId}
        defaultDate={selectedDate}
      />

      <DeleteUnavailabilityModal
        isOpen={unavailabilityToDelete !== null}
        onClose={() => setUnavailabilityToDelete(null)}
        onConfirm={() => {
          if (unavailabilityToDelete) {
            handleDeleteUnavailability(unavailabilityToDelete.id)
          }
        }}
        isPending={deleteUnavailabilityMutation.isPending}
      />
    </div>
  )
}
