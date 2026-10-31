import { describe, it, expect } from "vitest"
import { contrastRatio, evaluateContrast, parseColor, relativeLuminance } from "./contrast"

describe("contrast helpers", () => {
  it("should parse hex, rgb and srgb color strings", () => {
    expect(parseColor("#fff")).toEqual({ red: 255, green: 255, blue: 255 })
    expect(parseColor("#2563eb")).toEqual({ red: 37, green: 99, blue: 235 })
    expect(parseColor("rgb(31, 41, 55)")).toEqual({ red: 31, green: 41, blue: 55 })
    expect(parseColor("rgba(0 0 0 / 0.5)")).toEqual({ red: 0, green: 0, blue: 0 })
    expect(parseColor("color(srgb 1 0 0)")).toEqual({ red: 255, green: 0, blue: 0 })
  })

  it("should return null for unsupported values", () => {
    expect(parseColor("")).toBeNull()
    expect(parseColor("oklch(0.5 0.1 200)")).toBeNull()
  })

  it("should compute relative luminance bounds", () => {
    expect(relativeLuminance({ red: 0, green: 0, blue: 0 })).toBe(0)
    expect(relativeLuminance({ red: 255, green: 255, blue: 255 })).toBeCloseTo(1, 5)
  })

  it("should return 21 for black on white regardless of order", () => {
    const black = { red: 0, green: 0, blue: 0 }
    const white = { red: 255, green: 255, blue: 255 }
    expect(contrastRatio(black, white)).toBeCloseTo(21, 5)
    expect(contrastRatio(white, black)).toBeCloseTo(21, 5)
  })

  it("should return 1 for identical colors", () => {
    const color = { red: 120, green: 40, blue: 200 }
    expect(contrastRatio(color, color)).toBeCloseTo(1, 5)
  })

  it("should apply AA and AAA thresholds for text", () => {
    const strong = evaluateContrast("#000000", "#ffffff")
    expect(strong).toMatchObject({ passesAa: true, passesAaa: true })
    const medium = evaluateContrast("#767676", "#ffffff")
    expect(medium?.ratio).toBeCloseTo(4.54, 1)
    expect(medium).toMatchObject({ passesAa: true, passesAaa: false })
    const weak = evaluateContrast("#aaaaaa", "#ffffff")
    expect(weak).toMatchObject({ passesAa: false, passesAaa: false })
  })

  it("should use relaxed thresholds for ui components", () => {
    const evaluation = evaluateContrast("#949494", "#ffffff", "ui")
    expect(evaluation).toMatchObject({ passesAa: true, passesAaa: false })
  })

  it("should return null when a color cannot be parsed", () => {
    expect(evaluateContrast("", "#ffffff")).toBeNull()
  })
})
