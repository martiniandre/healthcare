import { portalApi } from "./api"

vi.mock("../../shared/utils/http", () => ({
  http: {
    get: vi.fn(),
  },
}))

import { http } from "../../shared/utils/http"

describe("portalApi", () => {
  afterEach(() => {
    vi.clearAllMocks()
  })

  it("should fetch the portal dashboard", async () => {
    vi.mocked(http.get).mockResolvedValue({ patient_info: { full_name: "Maria" } })

    await portalApi.getDashboard()

    expect(http.get).toHaveBeenCalledWith("/portal/dashboard")
  })

  it.each([
    ["getEncounters", "/portal/encounters"],
    ["getObservations", "/portal/observations"],
    ["getConditions", "/portal/conditions"],
    ["getMedications", "/portal/medications"],
    ["getReports", "/portal/reports"],
    ["getImaging", "/portal/imaging"],
  ])("%s should call %s", async (methodName, expectedUrl) => {
    vi.mocked(http.get).mockResolvedValue([])

    await portalApi[methodName as keyof typeof portalApi]()

    expect(http.get).toHaveBeenCalledWith(expectedUrl)
  })
})
