import { http } from "../../shared/utils/http"
import type { DiagnosticReport, DiagnosticReportVersion, Encounter, Observation, Patient, Condition, CreatePatientResponse, PatientsPage, AllergyIntolerance, MedicationRequest } from "./types"
import type { NewVitalSignsPanelFormData } from "./patient_schemas"

export const patientsApi = {
  getPatients: async (search?: string, sortField?: string, sortDirection?: string, page?: number, limit?: number): Promise<PatientsPage> => {
    const params = new URLSearchParams()
    if (search) params.append("search", search)
    if (sortField) params.append("sortField", sortField)
    if (sortDirection) params.append("sortDirection", sortDirection)
    if (page) params.append("page", page.toString())
    if (limit) params.append("limit", limit.toString())
    const queryString = params.toString()
    return http.get<PatientsPage>(`/patients${queryString ? `?${queryString}` : ""}`)
  },

  getPatient: async (patientFhirId: string): Promise<Patient> => {
    return http.get<Patient>(`/patients/${patientFhirId}`)
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
      reason_code: encounterData.reason_code ?? "",
      reason_display: encounterData.reason_display,
      practitioner_id: encounterData.practitioner_id,
    })
  },

  updateEncounter: async (encounterData: { encounter_fhir_id: string; status: "finished" | "cancelled" }): Promise<Encounter> => {
    return http.put<Encounter>(`/encounters/${encounterData.encounter_fhir_id}`, {
      status: encounterData.status,
    })
  },

  getObservations: async (encounterFhirId: string): Promise<Observation[]> => {
    return http.get<Observation[]>(`/encounters/${encounterFhirId}/observations`)
  },

  getAllPatientObservations: async (patientFhirId: string): Promise<Observation[]> => {
    return http.get<Observation[]>(`/patients/${patientFhirId}/observations`)
  },

  createVitalSignsBatch: async (encounterFhirId: string, patientFhirId: string, panelFormData: NewVitalSignsPanelFormData): Promise<Observation[]> => {
    return http.post<Observation[]>(`/encounters/${encounterFhirId}/observations/batch`, {
      patient_fhir_id: patientFhirId,
      panel: {
        heart_rate: panelFormData.heartRate ?? null,
        body_temperature: panelFormData.bodyTemperature ?? null,
        systolic_blood_pressure: panelFormData.systolicBloodPressure ?? null,
        diastolic_blood_pressure: panelFormData.diastolicBloodPressure ?? null,
        oxygen_saturation: panelFormData.oxygenSaturation ?? null,
        respiratory_rate: panelFormData.respiratoryRate ?? null,
        weight_kg: panelFormData.weightKg ?? null,
        height_cm: panelFormData.heightCm ?? null,
      },
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

  getDiagnosticReportVersions: async (reportFhirId: string): Promise<DiagnosticReportVersion[]> => {
    return http.get<DiagnosticReportVersion[]>(`/reports/${reportFhirId}/versions`)
  },

  createDiagnosticReport: async (reportData: Omit<DiagnosticReport, "fhir_id" | "created_at" | "status">): Promise<DiagnosticReport> => {
    return http.post<DiagnosticReport>(`/encounters/${reportData.encounter_fhir_id}/reports`, {
      patient_fhir_id: reportData.patient_fhir_id,
      report_code: reportData.report_code,
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
