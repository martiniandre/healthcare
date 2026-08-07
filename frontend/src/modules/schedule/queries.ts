import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query"
import { scheduleApi } from "./api"
import type { CreateAppointmentPayload } from "./types"

export const scheduleQueryKeys = {
  all: ["schedule"] as const,
  appointments: () => [...scheduleQueryKeys.all, "appointments"] as const,
  byStaff: (staffId: string, date: string) =>
    [...scheduleQueryKeys.appointments(), "staff", staffId, date] as const,
  byPatient: (patientFhirId: string) =>
    [...scheduleQueryKeys.appointments(), "patient", patientFhirId] as const,
  mine: () => [...scheduleQueryKeys.appointments(), "mine"] as const,
}

export const useStaffDayAppointmentsQuery = (staffId: string, date: string) => {
  return useQuery({
    queryKey: scheduleQueryKeys.byStaff(staffId, date),
    queryFn: () => scheduleApi.listByStaffOnDate(staffId, date),
    enabled: !!staffId && !!date,
  })
}

export const usePatientAppointmentsQuery = (patientFhirId: string) => {
  return useQuery({
    queryKey: scheduleQueryKeys.byPatient(patientFhirId),
    queryFn: () => scheduleApi.listByPatient(patientFhirId),
    enabled: !!patientFhirId,
  })
}

export const useMyAppointmentsQuery = () => {
  return useQuery({
    queryKey: scheduleQueryKeys.mine(),
    queryFn: scheduleApi.listMyAppointments,
  })
}

export const useCreateAppointmentMutation = () => {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (payload: CreateAppointmentPayload) => scheduleApi.createAppointment(payload),
    onSuccess: (_, variables) => {
      const bookedDate = variables.starts_at.slice(0, 10)
      queryClient.invalidateQueries({ queryKey: scheduleQueryKeys.byStaff(variables.staff_id, bookedDate) })
      queryClient.invalidateQueries({ queryKey: scheduleQueryKeys.byPatient(variables.patient_fhir_id) })
      queryClient.invalidateQueries({ queryKey: scheduleQueryKeys.mine() })
    },
  })
}

export const useCancelAppointmentMutation = () => {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (appointmentId: string) => scheduleApi.cancelAppointment(appointmentId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: scheduleQueryKeys.appointments() })
    },
  })
}
