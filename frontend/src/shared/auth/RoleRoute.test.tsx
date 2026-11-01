import { render, screen } from "@testing-library/react"
import { MemoryRouter, Route, Routes } from "react-router-dom"
import { RoleRoute } from "./RoleRoute"
import { Capability } from "./capabilities"
import { TestAbilityProvider } from "../../test/TestAbilityProvider"

const renderRoute = (capabilities: readonly string[], initialPath: string) =>
  render(
    <TestAbilityProvider capabilities={capabilities}>
      <MemoryRouter initialEntries={[initialPath]}>
        <Routes>
          <Route path="/" element={<RoleRoute requiredCapability={Capability.PatientRead}><p>patients page</p></RoleRoute>} />
          <Route path="/portal" element={<RoleRoute requiredCapability={Capability.PortalRead}><p>portal page</p></RoleRoute>} />
        </Routes>
      </MemoryRouter>
    </TestAbilityProvider>,
  )

describe("RoleRoute", () => {
  it("should render the page when the capability is granted", () => {
    renderRoute([Capability.PatientRead], "/")

    expect(screen.getByText("patients page")).toBeInTheDocument()
  })

  it("should redirect to the portal when the capability is missing and the portal is granted", () => {
    renderRoute([Capability.PortalRead], "/")

    expect(screen.getByText("portal page")).toBeInTheDocument()
  })

  it("should show a forbidden alert instead of redirecting forever when nothing is granted", () => {
    renderRoute([], "/")

    expect(screen.getByRole("alert")).toHaveTextContent("You do not have access to this area.")
  })
})
