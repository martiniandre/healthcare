import { describe, it, expect, vi } from "vitest"
import { renderHook } from "@testing-library/react"
import { useClinicalMedicationsColumns } from "./useClinicalMedicationsColumns"

const mockTranslateFunction = (key: string) => key

vi.mock("react-i18next", () => ({
  useTranslation: () => ({
    t: mockTranslateFunction,
  }),
}))

describe("useClinicalMedicationsColumns", () => {
  it("should return four column definitions", () => {
    const { result } = renderHook(() => useClinicalMedicationsColumns())
    expect(result.current).toHaveLength(4)
  })

  it("should translate every column header", () => {
    const { result } = renderHook(() => useClinicalMedicationsColumns())
    const headers = result.current.map((column) => column.header as string)
    expect(headers).toEqual([
      "details.medicationsCard.display",
      "details.medicationsCard.dosage",
      "details.medicationsCard.status",
      "details.medicationsCard.date",
    ])
  })

  it("should bind each column to its medication accessor", () => {
    const { result } = renderHook(() => useClinicalMedicationsColumns())
    const accessorKeys = result.current.map((column) => column.id)
    expect(accessorKeys).toEqual(["medication_name", "dosage_instructions", "status", "created_at"])
  })

  it("should keep the same array reference across rerenders with stable dependencies", () => {
    const { result, rerender } = renderHook(() => useClinicalMedicationsColumns())
    const firstReference = result.current
    rerender()
    expect(result.current).toBe(firstReference)
  })
})
