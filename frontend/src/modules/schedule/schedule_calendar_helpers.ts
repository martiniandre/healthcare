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
  "#2563eb",
  "#0d9488",
  "#7c3aed",
  "#db2777",
  "#ea580c",
  "#16a34a",
  "#ca8a04",
  "#0284c7",
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
    title: appointment.reason || appointment.patient_fhir_id,
    start: appointment.starts_at,
    end: appointment.ends_at,
    extendedProps: {
      appointment,
      staffColor,
    },
  }))
}
