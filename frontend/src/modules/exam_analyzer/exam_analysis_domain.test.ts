import { describe, it, expect } from "vitest"
import {
  QUALITY_VERDICT_THRESHOLDS,
  CONFIDENCE_BAND_THRESHOLDS,
  resolveQualityVerdict,
  resolveConfidenceBand,
  resolveFindingTriage,
  resolveUrgencyTone,
  resolveHistoryStatus,
  truncateIdentifier,
} from "./exam_analysis_domain"

describe("resolveQualityVerdict", () => {
  it("treats a high score as good", () => {
    expect(resolveQualityVerdict(0.9)).toEqual({
      tone: "good",
      translationKey: "card.qualityVerdict.good",
      percentage: 90,
    })
  })

  it("puts the good threshold in the good band", () => {
    expect(resolveQualityVerdict(QUALITY_VERDICT_THRESHOLDS.good).tone).toBe("good")
  })

  it("puts a score just below the good threshold in the fair band", () => {
    expect(
      resolveQualityVerdict(QUALITY_VERDICT_THRESHOLDS.good - 0.01).tone
    ).toBe("fair")
  })

  it("puts the fair threshold in the fair band", () => {
    expect(resolveQualityVerdict(QUALITY_VERDICT_THRESHOLDS.fair).tone).toBe("fair")
  })

  it("puts a score just below the fair threshold in the poor band", () => {
    expect(
      resolveQualityVerdict(QUALITY_VERDICT_THRESHOLDS.fair - 0.01).tone
    ).toBe("poor")
  })

  it("rounds the displayed percentage", () => {
    expect(resolveQualityVerdict(0.876).percentage).toBe(88)
  })

  it("clamps a score above one instead of overflowing the meter", () => {
    expect(resolveQualityVerdict(1.4).percentage).toBe(100)
    expect(resolveQualityVerdict(1.4).tone).toBe("good")
  })

  it("clamps a negative score", () => {
    expect(resolveQualityVerdict(-0.5).percentage).toBe(0)
    expect(resolveQualityVerdict(-0.5).tone).toBe("poor")
  })

  it("treats a non numeric score as poor rather than crashing", () => {
    expect(resolveQualityVerdict(Number.NaN).tone).toBe("poor")
  })
})

describe("resolveConfidenceBand", () => {
  it("puts the strong threshold in the strong band", () => {
    expect(resolveConfidenceBand(CONFIDENCE_BAND_THRESHOLDS.strong).band).toBe("strong")
  })

  it("puts a value just below the strong threshold in the moderate band", () => {
    expect(
      resolveConfidenceBand(CONFIDENCE_BAND_THRESHOLDS.strong - 0.01).band
    ).toBe("moderate")
  })

  it("puts the moderate threshold in the moderate band", () => {
    expect(resolveConfidenceBand(CONFIDENCE_BAND_THRESHOLDS.moderate).band).toBe("moderate")
  })

  it("puts a value just below the moderate threshold in the weak band", () => {
    expect(
      resolveConfidenceBand(CONFIDENCE_BAND_THRESHOLDS.moderate - 0.01).band
    ).toBe("weak")
  })

  it("exposes a translation key and rounded percentage", () => {
    expect(resolveConfidenceBand(0.92)).toEqual({
      band: "strong",
      translationKey: "card.confidenceBand.strong",
      percentage: 92,
    })
  })

  it("clamps an out of range confidence", () => {
    expect(resolveConfidenceBand(1.2).percentage).toBe(100)
  })
})

describe("resolveFindingTriage", () => {
  it("marks a high severity finding with strong confidence as a priority", () => {
    expect(resolveFindingTriage("high", 0.88)).toEqual({
      tier: "priority",
      translationKey: "card.triage.priority",
    })
  })

  it("asks for verification when severity is high but confidence is only moderate", () => {
    expect(resolveFindingTriage("high", 0.65).tier).toBe("verify")
  })

  it("asks for verification when severity is high but confidence is weak", () => {
    expect(resolveFindingTriage("high", 0.2).tier).toBe("verify")
  })

  it("asks for review when severity is medium and confidence is strong", () => {
    expect(resolveFindingTriage("medium", 0.9).tier).toBe("review")
  })

  it("asks for review when severity is medium and confidence is only moderate", () => {
    expect(resolveFindingTriage("medium", 0.5).tier).toBe("review")
  })

  it("keeps a medium severity finding with weak confidence as informative", () => {
    expect(resolveFindingTriage("medium", 0.3).tier).toBe("informative")
  })

  it("keeps any low severity finding as informative even with strong confidence", () => {
    expect(resolveFindingTriage("low", 0.99).tier).toBe("informative")
  })

  it("never escalates a low severity finding to priority", () => {
    expect(resolveFindingTriage("low", 1).tier).not.toBe("priority")
  })
})

describe("resolveUrgencyTone", () => {
  it("maps urgent to the urgent tone", () => {
    expect(resolveUrgencyTone("urgent")).toEqual({
      tone: "urgent",
      translationKey: "card.urgencyUrgent",
    })
  })

  it("maps medical followup to the followup tone", () => {
    expect(resolveUrgencyTone("medical_followup")).toEqual({
      tone: "followup",
      translationKey: "card.urgencyFollowup",
    })
  })

  it("maps normal to the routine tone", () => {
    expect(resolveUrgencyTone("normal")).toEqual({
      tone: "routine",
      translationKey: "card.urgencyNormal",
    })
  })
})

describe("resolveHistoryStatus", () => {
  it("labels a completed analysis", () => {
    expect(resolveHistoryStatus("completed")).toEqual({
      tone: "completed",
      translationKey: "history.statusCompleted",
    })
  })

  it("distinguishes pending from processing while keeping the same tone", () => {
    const pendingReading = resolveHistoryStatus("pending")
    const processingReading = resolveHistoryStatus("processing")

    expect(pendingReading.tone).toBe("inProgress")
    expect(processingReading.tone).toBe("inProgress")
    expect(pendingReading.translationKey).not.toBe(processingReading.translationKey)
  })

  it("marks insufficient data as needing attention", () => {
    expect(resolveHistoryStatus("insufficient_data")).toEqual({
      tone: "attention",
      translationKey: "history.statusQuality",
    })
  })

  it("marks a failed analysis", () => {
    expect(resolveHistoryStatus("failed")).toEqual({
      tone: "failed",
      translationKey: "history.statusFailed",
    })
  })
})

describe("truncateIdentifier", () => {
  it("shortens a long identifier to eight characters", () => {
    expect(truncateIdentifier("3f7c1a9e-2b44-4d18-9c0a-5e6f7a8b9c01")).toBe("3f7c1a9e…")
  })

  it("leaves a short identifier untouched", () => {
    expect(truncateIdentifier("ana-1")).toBe("ana-1")
  })

  it("leaves an identifier of exactly the visible length untouched", () => {
    expect(truncateIdentifier("12345678")).toBe("12345678")
  })

  it("returns an empty string unchanged", () => {
    expect(truncateIdentifier("")).toBe("")
  })
})
