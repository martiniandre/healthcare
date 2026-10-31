import { describe, it, expect } from "vitest"
import { resolveClinicalStatus } from "./clinicalStatus"

describe("resolveClinicalStatus", () => {
  it("should map known clinical statuses to presentation statuses", () => {
    expect(resolveClinicalStatus("allergy", "active")).toBe("attention")
    expect(resolveClinicalStatus("condition", "resolved")).toBe("normal")
    expect(resolveClinicalStatus("medication", "entered-in-error")).toBe("critical")
    expect(resolveClinicalStatus("report", "final")).toBe("normal")
    expect(resolveClinicalStatus("encounter", "in-progress")).toBe("info")
  })

  it("should fall back to neutral for unknown values", () => {
    expect(resolveClinicalStatus("report", "something-new")).toBe("neutral")
  })
})
