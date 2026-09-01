import { renderHook, act, waitFor } from "@testing-library/react"
import { useAuthInit } from "./useAuthInit"
import { authApi } from "../../modules/auth/api"
import { useAuthStore } from "../store/auth_store"

const csrfApi = vi.hoisted(() => ({
  setCsrfToken: vi.fn(),
  clearCsrfToken: vi.fn(),
}))

vi.mock("../services/api", () => csrfApi)

vi.mock("../../modules/auth/api", () => ({
  authApi: {
    me: vi.fn(),
    login: vi.fn(),
    logout: vi.fn(),
    register: vi.fn(),
  },
}))

describe("useAuthInit", () => {
  beforeEach(() => {
    useAuthStore.setState({
      isAuthenticated: false,
      userId: null,
      role: null,
      email: null,
      fullName: null,
      isActive: null,
    })
    vi.clearAllMocks()
  })

  it("should restore the session when me returns a valid user", async () => {
    vi.mocked(authApi.me).mockResolvedValue({
      userId: "user-1",
      role: "DOCTOR",
      email: "medico@clinica.com",
      fullName: "Médico Teste",
      isActive: true,
      csrfToken: "csrf-token-from-me",
    })

    const { result } = renderHook(() => useAuthInit())

    expect(result.current.isLoading).toBe(true)

    await waitFor(() => {
      expect(result.current.isLoading).toBe(false)
    })

    expect(useAuthStore.getState().isAuthenticated).toBe(true)
    expect(useAuthStore.getState().userId).toBe("user-1")
    expect(useAuthStore.getState().role).toBe("DOCTOR")
    expect(useAuthStore.getState().fullName).toBe("Médico Teste")
    expect(csrfApi.setCsrfToken).toHaveBeenCalledWith("csrf-token-from-me")
  })

  it("should finish loading without authenticating when me fails", async () => {
    vi.mocked(authApi.me).mockRejectedValue(new Error("session expired"))

    const { result } = renderHook(() => useAuthInit())

    await waitFor(() => {
      expect(result.current.isLoading).toBe(false)
    })

    expect(useAuthStore.getState().isAuthenticated).toBe(false)
    expect(useAuthStore.getState().userId).toBeNull()
  })

  it("should finish loading without a csrf token when me omits it", async () => {
    vi.mocked(authApi.me).mockResolvedValue({
      userId: "user-3",
      role: "RECEPTION",
      email: "recepcionista@hospital.com",
    })

    const { result } = renderHook(() => useAuthInit())

    await waitFor(() => {
      expect(result.current.isLoading).toBe(false)
    })

    expect(useAuthStore.getState().isAuthenticated).toBe(true)
    expect(csrfApi.setCsrfToken).not.toHaveBeenCalled()
  })

  it("should not login after unmount when the request resolves late", async () => {
    let resolveMe: (value: unknown) => void = () => {}
    vi.mocked(authApi.me).mockImplementation(
      () => new Promise((resolve) => {
        resolveMe = resolve
      })
    )

    const { unmount } = renderHook(() => useAuthInit())
    unmount()

    await act(async () => {
      resolveMe({
        userId: "user-2",
        role: "NURSE",
        email: "enfermeiro@hospital.com",
        isActive: true,
      })
      await Promise.resolve()
    })

    expect(useAuthStore.getState().isAuthenticated).toBe(false)
  })
})
