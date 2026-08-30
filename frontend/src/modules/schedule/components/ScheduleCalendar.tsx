import { useMemo } from "react"
import FullCalendar from "@fullcalendar/react"
import dayGridPlugin from "@fullcalendar/daygrid"
import timeGridPlugin from "@fullcalendar/timegrid"
import interactionPlugin from "@fullcalendar/interaction"
import multiMonthPlugin from "@fullcalendar/multimonth"
import type { EventContentArg, EventDropArg, DateSelectArg } from "@fullcalendar/core"
import { ScheduleEventChip } from "./ScheduleEventChip"
import type { CalendarEventShape } from "../schedule_calendar_helpers"
import type { Appointment } from "../types"

export type ScheduleViewMode = "week" | "day" | "month" | "year"

const VIEW_TYPE_BY_MODE: Record<ScheduleViewMode, string> = {
  week: "timeGridWeek",
  day: "timeGridDay",
  month: "dayGridMonth",
  year: "multiMonthYear",
}

interface ScheduleCalendarProps {
  events: CalendarEventShape[]
  viewMode: ScheduleViewMode
  onVisibleRangeChange: (rangeStart: string, rangeEnd: string) => void
  onCreateStart: (start: Date) => void
  onReschedule: (appointment: Appointment, newStart: Date, newEnd: Date) => void
}

const toIsoDate = (dateValue: Date): string => {
  const year = dateValue.getFullYear()
  const month = String(dateValue.getMonth() + 1).padStart(2, "0")
  const day = String(dateValue.getDate()).padStart(2, "0")
  return `${year}-${month}-${day}`
}

const headerAccentGradient =
  "linear-gradient(90deg, #2563eb 0%, #3b82f6 55%, #0d9488 100%)"

const isAlignedToHalfHour = (dateValue: Date): boolean => {
  return dateValue.getMinutes() === 0 || dateValue.getMinutes() === 30
}

export const ScheduleCalendar = ({
  events,
  viewMode,
  onVisibleRangeChange,
  onCreateStart,
  onReschedule,
}: ScheduleCalendarProps) => {
  const renderEventContent = useMemo(() => {
    return (eventContentArgument: EventContentArg) => {
      const extendedProps = eventContentArgument.event.extendedProps as unknown as CalendarEventShape["extendedProps"]
      if (!extendedProps?.appointment) {
        return null
      }
      return <ScheduleEventChip appointment={extendedProps.appointment} staffColor={extendedProps.staffColor} />
    }
  }, [])

  const handleSelect = (selection: DateSelectArg) => {
    const startDate = selection.start
    if (selection.allDay) {
      startDate.setHours(9, 0, 0, 0)
    }
    onCreateStart(startDate)
  }

  const handleDateClick = (clickInfo: { date: Date; allDay: boolean }) => {
    const clickedDate = new Date(clickInfo.date)
    if (clickInfo.allDay) {
      clickedDate.setHours(9, 0, 0, 0)
    }
    onCreateStart(clickedDate)
  }

  const handleEventDrop = (dropInfo: EventDropArg) => {
    const extendedProps = dropInfo.event.extendedProps as unknown as CalendarEventShape["extendedProps"]
    onReschedule(extendedProps.appointment, dropInfo.event.start as Date, dropInfo.event.end as Date)
  }

  return (
    <div className="schedule-calendar-container relative overflow-hidden rounded-2xl border border-border/70 bg-white shadow-[0_1px_2px_rgba(15,23,42,0.04),0_18px_40px_-26px_rgba(15,23,42,0.38)]">
      <div className="h-1.5 w-full" style={{ background: headerAccentGradient }} />
      <div className="p-3 sm:p-5">
        <FullCalendar
          plugins={[dayGridPlugin, timeGridPlugin, interactionPlugin, multiMonthPlugin]}
          initialView={VIEW_TYPE_BY_MODE[viewMode]}
          key={viewMode}
          locale="pt-br"
          firstDay={1}
          height="auto"
          headerToolbar={{ left: "prev,next today", center: "title", right: "" }}
          dayMaxEvents={3}
          events={events}
          eventContent={renderEventContent}
          eventDisplay="block"
          slotDuration="00:30:00"
          snapDuration="00:30:00"
          allDaySlot={false}
          selectable={true}
          selectMirror={true}
          selectAllow={(selectionInfo) => selectionInfo.end.getTime() - selectionInfo.start.getTime() <= 60 * 60 * 1000 && isAlignedToHalfHour(selectionInfo.start)}
          select={handleSelect}
          dateClick={handleDateClick}
          editable={true}
          eventDrop={handleEventDrop}
          datesSet={(dateSetInput) => {
            onVisibleRangeChange(toIsoDate(dateSetInput.start), toIsoDate(dateSetInput.end))
          }}
        />
      </div>
    </div>
  )
}
