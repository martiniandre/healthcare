import { describe, it, expect, vi, beforeEach } from "vitest"
import { render, screen, fireEvent } from "@testing-library/react"
import { AppointmentDetailsModal } from "./AppointmentDetailsModal"

vi.mock("react-i18next", () => ({
  useTranslation: () => ({
    t: (key: string) => key,
  }),
}))

vi.mock("../../staff/queries", () => ({
  useStaffListQuery: vi.fn(),
}))

vi.mock("../../patients/queries", () => ({
  usePatientQuery: vi.fn(),
}))

import { useStaffListQuery } from "../../staff/queries"
import { usePatientQuery } from "../../patients/queries"

const mockedUseStaffListQuery = vi.mocked(useStaffListQuery)
const mockedUsePatientQuery = vi.mocked(usePatientQuery)

const appointmentFixture = {
  id: "appointment-1",
  patient_fhir_id: "fhir-pat-1",
  staff_id: "emp-1",
  starts_at: "2027-05-20T09:00:00.000Z",
  ends_at: "2027-05-20T09:30:00.000Z",
  status: "scheduled",
  reason: "Consulta de rotina",
  version: 3,
  created_at: "2027-05-10T10:00:00.000Z",
}

describe("AppointmentDetailsModal", () => {
  beforeEach(() => {
    mockedUseStaffListQuery.mockReturnValue({
      data: [{ id: "emp-1", fullName: "Dr. André Silva", role: "DOCTOR" }],
    } as ReturnType<typeof useStaffListQuery>)
    mockedUsePatientQuery.mockReturnValue({
      data: { fhir_resource_id: "fhir-pat-1", full_name: "Guilherme de Souza Araujo" },
    } as ReturnType<typeof usePatientQuery>)
  })

  it("should render appointment information and action buttons", () => {
    render(
      <AppointmentDetailsModal
        isOpen
        appointment={appointmentFixture}
        onClose={vi.fn()}
        onEdit={vi.fn()}
        onCancel={vi.fn()}
      />
    )

    expect(screen.getByText("modals.details.title")).toBeDefined()
    expect(screen.getByText("Guilherme de Souza Araujo")).toBeDefined()
    expect(screen.getByText("Dr. André Silva — DOCTOR")).toBeDefined()
    expect(screen.getByText("Consulta de rotina")).toBeDefined()
    expect(screen.getByText("status.scheduled")).toBeDefined()
    expect(screen.getByRole("button", { name: "modals.details.edit" })).toBeDefined()
    expect(screen.getByRole("button", { name: "modals.details.cancel" })).toBeDefined()
    expect(screen.getByRole("button", { name: "modals.details.close" })).toBeDefined()
  })

  it("should call onEdit with the appointment when the edit button is clicked", () => {
    const onEdit = vi.fn()
    render(
      <AppointmentDetailsModal
        isOpen
        appointment={appointmentFixture}
        onClose={vi.fn()}
        onEdit={onEdit}
        onCancel={vi.fn()}
      />
    )

    fireEvent.click(screen.getByRole("button", { name: "modals.details.edit" }))
    expect(onEdit).toHaveBeenCalledWith(appointmentFixture)
  })

  it("should call onCancel with the appointment when the cancel button is clicked", () => {
    const onCancel = vi.fn()
    render(
      <AppointmentDetailsModal
        isOpen
        appointment={appointmentFixture}
        onClose={vi.fn()}
        onEdit={vi.fn()}
        onCancel={onCancel}
      />
    )

    fireEvent.click(screen.getByRole("button", { name: "modals.details.cancel" }))
    expect(onCancel).toHaveBeenCalledWith(appointmentFixture)
  })

  it("should call onClose when the close button is clicked", () => {
    const onClose = vi.fn()
    render(
      <AppointmentDetailsModal
        isOpen
        appointment={appointmentFixture}
        onClose={onClose}
        onEdit={vi.fn()}
        onCancel={vi.fn()}
      />
    )

    fireEvent.click(screen.getByRole("button", { name: "modals.details.close" }))
    expect(onClose).toHaveBeenCalled()
  })

  it("should render nothing when closed or without an appointment", () => {
    const { container: closedContainer } = render(
      <AppointmentDetailsModal
        isOpen={false}
        appointment={appointmentFixture}
        onClose={vi.fn()}
        onEdit={vi.fn()}
        onCancel={vi.fn()}
      />
    )
    expect(closedContainer).toBeEmptyDOMElement()

    const { container: emptyAppointmentContainer } = render(
      <AppointmentDetailsModal
        isOpen
        appointment={null}
        onClose={vi.fn()}
        onEdit={vi.fn()}
        onCancel={vi.fn()}
      />
    )
    expect(emptyAppointmentContainer).toBeEmptyDOMElement()
  })
})