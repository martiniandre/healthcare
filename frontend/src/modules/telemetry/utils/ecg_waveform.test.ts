import { describe, expect, it } from "vitest"
import { BedStatus, CardiacCondition } from "@/shared/types"
import { computeEcgSample, cyclePeriodFor, simulatedVitalsFor } from "./ecg_waveform"

describe("cyclePeriodFor", () => {
  it("should slow the cycle for bradycardia and speed it up for tachycardia", () => {
    expect(cyclePeriodFor(CardiacCondition.Bradycardia)).toBeGreaterThan(cyclePeriodFor(CardiacCondition.Normal))
    expect(cyclePeriodFor(CardiacCondition.Tachycardia)).toBeLessThan(cyclePeriodFor(CardiacCondition.Normal))
  })
})

describe("computeEcgSample", () => {
  const baseline = 100

  it("should stay on the baseline at the start of a cycle", () => {
    expect(computeEcgSample(0, baseline, CardiacCondition.Normal)).toBe(baseline)
  })

  it("should draw the R wave peak above the baseline", () => {
    expect(computeEcgSample(20, baseline, CardiacCondition.Normal)).toBeLessThan(baseline - 30)
  })

  it("should stay almost flat during cardiac arrest", () => {
    for (let traceOffset = 0; traceOffset < 120; traceOffset += 7) {
      expect(Math.abs(computeEcgSample(traceOffset, baseline, CardiacCondition.CardiacArrest) - baseline)).toBeLessThanOrEqual(0.8)
    }
  })
})

describe("simulatedVitalsFor", () => {
  it("should flag cardiac arrest as danger with no pulse", () => {
    expect(simulatedVitalsFor(CardiacCondition.CardiacArrest)).toEqual({ bpm: 0, spo2: 0, status: BedStatus.Danger })
  })

  it("should return normal vitals for a normal rhythm", () => {
    expect(simulatedVitalsFor(CardiacCondition.Normal).status).toBe(BedStatus.Normal)
  })
})
