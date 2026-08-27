import { describe, it, expect, vi } from "vitest"
import { renderHook } from "@testing-library/react"
import { useVitalSignsColumns } from "./useVitalSignsColumns"

const mockTranslateFunction = (key: string) => key

vi.mock("react-i18next", () => ({
  useTranslation: () => ({
    t: mockTranslateFunction,
  }),
}))

describe("useVitalSignsColumns", () => {
  it("should return four column definitions", () => {
    const { result } = renderHook(() => useVitalSignsColumns())
    expect(result.current).toHaveLength(4)
  })

  it("should translate every column header", () => {
    const { result } = renderHook(() => useVitalSignsColumns())
    const headers = result.current.map((column) => column.header as string)
    expect(headers).toEqual([
      "details.vitalsCard.display",
      "details.vitalsCard.code",
      "details.vitalsCard.value",
      "details.vitalsCard.date",
    ])
  })

  it("should bind each column to its observation accessor", () => {
    const { result } = renderHook(() => useVitalSignsColumns())
    const accessorKeys = result.current.map((column) => column.id)
    expect(accessorKeys).toEqual(["code_display", "loinc_code", "value_quantity", "created_at"])
  })

  it("should keep the same array reference across rerenders with stable dependencies", () => {
    const { result, rerender } = renderHook(() => useVitalSignsColumns())
    const firstReference = result.current
    rerender()
    expect(result.current).toBe(firstReference)
  })
})
