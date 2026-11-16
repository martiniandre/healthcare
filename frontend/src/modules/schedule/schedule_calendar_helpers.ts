import type { Appointment } from "./types"

export interface CalendarEventShape {
  id: string
  title: string
  start: string
  end: string
  backgroundColor?: string
  borderColor?: string
  extendedProps: {
    appointment: Appointment
    staffColor: string
  }
}

export const staffColorPalette = [
  "#4a6e52",
  "#9a5a3a",
  "#2f6f73",
  "#7a4f7d",
  "#9a6b12",
  "#a8475a",
  "#4c5f8a",
  "#6b6a2a",
]

export const staffColorForIndex = (staffIndex: number): string => {
  return staffColorPalette[staffIndex % staffColorPalette.length]
}

export const activeAppointmentsOnly = (appointments: Appointment[]): Appointment[] => {
  return appointments.filter((appointment) => appointment.status !== "cancelled")
}

export const appointmentsToCalendarEvents = (
  appointments: Appointment[],
  staffColor: string
): CalendarEventShape[] => {
  return activeAppointmentsOnly(appointments).map((appointment) => ({
    id: appointment.id,
    title: appointment.reason || "",
    start: appointment.starts_at,
    end: appointment.ends_at,
    extendedProps: {
      appointment,
      staffColor,
    },
  }))
}

export interface VisibleSlotRange {
  slotMinTime: string
  slotMaxTime: string
}

const DEFAULT_FIRST_VISIBLE_HOUR = 6
const DEFAULT_LAST_VISIBLE_HOUR = 22
const HOURS_IN_DAY = 24

const formatSlotHour = (hour: number): string => `${String(hour).padStart(2, "0")}:00:00`

export const resolveVisibleSlotRange = (events: Pick<CalendarEventShape, "start" | "end">[]): VisibleSlotRange => {
  let firstVisibleHour = DEFAULT_FIRST_VISIBLE_HOUR
  let lastVisibleHour = DEFAULT_LAST_VISIBLE_HOUR

  for (const calendarEvent of events) {
    const startDate = new Date(calendarEvent.start)
    const endDate = new Date(calendarEvent.end)
    if (Number.isNaN(startDate.getTime()) || Number.isNaN(endDate.getTime())) continue

    firstVisibleHour = Math.min(firstVisibleHour, startDate.getHours())

    const crossesMidnight = endDate.toDateString() !== startDate.toDateString()
    const endHourCeiling = crossesMidnight
      ? HOURS_IN_DAY
      : endDate.getHours() + (endDate.getMinutes() > 0 ? 1 : 0)
    lastVisibleHour = Math.max(lastVisibleHour, endHourCeiling)
  }

  return {
    slotMinTime: formatSlotHour(firstVisibleHour),
    slotMaxTime: formatSlotHour(Math.min(lastVisibleHour, HOURS_IN_DAY)),
  }
}
