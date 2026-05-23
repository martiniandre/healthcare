import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query"
import { imagingApi } from "./api"
import type { UploadImagingStudyPayload } from "./types"

export const imagingQueryKeys = {
  all: ["imaging"] as const,
  studies: (patientFhirId: string) => [...imagingQueryKeys.all, "studies", patientFhirId] as const,
  study: (imagingStudyId: string) => [...imagingQueryKeys.all, "study", imagingStudyId] as const,
}

export const useImagingStudyQuery = (studyId: string) => {
  return useQuery({
    queryKey: imagingQueryKeys.study(studyId),
    queryFn: () => imagingApi.getImagingStudy(studyId),
    enabled: !!studyId,
  })
}

export const useImagingStudiesQuery = (patientFhirId: string) => {
  return useQuery({
    queryKey: imagingQueryKeys.studies(patientFhirId),
    queryFn: () => imagingApi.getImagingStudies(patientFhirId),
    enabled: !!patientFhirId,
  })
}

export const useUploadImagingStudyMutation = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (payload: UploadImagingStudyPayload) => imagingApi.uploadImagingStudy(payload),
    onSuccess: (createdStudy, variables) => {
      queryClient.setQueryData(imagingQueryKeys.study(createdStudy.id), createdStudy)
      queryClient.invalidateQueries({ queryKey: imagingQueryKeys.studies(variables.patientFhirId) })
    },
  })
}
