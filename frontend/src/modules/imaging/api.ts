import { http } from "../../shared/utils/http"
import type { DicomUploadResult, ImagingStudy, UploadImagingStudyPayload } from "./types"

export const imagingApi = {
  getImagingStudies: async (patientFhirId: string): Promise<ImagingStudy[]> => {
    return http.get<ImagingStudy[]>(`/patients/${patientFhirId}/studies`)
  },

  getImagingStudy: async (imagingStudyId: string): Promise<ImagingStudy> => {
    return http.get<ImagingStudy>(`/studies/${imagingStudyId}`)
  },

  uploadImagingStudy: async (payload: UploadImagingStudyPayload): Promise<DicomUploadResult> => {
    const uploadFormData = new FormData()
    if (payload.title) {
      uploadFormData.append("title", payload.title)
    }
    if (payload.modality) {
      uploadFormData.append("modality", payload.modality)
    }
    uploadFormData.append("file", payload.dicomBlob, "study.dcm")

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
