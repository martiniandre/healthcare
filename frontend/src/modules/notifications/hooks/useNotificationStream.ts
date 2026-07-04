import { useEffect, useRef } from "react"
import { useQueryClient } from "@tanstack/react-query"
import { notificationKeys } from "../queries"
import { useAuthStore } from "../../../shared/store/auth_store"

export function useNotificationStream() {
  const queryClient = useQueryClient()
  const isAuthenticated = useAuthStore((state) => state.isAuthenticated)
  const eventSourceRef = useRef<EventSource | null>(null)

  useEffect(() => {
    if (!isAuthenticated) {
      return
    }

    const eventSource = new EventSource("/api/v1/notifications/stream")
    eventSourceRef.current = eventSource

    eventSource.addEventListener("notification", () => {
      queryClient.invalidateQueries({ queryKey: notificationKeys.all })
    })

    eventSource.onerror = () => {
      eventSource.close()
    }

    return () => {
      eventSource.close()
      eventSourceRef.current = null
    }
  }, [isAuthenticated, queryClient])
}
