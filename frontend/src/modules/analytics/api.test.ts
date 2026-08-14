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

  it("should fetch consultations per doctor", async () => {
    vi.mocked(http.get).mockResolvedValue([])

    await dashboardApi.getConsultationsPerDoctor()

    expect(http.get).toHaveBeenCalledWith("/analytics/dashboard/consultations-per-doctor")
  })

  it("should fetch the occupancy rate", async () => {
    vi.mocked(http.get).mockResolvedValue({ rate: 0.8, total_beds: 10, occupied_beds: 8 })

    await dashboardApi.getOccupancyRate()

    expect(http.get).toHaveBeenCalledWith("/analytics/dashboard/occupancy-rate")
  })

  it("should fetch the average wait time", async () => {
    vi.mocked(http.get).mockResolvedValue({ average_minutes: 15, by_department: [] })

    await dashboardApi.getAvgWaitTime()

    expect(http.get).toHaveBeenCalledWith("/analytics/dashboard/avg-wait-time")
  })

  it("should fetch the top diagnoses", async () => {
    vi.mocked(http.get).mockResolvedValue([])

    await dashboardApi.getTopDiagnoses()

    expect(http.get).toHaveBeenCalledWith("/analytics/dashboard/top-diagnoses")
  })
})
