import { describe, it, expect, vi, beforeEach } from "vitest"
import { render, screen, fireEvent, waitFor } from "@testing-library/react"
import { PatientModal } from "./PatientModal"

vi.mock("react-i18next", async (importOriginal) => {
  const actual = await importOriginal<typeof import("react-i18next")>()
  return {
    ...actual,
    useTranslation: () => ({ t: (key: string) => key, i18n: { language: "pt-BR" } }),
  }
})

const mutationState = vi.hoisted(() => ({
  mutateAsync: vi.fn(),
  reset: vi.fn(),
  isPending: false,
  isError: false,
}))

vi.mock("../queries", () => ({
  useCreatePatientMutation: () => mutationState,
}))

vi.mock("../../../shared/store/toast_store", () => ({
  toast: { success: vi.fn(), error: vi.fn(), info: vi.fn() },
}))

const fillValidPatient = () => {
  fireEvent.change(screen.getByLabelText("modal.name"), { target: { value: "Maria da Silva" } })
  fireEvent.change(screen.getByLabelText("modal.birthDate"), { target: { value: "01011990" } })
  fireEvent.change(screen.getByLabelText("modal.document"), { target: { value: "52998224725" } })
  fireEvent.change(screen.getByLabelText("modal.phone"), { target: { value: "11987654321" } })
}

describe("PatientModal", () => {
  beforeEach(() => {
    vi.clearAllMocks()
    mutationState.isPending = false
    mutationState.isError = false
    mutationState.mutateAsync.mockResolvedValue({})
  })

  it("should mark all four fields as required and use mobile friendly input hints", () => {
    render(<PatientModal isOpen onOpenChange={vi.fn()} />)
    expect(screen.getAllByText("form.required")).toHaveLength(4)
    expect(screen.getByLabelText("modal.birthDate")).toHaveAttribute("inputmode", "numeric")
    expect(screen.getByLabelText("modal.document")).toHaveAttribute("inputmode", "numeric")
    expect(screen.getByLabelText("modal.phone")).toHaveAttribute("type", "tel")
  })

  it("should not flash an error while typing the first time and validate after leaving the field", async () => {
    render(<PatientModal isOpen onOpenChange={vi.fn()} />)
    const nameInput = screen.getByLabelText("modal.name")
    fireEvent.change(nameInput, { target: { value: "Ma" } })
    expect(nameInput).toHaveAttribute("aria-invalid", "false")
    fireEvent.blur(nameInput)
    await waitFor(() => expect(nameInput).toHaveAttribute("aria-invalid", "true"))
    expect(nameInput.getAttribute("aria-describedby")).toContain("error")
  })

  it("should focus the first invalid field when submitting an empty form", async () => {
    render(<PatientModal isOpen onOpenChange={vi.fn()} />)
    fireEvent.click(screen.getByRole("button", { name: "modal.confirm" }))
    await waitFor(() => expect(screen.getByLabelText("modal.name")).toHaveFocus())
    expect(mutationState.mutateAsync).not.toHaveBeenCalled()
  })

  it("should submit the normalized payload for a valid patient", async () => {
    const onOpenChange = vi.fn()
    render(<PatientModal isOpen onOpenChange={onOpenChange} />)
    fillValidPatient()
    fireEvent.click(screen.getByRole("button", { name: "modal.confirm" }))
    await waitFor(() => expect(mutationState.mutateAsync).toHaveBeenCalledTimes(1))
    expect(mutationState.mutateAsync).toHaveBeenCalledWith({
      full_name: "Maria da Silva",
      birth_date: "1990-01-01",
      document_id: "529.982.247-25",
      phone_number: "(11) 98765-4321",
    })
    await waitFor(() => expect(onOpenChange).toHaveBeenCalledWith(false))
  })

  it("should ask before discarding typed data and close right away when pristine", async () => {
    const onOpenChange = vi.fn()
    render(<PatientModal isOpen onOpenChange={onOpenChange} />)
    fireEvent.click(screen.getByRole("button", { name: "modal.cancel" }))
    expect(onOpenChange).toHaveBeenCalledWith(false)

    onOpenChange.mockClear()
    fireEvent.change(screen.getByLabelText("modal.name"), { target: { value: "Maria" } })
    fireEvent.click(screen.getByRole("button", { name: "modal.cancel" }))
    expect(onOpenChange).not.toHaveBeenCalled()
    fireEvent.click(await screen.findByRole("button", { name: "discardChanges.discard" }))
    expect(onOpenChange).toHaveBeenCalledWith(false)
  })

  it("should show the server error at the top and keep the typed values", async () => {
    mutationState.mutateAsync.mockRejectedValue(new Error("network"))
    const { rerender } = render(<PatientModal isOpen onOpenChange={vi.fn()} />)
    fillValidPatient()
    fireEvent.click(screen.getByRole("button", { name: "modal.confirm" }))
    await waitFor(() => expect(mutationState.mutateAsync).toHaveBeenCalledTimes(1))

    mutationState.isError = true
    rerender(<PatientModal isOpen onOpenChange={vi.fn()} />)

    expect(await screen.findByTestId("form-submit-error")).toBeInTheDocument()
    expect(screen.getByLabelText("modal.name")).toHaveValue("Maria da Silva")
    expect(screen.getByLabelText("modal.document")).toHaveValue("529.982.247-25")
  })

  it("should disable the confirm button with a spinner while saving", () => {
    mutationState.isPending = true
    render(<PatientModal isOpen onOpenChange={vi.fn()} />)
    expect(screen.getByRole("button", { name: /modal.confirm/ })).toBeDisabled()
    expect(screen.getByRole("status")).toBeInTheDocument()
  })
})
