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
  usePatientQuery: vi.fn(),
}))

vi.mock("../hooks/useIdempotencyKey", () => ({
  useIdempotencyKey: () => ({
    getOrCreateKey: () => "fixed-idempotency-key",
    resetKey: vi.fn(),
  }),
}))

import { useStaffListQuery } from "../../staff/queries"
import { usePatientsQuery, usePatientQuery } from "../../patients/queries"

const mockedUseStaffListQuery = vi.mocked(useStaffListQuery)
const mockedUsePatientsQuery = vi.mocked(usePatientsQuery)
const mockedUsePatientQuery = vi.mocked(usePatientQuery)

const staffFixture = [
  { id: "emp-1", fullName: "Dr. André Silva", role: "DOCTOR" },
]
const patientsFixture = [
  { patient_id: "pat-1", fhir_resource_id: "fhir-pat-1", full_name: "Guilherme de Souza Araujo" },
]

const appointmentFixture = {
  id: "appointment-1",
  patient_fhir_id: "fhir-pat-1",
  staff_id: "emp-1",
  starts_at: new Date(2027, 4, 20, 9, 0).toISOString(),
  ends_at: new Date(2027, 4, 20, 9, 30).toISOString(),
  status: "scheduled",
  reason: "Consulta de rotina",
  version: 3,
  created_at: "2027-05-10T10:00:00.000Z",
} as const

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
  fireEvent.change(screen.getByRole("combobox", { name: "modals.create.startTime" }), {
    target: { value: startValue },
  })
  fireEvent.change(screen.getByRole("combobox", { name: "modals.create.endTime" }), {
    target: { value: endValue },
  })
}

const renderModal = (overrides: {
  onSubmit?: (payload: unknown) => Promise<void>
  isPending?: boolean
  defaultStaffId?: string
  defaultDate?: string
  defaultStartTime?: string
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
      defaultStartTime={overrides.defaultStartTime}
    />
  )
  return onSubmit
}

