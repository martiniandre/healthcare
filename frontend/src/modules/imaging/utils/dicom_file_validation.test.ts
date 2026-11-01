import { describe, it, expect } from "vitest"
import {
  formatDicomFileSize,
  validateDicomFile,
  DICOM_MAX_ARCHIVE_BYTES,
  DICOM_MAX_FILE_BYTES,
} from "./dicom_file_validation"

const ZIP_LOCAL_HEADER_BYTES = [0x50, 0x4b, 0x03, 0x04, 0x14, 0x00]

const buildFileWithSize = (fileName: string, sizeInBytes: number, contentBytes: number[] = []) => {
  const file = new File([new Uint8Array(contentBytes)], fileName)
  Object.defineProperty(file, "size", { value: sizeInBytes })
  return file
}

describe("validateDicomFile", () => {
  it("should accept dcm files regardless of case", async () => {
    expect(await validateDicomFile(new File(["data"], "study.DCM"))).toBeNull()
  })

  it("should accept a zip archive whose content starts with a ZIP local header", async () => {
    expect(await validateDicomFile(new File([new Uint8Array(ZIP_LOCAL_HEADER_BYTES)], "exam.ZIP"))).toBeNull()
  })

  it("should reject a zip name whose content is not a ZIP archive", async () => {
    expect(await validateDicomFile(new File(["not a zip"], "exam.zip"))).toBe("invalidZip")
  })

  it("should reject empty files", async () => {
    expect(await validateDicomFile(new File([""], "study.dcm"))).toBe("empty")
    expect(await validateDicomFile(new File([""], "exam.zip"))).toBe("empty")
  })

  it("should reject files that are neither dcm nor zip", async () => {
    expect(await validateDicomFile(new File(["data"], "photo.png"))).toBe("type")
    expect(await validateDicomFile(new File(["Rar!"], "exam.rar"))).toBe("type")
  })

  it("should reject dcm files above the single file limit", async () => {
    expect(await validateDicomFile(buildFileWithSize("study.dcm", DICOM_MAX_FILE_BYTES + 1))).toBe("tooLarge")
  })

  it("should reject zip archives above the archive limit", async () => {
    expect(
      await validateDicomFile(buildFileWithSize("exam.zip", DICOM_MAX_ARCHIVE_BYTES + 1, ZIP_LOCAL_HEADER_BYTES))
    ).toBe("zipTooLarge")
  })
})

describe("formatDicomFileSize", () => {
  it("should format sizes up to gigabytes", () => {
    expect(formatDicomFileSize(100)).toBe("100 B")
    expect(formatDicomFileSize(5 * 1024 * 1024)).toBe("5.0 MB")
    expect(formatDicomFileSize(3 * 1024 * 1024 * 1024)).toBe("3.00 GB")
  })
})
