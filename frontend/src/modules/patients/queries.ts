import { useQuery, useMutation, useQueryClient, useInfiniteQuery } from "@tanstack/react-query"
import { useCallback } from "react"
import { patientsApi } from "./api"
import { useDebounce } from "../../shared/hooks/useDebounce"
import type { MedicationSafetyAlert, CreateMedicationPayload } from "./types"
import type { NewVitalSignsPanelFormData } from "./patient_schemas"
import { resolveTimelineRange, TIMELINE_PAGE_SIZE, type TimelineFilterState } from "./components/timeline/timelineModel"

export const patientQueryKeys = {
  all: ["patients"] as const,
  lists: () => [...patientQueryKeys.all, "list"] as const,
  detail: (patientFhirId: string) => [...patientQueryKeys.all, "detail", patientFhirId] as const,
  encounters: (patientFhirId: string) => [...patientQueryKeys.all, "encounters", patientFhirId] as const,
  observations: (encounterFhirId: string) => [...patientQueryKeys.all, "observations", encounterFhirId] as const,
  reports: (encounterFhirId: string) => [...patientQueryKeys.all, "reports", encounterFhirId] as const,
  reportVersions: (reportFhirId: string) => [...patientQueryKeys.all, "reportVersions", reportFhirId] as const,
  timeline: (patientFhirId: string, filters: TimelineFilterState) => [...patientQueryKeys.all, "timeline", patientFhirId, filters] as const,
  medications: (encounterFhirId: string) => [...patientQueryKeys.all, "medications", encounterFhirId] as const,
  medicationSafety: (patientFhirId: string, medicationName: string) => [...patientQueryKeys.all, "medicationSafety", patientFhirId, medicationName] as const,
}

export const usePatientsQuery = (
  searchQueryValue?: string,
  sortField?: string,
  sortDirection?: string,
  page?: number,
  limit?: number
) => {
  return useQuery({
    queryKey: [...patientQueryKeys.lists(), { searchQueryValue, sortField, sortDirection, page, limit }],
    queryFn: () => patientsApi.getPatients(searchQueryValue, sortField, sortDirection, page, limit),
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
      reason_code?: string
      reason_display: string
      practitioner_id: string
    }) => patientsApi.createEncounter(payload),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: patientQueryKeys.encounters(variables.patient_fhir_id) })
    },
  })
}

export const useUpdateEncounterMutation = () => {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (payload: {
      encounter_fhir_id: string
      patient_fhir_id: string
      status: "finished" | "cancelled"
    }) => patientsApi.updateEncounter(payload),
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

export const useCreateVitalSignsPanelMutation = () => {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (payload: {
      encounter_fhir_id: string
      patient_fhir_id: string
      panel_form_data: NewVitalSignsPanelFormData
    }) => patientsApi.createVitalSignsBatch(payload.encounter_fhir_id, payload.patient_fhir_id, payload.panel_form_data),
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
      report_code: string
      report_display: string
      conclusion: string
    }) => patientsApi.createDiagnosticReport(payload),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: patientQueryKeys.reports(variables.encounter_fhir_id) })
    },
  })
}

export const useDiagnosticReportVersionsQuery = (reportFhirId: string) => {
  return useQuery({
    queryKey: patientQueryKeys.reportVersions(reportFhirId),
    queryFn: () => patientsApi.getDiagnosticReportVersions(reportFhirId),
    enabled: !!reportFhirId,
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
    mutationFn: (payload: CreateMedicationPayload) => patientsApi.createMedication(payload),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: patientQueryKeys.medications(variables.encounter_fhir_id) })
    },
  })
}

export type MedicationSafetyStatus = "idle" | "loading" | "error" | "ready"

export type MedicationSafetyVerification = { verified: true; alerts: MedicationSafetyAlert[] } | { verified: false }

const MEDICATION_SAFETY_MIN_LENGTH = 3
const MEDICATION_SAFETY_DEBOUNCE_MS = 500

export const useMedicationSafetyCheck = (patientId: string, medicationName: string, enabled = true) => {
  const queryClient = useQueryClient()
  const trimmedName = medicationName.trim()
  const debouncedName = useDebounce(trimmedName, MEDICATION_SAFETY_DEBOUNCE_MS)
  const isCheckable = enabled && !!patientId && debouncedName.length >= MEDICATION_SAFETY_MIN_LENGTH

  const query = useQuery({
    queryKey: patientQueryKeys.medicationSafety(patientId, debouncedName),
    queryFn: () => patientsApi.checkMedicationSafety(patientId, { medicationCode: "", medicationName: debouncedName }),
    enabled: isCheckable,
    retry: false,
    staleTime: 0,
    gcTime: 0,
  })

  const isDebouncing = enabled && trimmedName !== debouncedName && trimmedName.length >= MEDICATION_SAFETY_MIN_LENGTH

  let status: MedicationSafetyStatus = "idle"
  if (isCheckable) {
    if (isDebouncing || query.isPending || query.isFetching) status = "loading"
    else if (query.isError) status = "error"
    else status = "ready"
  } else if (isDebouncing) {
    status = "loading"
  }

  const alerts = status === "ready" ? query.data?.alerts ?? [] : []

  const verifyNow = useCallback(async (): Promise<MedicationSafetyVerification> => {
    try {
      const response = await queryClient.fetchQuery({
        queryKey: patientQueryKeys.medicationSafety(patientId, trimmedName),
        queryFn: () => patientsApi.checkMedicationSafety(patientId, { medicationCode: "", medicationName: trimmedName }),
        staleTime: 0,
      })
      return { verified: true, alerts: response.alerts }
    } catch {
      return { verified: false }
    }
  }, [queryClient, patientId, trimmedName])

  const retry = useCallback(() => {
    void query.refetch()
  }, [query])

  return { status, alerts, verifyNow, retry }
}

export type { DiagnosticReport, DiagnosticReportVersion, Encounter, Observation, Patient, PatientsPage, AllergyIntolerance, MedicationRequest } from "./types"

export const useTimelineQuery = (patientId: string, filters: TimelineFilterState, enabled = true) => {
  return useInfiniteQuery({
    queryKey: patientQueryKeys.timeline(patientId, filters),
    queryFn: ({ pageParam }) =>
      patientsApi.getTimeline(patientId, {
        ...resolveTimelineRange(filters, new Date()),
        types: filters.types,
        cursor: pageParam,
        limit: TIMELINE_PAGE_SIZE,
      }),
    initialPageParam: undefined as string | undefined,
    getNextPageParam: (lastPage) => lastPage.nextCursor ?? undefined,
    enabled: !!patientId && enabled,
  })
}
