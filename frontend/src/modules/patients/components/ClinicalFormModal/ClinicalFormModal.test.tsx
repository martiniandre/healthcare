import { describe, it, expect, vi } from "vitest"
import { render, screen, fireEvent, waitFor } from "@testing-library/react"
import { ClinicalFormModal } from "./ClinicalFormModal"
import { allergyFormConfig, vitalSignsPanelFormConfig } from "./clinicalFormConfigs"

vi.mock("react-i18next", async (importOriginal) => {
  const actual = await importOriginal()
  return {
    ...actual,
    useTranslation: () => ({
      t: (key: string) => key,
    }),
  }
})

const heartRateLabel = "modals.observation.panel.heartRate"
const temperatureLabel = "modals.observation.panel.bodyTemperature"
const systolicLabel = "modals.observation.panel.systolicBloodPressure"
const diastolicLabel = "modals.observation.panel.diastolicBloodPressure"

const fillAllergyForm = () => {
  fireEvent.change(screen.getByLabelText("modals.allergy.display"), { target: { value: "Penicillin allergy" } })
  fireEvent.change(screen.getByLabelText("modals.allergy.code"), { target: { value: "300916003" } })
  fireEvent.change(screen.getByLabelText("modals.allergy.reaction"), { target: { value: "Severe hives" } })
}

describe("ClinicalFormModal with the vital signs panel config", () => {
  it("should submit filled metrics as numbers leaving unmeasured ones undefined", async () => {
    const onSubmit = vi.fn()
    render(
      <ClinicalFormModal
        isOpen
        onClose={vi.fn()}
        onSubmit={onSubmit}
        isPending={false}
        config={vitalSignsPanelFormConfig}
      />
    )
    fireEvent.change(screen.getByLabelText(heartRateLabel), { target: { value: "72" } })
    fireEvent.change(screen.getByLabelText(systolicLabel), { target: { value: "120" } })
    fireEvent.change(screen.getByLabelText(diastolicLabel), { target: { value: "80" } })
    fireEvent.click(screen.getByRole("button", { name: "modals.observation.confirm" }))
    await waitFor(() => expect(onSubmit).toHaveBeenCalled())
    expect(onSubmit.mock.calls[0][0]).toEqual({
      heartRate: 72,
      systolicBloodPressure: 120,
      diastolicBloodPressure: 80,
    })
  })

  it("should block submission until at least one metric is measured", async () => {
    const onSubmit = vi.fn()
    render(
      <ClinicalFormModal
        isOpen
        onClose={vi.fn()}
        onSubmit={onSubmit}
        isPending={false}
        config={vitalSignsPanelFormConfig}
      />
    )
    fireEvent.click(screen.getByRole("button", { name: "modals.observation.confirm" }))
    await waitFor(() => expect(onSubmit).not.toHaveBeenCalled())
    expect(
      (
        await screen.findAllByText(
          /Informe pelo menos um sinal vital|Record at least one vital sign|Registra al menos un signo vital/
        )
      ).length
    ).toBeGreaterThan(0)
  })

  it("should block submission when a metric leaves its clinical range", async () => {
    const onSubmit = vi.fn()
    render(
      <ClinicalFormModal
        isOpen
        onClose={vi.fn()}
        onSubmit={onSubmit}
        isPending={false}
        config={vitalSignsPanelFormConfig}
      />
    )
    fireEvent.change(screen.getByLabelText(temperatureLabel), { target: { value: "50" } })
    fireEvent.click(screen.getByRole("button", { name: "modals.observation.confirm" }))
    await waitFor(() => expect(onSubmit).not.toHaveBeenCalled())
    expect(
      (await screen.findAllByText(/fora do intervalo clínico|outside the allowed clinical range|fuera del rango clínico/))
        .length
    ).toBeGreaterThan(0)
  })

  it("should summarise the errors at the top of long forms after a failed submit", async () => {
    render(
      <ClinicalFormModal
        isOpen
        onClose={vi.fn()}
        onSubmit={vi.fn()}
        isPending={false}
        config={vitalSignsPanelFormConfig}
      />
    )
    expect(screen.queryByTestId("form-error-summary")).not.toBeInTheDocument()
    fireEvent.click(screen.getByRole("button", { name: "modals.observation.confirm" }))
    expect(await screen.findByTestId("form-error-summary")).toBeInTheDocument()
  })

  it("should not mark optional vital sign fields as required", () => {
    render(
      <ClinicalFormModal
        isOpen
        onClose={vi.fn()}
        onSubmit={vi.fn()}
        isPending={false}
        config={vitalSignsPanelFormConfig}
      />
    )
    expect(screen.queryByText("form.required")).not.toBeInTheDocument()
    expect(screen.getByLabelText(heartRateLabel)).toHaveAttribute("inputmode", "numeric")
    expect(screen.getByLabelText(temperatureLabel)).toHaveAttribute("inputmode", "decimal")
  })
})

