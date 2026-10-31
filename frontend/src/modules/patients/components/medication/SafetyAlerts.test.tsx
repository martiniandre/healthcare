import { describe, it, expect, vi } from "vitest"
import { render, screen, fireEvent } from "@testing-library/react"
import { SafetyAlerts } from "./SafetyAlerts"
import type { MedicationSafetyAlert } from "../../types"

vi.mock("react-i18next", async (importOriginal) => {
  const actual = await importOriginal<typeof import("react-i18next")>()
  return {
    ...actual,
    useTranslation: () => ({
      t: (key: string, options?: Record<string, string>) =>
        options?.allergen ? `${key}:${options.allergen}` : key,
    }),
  }
})

const buildAlert = (overrides: Partial<MedicationSafetyAlert>): MedicationSafetyAlert => ({
  kind: "interaction",
  severity: "moderate",
  blocking: false,
  messageKey: "interaction.moderate",
  params: { medication: "Warfarin", drugClass: "NSAID" },
  relatedResourceId: "resource-1",
  ...overrides,
})

const renderPanel = (properties: Partial<React.ComponentProps<typeof SafetyAlerts>>) =>
  render(
    <SafetyAlerts
      status="ready"
      alerts={[]}
      onRetry={vi.fn()}
      unverifiedAcknowledged={false}
      onUnverifiedAcknowledgedChange={vi.fn()}
      {...properties}
    />,
  )

describe("SafetyAlerts", () => {
  it("should render nothing while idle", () => {
    const { container } = renderPanel({ status: "idle" })
    expect(container).toBeEmptyDOMElement()
  })

  it("should render a loading state while checking", () => {
    renderPanel({ status: "loading" })
    expect(screen.getByText("medicationSafety.checking")).toBeInTheDocument()
  })

  it("should render allergy alerts as destructive with a critical badge and localized message", () => {
    renderPanel({
      alerts: [
        buildAlert({
          kind: "allergy",
          severity: "contraindicated",
          blocking: true,
          messageKey: "allergy.activeIngredient",
          params: { allergen: "Penicillin", medication: "Amoxicillin" },
        }),
      ],
    })
    const alertElement = screen.getByTestId("safety-alert-allergy")
    expect(alertElement).toHaveTextContent("medicationSafety.messages.allergy.activeIngredient:Penicillin")
    expect(alertElement.querySelector("[data-status='critical']")).not.toBeNull()
    expect(alertElement.querySelector("svg")).not.toBeNull()
  })

  it("should render a moderate interaction with an attention badge", () => {
    renderPanel({ alerts: [buildAlert({ severity: "moderate" })] })
    const alertElement = screen.getByTestId("safety-alert-interaction")
    expect(alertElement.querySelector("[data-status='attention']")).not.toBeNull()
    expect(alertElement).toHaveTextContent("medicationSafety.severity.moderate")
  })

  it.each(["major", "contraindicated"] as const)("should render a %s interaction with a critical badge", (severity) => {
    renderPanel({ alerts: [buildAlert({ severity, messageKey: `interaction.${severity}` })] })
    const alertElement = screen.getByTestId("safety-alert-interaction")
    expect(alertElement.querySelector("[data-status='critical']")).not.toBeNull()
    expect(alertElement).toHaveTextContent(`medicationSafety.severity.${severity}`)
  })

  it("should render an info interaction with an info badge", () => {
    renderPanel({ alerts: [buildAlert({ severity: "info" })] })
    expect(screen.getByTestId("safety-alert-interaction").querySelector("[data-status='info']")).not.toBeNull()
  })

  it("should list allergy alerts before interaction alerts", () => {
    renderPanel({
      alerts: [
        buildAlert({ relatedResourceId: "interaction-1" }),
        buildAlert({ kind: "allergy", relatedResourceId: "allergy-1", messageKey: "allergy.sameClass" }),
      ],
    })
    const alertElements = screen.getAllByRole("alert")
    expect(alertElements[0]).toHaveAttribute("data-testid", "safety-alert-allergy")
  })

  it("should show the failure warning, retry and require the acknowledgement checkbox", () => {
    const onRetry = vi.fn()
    const onAcknowledgedChange = vi.fn()
    renderPanel({
      status: "error",
      onRetry,
      onUnverifiedAcknowledgedChange: onAcknowledgedChange,
      unverifiedAcknowledgementError: "required",
    })
    expect(screen.getByTestId("safety-check-failed")).toBeInTheDocument()
    expect(screen.getByText("required")).toBeInTheDocument()
    fireEvent.click(screen.getByRole("button", { name: "medicationSafety.checkFailed.retry" }))
    expect(onRetry).toHaveBeenCalledTimes(1)
    fireEvent.click(screen.getByRole("checkbox"))
    expect(onAcknowledgedChange).toHaveBeenCalledWith(true)
  })
})
