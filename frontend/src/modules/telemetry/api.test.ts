import { telemetryApi } from "./api"

vi.mock("../../shared/utils/http", () => ({
  http: {
    get: vi.fn(),
    post: vi.fn(),
  },
}))

import { http } from "../../shared/utils/http"

describe("telemetryApi", () => {
  afterEach(() => {
    vi.clearAllMocks()
  })

  it("should list telemetry rooms", async () => {
    vi.mocked(http.get).mockResolvedValue([])

    await telemetryApi.getRooms()

    expect(http.get).toHaveBeenCalledWith("/telemetry/rooms")
  })

  it("should unlock a room with the provided passcode", async () => {
    vi.mocked(http.post).mockResolvedValue({ unlocked: true })

    await telemetryApi.unlockRoom("room-1", "1234")

    expect(http.post).toHaveBeenCalledWith("/telemetry/rooms/room-1/unlock", { passcode: "1234" })
  })

  it("should list the beds of a room", async () => {
    vi.mocked(http.get).mockResolvedValue([])

    await telemetryApi.getBeds("room-1")

    expect(http.get).toHaveBeenCalledWith("/telemetry/rooms/room-1/beds")
  })

  it("should update a bed condition with the vital signs payload", async () => {
    vi.mocked(http.post).mockResolvedValue({ ok: true })

    await telemetryApi.updateBedCondition({
      bedId: "bed-1",
      bpm: 88,
      spo2: 96,
      temperature: 36.8,
      status: "normal",
      condition: "Normal",
    })

    expect(http.post).toHaveBeenCalledWith("/telemetry/beds/bed-1/condition", {
      bpm: 88,
      spo2: 96,
      temperature: 36.8,
      status: "normal",
      condition: "Normal",
    })
  })
})