describe("ClinicalFormModal with the allergy config", () => {
  const renderAllergyModal = (overrides: Partial<Parameters<typeof ClinicalFormModal>[0]> = {}) => {
    const properties = {
      isOpen: true,
      onClose: vi.fn(),
      onSubmit: vi.fn(),
      isPending: false,
      config: allergyFormConfig,
      ...overrides,
    }
    const view = render(<ClinicalFormModal {...(properties as Parameters<typeof ClinicalFormModal>[0])} />)
    return { ...view, properties }
  }

  it("should mark every field as required for assistive technology", () => {
    renderAllergyModal()
    expect(screen.getAllByText("form.required")).toHaveLength(3)
    expect(screen.getByLabelText("modals.allergy.display")).toHaveAttribute("aria-required", "true")
  })

  it("should connect hints and errors to the field through aria-describedby", async () => {
    renderAllergyModal()
    const codeInput = screen.getByLabelText("modals.allergy.code")
    expect(codeInput.getAttribute("aria-describedby")).toContain("description")
    fireEvent.blur(codeInput)
    await waitFor(() => expect(codeInput).toHaveAttribute("aria-invalid", "true"))
    expect(codeInput.getAttribute("aria-describedby")).toContain("error")
  })

  it("should not show an error before the field is touched", () => {
    renderAllergyModal()
    fireEvent.change(screen.getByLabelText("modals.allergy.display"), { target: { value: "Pe" } })
    expect(screen.getByLabelText("modals.allergy.display")).toHaveAttribute("aria-invalid", "false")
  })

  it("should focus the first invalid field when submitting an incomplete form", async () => {
    renderAllergyModal()
    fireEvent.click(screen.getByRole("button", { name: "modals.allergy.confirm" }))
    await waitFor(() => expect(screen.getByLabelText("modals.allergy.display")).toHaveFocus())
  })

  it("should close immediately when the form is pristine", () => {
    const { properties } = renderAllergyModal()
    fireEvent.click(screen.getByRole("button", { name: "modal.cancel" }))
    expect(properties.onClose).toHaveBeenCalledTimes(1)
  })

  it("should ask for confirmation before discarding a dirty form", async () => {
    const { properties } = renderAllergyModal()
    fireEvent.change(screen.getByLabelText("modals.allergy.display"), { target: { value: "Peanut" } })
    fireEvent.click(screen.getByRole("button", { name: "modal.cancel" }))
    expect(properties.onClose).not.toHaveBeenCalled()
    expect(await screen.findByText("discardChanges.title")).toBeInTheDocument()

    fireEvent.click(screen.getByRole("button", { name: "discardChanges.keepEditing" }))
    expect(properties.onClose).not.toHaveBeenCalled()
    expect(screen.getByLabelText("modals.allergy.display")).toHaveValue("Peanut")

    fireEvent.click(screen.getByRole("button", { name: "modal.cancel" }))
    fireEvent.click(await screen.findByRole("button", { name: "discardChanges.discard" }))
    expect(properties.onClose).toHaveBeenCalledTimes(1)
  })

  it("should ask for confirmation when a dirty form is closed with Escape", async () => {
    const { properties } = renderAllergyModal()
    fireEvent.change(screen.getByLabelText("modals.allergy.display"), { target: { value: "Peanut" } })
    fireEvent.keyDown(screen.getByRole("dialog"), { key: "Escape" })
    expect(await screen.findByText("discardChanges.title")).toBeInTheDocument()
    expect(properties.onClose).not.toHaveBeenCalled()
  })

  it("should submit only once when the form is submitted repeatedly while saving", async () => {
    const onSubmit = vi.fn().mockReturnValue(new Promise(() => undefined))
    renderAllergyModal({ onSubmit })
    fillAllergyForm()
    const submitButton = screen.getByRole("button", { name: "modals.allergy.confirm" })
    fireEvent.click(submitButton)
    await waitFor(() => expect(onSubmit).toHaveBeenCalledTimes(1))
    fireEvent.click(submitButton)
    fireEvent.submit(submitButton.closest("form") as HTMLFormElement)
    await new Promise((resolve) => setTimeout(resolve, 50))
    expect(onSubmit).toHaveBeenCalledTimes(1)
  })

  it("should disable the submit button with a spinner while pending", () => {
    renderAllergyModal({ isPending: true })
    expect(screen.getByRole("button", { name: /modals.allergy.confirm/ })).toBeDisabled()
    expect(screen.getByRole("status")).toBeInTheDocument()
  })

  it("should show the server error and keep what the user typed", async () => {
    const onSubmit = vi.fn().mockResolvedValue(undefined)
    const { rerender, properties } = renderAllergyModal({ onSubmit })
    fillAllergyForm()
    expect(screen.queryByTestId("form-submit-error")).not.toBeInTheDocument()
    fireEvent.click(screen.getByRole("button", { name: "modals.allergy.confirm" }))
    await waitFor(() => expect(onSubmit).toHaveBeenCalledTimes(1))

    rerender(<ClinicalFormModal {...(properties as Parameters<typeof ClinicalFormModal>[0])} hasSubmitError />)

    expect(await screen.findByTestId("form-submit-error")).toBeInTheDocument()
    expect(screen.getByLabelText("modals.allergy.display")).toHaveValue("Penicillin allergy")
    expect(screen.getByLabelText("modals.allergy.reaction")).toHaveValue("Severe hives")
  })
})
