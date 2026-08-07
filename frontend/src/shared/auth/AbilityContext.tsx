import { Can, AbilityProvider as CaslAbilityProvider } from "@casl/react"
import { defineAppAbility } from "./ability"
import { Action, Feature } from "./types"
import type { ReactNode } from "react"

export { Can, Action, Feature }

interface AbilityProviderProps {
  role: string | null
  children: ReactNode
}

export function AbilityProvider({ role, children }: AbilityProviderProps) {
  const ability = defineAppAbility(role)
  return <CaslAbilityProvider value={ability}>{children}</CaslAbilityProvider>
}
