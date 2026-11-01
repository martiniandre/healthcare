import { describe, it, expect } from "vitest"
import {
  resolveInstanceRejectionReasonKey,
  resolveUploadRejectionMessageKeys,
  resolveValidationMessageKey,
  resolveWarningMessageKey,
} from "./dicom_upload_messages"
import { DicomInstanceRejectionCode, DicomUploadErrorCode } from "../types"

describe("resolveUploadRejectionMessageKeys", () => {
  it.each([
    [DicomUploadErrorCode.PATIENT_MISMATCH, "upload.errors.patientMismatchTitle"],
    [DicomUploadErrorCode.INVALID_DICOM, "upload.errors.invalidDicomTitle"],
    [DicomUploadErrorCode.DUPLICATE_INSTANCE, "upload.errors.duplicateTitle"],
    [DicomUploadErrorCode.ZIP_INVALID, "upload.errors.zipInvalidTitle"],
    [DicomUploadErrorCode.ZIP_UNSAFE_ENTRY, "upload.errors.zipUnsafeTitle"],
    [DicomUploadErrorCode.ZIP_LIMITS_EXCEEDED, "upload.errors.zipLimitsTitle"],
    [DicomUploadErrorCode.NO_VALID_INSTANCES, "upload.errors.noValidInstancesTitle"],
  ])("should map %s to a title and a help message", (errorCode, expectedTitleKey) => {
    const keys = resolveUploadRejectionMessageKeys(errorCode)

    expect(keys?.title).toBe(expectedTitleKey)
    expect(keys?.help).toBeTruthy()
  })

  it("should return null for codes without a specific message", () => {
    expect(resolveUploadRejectionMessageKeys("SOMETHING_NEW")).toBeNull()
    expect(resolveUploadRejectionMessageKeys("constructor")).toBeNull()
    expect(resolveUploadRejectionMessageKeys(undefined)).toBeNull()
  })
})

describe("resolveInstanceRejectionReasonKey", () => {
  it.each([
    [DicomInstanceRejectionCode.DUPLICATE_INSTANCE, "upload.result.reasons.duplicate"],
    [DicomInstanceRejectionCode.INVALID_DICOM, "upload.result.reasons.invalid"],
    [DicomInstanceRejectionCode.UNSUPPORTED_TRANSFER_SYNTAX, "upload.result.reasons.unsupportedTransferSyntax"],
    [DicomInstanceRejectionCode.MISSING_IDENTIFIERS, "upload.result.reasons.missingIdentifiers"],
  ])("should map %s to a translated reason", (rejectionCode, expectedKey) => {
    expect(resolveInstanceRejectionReasonKey(rejectionCode)).toBe(expectedKey)
  })

  it("should fall back to the generic reason for unknown codes", () => {
    expect(resolveInstanceRejectionReasonKey("BRAND_NEW")).toBe("upload.result.reasons.unknown")
  })
})

describe("resolveWarningMessageKey and resolveValidationMessageKey", () => {
  it("should fall back to the generic warning for unknown codes", () => {
    expect(resolveWarningMessageKey("brand_new_code")).toBe("upload.warnings.unknown")
    expect(resolveWarningMessageKey("anonymized_study")).toBe("upload.warnings.anonymized_study")
  })

  it("should map each client validation error to its message", () => {
    expect(resolveValidationMessageKey("empty")).toBe("upload.errorEmpty")
    expect(resolveValidationMessageKey("type")).toBe("upload.errorType")
    expect(resolveValidationMessageKey("tooLarge")).toBe("upload.errorLimit")
    expect(resolveValidationMessageKey("zipTooLarge")).toBe("upload.errorZipLimit")
    expect(resolveValidationMessageKey("invalidZip")).toBe("upload.errorZipInvalid")
  })
})
