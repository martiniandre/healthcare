import { describe, it, expect, vi, beforeEach } from "vitest"
import { render, screen, fireEvent, waitFor } from "@testing-library/react"
import { AppointmentModal } from "./AppointmentModal"

vi.mock("react-i18next", () => ({
  useTranslation: () => ({
    t: (key: string) => key,
  }),
}))

vi.mock("../../staff/queries", () => ({
  useStaffListQuery: vi.fn(),
}))

vi.mock("../../patients/queries", () => ({
  usePatientsQuery: vi.fn(),
}))

vi.mock("../hooks/useIdempotencyKey", () => ({
  useIdempotencyKey: () => ({
    getOrCreateKey: () => "fixed-idempotency-key",
    resetKey: vi.fn(),
  }),
}))

import { useStaffListQuery } from "../../staff/queries"
import { usePatientsQuery } from "../../patients/queries"

const mockedUseStaffListQuery = vi.mocked(useStaffListQuery)
const mockedUsePatientsQuery = vi.mocked(usePatientsQuery)

const staffFixture = [
  { id: "emp-1", fullName: "Dr. André Silva", role: "DOCTOR", fhirResourceId: "fhir-staff-1" },
]
const patientsFixture = [
  { patient_id: "pat-1", fhir_resource_id: "fhir-pat-1", full_name: "Guilherme de Souza Araujo" },
]

const fillPatientAndStaff = async () => {
  const comboboxes = screen.getAllByRole("combobox")
  fireEvent.click(comboboxes[0])
  const patientOption = (await screen.findAllByRole("option", { name: "Guilherme de Souza Araujo" }))
    .find((option) => option.tagName === "DIV")
  if (!patientOption) {
    throw new Error("Patient option not found")
  }
  fireEvent.click(patientOption)

  const refreshedComboboxes = screen.getAllByRole("combobox")
  fireEvent.click(refreshedComboboxes[1])
  const staffOption = (await screen.findAllByRole("option", { name: /Dr\. André Silva/ }))
    .find((option) => option.tagName === "DIV")
  if (!staffOption) {
    throw new Error("Staff option not found")
  }
  fireEvent.click(staffOption)
}

const fillTimes = (startValue: string, endValue: string) => {
  const timeInputs = document.querySelectorAll('input[type="time"]')
  fireEvent.change(timeInputs[0], { target: { value: startValue } })
  fireEvent.change(timeInputs[1], { target: { value: endValue } })
}

const renderModal = (overrides: {
  onSubmit?: (payload: unknown) => Promise<void>
  isPending?: boolean
  defaultStaffId?: string
  defaultDate?: string
} = {}) => {
  const onSubmit = overrides.onSubmit ?? vi.fn().mockResolvedValue(undefined)
  render(
    <AppointmentModal
      isOpen
      onClose={vi.fn()}
      onSubmit={onSubmit}
      isPending={overrides.isPending ?? false}
      defaultStaffId={overrides.defaultStaffId}
      defaultDate={overrides.defaultDate}
    />
  )
  return onSubmit
}

describe("AppointmentModal", () => {
  beforeEach(() => {
    mockedUseStaffListQuery.mockReturnValue({ data: staffFixture } as ReturnType<typeof useStaffListQuery>)
    mockedUsePatientsQuery.mockReturnValue({ data: patientsFixture } as ReturnType<typeof usePatientsQuery>)
  })

  it("should show validation errors when required fields are missing", async () => {
    renderModal()
    fireEvent.click(screen.getByRole("button", { name: "modals.create.confirm" }))
    expect(await screen.findByText("validation.patientRequired")).toBeDefined()
    expect(await screen.findByText("validation.staffRequired")).toBeDefined()
    expect(await screen.findByText("validation.dateRequired")).toBeDefined()
    expect(await screen.findByText("validation.startTimeRequired")).toBeDefined()
    expect(await screen.findByText("validation.endTimeRequired")).toBeDefined()
  })

  it("should submit appointment with patient, staff, times and idempotency key", async () => {
    const onSubmit = renderModal({ defaultDate: "2026-05-20" })
    await fillPatientAndStaff()
    fillTimes("09:00", "10:00")
    fireEvent.change(screen.getByPlaceholderText("modals.create.reasonPlaceholder"), {
      target: { value: "Consulta de rotina" },
    })

    fireEvent.click(screen.getByRole("button", { name: "modals.create.confirm" }))
    await waitFor(() => expect(onSubmit).toHaveBeenCalled())

    const submittedPayload = (onSubmit as ReturnType<typeof vi.fn>).mock.calls[0][0]
    expect(submittedPayload).toMatchObject({
      patient_fhir_id: "fhir-pat-1",
      staff_id: "fhir-staff-1",
      reason: "Consulta de rotina",
      idempotency_key: "fixed-idempotency-key",
    })
    const parsedStart = new Date(submittedPayload.starts_at as string)
    const parsedEnd = new Date(submittedPayload.ends_at as string)
    expect(parsedStart.getTime()).toBe(new Date("2026-05-20T09:00").getTime())
    expect(parsedEnd.getTime()).toBe(new Date("2026-05-20T10:00").getTime())
  })

  it("should display conflict message when backend returns 409", async () => {
    const onSubmit = vi.fn().mockRejectedValue(
      Object.assign(new Error("Conflict"), { isAxiosError: true, response: { status: 409 } })
    )
    renderModal({ onSubmit, defaultDate: "2026-05-20" })
    await fillPatientAndStaff()
    fillTimes("09:00", "10:00")
    fireEvent.change(screen.getByPlaceholderText("modals.create.reasonPlaceholder"), {
      target: { value: "Consulta de rotina" },
    })

    fireEvent.click(screen.getByRole("button", { name: "modals.create.confirm" }))
    expect(await screen.findByText("errors.conflict")).toBeDefined()
  })

  it("should render nothing when closed", () => {
    const { container } = render(
      <AppointmentModal isOpen={false} onClose={vi.fn()} onSubmit={vi.fn()} isPending={false} />
    )
    expect(container).toBeEmptyDOMElement()
  })
})
