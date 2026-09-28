import { describe, it, expect } from "vitest"
import { resolveOccupancySeverity, resolveWaitTimeSeverity } from "./dashboard_thresholds"

describe("resolveOccupancySeverity", () => {
  it("treats a rate at or above the critical threshold as critical", () => {
    expect(resolveOccupancySeverity(85)).toBe("critical")
    expect(resolveOccupancySeverity(99)).toBe("critical")
  })

  it("treats a rate inside the elevated band as elevated", () => {
    expect(resolveOccupancySeverity(65)).toBe("elevated")
    expect(resolveOccupancySeverity(84.9)).toBe("elevated")
  })

  it("treats a rate below the elevated threshold as stable", () => {
    expect(resolveOccupancySeverity(64.9)).toBe("stable")
    expect(resolveOccupancySeverity(0)).toBe("stable")
  })
})

describe("resolveWaitTimeSeverity", () => {
  it("treats minutes at or above the critical threshold as critical", () => {
    expect(resolveWaitTimeSeverity(45)).toBe("critical")
    expect(resolveWaitTimeSeverity(120)).toBe("critical")
  })

  it("treats minutes inside the elevated band as elevated", () => {
    expect(resolveWaitTimeSeverity(25)).toBe("elevated")
    expect(resolveWaitTimeSeverity(44.9)).toBe("elevated")
  })

  it("treats minutes below the elevated threshold as stable", () => {
    expect(resolveWaitTimeSeverity(24.9)).toBe("stable")
    expect(resolveWaitTimeSeverity(0)).toBe("stable")
  })
})
