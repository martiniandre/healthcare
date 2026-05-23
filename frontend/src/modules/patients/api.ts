import { http } from "../../shared/utils/http"
import type { DiagnosticReport, Encounter, Observation, Patient, Condition } from "./types"

export const patientsApi = {
  getPatients: async (): Promise<Patient[]> => {
    return http.get<Patient[]>("/patients")
  },

  getPatient: async (patientFhirId: string): Promise<Patient | null> => {
    try {
      return await http.get<Patient>(`/patients/${patientFhirId}`)
    } catch {
      return null
    }
  },

  createPatient: async (patientData: Omit<Patient, "patient_id" | "fhir_resource_id">): Promise<Patient> => {
    const creationResponse = await http.post<{ patient_id: string; fhir_resource_id: string }>("/patients", patientData)
    return {
      patient_id: creationResponse.patient_id,
      fhir_resource_id: creationResponse.fhir_resource_id,
      ...patientData,
    }
  },

  getEncounters: async (patientFhirId: string): Promise<Encounter[]> => {
    return http.get<Encounter[]>(`/patients/${patientFhirId}/encounters`)
  },

  createEncounter: async (encounterData: Omit<Encounter, "fhir_id" | "created_at" | "status">): Promise<Encounter> => {
    return http.post<Encounter>(`/patients/${encounterData.patient_fhir_id}/encounters`, {
      reason_display: encounterData.reason_display,
      practitioner_id: encounterData.practitioner_id,
    })
  },

  getObservations: async (encounterFhirId: string): Promise<Observation[]> => {
    return http.get<Observation[]>(`/encounters/${encounterFhirId}/observations`)
  },

  getAllPatientObservations: async (patientFhirId: string): Promise<Observation[]> => {
    return http.get<Observation[]>(`/patients/${patientFhirId}/observations`)
  },

  createObservation: async (observationData: Omit<Observation, "fhir_id" | "created_at">): Promise<Observation> => {
    return http.post<Observation>(`/encounters/${observationData.encounter_fhir_id}/observations`, {
      patient_fhir_id: observationData.patient_fhir_id,
      loinc_code: observationData.loinc_code,
      code_display: observationData.code_display,
      value_quantity: observationData.value_quantity,
      value_unit: observationData.value_unit,
    })
  },

  getConditions: async (patientFhirId: string): Promise<Condition[]> => {
    return http.get<Condition[]>(`/patients/${patientFhirId}/conditions`)
  },

  createCondition: async (conditionData: Omit<Condition, "fhir_id" | "created_at">): Promise<Condition> => {
    return http.post<Condition>(`/patients/${conditionData.patient_fhir_id}/conditions`, {
      icd10_code: conditionData.icd10_code,
      code_display: conditionData.code_display,
    })
  },

  getDiagnosticReports: async (encounterFhirId: string): Promise<DiagnosticReport[]> => {
    return http.get<DiagnosticReport[]>(`/encounters/${encounterFhirId}/reports`)
  },

  createDiagnosticReport: async (reportData: Omit<DiagnosticReport, "fhir_id" | "created_at" | "status">): Promise<DiagnosticReport> => {
    return http.post<DiagnosticReport>(`/encounters/${reportData.encounter_fhir_id}/reports`, {
      patient_fhir_id: reportData.patient_fhir_id,
      report_display: reportData.report_display,
      conclusion: reportData.conclusion,
    })
  },
}
