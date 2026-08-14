import { describe, it, expect, vi } from "vitest"
import { render, screen, fireEvent, waitFor } from "@testing-library/react"
import { ClinicalFormModal } from "./ClinicalFormModal"
import { observationFormConfig } from "./clinicalFormConfigs"

vi.mock("react-i18next", async (importOriginal) => {
  const actual = await importOriginal()
  return {
    ...actual,
    useTranslation: () => ({
      t: (key: string) => key,
    }),
  }
})

describe("ClinicalFormModal", () => {
  it("should enrich observation submission with default metric metadata", async () => {
    const onSubmit = vi.fn()
    render(
      <ClinicalFormModal
        isOpen
        onClose={vi.fn()}
        onSubmit={onSubmit}
        isPending={false}
        config={observationFormConfig}
      />
    )
    fireEvent.change(screen.getByPlaceholderText("modals.observation.valuePlaceholder"), {
      target: { value: "72" },
    })
    fireEvent.click(screen.getByRole("button", { name: "modals.observation.confirm" }))
    await waitFor(() => expect(onSubmit).toHaveBeenCalled())
    expect(onSubmit.mock.calls[0][0]).toEqual({
      loincCode: "8867-4",
      valueQuantity: 72,
      codeDisplay: "Frequência Cardíaca",
      valueUnit: "bpm",
    })
  })

  it("should enrich observation submission with selected metric metadata", async () => {
    const onSubmit = vi.fn()
    render(
      <ClinicalFormModal
        isOpen
        onClose={vi.fn()}
        onSubmit={onSubmit}
        isPending={false}
        config={observationFormConfig}
      />
    )
    fireEvent.click(screen.getByRole("combobox"))
    const temperatureOption = (await screen.findAllByRole("option", { name: "modals.observation.temperature" }))
      .find((option) => option.tagName === "DIV")
    if (!temperatureOption) {
      throw new Error("Temperature option not found")
    }
    fireEvent.click(temperatureOption)
    fireEvent.change(screen.getByPlaceholderText("modals.observation.valuePlaceholder"), {
      target: { value: "37.5" },
    })
    fireEvent.click(screen.getByRole("button", { name: "modals.observation.confirm" }))
    await waitFor(() => expect(onSubmit).toHaveBeenCalled())
    expect(onSubmit.mock.calls[0][0]).toEqual({
      loincCode: "8310-5",
      valueQuantity: 37.5,
      codeDisplay: "Temperatura Corporal",
      valueUnit: "°C",
    })
  })

  it("should block submission when required numeric value is missing", async () => {
    const onSubmit = vi.fn()
    render(
      <ClinicalFormModal
        isOpen
        onClose={vi.fn()}
        onSubmit={onSubmit}
        isPending={false}
        config={observationFormConfig}
      />
    )
    fireEvent.click(screen.getByRole("button", { name: "modals.observation.confirm" }))
    await waitFor(() => expect(onSubmit).not.toHaveBeenCalled())
    expect(await screen.findByText(/expected number/i)).toBeDefined()
  })
})
