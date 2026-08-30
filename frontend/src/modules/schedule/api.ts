import { http } from "../../shared/utils/http"
import type {
  Appointment,
  CreateAppointmentPayload,
  RescheduleAppointmentPayload,
  StaffUnavailability,
  CreateUnavailabilityPayload,
} from "./types"

export const scheduleApi = {
  createAppointment: async (payload: CreateAppointmentPayload): Promise<Appointment> => {
    return http.post<Appointment>("/appointments", payload)
  },

  listByStaffOnDate: async (staffId: string, date: string): Promise<Appointment[]> => {
    const queryParameters = new URLSearchParams({ staff_id: staffId, date })
    return http.get<Appointment[]>(`/appointments?${queryParameters.toString()}`)
  },

  listByStaffInRange: async (staffId: string, startDate: string, endDate: string): Promise<Appointment[]> => {
    const queryParameters = new URLSearchParams({ staff_id: staffId, start_date: startDate, end_date: endDate })
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

  rescheduleAppointment: async (appointmentId: string, payload: RescheduleAppointmentPayload): Promise<Appointment> => {
    return http.put<Appointment>(`/appointments/${appointmentId}`, payload)
  },

  createUnavailability: async (payload: CreateUnavailabilityPayload): Promise<StaffUnavailability> => {
    return http.post<StaffUnavailability>("/schedule/unavailability", payload)
  },

  listUnavailabilityByStaff: async (staffId: string): Promise<StaffUnavailability[]> => {
    const queryParameters = new URLSearchParams({ staff_id: staffId })
    return http.get<StaffUnavailability[]>(`/schedule/unavailability?${queryParameters.toString()}`)
  },

  deleteUnavailability: async (unavailabilityId: string): Promise<StaffUnavailability> => {
    return http.delete<StaffUnavailability>(`/schedule/unavailability/${unavailabilityId}`)
  },
}
