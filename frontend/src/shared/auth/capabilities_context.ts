import { createContext, useContext } from "react"
import type { Capability } from "./capabilities"

export interface GrantedCapabilities {
  capabilities: ReadonlySet<string>
  hasCapability: (requiredCapability: Capability) => boolean
}

export function buildGrantedCapabilities(capabilities: readonly string[]): GrantedCapabilities {
  const capabilitySet = new Set(capabilities)
  return {
    capabilities: capabilitySet,
    hasCapability: (requiredCapability) => capabilitySet.has(requiredCapability),
  }
}

export const GrantedCapabilitiesContext = createContext<GrantedCapabilities>(buildGrantedCapabilities([]))

export function useCapabilities(): GrantedCapabilities {
  return useContext(GrantedCapabilitiesContext)
}
