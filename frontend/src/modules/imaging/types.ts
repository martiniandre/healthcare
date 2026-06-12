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
}

export interface UploadImagingStudyPayload {
  patientFhirId: string
  title: string
  modality: DicomModality
  dicomBlob: Blob
}
