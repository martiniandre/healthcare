import { allCapabilities, Capability, capabilityPattern } from "./capabilities"
import { capabilityAbilityRules } from "./ability"
import { backendMatrixCapabilitiesByRoleFixture } from "../../test/backend_matrix_role_fixtures"

describe("capability contract", () => {
  it("should use the resource:action syntax for every capability", () => {
    for (const capability of allCapabilities) {
      expect(capability).toMatch(capabilityPattern)
    }
  })

  it("should not declare the same capability string twice", () => {
    expect(new Set(allCapabilities).size).toBe(allCapabilities.length)
    expect(Object.keys(Capability)).toHaveLength(allCapabilities.length)
  })

  it("should map every capability to exactly one ability rule", () => {
    expect(Object.keys(capabilityAbilityRules).sort()).toEqual([...allCapabilities].sort())
  })

  it("should only use declared capabilities in the role fixtures", () => {
    for (const fixtureCapabilities of Object.values(backendMatrixCapabilitiesByRoleFixture)) {
      for (const capability of fixtureCapabilities) {
        expect(allCapabilities).toContain(capability)
      }
    }
  })
})
