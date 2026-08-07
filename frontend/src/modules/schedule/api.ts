import { http } from "../../shared/utils/http"
import type { Appointment, CreateAppointmentPayload } from "./types"

export const scheduleApi = {
  createAppointment: async (payload: CreateAppointmentPayload): Promise<Appointment> => {
    return http.post<Appointment>("/appointments", payload)
  },

  listByStaffOnDate: async (staffId: string, date: string): Promise<Appointment[]> => {
    const queryParameters = new URLSearchParams({ staff_id: staffId, date })
    return http.get<Appointment[]>(`/appointments?${queryParameters.toString()}`)
  },

  listByPatient: async (patientFhirId: string): Promise<Appointment[]> => {
    const queryParameters = new URLSearchParams({ patient_fhir_id: patientFhirId })
    return http.get<Appointment[]>(`/appointments?${queryParameters.toString()}`)
  },

  listMyAppointments: async (): Promise<Appointment[]> => {
    return http.get<Appointment[]>("/appointments/my")
  },

  getAppointment: async (appointmentId: string): Promise<Appointment> => {
    return http.get<Appointment>(`/appointments/${appointmentId}`)
  },

  cancelAppointment: async (appointmentId: string): Promise<Appointment> => {
    return http.post<Appointment>(`/appointments/${appointmentId}/cancel`)
  },
}
