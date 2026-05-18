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

const INITIAL_PATIENTS: MockPatient[] = [
  {
    patient_id: "pat-1",
    fhir_resource_id: "fhir-pat-1",
    full_name: "Guilherme de Souza Araujo",
    birth_date: "1988-04-12",
    document_id: "123.456.789-00",
    phone_number: "(11) 98765-4321",
  },
  {
    patient_id: "pat-2",
    fhir_resource_id: "fhir-pat-2",
    full_name: "Mariana Costa Silva",
    birth_date: "1995-11-23",
    document_id: "987.654.321-11",
    phone_number: "(21) 99999-8888",
  },
]

const INITIAL_ENCOUNTERS: MockEncounter[] = [
  {
    fhir_id: "enc-1",
    patient_fhir_id: "fhir-pat-1",
    status: "finished",
    reason_display: "Consulta de Rotina Geral",
    created_at: "2026-05-10T10:00:00Z",
  },
  {
    fhir_id: "enc-2",
    patient_fhir_id: "fhir-pat-1",
    status: "finished",
    reason_display: "Retorno Cardiológico",
    created_at: "2026-05-15T14:30:00Z",
  },
]

const INITIAL_OBSERVATIONS: MockObservation[] = [
  {
    fhir_id: "obs-1",
    encounter_fhir_id: "enc-1",
    patient_fhir_id: "fhir-pat-1",
    loinc_code: "8867-4",
    code_display: "Frequência Cardíaca",
    value_quantity: 72,
    value_unit: "bpm",
    created_at: "2026-05-10T10:05:00Z",
  },
  {
    fhir_id: "obs-2",
    encounter_fhir_id: "enc-1",
    patient_fhir_id: "fhir-pat-1",
    loinc_code: "85354-9",
    code_display: "Pressão Arterial Sistólica",
    value_quantity: 120,
    value_unit: "mmHg",
    created_at: "2026-05-10T10:05:00Z",
  },
  {
    fhir_id: "obs-3",
    encounter_fhir_id: "enc-1",
    patient_fhir_id: "fhir-pat-1",
    loinc_code: "8310-5",
    code_display: "Temperatura Corporal",
    value_quantity: 36.5,
    value_unit: "°C",
    created_at: "2026-05-10T10:05:00Z",
  },
  {
    fhir_id: "obs-4",
    encounter_fhir_id: "enc-2",
    patient_fhir_id: "fhir-pat-1",
    loinc_code: "8867-4",
    code_display: "Frequência Cardíaca",
    value_quantity: 85,
    value_unit: "bpm",
    created_at: "2026-05-15T14:35:00Z",
  },
  {
    fhir_id: "obs-5",
    encounter_fhir_id: "enc-2",
    patient_fhir_id: "fhir-pat-1",
    loinc_code: "85354-9",
    code_display: "Pressão Arterial Sistólica",
    value_quantity: 135,
    value_unit: "mmHg",
    created_at: "2026-05-15T14:35:00Z",
  },
]

const INITIAL_CONDITIONS: MockCondition[] = [
  {
    fhir_id: "cond-1",
    patient_fhir_id: "fhir-pat-1",
    icd10_code: "I10",
    code_display: "Hipertensão Essencial Primária",
    clinical_status: "active",
    created_at: "2026-05-15T14:40:00Z",
  },
]

