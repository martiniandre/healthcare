import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query"
import { examAnalyzerApi } from "./api"

export const examAnalyzerKeys = {
  all: ["exam_analyses"] as const,
  lists: () => [...examAnalyzerKeys.all, "list"] as const,
  detail: (id: string) => [...examAnalyzerKeys.all, "detail", id] as const,
}

export const useExamAnalysesQuery = (patientFhirId?: string) => {
  return useQuery({
    queryKey: [...examAnalyzerKeys.lists(), patientFhirId || "all"],
    queryFn: () => examAnalyzerApi.getAnalyses(patientFhirId),
  })
}

export const useExamAnalysisQuery = (
  id: string,
  options?: { enabled?: boolean; refetchInterval?: number }
) => {
  return useQuery({
    queryKey: examAnalyzerKeys.detail(id),
    queryFn: () => examAnalyzerApi.getAnalysis(id),
    enabled: options?.enabled ?? !!id,
    refetchInterval: options?.refetchInterval,
  })
}

export const useUploadExamMutation = () => {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (variables: {
      file: File
      consent: boolean
      anonymize: boolean
      patientFhirId?: string
      onUploadProgress?: (progressPercentage: number) => void
    }) =>
      examAnalyzerApi.uploadExamFile(
        variables.file,
        variables.consent,
        variables.anonymize,
        variables.patientFhirId,
        variables.onUploadProgress
      ),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: examAnalyzerKeys.all })
    },
  })
}

export const useDeleteAnalysisMutation = () => {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (id: string) => examAnalyzerApi.deleteAnalysis(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: examAnalyzerKeys.all })
    },
  })
}
