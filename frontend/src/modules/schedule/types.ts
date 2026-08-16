export type AppointmentStatus = "scheduled" | "confirmed" | "cancelled" | "finished"

export interface Appointment {
  id: string
  patient_fhir_id: string
  staff_id: string
  starts_at: string
  ends_at: string
  status: AppointmentStatus
  reason: string
  version: number
  created_at: string
}

export interface CreateAppointmentPayload {
  patient_fhir_id: string
  staff_id: string
  starts_at: string
  ends_at: string
  reason?: string
  idempotency_key: string
}

export const AppointmentStatusLabel = {
  scheduled: "scheduled",
  confirmed: "confirmed",
  cancelled: "cancelled",
  finished: "finished",
} as const

export interface StaffUnavailability {
  id: string
  staff_id: string
  starts_at: string
  ends_at: string
  reason: string
  created_at: string
}

export interface CreateUnavailabilityPayload {
  staff_id: string
  starts_at: string
  ends_at: string
  reason: string
}
