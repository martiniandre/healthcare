import { describe, it, expect, vi, afterEach } from "vitest"
import { renderHook, waitFor } from "@testing-library/react"
import type { ReactNode } from "react"
import { QueryClient, QueryClientProvider } from "@tanstack/react-query"
import { useCreateVitalSignsPanelMutation } from "./queries"

vi.mock("./api", () => ({
  patientsApi: {
    createVitalSignsBatch: vi.fn(),
  },
}))

import { patientsApi } from "./api"

const panelPayload = {
  encounter_fhir_id: "enc-1",
  patient_fhir_id: "patient-1",
  panel_form_data: { heartRate: 72 },
}

const renderMutationHook = () => {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false }, mutations: { retry: false } },
  })
  queryClient.invalidateQueries = vi.fn()
  const wrapper = ({ children }: { children: ReactNode }) => (
    <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
  )
  return { ...renderHook(() => useCreateVitalSignsPanelMutation(), { wrapper }), queryClient }
}

describe("useCreateVitalSignsPanelMutation", () => {
  afterEach(() => {
    vi.clearAllMocks()
  })

  it("should call the batch endpoint and invalidate the observations list on success", async () => {
    vi.mocked(patientsApi.createVitalSignsBatch).mockResolvedValue([])
    const { result, queryClient } = renderMutationHook()

    await result.current.mutateAsync(panelPayload)

    await waitFor(() => {
      expect(queryClient.invalidateQueries).toHaveBeenCalledWith({
        queryKey: ["patients", "observations", "enc-1"],
      })
    })
    expect(patientsApi.createVitalSignsBatch).toHaveBeenCalledWith("enc-1", "patient-1", { heartRate: 72 })
  })

  it("should propagate the error and skip invalidation when the batch fails", async () => {
    vi.mocked(patientsApi.createVitalSignsBatch).mockRejectedValue(new Error("batch failed"))
    const { result, queryClient } = renderMutationHook()

    await expect(result.current.mutateAsync(panelPayload)).rejects.toThrow("batch failed")
    expect(queryClient.invalidateQueries).not.toHaveBeenCalled()
  })
})
