import { describe, it, expect } from "vitest"
import { BedStatus } from "../../../shared/types"
import { appendBedSamples, resolveBedClinicalStatus, sortBedsBySeverity } from "./bed_helpers"

describe("bed_helpers", () => {
  it("should map bed status to clinical status", () => {
    expect(resolveBedClinicalStatus(BedStatus.Danger)).toBe("critical")
    expect(resolveBedClinicalStatus(BedStatus.Warning)).toBe("attention")
    expect(resolveBedClinicalStatus(BedStatus.Normal)).toBe("normal")
  })

  it("should sort critical first, then attention, then normal, stable by bed number", () => {
    const beds = [
      { bedNumber: "Leito 10", status: BedStatus.Normal },
      { bedNumber: "Leito 2", status: BedStatus.Normal },
      { bedNumber: "Leito 3", status: BedStatus.Warning },
      { bedNumber: "Leito 7", status: BedStatus.Danger },
      { bedNumber: "Leito 1", status: BedStatus.Danger },
    ]

    const sortedNumbers = sortBedsBySeverity(beds).map((bed) => bed.bedNumber)

    expect(sortedNumbers).toEqual(["Leito 1", "Leito 7", "Leito 3", "Leito 2", "Leito 10"])
  })

  it("should not mutate the input list", () => {
    const beds = [
      { bedNumber: "B", status: BedStatus.Normal },
      { bedNumber: "A", status: BedStatus.Danger },
    ]

    sortBedsBySeverity(beds)

    expect(beds[0].bedNumber).toBe("B")
  })

  it("should append samples per bed and bound the history", () => {
    let history = {}
    for (let reading = 0; reading < 5; reading++) {
      history = appendBedSamples(history, [{ id: "bed-1", bpm: 60 + reading, spo2: 98 }], 3)
    }

    expect(history).toEqual({
      "bed-1": [
        { bpm: 62, spo2: 98 },
        { bpm: 63, spo2: 98 },
        { bpm: 64, spo2: 98 },
      ],
    })
  })
})
