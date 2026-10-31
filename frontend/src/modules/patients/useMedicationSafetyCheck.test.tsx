import { describe, it, expect, vi, afterEach } from "vitest"
import { renderHook, waitFor } from "@testing-library/react"
import type { ReactNode } from "react"
import { QueryClient, QueryClientProvider } from "@tanstack/react-query"
import { useMedicationSafetyCheck } from "./queries"
import type { MedicationSafetyAlert } from "./types"

vi.mock("./api", () => ({
  patientsApi: {
    checkMedicationSafety: vi.fn(),
  },
}))

import { patientsApi } from "./api"

const allergyAlert: MedicationSafetyAlert = {
  kind: "allergy",
  severity: "major",
  blocking: true,
  messageKey: "allergy.activeIngredient",
  params: {},
  relatedResourceId: "allergy-1",
}

const renderSafetyHook = (medicationName: string) => {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false }, mutations: { retry: false } },
  })
  const wrapper = ({ children }: { children: ReactNode }) => (
    <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
  )
  return renderHook(() => useMedicationSafetyCheck("patient-1", medicationName), { wrapper })
}

describe("useMedicationSafetyCheck", () => {
  afterEach(() => {
    vi.clearAllMocks()
  })

  it("should stay idle for names shorter than the minimum", () => {
    const { result } = renderSafetyHook("am")
    expect(result.current.status).toBe("idle")
    expect(patientsApi.checkMedicationSafety).not.toHaveBeenCalled()
  })

  it("should call the check after the debounce and expose the alerts", async () => {
    vi.mocked(patientsApi.checkMedicationSafety).mockResolvedValue({ alerts: [allergyAlert] })
    const { result } = renderSafetyHook("Amoxicillin")
    expect(result.current.status).toBe("loading")
    await waitFor(() => expect(result.current.status).toBe("ready"), { timeout: 3000 })
    expect(patientsApi.checkMedicationSafety).toHaveBeenCalledWith("patient-1", {
      medicationCode: "",
      medicationName: "Amoxicillin",
    })
    expect(result.current.alerts).toEqual([allergyAlert])
  })

  it("should expose an error status when the check fails and report no alerts", async () => {
    vi.mocked(patientsApi.checkMedicationSafety).mockRejectedValue(new Error("down"))
    const { result } = renderSafetyHook("Amoxicillin")
    await waitFor(() => expect(result.current.status).toBe("error"), { timeout: 3000 })
    expect(result.current.alerts).toEqual([])
  })

  it("should report an unverified result from verifyNow when the request fails", async () => {
    vi.mocked(patientsApi.checkMedicationSafety).mockRejectedValue(new Error("down"))
    const { result } = renderSafetyHook("Amoxicillin")
    await expect(result.current.verifyNow()).resolves.toEqual({ verified: false })
  })

  it("should return the fresh alerts from verifyNow", async () => {
    vi.mocked(patientsApi.checkMedicationSafety).mockResolvedValue({ alerts: [allergyAlert] })
    const { result } = renderSafetyHook("Amoxicillin")
    await expect(result.current.verifyNow()).resolves.toEqual({ verified: true, alerts: [allergyAlert] })
  })
})
