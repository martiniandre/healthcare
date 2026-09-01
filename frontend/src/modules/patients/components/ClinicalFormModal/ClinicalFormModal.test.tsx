import { describe, it, expect, vi } from "vitest"
import { render, screen, fireEvent, waitFor } from "@testing-library/react"
import { ClinicalFormModal } from "./ClinicalFormModal"
import { vitalSignsPanelFormConfig } from "./clinicalFormConfigs"

vi.mock("react-i18next", async (importOriginal) => {
  const actual = await importOriginal()
  return {
    ...actual,
    useTranslation: () => ({
      t: (key: string) => key,
    }),
  }
})

const getNumericInputByIndex = (index: number) => screen.getAllByRole("spinbutton")[index]

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
    fireEvent.change(getNumericInputByIndex(0), { target: { value: "72" } })
    fireEvent.change(getNumericInputByIndex(2), { target: { value: "120" } })
    fireEvent.change(getNumericInputByIndex(3), { target: { value: "80" } })
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
      await screen.findByText(
        /Informe pelo menos um sinal vital|Record at least one vital sign|Registra al menos un signo vital/
      )
    ).toBeInTheDocument()
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
    fireEvent.change(getNumericInputByIndex(1), { target: { value: "50" } })
    fireEvent.click(screen.getByRole("button", { name: "modals.observation.confirm" }))
    await waitFor(() => expect(onSubmit).not.toHaveBeenCalled())
    expect(
      await screen.findByText(/fora do intervalo clínico|outside the allowed clinical range|fuera del rango clínico/)
    ).toBeDefined()
  })
})
