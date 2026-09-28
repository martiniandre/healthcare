import { describe, it, expect } from "vitest"
import {
  resolveFhirComplianceSeverity,
  resolveTrendDirection,
  resolveTrendReading,
  resolveDayTranslationKey,
  summarizeWeeklyConsultations,
} from "./stats_domain"
import type { ConsultationsDayData } from "./types"

const buildWeek = (counts: number[]): ConsultationsDayData[] => {
  const dayKeys = ["mon", "tue", "wed", "thu", "fri", "sat", "sun"]
  return counts.map((count, indexValue) => ({
    dayName: `analytics.days.${dayKeys[indexValue]}`,
    count,
  }))
}

describe("resolveFhirComplianceSeverity", () => {
  it("treats a rate at or above the elevated ceiling as compliant", () => {
    expect(resolveFhirComplianceSeverity(100)).toBe("stable")
    expect(resolveFhirComplianceSeverity(98)).toBe("stable")
    expect(resolveFhirComplianceSeverity(99.4)).toBe("stable")
  })

  it("treats a rate below the compliant ceiling but at or above the floor as elevated", () => {
    expect(resolveFhirComplianceSeverity(97.9)).toBe("elevated")
    expect(resolveFhirComplianceSeverity(90)).toBe("elevated")
  })

  it("treats a rate below the floor as critical", () => {
    expect(resolveFhirComplianceSeverity(89.9)).toBe("critical")
    expect(resolveFhirComplianceSeverity(0)).toBe("critical")
  })
})

describe("resolveTrendDirection", () => {
  it("reads a leading plus as a rising trend", () => {
    expect(resolveTrendDirection("+5%")).toBe("rising")
    expect(resolveTrendDirection("+12%")).toBe("rising")
    expect(resolveTrendDirection("  +1%")).toBe("rising")
  })

  it("reads a leading hyphen or unicode minus as a falling trend", () => {
    expect(resolveTrendDirection("-3%")).toBe("falling")
    expect(resolveTrendDirection("\u22122%")).toBe("falling")
  })

  it("treats an unprefixed value as flat", () => {
    expect(resolveTrendDirection("stable")).toBe("flat")
    expect(resolveTrendDirection("Estável")).toBe("flat")
  })

  it("treats a blank trend as unknown", () => {
    expect(resolveTrendDirection("")).toBe("unknown")
    expect(resolveTrendDirection("   ")).toBe("unknown")
  })
})

describe("resolveTrendReading", () => {
  it("escalates a large rise to critical", () => {
    expect(resolveTrendReading("+12%")).toEqual({
      direction: "rising",
      magnitude: 12,
      severity: "critical",
    })
    expect(resolveTrendReading("+10%").severity).toBe("critical")
  })

  it("treats a moderate rise as elevated", () => {
    expect(resolveTrendReading("+5%").severity).toBe("elevated")
    expect(resolveTrendReading("+7.5%").severity).toBe("elevated")
  })

  it("keeps a small rise neutral instead of alarming", () => {
    expect(resolveTrendReading("+0.1%").severity).toBe("neutral")
    expect(resolveTrendReading("+4.9%").severity).toBe("neutral")
  })

  it("treats any decline as stable regardless of magnitude", () => {
    expect(resolveTrendReading("-1%").severity).toBe("stable")
    expect(resolveTrendReading("-40%").severity).toBe("stable")
  })

  it("treats a flat or unparsable trend as neutral", () => {
    expect(resolveTrendReading("stable").severity).toBe("neutral")
    expect(resolveTrendReading("").severity).toBe("neutral")
  })

  it("parses a decimal magnitude", () => {
    expect(resolveTrendReading("+12.5%").magnitude).toBe(12.5)
  })
})

describe("resolveDayTranslationKey", () => {  it("normalizes a fully qualified key to the module namespace", () => {
    expect(resolveDayTranslationKey("analytics.days.mon")).toBe("days.mon")
  })

  it("normalizes a bare day name to the module namespace", () => {
    expect(resolveDayTranslationKey("mon")).toBe("days.mon")
  })

  it("normalizes mixed casing and padding", () => {
    expect(resolveDayTranslationKey("  analytics.days.FRI ")).toBe("days.fri")
  })
})

describe("summarizeWeeklyConsultations", () => {
  it("returns a zeroed summary for an empty week", () => {
    expect(summarizeWeeklyConsultations([])).toEqual({
      total: 0,
      min: 0,
      average: 0,
      peak: 0,
      peakDayKey: null,
    })
  })

  it("computes total, min, rounded average and peak", () => {
    const summary = summarizeWeeklyConsultations(buildWeek([12, 14, 13, 11, 15, 8, 6]))

    expect(summary.total).toBe(79)
    expect(summary.min).toBe(6)
    expect(summary.average).toBe(11)
    expect(summary.peak).toBe(15)
  })

  it("resolves the peak day to a translation key", () => {
    const summary = summarizeWeeklyConsultations(buildWeek([12, 14, 13, 11, 15, 8, 6]))
    expect(summary.peakDayKey).toBe("days.fri")
  })

  it("reports no peak day when every day is zero", () => {
    const summary = summarizeWeeklyConsultations(buildWeek([0, 0, 0, 0, 0, 0, 0]))

    expect(summary.peak).toBe(0)
    expect(summary.peakDayKey).toBeNull()
  })
})
