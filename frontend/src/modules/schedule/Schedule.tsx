import { useMemo, useState } from "react"
import { useTranslation } from "react-i18next"
import { isAxiosError } from "axios"
import { CalendarPlus, CalendarOff, TriangleAlert } from "lucide-react"
import { useStaffListQuery } from "../staff/queries"
import { Button } from "@/shared/components/ui/button"
import { Spinner } from "@/shared/components/ui/spinner"
import { Skeleton } from "@/shared/components/ui/skeleton"
import { Alert, AlertDescription, AlertTitle } from "@/shared/components/ui/alert"
import { EmptyState } from "@/shared/components/ui/EmptyState"
import { PageContainer, PageHeader } from "@/shared/components/ui/PageContainer"
import { STAFF_ROLE_LABEL_KEYS } from "../staff/staff_roles"
import type { StaffRole } from "@/shared/types"
import { AppointmentModal } from "./components/AppointmentModal"
import { CancelAppointmentModal } from "./components/CancelAppointmentModal"
import { AppointmentDetailsModal } from "./components/AppointmentDetailsModal"
import { UnavailabilityCard } from "./components/UnavailabilityCard"
import { CreateUnavailabilityModal } from "./components/CreateUnavailabilityModal"
import { DeleteUnavailabilityModal } from "./components/DeleteUnavailabilityModal"
import { StaffOverlaySidebar } from "./components/StaffOverlaySidebar"
import { ScheduleCalendar, type ScheduleViewMode } from "./components/ScheduleCalendar"
import { ScheduleViewToggle } from "./components/ScheduleViewToggle"
import {
  useCreateAppointmentMutation,
  useCancelAppointmentMutation,
  useUpdateAppointmentMutation,
  useStaffUnavailabilityQuery,
  useCreateUnavailabilityMutation,
  useDeleteUnavailabilityMutation,
  useStaffRangeAppointmentsQueries,
} from "./queries"
import { appointmentsToCalendarEvents, staffColorForIndex } from "./schedule_calendar_helpers"
import type { CalendarEventShape } from "./schedule_calendar_helpers"
import { toast } from "../../shared/store/toast_store"
import type { Appointment, CreateAppointmentPayload, UpdateAppointmentPayload, StaffUnavailability, CreateUnavailabilityPayload } from "./types"

const todayDate = (): string => {
  const localDate = new Date()
  const year = localDate.getFullYear()
  const month = String(localDate.getMonth() + 1).padStart(2, "0")
  const day = String(localDate.getDate()).padStart(2, "0")
  return `${year}-${month}-${day}`
}

const toLocalDateString = (dateValue: Date): string => {
  const year = dateValue.getFullYear()
  const month = String(dateValue.getMonth() + 1).padStart(2, "0")
  const day = String(dateValue.getDate()).padStart(2, "0")
  return `${year}-${month}-${day}`
}

const toLocalTimeString = (dateValue: Date): string => {
  const hour = String(dateValue.getHours()).padStart(2, "0")
  const minute = String(dateValue.getMinutes()).padStart(2, "0")
  return `${hour}:${minute}`
}

