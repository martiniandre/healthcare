export interface ModalityData {
  modality: string
  percentage: number
  count: number
  color: string
}

export interface ConsultationsDayData {
  dayName: string
  count: number
}

export interface PathologyData {
  code: string
  description: string
  category: string
  activeCases: number
  trend: string
}

export interface StatsResponse {
  total_patients: number
  fhir_compliance_rate: number
  avg_service_duration_minutes: number
  weekly_consultations: ConsultationsDayData[]
  exam_modalities: ModalityData[]
  pathology_cases: PathologyData[]
}
