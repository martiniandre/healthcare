import { renderHook } from "@testing-library/react"
import { MemoryRouter } from "react-router-dom"
import { usePageViewLogger } from "./usePageViewLogger"
import { auditLogsApi } from "../../modules/audit_logs/api"
import { useAuthStore } from "../store/auth_store"

vi.mock("../../modules/audit_logs/api", () => ({
  auditLogsApi: {
    createAuditLog: vi.fn(() => Promise.resolve()),
    listAuditLogs: vi.fn(),
  },
}))

describe("usePageViewLogger", () => {
  beforeEach(() => {
    useAuthStore.setState({ isAuthenticated: false })
    vi.clearAllMocks()
  })

  function renderLogger(initialPath: string) {
    return renderHook(() => usePageViewLogger(), {
      wrapper: ({ children }) => <MemoryRouter initialEntries={[initialPath]}>{children}</MemoryRouter>,
    })
  }

  it("should log a page view when the user is authenticated", () => {
    useAuthStore.setState({ isAuthenticated: true })
    vi.mocked(auditLogsApi.createAuditLog).mockResolvedValue({} as never)

    renderLogger("/patients?page=2")

    expect(auditLogsApi.createAuditLog).toHaveBeenCalledWith(
      expect.objectContaining({
        method: "PAGE_VIEW",
        correlation_id: "Viewed page: /patients?page=2",
        access_granted: true,
      })
    )
  })

  it("should not log page views for unauthenticated users", () => {
    useAuthStore.setState({ isAuthenticated: false })

    renderLogger("/login")

    expect(auditLogsApi.createAuditLog).not.toHaveBeenCalled()
  })

  it("should swallow logging errors gracefully", async () => {
    useAuthStore.setState({ isAuthenticated: true })
    const consoleErrorSpy = vi.spyOn(console, "error").mockImplementation(() => {})
    vi.mocked(auditLogsApi.createAuditLog).mockRejectedValue(new Error("audit api down"))

    renderLogger("/dashboard")

    await vi.waitFor(() => {
      expect(consoleErrorSpy).toHaveBeenCalledWith("Failed to log page view", expect.any(Error))
    })
    consoleErrorSpy.mockRestore()
  })
})
