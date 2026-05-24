import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query"
import { patientsApi } from "./api"

export const patientQueryKeys = {
  all: ["patients"] as const,
  lists: () => [...patientQueryKeys.all, "list"] as const,
  detail: (patientFhirId: string) => [...patientQueryKeys.all, "detail", patientFhirId] as const,
  encounters: (patientFhirId: string) => [...patientQueryKeys.all, "encounters", patientFhirId] as const,
  observations: (encounterFhirId: string) => [...patientQueryKeys.all, "observations", encounterFhirId] as const,
  reports: (encounterFhirId: string) => [...patientQueryKeys.all, "reports", encounterFhirId] as const,
  medications: (encounterFhirId: string) => [...patientQueryKeys.all, "medications", encounterFhirId] as const,
}

export const usePatientsQuery = (searchQueryValue?: string) => {
  return useQuery({
    queryKey: [...patientQueryKeys.lists(), searchQueryValue],
    queryFn: async () => {
      const allPatients = await patientsApi.getPatients()
      if (!searchQueryValue) {
        return allPatients
      }
      const lowerSearch = searchQueryValue.toLowerCase()
      return allPatients.filter((patientItem) => {
        return (
          patientItem.full_name.toLowerCase().includes(lowerSearch) ||
          patientItem.document_id.includes(lowerSearch) ||
          patientItem.phone_number.includes(lowerSearch)
        )
      })
    },
  })
}

export const usePatientQuery = (patientId: string) => {
  return useQuery({
    queryKey: patientQueryKeys.detail(patientId),
    queryFn: () => patientsApi.getPatient(patientId),
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
    }) => patientsApi.createPatient(payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: patientQueryKeys.lists() })
    },
  })
}

export const useEncountersQuery = (patientId: string) => {
  return useQuery({
    queryKey: patientQueryKeys.encounters(patientId),
    queryFn: () => patientsApi.getEncounters(patientId),
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
    }) => patientsApi.createEncounter(payload),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: patientQueryKeys.encounters(variables.patient_fhir_id) })
    },
  })
}

export const useObservationsQuery = (encounterId: string) => {
  return useQuery({
    queryKey: patientQueryKeys.observations(encounterId),
    queryFn: () => patientsApi.getObservations(encounterId),
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
    }) => patientsApi.createObservation(payload),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: patientQueryKeys.observations(variables.encounter_fhir_id) })
    },
  })
}

export const useDiagnosticReportsQuery = (encounterId: string) => {
  return useQuery({
    queryKey: patientQueryKeys.reports(encounterId),
    queryFn: () => patientsApi.getDiagnosticReports(encounterId),
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
    }) => patientsApi.createDiagnosticReport(payload),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: patientQueryKeys.reports(variables.encounter_fhir_id) })
    },
  })
}

export const usePatientConditionsQuery = (patientFhirId: string) => {
  return useQuery({
    queryKey: [...patientQueryKeys.all, "conditions", patientFhirId],
    queryFn: () => patientsApi.getConditions(patientFhirId),
    enabled: !!patientFhirId,
  })
}

export const useCreateConditionMutation = () => {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (payload: {
      patient_fhir_id: string
      icd10_code: string
      code_display: string
    }) => patientsApi.createCondition(payload),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({
        queryKey: [...patientQueryKeys.all, "conditions", variables.patient_fhir_id],
      })
    },
  })
}

export const usePatientAllergiesQuery = (patientFhirId: string) => {
  return useQuery({
    queryKey: [...patientQueryKeys.all, "allergies", patientFhirId],
    queryFn: () => patientsApi.getAllergies(patientFhirId),
    enabled: !!patientFhirId,
  })
}

export const useCreateAllergyMutation = () => {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (payload: {
      patient_fhir_id: string
      allergen_code: string
      allergen_display: string
      reaction: string
    }) => patientsApi.createAllergy(payload),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({
        queryKey: [...patientQueryKeys.all, "allergies", variables.patient_fhir_id],
      })
    },
  })
}

export const useMedicationsQuery = (encounterId: string) => {
  return useQuery({
    queryKey: patientQueryKeys.medications(encounterId),
    queryFn: () => patientsApi.getMedications(encounterId),
    enabled: !!encounterId,
  })
}

export const useCreateMedicationMutation = () => {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (payload: {
      encounter_fhir_id: string
      patient_fhir_id: string
      medication_display: string
      dosage_instruction: string
    }) => patientsApi.createMedication(payload),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: patientQueryKeys.medications(variables.encounter_fhir_id) })
    },
  })
}

export type { DiagnosticReport, Encounter, Observation, Patient, AllergyIntolerance, MedicationRequest } from "./types"
