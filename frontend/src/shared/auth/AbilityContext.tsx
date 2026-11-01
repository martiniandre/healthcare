import { useMemo } from "react"
import { Can, AbilityProvider as CaslAbilityProvider } from "@casl/react"
import { defineAppAbility } from "./ability"
import { buildGrantedCapabilities, GrantedCapabilitiesContext } from "./capabilities_context"
import { Action, Feature } from "./types"
import type { ReactNode } from "react"

export { Can, Action, Feature }

interface AbilityProviderProps {
  capabilities: readonly string[]
  children: ReactNode
}

export function AbilityProvider({ capabilities, children }: AbilityProviderProps) {
  const ability = useMemo(() => defineAppAbility(capabilities), [capabilities])
  const grantedCapabilities = useMemo(() => buildGrantedCapabilities(capabilities), [capabilities])

  return (
    <GrantedCapabilitiesContext.Provider value={grantedCapabilities}>
      <CaslAbilityProvider value={ability}>{children}</CaslAbilityProvider>
    </GrantedCapabilitiesContext.Provider>
  )
}
