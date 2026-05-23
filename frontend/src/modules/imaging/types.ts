export interface ImagingStudy {
  id: string
  patient_fhir_id: string
  title: string
  modality: string
  study_instance_uid: string
  status: string
  download_url?: string
  created_at: string
}

export interface UploadImagingStudyPayload {
  patientFhirId: string
  title: string
  modality: string
  dicomBlob: Blob
}
