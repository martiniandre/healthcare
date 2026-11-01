import { DicomModality, ImagingStudyStatus } from "../../shared/types"

export interface ImagingStudy {
  id: string
  patient_fhir_id: string
  title: string
  modality: DicomModality
  study_instance_uid: string
  status: ImagingStudyStatus
  download_url?: string
  created_at: string
  studyInstanceUid?: string
  studyDate?: string
  accessionNumber?: string
  seriesInstanceUid?: string
  sopInstanceUid?: string
}

export interface ImagingStudyInstance {
  id: string
  sopInstanceUid: string
  seriesInstanceUid: string
  instanceNumber?: number | null
  fileName: string
  downloadUrl: string
  expiresAt: string
}

export const DicomUploadWarningCode = {
  PATIENT_IDENTITY_MISSING: "patient_identity_missing",
  ANONYMIZED_STUDY: "anonymized_study",
  DUPLICATE_INSTANCE: "duplicate_instance",
  MODALITY_OVERRIDDEN: "modality_overridden",
} as const

export type DicomUploadWarningCode = typeof DicomUploadWarningCode[keyof typeof DicomUploadWarningCode]

export interface DicomUploadWarning {
  code: string
  message?: string
}

export const DicomInstanceRejectionCode = {
  DUPLICATE_INSTANCE: "DICOM_DUPLICATE_INSTANCE",
  INVALID_DICOM: "INVALID_DICOM",
  UNSUPPORTED_TRANSFER_SYNTAX: "DICOM_UNSUPPORTED_TRANSFER_SYNTAX",
  MISSING_IDENTIFIERS: "DICOM_MISSING_IDENTIFIERS",
} as const

export type DicomInstanceRejectionCode = typeof DicomInstanceRejectionCode[keyof typeof DicomInstanceRejectionCode]

export interface DicomRejectedInstance {
  fileName: string
  code: string
  reason?: string
}

export interface DicomUploadResult extends ImagingStudy {
  warnings?: DicomUploadWarning[]
  studyId?: string
  studyCreated?: boolean
  instancesAccepted?: number
  instancesRejected?: DicomRejectedInstance[]
}

export const DicomUploadErrorCode = {
  PATIENT_MISMATCH: "DICOM_PATIENT_MISMATCH",
  INVALID_DICOM: "INVALID_DICOM",
  DUPLICATE_INSTANCE: "DICOM_DUPLICATE_INSTANCE",
  ZIP_INVALID: "ZIP_INVALID",
  ZIP_UNSAFE_ENTRY: "ZIP_UNSAFE_ENTRY",
  ZIP_LIMITS_EXCEEDED: "ZIP_LIMITS_EXCEEDED",
  NO_VALID_INSTANCES: "DICOM_NO_VALID_INSTANCES",
} as const

export type DicomUploadErrorCode = typeof DicomUploadErrorCode[keyof typeof DicomUploadErrorCode]

export interface UploadImagingStudyPayload {
  patientFhirId: string
  title?: string
  modality?: DicomModality
  dicomBlob: Blob
  onUploadProgress?: (progressPercentage: number) => void
  signal?: AbortSignal
}
