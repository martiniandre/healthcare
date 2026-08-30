import { useMemo } from "react"
import FullCalendar from "@fullcalendar/react"
import dayGridPlugin from "@fullcalendar/daygrid"
import interactionPlugin from "@fullcalendar/interaction"
import type { EventContentArg } from "@fullcalendar/core"
import { ScheduleEventChip } from "./ScheduleEventChip"
import type { CalendarEventShape } from "../schedule_calendar_helpers"

interface ScheduleMonthCalendarProps {
  events: CalendarEventShape[]
  onVisibleRangeChange: (rangeStart: string, rangeEnd: string) => void
}

const toIsoDate = (dateValue: Date): string => {
  const year = dateValue.getFullYear()
  const month = String(dateValue.getMonth() + 1).padStart(2, "0")
  const day = String(dateValue.getDate()).padStart(2, "0")
  return `${year}-${month}-${day}`
}

export const ScheduleMonthCalendar = ({ events, onVisibleRangeChange }: ScheduleMonthCalendarProps) => {
  const renderEventContent = useMemo(() => {
    return (eventContentArgument: EventContentArg) => {
      const extendedProps = eventContentArgument.event.extendedProps as unknown as CalendarEventShape["extendedProps"]
      return <ScheduleEventChip appointment={extendedProps.appointment} staffColor={extendedProps.staffColor} />
    }
  }, [])

  return (
    <div className="bg-card border border-border rounded-xl p-4 schedule-calendar-container">
      <FullCalendar
        plugins={[dayGridPlugin, interactionPlugin]}
        initialView="dayGridMonth"
        locale="pt-br"
        firstDay={1}
        height="auto"
        headerToolbar={{ left: "prev,next today", center: "title", right: "" }}
        dayMaxEvents={3}
        events={events}
        eventContent={renderEventContent}
        eventDisplay="block"
        datesSet={(dateSetInput) => {
          onVisibleRangeChange(toIsoDate(dateSetInput.start), toIsoDate(dateSetInput.end))
        }}
      />
    </div>
  )
}
