import { describe, it, expect } from "vitest"
import { calculateTotalPages } from "./pagination"

describe("calculateTotalPages", () => {
  it("should return at least one page for an empty result set", () => {
    expect(calculateTotalPages(0, 5)).toBe(1)
  })

  it("should return one page when all patients fit in a single page", () => {
    expect(calculateTotalPages(5, 5)).toBe(1)
    expect(calculateTotalPages(3, 5)).toBe(1)
  })

  it("should round up partial pages", () => {
    expect(calculateTotalPages(6, 5)).toBe(2)
    expect(calculateTotalPages(11, 5)).toBe(3)
  })

  it("should handle exact multiples without an extra page", () => {
    expect(calculateTotalPages(10, 5)).toBe(2)
    expect(calculateTotalPages(15, 5)).toBe(3)
  })
})