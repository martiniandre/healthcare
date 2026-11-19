import { analyticsApi } from "./api"
import { dashboardApi } from "./dashboard_api"

vi.mock("../../shared/utils/http", () => ({
  http: {
    get: vi.fn(),
  },
}))

import { http } from "../../shared/utils/http"

describe("analyticsApi", () => {
  afterEach(() => {
    vi.clearAllMocks()
  })

  it("should fetch the analytics stats", async () => {
    vi.mocked(http.get).mockResolvedValue({ total_patients: 42 })

    const stats = await analyticsApi.getStats()

    expect(http.get).toHaveBeenCalledWith("/analytics")
    expect(stats).toEqual({ total_patients: 42 })
  })
})

describe("dashboardApi", () => {
  afterEach(() => {
    vi.clearAllMocks()
  })

  it("should fetch the dashboard overview", async () => {
    vi.mocked(http.get).mockResolvedValue({ cards: [] })

    await dashboardApi.getDashboard()

    expect(http.get).toHaveBeenCalledWith("/analytics/dashboard")
  })
})
