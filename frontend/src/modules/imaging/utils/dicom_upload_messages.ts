import { DicomInstanceRejectionCode, DicomUploadErrorCode } from "../types"
import type { DicomFileValidationError } from "./dicom_file_validation"

export interface UploadRejectionMessageKeys {
  title: string
  help: string
}

const UPLOAD_REJECTION_MESSAGE_KEYS: ReadonlyMap<string, UploadRejectionMessageKeys> = new Map([
  [
    DicomUploadErrorCode.PATIENT_MISMATCH,
    { title: "upload.errors.patientMismatchTitle", help: "upload.errors.patientMismatchHelp" },
  ],
  [DicomUploadErrorCode.INVALID_DICOM, { title: "upload.errors.invalidDicomTitle", help: "upload.errors.invalidDicomHelp" }],
  [DicomUploadErrorCode.DUPLICATE_INSTANCE, { title: "upload.errors.duplicateTitle", help: "upload.errors.duplicateHelp" }],
  [DicomUploadErrorCode.ZIP_INVALID, { title: "upload.errors.zipInvalidTitle", help: "upload.errors.zipInvalidHelp" }],
  [DicomUploadErrorCode.ZIP_UNSAFE_ENTRY, { title: "upload.errors.zipUnsafeTitle", help: "upload.errors.zipUnsafeHelp" }],
  [DicomUploadErrorCode.ZIP_LIMITS_EXCEEDED, { title: "upload.errors.zipLimitsTitle", help: "upload.errors.zipLimitsHelp" }],
  [
    DicomUploadErrorCode.NO_VALID_INSTANCES,
    { title: "upload.errors.noValidInstancesTitle", help: "upload.errors.noValidInstancesHelp" },
  ],
])

const INSTANCE_REJECTION_REASON_KEYS: ReadonlyMap<string, string> = new Map([
  [DicomInstanceRejectionCode.DUPLICATE_INSTANCE, "upload.result.reasons.duplicate"],
  [DicomInstanceRejectionCode.INVALID_DICOM, "upload.result.reasons.invalid"],
  [DicomInstanceRejectionCode.UNSUPPORTED_TRANSFER_SYNTAX, "upload.result.reasons.unsupportedTransferSyntax"],
  [DicomInstanceRejectionCode.MISSING_IDENTIFIERS, "upload.result.reasons.missingIdentifiers"],
])

const UNKNOWN_REJECTION_REASON_KEY = "upload.result.reasons.unknown"

const WARNING_MESSAGE_KEYS: ReadonlyMap<string, string> = new Map([
  ["patient_identity_missing", "upload.warnings.patient_identity_missing"],
  ["anonymized_study", "upload.warnings.anonymized_study"],
  ["duplicate_instance", "upload.warnings.duplicate_instance"],
  ["modality_overridden", "upload.warnings.modality_overridden"],
])

const UNKNOWN_WARNING_MESSAGE_KEY = "upload.warnings.unknown"

const VALIDATION_MESSAGE_KEYS: Record<DicomFileValidationError, string> = {
  empty: "upload.errorEmpty",
  type: "upload.errorType",
  tooLarge: "upload.errorLimit",
  zipTooLarge: "upload.errorZipLimit",
  invalidZip: "upload.errorZipInvalid",
}

export const resolveUploadRejectionMessageKeys = (errorCode: string | undefined): UploadRejectionMessageKeys | null => {
  return errorCode ? (UPLOAD_REJECTION_MESSAGE_KEYS.get(errorCode) ?? null) : null
}

export const resolveInstanceRejectionReasonKey = (rejectionCode: string): string => {
  return INSTANCE_REJECTION_REASON_KEYS.get(rejectionCode) ?? UNKNOWN_REJECTION_REASON_KEY
}

export const resolveWarningMessageKey = (warningCode: string): string => {
  return WARNING_MESSAGE_KEYS.get(warningCode) ?? UNKNOWN_WARNING_MESSAGE_KEY
}

export const resolveValidationMessageKey = (validationError: DicomFileValidationError): string => {
  return VALIDATION_MESSAGE_KEYS[validationError]
}
