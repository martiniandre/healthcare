import { ImagingStudyStatus } from "../../../shared/types"
import {
  createValidDicomBlob,
  normalizeDicomStudyDate,
  resolveStudyClinicalStatus,
  shortenDicomUid,
  waitForUploadFrame,
} from "./pacs_helpers"

describe("pacs_helpers", () => {
  it("should create a DICOM blob with a valid preamble signature", async () => {
    const dicomBlob = createValidDicomBlob()
    const preambleBytes = new Uint8Array(await dicomBlob.slice(128, 132).arrayBuffer())

    expect(dicomBlob.type).toBe("application/dicom")
    expect(Array.from(preambleBytes)).toEqual([68, 73, 67, 77])
  })

  it("should wait for the requested upload frame duration", async () => {
    vi.useFakeTimers()
    const completionMarker = vi.fn()
    const waitingPromise = waitForUploadFrame(300).then(completionMarker)

    await vi.advanceTimersByTimeAsync(300)
    await waitingPromise

    expect(completionMarker).toHaveBeenCalledTimes(1)
    vi.useRealTimers()
  })
})

describe("resolveStudyClinicalStatus", () => {
  it("should map each study status to a clinical status", () => {
    expect(resolveStudyClinicalStatus(ImagingStudyStatus.Available)).toBe("normal")
    expect(resolveStudyClinicalStatus(ImagingStudyStatus.Registered)).toBe("info")
    expect(resolveStudyClinicalStatus(ImagingStudyStatus.Cancelled)).toBe("critical")
    expect(resolveStudyClinicalStatus(ImagingStudyStatus.EnteredInError)).toBe("critical")
    expect(resolveStudyClinicalStatus(ImagingStudyStatus.Unknown)).toBe("neutral")
  })

  it("should normalize DICOM dates and keep other formats untouched", () => {
    expect(normalizeDicomStudyDate("20260315")).toBe("2026-03-15")
    expect(normalizeDicomStudyDate("2026-03-15")).toBe("2026-03-15")
  })

  it("should shorten long UIDs and keep short ones", () => {
    expect(shortenDicomUid("1.2.3")).toBe("1.2.3")
    expect(shortenDicomUid("1.2.840.113619.2.55.3.604688119.969.1234567890.123")).toBe("1.2.840.11…567890.123")
  })
})
