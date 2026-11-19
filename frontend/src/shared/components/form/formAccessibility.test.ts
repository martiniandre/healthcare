import { describe, expect, it } from "vitest"
import { joinDescribedBy } from "./formAccessibility"

describe("joinDescribedBy", () => {
  it("should join the present identifiers with spaces", () => {
    expect(joinDescribedBy("hint", false, undefined, "error")).toBe("hint error")
  })

  it("should return undefined when no identifier is present", () => {
    expect(joinDescribedBy(false, undefined)).toBeUndefined()
  })
})
