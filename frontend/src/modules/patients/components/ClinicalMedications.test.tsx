import { describe, it, expect, vi, beforeAll, beforeEach } from "vitest"
import { render, screen, fireEvent, waitFor } from "@testing-library/react"
import { QueryClient, QueryClientProvider } from "@tanstack/react-query"
import { AxiosError, AxiosHeaders, type AxiosResponse } from "axios"
import ClinicalMedications from "./ClinicalMedications"
import { TestAbilityProvider } from "../../../test/TestAbilityProvider"

vi.mock("react-i18next", async (importOriginal) => {
  const actual = await importOriginal<typeof import("react-i18next")>()
  return {
    ...actual,
    useTranslation: () => ({ t: (key: string) => key }),
  }
})

vi.mock("../api", () => ({
  patientsApi: {
    getMedications: vi.fn(),
    createMedication: vi.fn(),
    checkMedicationSafety: vi.fn(),
  },
}))

const toastMocks = vi.hoisted(() => ({ success: vi.fn(), error: vi.fn(), info: vi.fn() }))

vi.mock("../../../shared/store/toast_store", () => ({ toast: toastMocks }))

import { patientsApi } from "../api"

const allergyAlert = {
  kind: "allergy",
  severity: "contraindicated",
  blocking: true,
  messageKey: "allergy.activeIngredient",
  params: { allergen: "Penicillin", medication: "Amoxicillin" },
  relatedResourceId: "allergy-1",
}

const buildConflictError = () => {
  const response = {
    status: 409,
    data: { error: "blocked", code: "MEDICATION_ALLERGY_BLOCKED", alerts: [allergyAlert] },
    statusText: "",
    headers: {},
    config: { headers: new AxiosHeaders() },
  } as AxiosResponse
  return new AxiosError("conflict", "ERR_BAD_REQUEST", undefined, undefined, response)
}

const renderSection = () => {
  const queryClient = new QueryClient({ defaultOptions: { queries: { retry: false } } })
  render(
    <QueryClientProvider client={queryClient}>
      <TestAbilityProvider role="DOCTOR">
        <ClinicalMedications patientId="patient-1" encounterId="enc-1" />
      </TestAbilityProvider>
    </QueryClientProvider>,
  )
}

const openAndFill = async () => {
  fireEvent.click(await screen.findByRole("button", { name: "details.medicationsCard.add" }))
  fireEvent.change(await screen.findByLabelText("modals.medication.name"), { target: { value: "Amoxicillin 500mg" } })
  fireEvent.change(screen.getByLabelText("modals.medication.dosage"), { target: { value: "1 tablet every 8 hours" } })
  await waitFor(() => expect(screen.getByRole("button", { name: "modals.medication.confirm" })).toBeEnabled(), {
    timeout: 3000,
  })
}

class ResizeObserverStub {
  observe() {}
  unobserve() {}
  disconnect() {}
}

beforeAll(() => {
  vi.stubGlobal("ResizeObserver", ResizeObserverStub)
})

describe("ClinicalMedications prescription safety flow", () => {
  beforeEach(() => {
    vi.clearAllMocks()
    vi.mocked(patientsApi.getMedications).mockResolvedValue([])
    vi.mocked(patientsApi.checkMedicationSafety).mockResolvedValue({ alerts: [] })
  })

  it("should show the returned warnings after a successful prescription", async () => {
    vi.mocked(patientsApi.createMedication).mockResolvedValue({
      fhir_id: "med-1",
      warnings: [{ ...allergyAlert, kind: "interaction", severity: "moderate", blocking: false, messageKey: "interaction.moderate" }],
    } as never)
    renderSection()
    await openAndFill()
    fireEvent.click(screen.getByRole("button", { name: "modals.medication.confirm" }))
    await waitFor(() => expect(patientsApi.createMedication).toHaveBeenCalled())
    await waitFor(() => expect(toastMocks.info).toHaveBeenCalled())
    expect(toastMocks.success).toHaveBeenCalledWith("toast.medicationSuccess")
    expect(patientsApi.createMedication).toHaveBeenCalledWith({
      encounter_fhir_id: "enc-1",
      patient_fhir_id: "patient-1",
      medication_name: "Amoxicillin 500mg",
      dosage_instructions: "1 tablet every 8 hours",
      overrideJustification: undefined,
    })
  })

  it("should show the alerts of a 409 response and resubmit with the justification", async () => {
    vi.mocked(patientsApi.createMedication)
      .mockRejectedValueOnce(buildConflictError())
      .mockResolvedValueOnce({ fhir_id: "med-1" } as never)
    renderSection()
    await openAndFill()
    fireEvent.click(screen.getByRole("button", { name: "modals.medication.confirm" }))
    expect(await screen.findByTestId("safety-alert-allergy")).toBeInTheDocument()
    expect(toastMocks.error).toHaveBeenCalledWith("medicationSafety.blockedByServer")

    fireEvent.change(screen.getByLabelText("medicationSafety.override.justification"), {
      target: { value: "Benefit outweighs the risk" },
    })
    fireEvent.click(screen.getByRole("checkbox"))
    fireEvent.click(screen.getByRole("button", { name: "medicationSafety.override.confirm" }))
    await waitFor(() => expect(patientsApi.createMedication).toHaveBeenCalledTimes(2))
    expect(patientsApi.createMedication).toHaveBeenLastCalledWith(
      expect.objectContaining({ overrideJustification: "Benefit outweighs the risk" }),
    )
  })
})
