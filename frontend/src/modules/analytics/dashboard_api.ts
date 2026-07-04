import { http } from "../../shared/utils/http"
import type { DashboardData, DoctorConsultation, DiagnosisCount } from "./dashboard_types"

interface OccupancyRateResponse {
  rate: number
  total_beds: number
  occupied_beds: number
}

interface AvgWaitTimeResponse {
  average_minutes: number
  by_department: { department: string; minutes: number }[]
}

export const dashboardApi = {
  getDashboard: async (): Promise<DashboardData> => {
    return http.get<DashboardData>("/analytics/dashboard")
  },

  getConsultationsPerDoctor: async (): Promise<DoctorConsultation[]> => {
    return http.get<DoctorConsultation[]>("/analytics/dashboard/consultations-per-doctor")
  },

  getOccupancyRate: async (): Promise<OccupancyRateResponse> => {
    return http.get<OccupancyRateResponse>("/analytics/dashboard/occupancy-rate")
  },

  getAvgWaitTime: async (): Promise<AvgWaitTimeResponse> => {
    return http.get<AvgWaitTimeResponse>("/analytics/dashboard/avg-wait-time")
  },

  getTopDiagnoses: async (): Promise<DiagnosisCount[]> => {
    return http.get<DiagnosisCount[]>("/analytics/dashboard/top-diagnoses")
  },
}
