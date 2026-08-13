import { http } from "../../shared/utils/http"
import type { ImagingStudy, UploadImagingStudyPayload } from "./types"

export const imagingApi = {
  getImagingStudies: async (patientFhirId: string): Promise<ImagingStudy[]> => {
    return http.get<ImagingStudy[]>(`/patients/${patientFhirId}/studies`)
  },

  getImagingStudy: async (imagingStudyId: string): Promise<ImagingStudy> => {
    return http.get<ImagingStudy>(`/studies/${imagingStudyId}`)
  },

  uploadImagingStudy: async (payload: UploadImagingStudyPayload): Promise<ImagingStudy> => {
    const uploadFormData = new FormData()
    uploadFormData.append("title", payload.title)
    uploadFormData.append("modality", payload.modality)
    uploadFormData.append("file", payload.dicomBlob, "study.dcm")

    return http.post<ImagingStudy>(`/patients/${payload.patientFhirId}/studies`, uploadFormData)
  },
}