export const Schedule = () => {
  const { t } = useTranslation("schedule")
  const { t: translateSidebar } = useTranslation("sidebar")
  const { t: translateStaff } = useTranslation("staff")
  const {
    data: staffMembers = [],
    isLoading: isStaffLoading,
    isError: isStaffError,
    refetch: refetchStaff,
  } = useStaffListQuery()

  const [selectedStaffIds, setSelectedStaffIds] = useState<string[] | null>(null)
  const [visibleRange, setVisibleRange] = useState({ startDate: todayDate(), endDate: todayDate() })
  const [viewMode, setViewMode] = useState<ScheduleViewMode>("week")
  const [createDefaults, setCreateDefaults] = useState<{ date: string; startTime: string } | null>(null)
  const [isCreateModalOpen, setIsCreateModalOpen] = useState(false)
  const [appointmentToView, setAppointmentToView] = useState<Appointment | null>(null)
  const [appointmentToEdit, setAppointmentToEdit] = useState<Appointment | null>(null)
  const [appointmentToCancel, setAppointmentToCancel] = useState<Appointment | null>(null)
  const [isCreateUnavailabilityOpen, setIsCreateUnavailabilityOpen] = useState(false)
  const [unavailabilityToDelete, setUnavailabilityToDelete] = useState<StaffUnavailability | null>(null)

  const staffMemberOptions = useMemo(
    () => staffMembers.filter((staffMember) => staffMember.id),
    [staffMembers]
  )

  const effectiveSelectedStaffIds = useMemo(() => {
    if (selectedStaffIds !== null) {
      return selectedStaffIds
    }
    const defaultSelectedStaffId = staffMemberOptions[0]?.id
    return defaultSelectedStaffId ? [defaultSelectedStaffId] : []
  }, [selectedStaffIds, staffMemberOptions])

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
  const rangeQueriesFailed = rangeQueryResults.some((queryResult) => queryResult.isError)

  const handleRetryAppointments = () => {
    rangeQueryResults.forEach((queryResult) => {
      if (queryResult.isError) {
        void queryResult.refetch()
      }
    })
  }

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

  const {
    data: unavailabilityWindows = [],
    isLoading: isUnavailabilityLoading,
    isError: isUnavailabilityError,
    refetch: refetchUnavailability,
  } = useStaffUnavailabilityQuery(focusedStaffId)

  const createAppointmentMutation = useCreateAppointmentMutation()
  const cancelAppointmentMutation = useCancelAppointmentMutation()
  const updateAppointmentMutation = useUpdateAppointmentMutation()
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

  const handleCalendarSlotClick = (clickedStart: Date) => {
    const dateString = toLocalDateString(clickedStart)
    const timeString = toLocalTimeString(clickedStart)
    setCreateDefaults({ date: dateString, startTime: timeString })
    setIsCreateModalOpen(true)
  }

  const handleNewAppointmentButton = () => {
    setCreateDefaults(null)
    setIsCreateModalOpen(true)
  }

  const handleRescheduleAppointment = (appointment: Appointment, newStart: Date, newEnd: Date) => {
    updateAppointmentMutation.mutate(
      {
        appointmentId: appointment.id,
        payload: {
          patient_fhir_id: appointment.patient_fhir_id,
          staff_id: appointment.staff_id,
          starts_at: newStart.toISOString(),
          ends_at: newEnd.toISOString(),
          reason: appointment.reason,
        },
      },
      {
        onSuccess: () => {
          toast.success(t("toasts.rescheduleSuccess"))
        },
        onError: (updateError) => {
          if (isAxiosError(updateError) && updateError.response?.status === 409) {
            toast.error(t("errors.conflict"))
            return
          }
          toast.error(t("toasts.rescheduleError"))
        },
      }
    )
  }

  const handleUpdateAppointment = async (payload: UpdateAppointmentPayload) => {
    if (!appointmentToEdit) {
      return
    }
    try {
      await updateAppointmentMutation.mutateAsync({ appointmentId: appointmentToEdit.id, payload })
      toast.success(t("toasts.updateSuccess"))
    } catch (updateError) {
      if (isAxiosError(updateError) && updateError.response?.status === 409) {
        throw updateError
      }
      toast.error(t("toasts.updateError"))
      throw updateError
    }
  }

  const handleEventClick = (appointment: Appointment) => {
    setAppointmentToView(appointment)
  }

  const handleEditFromDetails = (appointment: Appointment) => {
    setAppointmentToView(null)
    setAppointmentToEdit(appointment)
  }

  const handleCancelFromDetails = (appointment: Appointment) => {
    setAppointmentToView(null)
    setAppointmentToCancel(appointment)
  }

  const handleCreateUnavailability = async (payload: CreateUnavailabilityPayload) => {
    try {
      await createUnavailabilityMutation.mutateAsync(payload)
      toast.success(t("unavailability.toasts.createSuccess"))
    } catch (createError) {
      toast.error(t("unavailability.toasts.createError"))
      throw createError
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
    <PageContainer>
      <PageHeader
        eyebrow={translateSidebar("topics.clinical")}
        title={t("title")}
        description={t("subtitle")}
        actions={
          <>
            <Button
              variant="outline"
              size="lg"
              onClick={() => setIsCreateUnavailabilityOpen(true)}
              disabled={!focusedStaffId}
            >
              <CalendarOff />
              {t("unavailability.createButton")}
            </Button>
            <Button size="lg" onClick={handleNewAppointmentButton} disabled={!focusedStaffId}>
              <CalendarPlus />
              {t("newAppointment")}
            </Button>
          </>
        }
      />

      <div className="grid grid-cols-1 items-start gap-5 lg:grid-cols-[272px_1fr]">
        <StaffOverlaySidebar
          viewControl={<ScheduleViewToggle value={viewMode} onChange={setViewMode} />}
          staffMembers={staffMemberOptions.map((staffMember) => ({
            id: staffMember.id,
            fullName: staffMember.fullName,
            role: translateStaff(STAFF_ROLE_LABEL_KEYS[staffMember.role as StaffRole] ?? staffMember.role, { defaultValue: staffMember.role }),
          }))}
          selectedStaffIds={effectiveSelectedStaffIds}
          onToggleStaff={handleToggleStaff}
          isLoading={isStaffLoading}
          isError={isStaffError}
          onRetry={() => refetchStaff()}
        />

        <div className="flex min-w-0 flex-col gap-5">
          {rangeQueriesFailed && (
            <Alert variant="destructive">
              <TriangleAlert aria-hidden="true" />
              <AlertTitle>{t("errors.loadAppointmentsTitle")}</AlertTitle>
              <AlertDescription>
                <p>{t("errors.loadAppointmentsDescription")}</p>
                <Button type="button" variant="outline" size="sm" onClick={handleRetryAppointments} className="mt-2">
                  {t("retry")}
                </Button>
              </AlertDescription>
            </Alert>
          )}
          <div className="relative">
            <ScheduleCalendar
              events={calendarEvents}
              viewMode={viewMode}
              onVisibleRangeChange={(rangeStart, rangeEnd) =>
                setVisibleRange((previousRange) =>
                  previousRange.startDate === rangeStart && previousRange.endDate === rangeEnd
                    ? previousRange
                    : { startDate: rangeStart, endDate: rangeEnd }
                )
              }
              onCreateStart={handleCalendarSlotClick}
              onReschedule={handleRescheduleAppointment}
              onEventClick={handleEventClick}
            />
            {rangeQueriesLoading && effectiveSelectedStaffIds.length > 0 && (
              <div className="absolute inset-0 z-10 flex items-center justify-center rounded-lg bg-background/40 backdrop-blur-sm">
                <Spinner className="size-7 text-primary" />
              </div>
            )}
          </div>

          {focusedStaffId && (
            <div className="flex flex-col gap-3">
              <h2 className="font-display text-xl font-semibold text-foreground">{t("unavailability.blockedPeriodsTitle")}</h2>
              {isUnavailabilityLoading ? (
                <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-4" role="status" aria-label={t("loading.unavailability")}>
                  {Array.from({ length: 3 }).map((_, skeletonIndex) => (
                    <Skeleton key={skeletonIndex} className="h-28 w-full rounded-lg" />
                  ))}
                </div>
              ) : isUnavailabilityError ? (
                <Alert variant="destructive">
                  <TriangleAlert aria-hidden="true" />
                  <AlertTitle>{t("errors.loadUnavailabilityTitle")}</AlertTitle>
                  <AlertDescription>
                    <Button type="button" variant="outline" size="sm" onClick={() => refetchUnavailability()}>
                      {t("retry")}
                    </Button>
                  </AlertDescription>
                </Alert>
              ) : upcomingUnavailabilityWindows.length === 0 ? (
                <EmptyState icon={CalendarOff} title={t("emptyState.noUnavailability")} className="p-6" />
              ) : (
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
              )}
            </div>
          )}
        </div>
      </div>

      <AppointmentModal
        isOpen={isCreateModalOpen}
        onClose={() => setIsCreateModalOpen(false)}
        onSubmit={(payload) => handleCreateAppointment(payload as CreateAppointmentPayload)}
        isPending={createAppointmentMutation.isPending}
        defaultStaffId={focusedStaffId}
        defaultDate={createDefaults?.date ?? todayDate()}
        defaultStartTime={createDefaults?.startTime}
      />

      <AppointmentDetailsModal
        isOpen={appointmentToView !== null}
        appointment={appointmentToView}
        onClose={() => setAppointmentToView(null)}
        onEdit={handleEditFromDetails}
        onCancel={handleCancelFromDetails}
      />

      <AppointmentModal
        isOpen={appointmentToEdit !== null}
        onClose={() => setAppointmentToEdit(null)}
        onSubmit={handleUpdateAppointment}
        isPending={updateAppointmentMutation.isPending}
        appointment={appointmentToEdit}
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
        staffName={staffMembers.find((staffMember) => staffMember.id === focusedStaffId)?.fullName}
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
    </PageContainer>
  )
}
