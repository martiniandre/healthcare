import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query"
import { imagingApi } from "./api"
import type { UploadImagingStudyPayload } from "./types"

export const imagingQueryKeys = {
  all: ["imaging"] as const,
  studies: (patientFhirId: string) => [...imagingQueryKeys.all, "studies", patientFhirId] as const,
  study: (imagingStudyId: string) => [...imagingQueryKeys.all, "study", imagingStudyId] as const,
  instances: (imagingStudyId: string) => [...imagingQueryKeys.all, "instances", imagingStudyId] as const,
}

export const useImagingStudyQuery = (studyId: string) => {
  return useQuery({
    queryKey: imagingQueryKeys.study(studyId),
    queryFn: () => imagingApi.getImagingStudy(studyId),
    enabled: !!studyId,
  })
}

export const useImagingStudyInstancesQuery = (studyId: string) => {
  return useQuery({
    queryKey: imagingQueryKeys.instances(studyId),
    queryFn: () => imagingApi.getImagingStudyInstances(studyId),
    enabled: !!studyId,
  })
}

export const useImagingStudiesQuery = (patientFhirId: string, enableQuery?: boolean) => {
  return useQuery({
    queryKey: imagingQueryKeys.studies(patientFhirId),
    queryFn: () => imagingApi.getImagingStudies(patientFhirId),
    enabled: !!patientFhirId && enableQuery !== false,
  })
}

export const useUploadImagingStudyMutation = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (payload: UploadImagingStudyPayload) => imagingApi.uploadImagingStudy(payload),
    onSuccess: (uploadedStudy, variables) => {
      const targetStudyId = uploadedStudy.studyId ?? uploadedStudy.id
      queryClient.setQueryData(imagingQueryKeys.study(targetStudyId), uploadedStudy)
      queryClient.invalidateQueries({ queryKey: imagingQueryKeys.studies(variables.patientFhirId) })
      queryClient.invalidateQueries({ queryKey: imagingQueryKeys.instances(targetStudyId) })
    },
  })
}
