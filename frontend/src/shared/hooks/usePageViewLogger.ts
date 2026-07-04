import { useEffect } from "react"
import { useLocation } from "react-router-dom"
import { useAuthStore } from "../store/auth_store"
import { auditLogsApi } from "../../modules/audit_logs/api"

export const usePageViewLogger = () => {
  const location = useLocation()
  const isAuthenticated = useAuthStore((state) => state.isAuthenticated)

  useEffect(() => {
    if (isAuthenticated) {
      auditLogsApi.createAuditLog({
        method: "PAGE_VIEW",
        correlation_id: `Viewed page: ${location.pathname}${location.search}`,
        access_granted: true,
      }).catch((logError) => {
        console.error("Failed to log page view", logError)
      })
    }
  }, [location.pathname, location.search, isAuthenticated])
}
