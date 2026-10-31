import { describe, it, expect, vi, beforeAll, beforeEach } from "vitest"
import { render, screen, fireEvent, waitFor } from "@testing-library/react"
import { MedicationModal } from "./MedicationModal"
import type { MedicationSafetyAlert } from "../../types"

vi.mock("react-i18next", async (importOriginal) => {
  const actual = await importOriginal<typeof import("react-i18next")>()
  return {
    ...actual,
    useTranslation: () => ({ t: (key: string) => key }),
  }
})

const abilityState = vi.hoisted(() => ({ canCreate: true }))

vi.mock("../../../../shared/auth/useAbility", () => ({
  useAbility: () => ({ can: () => abilityState.canCreate }),
}))

vi.mock("../../queries", () => ({
  useMedicationSafetyCheck: vi.fn(),
}))

import { useMedicationSafetyCheck } from "../../queries"

const mockedSafetyCheck = vi.mocked(useMedicationSafetyCheck)

const allergyAlert: MedicationSafetyAlert = {
  kind: "allergy",
  severity: "contraindicated",
  blocking: true,
  messageKey: "allergy.activeIngredient",
  params: { allergen: "Penicillin", medication: "Amoxicillin" },
  relatedResourceId: "allergy-1",
}

const interactionAlert: MedicationSafetyAlert = {
  kind: "interaction",
  severity: "moderate",
  blocking: false,
  messageKey: "interaction.moderate",
  params: { medication: "Amoxicillin", drugClass: "NSAID" },
  relatedResourceId: "med-9",
}

const configureSafetyCheck = (
  status: "idle" | "loading" | "error" | "ready",
  alerts: MedicationSafetyAlert[],
  verification: Awaited<ReturnType<ReturnType<typeof useMedicationSafetyCheck>["verifyNow"]>> = { verified: true, alerts },
) => {
  const verifyNow = vi.fn().mockResolvedValue(verification)
  mockedSafetyCheck.mockReturnValue({ status, alerts, verifyNow, retry: vi.fn() })
  return verifyNow
}

const renderModal = (onSubmit = vi.fn().mockResolvedValue({ status: "success" })) => {
  render(<MedicationModal patientId="patient-1" isOpen onClose={vi.fn()} onSubmit={onSubmit} isPending={false} />)
  return onSubmit
}

const fillPrescription = () => {
  fireEvent.change(screen.getByLabelText("modals.medication.name"), { target: { value: "Amoxicillin 500mg" } })
  fireEvent.change(screen.getByLabelText("modals.medication.dosage"), { target: { value: "1 tablet every 8 hours" } })
}

const submitButton = (name: string) => screen.getByRole("button", { name })

class ResizeObserverStub {
  observe() {}
  unobserve() {}
  disconnect() {}
}

beforeAll(() => {
  vi.stubGlobal("ResizeObserver", ResizeObserverStub)
})

