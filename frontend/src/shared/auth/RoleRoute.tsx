import type { ReactNode } from "react"
import { useTranslation } from "react-i18next"
import { Navigate, useLocation } from "react-router-dom"
import { Alert, AlertDescription } from "../components/ui/alert"
import type { Capability } from "./capabilities"
import { useCapabilities } from "./capabilities_context"
import { getHomePathForRole, isCapabilityGranted } from "./route_access"

interface RoleRouteProps {
  requiredCapability: Capability
  children: ReactNode
}

export function RoleRoute({ requiredCapability, children }: RoleRouteProps) {
  const { t } = useTranslation("common")
  const location = useLocation()
  const { capabilities } = useCapabilities()

  if (isCapabilityGranted(capabilities, requiredCapability)) {
    return <>{children}</>
  }

  const homePath = getHomePathForRole(capabilities)
  if (location.pathname === homePath) {
    return (
      <div className="p-6">
        <Alert variant="destructive">
          <AlertDescription>{t("permissions.forbidden")}</AlertDescription>
        </Alert>
      </div>
    )
  }

  return <Navigate to={homePath} replace />
}
