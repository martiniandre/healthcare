import { describe, it, expect } from "vitest"
import { formatNumber } from "./numbers"

describe("formatNumber", () => {
  it("should use the decimal separator of the explicit locale", () => {
    expect(formatNumber(1234.5, undefined, "pt-BR")).toBe("1.234,5")
    expect(formatNumber(1234.5, undefined, "en-US")).toBe("1,234.5")
  })

  it("should round to two fraction digits by default", () => {
    expect(formatNumber(1.23456, undefined, "en-US")).toBe("1.23")
  })

  it("should fall back to the plain string for an invalid locale", () => {
    expect(formatNumber(42, undefined, "not a locale")).toBe("42")
  })
})
