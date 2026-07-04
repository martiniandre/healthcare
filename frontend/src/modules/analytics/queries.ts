import { useQuery } from "@tanstack/react-query"
import { analyticsApi } from "./api"

export const analyticsKeys = {
  all: ["analytics"] as const,
}

export const useStatsQuery = () => {
  return useQuery({
    queryKey: analyticsKeys.all,
    queryFn: analyticsApi.getStats,
  })
}
