import { http } from "../../shared/utils/http"
import type { DiagnosticReport, Encounter, Observation, Patient, Condition, CreatePatientResponse, AllergyIntolerance, MedicationRequest } from "./types"

export const patientsApi = {
  getPatients: async (search?: string, sortField?: string, sortDirection?: string, page?: number, limit?: number): Promise<Patient[]> => {
    const params = new URLSearchParams()
    if (search) params.append("search", search)
    if (sortField) params.append("sortField", sortField)
    if (sortDirection) params.append("sortDirection", sortDirection)
    if (page) params.append("page", page.toString())
    if (limit) params.append("limit", limit.toString())
    const queryString = params.toString()
    return http.get<Patient[]>(`/patients${queryString ? `?${queryString}` : ""}`)
  },

  getPatient: async (patientFhirId: string): Promise<Patient | null> => {
    try {
      return await http.get<Patient>(`/patients/${patientFhirId}`)
    } catch {
      return null
    }
  },

  createPatient: async (patientData: Omit<Patient, "patient_id" | "fhir_resource_id">): Promise<Patient> => {
    const creationResponse = await http.post<CreatePatientResponse>("/patients", patientData)
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

  createCondition: async (conditionData: Omit<Condition, "fhir_id" | "created_at" | "clinical_status">): Promise<Condition> => {
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

  getAllergies: async (patientFhirId: string): Promise<AllergyIntolerance[]> => {
    return http.get<AllergyIntolerance[]>(`/patients/${patientFhirId}/allergies`)
  },

  createAllergy: async (allergyData: Omit<AllergyIntolerance, "fhir_id" | "created_at" | "clinical_status">): Promise<AllergyIntolerance> => {
    return http.post<AllergyIntolerance>(`/patients/${allergyData.patient_fhir_id}/allergies`, {
      allergen_code: allergyData.allergen_code,
      allergen_display: allergyData.allergen_display,
      reaction: allergyData.reaction,
    })
  },

  getMedications: async (encounterFhirId: string): Promise<MedicationRequest[]> => {
    return http.get<MedicationRequest[]>(`/encounters/${encounterFhirId}/medications`)
  },

  createMedication: async (medicationData: Omit<MedicationRequest, "fhir_id" | "created_at" | "status">): Promise<MedicationRequest> => {
    return http.post<MedicationRequest>(`/encounters/${medicationData.encounter_fhir_id}/medications`, {
      patient_fhir_id: medicationData.patient_fhir_id,
      medication_name: medicationData.medication_name,
      dosage_instructions: medicationData.dosage_instructions,
    })
  },
}
