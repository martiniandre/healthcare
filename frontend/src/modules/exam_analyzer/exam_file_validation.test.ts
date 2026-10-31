import { describe, it, expect } from "vitest"
import { formatFileSize, validateExamFile, EXAM_MAX_FILE_BYTES } from "./exam_file_validation"

const buildFile = (name: string, content: string) => new File([content], name)

describe("validateExamFile", () => {
  it("should accept png, jpg, jpeg and pdf files regardless of case", () => {
    expect(validateExamFile(buildFile("scan.PNG", "data"))).toBeNull()
    expect(validateExamFile(buildFile("scan.jpg", "data"))).toBeNull()
    expect(validateExamFile(buildFile("scan.jpeg", "data"))).toBeNull()
    expect(validateExamFile(buildFile("report.pdf", "data"))).toBeNull()
  })

  it("should reject empty files first", () => {
    expect(validateExamFile(buildFile("scan.txt", ""))).toBe("empty")
  })

  it("should reject unsupported extensions", () => {
    expect(validateExamFile(buildFile("study.dcm", "data"))).toBe("type")
    expect(validateExamFile(buildFile("anim.gif", "data"))).toBe("type")
  })

  it("should reject files above the limit", () => {
    const oversizedFile = buildFile("big.pdf", "data")
    Object.defineProperty(oversizedFile, "size", { value: EXAM_MAX_FILE_BYTES + 1 })
    expect(validateExamFile(oversizedFile)).toBe("tooLarge")
  })
})

describe("formatFileSize", () => {
  it("should format bytes, kilobytes and megabytes", () => {
    expect(formatFileSize(512)).toBe("512 B")
    expect(formatFileSize(2048)).toBe("2.0 KB")
    expect(formatFileSize(3 * 1024 * 1024)).toBe("3.00 MB")
  })
})