describe("AppointmentModal", () => {
  beforeEach(() => {
    mockedUseStaffListQuery.mockReturnValue({ data: staffFixture } as ReturnType<typeof useStaffListQuery>)
mockedUsePatientsQuery.mockReturnValue({
      data: { patients: patientsFixture, total: patientsFixture.length, page: 1, limit: 100 },
    } as ReturnType<typeof usePatientsQuery>)
    mockedUsePatientQuery.mockReturnValue({ data: undefined } as ReturnType<typeof usePatientQuery>)
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

  it("should submit appointment without a reason", async () => {
    const onSubmit = renderModal({ defaultDate: "2027-05-20" })
    await fillPatientAndStaff()
    fillTimes("09:00", "09:30")

    fireEvent.click(screen.getByRole("button", { name: "modals.create.confirm" }))
    await waitFor(() => expect(onSubmit).toHaveBeenCalled())

    const submittedPayload = (onSubmit as ReturnType<typeof vi.fn>).mock.calls[0][0]
    expect(submittedPayload).toMatchObject({
      patient_fhir_id: "fhir-pat-1",
      staff_id: "emp-1",
      reason: "",
      idempotency_key: "fixed-idempotency-key",
    })
    const parsedStart = new Date(submittedPayload.starts_at as string)
    const parsedEnd = new Date(submittedPayload.ends_at as string)
    expect(parsedStart.getTime()).toBe(new Date("2027-05-20T09:00").getTime())
    expect(parsedEnd.getTime()).toBe(new Date("2027-05-20T09:30").getTime())
  })

  it("should submit appointment with a reason", async () => {
    const onSubmit = renderModal({ defaultDate: "2027-05-20" })
    await fillPatientAndStaff()
    fillTimes("09:00", "09:45")
    fireEvent.change(screen.getByPlaceholderText("modals.create.reasonPlaceholder"), {
      target: { value: "Consulta de rotina" },
    })

    fireEvent.click(screen.getByRole("button", { name: "modals.create.confirm" }))
    await waitFor(() => expect(onSubmit).toHaveBeenCalled())

    const submittedPayload = (onSubmit as ReturnType<typeof vi.fn>).mock.calls[0][0]
    expect(submittedPayload).toMatchObject({
      patient_fhir_id: "fhir-pat-1",
      staff_id: "emp-1",
      reason: "Consulta de rotina",
      idempotency_key: "fixed-idempotency-key",
    })
  })

  it("should display conflict message when backend returns 409", async () => {
    const onSubmit = vi.fn().mockRejectedValue(
      Object.assign(new Error("Conflict"), { isAxiosError: true, response: { status: 409 } })
    )
    renderModal({ onSubmit, defaultDate: "2027-05-20" })
    await fillPatientAndStaff()
    fillTimes("09:00", "09:30")

    fireEvent.click(screen.getByRole("button", { name: "modals.create.confirm" }))
    expect(await screen.findByText("errors.conflict")).toBeDefined()
  })

  it("should prefill the start time from defaultStartTime", () => {
    renderModal({ defaultDate: "2027-05-20", defaultStartTime: "09:15" })
    const startTimeSelect = screen.getByRole("combobox", { name: "modals.create.startTime" }) as HTMLSelectElement
    expect(startTimeSelect.value).toBe("09:15")
  })

  it("should block past dates on the appointment date input", () => {
    render(
      <AppointmentModal isOpen onClose={vi.fn()} onSubmit={vi.fn()} isPending={false} defaultDate="2027-05-20" />
    )
    const dateInput = document.querySelector('input[type="date"]') as HTMLInputElement
    const now = new Date()
    const expectedMin = `${String(now.getFullYear()).padStart(4, "0")}-${String(now.getMonth() + 1).padStart(2, "0")}-${String(now.getDate()).padStart(2, "0")}`
    expect(dateInput.min).toBe(expectedMin)
  })

  it("should render nothing when closed", () => {
    const { container } = render(
      <AppointmentModal isOpen={false} onClose={vi.fn()} onSubmit={vi.fn()} isPending={false} />
    )
    expect(container).toBeEmptyDOMElement()
  })

  it("should prefill the form from the appointment in edit mode", () => {
    render(
      <AppointmentModal
        isOpen
        onClose={vi.fn()}
        onSubmit={vi.fn()}
        isPending={false}
        appointment={appointmentFixture}
      />
    )

    expect(screen.getByText("modals.edit.title")).toBeDefined()
    const comboboxes = screen.getAllByRole("combobox")
    expect(comboboxes[0]).toHaveTextContent("Guilherme de Souza Araujo")
    expect(comboboxes[1]).toHaveTextContent("Dr. André Silva")
    const dateInput = document.querySelector('input[type="date"]') as HTMLInputElement
    expect(dateInput.value).toBe("2027-05-20")
    const startTimeSelect = screen.getByRole("combobox", { name: "modals.create.startTime" }) as HTMLSelectElement
    expect(startTimeSelect.value).toBe("09:00")
    const endTimeSelect = screen.getByRole("combobox", { name: "modals.create.endTime" }) as HTMLSelectElement
    expect(endTimeSelect.value).toBe("09:30")
    expect(screen.getByPlaceholderText("modals.create.reasonPlaceholder")).toHaveValue("Consulta de rotina")
  })

  it("should submit an update payload without an idempotency key in edit mode", async () => {
    const onSubmit = vi.fn().mockResolvedValue(undefined)
    render(
      <AppointmentModal
        isOpen
        onClose={vi.fn()}
        onSubmit={onSubmit}
        isPending={false}
        appointment={appointmentFixture}
      />
    )

    fireEvent.click(screen.getByRole("button", { name: "modals.edit.confirm" }))
    await waitFor(() => expect(onSubmit).toHaveBeenCalled())

    const submittedPayload = (onSubmit as ReturnType<typeof vi.fn>).mock.calls[0][0]
    expect(submittedPayload).toEqual({
      patient_fhir_id: "fhir-pat-1",
      staff_id: "emp-1",
      starts_at: new Date(2027, 4, 20, 9, 0).toISOString(),
      ends_at: new Date(2027, 4, 20, 9, 30).toISOString(),
      reason: "Consulta de rotina",
    })
    expect(submittedPayload).not.toHaveProperty("idempotency_key")
  })

  it("should display conflict message when the update returns 409", async () => {
    const onSubmit = vi.fn().mockRejectedValue(
      Object.assign(new Error("Conflict"), { isAxiosError: true, response: { status: 409 } })
    )
    render(
      <AppointmentModal
        isOpen
        onClose={vi.fn()}
        onSubmit={onSubmit}
        isPending={false}
        appointment={appointmentFixture}
      />
    )

    fireEvent.click(screen.getByRole("button", { name: "modals.edit.confirm" }))
    expect(await screen.findByText("errors.conflict")).toBeDefined()
  })
})
