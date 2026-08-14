import { scheduleApi } from "./api"

vi.mock("../../shared/utils/http", () => ({
  http: {
    get: vi.fn(),
    post: vi.fn(),
  },
}))

import { http } from "../../shared/utils/http"

describe("scheduleApi", () => {
  afterEach(() => {
    vi.clearAllMocks()
  })

  it("should create an appointment with the given payload", async () => {
    vi.mocked(http.post).mockResolvedValue({ id: "appointment-1" })
    const payload = {
      patient_fhir_id: "patient-1",
      staff_id: "staff-1",
      starts_at: "2026-09-01T09:00:00Z",
      ends_at: "2026-09-01T09:30:00Z",
      reason: "Retorno",
    }

    await scheduleApi.createAppointment(payload)

    expect(http.post).toHaveBeenCalledWith("/appointments", payload)
  })

  it("should list appointments by staff and date", async () => {
    vi.mocked(http.get).mockResolvedValue([])

    await scheduleApi.listByStaffOnDate("staff-1", "2026-09-01")

    expect(http.get).toHaveBeenCalledWith("/appointments?staff_id=staff-1&date=2026-09-01")
  })

  it("should list appointments by patient", async () => {
    vi.mocked(http.get).mockResolvedValue([])

    await scheduleApi.listByPatient("patient-1")

    expect(http.get).toHaveBeenCalledWith("/appointments?patient_fhir_id=patient-1")
  })

  it("should list the current user appointments", async () => {
    vi.mocked(http.get).mockResolvedValue([])

    await scheduleApi.listMyAppointments()

    expect(http.get).toHaveBeenCalledWith("/appointments/my")
  })

  it("should get a single appointment", async () => {
    vi.mocked(http.get).mockResolvedValue({ id: "appointment-1" })

    await scheduleApi.getAppointment("appointment-1")

    expect(http.get).toHaveBeenCalledWith("/appointments/appointment-1")
  })

  it("should cancel an appointment", async () => {
    vi.mocked(http.post).mockResolvedValue({ id: "appointment-1", status: "cancelled" })

    await scheduleApi.cancelAppointment("appointment-1")

    expect(http.post).toHaveBeenCalledWith("/appointments/appointment-1/cancel")
  })
})
