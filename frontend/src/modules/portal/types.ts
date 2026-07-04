export interface PortalDashboard {
  patient_info: PatientInfo
  upcoming_encounters: PortalEncounter[]
  recent_observations: PortalObservation[]
  active_conditions: PortalCondition[]
  active_medications: PortalMedication[]
  recent_reports: PortalReport[]
  recent_imaging: PortalImaging[]
}

export interface PatientInfo {
  fhir_resource_id: string
  full_name: string
  birth_date: string
  document_id: string
}

export interface PortalEncounter {
  fhir_resource_id: string
  status: string
  reason_display: string
  started_at: string
  ended_at?: string
}

export interface PortalObservation {
  fhir_resource_id: string
  code_display: string
  loinc_code: string
  value_quantity: number
  value_unit: string
  observed_at: string
}

export interface PortalCondition {
  fhir_resource_id: string
  code_display: string
  icd10_code: string
  clinical_status: string
  onset_at: string
}

export interface PortalMedication {
  fhir_resource_id: string
  medication_name: string
  dosage_instructions: string
  status: string
  issued_at: string
}

export interface PortalReport {
  fhir_resource_id: string
  report_display: string
  status: string
  conclusion: string
  issued_at: string
}

export interface PortalImaging {
  fhir_resource_id: string
  title: string
  modality: string
  status: string
  created_at: string
}
