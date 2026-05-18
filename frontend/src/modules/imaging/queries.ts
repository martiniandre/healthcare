import { useQuery } from "@tanstack/react-query"
import { clinicApi } from "../../shared/utils/api_client"

export const useImagingStudyQuery = (studyId: string) => {
  return useQuery({
    queryKey: ["study", studyId],
    queryFn: () => clinicApi.getImagingStudy(studyId),
    enabled: !!studyId,
  })
}
