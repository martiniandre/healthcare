import { describe, it, expect } from "vitest"
import { formatDicomFileSize, validateDicomFile, DICOM_MAX_FILE_BYTES } from "./dicom_file_validation"

describe("validateDicomFile", () => {
  it("should accept dcm files regardless of case", () => {
    expect(validateDicomFile(new File(["data"], "study.DCM"))).toBeNull()
  })

  it("should reject empty files", () => {
    expect(validateDicomFile(new File([""], "study.dcm"))).toBe("empty")
  })

  it("should reject files that are not dcm", () => {
    expect(validateDicomFile(new File(["data"], "photo.png"))).toBe("type")
  })

  it("should reject files above the limit", () => {
    const oversizedFile = new File(["data"], "study.dcm")
    Object.defineProperty(oversizedFile, "size", { value: DICOM_MAX_FILE_BYTES + 1 })
    expect(validateDicomFile(oversizedFile)).toBe("tooLarge")
  })
})

describe("formatDicomFileSize", () => {
  it("should format sizes up to gigabytes", () => {
    expect(formatDicomFileSize(100)).toBe("100 B")
    expect(formatDicomFileSize(5 * 1024 * 1024)).toBe("5.0 MB")
    expect(formatDicomFileSize(3 * 1024 * 1024 * 1024)).toBe("3.00 GB")
  })
})
