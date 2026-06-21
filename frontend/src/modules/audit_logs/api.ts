import { http } from "../../shared/utils/http"
import type { AuditLog, AuditLogsFilter } from "./types"

export const auditLogsApi = {
  listAuditLogs: async (filters: AuditLogsFilter): Promise<AuditLog[]> => {
    const queryParameters = new URLSearchParams()
    if (filters.action && filters.action !== "All") {
      queryParameters.append("action", filters.action)
    }
    if (filters.email) {
      queryParameters.append("email", filters.email)
    }
    if (filters.status && filters.status !== "All") {
      queryParameters.append("status", filters.status)
    }
    if (filters.startDate) {
      queryParameters.append("startDate", filters.startDate)
    }
    if (filters.endDate) {
      queryParameters.append("endDate", filters.endDate)
    }
    const queryString = queryParameters.toString()
    return http.get<AuditLog[]>(`/audit-logs${queryString ? `?${queryString}` : ""}`)
  },

  createAuditLog: async (payload: { action: string; details: string; status: string }): Promise<void> => {
    return http.post<void>("/audit-logs", payload)
  },
}
