import { useMemo, useState } from "react"
import { useTranslation } from "react-i18next"
import { isAxiosError } from "axios"
import { CalendarPlus, CalendarOff, Loader2 } from "lucide-react"
import { useStaffListQuery } from "../staff/queries"
import { Button } from "../../shared/components/ui/Button"
import { AppointmentModal } from "./components/AppointmentModal"
import { CancelAppointmentModal } from "./components/CancelAppointmentModal"
import { UnavailabilityCard } from "./components/UnavailabilityCard"
import { CreateUnavailabilityModal } from "./components/CreateUnavailabilityModal"
import { DeleteUnavailabilityModal } from "./components/DeleteUnavailabilityModal"
import { StaffOverlaySidebar } from "./components/StaffOverlaySidebar"
import { ScheduleMonthCalendar } from "./components/ScheduleMonthCalendar"
import {
  useCreateAppointmentMutation,
  useCancelAppointmentMutation,
  useStaffUnavailabilityQuery,
  useCreateUnavailabilityMutation,
  useDeleteUnavailabilityMutation,
  useStaffRangeAppointmentsQueries,
} from "./queries"
import { appointmentsToCalendarEvents, staffColorForIndex } from "./schedule_calendar_helpers"
import type { CalendarEventShape } from "./schedule_calendar_helpers"
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

  const [selectedStaffIds, setSelectedStaffIds] = useState<string[] | null>(null)
  const [visibleRange, setVisibleRange] = useState({ startDate: todayDate(), endDate: todayDate() })
  const [isCreateModalOpen, setIsCreateModalOpen] = useState(false)
  const [appointmentToCancel, setAppointmentToCancel] = useState<Appointment | null>(null)
  const [isCreateUnavailabilityOpen, setIsCreateUnavailabilityOpen] = useState(false)
  const [unavailabilityToDelete, setUnavailabilityToDelete] = useState<StaffUnavailability | null>(null)

  const staffMemberOptions = useMemo(
    () => staffMembers.filter((staffMember) => staffMember.fhirResourceId),
    [staffMembers]
  )

  const effectiveSelectedStaffIds = useMemo(
    () =>
      selectedStaffIds !== null && selectedStaffIds.length > 0
        ? selectedStaffIds
        : selectedStaffIds !== null
          ? []
          : staffMemberOptions.length > 0
            ? [staffMemberOptions[0].fhirResourceId]
            : [],
    [selectedStaffIds, staffMemberOptions]
  )

  const rangeFilters = useMemo(
    () =>
      effectiveSelectedStaffIds.map((staffId) => ({
        staffId,
        startDate: visibleRange.startDate,
        endDate: visibleRange.endDate,
      })),
    [effectiveSelectedStaffIds, visibleRange]
  )

  const rangeQueryResults = useStaffRangeAppointmentsQueries(rangeFilters)
  const rangeQueriesLoading = rangeQueryResults.some((queryResult) => queryResult.isLoading)

  const calendarEvents = useMemo<CalendarEventShape[]>(() => {
    const mergedEvents: CalendarEventShape[] = []
    effectiveSelectedStaffIds.forEach((staffId, staffIndex) => {
      const queryResultIndex = rangeFilters.findIndex((rangeFilter) => rangeFilter.staffId === staffId)
      const staffAppointments = queryResultIndex >= 0 ? rangeQueryResults[queryResultIndex]?.data ?? [] : []
      mergedEvents.push(...appointmentsToCalendarEvents(staffAppointments, staffColorForIndex(staffIndex)))
    })
    return mergedEvents
  }, [effectiveSelectedStaffIds, rangeQueryResults, rangeFilters])

  const focusedStaffId = effectiveSelectedStaffIds[0] ?? ""

  const { data: unavailabilityWindows = [] } = useStaffUnavailabilityQuery(focusedStaffId)

  const createAppointmentMutation = useCreateAppointmentMutation()
  const cancelAppointmentMutation = useCancelAppointmentMutation()
  const createUnavailabilityMutation = useCreateUnavailabilityMutation()
  const deleteUnavailabilityMutation = useDeleteUnavailabilityMutation()

  const handleToggleStaff = (toggledStaffId: string) => {
    const baseSelection = effectiveSelectedStaffIds
    setSelectedStaffIds(
      baseSelection.includes(toggledStaffId)
        ? baseSelection.filter((selectedId) => selectedId !== toggledStaffId)
        : [...baseSelection, toggledStaffId]
    )
  }

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

  const upcomingUnavailabilityWindows = unavailabilityWindows

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
            disabled={!focusedStaffId}
          >
            <CalendarOff className="w-4 h-4" />
            {t("unavailability.createButton")}
          </Button>
          <Button onClick={() => setIsCreateModalOpen(true)} disabled={!focusedStaffId}>
            <CalendarPlus className="w-4 h-4" />
            {t("newAppointment")}
          </Button>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-[280px_1fr] gap-4 items-start">
        <StaffOverlaySidebar
          staffMembers={staffMemberOptions.map((staffMember) => ({
            id: staffMember.fhirResourceId,
            fullName: staffMember.fullName,
            role: staffMember.role,
          }))}
          selectedStaffIds={effectiveSelectedStaffIds}
          onToggleStaff={handleToggleStaff}
        />

        <div className="flex flex-col gap-4 min-w-0">
          {rangeQueriesLoading && effectiveSelectedStaffIds.length > 0 ? (
            <div className="flex items-center justify-center py-16 bg-white border border-border rounded-xl">
              <Loader2 className="w-6 h-6 text-primary animate-spin" />
            </div>
          ) : (
            <ScheduleMonthCalendar
              events={calendarEvents}
              onVisibleRangeChange={(rangeStart, rangeEnd) =>
                setVisibleRange((previousRange) =>
                  previousRange.startDate === rangeStart && previousRange.endDate === rangeEnd
                    ? previousRange
                    : { startDate: rangeStart, endDate: rangeEnd }
                )
              }
            />
          )}

          {upcomingUnavailabilityWindows.length > 0 && (
            <div className="flex flex-col gap-3">
              <h2 className="text-sm font-bold text-gray-700">{t("unavailability.blockedPeriodsTitle")}</h2>
              <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-4">
                {upcomingUnavailabilityWindows.map((unavailability) => (
                  <UnavailabilityCard
                    key={unavailability.id}
                    unavailability={unavailability}
                    onDelete={(id) => {
                      const found = upcomingUnavailabilityWindows.find((entry) => entry.id === id)
                      setUnavailabilityToDelete(found ?? null)
                    }}
                  />
                ))}
              </div>
            </div>
          )}
        </div>
      </div>

      <AppointmentModal
        isOpen={isCreateModalOpen}
        onClose={() => setIsCreateModalOpen(false)}
        onSubmit={handleCreateAppointment}
        isPending={createAppointmentMutation.isPending}
        defaultStaffId={focusedStaffId}
        defaultDate={todayDate()}
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
        staffId={focusedStaffId}
        defaultDate={todayDate()}
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
