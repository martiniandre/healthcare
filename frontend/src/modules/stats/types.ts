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
  descriptionKey: string
  categoryKey: string
  activeCases: number
  trend: string
}

export interface StatsResponse {
  totalRegisteredPatients: number
  fhirComplianceRate: number
  averageServiceDurationMinutes: number
  activeConsultationsTotal: number
  totalStudiesCount: number
  examModalitiesData: ModalityData[]
  consultationsWeeklyData: ConsultationsDayData[]
  pathologies: PathologyData[]
}
