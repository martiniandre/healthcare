import { describe, it, expect, vi } from "vitest"
import { render, screen } from "@testing-library/react"
import { MemoryRouter, Route, Routes } from "react-router-dom"
import { TestAbilityProvider } from "../../test/TestAbilityProvider"
import type { RoleFixtureName } from "../../test/backend_matrix_role_fixtures"
import { PatientDetails } from "./PatientDetails"

vi.mock("./queries", () => ({
  usePatientQuery: () => ({
    data: { patient_id: "1", fhir_resource_id: "patient-1", full_name: "Ana", birth_date: "1990-01-01", document_id: "1", phone_number: "1" },
    isLoading: false,
    isError: false,
    refetch: vi.fn(),
  }),
  useEncountersQuery: () => ({ data: [] }),
}))

vi.mock("../imaging/queries", () => ({
  useImagingStudiesQuery: () => ({ data: [], isLoading: false, isError: false, refetch: vi.fn() }),
}))

vi.mock("../imaging/ImagingWorkspace", () => ({ ImagingWorkspace: () => null }))
vi.mock("./components/PatientPageLayout", () => ({
  PatientPageLayout: ({ sidebarTop, children }: { sidebarTop: React.ReactNode; children: React.ReactNode }) => (
    <div>
      {sidebarTop}
      {children}
    </div>
  ),
}))
vi.mock("./components/EncounterSelectionDialog", () => ({ EncounterSelectionDialog: () => null }))
vi.mock("./components/modals/ExamAnalyzerModal", () => ({ ExamAnalyzerModal: () => null }))
vi.mock("./components/PACSStudies", () => ({ PACSStudies: () => null }))
vi.mock("./components/EncounterHistory", () => ({ default: () => <div>encounter history</div> }))
vi.mock("./components/ClinicalConditions", () => ({ default: () => null }))
vi.mock("./components/ClinicalAllergies", () => ({ default: () => null }))
vi.mock("./components/timeline/ClinicalTimeline", () => ({ default: () => <div>timeline content</div> }))

const renderDetails = (role: RoleFixtureName, url = "/patients/patient-1?tab=timeline") =>
  render(
    <TestAbilityProvider role={role}>
      <MemoryRouter initialEntries={[url]}>
        <Routes>
          <Route path="/patients/:id" element={<PatientDetails />} />
        </Routes>
      </MemoryRouter>
    </TestAbilityProvider>,
  )

describe("PatientDetails timeline tab", () => {
  it.each<RoleFixtureName>(["ADMIN", "DOCTOR", "NURSE"])("should show the timeline tab and content for %s", async (role) => {
    renderDetails(role)

    expect(screen.getByRole("tab", { name: "Timeline" })).toBeInTheDocument()
    expect(await screen.findByText("timeline content")).toBeInTheDocument()
  })

  it("should hide the timeline tab and fall back to encounters for reception", async () => {
    renderDetails("RECEPTION")

    expect(screen.queryByRole("tab", { name: "Timeline" })).toBeNull()
    expect(screen.queryByText("timeline content")).toBeNull()
    expect(await screen.findByText("encounter history")).toBeInTheDocument()
  })
})
