import type { ReactNode } from "react"
import { Navigate } from "react-router-dom"
import { useAuthStore } from "../store/auth_store"
import { getHomePathForRole, isRoleAllowed } from "./route_access"

interface RoleRouteProps {
  allowedRoles: readonly string[]
  children: ReactNode
}

export function RoleRoute({ allowedRoles, children }: RoleRouteProps) {
  const role = useAuthStore((state) => state.role)

  if (!isRoleAllowed(role, allowedRoles)) {
    return <Navigate to={getHomePathForRole(role)} replace />
  }

  return <>{children}</>
}
