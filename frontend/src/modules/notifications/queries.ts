import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query"
import { notificationsApi } from "./api"

export const notificationKeys = {
  all: ["notifications"] as const,
  list: () => [...notificationKeys.all, "list"] as const,
  unreadCount: () => [...notificationKeys.all, "unreadCount"] as const,
}

export function useNotificationsQuery(limit?: number, offset?: number) {
  return useQuery({
    queryKey: [...notificationKeys.list(), limit, offset],
    queryFn: () => notificationsApi.list(limit, offset),
  })
}

export function useUnreadCountQuery() {
  return useQuery({
    queryKey: notificationKeys.unreadCount(),
    queryFn: () => notificationsApi.getUnreadCount(),
    refetchInterval: 30000,
  })
}

export function useMarkReadMutation() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (notificationId: string) => notificationsApi.markRead(notificationId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: notificationKeys.all })
    },
  })
}
