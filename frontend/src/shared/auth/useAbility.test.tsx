import { renderHook } from "@testing-library/react"
import { AbilityProvider } from "./AbilityContext"
import { Action, Feature } from "./types"
import { useAbility } from "./useAbility"

describe("useAbility", () => {
  function renderAbility(role: string | null) {
    return renderHook(() => useAbility(), {
      wrapper: ({ children }) => <AbilityProvider role={role}>{children}</AbilityProvider>,
    })
  }

  it("should return an ability that grants the admin full management rights", () => {
    const { result } = renderAbility("ADMIN")

    expect(result.current.can(Action.Manage, Feature.All)).toBe(true)
    expect(result.current.can(Action.Read, Feature.Patient)).toBe(true)
  })

  it("should return an ability that forbids admin from creating medication requests", () => {
    const { result } = renderAbility("ADMIN")

    expect(result.current.can(Action.Create, Feature.MedicationRequest)).toBe(false)
  })

  it("should return an ability that allows doctors to read patients but not create them", () => {
    const { result } = renderAbility("DOCTOR")

    expect(result.current.can(Action.Read, Feature.Patient)).toBe(true)
    expect(result.current.can(Action.Create, Feature.Patient)).toBe(false)
    expect(result.current.can(Action.Read, Feature.Staff)).toBe(true)
  })

  it("should return an ability that restricts patients to their portal", () => {
    const { result } = renderAbility("PATIENT")

    expect(result.current.can(Action.Read, Feature.Portal)).toBe(true)
    expect(result.current.can(Action.Read, Feature.Patient)).toBe(false)
  })

  it("should return an empty ability for unknown roles", () => {
    const { result } = renderAbility(null)

    expect(result.current.can(Action.Read, Feature.Portal)).toBe(false)
    expect(result.current.can(Action.Manage, Feature.All)).toBe(false)
  })
})
