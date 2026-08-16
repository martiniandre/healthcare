import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query"
import { scheduleApi } from "./api"
import type { CreateAppointmentPayload, CreateUnavailabilityPayload } from "./types"

export const scheduleQueryKeys = {
  all: ["schedule"] as const,
  appointments: () => [...scheduleQueryKeys.all, "appointments"] as const,
  byStaff: (staffId: string, date: string) =>
    [...scheduleQueryKeys.appointments(), "staff", staffId, date] as const,
  byPatient: (patientFhirId: string) =>
    [...scheduleQueryKeys.appointments(), "patient", patientFhirId] as const,
  mine: () => [...scheduleQueryKeys.appointments(), "mine"] as const,
  unavailability: () => [...scheduleQueryKeys.all, "unavailability"] as const,
  unavailabilityByStaff: (staffId: string) =>
    [...scheduleQueryKeys.unavailability(), "staff", staffId] as const,
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

export const useStaffUnavailabilityQuery = (staffId: string) => {
  return useQuery({
    queryKey: scheduleQueryKeys.unavailabilityByStaff(staffId),
    queryFn: () => scheduleApi.listUnavailabilityByStaff(staffId),
    enabled: !!staffId,
  })
}

export const useCreateUnavailabilityMutation = () => {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (payload: CreateUnavailabilityPayload) => scheduleApi.createUnavailability(payload),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({
        queryKey: scheduleQueryKeys.unavailabilityByStaff(variables.staff_id),
      })
    },
  })
}

export const useDeleteUnavailabilityMutation = () => {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (unavailabilityId: string) => scheduleApi.deleteUnavailability(unavailabilityId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: scheduleQueryKeys.unavailability() })
    },
  })
}
