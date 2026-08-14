import { useQuery } from "@tanstack/react-query"
import { auditLogsApi } from "./api"
import type { AuditLogsFilter } from "./types"

export const auditLogsQueryKeys = {
  all: ["auditLogs"] as const,
  lists: () => [...auditLogsQueryKeys.all, "list"] as const,
}

export const useAuditLogsQuery = (filters: AuditLogsFilter) => {
  return useQuery({
    queryKey: [...auditLogsQueryKeys.lists(), filters],
    queryFn: () => auditLogsApi.listAuditLogs(filters),
  })
}
