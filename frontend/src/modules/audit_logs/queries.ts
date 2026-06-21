import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query"
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

export const useCreateAuditLogMutation = () => {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (payload: { method: string; correlation_id: string; access_granted: boolean }) =>
      auditLogsApi.createAuditLog(payload),
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: auditLogsQueryKeys.lists(),
      })
    },
  })
}
