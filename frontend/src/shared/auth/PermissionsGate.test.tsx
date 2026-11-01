import { fireEvent, render, screen, waitFor } from "@testing-library/react"
import { QueryClient, QueryClientProvider } from "@tanstack/react-query"
import { PermissionsGate } from "./PermissionsGate"
import { useCapabilities } from "./capabilities_context"
import { Capability } from "./capabilities"
import { authApi } from "../../modules/auth/api"
import { useAuthStore } from "../store/auth_store"

vi.mock("../../modules/auth/api", () => ({
  authApi: {
    getPermissions: vi.fn(),
    me: vi.fn(),
    login: vi.fn(),
    logout: vi.fn(),
    register: vi.fn(),
  },
}))

vi.mock("../services/api", () => ({
  setCsrfToken: vi.fn(),
  clearCsrfToken: vi.fn(),
}))

function GrantedCapabilitiesProbe() {
  const { hasCapability } = useCapabilities()
  return <p>{hasCapability(Capability.PatientRead) ? "can read patients" : "cannot read patients"}</p>
}

const renderGate = () => {
  const queryClient = new QueryClient({ defaultOptions: { queries: { retry: false } } })
  return render(
    <QueryClientProvider client={queryClient}>
      <PermissionsGate>
        <GrantedCapabilitiesProbe />
      </PermissionsGate>
    </QueryClientProvider>,
  )
}

describe("PermissionsGate", () => {
  beforeEach(() => {
    vi.clearAllMocks()
    useAuthStore.setState({ isAuthenticated: true, userId: "user-1", role: "DOCTOR" })
  })

  it("should show a loading state until the permissions arrive", async () => {
    vi.mocked(authApi.getPermissions).mockResolvedValue({ role: "DOCTOR", capabilities: [Capability.PatientRead] })

    renderGate()

    expect(screen.getByRole("status", { name: "Loading your permissions" })).toBeInTheDocument()
    expect(screen.queryByText("can read patients")).not.toBeInTheDocument()
    expect(await screen.findByText("can read patients")).toBeInTheDocument()
  })

  it("should expose only the capabilities returned by the api", async () => {
    vi.mocked(authApi.getPermissions).mockResolvedValue({ role: "PATIENT", capabilities: [Capability.PortalRead] })

    renderGate()

    expect(await screen.findByText("cannot read patients")).toBeInTheDocument()
  })

  it("should show an alert with a retry action when the permissions fail to load", async () => {
    vi.mocked(authApi.getPermissions).mockRejectedValueOnce(new Error("network"))
    vi.mocked(authApi.getPermissions).mockResolvedValueOnce({ role: "DOCTOR", capabilities: [Capability.PatientRead] })

    renderGate()

    expect(await screen.findByRole("alert")).toHaveTextContent("Could not load your permissions")
    expect(screen.queryByText("can read patients")).not.toBeInTheDocument()

    fireEvent.click(screen.getByRole("button", { name: "Retry" }))

    await waitFor(() => expect(screen.getByText("can read patients")).toBeInTheDocument())
    expect(authApi.getPermissions).toHaveBeenCalledTimes(2)
  })
})
