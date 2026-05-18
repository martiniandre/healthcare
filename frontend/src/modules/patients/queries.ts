import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query"
import { clinicApi } from "../../shared/utils/api_client"

export const usePatientsQuery = () => {
  return useQuery({
    queryKey: ["patients"],
    queryFn: () => clinicApi.getPatients(),
  })
}

export const usePatientQuery = (patientId: string) => {
  return useQuery({
    queryKey: ["patient", patientId],
    queryFn: () => clinicApi.getPatient(patientId),
    enabled: !!patientId,
  })
}

export const useCreatePatientMutation = () => {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (payload: {
      full_name: string
      birth_date: string
      document_id: string
      phone_number: string
    }) => clinicApi.createPatient(payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["patients"] })
    },
  })
}

export const useEncountersQuery = (patientId: string) => {
  return useQuery({
    queryKey: ["encounters", patientId],
    queryFn: () => clinicApi.getEncounters(patientId),
    enabled: !!patientId,
  })
}

export const useCreateEncounterMutation = () => {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (payload: {
      patient_fhir_id: string
      reason_display: string
      practitioner_id?: string
    }) => clinicApi.createEncounter(payload),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: ["encounters", variables.patient_fhir_id] })
    },
  })
}

export const useObservationsQuery = (encounterId: string) => {
  return useQuery({
    queryKey: ["observations", encounterId],
    queryFn: () => clinicApi.getObservations(encounterId),
    enabled: !!encounterId,
  })
}

export const useCreateObservationMutation = () => {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (payload: {
      encounter_fhir_id: string
      patient_fhir_id: string
      loinc_code: string
      code_display: string
      value_quantity: number
      value_unit: string
    }) => clinicApi.createObservation(payload),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: ["observations", variables.encounter_fhir_id] })
    },
  })
}

export const useDiagnosticReportsQuery = (encounterId: string) => {
  return useQuery({
    queryKey: ["reports", encounterId],
    queryFn: () => clinicApi.getDiagnosticReports(encounterId),
    enabled: !!encounterId,
  })
}

export const useCreateDiagnosticReportMutation = () => {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (payload: {
      encounter_fhir_id: string
      patient_fhir_id: string
      report_display: string
      conclusion: string
    }) => clinicApi.createDiagnosticReport(payload),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: ["reports", variables.encounter_fhir_id] })
    },
  })
}

export const useImagingStudiesQuery = (patientId: string) => {
  return useQuery({
    queryKey: ["studies", patientId],
    queryFn: () => clinicApi.getImagingStudies(patientId),
    enabled: !!patientId,
  })
}
