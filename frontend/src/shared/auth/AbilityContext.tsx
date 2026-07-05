import { Can, AbilityProvider as CaslAbilityProvider, useAbility as useCaslAbility } from "@casl/react"
import { defineAppAbility } from "./ability"
import { Action, Feature } from "./types"
import type { AppAbility } from "./types"
import type { ReactNode } from "react"

export { Can, Action, Feature }

export function useAbility(): AppAbility {
  return useCaslAbility<AppAbility>()
}

interface AbilityProviderProps {
  role: string | null
  children: ReactNode
}

export function AbilityProvider({ role, children }: AbilityProviderProps) {
  const ability = defineAppAbility(role)
  return <CaslAbilityProvider value={ability}>{children}</CaslAbilityProvider>
}
