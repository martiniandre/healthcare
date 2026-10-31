import { describe, it, expect } from "vitest"
import { formatDecimal, formatPercentage } from "./format_numbers"

describe("format_numbers", () => {
  it("formats decimals with the locale separator", () => {
    expect(formatDecimal(72.4, "en-US", 1)).toBe("72.4")
    expect(formatDecimal(72.4, "pt-BR", 1)).toBe("72,4")
  })

  it("formats grouped integers", () => {
    expect(formatDecimal(1234, "en-US")).toBe("1,234")
  })

  it("formats percentages from a 0-100 value", () => {
    expect(formatPercentage(71, "en-US")).toBe("71%")
    expect(formatPercentage(98.5, "en-US", 1)).toBe("98.5%")
  })

  it("falls back when the locale is invalid", () => {
    expect(formatDecimal(3.14159, "invalid_locale!!", 2)).toBe("3.14")
  })
})