const INITIAL_REPORTS: MockDiagnosticReport[] = [
  {
    fhir_id: "rep-1",
    encounter_fhir_id: "enc-2",
    patient_fhir_id: "fhir-pat-1",
    report_display: "Eletrocardiograma de Repouso",
    status: "final",
    conclusion: "Ritmo sinusal com leve taquicardia. Recomenda-se acompanhamento ambulatorial.",
    created_at: "2026-05-15T14:45:00Z",
  },
]

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
    return getStorageItem<MockPatient[]>("healthcare_patients", INITIAL_PATIENTS)
  },

  getPatient: async (patientFhirId: string): Promise<MockPatient | null> => {
    const activePatients = getStorageItem<MockPatient[]>("healthcare_patients", INITIAL_PATIENTS)
    return activePatients.find((patient) => patient.fhir_resource_id === patientFhirId) || null
  },

  createPatient: async (patientData: Omit<MockPatient, "patient_id" | "fhir_resource_id">): Promise<MockPatient> => {
    const activePatients = getStorageItem<MockPatient[]>("healthcare_patients", INITIAL_PATIENTS)
    const incrementalId = activePatients.length + 1
    const newPatient: MockPatient = {
      patient_id: `pat-${incrementalId}`,
      fhir_resource_id: `fhir-pat-${incrementalId}`,
      ...patientData,
    }
    activePatients.push(newPatient)
    setStorageItem("healthcare_patients", activePatients)
    return newPatient
  },

  getEncounters: async (patientFhirId: string): Promise<MockEncounter[]> => {
    const activeEncounters = getStorageItem<MockEncounter[]>("healthcare_encounters", INITIAL_ENCOUNTERS)
    return activeEncounters.filter((encounter) => encounter.patient_fhir_id === patientFhirId)
  },

  createEncounter: async (encounterData: Omit<MockEncounter, "fhir_id" | "created_at" | "status">): Promise<MockEncounter> => {
    const activeEncounters = getStorageItem<MockEncounter[]>("healthcare_encounters", INITIAL_ENCOUNTERS)
    const incrementalId = activeEncounters.length + 1
    const newEncounter: MockEncounter = {
      fhir_id: `enc-${incrementalId}`,
      status: "finished",
      created_at: new Date().toISOString(),
      ...encounterData,
    }
    activeEncounters.push(newEncounter)
    setStorageItem("healthcare_encounters", activeEncounters)
    return newEncounter
  },

  getObservations: async (encounterFhirId: string): Promise<MockObservation[]> => {
    const activeObservations = getStorageItem<MockObservation[]>("healthcare_observations", INITIAL_OBSERVATIONS)
    return activeObservations.filter((observation) => observation.encounter_fhir_id === encounterFhirId)
  },

  getAllPatientObservations: async (patientFhirId: string): Promise<MockObservation[]> => {
    const activeObservations = getStorageItem<MockObservation[]>("healthcare_observations", INITIAL_OBSERVATIONS)
    return activeObservations.filter((observation) => observation.patient_fhir_id === patientFhirId)
  },

  createObservation: async (observationData: Omit<MockObservation, "fhir_id" | "created_at">): Promise<MockObservation> => {
    const activeObservations = getStorageItem<MockObservation[]>("healthcare_observations", INITIAL_OBSERVATIONS)
    const incrementalId = activeObservations.length + 1
    const newObservation: MockObservation = {
      fhir_id: `obs-${incrementalId}`,
      created_at: new Date().toISOString(),
      ...observationData,
    }
    activeObservations.push(newObservation)
    setStorageItem("healthcare_observations", activeObservations)
    return newObservation
  },

  getConditions: async (patientFhirId: string): Promise<MockCondition[]> => {
    const activeConditions = getStorageItem<MockCondition[]>("healthcare_conditions", INITIAL_CONDITIONS)
    return activeConditions.filter((condition) => condition.patient_fhir_id === patientFhirId)
  },

  createCondition: async (conditionData: Omit<MockCondition, "fhir_id" | "created_at">): Promise<MockCondition> => {
    const activeConditions = getStorageItem<MockCondition[]>("healthcare_conditions", INITIAL_CONDITIONS)
    const incrementalId = activeConditions.length + 1
    const newCondition: MockCondition = {
      fhir_id: `cond-${incrementalId}`,
      created_at: new Date().toISOString(),
      ...conditionData,
    }
    activeConditions.push(newCondition)
    setStorageItem("healthcare_conditions", activeConditions)
    return newCondition
  },

  getDiagnosticReports: async (encounterFhirId: string): Promise<MockDiagnosticReport[]> => {
    const activeReports = getStorageItem<MockDiagnosticReport[]>("healthcare_reports", INITIAL_REPORTS)
    return activeReports.filter((report) => report.encounter_fhir_id === encounterFhirId)
  },

  createDiagnosticReport: async (reportData: Omit<MockDiagnosticReport, "fhir_id" | "created_at" | "status">): Promise<MockDiagnosticReport> => {
    const activeReports = getStorageItem<MockDiagnosticReport[]>("healthcare_reports", INITIAL_REPORTS)
    const incrementalId = activeReports.length + 1
    const newReport: MockDiagnosticReport = {
      fhir_id: `rep-${incrementalId}`,
      status: "final",
      created_at: new Date().toISOString(),
      ...reportData,
    }
    activeReports.push(newReport)
    setStorageItem("healthcare_reports", activeReports)
    return newReport
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

