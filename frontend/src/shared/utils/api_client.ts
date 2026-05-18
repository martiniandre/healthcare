import { http } from "./http"

interface MockPatient {
  patient_id: string
  fhir_resource_id: string
  full_name: string
  birth_date: string
  document_id: string
  phone_number: string
}

interface MockEncounter {
  fhir_id: string
  patient_fhir_id: string
  status: string
  reason_display: string
  practitioner_id?: string
  created_at: string
}

interface MockObservation {
  fhir_id: string
  encounter_fhir_id: string
  patient_fhir_id: string
  loinc_code: string
  code_display: string
  value_quantity: number
  value_unit: string
  created_at: string
}

interface MockCondition {
  fhir_id: string
  patient_fhir_id: string
  icd10_code: string
  code_display: string
  clinical_status: string
  created_at: string
}

interface MockDiagnosticReport {
  fhir_id: string
  encounter_fhir_id: string
  patient_fhir_id: string
  report_display: string
  status: string
  conclusion: string
  created_at: string
}

interface MockImagingStudy {
  imaging_study_id: string
  patient_fhir_id: string
  title: string
  modality: string
  study_instance_uid: string
  status: string
  created_at: string
}




const INITIAL_STUDIES: MockImagingStudy[] = [
  {
    imaging_study_id: "study-1",
    patient_fhir_id: "fhir-pat-1",
    title: "Tomografia Computadorizada de Tórax",
    modality: "CT",
    study_instance_uid: "1.2.840.10008.5.1.4.1.1.2.20260516.1",
    status: "completed",
    created_at: "2026-05-16T10:00:00Z",
  },
]

function getStorageItem<T>(key: string, defaultValue: T): T {
  const storedItem = localStorage.getItem(key)
  if (!storedItem) {
    localStorage.setItem(key, JSON.stringify(defaultValue))
    return defaultValue
  }
  return JSON.parse(storedItem)
}

function setStorageItem<T>(key: string, value: T): void {
  localStorage.setItem(key, JSON.stringify(value))
}

export const clinicApi = {
  getPatients: async (): Promise<MockPatient[]> => {
    return http.get<MockPatient[]>("/api/patients")
  },

  getPatient: async (patientFhirId: string): Promise<MockPatient | null> => {
    try {
      return await http.get<MockPatient>(`/api/patients/${patientFhirId}`)
    } catch {
      return null
    }
  },

  createPatient: async (patientData: Omit<MockPatient, "patient_id" | "fhir_resource_id">): Promise<MockPatient> => {
    const creationResponse = await http.post<{ patient_id: string; fhir_resource_id: string }>("/api/patients", patientData)
    return {
      patient_id: creationResponse.patient_id,
      fhir_resource_id: creationResponse.fhir_resource_id,
      ...patientData,
    }
  },

  getEncounters: async (patientFhirId: string): Promise<MockEncounter[]> => {
    return http.get<MockEncounter[]>(`/api/patients/${patientFhirId}/encounters`)
  },

  createEncounter: async (encounterData: Omit<MockEncounter, "fhir_id" | "created_at" | "status">): Promise<MockEncounter> => {
    return http.post<MockEncounter>(`/api/patients/${encounterData.patient_fhir_id}/encounters`, {
      reason_display: encounterData.reason_display,
      practitioner_id: encounterData.practitioner_id,
    })
  },

  getObservations: async (encounterFhirId: string): Promise<MockObservation[]> => {
    return http.get<MockObservation[]>(`/api/encounters/${encounterFhirId}/observations`)
  },

  getAllPatientObservations: async (patientFhirId: string): Promise<MockObservation[]> => {
    return http.get<MockObservation[]>(`/api/patients/${patientFhirId}/observations`)
  },

  createObservation: async (observationData: Omit<MockObservation, "fhir_id" | "created_at">): Promise<MockObservation> => {
    return http.post<MockObservation>(`/api/encounters/${observationData.encounter_fhir_id}/observations`, {
      patient_fhir_id: observationData.patient_fhir_id,
      loinc_code: observationData.loinc_code,
      code_display: observationData.code_display,
      value_quantity: observationData.value_quantity,
      value_unit: observationData.value_unit,
    })
  },

  getConditions: async (patientFhirId: string): Promise<MockCondition[]> => {
    return http.get<MockCondition[]>(`/api/patients/${patientFhirId}/conditions`)
  },

  createCondition: async (conditionData: Omit<MockCondition, "fhir_id" | "created_at">): Promise<MockCondition> => {
    return http.post<MockCondition>(`/api/patients/${conditionData.patient_fhir_id}/conditions`, {
      icd10_code: conditionData.icd10_code,
      code_display: conditionData.code_display,
    })
  },

  getDiagnosticReports: async (encounterFhirId: string): Promise<MockDiagnosticReport[]> => {
    return http.get<MockDiagnosticReport[]>(`/api/encounters/${encounterFhirId}/reports`)
  },

  createDiagnosticReport: async (reportData: Omit<MockDiagnosticReport, "fhir_id" | "created_at" | "status">): Promise<MockDiagnosticReport> => {
    return http.post<MockDiagnosticReport>(`/api/encounters/${reportData.encounter_fhir_id}/reports`, {
      patient_fhir_id: reportData.patient_fhir_id,
      report_display: reportData.report_display,
      conclusion: reportData.conclusion,
    })
  },

  getImagingStudies: async (patientFhirId: string): Promise<MockImagingStudy[]> => {
    const activeStudies = getStorageItem<MockImagingStudy[]>("healthcare_studies", INITIAL_STUDIES)
    return activeStudies.filter((study) => study.patient_fhir_id === patientFhirId)
  },

  getImagingStudy: async (imagingStudyId: string): Promise<MockImagingStudy | null> => {
    const activeStudies = getStorageItem<MockImagingStudy[]>("healthcare_studies", INITIAL_STUDIES)
    return activeStudies.find((study) => study.imaging_study_id === imagingStudyId) || null
  },

  createImagingStudy: async (studyData: Omit<MockImagingStudy, "imaging_study_id" | "created_at" | "status" | "study_instance_uid">): Promise<MockImagingStudy> => {
    const activeStudies = getStorageItem<MockImagingStudy[]>("healthcare_studies", INITIAL_STUDIES)
    const incrementalId = activeStudies.length + 1
    const newStudy: MockImagingStudy = {
      imaging_study_id: `study-${incrementalId}`,
      study_instance_uid: `1.2.840.10008.5.1.4.1.1.2.20260516.${incrementalId}`,
      status: "completed",
      created_at: new Date().toISOString(),
      ...studyData,
    }
    activeStudies.push(newStudy)
    setStorageItem("healthcare_studies", activeStudies)
    return newStudy
  },

  login: async (emailValue: string, passwordValue: string): Promise<{ token: string; userId: string; role: string; email: string }> => {
    try {
      return await http.post<{ token: string; userId: string; role: string; email: string }>("/api/auth/login", {
        email: emailValue,
        password: passwordValue,
      })
    } catch (apiError) {
      const axiosError = apiError as { response?: { data?: { error?: string } } }
      const errorMessage = axiosError.response?.data?.error || "Credenciais inválidas."
      throw new Error(errorMessage, { cause: apiError })
    }
  },

  logout: async (): Promise<void> => {
    try {
      await http.post("/api/auth/logout")
    } catch (logoutError) {
      console.warn("Logout endpoint failed:", logoutError)
    }
  },
}

