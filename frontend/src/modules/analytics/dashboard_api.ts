import { http } from "../../shared/utils/http"
import type { DashboardData } from "./dashboard_types"

export const dashboardApi = {
  getDashboard: async (): Promise<DashboardData> => {
    return http.get<DashboardData>("/analytics/dashboard")
  },
}
