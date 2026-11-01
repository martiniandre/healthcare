import { renderHook } from "@testing-library/react"
import { AbilityProvider } from "./AbilityContext"
import { Capability } from "./capabilities"
import { Action, Feature } from "./types"
import { useAbility } from "./useAbility"
import { TestAbilityProvider } from "../../test/TestAbilityProvider"
import type { RoleFixtureName } from "../../test/backend_matrix_role_fixtures"

describe("useAbility", () => {
  function renderRoleAbility(role: RoleFixtureName) {
    return renderHook(() => useAbility(), {
      wrapper: ({ children }) => <TestAbilityProvider role={role}>{children}</TestAbilityProvider>,
    })
  }

  function renderCapabilityAbility(capabilities: readonly string[]) {
    return renderHook(() => useAbility(), {
      wrapper: ({ children }) => <AbilityProvider capabilities={capabilities}>{children}</AbilityProvider>,
    })
  }

  it("should grant exactly the rules mapped from the loaded capabilities", () => {
    const { result } = renderCapabilityAbility([Capability.PatientRead, Capability.AllergyCreate])

    expect(result.current.can(Action.Read, Feature.Patient)).toBe(true)
    expect(result.current.can(Action.Create, Feature.Allergy)).toBe(true)
    expect(result.current.can(Action.Create, Feature.Patient)).toBe(false)
    expect(result.current.can(Action.Read, Feature.Allergy)).toBe(false)
  })

  it("should ignore capabilities the frontend does not know", () => {
    const { result } = renderCapabilityAbility(["unknown:capability", Capability.PatientRead])

    expect(result.current.can(Action.Read, Feature.Patient)).toBe(true)
    expect(result.current.rules).toHaveLength(1)
  })

  it("should map appointment cancellation to the delete action", () => {
    const { result } = renderCapabilityAbility([Capability.AppointmentCancel])

    expect(result.current.can(Action.Delete, Feature.Appointment)).toBe(true)
    expect(result.current.can(Action.Read, Feature.Appointment)).toBe(false)
  })

  it("should let admins create medication requests", () => {
    const { result } = renderRoleAbility("ADMIN")

    expect(result.current.can(Action.Create, Feature.MedicationRequest)).toBe(true)
    expect(result.current.can(Action.Read, Feature.Patient)).toBe(true)
    expect(result.current.can(Action.Read, Feature.AuditLog)).toBe(true)
  })

  it("should let doctors read patients but not create them", () => {
    const { result } = renderRoleAbility("DOCTOR")

    expect(result.current.can(Action.Read, Feature.Patient)).toBe(true)
    expect(result.current.can(Action.Create, Feature.Patient)).toBe(false)
    expect(result.current.can(Action.Read, Feature.Staff)).toBe(true)
    expect(result.current.can(Action.Create, Feature.MedicationRequest)).toBe(true)
  })

  it("should keep nurses from prescribing or deleting exam analyses", () => {
    const { result } = renderRoleAbility("NURSE")

    expect(result.current.can(Action.Read, Feature.MedicationRequest)).toBe(true)
    expect(result.current.can(Action.Create, Feature.MedicationRequest)).toBe(false)
    expect(result.current.can(Action.Delete, Feature.ExamAnalysis)).toBe(false)
  })

  it("should keep reception out of the clinical timeline", () => {
    const { result } = renderRoleAbility("RECEPTION")

    expect(result.current.can(Action.Read, Feature.Timeline)).toBe(false)
    expect(result.current.can(Action.Create, Feature.Patient)).toBe(true)
  })

  it("should restrict patients to their portal", () => {
    const { result } = renderRoleAbility("PATIENT")

    expect(result.current.can(Action.Read, Feature.Portal)).toBe(true)
    expect(result.current.can(Action.Read, Feature.Patient)).toBe(false)
  })

  it("should be empty when no capabilities are loaded", () => {
    const { result } = renderCapabilityAbility([])

    expect(result.current.can(Action.Read, Feature.Portal)).toBe(false)
    expect(result.current.rules).toHaveLength(0)
  })
})