describe("MedicationModal", () => {
  beforeEach(() => {
    abilityState.canCreate = true
    vi.clearAllMocks()
  })

  it("should submit the payload without override when there are no blocking alerts", async () => {
    configureSafetyCheck("ready", [interactionAlert])
    const onSubmit = renderModal()
    fillPrescription()
    expect(screen.getByTestId("safety-alert-interaction")).toBeInTheDocument()
    expect(screen.queryByTestId("override-section")).not.toBeInTheDocument()
    fireEvent.click(submitButton("modals.medication.confirm"))
    await waitFor(() => expect(onSubmit).toHaveBeenCalled())
    expect(onSubmit).toHaveBeenCalledWith({
      medicationDisplay: "Amoxicillin 500mg",
      dosageInstruction: "1 tablet every 8 hours",
      overrideJustification: undefined,
    })
  })

  it("should require a justification and risk acknowledgement when an alert is blocking", async () => {
    configureSafetyCheck("ready", [allergyAlert])
    const onSubmit = renderModal()
    fillPrescription()
    expect(screen.getByTestId("override-section")).toBeInTheDocument()
    fireEvent.click(submitButton("medicationSafety.override.confirm"))
    expect(await screen.findByText("validation.overrideJustificationMin")).toBeInTheDocument()
    expect(screen.getByText("validation.riskAcknowledgementRequired")).toBeInTheDocument()
    expect(onSubmit).not.toHaveBeenCalled()
  })

  it("should reject a justification shorter than 10 characters", async () => {
    configureSafetyCheck("ready", [allergyAlert])
    const onSubmit = renderModal()
    fillPrescription()
    fireEvent.change(screen.getByLabelText("medicationSafety.override.justification"), { target: { value: "too short" } })
    fireEvent.click(screen.getByRole("checkbox"))
    fireEvent.click(submitButton("medicationSafety.override.confirm"))
    expect(await screen.findByText("validation.overrideJustificationMin")).toBeInTheDocument()
    expect(onSubmit).not.toHaveBeenCalled()
  })

  it("should send the trimmed overrideJustification once justified and acknowledged", async () => {
    configureSafetyCheck("ready", [allergyAlert])
    const onSubmit = renderModal()
    fillPrescription()
    fireEvent.change(screen.getByLabelText("medicationSafety.override.justification"), {
      target: { value: "  Benefit outweighs the risk under supervision  " },
    })
    fireEvent.click(screen.getByRole("checkbox"))
    fireEvent.click(submitButton("medicationSafety.override.confirm"))
    await waitFor(() => expect(onSubmit).toHaveBeenCalled())
    expect(onSubmit).toHaveBeenCalledWith({
      medicationDisplay: "Amoxicillin 500mg",
      dosageInstruction: "1 tablet every 8 hours",
      overrideJustification: "Benefit outweighs the risk under supervision",
    })
  })

  it("should keep roles without prescribing permission from overriding blocks", () => {
    abilityState.canCreate = false
    configureSafetyCheck("ready", [allergyAlert])
    const onSubmit = renderModal()
    fillPrescription()
    expect(screen.getByTestId("override-not-permitted")).toBeInTheDocument()
    expect(screen.queryByLabelText("medicationSafety.override.justification")).not.toBeInTheDocument()
    expect(submitButton("modals.medication.confirm")).toBeDisabled()
    expect(onSubmit).not.toHaveBeenCalled()
  })

  it("should require an explicit acknowledgement when the check fails", async () => {
    configureSafetyCheck("error", [], { verified: false })
    const onSubmit = renderModal()
    fillPrescription()
    expect(screen.getByTestId("safety-check-failed")).toBeInTheDocument()
    fireEvent.click(submitButton("modals.medication.confirm"))
    expect(await screen.findByText("validation.unverifiedAcknowledgementRequired")).toBeInTheDocument()
    expect(onSubmit).not.toHaveBeenCalled()
    fireEvent.click(screen.getByRole("checkbox"))
    fireEvent.click(submitButton("modals.medication.confirm"))
    await waitFor(() => expect(onSubmit).toHaveBeenCalledTimes(1))
  })

  it("should verify again right before submitting and stop if a blocking alert appeared", async () => {
    const verifyNow = configureSafetyCheck("ready", [], { verified: true, alerts: [allergyAlert] })
    const onSubmit = renderModal()
    fillPrescription()
    fireEvent.click(submitButton("modals.medication.confirm"))
    await waitFor(() => expect(verifyNow).toHaveBeenCalled())
    expect(await screen.findByTestId("override-section")).toBeInTheDocument()
    expect(onSubmit).not.toHaveBeenCalled()
  })

  it("should show the alerts of a 409 response and ask for a justification", async () => {
    configureSafetyCheck("ready", [])
    const onSubmit = renderModal(vi.fn().mockResolvedValue({ status: "blocked", alerts: [allergyAlert] }))
    fillPrescription()
    fireEvent.click(submitButton("modals.medication.confirm"))
    await waitFor(() => expect(onSubmit).toHaveBeenCalledTimes(1))
    expect(await screen.findByTestId("safety-alert-allergy")).toBeInTheDocument()
    expect(screen.getByTestId("override-section")).toBeInTheDocument()
    expect(screen.getByText("medicationSafety.override.serverBlocked")).toBeInTheDocument()
  })

  it("should mark the prescription fields as required", () => {
    configureSafetyCheck("ready", [])
    renderModal()
    expect(screen.getAllByText("form.required")).toHaveLength(2)
    expect(screen.getByLabelText("modals.medication.name")).toHaveAttribute("aria-required", "true")
  })

  it("should ask for confirmation before discarding a dirty prescription", async () => {
    configureSafetyCheck("ready", [])
    const onClose = vi.fn()
    render(
      <MedicationModal
        patientId="patient-1"
        isOpen
        onClose={onClose}
        onSubmit={vi.fn().mockResolvedValue({ status: "success" })}
        isPending={false}
      />,
    )
    fillPrescription()
    fireEvent.click(screen.getByRole("button", { name: "modal.cancel" }))
    expect(onClose).not.toHaveBeenCalled()
    fireEvent.click(await screen.findByRole("button", { name: "discardChanges.discard" }))
    expect(onClose).toHaveBeenCalledTimes(1)
  })

  it("should close without confirmation when nothing was typed", () => {
    configureSafetyCheck("ready", [])
    const onClose = vi.fn()
    render(
      <MedicationModal
        patientId="patient-1"
        isOpen
        onClose={onClose}
        onSubmit={vi.fn()}
        isPending={false}
      />,
    )
    fireEvent.click(screen.getByRole("button", { name: "modal.cancel" }))
    expect(onClose).toHaveBeenCalledTimes(1)
  })

  it("should show a retry-friendly error and keep the input when saving fails", async () => {
    configureSafetyCheck("ready", [])
    const onSubmit = renderModal(vi.fn().mockResolvedValue({ status: "failed" }))
    fillPrescription()
    fireEvent.click(submitButton("modals.medication.confirm"))
    expect(await screen.findByTestId("form-submit-error")).toBeInTheDocument()
    expect(onSubmit).toHaveBeenCalledTimes(1)
    expect(screen.getByLabelText("modals.medication.name")).toHaveValue("Amoxicillin 500mg")
  })

  it("should not submit twice while the first submission is still running", async () => {
    configureSafetyCheck("ready", [])
    const onSubmit = renderModal(vi.fn().mockReturnValue(new Promise(() => undefined)))
    fillPrescription()
    const button = submitButton("modals.medication.confirm")
    fireEvent.click(button)
    await waitFor(() => expect(onSubmit).toHaveBeenCalledTimes(1))
    fireEvent.submit(button.closest("form") as HTMLFormElement)
    await new Promise((resolve) => setTimeout(resolve, 50))
    expect(onSubmit).toHaveBeenCalledTimes(1)
  })
})
