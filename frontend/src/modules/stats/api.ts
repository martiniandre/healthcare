import { http } from "../../shared/utils/http"
import type { StatsResponse } from "./types"

export const statsApi = {
  getStats: async (): Promise<StatsResponse> => {
    return http.get<StatsResponse>("/stats")
  },
}
