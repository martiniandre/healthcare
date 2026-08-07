import { useAbility as useCaslAbility } from "@casl/react"
import type { AppAbility } from "./types"

export function useAbility(): AppAbility {
  return useCaslAbility<AppAbility>()
}
