import { authApi } from "./api"

vi.mock("../../shared/utils/http", () => ({
  http: {
    get: vi.fn(),
    post: vi.fn(),
  },
}))

import { http } from "../../shared/utils/http"

describe("authApi", () => {
  afterEach(() => {
    vi.clearAllMocks()
  })

  it("should login with email and password", async () => {
    vi.mocked(http.post).mockResolvedValue({ userId: "user-1", role: "DOCTOR" })

    const session = await authApi.login({ email: "medico@clinica.com", password: "secret" })

    expect(http.post).toHaveBeenCalledWith("/auth/login", { email: "medico@clinica.com", password: "secret" })
    expect(session).toEqual({ userId: "user-1", role: "DOCTOR" })
  })

  it("should fetch the current session", async () => {
    vi.mocked(http.get).mockResolvedValue({ userId: "user-1" })

    await authApi.me()

    expect(http.get).toHaveBeenCalledWith("/auth/me")
  })

  it("should register a new user", async () => {
    vi.mocked(http.post).mockResolvedValue(undefined)

    await authApi.register({ email: "novo@mail.com", password: "secret" })

    expect(http.post).toHaveBeenCalledWith("/auth/register", { email: "novo@mail.com", password: "secret" })
  })

  it("should request logout and swallow endpoint failures", async () => {
    const consoleWarnSpy = vi.spyOn(console, "warn").mockImplementation(() => {})
    vi.mocked(http.post).mockRejectedValue(new Error("session already closed"))

    await authApi.logout()

    expect(http.post).toHaveBeenCalledWith("/auth/logout")
    consoleWarnSpy.mockRestore()
  })
})
