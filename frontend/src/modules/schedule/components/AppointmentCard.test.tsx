import { describe, it, expect, vi, beforeEach } from "vitest"
import { render, screen, fireEvent } from "@testing-library/react"
import { AppointmentCard } from "./AppointmentCard"

vi.mock("react-i18next", () => ({
  useTranslation: () => ({
    t: (key: string) => key,
  }),
}))

vi.mock("../../patients/queries", () => ({
  usePatientQuery: vi.fn(),
}))

import { usePatientQuery } from "../../patients/queries"

const mockedUsePatientQuery = vi.mocked(usePatientQuery)

const baseAppointment = {
  id: "appt-1",
  patient_fhir_id: "fhir-pat-1",
  staff_id: "fhir-staff-1",
  starts_at: "2026-05-20T09:00:00Z",
  ends_at: "2026-05-20T10:00:00Z",
  status: "scheduled" as const,
  reason: "Consulta de rotina",
  version: 1,
  created_at: "2026-05-10T10:00:00Z",
}

describe("AppointmentCard", () => {
  beforeEach(() => {
    mockedUsePatientQuery.mockReturnValue({
      data: { patient_id: "pat-1", fhir_resource_id: "fhir-pat-1", full_name: "Guilherme de Souza Araujo" },
      isLoading: false,
    } as ReturnType<typeof usePatientQuery>)
  })

  it("should render patient name, time range and reason", () => {
    render(<AppointmentCard appointment={baseAppointment} onCancel={vi.fn()} />)
    expect(screen.getByText("Guilherme de Souza Araujo")).toBeDefined()
    expect(screen.getByText("Consulta de rotina")).toBeDefined()
    expect(screen.getByText("status.scheduled")).toBeDefined()
  })

  it("should show cancel button for scheduled appointment and trigger onCancel", () => {
    const onCancel = vi.fn()
    render(<AppointmentCard appointment={baseAppointment} onCancel={onCancel} />)
    fireEvent.click(screen.getByRole("button", { name: "cards.cancelAppointment" }))
    expect(onCancel).toHaveBeenCalledWith(baseAppointment)
  })

  it("should not render cancel button for cancelled appointment", () => {
    render(
      <AppointmentCard
        appointment={{ ...baseAppointment, status: "cancelled" }}
        onCancel={vi.fn()}
      />
    )
    expect(screen.queryByRole("button", { name: "cards.cancelAppointment" })).toBeNull()
  })

  it("should fall back to patient fhir id when patient name is not available", () => {
    mockedUsePatientQuery.mockReturnValue({
      data: null,
      isLoading: false,
    } as ReturnType<typeof usePatientQuery>)
    render(<AppointmentCard appointment={baseAppointment} onCancel={vi.fn()} />)
    expect(screen.getByText("fhir-pat-1")).toBeDefined()
  })
})
