import { http } from "../../shared/utils/http"
import { api } from "../../shared/services/api"
import type { ExamAnalysis } from "./types"

export const examAnalyzerApi = {
  getAnalyses: async (patientFhirId?: string): Promise<ExamAnalysis[]> => {
    const targetUrl = patientFhirId ? `/exam-analyses?patientFhirId=${patientFhirId}` : "/exam-analyses"
    return http.get<ExamAnalysis[]>(targetUrl)
  },

  getAnalysis: async (id: string): Promise<ExamAnalysis> => {
    return http.get<ExamAnalysis>(`/exam-analyses/${id}`)
  },

  uploadExamFile: async (
    file: File,
    consent: boolean,
    anonymize: boolean,
    patientFhirId?: string,
    onUploadProgress?: (progressPercentage: number) => void
  ): Promise<ExamAnalysis> => {
    const formData = new FormData()
    formData.append("file", file)
    formData.append("consent", consent ? "true" : "false")
    formData.append("anonymize", anonymize ? "true" : "false")
    if (patientFhirId) {
      formData.append("patientFhirId", patientFhirId)
    }

    return http.post<ExamAnalysis, FormData>("/exam-analyses", formData, {
      headers: {
        "Content-Type": "multipart/form-data",
      },
      onUploadProgress: (progressEvent) => {
        if (onUploadProgress && progressEvent.total) {
          const percentage = Math.round((progressEvent.loaded * 100) / progressEvent.total)
          onUploadProgress(percentage)
        }
      },
    })
  },

  deleteAnalysis: async (id: string): Promise<{ success: string }> => {
    return api.delete<{ success: string }>(`/exam-analyses/${id}`).then((responseData) => responseData.data)
  },
}
