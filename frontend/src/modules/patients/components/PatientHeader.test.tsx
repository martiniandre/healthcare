import { describe, it, expect, vi, beforeEach } from "vitest"
import { render, screen, fireEvent } from "@testing-library/react"
import { PatientHeader } from "./PatientHeader"

vi.mock("../queries", () => ({
  usePatientAllergiesQuery: vi.fn(),
  usePatientConditionsQuery: vi.fn(),
}))

vi.mock("../../../shared/auth/useAbility", () => ({
  useAbility: () => ({ can: () => true }),
}))

import { usePatientAllergiesQuery, usePatientConditionsQuery } from "../queries"

const mockedAllergies = vi.mocked(usePatientAllergiesQuery)
const mockedConditions = vi.mocked(usePatientConditionsQuery)

const patient = {
  patient_id: "1",
  fhir_resource_id: "fhir-1",
  full_name: "Maria Silva",
  birth_date: "1990-01-01",
  document_id: "123.456.789-00",
  phone_number: "11999999999",
}

const queryResult = (overrides: Record<string, unknown>) =>
  ({ data: [], isLoading: false, isError: false, refetch: vi.fn(), ...overrides }) as never

describe("PatientHeader", () => {
  beforeEach(() => {
    mockedAllergies.mockReturnValue(queryResult({}))
    mockedConditions.mockReturnValue(queryResult({}))
  })

  it("should render name, document and a sticky container", () => {
    render(<PatientHeader patient={patient} onBack={vi.fn()} />)
    expect(screen.getByText("Maria Silva")).toBeInTheDocument()
    expect(screen.getAllByText(/123\.456\.789-00/).length).toBeGreaterThan(0)
    expect(screen.getByTestId("patient-sticky-header").className).toContain("sticky")
  })

  it("should list active allergies and conditions with counts", () => {
    mockedAllergies.mockReturnValue(
      queryResult({ data: [{ allergen_display: "Penicillin", clinical_status: "active" }, { allergen_display: "Latex", clinical_status: "resolved" }] })
    )
    mockedConditions.mockReturnValue(
      queryResult({ data: [{ code_display: "Diabetes", clinical_status: "active" }] })
    )
    render(<PatientHeader patient={patient} onBack={vi.fn()} />)
    expect(screen.getByText("Penicillin")).toBeInTheDocument()
    expect(screen.queryByText("Latex")).toBeNull()
    expect(screen.getByText("Allergies (1)")).toBeInTheDocument()
    expect(screen.getByText("Diabetes")).toBeInTheDocument()
    expect(screen.getByText("Allergies (1)")).toBeInTheDocument()
  })

  it("should show a neutral badge when there are no known allergies", () => {
    render(<PatientHeader patient={patient} onBack={vi.fn()} />)
    const emptyBadge = screen.getByText("No known allergies")
    expect(emptyBadge.closest("[data-status]")?.getAttribute("data-status")).toBe("neutral")
  })

  it("should show skeletons while alerts load", () => {
    mockedAllergies.mockReturnValue(queryResult({ isLoading: true, data: undefined }))
    render(<PatientHeader patient={patient} onBack={vi.fn()} />)
    expect(screen.getByTestId("patient-header-allergies-loading")).toBeInTheDocument()
    expect(screen.getByTestId("patient-header-chip-loading")).toBeInTheDocument()
  })

  it("should offer a retry action when allergies fail to load", () => {
    const refetch = vi.fn()
    mockedAllergies.mockReturnValue(queryResult({ isError: true, data: undefined, refetch }))
    render(<PatientHeader patient={patient} onBack={vi.fn()} />)
    fireEvent.click(screen.getByRole("button", { name: "Try again" }))
    expect(refetch).toHaveBeenCalledTimes(1)
  })

  it("should call onBack from the back button", () => {
    const onBack = vi.fn()
    render(<PatientHeader patient={patient} onBack={onBack} />)
    fireEvent.click(screen.getByRole("button", { name: "Back" }))
    expect(onBack).toHaveBeenCalledTimes(1)
  })

  it("should open the mobile alerts popover from the count chip", () => {
    mockedAllergies.mockReturnValue(
      queryResult({ data: [{ allergen_display: "Penicillin", clinical_status: "active" }] })
    )
    render(<PatientHeader patient={patient} onBack={vi.fn()} />)
    fireEvent.click(screen.getByRole("button", { name: "View clinical alerts" }))
    expect(screen.getAllByText("Penicillin").length).toBeGreaterThan(1)
  })
})
