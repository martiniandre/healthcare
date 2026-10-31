import { describe, it, expect } from "vitest"
import { calculateAge } from "./patientAge"

const referenceDate = new Date(2026, 5, 15)

describe("calculateAge", () => {
  it("should count full years when the birthday already happened this year", () => {
    expect(calculateAge("1990-03-10", referenceDate)).toBe(36)
  })

  it("should subtract one year when the birthday has not happened yet", () => {
    expect(calculateAge("1990-12-01", referenceDate)).toBe(35)
  })

  it("should count the birthday itself as completed", () => {
    expect(calculateAge("2000-06-15", referenceDate)).toBe(26)
  })

  it("should return null for invalid or future dates", () => {
    expect(calculateAge("not-a-date", referenceDate)).toBeNull()
    expect(calculateAge("2030-01-01", referenceDate)).toBeNull()
  })
})
