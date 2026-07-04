import { http } from "../../shared/utils/http"
import type { NotificationListResponse, UnreadCountResponse } from "./types"

export const notificationsApi = {
  list: async (limit = 50, offset = 0): Promise<NotificationListResponse> => {
    return http.get<NotificationListResponse>(`/notifications?limit=${limit}&offset=${offset}`)
  },

  markRead: async (notificationId: string): Promise<void> => {
    await http.post<void>(`/notifications/${notificationId}/read`)
  },

  getUnreadCount: async (): Promise<UnreadCountResponse> => {
    return http.get<UnreadCountResponse>("/notifications/unread-count")
  },
}
