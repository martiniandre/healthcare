import { useQuery } from "@tanstack/react-query"
import { statsApi } from "./api"

export const statsKeys = {
  all: ["stats"] as const,
}

export const useStatsQuery = () => {
  return useQuery({
    queryKey: statsKeys.all,
    queryFn: statsApi.getStats,
  })
}
