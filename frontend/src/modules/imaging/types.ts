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

export interface DicomUploadResult extends ImagingStudy {
  warnings?: DicomUploadWarning[]
}

export const DicomUploadErrorCode = {
  PATIENT_MISMATCH: "DICOM_PATIENT_MISMATCH",
  INVALID_DICOM: "INVALID_DICOM",
  DUPLICATE_INSTANCE: "DICOM_DUPLICATE_INSTANCE",
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
