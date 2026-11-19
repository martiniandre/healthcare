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

