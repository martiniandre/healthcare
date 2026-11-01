import type { ReactNode } from "react"
import { AbilityProvider } from "../shared/auth/AbilityContext"
import { backendMatrixCapabilitiesByRoleFixture, type RoleFixtureName } from "./backend_matrix_role_fixtures"

interface TestAbilityProviderProps {
  role?: RoleFixtureName
  capabilities?: readonly string[]
  children: ReactNode
}

export function TestAbilityProvider({ role, capabilities, children }: TestAbilityProviderProps) {
  const resolvedCapabilities = capabilities ?? (role ? backendMatrixCapabilitiesByRoleFixture[role] : [])
  return <AbilityProvider capabilities={resolvedCapabilities}>{children}</AbilityProvider>
}
