import { ImagingStudyStatus } from "../../../shared/types"
import { createValidDicomBlob, resolveStudyClinicalStatus, waitForUploadFrame } from "./pacs_helpers"

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
})
