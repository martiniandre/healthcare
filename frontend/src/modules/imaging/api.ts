import { http } from "../../shared/utils/http"
import type { DicomUploadResult, ImagingStudy, ImagingStudyInstance, UploadImagingStudyPayload } from "./types"

const DEFAULT_DICOM_FILE_NAME = "study.dcm"

const resolveUploadFileName = (dicomBlob: Blob): string => {
  return dicomBlob instanceof File && dicomBlob.name !== "" ? dicomBlob.name : DEFAULT_DICOM_FILE_NAME
}

export const imagingApi = {
  getImagingStudies: async (patientFhirId: string): Promise<ImagingStudy[]> => {
    return http.get<ImagingStudy[]>(`/patients/${patientFhirId}/studies`)
  },

  getImagingStudy: async (imagingStudyId: string): Promise<ImagingStudy> => {
    return http.get<ImagingStudy>(`/studies/${imagingStudyId}`)
  },

  getImagingStudyInstances: async (imagingStudyId: string): Promise<ImagingStudyInstance[]> => {
    return http.get<ImagingStudyInstance[]>(`/studies/${imagingStudyId}/instances`)
  },

  uploadImagingStudy: async (payload: UploadImagingStudyPayload): Promise<DicomUploadResult> => {
    const uploadFormData = new FormData()
    if (payload.title) {
      uploadFormData.append("title", payload.title)
    }
    if (payload.modality) {
      uploadFormData.append("modality", payload.modality)
    }
    uploadFormData.append("file", payload.dicomBlob, resolveUploadFileName(payload.dicomBlob))

    return http.post<DicomUploadResult>(`/patients/${payload.patientFhirId}/studies`, uploadFormData, {
      signal: payload.signal,
      onUploadProgress: (progressEvent) => {
        if (payload.onUploadProgress && progressEvent.total) {
          payload.onUploadProgress(Math.round((progressEvent.loaded * 100) / progressEvent.total))
        }
      },
    })
  },
}
