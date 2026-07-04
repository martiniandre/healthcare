import { useQuery } from "@tanstack/react-query"
import { dashboardApi } from "./dashboard_api"

export const dashboardKeys = {
  all: ["dashboard"] as const,
  consultationsPerDoctor: ["dashboard", "consultations-per-doctor"] as const,
  occupancyRate: ["dashboard", "occupancy-rate"] as const,
  avgWaitTime: ["dashboard", "avg-wait-time"] as const,
  topDiagnoses: ["dashboard", "top-diagnoses"] as const,
}

export const useDashboardQuery = () => {
  return useQuery({
    queryKey: dashboardKeys.all,
    queryFn: dashboardApi.getDashboard,
  })
}

export const useConsultationsPerDoctorQuery = () => {
  return useQuery({
    queryKey: dashboardKeys.consultationsPerDoctor,
    queryFn: dashboardApi.getConsultationsPerDoctor,
  })
}

export const useOccupancyRateQuery = () => {
  return useQuery({
    queryKey: dashboardKeys.occupancyRate,
    queryFn: dashboardApi.getOccupancyRate,
  })
}

export const useAvgWaitTimeQuery = () => {
  return useQuery({
    queryKey: dashboardKeys.avgWaitTime,
    queryFn: dashboardApi.getAvgWaitTime,
  })
}

export const useTopDiagnosesQuery = () => {
  return useQuery({
    queryKey: dashboardKeys.topDiagnoses,
    queryFn: dashboardApi.getTopDiagnoses,
  })
}
