import { http } from "../../shared/utils/http"
import type { ExamAnalysis, UploadExamPayload } from "./types"

export const examAnalyzerApi = {
  getAnalyses: async (patientFhirId?: string): Promise<ExamAnalysis[]> => {
    const targetUrl = patientFhirId ? `/exam-analyses?patientFhirId=${patientFhirId}` : "/exam-analyses"
    return http.get<ExamAnalysis[]>(targetUrl)
  },

  getAnalysis: async (id: string): Promise<ExamAnalysis> => {
    return http.get<ExamAnalysis>(`/exam-analyses/${id}`)
  },

  uploadExamFile: async (payload: UploadExamPayload): Promise<ExamAnalysis> => {
    const formData = new FormData()
    formData.append("file", payload.file)
    formData.append("consent", payload.consent ? "true" : "false")
    formData.append("anonymize", payload.anonymize ? "true" : "false")
    if (!payload.anonymize && payload.acknowledgeExternalProcessing) {
      formData.append("acknowledgeExternalProcessing", "true")
    }
    if (payload.patientFhirId) {
      formData.append("patientFhirId", payload.patientFhirId)
    }

    return http.post<ExamAnalysis, FormData>("/exam-analyses", formData, {
      signal: payload.signal,
      onUploadProgress: (progressEvent) => {
        if (payload.onUploadProgress && progressEvent.total) {
          const percentage = Math.round((progressEvent.loaded * 100) / progressEvent.total)
          payload.onUploadProgress(percentage)
        }
      },
    })
  },

  deleteAnalysis: async (id: string): Promise<{ success: string }> => {
    return http.delete<{ success: string }>(`/exam-analyses/${id}`)
  },
}
