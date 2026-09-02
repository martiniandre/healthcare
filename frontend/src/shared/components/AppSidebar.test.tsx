import { describe, it, expect, beforeEach, afterEach, vi } from "vitest"
import { render, screen, fireEvent, waitFor } from "@testing-library/react"
import { AppSidebar } from "./AppSidebar"
import { useAuthStore } from "../store/auth_store"
import { useLayoutStore } from "../store/layout_store"

const mockTranslateFunction = (key: string) => key
const mockChangeLanguage = vi.fn()
const mockNavigate = vi.fn()
const mockLocation = { pathname: "/", search: "", hash: "", state: null, key: "default" }

vi.mock("react-i18next", () => ({
  useTranslation: () => ({
    t: mockTranslateFunction,
    i18n: {
      language: "en-US",
      changeLanguage: mockChangeLanguage,
    },
  }),
}))

vi.mock("react-router-dom", () => ({
  useNavigate: () => mockNavigate,
  useLocation: () => mockLocation,
}))

vi.mock("../../modules/auth/api", () => ({
  authApi: {
    login: vi.fn(),
    me: vi.fn(),
    logout: vi.fn().mockResolvedValue(undefined),
    register: vi.fn(),
  },
}))

vi.mock("../services/api", () => ({
  clearCsrfToken: vi.fn(),
}))

describe("AppSidebar", () => {
  afterEach(() => {
    vi.unstubAllEnvs()
  })

  beforeEach(() => {
    window.localStorage.clear()
    mockNavigate.mockReset()
    mockLocation.pathname = "/"
    useAuthStore.setState({
      isAuthenticated: true,
      userId: "user-1",
      role: "DOCTOR",
      email: "doctor@hospital.com",
      fullName: "Dr. Test",
      isActive: true,
    })
    useLayoutStore.setState({ isMobileSidebarOpen: false, theme: "light" })
  })

  it("should render navigation grouped by clinical topics for staff roles", () => {
    render(<AppSidebar />)

    expect(screen.getByRole("heading", { name: "topics.clinical" })).toBeInTheDocument()
    expect(screen.getByRole("heading", { name: "topics.diagnostics" })).toBeInTheDocument()
    expect(screen.getByRole("heading", { name: "topics.operations" })).toBeInTheDocument()
    expect(screen.queryByRole("heading", { name: "topics.system" })).not.toBeInTheDocument()
    expect(screen.queryByRole("heading", { name: "topics.patient" })).not.toBeInTheDocument()

    expect(screen.getByRole("button", { name: "patients" })).toBeInTheDocument()
    expect(screen.getByRole("button", { name: "dashboard" })).toBeInTheDocument()
    expect(screen.getByRole("button", { name: "schedule" })).toBeInTheDocument()
    expect(screen.getByRole("button", { name: "telemetry" })).toBeInTheDocument()
    expect(screen.getByRole("button", { name: "examAnalyzer" })).toBeInTheDocument()
    expect(screen.getByRole("button", { name: "analytics" })).toBeInTheDocument()
    expect(screen.getByRole("button", { name: "staffManagement" })).toBeInTheDocument()
    expect(screen.getByRole("button", { name: "appearance" })).toBeInTheDocument()
    expect(screen.getByRole("button", { name: "English" })).toBeInTheDocument()

    expect(screen.queryByRole("button", { name: "portal" })).not.toBeInTheDocument()
    expect(screen.queryByRole("button", { name: "auditLogs" })).not.toBeInTheDocument()
    expect(screen.queryByRole("button", { name: "settings" })).not.toBeInTheDocument()
  })

  it("should show audit logs topic item only for admin users", () => {
    useAuthStore.setState({ role: "DOCTOR" })
    const { unmount } = render(<AppSidebar />)
    expect(screen.queryByRole("button", { name: "auditLogs" })).not.toBeInTheDocument()
    unmount()

    useAuthStore.setState({ role: "ADMIN" })
    render(<AppSidebar />)
    expect(screen.getByRole("button", { name: "auditLogs" })).toBeInTheDocument()
  })

  it("should render only the patient access topic for patient roles", () => {
    useAuthStore.setState({ role: "PATIENT" })
    render(<AppSidebar />)

    expect(screen.getByRole("heading", { name: "topics.patient" })).toBeInTheDocument()
    expect(screen.getByRole("button", { name: "portal" })).toBeInTheDocument()

    expect(screen.queryByRole("heading", { name: "topics.clinical" })).not.toBeInTheDocument()
    expect(screen.queryByRole("heading", { name: "topics.diagnostics" })).not.toBeInTheDocument()
    expect(screen.queryByRole("heading", { name: "topics.operations" })).not.toBeInTheDocument()
    expect(screen.queryByRole("button", { name: "patients" })).not.toBeInTheDocument()
    expect(screen.queryByRole("button", { name: "settings" })).not.toBeInTheDocument()
  })

  it("should mark the active navigation item with aria-current", () => {
    mockLocation.pathname = "/analytics"
    render(<AppSidebar />)

    expect(screen.getByRole("button", { name: "analytics" })).toHaveAttribute("aria-current", "page")
    expect(screen.getByRole("button", { name: "patients" })).not.toHaveAttribute("aria-current")
  })

  it("should navigate on click and close the mobile sidebar", () => {
    useLayoutStore.setState({ isMobileSidebarOpen: true })
    render(<AppSidebar />)

    fireEvent.click(screen.getByRole("button", { name: "telemetry" }))
    expect(mockNavigate).toHaveBeenCalledWith("/telemetry")
    expect(useLayoutStore.getState().isMobileSidebarOpen).toBe(false)
  })

  it("should toggle the dark theme via the appearance control", () => {
    render(<AppSidebar />)

    const appearanceToggle = screen.getByRole("button", { name: "appearance" })
    expect(appearanceToggle).toHaveAttribute("aria-pressed", "false")
    expect(useLayoutStore.getState().theme).toBe("light")

    fireEvent.click(appearanceToggle)

    expect(useLayoutStore.getState().theme).toBe("dark")
    expect(appearanceToggle).toHaveAttribute("aria-pressed", "true")
    expect(window.localStorage.getItem("healthcare.theme")).toBe("dark")
  })

  it("should hide the schedule item in production mode", () => {
    vi.stubEnv("PROD", true)

    render(<AppSidebar />)

    expect(screen.queryByRole("button", { name: "schedule" })).not.toBeInTheDocument()
    expect(screen.getByRole("button", { name: "patients" })).toBeInTheDocument()
    expect(screen.getByRole("button", { name: "dashboard" })).toBeInTheDocument()
    expect(screen.getByRole("button", { name: "telemetry" })).toBeInTheDocument()
  })

  it("should logout and close the sidebar", async () => {
    useLayoutStore.setState({ isMobileSidebarOpen: true })
    render(<AppSidebar />)

    fireEvent.click(screen.getByRole("button", { name: "logout" }))
    expect(useLayoutStore.getState().isMobileSidebarOpen).toBe(false)
    await waitFor(() => {
      expect(useAuthStore.getState().isAuthenticated).toBe(false)
    })
  })
})