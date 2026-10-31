import { describe, it, expect, vi } from "vitest"
import { render, screen, fireEvent, waitFor } from "@testing-library/react"
import { CreateUnavailabilityModal } from "./CreateUnavailabilityModal"

const translationApi = {
  t: (key: string, options?: Record<string, unknown>) => (options ? `${key}:${JSON.stringify(options)}` : key),
  i18n: { language: "en-US" },
}

vi.mock("react-i18next", () => ({
  useTranslation: () => translationApi,
}))

describe("CreateUnavailabilityModal", () => {
  it("should render nothing when closed", () => {
    const { container } = render(
      <CreateUnavailabilityModal
        isOpen={false}
        onClose={vi.fn()}
        onSubmit={vi.fn()}
        isPending={false}
        staffId="staff-1"
        defaultDate="2026-09-01"
      />
    )
    expect(container.innerHTML).toBe("")
  })

  it("should render the modal title when open", () => {
    render(
      <CreateUnavailabilityModal
        isOpen={true}
        onClose={vi.fn()}
        onSubmit={vi.fn()}
        isPending={false}
        staffId="staff-1"
        defaultDate="2026-09-01"
      />
    )
    expect(screen.getByText("unavailability.modal.title")).toBeInTheDocument()
  })

  it("should render date, time, and reason inputs", () => {
    render(
      <CreateUnavailabilityModal
        isOpen={true}
        onClose={vi.fn()}
        onSubmit={vi.fn()}
        isPending={false}
        staffId="staff-1"
        defaultDate="2026-09-01"
      />
    )
    expect(screen.getByText("unavailability.modal.date")).toBeInTheDocument()
    expect(screen.getByText("unavailability.modal.startTime")).toBeInTheDocument()
    expect(screen.getByText("unavailability.modal.endTime")).toBeInTheDocument()
    expect(screen.getByText("unavailability.modal.reason")).toBeInTheDocument()
  })

  it("should render cancel and confirm buttons", () => {
    render(
      <CreateUnavailabilityModal
        isOpen={true}
        onClose={vi.fn()}
        onSubmit={vi.fn()}
        isPending={false}
        staffId="staff-1"
        defaultDate="2026-09-01"
      />
    )
    expect(screen.getByText("unavailability.modal.cancel")).toBeInTheDocument()
    expect(screen.getByText("unavailability.modal.confirm")).toBeInTheDocument()
  })

  it("should block past dates on the unavailability date input", () => {
    render(
      <CreateUnavailabilityModal
        isOpen={true}
        onClose={vi.fn()}
        onSubmit={vi.fn()}
        isPending={false}
        staffId="staff-1"
        defaultDate="2026-09-01"
      />
    )
    const dateInput = document.querySelector('input[type="date"]') as HTMLInputElement
    const now = new Date()
    const expectedMin = `${String(now.getFullYear()).padStart(4, "0")}-${String(now.getMonth() + 1).padStart(2, "0")}-${String(now.getDate()).padStart(2, "0")}`
    expect(dateInput.min).toBe(expectedMin)
  })

  const renderOpen = (overrides: { onClose?: () => void; onSubmit?: (payload: unknown) => Promise<void> } = {}) => {
    render(
      <CreateUnavailabilityModal
        isOpen={true}
        onClose={overrides.onClose ?? vi.fn()}
        onSubmit={overrides.onSubmit ?? vi.fn().mockResolvedValue(undefined)}
        isPending={false}
        staffId="staff-1"
        staffName="Dr. André Silva"
        defaultDate="2027-05-20"
      />
    )
  }

  const chooseSlot = async (comboboxName: RegExp, timeValue: string) => {
    fireEvent.click(screen.getByRole("combobox", { name: comboboxName }))
    const timeOption = (await screen.findAllByRole("option", { name: timeValue })).find(
      (option) => option.tagName === "DIV"
    )
    if (!timeOption) {
      throw new Error("Time option not found")
    }
    fireEvent.click(timeOption)
  }

  it("should fill the end time one hour after the start and summarize the block", async () => {
    renderOpen()
    await chooseSlot(/unavailability\.modal\.startTime/, "12:00")
    expect(screen.getByRole("combobox", { name: /unavailability\.modal\.endTime/ })).toHaveTextContent("13:00")
    const summary = screen.getByTestId("selection-summary")
    expect(summary).toHaveTextContent("unavailability.summary.forStaff")
    expect(summary).toHaveTextContent("Dr. André Silva")
  })

  it("should focus the first invalid field and list errors on failed submit", async () => {
    renderOpen()
    fireEvent.click(screen.getByRole("button", { name: "unavailability.modal.confirm" }))
    expect(await screen.findByTestId("form-error-summary")).toHaveTextContent("validation.startTimeRequired")
    await waitFor(() =>
      expect(screen.getByRole("combobox", { name: /unavailability\.modal\.startTime/ })).toHaveFocus()
    )
  })

  it("should ask for confirmation before closing a dirty form", async () => {
    const onClose = vi.fn()
    renderOpen({ onClose })
    fireEvent.change(screen.getByPlaceholderText("unavailability.modal.reasonPlaceholder"), {
      target: { value: "Treinamento" },
    })
    fireEvent.click(screen.getByRole("button", { name: "unavailability.modal.cancel" }))
    expect(await screen.findByText("discardChanges.title")).toBeDefined()
    expect(onClose).not.toHaveBeenCalled()
    fireEvent.click(screen.getByRole("button", { name: "discardChanges.discard" }))
    expect(onClose).toHaveBeenCalledTimes(1)
  })

  it("should keep the form open and show the conflict alert on 409", async () => {
    const onClose = vi.fn()
    const onSubmit = vi.fn().mockRejectedValue(
      Object.assign(new Error("Conflict"), { isAxiosError: true, response: { status: 409 } })
    )
    renderOpen({ onClose, onSubmit })
    await chooseSlot(/unavailability\.modal\.startTime/, "12:00")
    fireEvent.click(screen.getByRole("button", { name: "unavailability.modal.confirm" }))
    expect(await screen.findByTestId("conflict-alert")).toBeDefined()
    expect(onClose).not.toHaveBeenCalled()
  })

  it("should not submit twice while a submission is in flight", async () => {
    let resolveSubmit: () => void = () => undefined
    const onSubmit = vi.fn().mockImplementation(() => new Promise<void>((resolve) => { resolveSubmit = resolve }))
    renderOpen({ onSubmit })
    await chooseSlot(/unavailability\.modal\.startTime/, "12:00")
    const form = document.querySelector("form") as HTMLFormElement
    fireEvent.submit(form)
    fireEvent.submit(form)
    await waitFor(() => expect(onSubmit).toHaveBeenCalledTimes(1))
    resolveSubmit()
  })
})
