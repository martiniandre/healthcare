import { http } from "./http"

vi.mock("../services/api", () => ({
  api: {
    get: vi.fn(),
    post: vi.fn(),
    put: vi.fn(),
    delete: vi.fn(),
  },
}))

import { api } from "../services/api"

describe("http helper", () => {
  afterEach(() => {
    vi.clearAllMocks()
  })

  it("should return the data payload of a get request", async () => {
    vi.mocked(api.get).mockResolvedValue({ data: { count: 3 } } as never)

    const result = await http.get<{ count: number }>("/patients")

    expect(api.get).toHaveBeenCalledWith("/patients", undefined)
    expect(result).toEqual({ count: 3 })
  })

  it("should forward the configuration on get requests", async () => {
    vi.mocked(api.get).mockResolvedValue({ data: [] } as never)
    const requestConfig = { params: { page: 1 } }

    await http.get("/patients", requestConfig)

    expect(api.get).toHaveBeenCalledWith("/patients", requestConfig)
  })

  it("should return the data payload of a post request with body", async () => {
    vi.mocked(api.post).mockResolvedValue({ data: { id: "report-1" } } as never)

    const result = await http.post<{ id: string }>("/reports", { conclusion: "ok" })

    expect(api.post).toHaveBeenCalledWith("/reports", { conclusion: "ok" }, undefined)
    expect(result).toEqual({ id: "report-1" })
  })

  it("should return the data payload of a put request", async () => {
    vi.mocked(api.put).mockResolvedValue({ data: { status: "final" } } as never)

    const result = await http.put<{ status: string }>("/reports/1", { status: "final" })

    expect(api.put).toHaveBeenCalledWith("/reports/1", { status: "final" }, undefined)
    expect(result).toEqual({ status: "final" })
  })

  it("should return the data payload of a delete request", async () => {
    vi.mocked(api.delete).mockResolvedValue({ data: { success: true } } as never)

    const result = await http.delete<{ success: boolean }>("/exam-analyses/1")

    expect(api.delete).toHaveBeenCalledWith("/exam-analyses/1", undefined)
    expect(result).toEqual({ success: true })
  })
})
