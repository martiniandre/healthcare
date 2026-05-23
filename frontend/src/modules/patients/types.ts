export interface Patient {
  patient_id: string
  fhir_resource_id: string
  full_name: string
  birth_date: string
  document_id: string
  phone_number: string
}

export interface Encounter {
  fhir_id: string
  patient_fhir_id: string
  status: string
  reason_display: string
  practitioner_id?: string
  created_at: string
}

export interface Observation {
  fhir_id: string
  encounter_fhir_id: string
  patient_fhir_id: string
  loinc_code: string
  code_display: string
  value_quantity: number
  value_unit: string
  created_at: string
}

export interface Condition {
  fhir_id: string
  patient_fhir_id: string
  icd10_code: string
  code_display: string
  clinical_status: string
  created_at: string
}

export interface DiagnosticReport {
  fhir_id: string
  encounter_fhir_id: string
  patient_fhir_id: string
  report_display: string
  status: string
  conclusion: string
  created_at: string
}
